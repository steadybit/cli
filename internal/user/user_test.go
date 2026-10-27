// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package user_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/steadybit/cli/v6/internal/platformtest"
	"github.com/steadybit/cli/v6/internal/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func TestInvite(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/users/invite", platformtest.Reply{})

	out, err := platformtest.Stdout(t, func() error {
		return user.Invite(ctx, p.Client, user.InviteOptions{Emails: []string{"jane@example.com", "joe@example.com"}, Role: "user", Team: "ADM"})
	})

	require.NoError(t, err)
	assert.Equal(t, "2 user(s) invited. They receive an email to join.\n", out)
	assert.Equal(t, map[string]any{"invitations": []any{
		map[string]any{"email": "jane@example.com", "role": "USER", "teamKey": "ADM"},
		map[string]any{"email": "joe@example.com", "role": "USER", "teamKey": "ADM"},
	}}, p.Requests("POST /api/users/invite")[0].JSON(t))
}

func TestInviteLeavesOutWhatIsNotGiven(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/users/invite", platformtest.Reply{})

	_, err := platformtest.Stdout(t, func() error {
		return user.Invite(ctx, p.Client, user.InviteOptions{Emails: []string{"jane@example.com"}})
	})

	require.NoError(t, err)
	assert.Equal(t, map[string]any{"invitations": []any{map[string]any{"email": "jane@example.com"}}}, p.Requests("POST /api/users/invite")[0].JSON(t))
}

func TestInviteRefusals(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/users/invite", platformtest.Reply{Status: http.StatusForbidden})

	assert.EqualError(t, user.Invite(ctx, p.Client, user.InviteOptions{Emails: []string{"jane@example.com"}}), "Inviting users needs an admin access token.")
	assert.EqualError(t, user.Invite(ctx, p.Client, user.InviteOptions{Emails: []string{"jane"}}), "'jane' is not an email address.")
	assert.EqualError(t, user.Invite(ctx, p.Client, user.InviteOptions{Emails: []string{"jane@example.com"}, Role: "owner"}), "--role must be USER or ADMIN, not 'owner'.")
	assert.EqualError(t, user.Invite(ctx, p.Client, user.InviteOptions{}), "No one to invite. Pass --email.")
	assert.Len(t, p.Requests("POST /api/users/invite"), 1)
}

var users = map[string]any{"items": []any{
	map[string]any{"username": "u-1", "name": "Jane Doe", "email": "jane@example.com", "role": "ADMIN", "teams": []any{"ADM"}, "lastLogin": "2026-10-13T08:12:00Z"},
	map[string]any{"username": "u-2", "name": "Joe", "email": "joe@example.com", "role": "USER", "teams": []any{}, "managedBy": "SCIM"},
}}

func TestList(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/users", platformtest.Reply{JSON: users})

	out, err := platformtest.Stdout(t, func() error { return user.List(ctx, p.Client, user.ListOptions{Query: "jo"}) })

	require.NoError(t, err)
	assert.Contains(t, out, "Jane Doe")
	assert.Contains(t, out, "2026-10-13T08:12:00Z")
	assert.Contains(t, out, "Steadybit")
	assert.Contains(t, out, "SCIM")
	assert.Contains(t, out, "never")
	assert.Equal(t, []string{"jo"}, p.Requests("GET /api/users")[0].Query["query"])
}

func TestListAsJSON(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/users", platformtest.Reply{JSON: users})

	out, err := platformtest.Stdout(t, func() error { return user.List(ctx, p.Client, user.ListOptions{Type: "json"}) })

	require.NoError(t, err)
	assert.Contains(t, out, `"username": "u-2"`)
	assert.Contains(t, out, `"managedBy": "SCIM"`)
	assert.Empty(t, p.Requests("GET /api/users")[0].Query["query"])
}

func TestListNeedsAnAdminToken(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/users", platformtest.Reply{Status: http.StatusForbidden})

	assert.EqualError(t, user.List(ctx, p.Client, user.ListOptions{}), "Listing users needs an admin access token.")
}

func TestDelete(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("DELETE /api/users", platformtest.Reply{Status: http.StatusNoContent})
	// an ADFS subject: base64, so it contains '/', '+' and '='
	username := "Xr2m/ZO0F+ab="

	out, err := platformtest.Stdout(t, func() error { return user.Delete(ctx, p.Client, user.DeleteOptions{Username: username, Yes: true}) })

	require.NoError(t, err)
	assert.Equal(t, "User "+username+" removed from the tenant.\n", out)
	require.Len(t, p.Requests("DELETE /api/users"), 1)
	assert.Equal(t, []string{username}, p.Requests("DELETE /api/users")[0].Query["username"])
}

func TestDeleteRefusals(t *testing.T) {
	p := platformtest.New(t)
	p.Handle("DELETE /api/users", func(r platformtest.Request) platformtest.Reply {
		if q := r.Query["username"]; len(q) == 1 && q[0] == "nobody" {
			return platformtest.Reply{Status: http.StatusNotFound}
		}
		return platformtest.Reply{Status: http.StatusForbidden, JSON: map[string]any{"title": "Forbidden", "detail": "The user is managed by SCIM"}}
	})

	assert.EqualError(t, user.Delete(ctx, p.Client, user.DeleteOptions{Username: "nobody", Yes: true}), "User nobody is not a member of the tenant.")
	err := user.Delete(ctx, p.Client, user.DeleteOptions{Username: "u-2", Yes: true})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to remove user u-2")
	assert.Contains(t, err.Error(), "managed by SCIM")
}
