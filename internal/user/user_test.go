// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package user_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/steadybit/cli/internal/platformtest"
	"github.com/steadybit/cli/internal/user"
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
