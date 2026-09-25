// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package template implements `template list` and `template get`.
package template

import (
	"context"
	"fmt"
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/steadybit/cli/api"
	"github.com/steadybit/cli/internal/jsyaml"
	"github.com/steadybit/cli/internal/output"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/resource"
	"github.com/steadybit/cli/internal/table"
)

type ListOptions struct {
	Tags, TargetTypes, Actions, Search []string
}

func optional(values []string) *[]string {
	if len(values) == 0 {
		return nil
	}
	return &values
}

func List(ctx context.Context, c *platform.Client, o ListOptions) error {
	var summaries struct {
		Templates []struct {
			ID            string `json:"id"`
			TemplateTitle string `json:"templateTitle"`
		} `json:"templates"`
	}
	resp, err := c.GetExperimentTemplates(ctx, &api.GetExperimentTemplatesParams{
		Tag: optional(o.Tags), TargetType: optional(o.TargetTypes), Action: optional(o.Actions), FreeTextPhrases: optional(o.Search),
	})
	if _, err := platform.Decode(resp, err, &summaries); err != nil {
		return platform.Failed(err, "Failed to get the experiment templates")
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
	if err := resource.Output(doc, o.File, o.Type); err != nil {
		return err
	}
	if o.File != "" {
		fmt.Printf("Experiment template %s written to %s.\n", o.ID, o.File)
	}
	return nil
}
