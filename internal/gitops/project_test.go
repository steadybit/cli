// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package gitops_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/steadybit/cli/internal/gitops"
	"github.com/steadybit/cli/internal/platformtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	serviceID = "019cd80d-a4c9-775b-bdf8-2672a280ce7c"
	profileID = "019eacd7-fb2c-733a-bed5-99a935323db5"
)

func fakeProject(t *testing.T) *platformtest.Platform {
	p := platformtest.New(t)
	p.Reply("GET /api/experiments", platformtest.Reply{JSON: map[string]any{"experiments": []any{map[string]any{"key": "ADM-1"}}}})
	p.Reply("GET /api/experiments/ADM-1", platformtest.Reply{Body: `{"key":"ADM-1","version":2,"name":"Survives","team":"ADM","created":"c","lanes":[]}`})
	p.Reply("GET /api/experiments/schedules/v2", platformtest.Reply{Body: `[{"id":"01951394-727f-76a0-8675-c7519ebd0ff5","experimentKey":"ADM-1","cron":"0 0 9 ? * *","editedBy":{},"lastUpdated":"x"}]`})
	p.Reply("GET /api/experiments/schedules/*", platformtest.Reply{Body: `{"id":"01951394-727f-76a0-8675-c7519ebd0ff5","experimentKey":"ADM-1","cron":"0 0 9 ? * *","editedBy":{},"lastUpdated":"x"}`})
	p.Reply("GET /api/services", platformtest.Reply{JSON: map[string]any{"items": []any{map[string]any{"id": serviceID}}}})
	p.Reply("GET /api/services/"+serviceID, platformtest.Reply{Body: `{"id":"` + serviceID + `","name":"Checkout / API","team":"ADM","serviceProfile":"Shop","version":1}`})
	p.Reply("GET /api/services/profiles", platformtest.Reply{JSON: map[string]any{"items": []any{map[string]any{"id": profileID, "name": "Shop"}}}})
	p.Reply("GET /api/services/profiles/"+profileID, platformtest.Reply{Body: `{"id":"` + profileID + `","name":"Shop","origin":"CUSTOM","templates":[],"defaultProfile":false}`})
	return p
}

func TestExportWritesAProjectThatMatchesThePlatform(t *testing.T) {
	p := fakeProject(t)
	dir := t.TempDir()

	out, err := platformtest.Stdout(t, func() error { return gitops.Export(ctx, p.Client, gitops.ExportOptions{Directory: dir, Team: "ADM"}) })

	require.NoError(t, err)
	assert.Equal(t, "Exported team ADM to "+dir+": 1 experiments, 1 schedules, 1 services, 1 service profiles.\n", out)
	for file, content := range map[string]string{
		"experiments/adm-1.yaml":        "key: ADM-1\nname: Survives\nteam: ADM\nlanes: []\n",
		"schedules/adm-1-9ebd0ff5.yaml": "id: 01951394-727f-76a0-8675-c7519ebd0ff5\nexperimentKey: ADM-1\ncron: 0 0 9 ? * *\n",
		"services/checkout-api.yaml":    "id: " + serviceID + "\nname: Checkout / API\nteam: ADM\nserviceProfile: Shop\n",
		"service-profiles/shop.yaml":    "id: " + profileID + "\nname: Shop\norigin: CUSTOM\ntemplates: []\n",
	} {
		written, err := os.ReadFile(filepath.Join(dir, file))
		require.NoError(t, err, file)
		assert.Equal(t, content, string(written), file)
	}

	out, err = platformtest.Stdout(t, func() error { return gitops.DiffProject(ctx, p.Client, dir) })
	require.NoError(t, err)
	assert.Equal(t, 4, strings.Count(out, "match the platform"))
}

func TestAProjectEditedLocallyDrifts(t *testing.T) {
	p := fakeProject(t)
	dir := t.TempDir()
	_, err := platformtest.Stdout(t, func() error { return gitops.Export(ctx, p.Client, gitops.ExportOptions{Directory: dir, Team: "ADM"}) })
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "experiments", "adm-1.yaml"), []byte("key: ADM-1\nname: Renamed\nteam: ADM\nlanes: []\n"), 0o644))

	out, err := platformtest.Stdout(t, func() error { return gitops.DiffProject(ctx, p.Client, dir) })
	assert.ErrorIs(t, err, gitops.ErrDifferent)
	assert.Contains(t, out, "-name: Survives\n+name: Renamed\n")

	out, err = platformtest.Stdout(t, func() error {
		return gitops.ApplyProject(ctx, p.Client, gitops.ApplyOptions{Directory: dir, DryRun: true})
	})
	require.NoError(t, err)
	assert.Contains(t, out, "would update experiment ADM-1 (2 lines changed).")
	assert.Empty(t, p.Requests("POST /api/experiments/ADM-1"), "a dry run changes nothing")
}

func TestApplyNeedsAProject(t *testing.T) {
	err := gitops.ApplyProject(ctx, nil, gitops.ApplyOptions{Directory: t.TempDir()})

	assert.ErrorContains(t, err, "holds none of experiments/, schedules/, services/ or service-profiles/.")
}
