// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package report_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/steadybit/cli/internal/platformtest"
	"github.com/steadybit/cli/internal/report"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func kind(t *testing.T, name string) report.Kind {
	for _, k := range report.Kinds {
		if k.Name == name {
			return k
		}
	}
	t.Fatalf("no report %s", name)
	return report.Kind{}
}

const series = `{"from":"2026-08-01","to":"2026-09-01","rollup":"MONTHLY","groupBy":"NONE","series":[{"name":"users","values":[["2026-08-01",16]]}]}`

func TestEveryReportPostsToItsEndpoint(t *testing.T) {
	paths := map[string]string{
		"users": "/api/reports/users", "teams": "/api/reports/teams", "environments": "/api/reports/environments",
		"experiments-executed": "/api/reports/experiments/executed", "experiments-created": "/api/reports/experiments/created",
		"services-distribution": "/api/reports/services/distribution", "services-by-category": "/api/reports/services/by-category",
		"services-average": "/api/reports/services/average",
	}
	require.Len(t, report.Kinds, len(paths))
	for _, k := range report.Kinds {
		t.Run(k.Name, func(t *testing.T) {
			p := platformtest.New(t)
			p.Reply("POST "+paths[k.Name], platformtest.Reply{Body: series})

			out, err := platformtest.Stdout(t, func() error {
				return report.Run(ctx, p.Client, k, report.Options{From: "2026-08-01", To: "2026-09-01"})
			})

			require.NoError(t, err)
			assert.Contains(t, out, "series:\n  - name: users\n")
			assert.Equal(t, map[string]any{"from": "2026-08-01", "to": "2026-09-01"}, p.Requests("POST " + paths[k.Name])[0].JSON(t))
		})
	}
}

func TestDefaultsToTheLast30Days(t *testing.T) {
	report.Today = func() time.Time { return time.Date(2026, 3, 31, 15, 0, 0, 0, time.UTC) }
	t.Cleanup(func() { report.Today = func() time.Time { return time.Now().UTC() } })
	p := platformtest.New(t)
	p.Reply("POST /api/reports/users", platformtest.Reply{Body: series})

	_, err := platformtest.Stdout(t, func() error { return report.Run(ctx, p.Client, kind(t, "users"), report.Options{Rollup: "daily"}) })

	require.NoError(t, err)
	assert.Equal(t, map[string]any{"from": "2026-03-01", "to": "2026-03-31", "rollup": "DAILY"}, p.Requests("POST /api/reports/users")[0].JSON(t))
}

func TestFiltersAreSentAsIds(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/teams/ADM", platformtest.Reply{JSON: map[string]any{"id": "team-1", "key": "ADM"}})
	p.Reply("GET /api/environments", platformtest.Reply{JSON: map[string]any{"environments": []any{map[string]any{"id": "env-1", "name": "Global"}, map[string]any{"id": "env-2", "name": "Global 2"}}}})
	p.Reply("POST /api/reports/services/average", platformtest.Reply{Body: series})

	_, err := platformtest.Stdout(t, func() error {
		return report.Run(ctx, p.Client, kind(t, "services-average"), report.Options{
			From: "2026-08-01", To: "2026-09-01", Teams: []string{"ADM"}, Environments: []string{"Global"}, Services: []string{"svc-1"},
			ServiceProperties: []string{"tier=gold", "tier=silver", "region=eu"}, Categories: []string{"Redundancy"}, Type: "json",
		})
	})

	require.NoError(t, err)
	assert.Equal(t, map[string]any{
		"from": "2026-08-01", "to": "2026-09-01", "teamIds": []any{"team-1"}, "environmentIds": []any{"env-1"}, "serviceIds": []any{"svc-1"},
		"serviceProperties": map[string]any{"tier": []any{"gold", "silver"}, "region": []any{"eu"}}, "categoryKeys": []any{"Redundancy"},
	}, p.Requests("POST /api/reports/services/average")[0].JSON(t))
}

func TestGroupBy(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/reports/experiments/executed", platformtest.Reply{Body: series})

	_, err := platformtest.Stdout(t, func() error {
		return report.Run(ctx, p.Client, kind(t, "experiments-executed"), report.Options{From: "2026-08-01", To: "2026-09-01", GroupBy: "state"})
	})

	require.NoError(t, err)
	assert.Equal(t, []string{"STATE"}, p.Requests("POST /api/reports/experiments/executed")[0].Query["groupBy"])
	err = report.Run(ctx, p.Client, kind(t, "experiments-created"), report.Options{GroupBy: "state"})
	assert.EqualError(t, err, "--group-by must be one of NONE, CREATED_VIA, ORIGIN, not 'state'.")
}

func TestRefusals(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/teams/NOPE", platformtest.Reply{Status: http.StatusNotFound})
	p.Reply("GET /api/environments", platformtest.Reply{JSON: map[string]any{"environments": []any{}}})
	users, executed := kind(t, "users"), kind(t, "experiments-executed")

	assert.EqualError(t, report.Run(ctx, p.Client, users, report.Options{From: "yesterday"}), "--from 'yesterday' is not a date like 2026-09-01.")
	assert.EqualError(t, report.Run(ctx, p.Client, users, report.Options{Rollup: "weekly"}), "--rollup must be DAILY or MONTHLY, not 'weekly'.")
	assert.EqualError(t, report.Run(ctx, p.Client, executed, report.Options{Teams: []string{"NOPE"}}), "Team NOPE not found.")
	assert.EqualError(t, report.Run(ctx, p.Client, executed, report.Options{Environments: []string{"Nowhere"}}), "Environment Nowhere not found.")
	assert.EqualError(t, report.Run(ctx, p.Client, kind(t, "services-average"), report.Options{ServiceProperties: []string{"tier"}}), "'tier' is not in the form KEY=VALUE.")
}
