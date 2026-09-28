// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package integration_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/steadybit/cli/v6/internal/integration"
	"github.com/steadybit/cli/v6/internal/platformtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

const id = "0190d7b2-7d3e-7a4b-8c5d-6e7f8a9b0c1d"

func TestEveryKindListsFromItsOwnEndpoint(t *testing.T) {
	paths := map[string]string{
		"webhook": "/api/integrations/webhook", "slack": "/api/integrations/slack",
		"preflight": "/api/integrations/preflight", "preflight-action": "/api/integrations/preflight-action",
	}
	for _, k := range integration.Kinds {
		t.Run(k.Name, func(t *testing.T) {
			p := platformtest.New(t)
			p.Reply("GET "+paths[k.Name], platformtest.Reply{JSON: map[string]any{"content": []any{
				map[string]any{"id": id, "name": "Notify", "scope": "TEAM", "team": "ADM", k.Column: "where"},
			}}})

			out, err := platformtest.Stdout(t, func() error { return integration.List(ctx, p.Client, k, "") })

			require.NoError(t, err)
			assert.Contains(t, out, "│ "+id+" │ Notify │ TEAM  │ ADM  │ where")
		})
	}
}

func TestListSaysWhenThereIsNone(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/integrations/slack", platformtest.Reply{JSON: map[string]any{"content": []any{}}})

	out, err := platformtest.Stdout(t, func() error { return integration.List(ctx, p.Client, integration.Slack, "") })

	require.NoError(t, err)
	assert.Equal(t, "No Slack integrations found.\n", out)
}

func TestGetAndApplyRoundTrip(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/integrations/preflight-action/"+id, platformtest.Reply{Body: `{"id":"` + id + `","version":3,"scope":"GLOBAL","name":"Gate","preflightActionId":"com.example.gate"}`})
	p.Reply("POST /api/integrations/preflight-action", platformtest.Reply{JSON: map[string]any{"id": id, "name": "Gate"}})
	file := filepath.Join(t.TempDir(), "gate.yml")

	out, err := platformtest.Stdout(t, func() error {
		if err := integration.Get(ctx, p.Client, integration.PreflightAction, integration.GetOptions{ID: id, File: file}); err != nil {
			return err
		}
		return integration.Apply(ctx, p.Client, integration.PreflightAction, integration.ApplyOptions{Files: []string{file}})
	})

	require.NoError(t, err)
	assert.Equal(t, "Preflight action integration "+id+" written to "+file+".\nPreflight action integration Gate ("+id+") updated.\n", out)
	assert.Equal(t, map[string]any{"id": id, "scope": "GLOBAL", "name": "Gate", "preflightActionId": "com.example.gate"},
		p.Requests("POST /api/integrations/preflight-action")[0].JSON(t))
}

func TestApplyCreatesAndRefusesAMaskedSecret(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/integrations/webhook", platformtest.Reply{Status: http.StatusCreated, JSON: map[string]any{"id": id, "name": "Notify"}})
	dir := t.TempDir()
	file := filepath.Join(dir, "webhook.yml")
	require.NoError(t, os.WriteFile(file, []byte("name: Notify\nscope: GLOBAL\nurl: https://example.com\nsecret: s3cret\n"), 0o644))
	masked := filepath.Join(dir, "masked.yml")
	require.NoError(t, os.WriteFile(masked, []byte("id: "+id+"\nname: Notify\nsecret: '******'\n"), 0o644))

	out, err := platformtest.Stdout(t, func() error {
		return integration.Apply(ctx, p.Client, integration.Webhook, integration.ApplyOptions{Files: []string{file}})
	})

	require.NoError(t, err)
	assert.Equal(t, "Webhook integration Notify ("+id+") created.\n", out)
	content, _ := os.ReadFile(file)
	assert.Equal(t, "id: "+id+"\nname: Notify\nscope: GLOBAL\nurl: https://example.com\nsecret: s3cret\n", string(content))
	err = integration.Apply(ctx, p.Client, integration.Webhook, integration.ApplyOptions{Files: []string{masked}})
	assert.EqualError(t, err, "Webhook integration file '"+masked+"' holds the masked secret `get` writes. Put the secret in, or remove it for none.")
	assert.Len(t, p.Requests("POST /api/integrations/webhook"), 1)
}

func TestDelete(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("DELETE /api/integrations/slack/"+id, platformtest.Reply{JSON: map[string]any{"id": id}})
	p.Reply("DELETE /api/integrations/preflight/"+id, platformtest.Reply{Status: http.StatusNotFound})

	out, err := platformtest.Stdout(t, func() error {
		return integration.Delete(ctx, p.Client, integration.Slack, integration.DeleteOptions{ID: id, Yes: true})
	})

	require.NoError(t, err)
	assert.Equal(t, "Slack integration "+id+" deleted.\n", out)
	assert.EqualError(t, integration.Delete(ctx, p.Client, integration.Preflight, integration.DeleteOptions{ID: id, Yes: true}), "Preflight webhook "+id+" not found.")
	assert.EqualError(t, integration.Delete(ctx, p.Client, integration.Preflight, integration.DeleteOptions{ID: "x", Yes: true}), "Preflight webhook x not found.")
}
