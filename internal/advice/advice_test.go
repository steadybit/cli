// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package advice_test

import (
	"context"
	"testing"

	"github.com/steadybit/cli/v6/internal/advice"
	"github.com/steadybit/cli/v6/internal/platformtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func item(ref, label, status string) map[string]any {
	return map[string]any{"target": map[string]any{"reference": ref}, "advice": map[string]any{"label": label, "status": status}}
}

func TestPagesThroughAdviceAndReportsMismatches(t *testing.T) {
	p := platformtest.New(t)
	p.Handle("POST /api/advice", func(r platformtest.Request) platformtest.Reply {
		if r.JSON(t).(map[string]any)["offset"] == float64(0) {
			return platformtest.Reply{JSON: map[string]any{"totalItems": 3, "nextOffset": 2, "items": []any{item("t-1", "a-1", "IMPLEMENTED"), item("t-2", "a-2", "ACTION_NEEDED")}}}
		}
		return platformtest.Reply{JSON: map[string]any{"totalItems": 3, "items": []any{item("t-3", "a-3", "IMPLEMENTED")}}}
	})

	out, err := platformtest.Stdout(t, func() error {
		return advice.ValidateStatus(context.Background(), p.Client, advice.Options{Environment: "Global", Query: "a=b", Status: "Implemented"})
	})

	assert.EqualError(t, err, "1 of 3 advice did not match the expected status.")
	assert.Contains(t, out, "Fetched 2 of 3 matching advice.\nFetched 3 of 3 matching advice.\n")
	first := p.Requests("POST /api/advice")[0].JSON(t)
	assert.Equal(t, map[string]any{"environmentName": "Global", "offset": float64(0), "query": "a=b"}, first)
}

// The platform reports IMPLEMENTED while the flag defaults to Implemented; case and the
// separator are ignored, so `action needed` matches too.
func TestMatchesStatusesRegardlessOfCaseAndSeparator(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/advice", platformtest.Reply{JSON: map[string]any{"totalItems": 1, "items": []any{item("t", "a", "ACTION_NEEDED")}}})

	_, err := platformtest.Stdout(t, func() error {
		return advice.ValidateStatus(context.Background(), p.Client, advice.Options{Environment: "Global", Status: "action needed"})
	})

	require.NoError(t, err)
}
