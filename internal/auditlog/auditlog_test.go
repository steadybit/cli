// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package auditlog_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/steadybit/cli/internal/auditlog"
	"github.com/steadybit/cli/internal/platformtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

const entries = `[
 {"id":"a","eventName":"experiment.created","eventTime":"2026-09-01T10:00:00Z","principal":{"principalType":"USER","name":"Jane Doe","username":"u-1"},"tenant":{"key":"demo","name":"Demo"},"team":{"id":"t","key":"ADM","name":"Admins"}},
 {"id":"b","eventName":"environment.deleted","eventTime":"2026-09-02T11:00:00Z","principal":{"principalType":"ACCESS_TOKEN","name":"CI/CD","tokenType":"ADMIN"},"tenant":{"key":"demo","name":"Demo"},"environment":{"id":"e","name":"Prod","predicate":{}}},
 {"id":"c","eventName":"advice.updated","eventTime":"2026-09-03T12:00:00Z","principal":{"principalType":"BATCH_JOB","username":"system"},"tenant":{"key":"demo","name":"Demo"}}
]`

func TestShowsWhoDidWhat(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/audit-log", platformtest.Reply{Body: entries})

	out, err := platformtest.Stdout(t, func() error {
		return auditlog.Show(ctx, p.Client, auditlog.Options{From: "2026-09-01", To: "2026-09-08T12:00:00+02:00"})
	})

	require.NoError(t, err)
	assert.Contains(t, out, "│ 2026-09-01T10:00:00Z │ experiment.created  │ Jane Doe    │ ADM  │             │")
	assert.Contains(t, out, "│ 2026-09-02T11:00:00Z │ environment.deleted │ token CI/CD │      │ Prod        │")
	assert.Contains(t, out, "│ 2026-09-03T12:00:00Z │ advice.updated      │ batch job   │      │             │")
	q := p.Requests("GET /api/audit-log")[0].Query
	assert.Equal(t, []string{"2026-09-01T00:00:00Z"}, q["from"])
	assert.Equal(t, []string{"2026-09-08T12:00:00+02:00"}, q["to"])
}

func TestPrintsTheEntriesAsJSON(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/audit-log", platformtest.Reply{Body: `[{"id":"a","eventName":"x"}]`})

	out, err := platformtest.Stdout(t, func() error { return auditlog.Show(ctx, p.Client, auditlog.Options{Type: "json"}) })

	require.NoError(t, err)
	assert.Equal(t, "[\n  {\n    \"id\": \"a\",\n    \"eventName\": \"x\"\n  }\n]\n", out)
	assert.Empty(t, p.Requests("GET /api/audit-log")[0].Query)
}

func TestRefusals(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/audit-log", platformtest.Reply{Status: http.StatusForbidden})

	assert.EqualError(t, auditlog.Show(ctx, p.Client, auditlog.Options{}), "The audit log needs an admin access token.")
	assert.EqualError(t, auditlog.Show(ctx, p.Client, auditlog.Options{From: "last week"}), "--from 'last week' is neither a date like 2026-09-01 nor a time like 2026-09-01T12:00:00Z.")
}
