// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package license_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/steadybit/cli/v6/internal/license"
	"github.com/steadybit/cli/v6/internal/output"
	"github.com/steadybit/cli/v6/internal/platformtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

const summary = `{"license":{"id":1,"licenseType":"ENTERPRISE","orderNumber":"1234-2","validFrom":"2026-01-22","validTo":"2099-01-01"},
 "expires":"2099-01-01T00:00:00Z","tenantKey":"demo","features":[
 {"name":"TEMPLATES","type":"SIMPLE","usage":0},
 {"name":"SERVICES","type":"HARD_LIMIT","usage":60,"hardLimit":50},
 {"name":"AUDIT_LOG","type":"SIMPLE","usage":0},
 {"name":"ENVIRONMENT_SIZE","type":"SOFT_LIMIT","usage":45,"softLimit":100},
 {"name":"USER_SIZE","type":"HARD_LIMIT","usage":22}]}`

func TestShowsTheLicenseAndItsLimits(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/license", platformtest.Reply{Body: summary})

	out, err := platformtest.Stdout(t, func() error { return license.Show(ctx, p.Client, license.ShowOptions{}) })

	require.NoError(t, err)
	assert.Equal(t, `Enterprise license 1234-2 of tenant demo, valid from 2026-01-22 to 2099-01-01.
┌──────────────────┬──────┬────────────┐
│ Limit            │ Used │   Licensed │
├──────────────────┼──────┼────────────┤
│ ENVIRONMENT_SIZE │   45 │ 100 (soft) │
│ SERVICES         │   60 │         50 │
│ USER_SIZE        │   22 │  unlimited │
└──────────────────┴──────┴────────────┘
Included: AUDIT_LOG, TEMPLATES
`, out)
}

func TestSaysWhenTheLicenseHasExpired(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/license", platformtest.Reply{Body: `{"license":{"licenseType":"TRIAL","orderNumber":"7","validFrom":"2020-01-01","validTo":"2020-02-01"},"expires":"2020-02-01T00:00:00Z","tenantKey":"demo","features":[]}`})

	out, err := platformtest.Stdout(t, func() error { return license.Show(ctx, p.Client, license.ShowOptions{}) })

	require.NoError(t, err)
	assert.Equal(t, "Trial license 7 of tenant demo, valid from 2020-01-01 to 2020-02-01, expired.\n", out)
}

func TestWarnsOfALicenseAboutToExpire(t *testing.T) {
	p := platformtest.New(t)
	expires := time.Now().Add(10 * 24 * time.Hour).UTC().Format(time.RFC3339)
	p.Reply("GET /api/license", platformtest.Reply{Body: `{"license":{"licenseType":"PROFESSIONAL","orderNumber":"8","validFrom":"2020-01-01","validTo":"` + expires[:10] + `"},"expires":"` + expires + `","tenantKey":"demo"}`})

	out, err := platformtest.Stdout(t, func() error { return license.Show(ctx, p.Client, license.ShowOptions{}) })

	require.NoError(t, err)
	assert.Equal(t, "Professional license 8 of tenant demo, valid from 2020-01-01 to "+expires[:10]+", expires in 10 days.\n", out)
}

func TestSaysWhenThereIsNoLicense(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/license", platformtest.Reply{Body: `{"license":null,"tenantKey":"demo"}`})

	out, err := platformtest.Stdout(t, func() error { return license.Show(ctx, p.Client, license.ShowOptions{}) })

	require.NoError(t, err)
	assert.Equal(t, "The tenant has no license.\n", out)
}

func TestPrintsTheLicenseAsThePlatformSendsIt(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/license", platformtest.Reply{Body: `{"tenantKey":"demo","expires":"2099-01-01T00:00:00Z"}`})

	yaml, err := platformtest.Stdout(t, func() error { return license.Show(ctx, p.Client, license.ShowOptions{Type: "yaml"}) })
	require.NoError(t, err)
	output.JQ = ".expires"
	t.Cleanup(func() { output.JQ = "" })
	jq, err := platformtest.Stdout(t, func() error { return license.Show(ctx, p.Client, license.ShowOptions{}) })
	require.NoError(t, err)

	assert.Equal(t, "tenantKey: demo\nexpires: '2099-01-01T00:00:00Z'\n", yaml)
	assert.Equal(t, "2099-01-01T00:00:00Z\n", jq)
}

func TestDownloadsTheReportUnderItsName(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/license/report", platformtest.Reply{Body: "PK-zip", Headers: map[string]string{
		"Content-Disposition": `attachment; filename="../license-usage-reports-demo.zip"`,
	}})
	t.Chdir(t.TempDir())

	out, err := platformtest.Stdout(t, func() error { return license.Report(ctx, p.Client, license.ReportOptions{}) })
	require.NoError(t, err)
	named, _ := os.ReadFile("license-usage-reports-demo.zip")
	file := filepath.Join("reports", "usage.zip")
	_, err = platformtest.Stdout(t, func() error { return license.Report(ctx, p.Client, license.ReportOptions{Output: file}) })
	require.NoError(t, err)
	given, _ := os.ReadFile(file)

	assert.Equal(t, "License report written to license-usage-reports-demo.zip.\n", out)
	assert.Equal(t, "PK-zip", string(named))
	assert.Equal(t, "PK-zip", string(given))
}

func TestDownloadsTheReportWithoutAName(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/license/report", platformtest.Reply{Body: "PK-zip"})
	t.Chdir(t.TempDir())

	out, err := platformtest.Stdout(t, func() error { return license.Report(ctx, p.Client, license.ReportOptions{}) })

	require.NoError(t, err)
	assert.Equal(t, "License report written to license-report.zip.\n", out)
	assert.FileExists(t, "license-report.zip")
}

func TestNeedsAnAdminToken(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/license", platformtest.Reply{Status: http.StatusForbidden})
	p.Reply("GET /api/license/report", platformtest.Reply{Status: http.StatusForbidden})
	t.Chdir(t.TempDir())

	assert.EqualError(t, license.Show(ctx, p.Client, license.ShowOptions{}), "The license needs an admin access token.")
	assert.EqualError(t, license.Report(ctx, p.Client, license.ReportOptions{}), "The license needs an admin access token.")
	assert.NoFileExists(t, "license-report.zip")
}
