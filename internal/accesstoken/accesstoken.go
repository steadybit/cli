// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package accesstoken implements the `access-token` commands, on the v2 endpoints only.
package accesstoken

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/steadybit/cli/api"
	"github.com/steadybit/cli/internal/jsyaml"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/resource"
	"github.com/steadybit/cli/internal/table"
)

func notFoundOr(err error, id, format string) error {
	if platform.IsStatus(err, http.StatusNotFound) {
		return fmt.Errorf("Access token %s not found.", id)
	}
	return platform.Failed(err, format, id)
}

type ListOptions struct {
	Name, CreatedBy, Type string
	Teams                 []string
	Expired               *bool
	// Output is -t for the other listings; --type already names the token type here.
	Output string
}

func List(ctx context.Context, c *platform.Client, o ListOptions) error {
	params := api.GetAccessTokens1Params{Teams: resource.Optional(o.Teams), Expired: o.Expired}
	if o.Name != "" {
		params.Name = &o.Name
	}
	if o.CreatedBy != "" {
		params.CreatedBy = &o.CreatedBy
	}
	if o.Type != "" {
		t := api.GetAccessTokens1ParamsType(strings.ToUpper(o.Type))
		if !t.Valid() {
			return fmt.Errorf("--type must be ADMIN, TEAM or WILDCARD, not '%s'.", o.Type)
		}
		params.Type = &t
	}
	type summary struct {
		ID, Name, Type string
		Teams          []string
		ExpiresAt      *string `json:"expiresAt"`
		LastUsed       *string `json:"lastUsed"`
	}
	raw, err := platform.AllPagesRaw(func(page, size int32) (*http.Response, error) {
		params.PageRequest = api.PageRequestAO{Page: &page, Size: &size}
		return c.GetAccessTokens1(ctx, &params)
	})
	if err != nil {
		return platform.Failed(err, "Failed to get the access tokens")
	}
	if resource.Machine(o.Output) {
		return resource.List(raw, o.Output, nil)
	}
	var tokens []summary
	if err := resource.DecodeEach(raw, &tokens); err != nil {
		return err
	}
	if len(tokens) == 0 {
		fmt.Println("No access tokens found.")
		return nil
	}
	t := table.New(
		table.Column{Name: "id", Title: "Id", Alignment: table.Left},
		table.Column{Name: "name", Title: "Name", Alignment: table.Left},
		table.Column{Name: "type", Title: "Type", Alignment: table.Left},
		table.Column{Name: "teams", Title: "Teams", Alignment: table.Left},
		table.Column{Name: "expiresAt", Title: "Expires", Alignment: table.Left},
		table.Column{Name: "lastUsed", Title: "Last used", Alignment: table.Left},
	)
	or := func(s *string, fallback string) string {
		if s == nil {
			return fallback
		}
		return *s
	}
	for _, token := range tokens {
		t.AddRow(table.Default, table.Cell("id", token.ID), table.Cell("name", token.Name), table.Cell("type", token.Type),
			table.Cell("teams", strings.Join(token.Teams, ", ")), table.Cell("expiresAt", or(token.ExpiresAt, "never")), table.Cell("lastUsed", or(token.LastUsed, "never")))
	}
	t.Print()
	return nil
}

type created struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

// printToken shows a new token, the only time it can be seen. With a type, only the id
// and token are printed, for a pipeline to read.
func printToken(token created, datatype, headline string) error {
	if datatype != "" {
		m := jsyaml.NewMap()
		m.Set("id", token.ID)
		m.Set("token", token.Token)
		return resource.OutputValue(m, "", datatype)
	}
	fmt.Printf("%s Store it now, it cannot be shown again:\n%s\n", headline, token.Token)
	return nil
}

type CreateOptions struct {
	Name, Type, ExpiresAt string
	Teams                 []string
	Output                string
}

func Create(ctx context.Context, c *platform.Client, o CreateOptions) error {
	kind := api.CreateAccessTokenRequestV2AOType(strings.ToUpper(o.Type))
	if !kind.Valid() {
		return fmt.Errorf("--type must be ADMIN, TEAM or WILDCARD, not '%s'.", o.Type)
	}
	expiresAt, err := resource.Time("expires-at", o.ExpiresAt)
	if err != nil {
		return err
	}
	var token created
	resp, err := c.CreateAccessToken1(ctx, api.CreateAccessTokenRequestV2AO{Name: o.Name, Type: kind, Teams: resource.Optional(o.Teams), ExpiresAt: expiresAt})
	if _, err := platform.Decode(resp, err, &token); err != nil {
		return platform.Failed(err, "Failed to create access token %s", o.Name)
	}
	return printToken(token, o.Output, fmt.Sprintf("Access token %s (%s) created.", o.Name, token.ID))
}

type RecreateOptions struct {
	ID, ExpiresAt string
	Output        string
	Yes           bool
}

func Recreate(ctx context.Context, c *platform.Client, o RecreateOptions) error {
	expiresAt, err := resource.Time("expires-at", o.ExpiresAt)
	if err != nil {
		return err
	}
	if ok, err := resource.Confirmed(o.Yes, fmt.Sprintf("Recreate access token %s? The current token stops working.", o.ID)); !ok || err != nil {
		return err
	}
	var token created
	resp, err := c.RecreateAccessToken(ctx, o.ID, api.RecreateAccessTokenRequestV2AO{ExpiresAt: expiresAt})
	if _, err := platform.Decode(resp, err, &token); err != nil {
		return notFoundOr(err, o.ID, "Failed to recreate access token %s")
	}
	return printToken(token, o.Output, fmt.Sprintf("Access token %s recreated as %s.", o.ID, token.ID))
}

type DeleteOptions struct {
	ID  string
	Yes bool
}

func Delete(ctx context.Context, c *platform.Client, o DeleteOptions) error {
	if ok, err := resource.Confirmed(o.Yes, fmt.Sprintf("Delete access token %s? Everything using it loses access.", o.ID)); !ok || err != nil {
		return err
	}
	if _, _, err := platform.Read(c.DeleteAccessToken1(ctx, o.ID)); err != nil {
		return notFoundOr(err, o.ID, "Failed to delete access token %s")
	}
	fmt.Printf("Access token %s deleted.\n", o.ID)
	return nil
}
