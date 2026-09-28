// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package accesstoken_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/steadybit/cli/v6/internal/accesstoken"
	"github.com/steadybit/cli/v6/internal/platformtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

const id = "Xy12AbCd"

func TestListWalksEveryPageWithTheFilters(t *testing.T) {
	p := platformtest.New(t)
	p.Handle("GET /api/access-tokens/v2", func(r platformtest.Request) platformtest.Reply {
		if r.Query["page"][0] == "0" {
			return platformtest.Reply{JSON: map[string]any{"items": []any{map[string]any{"id": id, "name": "ci", "type": "TEAM", "teams": []any{"ADM", "OPS"}}}, "nextPage": 1}}
		}
		return platformtest.Reply{JSON: map[string]any{"items": []any{map[string]any{"id": "b", "name": "admin", "type": "ADMIN", "expiresAt": "2026-12-31T00:00:00Z", "lastUsed": "2026-09-01T10:00:00Z"}}}}
	})
	no := false

	out, err := platformtest.Stdout(t, func() error {
		return accesstoken.List(ctx, p.Client, accesstoken.ListOptions{Type: "team", Teams: []string{"ADM", "OPS"}, Expired: &no})
	})

	require.NoError(t, err)
	assert.Contains(t, out, "│ "+id+" │ ci    │ TEAM  │ ADM, OPS │ never                │ never                │")
	assert.Contains(t, out, "│ b        │ admin │ ADMIN │          │ 2026-12-31T00:00:00Z │ 2026-09-01T10:00:00Z │")
	q := p.Requests("GET /api/access-tokens/v2")[0].Query
	assert.Equal(t, []string{"TEAM"}, q["type"])
	assert.Equal(t, []string{"ADM", "OPS"}, q["teams"])
	assert.Equal(t, []string{"false"}, q["expired"])
	assert.Equal(t, []string{"100"}, q["size"])
	assert.EqualError(t, accesstoken.List(ctx, p.Client, accesstoken.ListOptions{Type: "user"}), "--type must be ADMIN, TEAM or WILDCARD, not 'user'.")
}

func TestCreatePrintsTheTokenOnce(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/access-tokens/v2", platformtest.Reply{JSON: map[string]any{"id": id, "token": "secret-value"}})

	out, err := platformtest.Stdout(t, func() error {
		return accesstoken.Create(ctx, p.Client, accesstoken.CreateOptions{Name: "ci", Type: "team", Teams: []string{"ADM"}, ExpiresAt: "2026-12-31"})
	})

	require.NoError(t, err)
	assert.Equal(t, "Access token ci ("+id+") created. Store it now, it cannot be shown again:\nsecret-value\n", out)
	assert.Equal(t, map[string]any{"name": "ci", "type": "TEAM", "teams": []any{"ADM"}, "expiresAt": "2026-12-31T00:00:00Z"}, p.Requests("POST /api/access-tokens/v2")[0].JSON(t))

	out, err = platformtest.Stdout(t, func() error {
		return accesstoken.Create(ctx, p.Client, accesstoken.CreateOptions{Name: "ci", Type: "ADMIN", Output: "json"})
	})
	require.NoError(t, err)
	assert.Equal(t, "{\n  \"id\": \""+id+"\",\n  \"token\": \"secret-value\"\n}\n", out)
	assert.Equal(t, map[string]any{"name": "ci", "type": "ADMIN"}, p.Requests("POST /api/access-tokens/v2")[1].JSON(t))

	assert.EqualError(t, accesstoken.Create(ctx, p.Client, accesstoken.CreateOptions{Name: "ci", Type: "ADMIN", ExpiresAt: "tomorrow"}),
		"--expires-at 'tomorrow' is neither a date like 2026-09-01 nor a time like 2026-09-01T12:00:00Z.")
}

func TestRecreateAndDelete(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/access-tokens/v2/"+id+"/recreate", platformtest.Reply{JSON: map[string]any{"id": "new-id", "token": "new-secret"}})
	p.Reply("DELETE /api/access-tokens/v2/"+id, platformtest.Reply{})
	p.Reply("DELETE /api/access-tokens/v2/nope", platformtest.Reply{Status: http.StatusNotFound})

	out, err := platformtest.Stdout(t, func() error {
		if err := accesstoken.Recreate(ctx, p.Client, accesstoken.RecreateOptions{ID: id, ExpiresAt: "2027-06-30T12:00:00+02:00", Yes: true}); err != nil {
			return err
		}
		return accesstoken.Delete(ctx, p.Client, accesstoken.DeleteOptions{ID: id, Yes: true})
	})

	require.NoError(t, err)
	assert.Equal(t, "Access token "+id+" recreated as new-id. Store it now, it cannot be shown again:\nnew-secret\nAccess token "+id+" deleted.\n", out)
	assert.Equal(t, map[string]any{"expiresAt": "2027-06-30T12:00:00+02:00"}, p.Requests("POST /api/access-tokens/v2/" + id + "/recreate")[0].JSON(t))
	assert.EqualError(t, accesstoken.Delete(ctx, p.Client, accesstoken.DeleteOptions{ID: "nope", Yes: true}), "Access token nope not found.")
}

// Recreated without --expires-at, a token expires when the one it replaces would have,
// not never, which the platform would make it.
func TestRecreateKeepsTheExpiry(t *testing.T) {
	future := time.Now().Add(48 * time.Hour).UTC().Truncate(time.Second)
	past := time.Now().Add(-48 * time.Hour).UTC().Truncate(time.Second)
	p := platformtest.New(t)
	p.Reply("GET /api/access-tokens/v2", platformtest.Reply{JSON: map[string]any{"items": []any{
		map[string]any{"id": "soon", "expiresAt": future.Format(time.RFC3339)},
		map[string]any{"id": "gone", "expiresAt": past.Format(time.RFC3339)},
		map[string]any{"id": "forever"},
	}}})
	for _, token := range []string{"soon", "forever", "gone"} {
		p.Reply("POST /api/access-tokens/v2/"+token+"/recreate", platformtest.Reply{JSON: map[string]any{"id": "new", "token": "s"}})
	}
	recreate := func(token, expiresAt string) error {
		_, err := platformtest.Stdout(t, func() error {
			return accesstoken.Recreate(ctx, p.Client, accesstoken.RecreateOptions{ID: token, ExpiresAt: expiresAt, Yes: true})
		})
		return err
	}

	require.NoError(t, recreate("soon", ""))
	require.NoError(t, recreate("forever", ""))
	assert.EqualError(t, recreate("gone", ""), "Access token gone expired at "+past.Format(time.RFC3339)+". Pass --expires-at for the new one, or --expires-at never.")
	require.NoError(t, recreate("gone", "never"))
	assert.EqualError(t, recreate("missing", ""), "Access token missing not found.")

	assert.Equal(t, map[string]any{"expiresAt": future.Format(time.RFC3339)}, p.Requests("POST /api/access-tokens/v2/soon/recreate")[0].JSON(t))
	assert.Equal(t, map[string]any{}, p.Requests("POST /api/access-tokens/v2/forever/recreate")[0].JSON(t))
	assert.Equal(t, map[string]any{}, p.Requests("POST /api/access-tokens/v2/gone/recreate")[0].JSON(t))
}
