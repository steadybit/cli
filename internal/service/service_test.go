// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package service_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/steadybit/cli/internal/experiment"
	"github.com/steadybit/cli/internal/jsyaml"
	"github.com/steadybit/cli/internal/platformtest"
	"github.com/steadybit/cli/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

const id = "019cd80d-a4c9-775b-bdf8-2672a280ce7c"

const stored = `{"name":"Checkout","environment":"Global","team":"ADM","query":"x","validations":[],"serviceProfile":"P","variables":{"region":"eu"},"id":"` + id + `","version":3,"created":"c","createdBy":{},"edited":"e","editedBy":{}}`

func TestListWalksEveryPage(t *testing.T) {
	p := platformtest.New(t)
	p.Handle("GET /api/services", func(r platformtest.Request) platformtest.Reply {
		if r.Query["page"][0] == "0" {
			return platformtest.Reply{JSON: map[string]any{"items": []any{map[string]any{"id": "a", "name": "Checkout"}}, "nextPage": 1}}
		}
		return platformtest.Reply{JSON: map[string]any{"items": []any{map[string]any{"id": "b", "name": "Catalog"}}}}
	})

	out, err := platformtest.Stdout(t, func() error { return service.List(ctx, p.Client, service.ListOptions{Teams: []string{"ADM"}}) })

	require.NoError(t, err)
	assert.Contains(t, out, "Checkout")
	assert.Contains(t, out, "Catalog")
	first := p.Requests("GET /api/services")[0]
	assert.Equal(t, []string{"ADM"}, first.Query["teamKey"])
	assert.Equal(t, []string{"100"}, first.Query["size"])
}

func TestGetAndApplyRoundTrip(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/services/"+id, platformtest.Reply{Body: stored})
	p.Reply("POST /api/services", platformtest.Reply{JSON: map[string]any{"id": id, "name": "Checkout"}})
	file := filepath.Join(t.TempDir(), "service.yml")

	_, err := platformtest.Stdout(t, func() error {
		if err := service.Get(ctx, p.Client, service.GetOptions{ID: id, File: file}); err != nil {
			return err
		}
		return service.Apply(ctx, p.Client, service.ApplyOptions{Files: []string{file}})
	})

	require.NoError(t, err)
	content, _ := os.ReadFile(file)
	for _, field := range []string{"version", "created", "edited"} {
		assert.NotContains(t, string(content), field+":")
	}
	sent := p.Requests("POST /api/services")[0]
	assert.Equal(t, []string{"false"}, sent.Query["deleteExperiments"])
	assert.NotContains(t, sent.JSON(t), "version")
}

func TestApplyPointsAtDeleteExperiments(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/services", platformtest.Reply{Status: 422, JSON: map[string]any{"violations": []any{map[string]any{"message": "Cannot remove templates without setting `deleteExperiments` to true."}}}})
	file := filepath.Join(t.TempDir(), "service.yml")
	require.NoError(t, os.WriteFile(file, []byte("name: Checkout\n"), 0o644))

	err := service.Apply(ctx, p.Client, service.ApplyOptions{Files: []string{file}})

	assert.EqualError(t, err, "Service Checkout was not saved: the change would remove provided experiments. Pass --delete-experiments to delete them.")
}

func TestRiskGate(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/services/"+id+"/risk", platformtest.Reply{JSON: map[string]any{
		"risk": 42, "categoryRisks": map[string]any{"Redundancy": map[string]any{"total": 40, "experiment": 50}},
		"experimentRisks": []any{map[string]any{"experimentKey": "ADM-1", "risk": 60}}, "lastCalculated": "then",
	}})
	at := func(n int) *int { return &n }

	out, err := platformtest.Stdout(t, func() error { return service.Risk(ctx, p.Client, service.RiskOptions{ID: id, FailAbove: at(42)}) })

	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(out, "Risk of service "+id+": 42 (calculated then)\n"))
	assert.Contains(t, out, "│ Redundancy │    40 │          50 │        │")
	_, err = platformtest.Stdout(t, func() error { return service.Risk(ctx, p.Client, service.RiskOptions{ID: id, FailAbove: at(41)}) })
	assert.EqualError(t, err, "Risk of service "+id+" is 42, above the accepted 41.")
}

func TestExperimentsOfAService(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/services/"+id+"/experiments", platformtest.Reply{JSON: map[string]any{"items": []any{
		map[string]any{"templateId": "t-1", "category": "Scalability", "associationType": "PROVIDED"},
	}}})
	p.Reply("POST /api/services/"+id+"/experiments/custom", platformtest.Reply{Status: http.StatusCreated})
	p.Reply("DELETE /api/services/"+id+"/experiments/custom", platformtest.Reply{})
	p.Reply("POST /api/services/"+id+"/experiments/provided", platformtest.Reply{Status: http.StatusCreated, Headers: map[string]string{"Location": p.URL + "/api/experiments/ADM-9"}})

	out, err := platformtest.Stdout(t, func() error {
		if err := service.ListExperiments(ctx, p.Client, service.ExperimentListOptions{ID: id, Types: []string{"provided"}}); err != nil {
			return err
		}
		if err := service.Link(ctx, p.Client, id, "ADM-1", "Redundancy"); err != nil {
			return err
		}
		if err := service.Unlink(ctx, p.Client, id, "ADM-1"); err != nil {
			return err
		}
		placeholders := jsyaml.NewMap()
		placeholders.Set("REPLICAS", "3")
		return service.Provide(ctx, p.Client, service.ProvideOptions{ID: id, TemplateOptions: experiment.TemplateOptions{
			Template: "d7e65100-1d20-4980-be87-c351704910b8", Placeholder: placeholders, ResetProperties: true}})
	})

	require.NoError(t, err)
	assert.Contains(t, out, "(not created)")
	assert.Contains(t, out, "Provided experiment ADM-9 of service "+id+" created from template d7e65100-1d20-4980-be87-c351704910b8.")
	assert.Equal(t, []string{"PROVIDED"}, p.Requests("GET /api/services/" + id + "/experiments")[0].Query["type"])
	assert.Equal(t, map[string]any{"experimentKey": "ADM-1", "category": "Redundancy"}, p.Requests("POST /api/services/" + id + "/experiments/custom")[0].JSON(t))
	assert.Equal(t, []string{"ADM-1"}, p.Requests("DELETE /api/services/" + id + "/experiments/custom")[0].Query["experimentKey"])
}

func TestVariables(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("PATCH /api/services/"+id+"/variables", platformtest.Reply{})
	p.Reply("PUT /api/services/"+id+"/variables", platformtest.Reply{})
	file := filepath.Join(t.TempDir(), "vars.yml")
	require.NoError(t, os.WriteFile(file, []byte("hosts: [a, b]\nregion: us\n"), 0o644))

	_, err := platformtest.Stdout(t, func() error {
		if err := service.SetVariables(ctx, p.Client, []string{"region=eu", "url=http://x?a=b"}, service.VariableSetOptions{ID: id, File: file}); err != nil {
			return err
		}
		return service.SetVariables(ctx, p.Client, nil, service.VariableSetOptions{ID: id, Replace: true})
	})

	require.NoError(t, err)
	assert.Equal(t, map[string]any{"hosts": []any{"a", "b"}, "region": "eu", "url": "http://x?a=b"}, p.Requests("PATCH /api/services/" + id + "/variables")[0].JSON(t))
	assert.Equal(t, map[string]any{}, p.Requests("PUT /api/services/" + id + "/variables")[0].JSON(t))
	assert.EqualError(t, service.SetVariables(ctx, p.Client, []string{"novalue"}, service.VariableSetOptions{ID: id}), "'novalue' is not in the form KEY=VALUE.")
	assert.EqualError(t, service.SetVariables(ctx, p.Client, nil, service.VariableSetOptions{ID: id}), "No variables given. Pass KEY=VALUE arguments or --file.")
}

func TestAMalformedIdReadsAsNotFound(t *testing.T) {
	assert.EqualError(t, service.Delete(ctx, nil, "nope"), "Service nope not found.")
}
