// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package environment_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/steadybit/cli/internal/environment"
	"github.com/steadybit/cli/internal/platformtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

const id = "0190d7b2-1c5e-7f3a-8e4b-2d6f9a1c3e57"

func TestListSearches(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/environments", platformtest.Reply{JSON: map[string]any{"environments": []any{map[string]any{"id": id, "name": "Prod", "state": "READY"}}}})

	out, err := platformtest.Stdout(t, func() error { return environment.List(ctx, p.Client, environment.ListOptions{Search: "pro"}) })

	require.NoError(t, err)
	assert.Contains(t, out, "│ "+id+" │ Prod │ READY │")
	assert.Equal(t, []string{"pro"}, p.Requests("GET /api/environments")[0].Query["search"])
}

func TestListSaysWhenThereIsNone(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/environments", platformtest.Reply{JSON: map[string]any{"environments": []any{}}})

	out, err := platformtest.Stdout(t, func() error { return environment.List(ctx, p.Client, environment.ListOptions{}) })

	require.NoError(t, err)
	assert.Equal(t, "No environments found.\n", out)
	assert.NotContains(t, p.Requests("GET /api/environments")[0].Query, "search")
}

func TestGetAndApplyRoundTrip(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/environments/"+id, platformtest.Reply{Body: `{"id":"` + id + `","name":"Prod","version":4,"query":"k8s.cluster-name=\"prod\"","state":"READY"}`})
	p.Reply("POST /api/environments", platformtest.Reply{JSON: map[string]any{"id": id, "name": "Prod"}})
	file := filepath.Join(t.TempDir(), "environment.yml")

	out, err := platformtest.Stdout(t, func() error {
		if err := environment.Get(ctx, p.Client, environment.GetOptions{ID: id, File: file}); err != nil {
			return err
		}
		return environment.Apply(ctx, p.Client, environment.ApplyOptions{Files: []string{file}})
	})

	require.NoError(t, err)
	assert.Contains(t, out, "Environment Prod ("+id+") updated.")
	content, _ := os.ReadFile(file)
	assert.Equal(t, "id: "+id+"\nname: Prod\nquery: k8s.cluster-name=\"prod\"\n", string(content))
	assert.Equal(t, map[string]any{"id": id, "name": "Prod", "query": `k8s.cluster-name="prod"`}, p.Requests("POST /api/environments")[0].JSON(t))
}

func TestApplyCreatesAndWritesTheIdBack(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/environments", platformtest.Reply{Status: http.StatusCreated, JSON: map[string]any{"id": id, "name": "New"}})
	file := filepath.Join(t.TempDir(), "new.json")
	require.NoError(t, os.WriteFile(file, []byte(`{"name":"New","query":"x"}`), 0o644))

	out, err := platformtest.Stdout(t, func() error { return environment.Apply(ctx, p.Client, environment.ApplyOptions{Files: []string{file}}) })

	require.NoError(t, err)
	assert.Equal(t, "Environment New ("+id+") created.\n", out)
	content, _ := os.ReadFile(file)
	assert.Equal(t, "{\n  \"id\": \""+id+"\",\n  \"name\": \"New\",\n  \"query\": \"x\"\n}", string(content))
}

func TestDelete(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("DELETE /api/environments/"+id, platformtest.Reply{Status: http.StatusNotFound})

	err := environment.Delete(ctx, p.Client, environment.DeleteOptions{ID: id, Yes: true})

	assert.EqualError(t, err, "Environment "+id+" not found.")
	assert.EqualError(t, environment.Delete(ctx, p.Client, environment.DeleteOptions{ID: "Prod", Yes: true}), "Environment Prod not found.")
}

func TestVariablesMergeUnlessReplaced(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/environments/"+id+"/variables", platformtest.Reply{Body: `{"region":"eu"}`})
	p.Reply("PUT /api/environments/"+id+"/variables", platformtest.Reply{})
	p.Reply("POST /api/environments/"+id+"/variables", platformtest.Reply{})

	out, err := platformtest.Stdout(t, func() error {
		if err := environment.GetVariables(ctx, p.Client, environment.VariableGetOptions{ID: id, Type: "json"}); err != nil {
			return err
		}
		if err := environment.SetVariables(ctx, p.Client, []string{"region=us"}, environment.VariableSetOptions{ID: id}); err != nil {
			return err
		}
		return environment.SetVariables(ctx, p.Client, nil, environment.VariableSetOptions{ID: id, Replace: true})
	})

	require.NoError(t, err)
	assert.Equal(t, "{\n  \"region\": \"eu\"\n}\n1 variable(s) of environment "+id+" set.\n0 variable(s) of environment "+id+" set, all others removed.\n", out)
	assert.Equal(t, map[string]any{"region": "us"}, p.Requests("PUT /api/environments/" + id + "/variables")[0].JSON(t))
	assert.Equal(t, map[string]any{}, p.Requests("POST /api/environments/" + id + "/variables")[0].JSON(t))
	assert.EqualError(t, environment.SetVariables(ctx, p.Client, nil, environment.VariableSetOptions{ID: id}), "No variables given. Pass KEY=VALUE arguments or --file.")
}
