// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package target_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/steadybit/cli/v6/internal/platformtest"
	"github.com/steadybit/cli/v6/internal/target"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func container(name string, attributes ...[2]string) map[string]any {
	list := []any{}
	for _, a := range attributes {
		list = append(list, map[string]any{"key": a[0], "value": a[1]})
	}
	return map[string]any{"name": name, "type": "container", "attributes": list}
}

func TestQueryFollowsTheCursorUpToTheLimit(t *testing.T) {
	p := platformtest.New(t)
	p.Handle("GET /api/targets", func(r platformtest.Request) platformtest.Reply {
		if r.Query["cursor"] == nil {
			return platformtest.Reply{JSON: map[string]any{"items": []any{
				container("a", [2]string{"k8s.namespace", "shop"}, [2]string{"k8s.pod.name", "a-1"}, [2]string{"k8s.pod.name", "a-2"}),
				container("b"),
			}, "hasNext": true, "nextCursor": "c-2"}}
		}
		return platformtest.Reply{JSON: map[string]any{"items": []any{container("c", [2]string{"k8s.namespace", "shop"})}, "hasNext": true, "nextCursor": "c-3"}}
	})

	out, err := platformtest.Stdout(t, func() error {
		return target.Query(ctx, p.Client, target.QueryOptions{
			Environment: "Global", TargetType: "container", Query: `k8s.namespace="shop"`, Attributes: []string{"k8s.namespace", "k8s.pod.name"}, Limit: 3,
		})
	})

	require.NoError(t, err)
	assert.Contains(t, out, "│ a    │ container │ shop          │ a-1, a-2     │")
	assert.Contains(t, out, "│ c    │ container │ shop          │              │")
	requests := p.Requests("GET /api/targets")
	require.Len(t, requests, 2)
	assert.Equal(t, []string{"Global"}, requests[0].Query["environment"])
	assert.Equal(t, []string{"container"}, requests[0].Query["targetType"])
	assert.Equal(t, []string{`k8s.namespace="shop"`}, requests[0].Query["query"])
	assert.Equal(t, []string{"k8s.namespace", "k8s.pod.name"}, requests[0].Query["attribute"])
	assert.Equal(t, []string{"3"}, requests[0].Query["size"])
	assert.Equal(t, []string{"1"}, requests[1].Query["size"])
	assert.Equal(t, []string{"c-2"}, requests[1].Query["cursor"])
}

func TestQueryPrintsAsJSON(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/targets", platformtest.Reply{JSON: map[string]any{"items": []any{}, "hasNext": false}})

	out, err := platformtest.Stdout(t, func() error {
		return target.Query(ctx, p.Client, target.QueryOptions{Environment: "Global", Type: "json"})
	})

	require.NoError(t, err)
	assert.Equal(t, "[]\n", out)
	assert.Equal(t, []string{"1000"}, p.Requests("GET /api/targets")[0].Query["size"])
	assert.NotContains(t, p.Requests("GET /api/targets")[0].Query, "targetType")
}

func TestQueryReportsAMissingEnvironment(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/targets", platformtest.Reply{Status: http.StatusNotFound})

	assert.EqualError(t, target.Query(ctx, p.Client, target.QueryOptions{Environment: "Nowhere"}), "Environment Nowhere not found.")
}

func TestAttributeKeysAndValues(t *testing.T) {
	p := platformtest.New(t)
	p.Handle("GET /api/targets/attributes/keys", func(r platformtest.Request) platformtest.Reply {
		if r.Query["page"][0] == "0" {
			return platformtest.Reply{JSON: map[string]any{"items": []any{"k8s.namespace"}, "nextPage": 1}}
		}
		return platformtest.Reply{JSON: map[string]any{"items": []any{"k8s.pod.name"}}}
	})
	p.Reply("GET /api/targets/attributes/values", platformtest.Reply{JSON: map[string]any{"items": []any{"shop", "kube-system"}}})

	out, err := platformtest.Stdout(t, func() error {
		if err := target.AttributeKeys(ctx, p.Client, target.AttributeOptions{Environment: "Global", TargetType: "container"}); err != nil {
			return err
		}
		return target.AttributeValues(ctx, p.Client, target.AttributeOptions{Environment: "Global", Action: "com.example.attack", Key: "k8s.namespace", Type: "yaml"})
	})

	require.NoError(t, err)
	assert.Contains(t, out, "│ k8s.namespace │\n│ k8s.pod.name  │")
	assert.Contains(t, out, "- shop\n- kube-system\n")
	values := p.Requests("GET /api/targets/attributes/values")[0].Query
	assert.Equal(t, []string{"k8s.namespace"}, values["attributeKey"])
	assert.Equal(t, []string{"com.example.attack"}, values["actionId"])
	assert.EqualError(t, target.AttributeKeys(ctx, p.Client, target.AttributeOptions{Environment: "Global"}), "Either --target-type or --action must be specified.")
}
