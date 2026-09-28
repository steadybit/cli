// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package gitops_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/steadybit/cli/v6/internal/gitops"
	"github.com/steadybit/cli/v6/internal/output"
	"github.com/steadybit/cli/v6/internal/platformtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func doc(t *testing.T, text string) *output.Document {
	t.Helper()
	d, err := output.ParseDocument([]byte(text))
	require.NoError(t, err)
	return d
}

// What the platform fills in with defaults is not a difference; what it holds besides is.
func TestIgnoresPlatformDefaultsButNotRealValues(t *testing.T) {
	local := doc(t, "name: x\nlanes:\n  - steps:\n      - type: wait\n")
	remote := doc(t, `{"name":"x","sharedTeams":[],"properties":{},"lanes":[{"steps":[{"type":"wait","ignoreFailure":false}]}]}`)

	diff, err := gitops.Diff("f.yml", local.Value(), remote.Value())
	require.NoError(t, err)
	assert.Empty(t, diff)

	remote = doc(t, `{"name":"x","hypothesis":"set in the UI","lanes":[{"steps":[{"type":"wait","ignoreFailure":true}]}]}`)
	diff, err = gitops.Diff("f.yml", local.Value(), remote.Value())
	require.NoError(t, err)
	assert.Equal(t, `--- platform
+++ f.yml
@@ -2,5 +2,3 @@
 lanes:
   - steps:
       - type: wait
-        ignoreFailure: true
-hypothesis: set in the UI
`, diff)
}

func TestReportsFieldsTheFileChanges(t *testing.T) {
	local := doc(t, "name: new name\nteam: ADM\n")
	remote := doc(t, `{"name":"old name","team":"ADM"}`)

	diff, err := gitops.Diff("f.yml", local.Value(), remote.Value())

	require.NoError(t, err)
	assert.Contains(t, diff, "-name: old name\n+name: new name\n")
	assert.Equal(t, 2, gitops.Changes(diff))
}

func write(t *testing.T, name, content string) string {
	t.Helper()
	file := filepath.Join(t.TempDir(), name)
	require.NoError(t, os.WriteFile(file, []byte(content), 0o644))
	return file
}

const stored = `{"key":"TST-1","version":7,"name":"x","team":"TST","created":"c","createdBy":{},"edited":"e","editedBy":{},"lanes":[]}`

func TestDiffExitsWithDifferentWhenFilesDrift(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/experiments/TST-1", platformtest.Reply{Body: stored})
	same := write(t, "same.yml", "key: TST-1\nname: x\nteam: TST\nlanes: []\n")
	changed := write(t, "changed.yml", "key: TST-1\nname: y\nteam: TST\nlanes: []\n")
	fresh := write(t, "new.yml", "name: z\n")

	out, err := platformtest.Stdout(t, func() error { return gitops.DiffFiles(ctx, p.Client, gitops.Experiment, []string{same}, false) })
	require.NoError(t, err)
	assert.Equal(t, "1 experiment file(s) match the platform.\n", out)

	out, err = platformtest.Stdout(t, func() error {
		return gitops.DiffFiles(ctx, p.Client, gitops.Experiment, []string{same, changed, fresh}, false)
	})
	assert.ErrorIs(t, err, gitops.ErrDifferent)
	assert.Contains(t, out, "-name: x\n+name: y\n")
	assert.Contains(t, out, fresh+" would create a new experiment.\n")
	assert.Contains(t, out, "2 of 3 experiment file(s) differ from the platform.\n")
}

func TestDryRunFindsExperimentsByExternalID(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/experiments", platformtest.Reply{JSON: map[string]any{"experiments": []any{map[string]any{"key": "TST-1"}}}})
	p.Reply("GET /api/experiments/TST-1", platformtest.Reply{Body: stored})
	file := write(t, "e.yml", "externalId: shop-latency\nname: renamed\n")

	out, err := platformtest.Stdout(t, func() error { return gitops.DryRun(ctx, p.Client, gitops.Experiment, []string{file}, false) })

	require.NoError(t, err)
	assert.Equal(t, file+" would update experiment TST-1 (4 lines changed).\n", out)
	assert.Equal(t, []string{"shop-latency"}, p.Requests("GET /api/experiments")[0].Query["externalId"])
}

func TestAScheduleWithAnUnknownIdWouldBeCreated(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/experiments/schedules/gone", platformtest.Reply{Status: http.StatusNotFound})
	file := write(t, "s.yml", "id: gone\nexperimentKey: TST-1\n")

	out, err := platformtest.Stdout(t, func() error { return gitops.DryRun(ctx, p.Client, gitops.Schedule, []string{file}, false) })

	require.NoError(t, err)
	assert.Equal(t, file+" would create a new experiment schedule.\n", out)
}

func TestProfilesAreFoundByName(t *testing.T) {
	p := platformtest.New(t)
	id := "019eacd7-fb2c-733a-bed5-99a935323db5"
	p.Reply("GET /api/services/profiles", platformtest.Reply{JSON: map[string]any{"items": []any{
		map[string]any{"id": "019eacd7-fb2c-733a-bed5-99a935323db6", "name": "High Redundancy 2"},
		map[string]any{"id": id, "name": "High Redundancy"},
	}}})
	p.Reply("GET /api/services/profiles/"+id, platformtest.Reply{Body: `{"id":"` + id + `","name":"High Redundancy","origin":"CUSTOM","templates":[],"version":3,"defaultProfile":false}`})
	file := write(t, "p.yml", "name: High Redundancy\norigin: CUSTOM\ntemplates: []\n")

	out, err := platformtest.Stdout(t, func() error { return gitops.DryRun(ctx, p.Client, gitops.ServiceProfile, []string{file}, false) })

	require.NoError(t, err)
	assert.Equal(t, file+" matches service profile "+id+".\n", out)
}

// A schedule written by hand, right after it was applied: the platform added its default
// and gave the start time its own spelling. Neither is a difference.
func TestAHandWrittenScheduleMatchesRightAfterApply(t *testing.T) {
	p := platformtest.New(t)
	stored := `{"id":"s1","experimentKey":"TST-1","startAt":"2030-07-01T09:00:00Z","enabled":false,"allowParallel":true}`
	p.Reply("GET /api/experiments/schedules/s1", platformtest.Reply{Body: stored})
	file := write(t, "s.yml", "id: s1\nexperimentKey: TST-1\nstartAt: 2030-07-01T09:00:00Z\nenabled: false\n")

	out, err := platformtest.Stdout(t, func() error { return gitops.DiffFiles(ctx, p.Client, gitops.Schedule, []string{file}, false) })

	require.NoError(t, err)
	assert.Equal(t, "1 experiment schedule file(s) match the platform.\n", out)
}

// A default the file leaves out is a difference once it was changed on the platform:
// applying the file would reset it.
func TestAChangedDefaultIsADifference(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/experiments/schedules/s1", platformtest.Reply{Body: `{"id":"s1","experimentKey":"TST-1","startAt":"2030-07-01T10:00:00Z","allowParallel":false}`})
	file := write(t, "s.yml", "id: s1\nexperimentKey: TST-1\nstartAt: 2030-07-01T09:00:00Z\n")

	out, err := platformtest.Stdout(t, func() error { return gitops.DiffFiles(ctx, p.Client, gitops.Schedule, []string{file}, false) })

	assert.ErrorIs(t, err, gitops.ErrDifferent)
	assert.Contains(t, out, "-startAt: '2030-07-01T10:00:00Z'\n")
	assert.Contains(t, out, "+startAt: 2030-07-01T09:00:00.000Z\n")
	assert.Contains(t, out, "-allowParallel: false\n")
}

func TestAHandWrittenServiceAndProfileMatchRightAfterApply(t *testing.T) {
	p := platformtest.New(t)
	id := "019eacd7-fb2c-733a-bed5-99a935323db5"
	p.Reply("GET /api/services/"+id, platformtest.Reply{Body: `{"id":"` + id + `","name":"svc","query":"a=b","validations":[],"logoId":"service","logoColor":"blue"}`})
	p.Reply("GET /api/services/profiles/"+id, platformtest.Reply{Body: `{"id":"` + id + `","name":"p","templates":[],"origin":"CUSTOM"}`})
	service := write(t, "svc.yml", "id: "+id+"\nname: svc\nquery: a=b\nvalidations: []\n")
	profile := write(t, "p.yml", "id: "+id+"\nname: p\ntemplates: []\n")

	_, err := platformtest.Stdout(t, func() error { return gitops.DiffFiles(ctx, p.Client, gitops.Service, []string{service}, false) })
	require.NoError(t, err)
	_, err = platformtest.Stdout(t, func() error { return gitops.DiffFiles(ctx, p.Client, gitops.ServiceProfile, []string{profile}, false) })
	require.NoError(t, err)
}
