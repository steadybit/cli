// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package badge_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/steadybit/cli/v6/internal/badge"
	"github.com/steadybit/cli/v6/internal/output"
	"github.com/steadybit/cli/v6/internal/platformtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

var svg = platformtest.Reply{Body: "<svg/>", Headers: map[string]string{"Content-Type": "image/svg+xml;charset=UTF-8"}}

func platformWithExperiment(t *testing.T) *platformtest.Platform {
	p := platformtest.New(t)
	p.Reply("GET /api/license", platformtest.Reply{JSON: map[string]any{"tenantKey": "demo"}})
	p.Reply("GET /api/experiments/ADM-1", platformtest.Reply{JSON: map[string]any{"key": "ADM-1", "name": "Shop", "team": "ADM"}})
	p.Reply("GET /api/experiments/ADM-1/badge.svg", svg)
	return p
}

func TestPrintsTheBadgeOfAnExperimentAsMarkdown(t *testing.T) {
	p := platformWithExperiment(t)

	out, err := platformtest.Stdout(t, func() error { return badge.Print(ctx, p.Client, badge.Options{Key: "ADM-1"}) })

	require.NoError(t, err)
	assert.Equal(t, "[![ADM-1]("+p.URL+"/api/experiments/ADM-1/badge.svg?tenantKey=demo)]("+p.URL+"/experiments/edit/ADM-1?team=ADM&tenant=demo)\n", out)
	checked := p.Requests("GET /api/experiments/ADM-1/badge.svg")[0]
	assert.Equal(t, []string{"demo"}, checked.Query["tenantKey"])
	// Fetched as a README does: with the token, the platform ignores a wrong tenant key.
	assert.Empty(t, checked.Header.Get("Authorization"))
	assert.NotEmpty(t, p.Requests("GET /api/experiments/ADM-1")[0].Header.Get("Authorization"))
}

func TestPrintsTheBadgeAsHTMLOrURL(t *testing.T) {
	p := platformWithExperiment(t)

	html, err := platformtest.Stdout(t, func() error {
		return badge.Print(ctx, p.Client, badge.Options{Key: "ADM-1", Tenant: "shop", Scale: 2, Format: "html"})
	})
	require.NoError(t, err)
	url, err := platformtest.Stdout(t, func() error { return badge.Print(ctx, p.Client, badge.Options{Key: "ADM-1", Format: "url"}) })
	require.NoError(t, err)

	assert.Equal(t, `<a href="`+p.URL+`/experiments/edit/ADM-1?team=ADM&amp;tenant=shop"><img alt="ADM-1" src="`+p.URL+`/api/experiments/ADM-1/badge.svg?scale=2&amp;tenantKey=shop"></a>`+"\n", html)
	assert.Equal(t, p.URL+"/api/experiments/ADM-1/badge.svg?tenantKey=demo\n", url)
	// A given tenant key is taken as it is; only the other run reads the license.
	assert.Len(t, p.Requests("GET /api/license"), 1)
}

func TestPrintsTheBadgeOfATag(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/badges/linked-badge.svg", svg)

	out, err := platformtest.Stdout(t, func() error {
		return badge.Print(ctx, p.Client, badge.Options{Tag: "INCIDENT-100", CreateCaption: "Create one (now)", Tenant: "demo"})
	})

	require.NoError(t, err)
	assert.Equal(t, "[![INCIDENT-100]("+p.URL+"/api/badges/linked-badge.svg?createCaption=Create%20one%20%28now%29&tag=INCIDENT-100&tenantKey=demo)]("+
		p.URL+"/api/badges/link?tag=INCIDENT-100&tenantKey=demo)\n", out)
	assert.Equal(t, []string{"Create one (now)"}, p.Requests("GET /api/badges/linked-badge.svg")[0].Query["createCaption"])
}

func TestPrintsTheBadgeAsJSON(t *testing.T) {
	p := platformWithExperiment(t)
	output.JQ = ".image"
	t.Cleanup(func() { output.JQ = "" })

	out, err := platformtest.Stdout(t, func() error { return badge.Print(ctx, p.Client, badge.Options{Key: "ADM-1"}) })

	require.NoError(t, err)
	assert.Equal(t, p.URL+"/api/experiments/ADM-1/badge.svg?tenantKey=demo\n", out)
}

func TestReportsAWrongTenantOrExperiment(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/experiments/ADM-1", platformtest.Reply{JSON: map[string]any{"key": "ADM-1", "team": "ADM"}})
	p.Reply("GET /api/experiments/ADM-2", platformtest.Reply{Status: http.StatusNotFound})
	p.Reply("GET /api/experiments/ADM-1/badge.svg", platformtest.Reply{Status: http.StatusBadRequest,
		Body: `{"type":"https://steadybit.com/problems/missing-tenant-exception","title":"A tenant must be set","status":400}`})

	assert.EqualError(t, badge.Print(ctx, p.Client, badge.Options{Key: "ADM-1", Tenant: "nosuch"}), "Tenant nosuch not found.")
	assert.EqualError(t, badge.Print(ctx, p.Client, badge.Options{Key: "ADM-2", Tenant: "demo"}), "Experiment ADM-2 not found.")
}

func TestFindingTheTenantNeedsAnAdminToken(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/license", platformtest.Reply{Status: http.StatusForbidden})

	assert.EqualError(t, badge.Print(ctx, p.Client, badge.Options{Key: "ADM-1"}),
		"Finding the tenant key needs an admin access token. Pass it with --tenant: it is the tenant= of a platform URL.")
}

func TestRefusals(t *testing.T) {
	p := platformtest.New(t)

	assert.EqualError(t, badge.Print(ctx, p.Client, badge.Options{}), "Either --key or --tag must be specified.")
	assert.EqualError(t, badge.Print(ctx, p.Client, badge.Options{Key: "ADM-1", CreateCaption: "x"}), "--create-caption only applies to a badge for --tag.")
	assert.EqualError(t, badge.Print(ctx, p.Client, badge.Options{Key: "ADM-1", Format: "svg"}), `Unsupported badge format 'svg'. Use "markdown", "html" or "url".`)
}
