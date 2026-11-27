// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package accesstoken implements the `access-token` commands, on the v2 endpoints only.
package accesstoken

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/steadybit/cli/v6/api"
	"github.com/steadybit/cli/v6/internal/jsyaml"
	"github.com/steadybit/cli/v6/internal/platform"
	"github.com/steadybit/cli/v6/internal/resource"
	"github.com/steadybit/cli/v6/internal/table"
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

type RotateOptions struct {
	ID, ExpiresAt, GracePeriod string
	Output                     string
	Yes                        bool
}

// Rotate gives a token a new secret in place: its id, name, type, teams and expiry stay
// unless --expires-at says otherwise.
func Rotate(ctx context.Context, c *platform.Client, o RotateOptions) error {
	params := api.RotateAccessTokenParams{}
	if o.GracePeriod != "" {
		params.GracePeriod = &o.GracePeriod
	}
	if o.ExpiresAt != "" {
		var expiresAt api.RotateAccessTokenParamsExpiresAt
		if o.ExpiresAt == "never" {
			if err := expiresAt.FromRotateAccessTokenParamsExpiresAt1(api.Never); err != nil {
				return err
			}
		} else {
			at, err := resource.Time("expires-at", o.ExpiresAt)
			if err != nil {
				return err
			}
			if err := expiresAt.FromRotateAccessTokenParamsExpiresAt0(*at); err != nil {
				return err
			}
		}
		params.ExpiresAt = &expiresAt
	}
	question := fmt.Sprintf("Rotate access token %s? Its current secret stops working.", o.ID)
	if d, err := time.ParseDuration(o.GracePeriod); o.GracePeriod != "" && (err != nil || d > 0) {
		question = fmt.Sprintf("Rotate access token %s? Its current secret keeps working for %s.", o.ID, o.GracePeriod)
	}
	if ok, err := resource.Confirmed(o.Yes, question); !ok || err != nil {
		return err
	}
	var rotated struct {
		created
		PreviousSecretExpiresAt *string `json:"previousSecretExpiresAt"`
	}
	resp, err := c.RotateAccessToken(ctx, o.ID, &params)
	if _, err := platform.Decode(resp, err, &rotated); err != nil {
		// A platform without the rotate endpoint answers 404 too; the token list tells the two apart.
		if platform.IsStatus(err, http.StatusNotFound) {
			if token, lookupErr := listedToken(ctx, c, o.ID); lookupErr == nil && token != nil {
				return fmt.Errorf("The platform can't rotate access tokens yet. Update it, or use access-token recreate until then.")
			}
		}
		return notFoundOr(err, o.ID, "Failed to rotate access token %s")
	}
	var until any
	previous := "The previous secret stopped working."
	if rotated.PreviousSecretExpiresAt != nil {
		until = *rotated.PreviousSecretExpiresAt
		previous = fmt.Sprintf("The previous secret works until %s.", until)
	}
	if o.Output != "" {
		m := jsyaml.NewMap()
		m.Set("id", rotated.ID)
		m.Set("token", rotated.Token)
		m.Set("previousSecretExpiresAt", until)
		return resource.OutputValue(m, "", o.Output)
	}
	fmt.Printf("Access token %s rotated. %s Store the new one now, it cannot be shown again:\n%s\n", o.ID, previous, rotated.Token)
	return nil
}

type RecreateOptions struct {
	ID, ExpiresAt string
	Output        string
	Yes           bool
}

func Recreate(ctx context.Context, c *platform.Client, o RecreateOptions) error {
	var expiresAt *time.Time
	switch o.ExpiresAt {
	case "never":
	case "":
		// The platform gives a recreated token no expiry unless told one, which would
		// quietly turn a token meant to expire into one that never does.
		kept, err := currentExpiry(ctx, c, o.ID)
		if err != nil {
			return err
		}
		expiresAt = kept
	default:
		var err error
		if expiresAt, err = resource.Time("expires-at", o.ExpiresAt); err != nil {
			return err
		}
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

type listed struct {
	ID        string  `json:"id"`
	ExpiresAt *string `json:"expiresAt"`
}

// listedToken is the token with the id as the token list shows it, or nil when there is none.
func listedToken(ctx context.Context, c *platform.Client, id string) (*listed, error) {
	raw, err := platform.AllPagesRaw(func(page, size int32) (*http.Response, error) {
		return c.GetAccessTokens1(ctx, &api.GetAccessTokens1Params{PageRequest: api.PageRequestAO{Page: &page, Size: &size}})
	})
	if err != nil {
		return nil, platform.Failed(err, "Failed to get access token %s", id)
	}
	var tokens []listed
	if err := resource.DecodeEach(raw, &tokens); err != nil {
		return nil, err
	}
	for _, token := range tokens {
		if token.ID == id {
			return &token, nil
		}
	}
	return nil, nil
}

// currentExpiry is when the token expires now, for the new one to expire then too, or
// nil when it never does. One that has already expired needs an explicit new expiry.
func currentExpiry(ctx context.Context, c *platform.Client, id string) (*time.Time, error) {
	token, err := listedToken(ctx, c, id)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, fmt.Errorf("Access token %s not found.", id)
	}
	if token.ExpiresAt == nil || *token.ExpiresAt == "" {
		return nil, nil
	}
	expiry, err := time.Parse(time.RFC3339Nano, *token.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("Access token %s expires at '%s', which is not a time the CLI can read. Pass --expires-at.", id, *token.ExpiresAt)
	}
	if !expiry.After(time.Now()) {
		return nil, fmt.Errorf("Access token %s expired at %s. Pass --expires-at for the new one, or --expires-at never.", id, *token.ExpiresAt)
	}
	return &expiry, nil
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
