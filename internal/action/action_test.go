// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package action_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/steadybit/cli/internal/action"
	"github.com/steadybit/cli/internal/platformtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func TestListWalksEveryPageAndFiltersByKind(t *testing.T) {
	p := platformtest.New(t)
	p.Handle("GET /api/actions", func(r platformtest.Request) platformtest.Reply {
		if r.Query["page"][0] == "0" {
			return platformtest.Reply{JSON: map[string]any{"actions": []any{map[string]any{"id": "stress-cpu", "name": "Stress CPU", "kind": "ATTACK", "category": "resource"}}, "nextPage": 1}}
		}
		return platformtest.Reply{JSON: map[string]any{"actions": []any{map[string]any{"id": "http-check", "name": "HTTP Check", "kind": "CHECK"}}}}
	})

	out, err := platformtest.Stdout(t, func() error { return action.List(ctx, p.Client, action.ListOptions{Kinds: []string{"check"}}) })

	require.NoError(t, err)
	assert.Contains(t, out, "│ http-check │ HTTP Check │ CHECK │          │")
	assert.NotContains(t, out, "stress-cpu")
	assert.Len(t, p.Requests("GET /api/actions"), 2)
	assert.Equal(t, []string{"100"}, p.Requests("GET /api/actions")[0].Query["size"])
}

func TestListSaysWhenThereIsNone(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/actions", platformtest.Reply{JSON: map[string]any{"actions": []any{}}})

	out, err := platformtest.Stdout(t, func() error { return action.List(ctx, p.Client, action.ListOptions{}) })

	require.NoError(t, err)
	assert.Equal(t, "No actions found.\n", out)
}

func TestGet(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/actions/stress-cpu", platformtest.Reply{Body: `{"id":"stress-cpu","name":"Stress CPU","parameters":[{"name":"duration","type":"duration"}]}`})
	p.Reply("GET /api/actions/nope", platformtest.Reply{Status: http.StatusNotFound})

	out, err := platformtest.Stdout(t, func() error { return action.Get(ctx, p.Client, action.GetOptions{ID: "stress-cpu"}) })

	require.NoError(t, err)
	assert.Equal(t, "id: stress-cpu\nname: Stress CPU\nparameters:\n  - name: duration\n    type: duration\n\n", out)
	assert.EqualError(t, action.Get(ctx, p.Client, action.GetOptions{ID: "nope"}), "Action nope not found.")
}
