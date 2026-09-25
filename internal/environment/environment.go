// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package environment implements the `environment` commands.
package environment

import (
	"context"
	"fmt"
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/steadybit/cli/api"
	"github.com/steadybit/cli/internal/output"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/resource"
	"github.com/steadybit/cli/internal/table"
)

// The state is the platform's, and the version is dropped as `service get` drops it: kept
// in a file, it turns every apply after an edit in the UI into a conflict.
var readOnly = []string{"version", "state"}

func uuid(id string) (openapi_types.UUID, error) {
	u, ok := resource.UUID(id)
	if !ok {
		return u, fmt.Errorf("Environment %s not found.", id)
	}
	return u, nil
}

func notFoundOr(err error, id, format string) error {
	if platform.IsStatus(err, http.StatusNotFound) {
		return fmt.Errorf("Environment %s not found.", id)
	}
	return platform.Failed(err, format, id)
}

type ListOptions struct {
	Search string
}

func List(ctx context.Context, c *platform.Client, o ListOptions) error {
	var summaries struct {
		Environments []struct {
			ID, Name, State, Query string
		} `json:"environments"`
	}
	params := &api.GetEnvironmentsParams{}
	if o.Search != "" {
		params.Search = &o.Search
	}
	resp, err := c.GetEnvironments(ctx, params)
	if _, err := platform.Decode(resp, err, &summaries); err != nil {
		return platform.Failed(err, "Failed to get the environments")
	}
	if len(summaries.Environments) == 0 {
		fmt.Println("No environments found.")
		return nil
	}
	t := table.New(
		table.Column{Name: "id", Title: "Id", Alignment: table.Left},
		table.Column{Name: "name", Title: "Name", Alignment: table.Left},
		table.Column{Name: "state", Title: "State", Alignment: table.Left},
	)
	for _, e := range summaries.Environments {
		t.AddRow(table.Default, table.Cell("id", e.ID), table.Cell("name", e.Name), table.Cell("state", e.State))
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
	doc, _, err := platform.ReadDocument(c.GetEnvironment(ctx, id))
	if err != nil {
		return notFoundOr(err, o.ID, "Failed to get environment %s")
	}
	if err := resource.Output(resource.Strip(doc, readOnly...), o.File, o.Type); err != nil {
		return err
	}
	if o.File != "" {
		fmt.Printf("Environment %s written to %s.\n", o.ID, o.File)
	}
	return nil
}

type ApplyOptions struct {
	Files     []string
	Recursive bool
}

func Apply(ctx context.Context, c *platform.Client, o ApplyOptions) error {
	return resource.ApplyFiles(o.Files, o.Recursive, "environment", func(file string, doc *output.Document) (resource.Applied, error) {
		name, _ := doc.Get("name")
		if name == "" {
			return resource.Applied{}, fmt.Errorf("Environment file '%s' does not name the environment.", file)
		}
		var saved struct{ ID, Name string }
		resp, err := c.UpsertEnvironmentWithBody(ctx, "application/json", resource.Body(resource.Strip(doc, readOnly...).Value()))
		resp, err = platform.Decode(resp, err, &saved)
		if err != nil {
			return resource.Applied{}, platform.Failed(err, "Failed to save environment %s", name)
		}
		created := resp.StatusCode == http.StatusCreated
		fmt.Printf("Environment %s (%s) %s.\n", saved.Name, saved.ID, resource.CreatedOrUpdated(created))
		return resource.Applied{ID: saved.ID, Created: created}, nil
	})
}

type DeleteOptions struct {
	ID  string
	Yes bool
}

func Delete(ctx context.Context, c *platform.Client, o DeleteOptions) error {
	id, err := uuid(o.ID)
	if err != nil {
		return err
	}
	if ok, err := resource.Confirmed(o.Yes, fmt.Sprintf("Delete environment %s?", o.ID)); !ok || err != nil {
		return err
	}
	if _, _, err := platform.Read(c.DeleteEnvironment(ctx, id)); err != nil {
		return notFoundOr(err, o.ID, "Failed to delete environment %s")
	}
	fmt.Printf("Environment %s deleted.\n", o.ID)
	return nil
}

type VariableGetOptions struct {
	ID, Type string
}

func GetVariables(ctx context.Context, c *platform.Client, o VariableGetOptions) error {
	id, err := uuid(o.ID)
	if err != nil {
		return err
	}
	// The spec declares a string; the platform answers with the map of variables.
	doc, _, err := platform.ReadDocument(c.GetEnvironmentVariables(ctx, id))
	if err != nil {
		return notFoundOr(err, o.ID, "Failed to get the variables of environment %s")
	}
	return resource.Output(doc, "", o.Type)
}

type VariableSetOptions struct {
	ID      string
	File    string
	Replace bool
}

// SetVariables merges KEY=VALUE arguments over a file's variables with PUT; with
// --replace, POST makes them the only ones.
func SetVariables(ctx context.Context, c *platform.Client, pairs []string, o VariableSetOptions) error {
	variables, err := resource.Variables(pairs, o.File, o.Replace)
	if err != nil {
		return err
	}
	id, err := uuid(o.ID)
	if err != nil {
		return err
	}
	body := resource.Body(variables)
	if o.Replace {
		_, _, err = platform.Read(c.SetEnvironmentVariablesWithBody(ctx, id, "application/json", body))
	} else {
		_, _, err = platform.Read(c.UpdateEnvironmentVariablesWithBody(ctx, id, "application/json", body))
	}
	if err != nil {
		return notFoundOr(err, o.ID, "Failed to update the variables of environment %s")
	}
	fmt.Printf("%d variable(s) of environment %s %s.\n", variables.Len(), o.ID, resource.VariablesOutcome(o.Replace))
	return nil
}
