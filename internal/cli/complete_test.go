// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/v6/internal/platformtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompletesTeamsBeforeExperimentKeys(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/teams", platformtest.Reply{JSON: map[string]any{"teams": []any{map[string]any{"key": "ADM", "name": "Admins"}, map[string]any{"key": "SHOP", "name": "Shop"}}}})
	p.Reply("GET /api/experiments", platformtest.Reply{JSON: map[string]any{"experiments": []any{
		map[string]any{"key": "ADM-1", "name": "One"}, map[string]any{"key": "ADM-12", "name": "Twelve"}, map[string]any{"key": "ADM-2", "name": "Two"},
	}}})

	teams, directive := completeExperimentKeys(context.Background(), p.Client, "a")
	assert.Equal(t, []string{"ADM-\tAdmins"}, teams)
	assert.NotZero(t, directive&cobra.ShellCompDirectiveNoSpace)

	keys, _ := completeExperimentKeys(context.Background(), p.Client, "ADM-1")
	assert.Equal(t, []string{"ADM-1\tOne", "ADM-12\tTwelve"}, keys)
	assert.Equal(t, []string{"ADM"}, p.Requests("GET /api/experiments")[0].Query["team"])
}

// Every flag that names something on the platform completes it.
func TestIdFlagsComplete(t *testing.T) {
	root := newRoot()
	for _, path := range [][]string{
		{"experiment", "get", "--key"}, {"experiment", "run", "--template"}, {"schedule", "delete", "--id"},
		{"service", "risk", "--id"}, {"service-profile", "get", "--id"}, {"template", "get", "--id"},
		{"export", "--team"}, {"schedule", "create", "--experiment"},
	} {
		cmd, _, err := root.Find(path[:len(path)-1])
		require.NoError(t, err, strings.Join(path, " "))
		_, ok := cmd.GetFlagCompletionFunc(strings.TrimPrefix(path[len(path)-1], "--"))
		assert.True(t, ok, strings.Join(path, " "))
	}
}

func TestCompletionStaysQuietWithoutAccess(t *testing.T) {
	platformtest.Home(t)
	t.Setenv("STEADYBIT_TOKEN", "")
	root := newRoot()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&bytes.Buffer{})
	root.SetArgs([]string{"__complete", "template", "get", "-i", ""})

	require.NoError(t, root.Execute())
	assert.Equal(t, ":4\n", out.String())
}
