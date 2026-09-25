// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package gitops

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/steadybit/cli/v6/api"
	"github.com/steadybit/cli/v6/internal/experiment"
	"github.com/steadybit/cli/v6/internal/jsyaml"
	"github.com/steadybit/cli/v6/internal/output"
	"github.com/steadybit/cli/v6/internal/platform"
	"github.com/steadybit/cli/v6/internal/schedule"
	"github.com/steadybit/cli/v6/internal/service"
	"github.com/steadybit/cli/v6/internal/serviceprofile"
)

// A project is a directory holding what a team keeps in Git, one kind per directory,
// in the order applying them has to follow: services name their profile, schedules
// their experiment.
var projectKinds = []struct {
	dir  string
	kind Kind
}{
	{"service-profiles", ServiceProfile},
	{"services", Service},
	{"experiments", Experiment},
	{"schedules", Schedule},
}

var unsafe = regexp.MustCompile(`[^a-z0-9._-]+`)

// fileName makes a name safe and readable as a file name, and unique within dir.
func fileName(dir, name string, taken map[string]bool) string {
	base := strings.Trim(unsafe.ReplaceAllString(strings.ToLower(name), "-"), "-.")
	if base == "" {
		base = "unnamed"
	}
	candidate := base
	for i := 2; taken[filepath.Join(dir, candidate)]; i++ {
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
	taken[filepath.Join(dir, candidate)] = true
	return filepath.Join(dir, candidate+".yaml")
}

func writeDocument(file string, doc *jsyaml.Map, readOnly []string) error {
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		return err
	}
	return os.WriteFile(file, []byte(jsyaml.Dump(strip(doc, readOnly))), 0o644)
}

type ExportOptions struct {
	Directory string
	Team      string
}

// Export writes a team's experiments, schedules and services, and the custom service
// profiles those services use, as a project. Files are only written, never removed, so
// something deleted on the platform keeps its file until it is removed by hand.
func Export(ctx context.Context, c *platform.Client, o ExportOptions) error {
	taken := map[string]bool{}
	counts := map[string]int{}
	team := []string{o.Team}

	var experiments struct {
		Experiments []struct {
			Key string `json:"key"`
		} `json:"experiments"`
	}
	resp, err := c.GetExperiments(ctx, &api.GetExperimentsParams{Team: &team})
	if _, err := platform.Decode(resp, err, &experiments); err != nil {
		return platform.Failed(err, "Failed to get the experiments of team %s", o.Team)
	}
	for _, e := range experiments.Experiments {
		doc, _, err := platform.ReadDocument(c.GetExperiment(ctx, e.Key))
		if err != nil {
			return platform.Failed(err, "Failed to get experiment %s", e.Key)
		}
		if err := writeDocument(fileName(filepath.Join(o.Directory, "experiments"), e.Key, taken), doc.Value(), Experiment.ReadOnly); err != nil {
			return err
		}
		counts["experiments"]++
	}

	var schedules []json.RawMessage
	resp, err = c.GetAllSchedulesV2(ctx, &api.GetAllSchedulesV2Params{Team: &team})
	if _, err := platform.Decode(resp, err, &schedules); err != nil {
		return platform.Failed(err, "Failed to get the experiment schedules of team %s", o.Team)
	}
	for _, raw := range schedules {
		doc, err := output.ParseDocument(raw)
		if err != nil {
			return err
		}
		key, id := str(doc.Value(), "experimentKey"), str(doc.Value(), "id")
		if len(id) > 8 {
			id = id[len(id)-8:]
		}
		if err := writeDocument(fileName(filepath.Join(o.Directory, "schedules"), key+"-"+id, taken), doc.Value(), Schedule.ReadOnly); err != nil {
			return err
		}
		counts["schedules"]++
	}

	type summary struct {
		ID string `json:"id"`
	}
	services, err := platform.AllPages[summary](func(page, size int32) (*http.Response, error) {
		return c.GetServiceList(ctx, &api.GetServiceListParams{TeamKey: &team, Page: api.PageRequestAO{Page: &page, Size: &size}})
	})
	if err != nil {
		return platform.Failed(err, "Failed to get the services of team %s", o.Team)
	}
	profiles := map[string]bool{}
	for _, s := range services {
		u, _ := uuid(s.ID)
		doc, _, err := platform.ReadDocument(c.GetService(ctx, u))
		if err != nil {
			return platform.Failed(err, "Failed to get service %s", s.ID)
		}
		profiles[str(doc.Value(), "serviceProfile")] = true
		if err := writeDocument(fileName(filepath.Join(o.Directory, "services"), str(doc.Value(), "name"), taken), doc.Value(), Service.ReadOnly); err != nil {
			return err
		}
		counts["services"]++
	}

	// Profiles are shared by all teams; only the custom ones this team's services use
	// belong to its project. Steadybit's own come with the platform.
	for name := range profiles {
		if name == "" {
			continue
		}
		_, profile, err := ServiceProfile.Remote(ctx, c, mapWith("name", name))
		if err != nil {
			return platform.Failed(err, "Failed to get service profile %s", name)
		}
		if profile == nil || str(profile, "origin") != "CUSTOM" {
			continue
		}
		if err := writeDocument(fileName(filepath.Join(o.Directory, "service-profiles"), name, taken), profile, ServiceProfile.ReadOnly); err != nil {
			return err
		}
		counts["service-profiles"]++
	}

	fmt.Printf("Exported team %s to %s: %d experiments, %d schedules, %d services, %d service profiles.\n",
		o.Team, o.Directory, counts["experiments"], counts["schedules"], counts["services"], counts["service-profiles"])
	return nil
}

func mapWith(key, value string) *jsyaml.Map {
	m := jsyaml.NewMap()
	m.Set(key, value)
	return m
}

// projectDirs are the kinds a project directory holds, in the order to apply them.
func projectDirs(dir string) (map[string]string, error) {
	found := map[string]string{}
	for _, pk := range projectKinds {
		path := filepath.Join(dir, pk.dir)
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			found[pk.dir] = path
		}
	}
	if len(found) == 0 {
		return nil, fmt.Errorf("'%s' holds none of experiments/, schedules/, services/ or service-profiles/.", dir)
	}
	return found, nil
}

type ApplyOptions struct {
	Directory         string
	DeleteExperiments bool
	DryRun            bool
}

// ApplyProject applies every kind in a project, profiles first and schedules last.
func ApplyProject(ctx context.Context, c *platform.Client, o ApplyOptions) error {
	dirs, err := projectDirs(o.Directory)
	if err != nil {
		return err
	}
	for _, pk := range projectKinds {
		path, ok := dirs[pk.dir]
		if !ok {
			continue
		}
		if o.DryRun {
			err = DryRun(ctx, c, pk.kind, []string{path}, true)
		} else {
			switch pk.dir {
			case "service-profiles":
				err = serviceprofile.Apply(ctx, c, serviceprofile.ApplyOptions{Files: []string{path}, Recursive: true, DeleteExperiments: o.DeleteExperiments})
			case "services":
				err = service.Apply(ctx, c, service.ApplyOptions{Files: []string{path}, Recursive: true, DeleteExperiments: o.DeleteExperiments})
			case "experiments":
				err = experiment.Apply(ctx, c, experiment.ApplyOptions{Files: []string{path}, Recursive: true})
			case "schedules":
				err = schedule.Apply(ctx, c, schedule.ApplyOptions{Files: []string{path}, Recursive: true})
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// DiffProject diffs every kind in a project, and reports drift once all were compared.
func DiffProject(ctx context.Context, c *platform.Client, dir string) error {
	dirs, err := projectDirs(dir)
	if err != nil {
		return err
	}
	drift := false
	for _, pk := range projectKinds {
		if path, ok := dirs[pk.dir]; ok {
			err := DiffFiles(ctx, c, pk.kind, []string{path}, true)
			if errors.Is(err, ErrDifferent) {
				drift = true
			} else if err != nil {
				return err
			}
		}
	}
	if drift {
		return ErrDifferent
	}
	return nil
}
