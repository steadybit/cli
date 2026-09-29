// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package execution_test

import (
	"fmt"
	"testing"

	"github.com/steadybit/cli/v6/internal/execution"
	"github.com/steadybit/cli/v6/internal/output"
	"github.com/steadybit/cli/v6/internal/platformtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const executions = "POST /api/experiments/executions"

func pageItem(id int, state string) map[string]any {
	return map[string]any{"id": id, "experimentKey": "ADM-1", "name": "Shop survives", "state": state,
		"started": "2026-09-28T10:00:00.123456Z", "ended": "2026-09-28T10:01:30.5Z", "createdBy": "jane"}
}

func TestListWalksEveryPageWithTheFilters(t *testing.T) {
	p := platformtest.New(t)
	p.Handle(executions, func(r platformtest.Request) platformtest.Reply {
		if r.JSON(t).(map[string]any)["page"] == 0.0 {
			scheduled := pageItem(2, "FAILED")
			scheduled["scheduled"] = true
			return platformtest.Reply{JSON: map[string]any{"items": []any{scheduled}, "totalItems": 2, "nextPage": 1}}
		}
		named := pageItem(1, "COMPLETED")
		named["createdByDetails"] = map[string]any{"username": "jane", "name": "Jane Doe"}
		return platformtest.Reply{JSON: map[string]any{"items": []any{named}, "totalItems": 2, "nextPage": nil}}
	})

	out, err := platformtest.Stdout(t, func() error {
		return execution.List(ctx, p.Client, execution.ListOptions{
			Name: "shop", Teams: []string{"ADM", "OPS"}, Experiments: []string{"ADM-1"}, Environments: []string{"Global"},
			Services: []string{"shop"}, States: []string{"failed", "COMPLETED"}, From: "2026-09-28", To: "2026-09-29T12:00:00Z",
		})
	})

	require.NoError(t, err)
	assert.Contains(t, out, "│  2 │ ADM-1      │ Shop survives │ failed    │ 2026-09-28T10:00:00Z │    1m30s │ schedule │")
	assert.Contains(t, out, "│  1 │ ADM-1      │ Shop survives │ completed │ 2026-09-28T10:00:00Z │    1m30s │ Jane Doe │")
	assert.NotContains(t, out, "Showing")
	requests := p.Requests(executions)
	require.Len(t, requests, 2)
	assert.Equal(t, map[string]any{
		"page": 0.0, "size": 100.0, "name": "shop", "teamKeys": []any{"ADM", "OPS"}, "experimentKeys": []any{"ADM-1"},
		"environments": []any{"Global"}, "services": []any{"shop"}, "states": []any{"FAILED", "COMPLETED"},
		"createdFrom": "2026-09-28T00:00:00Z", "createdTo": "2026-09-29T12:00:00Z",
	}, requests[0].JSON(t))
	assert.Equal(t, 1.0, requests[1].JSON(t).(map[string]any)["page"])
}

// A tenant keeps thousands of runs, so a listing stops at the limit instead of paging
// through all of them, and says how many it left out.
func TestListStopsAtTheLimit(t *testing.T) {
	p := platformtest.New(t)
	p.Handle(executions, func(r platformtest.Request) platformtest.Reply {
		body := r.JSON(t).(map[string]any)
		page, size := int(body["page"].(float64)), int(body["size"].(float64))
		var items []any
		for i := range size {
			items = append(items, pageItem(1000-page*size-i, "COMPLETED"))
		}
		return platformtest.Reply{JSON: map[string]any{"items": items, "totalItems": 1000, "nextPage": page + 1}}
	})

	out, err := platformtest.Stdout(t, func() error { return execution.List(ctx, p.Client, execution.ListOptions{Limit: 3}) })
	require.NoError(t, err)
	assert.Contains(t, out, "│  998 │")
	assert.NotContains(t, out, "│  997 │")
	assert.Contains(t, out, "Showing the 3 most recent of 1000 experiment runs. Raise --limit, or 0 for all.\n")
	assert.Equal(t, 3.0, p.Requests(executions)[0].JSON(t).(map[string]any)["size"])

	out, err = platformtest.Stdout(t, func() error {
		return execution.List(ctx, p.Client, execution.ListOptions{Limit: 150, Type: "json"})
	})
	require.NoError(t, err)
	assert.Contains(t, out, `"id": 851`)
	assert.NotContains(t, out, `"id": 850`)
	requests := p.Requests(executions)
	require.Len(t, requests, 3, "150 runs are two pages of 100")
	assert.Equal(t, 1.0, requests[2].JSON(t).(map[string]any)["page"])
}

func TestListPrintsThePlatformItems(t *testing.T) {
	p := platformtest.New(t)
	p.Reply(executions, platformtest.Reply{Body: `{"items":[{"id":7,"state":"FAILED","properties":{"ticket":"SHOP-1"}}],"totalItems":1,"nextPage":null}`})

	out, err := platformtest.Stdout(t, func() error { return execution.List(ctx, p.Client, execution.ListOptions{Type: "yaml"}) })
	require.NoError(t, err)
	assert.Equal(t, "- id: 7\n  state: FAILED\n  properties:\n    ticket: SHOP-1\n", out)

	output.JQ = ".[].properties.ticket"
	t.Cleanup(func() { output.JQ = "" })
	out, err = platformtest.Stdout(t, func() error { return execution.List(ctx, p.Client, execution.ListOptions{}) })
	require.NoError(t, err)
	assert.Equal(t, "SHOP-1\n", out)
}

// A pipeline gates on the search: any matching run fails it, counted over all pages
// even when the limit printed fewer.
func TestFailOnMatch(t *testing.T) {
	p := platformtest.New(t)
	total := 0
	p.Handle(executions, func(platformtest.Request) platformtest.Reply {
		var items []any
		for i := range min(total, 1) {
			items = append(items, pageItem(i+1, "FAILED"))
		}
		return platformtest.Reply{JSON: map[string]any{"items": items, "totalItems": total}}
	})
	list := func() (string, error) {
		return platformtest.Stdout(t, func() error {
			return execution.List(ctx, p.Client, execution.ListOptions{States: []string{"FAILED"}, Limit: 1, FailOnMatch: true})
		})
	}

	out, err := list()
	require.NoError(t, err)
	assert.Equal(t, "No experiment runs found.\n", out)

	total = 1
	_, err = list()
	assert.EqualError(t, err, "1 experiment run matches.")

	total = 12
	out, err = list()
	assert.EqualError(t, err, "12 experiment runs match.")
	assert.Contains(t, out, "│  1 │ ADM-1", "the matching runs are still listed")
}

func TestListRejectsBadFilters(t *testing.T) {
	p := platformtest.New(t)
	for options, message := range map[*execution.ListOptions]string{
		{States: []string{"DONE"}}: "--state must be one of CREATED, PREPARED, RUNNING, FAILED, CANCELED, COMPLETED, ERRORED, not 'DONE'.",
		{From: "yesterday"}:        "--from 'yesterday' is neither a date like 2026-09-01 nor a time like 2026-09-01T12:00:00Z.",
		{Limit: -1}:                "--limit must be 0 or more, not -1.",
	} {
		assert.EqualError(t, execution.List(ctx, p.Client, *options), message, fmt.Sprint(*options))
	}
	assert.Empty(t, p.Requests(executions))
}
