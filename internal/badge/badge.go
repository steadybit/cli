// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package badge implements the `experiment badge` command: a status badge to embed in a
// README, which is served without an access token.
package badge

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"strings"

	"github.com/steadybit/cli/v6/internal/experiment"
	"github.com/steadybit/cli/v6/internal/output"
	"github.com/steadybit/cli/v6/internal/platform"
	"github.com/steadybit/cli/v6/internal/resource"
)

type Options struct {
	Key, Tag, CreateCaption string
	Tenant                  string
	Scale                   int
	// "markdown", "html" or "url".
	Format string
	Type   string
}

type badge struct {
	Image    string `json:"image"`
	Link     string `json:"link"`
	Markdown string `json:"markdown"`
	HTML     string `json:"html"`
}

// Print writes the snippet that embeds the badge. The badge URLs carry the tenant key
// and never the access token: a README is read by anyone, and the platform serves
// badges to anyone who knows the tenant key.
func Print(ctx context.Context, c *platform.Client, o Options) error {
	if (o.Key == "") == (o.Tag == "") {
		return errors.New("Either --key or --tag must be specified.")
	}
	if o.CreateCaption != "" && o.Tag == "" {
		return errors.New("--create-caption only applies to a badge for --tag.")
	}
	if o.Scale < 0 {
		return errors.New("--scale cannot be negative.")
	}
	// Checked before any request: the badge takes up to three.
	if _, err := output.ResolveDatatype(o.Type, ""); err != nil {
		return err
	}
	format := o.Format
	if format != "" && resource.Machine(o.Type) {
		return errors.New("--format cannot be combined with -t or --jq, which print every format.")
	}
	if format == "" {
		format = "markdown"
	}
	if format != "markdown" && format != "html" && format != "url" {
		return fmt.Errorf("Unsupported badge format '%s'. Use \"markdown\", \"html\" or \"url\".", format)
	}
	tenant, err := tenantKey(ctx, c, o.Tenant)
	if err != nil {
		return err
	}

	var imagePath, linkPath, alt string
	image := url.Values{"tenantKey": {tenant}}
	if o.Scale > 0 {
		image.Set("scale", fmt.Sprint(o.Scale))
	}
	if o.Key != "" {
		// The badge of a key that does not exist is an image saying "not found", with 200.
		doc, err := experiment.Fetch(ctx, c, o.Key)
		if err != nil {
			return err
		}
		team, _ := doc.Get("team")
		imagePath = "/api/experiments/" + url.PathEscape(o.Key) + "/badge.svg?" + query(image)
		linkPath = "/experiments/edit/" + url.PathEscape(o.Key) + "?" + query(url.Values{"tenant": {tenant}, "team": {team}})
		alt = o.Key
	} else {
		image.Set("tag", o.Tag)
		if o.CreateCaption != "" {
			image.Set("createCaption", o.CreateCaption)
		}
		imagePath = "/api/badges/linked-badge.svg?" + query(image)
		linkPath = "/api/badges/link?" + query(url.Values{"tenantKey": {tenant}, "tag": {o.Tag}})
		alt = o.Tag
	}
	if err := check(ctx, c, imagePath, tenant); err != nil {
		return err
	}

	b := badge{Image: c.BaseURL + imagePath, Link: c.BaseURL + linkPath}
	b.Markdown = fmt.Sprintf("[![%s](%s)](%s)", markdownText(alt), b.Image, b.Link)
	b.HTML = fmt.Sprintf(`<a href="%s"><img alt="%s" src="%s"></a>`, html.EscapeString(b.Link), html.EscapeString(alt), html.EscapeString(b.Image))
	if resource.Machine(o.Type) {
		raw, _ := json.Marshal(b)
		return resource.PrintJSONValue(raw, o.Type)
	}
	switch format {
	case "html":
		fmt.Println(b.HTML)
	case "url":
		fmt.Println(b.Image)
	default:
		fmt.Println(b.Markdown)
	}
	return nil
}

// query encodes spaces as %20: a badge caption is shown as written, and not every
// Markdown renderer or server reads + as a space.
func query(v url.Values) string { return strings.ReplaceAll(v.Encode(), "+", "%20") }

func markdownText(s string) string {
	return strings.NewReplacer(`\`, `\\`, "[", `\[`, "]", `\]`).Replace(s)
}

// tenantKey is the one given, or the one the license names. The access token does not
// say which tenant it belongs to, and the license is the only other place that does.
// A given key is still compared with the license when it can be read: the badge of
// another tenant's experiment is an image saying "not found", with 200, which the check
// of the badge cannot tell from a real one.
func tenantKey(ctx context.Context, c *platform.Client, given string) (string, error) {
	var summary struct {
		TenantKey string `json:"tenantKey"`
	}
	resp, err := c.GetLicenseSummary(ctx)
	_, err = platform.Decode(resp, err, &summary)
	if given != "" {
		// Without an admin token the license cannot be read, and the given key is taken as it is.
		if err == nil && summary.TenantKey != "" && summary.TenantKey != given {
			return "", fmt.Errorf("The access token belongs to tenant %s, not %s: the badge would show \"not found\". Leave out --tenant, or use an access token of tenant %s.",
				summary.TenantKey, given, given)
		}
		return given, nil
	}
	if platform.IsStatus(err, http.StatusForbidden) {
		return "", errors.New("Finding the tenant key needs an admin access token. Pass it with --tenant: it is the tenant= of a platform URL.")
	}
	if err != nil {
		return "", platform.Failed(err, "Failed to find the tenant key")
	}
	if summary.TenantKey == "" {
		return "", errors.New("The platform did not name the tenant. Pass its key with --tenant: it is the tenant= of a platform URL.")
	}
	return summary.TenantKey, nil
}

// check fetches the badge as a README would show it, without the token, so that a
// wrong tenant key fails here and not as a broken image.
func check(ctx context.Context, c *platform.Client, path, tenant string) error {
	_, resp, err := platform.Read(c.GetAnonymously(ctx, path))
	var apiErr *platform.APIError
	if errors.As(err, &apiErr) && apiErr.Status == http.StatusBadRequest && strings.HasSuffix(apiErr.ProblemType(), "/missing-tenant-exception") {
		return fmt.Errorf("Tenant %s not found.", tenant)
	}
	if err != nil {
		return platform.Failed(err, "Failed to get the badge")
	}
	if kind := resp.Header.Get("Content-Type"); !strings.HasPrefix(kind, "image/svg+xml") {
		return fmt.Errorf("The platform sent %s instead of a badge image.", kind)
	}
	return nil
}
