// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package template_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/steadybit/cli/internal/platformtest"
	"github.com/steadybit/cli/internal/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const id = "d7e65100-1d20-4980-be87-c351704910b8"

func TestListSendsTheFilters(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/experiments/templates", platformtest.Reply{JSON: map[string]any{"templates": []any{map[string]any{"id": id, "templateTitle": "Shop survives"}}}})

	out, err := platformtest.Stdout(t, func() error {
		return template.List(context.Background(), p.Client, template.ListOptions{Search: []string{"shop"}, Tags: []string{"k8s", "db"}})
	})

	require.NoError(t, err)
	assert.Contains(t, out, "Shop survives")
	q := p.Requests("GET /api/experiments/templates")[0].Query
	assert.Equal(t, []string{"shop"}, q["freeTextPhrases"])
	assert.Equal(t, []string{"k8s", "db"}, q["tag"])
}

func TestGetPlaceholdersAsAFileToFillIn(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/experiments/templates/"+id, platformtest.Reply{JSON: map[string]any{"id": id, "placeholders": []any{map[string]any{"key": "CLUSTER"}, map[string]any{"key": "NAMESPACE"}}}})

	out, err := platformtest.Stdout(t, func() error {
		return template.Get(context.Background(), p.Client, template.GetOptions{ID: id, Placeholders: true})
	})

	require.NoError(t, err)
	assert.Equal(t, "CLUSTER: ''\nNAMESPACE: ''\n\n", out)
}

func TestGetReportsAMissingTemplate(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/experiments/templates/"+id, platformtest.Reply{Status: http.StatusNotFound})

	assert.EqualError(t, template.Get(context.Background(), p.Client, template.GetOptions{ID: id}), "Experiment template "+id+" not found.")
}

func TestGetAndApplyRoundTrip(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/experiments/templates/"+id, platformtest.Reply{Body: `{"id":"` + id + `","version":2,"templateTitle":"Shop survives","templateDescription":"d","lanes":[],"futureField":1,"created":"c","createdBy":{},"edited":"e","editedBy":{}}`})
	p.Reply("POST /api/experiments/templates", platformtest.Reply{JSON: map[string]any{"id": id, "templateTitle": "Shop survives"}})
	file := filepath.Join(t.TempDir(), "template.yml")

	out, err := platformtest.Stdout(t, func() error {
		if err := template.Get(context.Background(), p.Client, template.GetOptions{ID: id, File: file}); err != nil {
			return err
		}
		return template.Apply(context.Background(), p.Client, template.ApplyOptions{Files: []string{file}})
	})

	require.NoError(t, err)
	assert.Contains(t, out, "Experiment template Shop survives ("+id+") updated.")
	content, _ := os.ReadFile(file)
	assert.Equal(t, "id: "+id+"\ntemplateTitle: Shop survives\ntemplateDescription: d\nlanes: []\nfutureField: 1\n", string(content))
	// A field the spec does not know yet is sent back, not dropped.
	assert.Equal(t, map[string]any{"id": id, "templateTitle": "Shop survives", "templateDescription": "d", "lanes": []any{}, "futureField": float64(1)},
		p.Requests("POST /api/experiments/templates")[0].JSON(t))
}

func TestApplyCreatesAndWritesTheIdBack(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/templates", platformtest.Reply{Status: http.StatusCreated, JSON: map[string]any{"id": id, "templateTitle": "New"}})
	dir := t.TempDir()
	file := filepath.Join(dir, "new.yml")
	require.NoError(t, os.WriteFile(file, []byte("templateTitle: New\nlanes: []\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "untitled.yml"), []byte("lanes: []\n"), 0o644))

	out, err := platformtest.Stdout(t, func() error {
		return template.Apply(context.Background(), p.Client, template.ApplyOptions{Files: []string{file}})
	})

	require.NoError(t, err)
	assert.Equal(t, "Experiment template New ("+id+") created.\n", out)
	content, _ := os.ReadFile(file)
	assert.Equal(t, "id: "+id+"\ntemplateTitle: New\nlanes: []\n", string(content))
	err = template.Apply(context.Background(), p.Client, template.ApplyOptions{Files: []string{filepath.Join(dir, "untitled.yml")}})
	assert.EqualError(t, err, "Template file '"+filepath.Join(dir, "untitled.yml")+"' does not name a templateTitle.")
}

func TestDelete(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("DELETE /api/experiments/templates/"+id, platformtest.Reply{})

	out, err := platformtest.Stdout(t, func() error {
		return template.Delete(context.Background(), p.Client, template.DeleteOptions{ID: id, Yes: true})
	})

	require.NoError(t, err)
	assert.Equal(t, "Experiment template "+id+" deleted.\n", out)
	assert.EqualError(t, template.Delete(context.Background(), p.Client, template.DeleteOptions{ID: "nope", Yes: true}), "Experiment template nope not found.")
}

func TestImportFromAHub(t *testing.T) {
	const hub = "0194a7d4-0d1f-7b21-9c64-5b6e3c1f2a10"
	p := platformtest.New(t)
	p.Handle("POST /api/experiments/templates/imports", func(r platformtest.Request) platformtest.Reply {
		if r.Query["overwrite"][0] == "false" {
			return platformtest.Reply{Status: http.StatusConflict}
		}
		return platformtest.Reply{}
	})

	out, err := platformtest.Stdout(t, func() error {
		return template.Import(context.Background(), p.Client, template.ImportOptions{Hub: hub, Templates: []string{id}, Overwrite: true})
	})

	require.NoError(t, err)
	assert.Equal(t, "1 experiment template(s) imported from hub "+hub+".\n", out)
	assert.Equal(t, map[string]any{"hubId": hub, "templateIds": []any{id}}, p.Requests("POST /api/experiments/templates/imports")[0].JSON(t))
	assert.EqualError(t, template.Import(context.Background(), p.Client, template.ImportOptions{Hub: hub, Templates: []string{id}}),
		"Some of the templates exist already. Pass --overwrite to replace them.")
}
