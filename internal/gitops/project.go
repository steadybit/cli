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
	"github.com/steadybit/cli/v6/internal/environment"
	"github.com/steadybit/cli/v6/internal/experiment"
	"github.com/steadybit/cli/v6/internal/hub"
	"github.com/steadybit/cli/v6/internal/integration"
	"github.com/steadybit/cli/v6/internal/jsyaml"
	"github.com/steadybit/cli/v6/internal/output"
	"github.com/steadybit/cli/v6/internal/platform"
	"github.com/steadybit/cli/v6/internal/property"
	"github.com/steadybit/cli/v6/internal/schedule"
	"github.com/steadybit/cli/v6/internal/service"
	"github.com/steadybit/cli/v6/internal/serviceprofile"
	"github.com/steadybit/cli/v6/internal/team"
	"github.com/steadybit/cli/v6/internal/template"
)

// A project is a directory holding what a team or the tenant keeps in Git, one kind per
// directory, in the order applying them has to follow. Templates and experiments carry
// the properties definitions define; teams name their environments, integrations their
// team, service profiles their templates, services their profile, experiments their
// team and environment, schedules their experiment. Hubs depend on nothing.
var projectKinds = []struct {
	dir   string
	kind  Kind
	apply func(ctx context.Context, c *platform.Client, path string, o ApplyOptions) error
}{
	{"property-definitions", PropertyDefinition, func(ctx context.Context, c *platform.Client, path string, _ ApplyOptions) error {
		return property.ApplyDefinitions(ctx, c, property.ApplyDefinitionOptions{Files: []string{path}, Recursive: true})
	}},
	{"environments", Environment, func(ctx context.Context, c *platform.Client, path string, _ ApplyOptions) error {
		return environment.Apply(ctx, c, environment.ApplyOptions{Files: []string{path}, Recursive: true})
	}},
	{"teams", Team, func(ctx context.Context, c *platform.Client, path string, _ ApplyOptions) error {
		return team.Apply(ctx, c, team.ApplyOptions{Files: []string{path}, Recursive: true})
	}},
	// Synchronized, as `hub apply --synchronize` does: service profiles name the templates
	// a hub brings, so on a restored tenant those have to be there first. A hub that cannot
	// be synchronized ends the apply.
	{"hubs", Hub, func(ctx context.Context, c *platform.Client, path string, _ ApplyOptions) error {
		return hub.Apply(ctx, c, hub.ApplyOptions{Files: []string{path}, Recursive: true, Synchronize: true})
	}},
	{"templates", Template, func(ctx context.Context, c *platform.Client, path string, _ ApplyOptions) error {
		return template.Apply(ctx, c, template.ApplyOptions{Files: []string{path}, Recursive: true})
	}},
	{"integrations/webhook", Integrations[integration.Webhook.Name], applyIntegrations(integration.Webhook)},
	{"integrations/slack", Integrations[integration.Slack.Name], applyIntegrations(integration.Slack)},
	{"integrations/preflight", Integrations[integration.Preflight.Name], applyIntegrations(integration.Preflight)},
	{"integrations/preflight-action", Integrations[integration.PreflightAction.Name], applyIntegrations(integration.PreflightAction)},
	{"service-profiles", ServiceProfile, func(ctx context.Context, c *platform.Client, path string, o ApplyOptions) error {
		return serviceprofile.Apply(ctx, c, serviceprofile.ApplyOptions{Files: []string{path}, Recursive: true, DeleteExperiments: o.DeleteExperiments})
	}},
	{"services", Service, func(ctx context.Context, c *platform.Client, path string, o ApplyOptions) error {
		return service.Apply(ctx, c, service.ApplyOptions{Files: []string{path}, Recursive: true, DeleteExperiments: o.DeleteExperiments})
	}},
	{"experiments", Experiment, func(ctx context.Context, c *platform.Client, path string, _ ApplyOptions) error {
		return experiment.Apply(ctx, c, experiment.ApplyOptions{Files: []string{path}, Recursive: true})
	}},
	{"schedules", Schedule, func(ctx context.Context, c *platform.Client, path string, _ ApplyOptions) error {
		return schedule.Apply(ctx, c, schedule.ApplyOptions{Files: []string{path}, Recursive: true})
	}},
}

// applyIntegrations leaves out the files that match the platform. Those holding a masked
// secret could not be applied, the platform keeping no secret it is not sent, and the
// others need not be.
func applyIntegrations(k integration.Kind) func(ctx context.Context, c *platform.Client, path string, o ApplyOptions) error {
	return func(ctx context.Context, c *platform.Client, path string, _ ApplyOptions) error {
		gk := Integrations[k.Name]
		results, err := compareAll(ctx, c, gk, []string{path}, true)
		if err != nil {
			return err
		}
		var changed []string
		for _, r := range results {
			if r.State == Unchanged {
				fmt.Println(r.Describe(gk))
			} else {
				changed = append(changed, r.File)
			}
		}
		if len(changed) == 0 {
			return nil
		}
		return integration.Apply(ctx, c, k, integration.ApplyOptions{Files: changed})
	}
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
	Tenant    bool
}

// Export writes a team's experiments, schedules and services, and the custom service
// profiles those services use, as a project; or the tenant's configuration. Files are
// only written, never removed, so something deleted on the platform keeps its file until
// it is removed by hand.
func Export(ctx context.Context, c *platform.Client, o ExportOptions) error {
	if o.Tenant {
		return exportTenant(ctx, c, o)
	}
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
	var names []string
	for _, pk := range projectKinds {
		names = append(names, pk.dir+"/")
		path := filepath.Join(dir, filepath.FromSlash(pk.dir))
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			found[pk.dir] = path
		}
	}
	if len(found) == 0 {
		return nil, fmt.Errorf("'%s' holds none of %s or %s.", dir, strings.Join(names[:len(names)-1], ", "), names[len(names)-1])
	}
	return found, nil
}

type ApplyOptions struct {
	Directory         string
	DeleteExperiments bool
	DryRun            bool
}

// ApplyProject applies every kind in a project, in the order of projectKinds.
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
			err = pk.apply(ctx, c, path, o)
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
