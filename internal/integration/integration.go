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
	"github.com/steadybit/cli/v6/internal/jsyaml"
	"github.com/steadybit/cli/v6/internal/output"
	"github.com/steadybit/cli/v6/internal/platform"
	"github.com/steadybit/cli/v6/internal/resource"
	"github.com/steadybit/cli/v6/internal/table"
)

// The version is dropped as `service get` drops it: kept in a file, it turns every apply
// after an edit in the UI into a conflict.
var ReadOnly = []string{"version"}

type Kind struct {
	Name   string // as the command names it
	Title  string // as messages name it
	Plural string
	// The column that says where the integration reports to.
	Column, ColumnTitle string
	// Secrets are the fields holding credentials; each value of a map among them is one.
	// `export` masks them, so that a tenant kept in Git does not hold them.
	Secrets []string

	list   func(ctx context.Context, c *platform.Client) (*http.Response, error)
	get    func(ctx context.Context, c *platform.Client, id openapi_types.UUID) (*http.Response, error)
	upsert func(ctx context.Context, c *platform.Client, body io.Reader) (*http.Response, error)
	delete func(ctx context.Context, c *platform.Client, id openapi_types.UUID) (*http.Response, error)
}

var (
	Webhook = Kind{
		Name: "webhook", Title: "Webhook integration", Plural: "webhook integrations", Column: "url", ColumnTitle: "URL",
		// Headers carry API keys and tokens as often as anything else.
		Secrets: []string{"secret", "headers"},
		list:    func(ctx context.Context, c *platform.Client) (*http.Response, error) { return c.GetCustomWebhooks(ctx) },
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
		// The URL of a Slack incoming webhook is its credential: whoever has it can post.
		Secrets: []string{"url"},
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
		Secrets: []string{"secret", "headers"},
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

// Fetch gets one integration as `get` reads it.
func (k Kind) Fetch(ctx context.Context, c *platform.Client, id openapi_types.UUID) (*http.Response, error) {
	return k.get(ctx, c, id)
}

// FetchAll lists the integrations of the kind, as `list` reads them.
func (k Kind) FetchAll(ctx context.Context, c *platform.Client) (*http.Response, error) {
	return k.list(ctx, c)
}

// Mask is what `export` writes in place of a credential. Like the platform's own mask of
// a secret, it is only asterisks, which is how apply and diff tell one.
const Mask = "********"

func Masked(v any) bool {
	s, ok := v.(string)
	return ok && s != "" && strings.Trim(s, "*") == ""
}

// MaskSecrets replaces the credentials in an integration with the mask.
func (k Kind) MaskSecrets(doc *jsyaml.Map) {
	for _, field := range k.Secrets {
		switch v, _ := doc.Get(field); x := v.(type) {
		case string:
			if x != "" {
				doc.Set(field, Mask)
			}
		case *jsyaml.Map:
			for _, key := range x.Keys() {
				if s, _ := x.Get(key); s != "" {
					x.Set(key, Mask)
				}
			}
		}
	}
}

// stored reads, once and only when asked, what the platform holds for the integration
// the file names: nil when it names none.
func (k Kind) stored(ctx context.Context, c *platform.Client, doc *jsyaml.Map) func() (*jsyaml.Map, error) {
	var stored *jsyaml.Map
	loaded := false
	return func() (*jsyaml.Map, error) {
		if !loaded {
			loaded = true
			id, _ := doc.Get("id")
			if u, ok := resource.UUID(fmt.Sprint(id)); ok {
				d, _, err := platform.ReadDocument(k.get(ctx, c, u))
				if err != nil && !platform.IsStatus(err, http.StatusNotFound) {
					return nil, platform.Failed(err, "Failed to get %s %s", k.Title, id)
				}
				if d != nil {
					stored = d.Value()
				}
			}
		}
		return stored, nil
	}
}

// unchanged tells a file that holds what the platform does, its secret aside, which a
// mask in the file stands for.
func unchanged(doc, stored *jsyaml.Map) bool {
	if stored == nil {
		return false
	}
	if secret, _ := stored.Get("secret"); secret == nil || secret == "" {
		return false
	}
	without := func(m *jsyaml.Map) string {
		c := jsyaml.Clone(m).(*jsyaml.Map)
		for _, f := range append([]string{"secret"}, ReadOnly...) {
			c.Delete(f)
		}
		return jsyaml.CompactJSON(c)
	}
	return without(doc) == without(stored)
}

// keepStored puts back what the platform holds in place of the masks `export` writes
// for the credentials it reads back in the clear, so that an exported file applies as
// it is. The secret it never reads back; a masked one is refused afterwards.
func (k Kind) keepStored(file string, doc *jsyaml.Map, load func() (*jsyaml.Map, error)) error {
	value := func(field, key string) (any, error) {
		stored, err := load()
		if err != nil {
			return nil, err
		}
		name := field
		var v any
		if stored != nil {
			v, _ = stored.Get(field)
			if key != "" {
				name = field + "." + key
				m, _ := v.(*jsyaml.Map)
				v = nil
				if m != nil {
					v, _ = m.Get(key)
				}
			}
		}
		if v == nil || v == "" || Masked(v) {
			return nil, fmt.Errorf("%s file '%s' holds a masked %s, and the platform has none to keep. Put the value in.", k.Title, file, name)
		}
		return v, nil
	}
	for _, field := range k.Secrets {
		if field == "secret" {
			continue
		}
		switch v, _ := doc.Get(field); x := v.(type) {
		case string:
			if Masked(x) {
				kept, err := value(field, "")
				if err != nil {
					return err
				}
				doc.Set(field, kept)
			}
		case *jsyaml.Map:
			for _, key := range x.Keys() {
				if s, _ := x.Get(key); Masked(s) {
					kept, err := value(field, key)
					if err != nil {
						return err
					}
					x.Set(key, kept)
				}
			}
		}
	}
	return nil
}

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

func List(ctx context.Context, c *platform.Client, k Kind, explicitType string) error {
	var result struct {
		Content []map[string]any `json:"content"`
	}
	resp, err := k.list(ctx, c)
	raw, err := resource.DecodeListed(resp, err, "content", &result)
	if err != nil {
		return platform.Failed(err, "Failed to get the %s", k.Plural)
	}
	if resource.Machine(explicitType) {
		return resource.List(raw, explicitType, nil)
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
	if err := resource.Output(resource.Strip(doc, ReadOnly...), o.File, o.Type); err != nil {
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
		stored := k.stored(ctx, c, doc.Value())
		if err := k.keepStored(file, doc.Value(), stored); err != nil {
			return resource.Applied{}, err
		}
		// The platform masks secrets when reading them back and rejects the mask, while
		// leaving the secret out removes it. Neither is what a file from `get` means. Left
		// as it was written, though, the file has nothing to apply, as `apply -d` finds.
		if secret, _ := doc.Value().Get("secret"); Masked(secret) {
			held, err := stored()
			if err != nil {
				return resource.Applied{}, err
			}
			if unchanged(doc.Value(), held) {
				id, _ := doc.Get("id")
				fmt.Printf("%s %s (%s) unchanged.\n", k.Title, name, id)
				return resource.Applied{}, nil
			}
			return resource.Applied{}, fmt.Errorf("%s file '%s' holds the masked secret `get` and `export` write. Put the secret in, or remove it for none.", k.Title, file)
		}
		var saved struct{ ID, Name string }
		resp, err := k.upsert(ctx, c, resource.Body(resource.Strip(doc, ReadOnly...).Value()))
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
