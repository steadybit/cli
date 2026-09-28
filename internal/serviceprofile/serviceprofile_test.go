// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package serviceprofile_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/steadybit/cli/v6/internal/platformtest"
	"github.com/steadybit/cli/v6/internal/serviceprofile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

const id = "019eacd7-fb2c-733a-bed5-99a935323db5"

func TestListFilters(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/services/profiles", platformtest.Reply{JSON: map[string]any{"items": []any{map[string]any{
		"id": id, "name": "High Redundancy", "origin": "CUSTOM", "templates": []any{map[string]any{"templateIds": []any{"a", "b"}}}}}}})

	out, err := platformtest.Stdout(t, func() error {
		return serviceprofile.List(ctx, p.Client, serviceprofile.ListOptions{Name: "Redund", Origins: []string{"custom"}, Default: true})
	})

	require.NoError(t, err)
	assert.Contains(t, out, "│ "+id+" │ High Redundancy │ CUSTOM │ false   │         2 │")
	q := p.Requests("GET /api/services/profiles")[0].Query
	assert.Equal(t, []string{"Redund"}, q["name"])
	assert.Equal(t, []string{"CUSTOM"}, q["origin"])
	assert.Equal(t, []string{"true"}, q["defaultProfile"])
}

func TestGetWritesWhatApplyCanSend(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/services/profiles/"+id, platformtest.Reply{Body: `{"name":"P","origin":"CUSTOM","templates":[],"id":"` + id + `","defaultProfile":false,"version":1,"created":"c","createdBy":"u","edited":"e","editedBy":"u"}`})
	file := filepath.Join(t.TempDir(), "p.yml")

	_, err := platformtest.Stdout(t, func() error { return serviceprofile.Get(ctx, p.Client, serviceprofile.GetOptions{ID: id, File: file}) })

	require.NoError(t, err)
	content, _ := os.ReadFile(file)
	assert.Equal(t, "name: P\norigin: CUSTOM\ntemplates: []\nid: "+id+"\n", string(content))
}

func TestApplyDefaultsToACustomProfile(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/services/profiles", platformtest.Reply{Status: http.StatusCreated, JSON: map[string]any{"id": id, "name": "P"}})
	file := filepath.Join(t.TempDir(), "p.yml")
	require.NoError(t, os.WriteFile(file, []byte("name: P\ntemplates: []\n"), 0o644))

	_, err := platformtest.Stdout(t, func() error {
		return serviceprofile.Apply(ctx, p.Client, serviceprofile.ApplyOptions{Files: []string{file}})
	})

	require.NoError(t, err)
	assert.Equal(t, map[string]any{"name": "P", "templates": []any{}, "origin": "CUSTOM"}, p.Requests("POST /api/services/profiles")[0].JSON(t))
	content, _ := os.ReadFile(file)
	assert.Equal(t, "id: "+id+"\nname: P\ntemplates: []\n", string(content))
}

func TestDeleteExplainsProvidedProfiles(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("DELETE /api/services/profiles/"+id, platformtest.Reply{Status: http.StatusUnprocessableEntity})

	err := serviceprofile.Delete(ctx, p.Client, id)

	assert.EqualError(t, err, "Service profile "+id+" is provided by Steadybit and cannot be deleted.")
}
