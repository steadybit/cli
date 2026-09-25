// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package hub_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/steadybit/cli/internal/hub"
	"github.com/steadybit/cli/internal/platformtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

const id = "0194a7d4-0d1f-7b21-9c64-5b6e3c1f2a10"

const stored = `{"hubName":"Reliability Hub","hubLink":"https://hub.example.com","repositoryUrl":"https://example.com/index.json","id":"` + id + `","version":2,` +
	`"templates":[{"id":"t-1"}],"lastSync":"s","lastRepositoryChange":"r","syncError":null,"created":"c","createdBy":{},"edited":"e","editedBy":{}}`

func TestList(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/hubs", platformtest.Reply{JSON: map[string]any{"hubs": []any{map[string]any{"id": id, "hubName": "Reliability Hub"}}}})

	out, err := platformtest.Stdout(t, func() error { return hub.List(ctx, p.Client) })

	require.NoError(t, err)
	assert.Contains(t, out, "│ "+id+" │ Reliability Hub │")
}

func TestGetAndApplyRoundTrip(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/hubs/"+id, platformtest.Reply{Body: stored})
	p.Reply("POST /api/hubs", platformtest.Reply{Body: stored})
	file := filepath.Join(t.TempDir(), "hub.yml")

	out, err := platformtest.Stdout(t, func() error {
		if err := hub.Get(ctx, p.Client, hub.GetOptions{ID: id, File: file}); err != nil {
			return err
		}
		return hub.Apply(ctx, p.Client, hub.ApplyOptions{Files: []string{file}, Synchronize: true})
	})

	require.NoError(t, err)
	assert.Contains(t, out, "Hub Reliability Hub ("+id+") updated.")
	content, _ := os.ReadFile(file)
	assert.Equal(t, "hubName: Reliability Hub\nhubLink: https://hub.example.com\nrepositoryUrl: https://example.com/index.json\nid: "+id+"\n", string(content))
	sent := p.Requests("POST /api/hubs")[0]
	assert.Equal(t, []string{"true"}, sent.Query["synchronize"])
	assert.Equal(t, map[string]any{"hubName": "Reliability Hub", "hubLink": "https://hub.example.com", "repositoryUrl": "https://example.com/index.json", "id": id}, sent.JSON(t))
}

func TestDeleteKeepsTemplatesUnlessAsked(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("DELETE /api/hubs/"+id, platformtest.Reply{})

	out, err := platformtest.Stdout(t, func() error {
		if err := hub.Delete(ctx, p.Client, hub.DeleteOptions{ID: id, Yes: true}); err != nil {
			return err
		}
		return hub.Delete(ctx, p.Client, hub.DeleteOptions{ID: id, Templates: true, Yes: true})
	})

	require.NoError(t, err)
	assert.Equal(t, "Hub "+id+" deleted.\nHub "+id+" deleted.\n", out)
	requests := p.Requests("DELETE /api/hubs/" + id)
	assert.Equal(t, []string{"false"}, requests[0].Query["deleteImportedTemplates"])
	assert.Equal(t, []string{"true"}, requests[1].Query["deleteImportedTemplates"])
}

func TestResync(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/hubs/"+id+"/resync", platformtest.Reply{Body: stored})
	p.Reply("POST /api/hubs/0194a7d4-0d1f-7b21-9c64-5b6e3c1f2a11/resync", platformtest.Reply{Body: `{"hubName":"Broken","templates":[],"syncError":"index.json not found"}`})
	p.Reply("POST /api/hubs/0194a7d4-0d1f-7b21-9c64-5b6e3c1f2a12/resync", platformtest.Reply{Status: http.StatusNotFound})

	out, err := platformtest.Stdout(t, func() error { return hub.Resync(ctx, p.Client, id) })

	require.NoError(t, err)
	assert.Equal(t, "Hub Reliability Hub synchronized, 1 template(s).\n", out)
	assert.EqualError(t, hub.Resync(ctx, p.Client, "0194a7d4-0d1f-7b21-9c64-5b6e3c1f2a11"), "Hub Broken could not be synchronized: index.json not found")
	assert.EqualError(t, hub.Resync(ctx, p.Client, "0194a7d4-0d1f-7b21-9c64-5b6e3c1f2a12"), "Hub 0194a7d4-0d1f-7b21-9c64-5b6e3c1f2a12 not found.")
}
