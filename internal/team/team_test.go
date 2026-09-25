// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package team_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/steadybit/cli/internal/platformtest"
	"github.com/steadybit/cli/internal/team"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

const stored = `{"id":"0190d7b2-0000-7000-8000-000000000001","key":"OPS","name":"Operations","version":7,"allowedActions":[],"allowedEnvironments":["Global"],"managedBy":"MANUAL",` +
	`"members":[{"username":"u-1","name":"Jane","pictureUrl":"p","email":"jane@example.com","role":"OWNER","managedBy":"MANUAL"}]}`

func TestList(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/teams", platformtest.Reply{Body: `{"teams":[` + stored + `]}`})

	out, err := platformtest.Stdout(t, func() error { return team.List(ctx, p.Client, team.ListOptions{Search: "ops"}) })

	require.NoError(t, err)
	assert.Contains(t, out, "│ OPS │ Operations │       1 │            1 │")
	assert.Equal(t, []string{"ops"}, p.Requests("GET /api/teams")[0].Query["search"])
}

func TestGetAndApplyRoundTripByKey(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/teams/OPS", platformtest.Reply{Body: stored})
	p.Reply("POST /api/teams", platformtest.Reply{Body: stored})
	file := filepath.Join(t.TempDir(), "team.yml")

	out, err := platformtest.Stdout(t, func() error {
		if err := team.Get(ctx, p.Client, team.GetOptions{Key: "OPS", File: file}); err != nil {
			return err
		}
		return team.Apply(ctx, p.Client, team.ApplyOptions{Files: []string{file}})
	})

	require.NoError(t, err)
	assert.Contains(t, out, "Team OPS updated.")
	before, _ := os.ReadFile(file)
	assert.NotContains(t, string(before), "version")
	assert.NotContains(t, string(before), "pictureUrl")
	sent := p.Requests("POST /api/teams")[0].JSON(t).(map[string]any)
	assert.Equal(t, []any{map[string]any{"username": "u-1", "email": "jane@example.com", "role": "OWNER"}}, sent["members"])
	assert.Equal(t, "MANUAL", sent["managedBy"])
}

func TestApplyWritesNothingBack(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/teams", platformtest.Reply{Status: http.StatusCreated, Body: stored})
	file := filepath.Join(t.TempDir(), "team.yml")
	original := "key: OPS\nname: Operations\nallowedActions: []\nallowedEnvironments: []\n"
	require.NoError(t, os.WriteFile(file, []byte(original), 0o644))

	out, err := platformtest.Stdout(t, func() error { return team.Apply(ctx, p.Client, team.ApplyOptions{Files: []string{file}}) })

	require.NoError(t, err)
	assert.Equal(t, "Team OPS created.\n", out)
	content, _ := os.ReadFile(file)
	assert.Equal(t, original, string(content))
}

func TestDeletePurgesOnlyWhenAsked(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("DELETE /api/teams/OPS", platformtest.Reply{Body: stored})
	p.Reply("DELETE /api/teams/NOPE", platformtest.Reply{Status: http.StatusNotFound})

	out, err := platformtest.Stdout(t, func() error {
		if err := team.Delete(ctx, p.Client, team.DeleteOptions{Key: "OPS", Yes: true}); err != nil {
			return err
		}
		return team.Delete(ctx, p.Client, team.DeleteOptions{Key: "OPS", Experiments: true, Yes: true})
	})

	require.NoError(t, err)
	assert.Equal(t, "Team OPS deleted.\nTeam OPS deleted.\n", out)
	requests := p.Requests("DELETE /api/teams/OPS")
	assert.Equal(t, []string{"false"}, requests[0].Query["purgeIncludingExperiments"])
	assert.Equal(t, []string{"true"}, requests[1].Query["purgeIncludingExperiments"])
	assert.EqualError(t, team.Delete(ctx, p.Client, team.DeleteOptions{Key: "NOPE", Yes: true}), "Team NOPE not found.")
}

func TestMembers(t *testing.T) {
	p := platformtest.New(t)
	result := platformtest.Reply{Body: `{"members":[{"username":"u-1","name":"Jane","email":"jane@example.com","role":"OWNER"}]}`}
	p.Reply("GET /api/teams/OPS/members", result)
	p.Reply("POST /api/teams/OPS/members/add", result)
	p.Reply("POST /api/teams/OPS/members/remove", result)
	p.Reply("PUT /api/teams/OPS/members", result)

	out, err := platformtest.Stdout(t, func() error {
		if err := team.ListMembers(ctx, p.Client, "OPS"); err != nil {
			return err
		}
		if err := team.AddMembers(ctx, p.Client, team.MemberOptions{Key: "OPS", Usernames: []string{"u-1"}, Emails: []string{"joe@example.com"}, Role: "owner"}); err != nil {
			return err
		}
		if err := team.RemoveMembers(ctx, p.Client, team.MemberOptions{Key: "OPS", Emails: []string{"joe@example.com"}, Yes: true}); err != nil {
			return err
		}
		return team.SetMembers(ctx, p.Client, team.MemberOptions{Key: "OPS", Usernames: []string{"u-1"}, Role: "OWNER", Validate: true, Yes: true})
	})

	require.NoError(t, err)
	assert.Contains(t, out, "│ u-1      │ Jane │ jane@example.com │ OWNER │")
	assert.Contains(t, out, "Team OPS now has 1 member(s).")
	add := p.Requests("POST /api/teams/OPS/members/add")[0]
	assert.Equal(t, map[string]any{"members": []any{map[string]any{"username": "u-1", "role": "OWNER"}, map[string]any{"email": "joe@example.com", "role": "OWNER"}}}, add.JSON(t))
	assert.Equal(t, []string{"false"}, add.Query["validateMembers"])
	assert.Equal(t, map[string]any{"emails": []any{"joe@example.com"}}, p.Requests("POST /api/teams/OPS/members/remove")[0].JSON(t))
	set := p.Requests("PUT /api/teams/OPS/members")[0]
	assert.Equal(t, map[string]any{"members": []any{map[string]any{"username": "u-1", "role": "OWNER"}}}, set.JSON(t))
	assert.Equal(t, []string{"true"}, set.Query["validateMembers"])

	assert.EqualError(t, team.AddMembers(ctx, p.Client, team.MemberOptions{Key: "OPS"}), "No members given. Pass --username or --email.")
	assert.EqualError(t, team.AddMembers(ctx, p.Client, team.MemberOptions{Key: "OPS", Usernames: []string{"x"}, Role: "admin"}), "--role must be MEMBER or OWNER, not 'admin'.")
}

func TestEnvironments(t *testing.T) {
	p := platformtest.New(t)
	result := platformtest.Reply{Body: `{"environments":[{"name":"Global"},{"name":"Prod"}]}`}
	p.Reply("GET /api/teams/OPS/environments", result)
	p.Reply("POST /api/teams/OPS/environments/add", result)
	p.Reply("POST /api/teams/OPS/environments/remove", result)
	p.Reply("PUT /api/teams/OPS/environments", result)
	p.Reply("GET /api/teams/NOPE/environments", platformtest.Reply{Status: http.StatusNotFound})

	out, err := platformtest.Stdout(t, func() error {
		if err := team.ListEnvironments(ctx, p.Client, "OPS"); err != nil {
			return err
		}
		o := team.EnvironmentOptions{Key: "OPS", Environments: []string{"Prod"}, Yes: true}
		if err := team.AddEnvironments(ctx, p.Client, o); err != nil {
			return err
		}
		if err := team.RemoveEnvironments(ctx, p.Client, o); err != nil {
			return err
		}
		return team.SetEnvironments(ctx, p.Client, o)
	})

	require.NoError(t, err)
	assert.Contains(t, out, "│ Prod        │")
	assert.Contains(t, out, "Team OPS now has 2 environment(s): Global, Prod.")
	want := map[string]any{"environments": []any{map[string]any{"name": "Prod"}}}
	assert.Equal(t, want, p.Requests("POST /api/teams/OPS/environments/add")[0].JSON(t))
	assert.Equal(t, want, p.Requests("POST /api/teams/OPS/environments/remove")[0].JSON(t))
	assert.Equal(t, want, p.Requests("PUT /api/teams/OPS/environments")[0].JSON(t))
	assert.EqualError(t, team.ListEnvironments(ctx, p.Client, "NOPE"), "Team NOPE not found.")
}
