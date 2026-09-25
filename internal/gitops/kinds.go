// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package gitops

import (
	"context"
	"fmt"
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/steadybit/cli/api"
	"github.com/steadybit/cli/internal/jsyaml"
	"github.com/steadybit/cli/internal/output"
	"github.com/steadybit/cli/internal/platform"
)

// Kind is a type of file kept in Git and how to find its counterpart on the platform.
type Kind struct {
	// Name is what messages call it, "experiment" or "service profile".
	Name string
	// ReadOnly fields are reported by the platform but not part of the file.
	ReadOnly []string
	// Identity is the field that names it on the platform. A file matched otherwise, by
	// an externalId or a name, has none yet, which is not a difference.
	Identity string
	// Remote finds the platform's version of a file: its id and content, or nil when
	// applying the file would create something new.
	Remote func(ctx context.Context, c *platform.Client, local *jsyaml.Map) (string, *jsyaml.Map, error)
}

func str(m *jsyaml.Map, key string) string {
	v, _ := m.Get(key)
	s, _ := v.(string)
	return s
}

// fetch reads one document; a 404 means there is nothing yet.
func fetch(resp *http.Response, err error) (*jsyaml.Map, error) {
	doc, _, err := platform.ReadDocument(resp, err)
	if platform.IsStatus(err, http.StatusNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return doc.Value(), nil
}

func uuid(id string) (openapi_types.UUID, bool) {
	var u openapi_types.UUID
	return u, u.UnmarshalText([]byte(id)) == nil
}

var Experiment = Kind{
	Name:     "experiment",
	ReadOnly: []string{"version", "created", "createdBy", "edited", "editedBy"},
	Identity: "key",
	// By its key, or, for a file that has none yet, by the externalId an apply would
	// match it with.
	Remote: func(ctx context.Context, c *platform.Client, local *jsyaml.Map) (string, *jsyaml.Map, error) {
		key := str(local, "key")
		if key == "" {
			externalID := str(local, "externalId")
			if externalID == "" {
				return "", nil, nil
			}
			var list struct {
				Experiments []struct {
					Key string `json:"key"`
				} `json:"experiments"`
			}
			ids := []string{externalID}
			resp, err := c.GetExperiments(ctx, &api.GetExperimentsParams{ExternalId: &ids})
			if _, err := platform.Decode(resp, err, &list); err != nil {
				return "", nil, err
			}
			if len(list.Experiments) == 0 {
				return "", nil, nil
			}
			key = list.Experiments[0].Key
		}
		remote, err := fetch(c.GetExperiment(ctx, key))
		return key, remote, err
	},
}

var Schedule = Kind{
	Name:     "experiment schedule",
	ReadOnly: []string{"editedBy", "lastUpdated", "nextExecution"},
	Identity: "id",
	Remote: func(ctx context.Context, c *platform.Client, local *jsyaml.Map) (string, *jsyaml.Map, error) {
		id := str(local, "id")
		if id == "" {
			return "", nil, nil
		}
		remote, err := fetch(c.GetSchedules(ctx, id))
		return id, remote, err
	},
}

var Service = Kind{
	Name:     "service",
	ReadOnly: []string{"version", "created", "createdBy", "edited", "editedBy"},
	Identity: "id",
	Remote: func(ctx context.Context, c *platform.Client, local *jsyaml.Map) (string, *jsyaml.Map, error) {
		id := str(local, "id")
		u, ok := uuid(id)
		if !ok {
			return "", nil, nil
		}
		remote, err := fetch(c.GetService(ctx, u))
		return id, remote, err
	},
}

var ServiceProfile = Kind{
	Name:     "service profile",
	ReadOnly: []string{"version", "created", "createdBy", "edited", "editedBy", "defaultProfile"},
	Identity: "id",
	// By its id, or by its name, which is unique among profiles.
	Remote: func(ctx context.Context, c *platform.Client, local *jsyaml.Map) (string, *jsyaml.Map, error) {
		id := str(local, "id")
		if u, ok := uuid(id); ok {
			remote, err := fetch(c.GetProfile(ctx, u))
			return id, remote, err
		}
		name := str(local, "name")
		if name == "" {
			return "", nil, nil
		}
		type profile struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}
		profiles, err := platform.AllPages[profile](func(page, size int32) (*http.Response, error) {
			return c.GetProfiles(ctx, &api.GetProfilesParams{Name: &name, Page: api.PageRequestAO{Page: &page, Size: &size}})
		})
		if err != nil {
			return "", nil, err
		}
		for _, p := range profiles {
			if p.Name == name {
				u, _ := uuid(p.ID)
				remote, err := fetch(c.GetProfile(ctx, u))
				return p.ID, remote, err
			}
		}
		return "", nil, nil
	},
}

func strip(m *jsyaml.Map, fields []string) *jsyaml.Map {
	c := jsyaml.Clone(m).(*jsyaml.Map)
	for _, f := range fields {
		c.Delete(f)
	}
	return c
}

// State is what applying a file would do.
type State int

const (
	Unchanged State = iota
	Changed
	New
)

type Result struct {
	File  string
	ID    string
	State State
	Diff  string
}

// Compare works out what applying the file would change on the platform.
func Compare(ctx context.Context, c *platform.Client, k Kind, file string, local *output.Document) (Result, error) {
	id, remote, err := k.Remote(ctx, c, local.Value())
	if err != nil {
		return Result{}, platform.Failed(err, "Failed to get the %s for %s", k.Name, file)
	}
	if remote == nil {
		return Result{File: file, ID: id, State: New}, nil
	}
	ignored := k.ReadOnly
	if _, has := local.Value().Get(k.Identity); !has {
		ignored = append(append([]string{}, ignored...), k.Identity)
	}
	diff, err := Diff(file, strip(local.Value(), ignored), strip(remote, ignored))
	if err != nil {
		return Result{}, err
	}
	if diff == "" {
		return Result{File: file, ID: id, State: Unchanged}, nil
	}
	return Result{File: file, ID: id, State: Changed, Diff: diff}, nil
}

// Describe is the one-line outcome of a comparison, as a dry run reports it.
func (r Result) Describe(k Kind) string {
	switch r.State {
	case New:
		return fmt.Sprintf("%s would create a new %s.", r.File, k.Name)
	case Changed:
		return fmt.Sprintf("%s would update %s %s (%d lines changed).", r.File, k.Name, r.ID, Changes(r.Diff))
	}
	return fmt.Sprintf("%s matches %s %s.", r.File, k.Name, r.ID)
}
