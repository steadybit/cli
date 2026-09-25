// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package schedule_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/steadybit/cli/internal/jsyaml"
	"github.com/steadybit/cli/internal/platformtest"
	"github.com/steadybit/cli/internal/schedule"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

const id = "01951394-727f-76a0-8675-c7519ebd0ff5"

var saved = map[string]any{
	"id": id, "experimentKey": "ADM-1", "cron": "0 0 9 ? * MON-FRI", "timezone": "Europe/Berlin", "enabled": true,
	"allowParallel": false, "editedBy": map[string]any{"username": "jane"}, "lastUpdated": "2026-09-01T10:00:00Z", "nextExecution": "2026-09-02T07:00:00Z",
}

func TestListFiltersWithRepeatedParameters(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/experiments/schedules/v2", platformtest.Reply{JSON: []any{saved}})

	out, err := platformtest.Stdout(t, func() error {
		return schedule.List(ctx, p.Client, schedule.ListOptions{Teams: []string{"ADM", "OPS"}, Experiments: []string{"ADM-1"}})
	})

	require.NoError(t, err)
	assert.Equal(t, []string{"ADM", "OPS"}, p.Requests("GET /api/experiments/schedules/v2")[0].Query["team"])
	assert.Contains(t, out, "0 0 9 ? * MON-FRI (Europe/Berlin)")
}

func TestGetLeavesOutWhatCannotBeSentBack(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/experiments/schedules/"+id, platformtest.Reply{Body: `{"id":"` + id + `","experimentKey":"ADM-1","cron":"x","editedBy":{},"lastUpdated":"y","nextExecution":"z"}`})
	file := filepath.Join(t.TempDir(), "schedule.yml")

	_, err := platformtest.Stdout(t, func() error { return schedule.Get(ctx, p.Client, schedule.GetOptions{ID: id, File: file}) })

	require.NoError(t, err)
	content, _ := os.ReadFile(file)
	assert.Equal(t, "id: "+id+"\nexperimentKey: ADM-1\ncron: x\n", string(content))
}

func TestApplyCreatesAndWritesTheIdFirst(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/schedules", platformtest.Reply{Status: http.StatusCreated, JSON: saved})
	file := filepath.Join(t.TempDir(), "new.yml")
	require.NoError(t, os.WriteFile(file, []byte("experimentKey: ADM-1\nstartAt: 2030-06-01T09:00:00Z\n"), 0o644))

	out, err := platformtest.Stdout(t, func() error { return schedule.Apply(ctx, p.Client, schedule.ApplyOptions{Files: []string{file}}) })

	require.NoError(t, err)
	assert.Equal(t, "Experiment schedule "+id+" for ADM-1 created.\n", out)
	// A YAML timestamp was a JavaScript Date: sent, and written back, as an ISO string.
	assert.Equal(t, map[string]any{"experimentKey": "ADM-1", "startAt": "2030-06-01T09:00:00.000Z"}, p.Requests("POST /api/experiments/schedules")[0].JSON(t))
	content, _ := os.ReadFile(file)
	assert.Equal(t, "id: "+id+"\nexperimentKey: ADM-1\nstartAt: 2030-06-01T09:00:00.000Z\n", string(content))
}

func TestApplyKeepsAFileThatHasAnId(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/schedules", platformtest.Reply{JSON: saved})
	file := filepath.Join(t.TempDir(), "s.json")
	original := `{"id":"` + id + `","experimentKey":"ADM-1","editedBy":{"username":"jane"}}`
	require.NoError(t, os.WriteFile(file, []byte(original), 0o644))

	_, err := platformtest.Stdout(t, func() error { return schedule.Apply(ctx, p.Client, schedule.ApplyOptions{Files: []string{file}}) })

	require.NoError(t, err)
	assert.Equal(t, map[string]any{"id": id, "experimentKey": "ADM-1"}, p.Requests("POST /api/experiments/schedules")[0].JSON(t))
	content, _ := os.ReadFile(file)
	assert.Equal(t, original, string(content))
}

func TestCreateAndItsRules(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/schedules", platformtest.Reply{Status: http.StatusCreated, JSON: saved})
	vars := jsyaml.NewMap()
	vars.Set("region", "eu")
	no := false

	_, err := platformtest.Stdout(t, func() error {
		return schedule.Create(ctx, p.Client, schedule.CreateOptions{Experiment: "ADM-1", Fields: schedule.Fields{Cron: "0 0 9 ? * *", Timezone: "Europe/Berlin", AllowParallel: &no, Variables: vars}})
	})

	require.NoError(t, err)
	assert.Equal(t, map[string]any{"experimentKey": "ADM-1", "cron": "0 0 9 ? * *", "timezone": "Europe/Berlin", "allowParallel": false, "variables": map[string]any{"region": "eu"}, "enabled": true},
		p.Requests("POST /api/experiments/schedules")[0].JSON(t))
	assert.EqualError(t, schedule.Create(ctx, p.Client, schedule.CreateOptions{Experiment: "ADM-1"}), "Either --cron or --start-at must be specified.")
	assert.EqualError(t, schedule.Create(ctx, p.Client, schedule.CreateOptions{Experiment: "ADM-1", Fields: schedule.Fields{Cron: "x", StartAt: "y"}}), "--cron and --start-at cannot be combined.")
}

func TestUpdateSendsOnlyWhatWasGiven(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("PATCH /api/experiments/schedules/"+id, platformtest.Reply{JSON: saved})

	out, err := platformtest.Stdout(t, func() error {
		if err := schedule.Update(ctx, p.Client, schedule.UpdateOptions{ID: id, Fields: schedule.Fields{Cron: "0 30 8 ? * *"}}); err != nil {
			return err
		}
		return schedule.SetEnabled(ctx, p.Client, id, false)
	})

	require.NoError(t, err)
	requests := p.Requests("PATCH /api/experiments/schedules/" + id)
	assert.Equal(t, map[string]any{"cron": "0 30 8 ? * *"}, requests[0].JSON(t))
	assert.Equal(t, map[string]any{"enabled": false}, requests[1].JSON(t))
	assert.Contains(t, out, "Experiment schedule "+id+" for ADM-1 disabled.")
	assert.EqualError(t, schedule.Update(ctx, p.Client, schedule.UpdateOptions{ID: id}), "Nothing to update. Pass at least one of the options, see --help.")
}

func TestDelete(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("DELETE /api/experiments/schedules/"+id, platformtest.Reply{})
	p.Reply("DELETE /api/experiments/schedules/nope", platformtest.Reply{Status: http.StatusNotFound})

	out, err := platformtest.Stdout(t, func() error { return schedule.Delete(ctx, p.Client, id) })

	require.NoError(t, err)
	assert.Equal(t, "Experiment schedule "+id+" deleted.\n", out)
	assert.EqualError(t, schedule.Delete(ctx, p.Client, "nope"), "Experiment schedule nope not found.")
}
