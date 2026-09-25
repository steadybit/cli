// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package template_test

import (
	"context"
	"net/http"
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
