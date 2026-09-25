// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package hub implements the `hub` commands.
package hub

import (
	"context"
	"fmt"
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/steadybit/cli/v6/api"
	"github.com/steadybit/cli/v6/internal/output"
	"github.com/steadybit/cli/v6/internal/platform"
	"github.com/steadybit/cli/v6/internal/resource"
	"github.com/steadybit/cli/v6/internal/table"
)

// What the last synchronisation found and who edited the hub is the platform's. The
// version is dropped as `service get` drops it.
var readOnly = []string{"version", "templates", "lastSync", "lastRepositoryChange", "syncError", "created", "createdBy", "edited", "editedBy"}

func uuid(id string) (openapi_types.UUID, error) {
	u, ok := resource.UUID(id)
	if !ok {
		return u, fmt.Errorf("Hub %s not found.", id)
	}
	return u, nil
}

func notFoundOr(err error, id, format string) error {
	if platform.IsStatus(err, http.StatusNotFound) {
		return fmt.Errorf("Hub %s not found.", id)
	}
	return platform.Failed(err, format, id)
}

func List(ctx context.Context, c *platform.Client, explicitType string) error {
	var summaries struct {
		Hubs []struct {
			ID      string `json:"id"`
			HubName string `json:"hubName"`
		} `json:"hubs"`
	}
	resp, err := c.GetHubs(ctx)
	raw, err := resource.DecodeListed(resp, err, "hubs", &summaries)
	if err != nil {
		return platform.Failed(err, "Failed to get the hubs")
	}
	if resource.Machine(explicitType) {
		return resource.List(raw, explicitType, nil)
	}
	if len(summaries.Hubs) == 0 {
		fmt.Println("No hubs found.")
		return nil
	}
	t := table.New(table.Column{Name: "id", Title: "Id", Alignment: table.Left}, table.Column{Name: "name", Title: "Name", Alignment: table.Left})
	for _, h := range summaries.Hubs {
		t.AddRow(table.Default, table.Cell("id", h.ID), table.Cell("name", h.HubName))
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
	doc, _, err := platform.ReadDocument(c.GetHubById(ctx, id))
	if err != nil {
		return notFoundOr(err, o.ID, "Failed to get hub %s")
	}
	if err := resource.Output(resource.Strip(doc, readOnly...), o.File, o.Type); err != nil {
		return err
	}
	if o.File != "" {
		fmt.Printf("Hub %s written to %s.\n", o.ID, o.File)
	}
	return nil
}

type ApplyOptions struct {
	Files       []string
	Recursive   bool
	Synchronize bool
}

func Apply(ctx context.Context, c *platform.Client, o ApplyOptions) error {
	return resource.ApplyFiles(o.Files, o.Recursive, "hub", func(file string, doc *output.Document) (resource.Applied, error) {
		name, _ := doc.Get("hubName")
		if name == "" {
			return resource.Applied{}, fmt.Errorf("Hub file '%s' does not name a hubName.", file)
		}
		var saved struct {
			ID      string `json:"id"`
			HubName string `json:"hubName"`
		}
		resp, err := c.UpsertHubWithBody(ctx, &api.UpsertHubParams{Synchronize: &o.Synchronize}, "application/json", resource.Body(resource.Strip(doc, readOnly...).Value()))
		resp, err = platform.Decode(resp, err, &saved)
		if err != nil {
			return resource.Applied{}, platform.Failed(err, "Failed to save hub %s", name)
		}
		created := resp.StatusCode == http.StatusCreated
		fmt.Printf("Hub %s (%s) %s.\n", saved.HubName, saved.ID, resource.CreatedOrUpdated(created))
		return resource.Applied{ID: saved.ID, Created: created}, nil
	})
}

type DeleteOptions struct {
	ID        string
	Templates bool
	Yes       bool
}

func Delete(ctx context.Context, c *platform.Client, o DeleteOptions) error {
	id, err := uuid(o.ID)
	if err != nil {
		return err
	}
	question := fmt.Sprintf("Delete hub %s? The templates imported from it are kept.", o.ID)
	if o.Templates {
		question = fmt.Sprintf("Delete hub %s and the templates imported from it?", o.ID)
	}
	if ok, err := resource.Confirmed(o.Yes, question); !ok || err != nil {
		return err
	}
	if _, _, err := platform.Read(c.DeleteHub(ctx, id, &api.DeleteHubParams{DeleteImportedTemplates: &o.Templates})); err != nil {
		return notFoundOr(err, o.ID, "Failed to delete hub %s")
	}
	fmt.Printf("Hub %s deleted.\n", o.ID)
	return nil
}

// Resync fetches the hub's repository again; the platform answers once it is done.
func Resync(ctx context.Context, c *platform.Client, idText string) error {
	id, err := uuid(idText)
	if err != nil {
		return err
	}
	var hub struct {
		HubName   string `json:"hubName"`
		Templates []any  `json:"templates"`
		SyncError string `json:"syncError"`
	}
	resp, err := c.ResyncHub(ctx, id)
	if _, err := platform.Decode(resp, err, &hub); err != nil {
		return notFoundOr(err, idText, "Failed to resynchronize hub %s")
	}
	if hub.SyncError != "" {
		return fmt.Errorf("Hub %s could not be synchronized: %s", hub.HubName, hub.SyncError)
	}
	fmt.Printf("Hub %s synchronized, %d template(s).\n", hub.HubName, len(hub.Templates))
	return nil
}
