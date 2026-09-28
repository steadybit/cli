// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package template implements the `template` commands.
package template

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/steadybit/cli/v6/api"
	"github.com/steadybit/cli/v6/internal/jsyaml"
	"github.com/steadybit/cli/v6/internal/output"
	"github.com/steadybit/cli/v6/internal/platform"
	"github.com/steadybit/cli/v6/internal/resource"
	"github.com/steadybit/cli/v6/internal/table"
)

// Who created and edited a template cannot be sent back. The version is dropped as
// `service get` drops it, so that an edit in the UI does not turn the next apply into a
// conflict.
var readOnly = []string{"created", "createdBy", "edited", "editedBy", "version"}

func notFoundOr(err error, id, format string) error {
	if platform.IsStatus(err, http.StatusNotFound) {
		return fmt.Errorf("Experiment template %s not found.", id)
	}
	return platform.Failed(err, format, id)
}

type ListOptions struct {
	Tags, TargetTypes, Actions, Search []string
	Type                               string
}

func optional(values []string) *[]string {
	if len(values) == 0 {
		return nil
	}
	return &values
}

func List(ctx context.Context, c *platform.Client, o ListOptions) error {
	var raw struct {
		Templates []json.RawMessage `json:"templates"`
	}
	resp, err := c.GetExperimentTemplates(ctx, &api.GetExperimentTemplatesParams{
		Tag: optional(o.Tags), TargetType: optional(o.TargetTypes), Action: optional(o.Actions), FreeTextPhrases: optional(o.Search),
	})
	if _, err := platform.Decode(resp, err, &raw); err != nil {
		return platform.Failed(err, "Failed to get the experiment templates")
	}
	if resource.Machine(o.Type) {
		return resource.List(raw.Templates, o.Type, nil)
	}
	var summaries struct {
		Templates []struct {
			ID            string `json:"id"`
			TemplateTitle string `json:"templateTitle"`
		}
	}
	if err := resource.DecodeEach(raw.Templates, &summaries.Templates); err != nil {
		return err
	}
	if len(summaries.Templates) == 0 {
		fmt.Println("No experiment templates found.")
		return nil
	}
	t := table.New(table.Column{Name: "id", Title: "Id", Alignment: table.Left}, table.Column{Name: "templateTitle", Title: "Title", Alignment: table.Left})
	for _, s := range summaries.Templates {
		t.AddRow(table.Default, table.Cell("id", s.ID), table.Cell("templateTitle", s.TemplateTitle))
	}
	t.Print()
	return nil
}

// Fetch gets a template, reporting one that does not exist by name.
func Fetch(ctx context.Context, c *platform.Client, id string) (*output.Document, error) {
	var uuid openapi_types.UUID
	if err := uuid.UnmarshalText([]byte(id)); err != nil {
		return nil, fmt.Errorf("Experiment template %s not found.", id)
	}
	doc, _, err := platform.ReadDocument(c.GetExperimentTemplate(ctx, uuid))
	if platform.IsStatus(err, http.StatusNotFound) {
		return nil, fmt.Errorf("Experiment template %s not found.", id)
	}
	if err != nil {
		return nil, platform.Failed(err, "Failed to get experiment template %s", id)
	}
	return doc, nil
}

type GetOptions struct {
	ID, File, Type string
	Placeholders   bool
}

func Get(ctx context.Context, c *platform.Client, o GetOptions) error {
	doc, err := Fetch(ctx, c, o.ID)
	if err != nil {
		return err
	}
	if o.Placeholders {
		// A starting point for --placeholders: every key the template asks for, empty.
		values := jsyaml.NewMap()
		placeholders, _ := doc.Value().Get("placeholders")
		list, _ := placeholders.([]any)
		for _, p := range list {
			if m, ok := p.(*jsyaml.Map); ok {
				if key, ok := m.Get("key"); ok {
					values.Set(fmt.Sprint(key), "")
				}
			}
		}
		doc = output.NewDocument(values)
	}
	if err := resource.Output(resource.Strip(doc, readOnly...), o.File, o.Type); err != nil {
		return err
	}
	if o.File != "" {
		fmt.Printf("Experiment template %s written to %s.\n", o.ID, o.File)
	}
	return nil
}

type ApplyOptions struct {
	Files     []string
	Recursive bool
}

func Apply(ctx context.Context, c *platform.Client, o ApplyOptions) error {
	return resource.ApplyFiles(o.Files, o.Recursive, "template", func(file string, doc *output.Document) (resource.Applied, error) {
		title, _ := doc.Get("templateTitle")
		if title == "" {
			return resource.Applied{}, fmt.Errorf("Template file '%s' does not name a templateTitle.", file)
		}
		var saved struct{ ID, TemplateTitle string }
		resp, err := c.UpsertExperimentTemplateWithBody(ctx, "application/json", resource.Body(resource.Strip(doc, readOnly...).Value()))
		resp, err = platform.Decode(resp, err, &saved)
		if err != nil {
			return resource.Applied{}, platform.Failed(err, "Failed to save experiment template %s", title)
		}
		created := resp.StatusCode == http.StatusCreated
		fmt.Printf("Experiment template %s (%s) %s.\n", saved.TemplateTitle, saved.ID, resource.CreatedOrUpdated(created))
		return resource.Applied{ID: saved.ID, Created: created}, nil
	})
}

type DeleteOptions struct {
	ID  string
	Yes bool
}

func Delete(ctx context.Context, c *platform.Client, o DeleteOptions) error {
	id, ok := resource.UUID(o.ID)
	if !ok {
		return fmt.Errorf("Experiment template %s not found.", o.ID)
	}
	// Deleting a template also deletes the experiments service profiles provided from it.
	if ok, err := resource.Confirmed(o.Yes, fmt.Sprintf("Delete experiment template %s, and the service experiments provided from it?", o.ID)); !ok || err != nil {
		return err
	}
	if _, _, err := platform.Read(c.DeleteExperimentTemplate(ctx, id)); err != nil {
		return notFoundOr(err, o.ID, "Failed to delete experiment template %s")
	}
	fmt.Printf("Experiment template %s deleted.\n", o.ID)
	return nil
}

type ImportOptions struct {
	Hub       string
	Templates []string
	Overwrite bool
}

func Import(ctx context.Context, c *platform.Client, o ImportOptions) error {
	hub, ok := resource.UUID(o.Hub)
	if !ok {
		return fmt.Errorf("Hub %s not found.", o.Hub)
	}
	ids := make([]openapi_types.UUID, len(o.Templates))
	for i, t := range o.Templates {
		if ids[i], ok = resource.UUID(t); !ok {
			return fmt.Errorf("Experiment template %s not found.", t)
		}
	}
	request := api.ExperimentTemplatesImportAO{HubId: hub, TemplateIds: &ids}
	_, _, err := platform.Read(c.ImportFromHub(ctx, &api.ImportFromHubParams{Overwrite: &o.Overwrite}, request))
	switch {
	case platform.IsStatus(err, http.StatusConflict):
		return errors.New("Some of the templates exist already. Pass --overwrite to replace them.")
	case platform.IsStatus(err, http.StatusNotFound):
		return fmt.Errorf("Hub %s not found.", o.Hub)
	case err != nil:
		return platform.Failed(err, "Failed to import experiment templates from hub %s", o.Hub)
	}
	fmt.Printf("%d experiment template(s) imported from hub %s.\n", len(o.Templates), o.Hub)
	return nil
}
