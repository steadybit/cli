// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package serviceprofile implements the `service-profile` commands.
package serviceprofile

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/steadybit/cli/api"
	"github.com/steadybit/cli/internal/jsyaml"
	"github.com/steadybit/cli/internal/output"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/resource"
	"github.com/steadybit/cli/internal/service"
	"github.com/steadybit/cli/internal/table"
)

// Only what can be sent back is kept. Whether a profile is the default is changed in the
// platform, not through the file.
var readOnly = []string{"created", "createdBy", "edited", "editedBy", "version", "defaultProfile"}

func uuid(id string) (openapi_types.UUID, error) {
	var u openapi_types.UUID
	if err := u.UnmarshalText([]byte(id)); err != nil {
		return u, fmt.Errorf("Service profile %s not found.", id)
	}
	return u, nil
}

func notFoundOr(err error, id, format string) error {
	if platform.IsStatus(err, http.StatusNotFound) {
		return fmt.Errorf("Service profile %s not found.", id)
	}
	return platform.Failed(err, format, id)
}

type ListOptions struct {
	Name    string
	Origins []string
	Default bool
	Type    string
}

func List(ctx context.Context, c *platform.Client, o ListOptions) error {
	type profile struct {
		ID             string `json:"id"`
		Name           string `json:"name"`
		Origin         string `json:"origin"`
		DefaultProfile bool   `json:"defaultProfile"`
		Templates      []struct {
			TemplateIDs []string `json:"templateIds"`
		} `json:"templates"`
	}
	params := api.GetProfilesParams{}
	if o.Name != "" {
		params.Name = &o.Name
	}
	if len(o.Origins) > 0 {
		origins := make([]string, len(o.Origins))
		for i, origin := range o.Origins {
			origins[i] = strings.ToUpper(origin)
		}
		params.Origin = &origins
	}
	if o.Default {
		params.DefaultProfile = &o.Default
	}
	raw, err := platform.AllPagesRaw(func(page, size int32) (*http.Response, error) {
		p := params
		p.Page = api.PageRequestAO{Page: &page, Size: &size}
		return c.GetProfiles(ctx, &p)
	})
	if err != nil {
		return platform.Failed(err, "Failed to get the service profiles")
	}
	if resource.Machine(o.Type) {
		return resource.List(raw, o.Type, nil)
	}
	var profiles []profile
	if err := resource.DecodeEach(raw, &profiles); err != nil {
		return err
	}
	if len(profiles) == 0 {
		fmt.Println("No service profiles found.")
		return nil
	}
	t := table.New(
		table.Column{Name: "id", Title: "Id", Alignment: table.Left},
		table.Column{Name: "name", Title: "Name", Alignment: table.Left},
		table.Column{Name: "origin", Title: "Origin", Alignment: table.Left},
		table.Column{Name: "defaultProfile", Title: "Default", Alignment: table.Left},
		table.Column{Name: "templates", Title: "Templates"},
	)
	for _, p := range profiles {
		count := 0
		for _, category := range p.Templates {
			count += len(category.TemplateIDs)
		}
		t.AddRow(table.Default, table.Cell("id", p.ID), table.Cell("name", p.Name), table.Cell("origin", p.Origin),
			table.Cell("defaultProfile", p.DefaultProfile), table.Cell("templates", count))
	}
	t.Print()
	return nil
}

type GetOptions struct {
	ID, File, Type string
}

func Get(ctx context.Context, c *platform.Client, o GetOptions) error {
	id, err := uuid(o.ID)
	if err != nil {
		return err
	}
	doc, _, err := platform.ReadDocument(c.GetProfile(ctx, id))
	if err != nil {
		return notFoundOr(err, o.ID, "Failed to get service profile %s")
	}
	if err := resource.Output(resource.Strip(doc, readOnly...), o.File, o.Type); err != nil {
		return err
	}
	if o.File != "" {
		fmt.Printf("Service profile %s written to %s.\n", o.ID, o.File)
	}
	return nil
}

type ApplyOptions struct {
	Files             []string
	Recursive         bool
	DeleteExperiments bool
}

func Apply(ctx context.Context, c *platform.Client, o ApplyOptions) error {
	return resource.ApplyFiles(o.Files, o.Recursive, "service profile", func(file string, doc *output.Document) (resource.Applied, error) {
		name, _ := doc.Get("name")
		if name == "" {
			return resource.Applied{}, fmt.Errorf("Service profile file '%s' does not name the profile.", file)
		}
		profile := resource.Strip(doc, readOnly...).Value()
		// Profiles written by hand are the team's own; PROVIDED ones come from Steadybit.
		if origin, ok := profile.Get("origin"); !ok || origin == nil {
			profile.Set("origin", "CUSTOM")
		}
		var saved struct{ ID, Name string }
		resp, err := c.UpsertProfileWithBody(ctx, &api.UpsertProfileParams{DeleteExperiments: &o.DeleteExperiments}, "application/json",
			bytes.NewReader([]byte(jsyaml.CompactJSON(profile))))
		resp, err = platform.Decode(resp, err, &saved)
		if err != nil {
			if !o.DeleteExperiments && service.RefusedForProvidedExperiments(err) {
				return resource.Applied{}, fmt.Errorf("Service profile %s was not saved: the change would remove provided experiments. Pass --delete-experiments to delete them.", name)
			}
			return resource.Applied{}, platform.Failed(err, "Failed to save service profile %s", name)
		}
		created := resp.StatusCode == http.StatusCreated
		fmt.Printf("Service profile %s (%s) %s.\n", saved.Name, saved.ID, resource.CreatedOrUpdated(created))
		return resource.Applied{ID: saved.ID, Created: created}, nil
	})
}

func Delete(ctx context.Context, c *platform.Client, idText string) error {
	id, err := uuid(idText)
	if err != nil {
		return err
	}
	if _, _, err := platform.Read(c.DeleteProfile(ctx, id)); err != nil {
		if platform.IsStatus(err, http.StatusUnprocessableEntity) {
			return fmt.Errorf("Service profile %s is provided by Steadybit and cannot be deleted.", idText)
		}
		return notFoundOr(err, idText, "Failed to delete service profile %s")
	}
	fmt.Printf("Service profile %s deleted.\n", idText)
	return nil
}
