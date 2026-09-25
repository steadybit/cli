// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package integration implements the `integration` commands. The four kinds of integration
// have the same endpoints, so one implementation serves them all.
package integration

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/steadybit/cli/internal/output"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/resource"
	"github.com/steadybit/cli/internal/table"
)

// The version is dropped as `service get` drops it: kept in a file, it turns every apply
// after an edit in the UI into a conflict.
var readOnly = []string{"version"}

type Kind struct {
	Name   string // as the command names it
	Title  string // as messages name it
	Plural string
	// The column that says where the integration reports to.
	Column, ColumnTitle string

	list   func(ctx context.Context, c *platform.Client) (*http.Response, error)
	get    func(ctx context.Context, c *platform.Client, id openapi_types.UUID) (*http.Response, error)
	upsert func(ctx context.Context, c *platform.Client, body io.Reader) (*http.Response, error)
	delete func(ctx context.Context, c *platform.Client, id openapi_types.UUID) (*http.Response, error)
}

var (
	Webhook = Kind{
		Name: "webhook", Title: "Webhook integration", Plural: "webhook integrations", Column: "url", ColumnTitle: "URL",
		list: func(ctx context.Context, c *platform.Client) (*http.Response, error) { return c.GetCustomWebhooks(ctx) },
		get: func(ctx context.Context, c *platform.Client, id openapi_types.UUID) (*http.Response, error) {
			return c.GetCustomWebhook(ctx, id)
		},
		upsert: func(ctx context.Context, c *platform.Client, body io.Reader) (*http.Response, error) {
			return c.UpsertCustomWebhookWithBody(ctx, "application/json", body)
		},
		delete: func(ctx context.Context, c *platform.Client, id openapi_types.UUID) (*http.Response, error) {
			return c.DeleteCustomWebhook(ctx, id)
		},
	}
	Slack = Kind{
		Name: "slack", Title: "Slack integration", Plural: "Slack integrations", Column: "channel", ColumnTitle: "Channel",
		list: func(ctx context.Context, c *platform.Client) (*http.Response, error) {
			return c.GetSlackIntegrations(ctx)
		},
		get: func(ctx context.Context, c *platform.Client, id openapi_types.UUID) (*http.Response, error) {
			return c.GetSlackIntegration(ctx, id)
		},
		upsert: func(ctx context.Context, c *platform.Client, body io.Reader) (*http.Response, error) {
			return c.UpsertSlackIntegrationWithBody(ctx, "application/json", body)
		},
		delete: func(ctx context.Context, c *platform.Client, id openapi_types.UUID) (*http.Response, error) {
			return c.DeleteSlackIntegration(ctx, id)
		},
	}
	Preflight = Kind{
		Name: "preflight", Title: "Preflight webhook", Plural: "preflight webhooks", Column: "url", ColumnTitle: "URL",
		list: func(ctx context.Context, c *platform.Client) (*http.Response, error) {
			return c.GetPreflightWebhooks(ctx)
		},
		get: func(ctx context.Context, c *platform.Client, id openapi_types.UUID) (*http.Response, error) {
			return c.GetPreflightWebhook(ctx, id)
		},
		upsert: func(ctx context.Context, c *platform.Client, body io.Reader) (*http.Response, error) {
			return c.UpsertPreflightWebhookWithBody(ctx, "application/json", body)
		},
		delete: func(ctx context.Context, c *platform.Client, id openapi_types.UUID) (*http.Response, error) {
			return c.DeletePreflightWebhook(ctx, id)
		},
	}
	PreflightAction = Kind{
		Name: "preflight-action", Title: "Preflight action integration", Plural: "preflight action integrations", Column: "preflightActionId", ColumnTitle: "Preflight action",
		list: func(ctx context.Context, c *platform.Client) (*http.Response, error) {
			return c.GetPreflightActionIntegrations(ctx)
		},
		get: func(ctx context.Context, c *platform.Client, id openapi_types.UUID) (*http.Response, error) {
			return c.GetPreflightActionIntegration(ctx, id)
		},
		upsert: func(ctx context.Context, c *platform.Client, body io.Reader) (*http.Response, error) {
			return c.UpsertPreflightActionIntegrationWithBody(ctx, "application/json", body)
		},
		delete: func(ctx context.Context, c *platform.Client, id openapi_types.UUID) (*http.Response, error) {
			return c.DeletePreflightActionIntegration(ctx, id)
		},
	}
	Kinds = []Kind{Webhook, Slack, Preflight, PreflightAction}
)

func (k Kind) uuid(id string) (openapi_types.UUID, error) {
	u, ok := resource.UUID(id)
	if !ok {
		return u, fmt.Errorf("%s %s not found.", k.Title, id)
	}
	return u, nil
}

func (k Kind) notFoundOr(err error, id, format string) error {
	if platform.IsStatus(err, http.StatusNotFound) {
		return fmt.Errorf("%s %s not found.", k.Title, id)
	}
	return platform.Failed(err, format, k.Title, id)
}

func List(ctx context.Context, c *platform.Client, k Kind) error {
	var result struct {
		Content []map[string]any `json:"content"`
	}
	resp, err := k.list(ctx, c)
	if _, err := platform.Decode(resp, err, &result); err != nil {
		return platform.Failed(err, "Failed to get the %s", k.Plural)
	}
	if len(result.Content) == 0 {
		fmt.Printf("No %s found.\n", k.Plural)
		return nil
	}
	t := table.New(
		table.Column{Name: "id", Title: "Id", Alignment: table.Left},
		table.Column{Name: "name", Title: "Name", Alignment: table.Left},
		table.Column{Name: "scope", Title: "Scope", Alignment: table.Left},
		table.Column{Name: "team", Title: "Team", Alignment: table.Left},
		table.Column{Name: k.Column, Title: k.ColumnTitle, Alignment: table.Left},
	)
	for _, i := range result.Content {
		t.AddRow(table.Default, table.Cell("id", i["id"]), table.Cell("name", i["name"]), table.Cell("scope", i["scope"]),
			table.Cell("team", i["team"]), table.Cell(k.Column, i[k.Column]))
	}
	t.Print()
	return nil
}

type GetOptions struct {
	ID, File, Type string
}

func Get(ctx context.Context, c *platform.Client, k Kind, o GetOptions) error {
	id, err := k.uuid(o.ID)
	if err != nil {
		return err
	}
	doc, _, err := platform.ReadDocument(k.get(ctx, c, id))
	if err != nil {
		return k.notFoundOr(err, o.ID, "Failed to get %s %s")
	}
	if err := resource.Output(resource.Strip(doc, readOnly...), o.File, o.Type); err != nil {
		return err
	}
	if o.File != "" {
		fmt.Printf("%s %s written to %s.\n", k.Title, o.ID, o.File)
	}
	return nil
}

type ApplyOptions struct {
	Files     []string
	Recursive bool
}

func Apply(ctx context.Context, c *platform.Client, k Kind, o ApplyOptions) error {
	return resource.ApplyFiles(o.Files, o.Recursive, k.Name+" integration", func(file string, doc *output.Document) (resource.Applied, error) {
		name, _ := doc.Get("name")
		if name == "" {
			return resource.Applied{}, fmt.Errorf("%s file '%s' does not name the integration.", k.Title, file)
		}
		// The platform masks secrets when reading them back and rejects the mask, while
		// leaving the secret out removes it. Neither is what a file from `get` means.
		if secret, _ := doc.Get("secret"); secret != "" && strings.Trim(secret, "*") == "" {
			return resource.Applied{}, fmt.Errorf("%s file '%s' holds the masked secret `get` writes. Put the secret in, or remove it for none.", k.Title, file)
		}
		var saved struct{ ID, Name string }
		resp, err := k.upsert(ctx, c, resource.Body(resource.Strip(doc, readOnly...).Value()))
		resp, err = platform.Decode(resp, err, &saved)
		if err != nil {
			return resource.Applied{}, platform.Failed(err, "Failed to save %s %s", k.Title, name)
		}
		created := resp.StatusCode == http.StatusCreated
		fmt.Printf("%s %s (%s) %s.\n", k.Title, saved.Name, saved.ID, resource.CreatedOrUpdated(created))
		return resource.Applied{ID: saved.ID, Created: created}, nil
	})
}

type DeleteOptions struct {
	ID  string
	Yes bool
}

func Delete(ctx context.Context, c *platform.Client, k Kind, o DeleteOptions) error {
	id, err := k.uuid(o.ID)
	if err != nil {
		return err
	}
	if ok, err := resource.Confirmed(o.Yes, fmt.Sprintf("Delete %s %s?", k.Title, o.ID)); !ok || err != nil {
		return err
	}
	if _, _, err := platform.Read(k.delete(ctx, c, id)); err != nil {
		return k.notFoundOr(err, o.ID, "Failed to delete %s %s")
	}
	fmt.Printf("%s %s deleted.\n", k.Title, o.ID)
	return nil
}
