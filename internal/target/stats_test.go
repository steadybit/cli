// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package target_test

import (
	"net/http"
	"testing"

	"github.com/steadybit/cli/v6/internal/output"
	"github.com/steadybit/cli/v6/internal/platformtest"
	"github.com/steadybit/cli/v6/internal/target"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const stats = `{"com.steadybit.extension_kubernetes.kubernetes-pod":11108,"com.steadybit.extension_aws.zone":12,"com.steadybit.extension_host.host":25}`

func TestStatsCountsEveryTypeInOrder(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/target-stats", platformtest.Reply{Body: stats})

	out, err := platformtest.Stdout(t, func() error { return target.Stats(ctx, p.Client, target.StatsOptions{}) })

	require.NoError(t, err)
	assert.Contains(t, out, "│ com.steadybit.extension_aws.zone                  │      12 │\n"+
		"│ com.steadybit.extension_host.host                 │      25 │\n"+
		"│ com.steadybit.extension_kubernetes.kubernetes-pod │   11108 │")
}

func TestStatsOfAQuery(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/target-stats", platformtest.Reply{Body: `{"com.steadybit.extension_host.host":3}`})

	out, err := platformtest.Stdout(t, func() error {
		return target.Stats(ctx, p.Client, target.StatsOptions{Query: `k8s.namespace="shop"`, Type: "yaml"})
	})

	require.NoError(t, err)
	assert.Equal(t, "com.steadybit.extension_host.host: 3\n", out)
	assert.Equal(t, map[string]any{"query": `k8s.namespace="shop"`}, p.Requests("POST /api/target-stats")[0].JSON(t))
}

func TestStatsPrintsThePlatformsValue(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/target-stats", platformtest.Reply{Body: stats})

	json, err := platformtest.Stdout(t, func() error { return target.Stats(ctx, p.Client, target.StatsOptions{Type: "json"}) })
	require.NoError(t, err)
	output.JQ = `.["com.steadybit.extension_host.host"]`
	t.Cleanup(func() { output.JQ = "" })
	jq, err := platformtest.Stdout(t, func() error { return target.Stats(ctx, p.Client, target.StatsOptions{}) })
	require.NoError(t, err)

	assert.Equal(t, "{\n  \"com.steadybit.extension_kubernetes.kubernetes-pod\": 11108,\n  \"com.steadybit.extension_aws.zone\": 12,\n  \"com.steadybit.extension_host.host\": 25\n}\n", json)
	assert.Equal(t, "25\n", jq)
}

func TestStatsSaysWhenNothingMatches(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/target-stats", platformtest.Reply{Body: `{}`})

	out, err := platformtest.Stdout(t, func() error { return target.Stats(ctx, p.Client, target.StatsOptions{Query: `k8s.namespace="none"`}) })

	require.NoError(t, err)
	assert.Equal(t, "No targets found.\n", out)
}

func TestStatsReportsAnInvalidQuery(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/target-stats", platformtest.Reply{Status: http.StatusUnprocessableEntity,
		Body: `{"title":"Constraint Violation","status":422,"violations":[{"field":"query","message":"Failed to parse query"}]}`})

	err := target.Stats(ctx, p.Client, target.StatsOptions{Query: "k8s.namespace="})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to get the target statistics: ")
	assert.Contains(t, err.Error(), "Failed to parse query")
}
