// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package gitops

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/steadybit/cli/v6/api"
	"github.com/steadybit/cli/v6/internal/integration"
	"github.com/steadybit/cli/v6/internal/jsyaml"
	"github.com/steadybit/cli/v6/internal/platform"
)

// The hubs every tenant is connected to by the platform itself, with the same ids
// everywhere. Kept in Git, they would be connected a second time on another platform.
var providedHubs = map[string]bool{
	"6a55640d-72d4-4194-a058-0afcf731dfac": true, // Steadybit Reliability Hub
	"738c90be-bdc4-4e0b-a4de-c8f1144b25c7": true, // Steadybit Service Templates
}

// tenantExport writes one kind of the tenant's configuration.
type tenantExport struct {
	o      ExportOptions
	taken  map[string]bool
	counts map[string]int
}

func (e *tenantExport) write(dir, name string, doc *jsyaml.Map, k Kind) error {
	e.counts[dir]++
	return writeDocument(fileName(filepath.Join(e.o.Directory, dir), name, e.taken), doc, k.ReadOnly)
}

func (e *tenantExport) each(ids []string, what string, get func(id string) (*http.Response, error), write func(doc *jsyaml.Map) error) error {
	for _, id := range ids {
		doc, _, err := platform.ReadDocument(get(id))
		if err != nil {
			return platform.Failed(err, "Failed to get %s %s", what, id)
		}
		if err := write(doc.Value()); err != nil {
			return err
		}
	}
	return nil
}

func byUUID(get func(openapi_types.UUID) (*http.Response, error)) func(string) (*http.Response, error) {
	return func(id string) (*http.Response, error) {
		u, ok := uuid(id)
		if !ok {
			return nil, errors.New("not a valid id")
		}
		return get(u)
	}
}

// exportTenant writes what the tenant's admins configure. What the platform provides
// is left out, as it comes with every platform: the hubs it connects, its service
// profiles, and the templates imported from a hub, which importing again brings back.
func exportTenant(ctx context.Context, c *platform.Client, o ExportOptions) error {
	e := &tenantExport{o: o, taken: map[string]bool{}, counts: map[string]int{}}

	type keyed struct {
		Key string `json:"key"`
	}
	definitions, err := platform.AllPages[keyed](func(page, size int32) (*http.Response, error) {
		return c.GetPropertyDefinitions(ctx, &api.GetPropertyDefinitionsParams{Page: api.PageRequestAO{Page: &page, Size: &size}})
	})
	if err != nil {
		return platform.Failed(err, "Failed to get the property definitions")
	}
	var keys []string
	for _, d := range definitions {
		keys = append(keys, d.Key)
	}
	if err := e.each(keys, "property definition", func(key string) (*http.Response, error) { return c.GetPropertyDefinition(ctx, key) },
		func(doc *jsyaml.Map) error {
			return e.write("property-definitions", str(doc, "key"), doc, PropertyDefinition)
		}); err != nil {
		return err
	}

	var environments struct {
		Environments []struct {
			ID string `json:"id"`
		} `json:"environments"`
	}
	resp, err := c.GetEnvironments(ctx, &api.GetEnvironmentsParams{})
	if _, err := platform.Decode(resp, err, &environments); err != nil {
		return platform.Failed(err, "Failed to get the environments")
	}
	var ids []string
	for _, env := range environments.Environments {
		ids = append(ids, env.ID)
	}
	if err := e.each(ids, "environment", byUUID(func(id openapi_types.UUID) (*http.Response, error) { return c.GetEnvironment(ctx, id) }),
		func(doc *jsyaml.Map) error { return e.write("environments", str(doc, "name"), doc, Environment) }); err != nil {
		return err
	}

	var teams struct {
		Teams []keyed `json:"teams"`
	}
	resp, err = c.GetTeams(ctx, &api.GetTeamsParams{})
	if _, err := platform.Decode(resp, err, &teams); err != nil {
		return platform.Failed(err, "Failed to get the teams")
	}
	keys = nil
	for _, t := range teams.Teams {
		keys = append(keys, t.Key)
	}
	if err := e.each(keys, "team", func(key string) (*http.Response, error) { return c.GetTeam(ctx, key) },
		func(doc *jsyaml.Map) error { return e.write("teams", str(doc, "key"), doc, Team) }); err != nil {
		return err
	}

	var hubs struct {
		Hubs []struct {
			ID string `json:"id"`
		} `json:"hubs"`
	}
	resp, err = c.GetHubs(ctx)
	if _, err := platform.Decode(resp, err, &hubs); err != nil {
		return platform.Failed(err, "Failed to get the hubs")
	}
	ids = nil
	for _, h := range hubs.Hubs {
		ids = append(ids, h.ID)
	}
	// An imported template keeps the id it has in its hub, which is how one is told.
	imported := map[string]bool{}
	if err := e.each(ids, "hub", byUUID(func(id openapi_types.UUID) (*http.Response, error) { return c.GetHubById(ctx, id) }),
		func(doc *jsyaml.Map) error {
			templates, _ := doc.Get("templates")
			list, _ := templates.([]any)
			for _, t := range list {
				if m, ok := t.(*jsyaml.Map); ok {
					imported[str(m, "id")] = true
				}
			}
			if providedHubs[str(doc, "id")] {
				return nil
			}
			return e.write("hubs", str(doc, "hubName"), doc, Hub)
		}); err != nil {
		return err
	}

	var templates struct {
		Templates []struct {
			ID string `json:"id"`
		} `json:"templates"`
	}
	// Hidden templates, and those whose actions, target types or property definitions are
	// not available right now, are the tenant's too; the platform lists them only if asked.
	all := true
	resp, err = c.GetExperimentTemplates(ctx, &api.GetExperimentTemplatesParams{IncludeHidden: &all, IncludeNonAvailable: &all})
	if _, err := platform.Decode(resp, err, &templates); err != nil {
		return platform.Failed(err, "Failed to get the experiment templates")
	}
	ids = nil
	for _, t := range templates.Templates {
		if !imported[t.ID] {
			ids = append(ids, t.ID)
		}
	}
	if err := e.each(ids, "experiment template", byUUID(func(id openapi_types.UUID) (*http.Response, error) { return c.GetExperimentTemplate(ctx, id) }),
		func(doc *jsyaml.Map) error { return e.write("templates", str(doc, "templateTitle"), doc, Template) }); err != nil {
		return err
	}

	for _, k := range integration.Kinds {
		var list struct {
			Content []struct {
				ID string `json:"id"`
			} `json:"content"`
		}
		resp, err := k.FetchAll(ctx, c)
		if _, err := platform.Decode(resp, err, &list); err != nil {
			return platform.Failed(err, "Failed to get the %s", k.Plural)
		}
		ids = nil
		for _, i := range list.Content {
			ids = append(ids, i.ID)
		}
		dir := "integrations/" + k.Name
		if err := e.each(ids, k.Title, byUUID(func(id openapi_types.UUID) (*http.Response, error) { return k.Fetch(ctx, c, id) }),
			func(doc *jsyaml.Map) error {
				k.MaskSecrets(doc)
				e.counts["integrations"]++
				return e.write(dir, str(doc, "name"), doc, Integrations[k.Name])
			}); err != nil {
			return err
		}
	}

	// Profiles Steadybit provides come with the platform.
	type profile struct {
		ID     string `json:"id"`
		Origin string `json:"origin"`
	}
	profiles, err := platform.AllPages[profile](func(page, size int32) (*http.Response, error) {
		return c.GetProfiles(ctx, &api.GetProfilesParams{Page: api.PageRequestAO{Page: &page, Size: &size}})
	})
	if err != nil {
		return platform.Failed(err, "Failed to get the service profiles")
	}
	ids = nil
	for _, p := range profiles {
		if p.Origin == "CUSTOM" {
			ids = append(ids, p.ID)
		}
	}
	if err := e.each(ids, "service profile", byUUID(func(id openapi_types.UUID) (*http.Response, error) { return c.GetProfile(ctx, id) }),
		func(doc *jsyaml.Map) error { return e.write("service-profiles", str(doc, "name"), doc, ServiceProfile) }); err != nil {
		return err
	}

	fmt.Printf("Exported the tenant to %s: %d experiment templates, %d environments, %d teams, %d property definitions, %d hubs, %d integrations, %d service profiles.\n",
		o.Directory, e.counts["templates"], e.counts["environments"], e.counts["teams"], e.counts["property-definitions"], e.counts["hubs"],
		e.counts["integrations"], e.counts["service-profiles"])
	return nil
}
