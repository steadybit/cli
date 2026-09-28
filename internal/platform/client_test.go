// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package platform_test

import (
	"context"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/steadybit/cli/v6/internal/platform"
	"github.com/steadybit/cli/v6/internal/platformtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	platform.RetryUnit = time.Millisecond
}

func TestSendsTheTokenAndUserAgent(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/teams", platformtest.Reply{JSON: map[string]any{"teams": []any{}}})

	_, _, err := platform.Read(p.Client.GetTeams(context.Background(), nil))

	require.NoError(t, err)
	r := p.Requests("GET /api/teams")[0]
	assert.Equal(t, "accessToken test-token", r.Header.Get("Authorization"))
	assert.True(t, strings.HasPrefix(r.Header.Get("User-Agent"), "steadybit@"))
}

func TestWaitsOutRateLimitsForAnyMethod(t *testing.T) {
	p := platformtest.New(t)
	var calls atomic.Int32
	p.Handle("POST /api/experiments", func(platformtest.Request) platformtest.Reply {
		if calls.Add(1) <= 3 {
			return platformtest.Reply{Status: http.StatusTooManyRequests, Headers: map[string]string{"RateLimit-Reset": "1"}}
		}
		return platformtest.Reply{Status: http.StatusCreated}
	})

	_, resp, err := platform.Read(p.Client.CreateOrUpdateExperimentWithBody(context.Background(), "application/json", strings.NewReader("{}")))

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.EqualValues(t, 4, calls.Load())
	// The body is sent again on every attempt, not only the first.
	for _, r := range p.Requests("POST /api/experiments") {
		assert.Equal(t, "{}", string(r.Body))
	}
}

func TestGivesUpOnRateLimitsOnceTheBudgetIsSpent(t *testing.T) {
	p := platformtest.New(t)
	original := platform.MaxRateLimitWait
	platform.MaxRateLimitWait = 5 * time.Millisecond
	t.Cleanup(func() { platform.MaxRateLimitWait = original })
	p.Reply("GET /api/teams", platformtest.Reply{Status: http.StatusTooManyRequests, Body: "slow down"})

	_, _, err := platform.Read(p.Client.GetTeams(context.Background(), nil))

	assert.ErrorContains(t, err, "responded with unexpected status code: 429 - slow down")
}

// A POST that failed in transit may still have started a run, so only methods defined to
// be idempotent are repeated.
func TestRetriesTransportFailuresOnlyForIdempotentMethods(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	var accepted atomic.Int32
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			accepted.Add(1)
			_ = conn.Close() // every connection fails before a response
		}
	}()
	t.Cleanup(func() { _ = listener.Close() })
	platformtest.Home(t)
	t.Setenv("STEADYBIT_URL", "http://"+listener.Addr().String())
	t.Setenv("STEADYBIT_TOKEN", "t")
	client, err := platform.New()
	require.NoError(t, err)

	_, _, err = platform.Read(client.GetTeams(context.Background(), nil))
	assert.ErrorContains(t, err, "Failed to call Steadybit API at GET")
	gets := accepted.Load()

	_, _, err = platform.Read(client.CreateOrUpdateExperimentWithBody(context.Background(), "application/json", strings.NewReader("{}")))
	assert.Error(t, err)

	assert.EqualValues(t, 4, gets)
	assert.EqualValues(t, 1, accepted.Load()-gets)
}

func TestReportsTheProblemBodyAsTheTypeScriptCLIDid(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/teams", platformtest.Reply{Status: 422, Body: `{"title":"bad","violations":[]}`})

	_, _, err := platform.Read(p.Client.GetTeams(context.Background(), nil))
	err = platform.Failed(err, "Failed to get %s", "teams")

	assert.Equal(t, "Failed to get teams: Steadybit API at GET "+p.URL+`/api/teams responded with unexpected status code: 422 - {"title":"bad","violations":[]}: {
  "title": "bad",
  "violations": []
}`, err.Error())
}
