// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package gitops

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/steadybit/cli/v6/api"
	"github.com/steadybit/cli/v6/internal/environment"
	"github.com/steadybit/cli/v6/internal/hub"
	"github.com/steadybit/cli/v6/internal/integration"
	"github.com/steadybit/cli/v6/internal/jsyaml"
	"github.com/steadybit/cli/v6/internal/output"
	"github.com/steadybit/cli/v6/internal/platform"
	"github.com/steadybit/cli/v6/internal/property"
	"github.com/steadybit/cli/v6/internal/team"
	"github.com/steadybit/cli/v6/internal/template"
)

// Kind is a type of file kept in Git and how to find its counterpart on the platform.
type Kind struct {
	// Name is what messages call it, "experiment" or "service profile".
	Name string
	// ReadOnly fields are reported by the platform but not part of the file. In
	// "members.name", the field is dropped from each member.
	ReadOnly []string
	// Defaults are the values the platform gives fields a file leaves out, or leaves as
	// an empty list or map. Holding exactly that, such a field is not a difference;
	// holding anything else, applying the file would reset it, which is.
	Defaults map[string]any
	// Secrets are fields holding credentials, as integration.Kind names them. A mask in
	// a file stands for whatever the platform holds, and no value is ever printed.
	Secrets []string
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
	Defaults: map[string]any{"allowParallel": true},
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
	Defaults: map[string]any{"logoId": "service", "logoColor": "blue"},
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
	Defaults: map[string]any{"origin": "CUSTOM"},
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

// byID finds the platform's version of a file by the id in it. A file without one is new;
// one whose id is no UUID names nothing on the platform, and applying it would not work.
func byID(get func(ctx context.Context, c *platform.Client, id openapi_types.UUID) (*http.Response, error)) func(context.Context, *platform.Client, *jsyaml.Map) (string, *jsyaml.Map, error) {
	return func(ctx context.Context, c *platform.Client, local *jsyaml.Map) (string, *jsyaml.Map, error) {
		id := str(local, "id")
		if id == "" {
			return "", nil, nil
		}
		u, ok := uuid(id)
		if !ok {
			return "", nil, fmt.Errorf("%s is not a valid id", id)
		}
		remote, err := fetch(get(ctx, c, u))
		return id, remote, err
	}
}

var Template = Kind{
	Name:     "experiment template",
	ReadOnly: template.ReadOnly,
	Identity: "id",
	Remote: byID(func(ctx context.Context, c *platform.Client, id openapi_types.UUID) (*http.Response, error) {
		return c.GetExperimentTemplate(ctx, id)
	}),
}

var Environment = Kind{
	Name:     "environment",
	ReadOnly: environment.ReadOnly,
	Identity: "id",
	// By its id, or by its name, which the platform matches a file without an id by.
	Remote: func(ctx context.Context, c *platform.Client, local *jsyaml.Map) (string, *jsyaml.Map, error) {
		if id := str(local, "id"); id != "" {
			return byID(func(ctx context.Context, c *platform.Client, id openapi_types.UUID) (*http.Response, error) {
				return c.GetEnvironment(ctx, id)
			})(ctx, c, local)
		}
		name := str(local, "name")
		if name == "" {
			return "", nil, nil
		}
		var list struct {
			Environments []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"environments"`
		}
		resp, err := c.GetEnvironments(ctx, &api.GetEnvironmentsParams{Search: &name})
		if _, err := platform.Decode(resp, err, &list); err != nil {
			return "", nil, err
		}
		for _, e := range list.Environments {
			if e.Name == name {
				u, _ := uuid(e.ID)
				remote, err := fetch(c.GetEnvironment(ctx, u))
				return e.ID, remote, err
			}
		}
		return "", nil, nil
	},
}

// Team is matched by its key; the id `get` writes is not needed, so a file without one
// is not a difference.
var Team = Kind{
	Name:     "team",
	ReadOnly: append(append([]string{}, team.ReadOnly...), prefixed("members.", team.MemberReadOnly)...),
	// Sent none, a team may still wait and validate services.
	Defaults: map[string]any{"allowedActions": []any{"wait", "service-validation"}, "managedBy": "MANUAL"},
	Identity: "id",
	Remote: func(ctx context.Context, c *platform.Client, local *jsyaml.Map) (string, *jsyaml.Map, error) {
		key := str(local, "key")
		if key == "" {
			return "", nil, nil
		}
		remote, err := fetch(c.GetTeam(ctx, key))
		return key, remote, err
	},
}

var PropertyDefinition = Kind{
	Name:     "property definition",
	ReadOnly: property.ReadOnly,
	Identity: "key",
	Remote: func(ctx context.Context, c *platform.Client, local *jsyaml.Map) (string, *jsyaml.Map, error) {
		key := str(local, "key")
		if key == "" {
			return "", nil, nil
		}
		remote, err := fetch(c.GetPropertyDefinition(ctx, key))
		return key, remote, err
	},
}

var Hub = Kind{
	Name:     "hub",
	ReadOnly: hub.ReadOnly,
	Identity: "id",
	Remote: byID(func(ctx context.Context, c *platform.Client, id openapi_types.UUID) (*http.Response, error) {
		return c.GetHubById(ctx, id)
	}),
}

// Integrations are the kinds of integration by their command's name.
var Integrations = map[string]Kind{
	integration.Webhook.Name:         integrationKind(integration.Webhook, "webhook integration", allTargetAttributes),
	integration.Slack.Name:           integrationKind(integration.Slack, "Slack integration", nil),
	integration.Preflight.Name:       integrationKind(integration.Preflight, "preflight webhook", allTargetAttributes),
	integration.PreflightAction.Name: integrationKind(integration.PreflightAction, "preflight action integration", nil),
}

// Sent none, a webhook reports every target attribute.
var allTargetAttributes = map[string]any{"targetAttributeIncludes": []any{"*"}}

func integrationKind(k integration.Kind, name string, defaults map[string]any) Kind {
	return Kind{
		Name:     name,
		ReadOnly: integration.ReadOnly,
		Defaults: defaults,
		Secrets:  k.Secrets,
		Identity: "id",
		Remote:   byID(k.Fetch),
	}
}

func prefixed(prefix string, fields []string) []string {
	out := make([]string, len(fields))
	for i, f := range fields {
		out[i] = prefix + f
	}
	return out
}

func strip(m *jsyaml.Map, fields []string) *jsyaml.Map {
	c := jsyaml.Clone(m).(*jsyaml.Map)
	for _, f := range fields {
		deletePath(c, strings.Split(f, "."))
	}
	return c
}

// deletePath removes a field, going into every item of the lists on its way.
func deletePath(v any, path []string) {
	switch x := v.(type) {
	case *jsyaml.Map:
		if len(path) == 1 {
			x.Delete(path[0])
		} else if child, ok := x.Get(path[0]); ok {
			deletePath(child, path[1:])
		}
	case []any:
		for _, item := range x {
			deletePath(item, path)
		}
	}
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
	// A field with a default is compared against that default: false is a value to
	// report when the default is true.
	var defaulted []string
	for key := range k.Defaults {
		defaulted = append(defaulted, key)
	}
	l, r := strip(local.Value(), ignored), withoutDefaults(local.Value(), strip(remote, ignored), k.Defaults)
	hideSecrets(l, r, k.Secrets)
	diff, err := Diff(file, l, r, defaulted...)
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

// withoutDefaults drops the fields the file leaves out that hold the platform's default,
// and takes the file's empty list or map for the default the platform turned it into.
func withoutDefaults(local, remote *jsyaml.Map, defaults map[string]any) *jsyaml.Map {
	out := jsyaml.Clone(remote).(*jsyaml.Map)
	for key, value := range defaults {
		held, ok := out.Get(key)
		if !ok || jsyaml.CompactJSON(held) != jsyaml.CompactJSON(value) {
			continue
		}
		switch set, has := local.Get(key); {
		case !has:
			out.Delete(key)
		case isEmptyCollection(set):
			out.Set(key, set)
		}
	}
	return out
}

func isEmptyCollection(v any) bool {
	switch x := v.(type) {
	case []any:
		return len(x) == 0
	case *jsyaml.Map:
		return x.Len() == 0
	}
	return false
}

// hideSecrets makes credentials comparable without printing them. The platform's value
// shows as the mask, which a mask in the file matches, as does the value itself; any
// other value in the file is what applying it would set, and shows as that.
func hideSecrets(local, remote *jsyaml.Map, secrets []string) {
	for _, field := range secrets {
		lv, _ := local.Get(field)
		rv, _ := remote.Get(field)
		lm, lIsMap := lv.(*jsyaml.Map)
		rm, rIsMap := rv.(*jsyaml.Map)
		if lIsMap || rIsMap {
			if lm == nil {
				lm = jsyaml.NewMap()
			}
			if rm == nil {
				rm = jsyaml.NewMap()
			}
			for _, key := range lm.Keys() {
				hideSecret(lm, rm, key)
			}
			for _, key := range rm.Keys() {
				hideSecret(lm, rm, key)
			}
			continue
		}
		hideSecret(local, remote, field)
	}
}

// hideSecret masks any value a secret field holds, whatever its type: an unquoted number
// in a file is as much a credential as a string.
func hideSecret(local, remote *jsyaml.Map, key string) {
	lv, lok := local.Get(key)
	rv, rok := remote.Get(key)
	if rok && !blank(rv) {
		remote.Set(key, integration.Mask)
	}
	switch {
	case !lok || blank(lv) || lv == integration.Mask:
	case integration.Masked(lv) || (rok && jsyaml.CompactJSON(lv) == jsyaml.CompactJSON(rv)):
		local.Set(key, integration.Mask)
	default:
		local.Set(key, integration.Mask+" (from the file)")
	}
}

func blank(v any) bool {
	return v == nil || v == ""
}
