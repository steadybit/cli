// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package killswitch_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/steadybit/cli/v6/internal/killswitch"
	"github.com/steadybit/cli/v6/internal/platformtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func TestStatus(t *testing.T) {
	p := platformtest.New(t)
	active := false
	p.Handle("GET /api/killswitch", func(platformtest.Request) platformtest.Reply {
		if !active {
			return platformtest.Reply{Body: `{"active":false}`}
		}
		return platformtest.Reply{Body: `{"active":true,"engagedBy":"u-1","engagedByDetails":{"username":"u-1","name":"Jane Doe"},"engaged":"2026-09-01T10:00:00Z"}`}
	})

	out, err := platformtest.Stdout(t, func() error {
		if err := killswitch.Status(ctx, p.Client, killswitch.StatusOptions{}); err != nil {
			return err
		}
		active = true
		if err := killswitch.Status(ctx, p.Client, killswitch.StatusOptions{}); err != nil {
			return err
		}
		return killswitch.Status(ctx, p.Client, killswitch.StatusOptions{Type: "yaml"})
	})

	require.NoError(t, err)
	assert.Equal(t, "The kill switch is inactive: experiments can run.\n"+
		"The kill switch is active since 2026-09-01T10:00:00Z, activated by Jane Doe: no experiment can run.\n"+
		"active: true\nengagedBy: u-1\nengagedByDetails:\n  username: u-1\n  name: Jane Doe\nengaged: '2026-09-01T10:00:00Z'\n\n", out)
}

func TestActivateAndDeactivate(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/killswitch", platformtest.Reply{})
	p.Reply("DELETE /api/killswitch", platformtest.Reply{})

	out, err := platformtest.Stdout(t, func() error {
		if err := killswitch.Activate(ctx, p.Client, true); err != nil {
			return err
		}
		return killswitch.Deactivate(ctx, p.Client)
	})

	require.NoError(t, err)
	assert.Equal(t, "Kill switch activated. Running experiments are stopped, and none can run until `steadybit killswitch deactivate`.\n"+
		"Kill switch deactivated. Experiments can run again; those it stopped are not restarted.\n", out)
	assert.Len(t, p.Requests("POST /api/killswitch"), 1)
	assert.Len(t, p.Requests("DELETE /api/killswitch"), 1)
}

func TestActivateReportsAFailure(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/killswitch", platformtest.Reply{Status: http.StatusForbidden})

	err := killswitch.Activate(ctx, p.Client, true)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to activate the kill switch: ")
}
