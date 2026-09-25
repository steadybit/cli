// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package property implements the `property` commands: the definitions of the properties
// experiments and services carry, and their associations, which say who carries them.
package property

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/steadybit/cli/v6/api"
	"github.com/steadybit/cli/v6/internal/output"
	"github.com/steadybit/cli/v6/internal/platform"
	"github.com/steadybit/cli/v6/internal/resource"
	"github.com/steadybit/cli/v6/internal/table"
)

// The version is dropped as `service get` drops it: kept in a file, it turns every apply
// after an edit in the UI into a conflict.
var readOnly = []string{"version"}

func definitionNotFoundOr(err error, key, format string) error {
	if platform.IsStatus(err, http.StatusNotFound) {
		return fmt.Errorf("Property definition %s not found.", key)
	}
	return platform.Failed(err, format, key)
}

func ListDefinitions(ctx context.Context, c *platform.Client, explicitType string) error {
	type definition struct {
		Key, Label, DataType string
		EnumValues           []string `json:"enumValues"`
	}
	raw, err := platform.AllPagesRaw(func(page, size int32) (*http.Response, error) {
		return c.GetPropertyDefinitions(ctx, &api.GetPropertyDefinitionsParams{Page: api.PageRequestAO{Page: &page, Size: &size}})
	})
	if err != nil {
		return platform.Failed(err, "Failed to get the property definitions")
	}
	if resource.Machine(explicitType) {
		return resource.List(raw, explicitType, nil)
	}
	var definitions []definition
	if err := resource.DecodeEach(raw, &definitions); err != nil {
		return err
	}
	if len(definitions) == 0 {
		fmt.Println("No property definitions found.")
		return nil
	}
	t := table.New(
		table.Column{Name: "key", Title: "Key", Alignment: table.Left},
		table.Column{Name: "label", Title: "Label", Alignment: table.Left},
		table.Column{Name: "dataType", Title: "Type", Alignment: table.Left},
		table.Column{Name: "enumValues", Title: "Values", Alignment: table.Left},
	)
	for _, d := range definitions {
		t.AddRow(table.Default, table.Cell("key", d.Key), table.Cell("label", d.Label), table.Cell("dataType", d.DataType),
			table.Cell("enumValues", values(d.EnumValues)))
	}
	t.Print()
	return nil
}

// values shows the first few enum values; `get` has them all.
func values(enum []string) string {
	const shown = 5
	if len(enum) > shown {
		return strings.Join(enum[:shown], ", ") + fmt.Sprintf(" and %d more", len(enum)-shown)
	}
	return strings.Join(enum, ", ")
}

type GetDefinitionOptions struct {
	Key, File, Type string
}

func GetDefinition(ctx context.Context, c *platform.Client, o GetDefinitionOptions) error {
	doc, _, err := platform.ReadDocument(c.GetPropertyDefinition(ctx, o.Key))
	if err != nil {
		return definitionNotFoundOr(err, o.Key, "Failed to get property definition %s")
	}
	if err := resource.Output(resource.Strip(doc, readOnly...), o.File, o.Type); err != nil {
		return err
	}
	if o.File != "" {
		fmt.Printf("Property definition %s written to %s.\n", o.Key, o.File)
	}
	return nil
}

type ApplyDefinitionOptions struct {
	Files        []string
	Recursive    bool
	DeleteValues bool
}

// ApplyDefinitions upserts definitions by their key; there is no id to write back.
func ApplyDefinitions(ctx context.Context, c *platform.Client, o ApplyDefinitionOptions) error {
	return resource.ApplyFiles(o.Files, o.Recursive, "property definition", func(file string, doc *output.Document) (resource.Applied, error) {
		key, _ := doc.Get("key")
		if key == "" {
			return resource.Applied{}, fmt.Errorf("Property definition file '%s' does not name the key.", file)
		}
		resp, err := c.UpsertPropertyDefinitionWithBody(ctx, &api.UpsertPropertyDefinitionParams{DeleteValues: &o.DeleteValues}, "application/json",
			resource.Body(resource.Strip(doc, readOnly...).Value()))
		_, resp, err = platform.Read(resp, err)
		if err != nil {
			return resource.Applied{}, platform.Failed(err, "Failed to save property definition %s", key)
		}
		created := resp.StatusCode == http.StatusCreated
		fmt.Printf("Property definition %s %s.\n", key, resource.CreatedOrUpdated(created))
		return resource.Applied{Created: created}, nil
	})
}

type DeleteDefinitionOptions struct {
	Key          string
	Associations bool
	Yes          bool
}

func DeleteDefinition(ctx context.Context, c *platform.Client, o DeleteDefinitionOptions) error {
	question := fmt.Sprintf("Delete property definition %s?", o.Key)
	if o.Associations {
		question = fmt.Sprintf("Delete property definition %s and its associations?", o.Key)
	}
	if ok, err := resource.Confirmed(o.Yes, question); !ok || err != nil {
		return err
	}
	params := &api.DeletePropertyDefinitionParams{DeleteAssociations: &o.Associations}
	if _, _, err := platform.Read(c.DeletePropertyDefinition(ctx, o.Key, params)); err != nil {
		return definitionNotFoundOr(err, o.Key, "Failed to delete property definition %s")
	}
	fmt.Printf("Property definition %s deleted.\n", o.Key)
	return nil
}

func associationUUID(id string) (openapi_types.UUID, error) {
	u, ok := resource.UUID(id)
	if !ok {
		return u, fmt.Errorf("Property association %s not found.", id)
	}
	return u, nil
}

func associationNotFoundOr(err error, id, format string) error {
	if platform.IsStatus(err, http.StatusNotFound) {
		return fmt.Errorf("Property association %s not found.", id)
	}
	return platform.Failed(err, format, id)
}

type ListAssociationOptions struct {
	Key, Experiment, Service, Type string
	// Output is -t for the other listings; --type already names the association type here.
	Output string
}

func ListAssociations(ctx context.Context, c *platform.Client, o ListAssociationOptions) error {
	params := api.GetAssociationsParams{}
	if o.Key != "" {
		params.Key = &o.Key
	}
	if o.Experiment != "" {
		params.ExperimentKey = &o.Experiment
	}
	if o.Service != "" {
		service, ok := resource.UUID(o.Service)
		if !ok {
			return fmt.Errorf("Service %s not found.", o.Service)
		}
		params.ServiceId = &service
	}
	if o.Type != "" {
		kind := api.GetAssociationsParamsAssociationTypeAO(strings.ToUpper(o.Type))
		if !kind.Valid() {
			return fmt.Errorf("--type must be EXPERIMENT or SERVICE, not '%s'.", o.Type)
		}
		params.AssociationTypeAO = &kind
	}
	type association struct {
		ID, Key, AssociationType      string
		ExperimentKey, ServiceID      *string
		Required, EditableInExecution *bool
	}
	// The page size is the platform's; it takes only the page number.
	raw, err := platform.AllPagesRaw(func(page, _ int32) (*http.Response, error) {
		params.Page = &page
		return c.GetAssociations(ctx, &params)
	})
	if err != nil {
		return platform.Failed(err, "Failed to get the property associations")
	}
	if resource.Machine(o.Output) {
		return resource.List(raw, o.Output, nil)
	}
	var associations []association
	if err := resource.DecodeEach(raw, &associations); err != nil {
		return err
	}
	if len(associations) == 0 {
		fmt.Println("No property associations found.")
		return nil
	}
	t := table.New(
		table.Column{Name: "id", Title: "Id", Alignment: table.Left},
		table.Column{Name: "key", Title: "Key", Alignment: table.Left},
		table.Column{Name: "type", Title: "Type", Alignment: table.Left},
		table.Column{Name: "for", Title: "For", Alignment: table.Left},
		table.Column{Name: "required", Title: "Required", Alignment: table.Left},
		table.Column{Name: "editable", Title: "Editable in runs", Alignment: table.Left},
	)
	boolOr := func(b *bool) string { return fmt.Sprint(b != nil && *b) }
	for _, a := range associations {
		// Without an experiment or service, an association is for all of them.
		target := "all"
		switch {
		case a.ExperimentKey != nil:
			target = *a.ExperimentKey
		case a.ServiceID != nil:
			target = *a.ServiceID
		}
		t.AddRow(table.Default, table.Cell("id", a.ID), table.Cell("key", a.Key), table.Cell("type", a.AssociationType), table.Cell("for", target),
			table.Cell("required", boolOr(a.Required)), table.Cell("editable", boolOr(a.EditableInExecution)))
	}
	t.Print()
	return nil
}

type GetAssociationOptions struct {
	ID, File, Type string
}

func GetAssociation(ctx context.Context, c *platform.Client, o GetAssociationOptions) error {
	id, err := associationUUID(o.ID)
	if err != nil {
		return err
	}
	doc, _, err := platform.ReadDocument(c.GetPropertyDefinition1(ctx, id))
	if err != nil {
		return associationNotFoundOr(err, o.ID, "Failed to get property association %s")
	}
	if err := resource.Output(resource.Strip(doc, readOnly...), o.File, o.Type); err != nil {
		return err
	}
	if o.File != "" {
		fmt.Printf("Property association %s written to %s.\n", o.ID, o.File)
	}
	return nil
}

type ApplyAssociationOptions struct {
	Files     []string
	Recursive bool
}

func ApplyAssociations(ctx context.Context, c *platform.Client, o ApplyAssociationOptions) error {
	return resource.ApplyFiles(o.Files, o.Recursive, "property association", func(file string, doc *output.Document) (resource.Applied, error) {
		key, _ := doc.Get("key")
		if key == "" {
			return resource.Applied{}, fmt.Errorf("Property association file '%s' does not name the property key.", file)
		}
		var saved struct{ ID, Key string }
		resp, err := c.UpsertPropertyAssociationWithBody(ctx, "application/json", resource.Body(resource.Strip(doc, readOnly...).Value()))
		resp, err = platform.Decode(resp, err, &saved)
		if err != nil {
			return resource.Applied{}, platform.Failed(err, "Failed to save the property association of %s", key)
		}
		created := resp.StatusCode == http.StatusCreated
		fmt.Printf("Property association %s of %s %s.\n", saved.ID, saved.Key, resource.CreatedOrUpdated(created))
		return resource.Applied{ID: saved.ID, Created: created}, nil
	})
}

type DeleteAssociationOptions struct {
	ID           string
	DeleteValues bool
	Yes          bool
}

func DeleteAssociation(ctx context.Context, c *platform.Client, o DeleteAssociationOptions) error {
	id, err := associationUUID(o.ID)
	if err != nil {
		return err
	}
	question := fmt.Sprintf("Delete property association %s?", o.ID)
	if o.DeleteValues {
		question = fmt.Sprintf("Delete property association %s and the values experiments and schedules have for it?", o.ID)
	}
	if ok, err := resource.Confirmed(o.Yes, question); !ok || err != nil {
		return err
	}
	if _, _, err := platform.Read(c.DeletePropertyAssociation(ctx, id, &api.DeletePropertyAssociationParams{DeleteValues: &o.DeleteValues})); err != nil {
		return associationNotFoundOr(err, o.ID, "Failed to delete property association %s")
	}
	fmt.Printf("Property association %s deleted.\n", o.ID)
	return nil
}
