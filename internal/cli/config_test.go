// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"os"
	"strings"
	"testing"

	"github.com/steadybit/cli/v6/internal/config"
	"github.com/steadybit/cli/v6/internal/platformtest"
	"github.com/steadybit/cli/v6/internal/prompt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// profiles gives the test a home of its own, with these profiles in it.
func profiles(t *testing.T, names ...string) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("USERPROFILE", os.Getenv("HOME"))
	for _, name := range names {
		require.NoError(t, config.AddProfile(config.Profile{Name: name, BaseURL: "https://example.com", APIAccessToken: "t-" + name}))
	}
}

// answering scripts the prompts, as on a terminal.
func answering(t *testing.T, input string) {
	original := prompt.Interactive
	prompt.Interactive = func() bool { return true }
	t.Cleanup(func() { prompt.Interactive = original; prompt.UseInput(os.Stdin) })
	prompt.UseInput(strings.NewReader(input))
}

func run(t *testing.T, args ...string) (string, error) {
	t.Helper()
	return platformtest.Stdout(t, func() error {
		root := newRoot()
		root.SetArgs(args)
		root.SilenceErrors = true
		return root.Execute()
	})
}

func names(t *testing.T) []string {
	t.Helper()
	list, err := config.Profiles()
	require.NoError(t, err)
	var out []string
	for _, p := range list {
		out = append(out, p.Name)
	}
	return out
}

// The options given are kept, and only the missing ones are asked for.
func TestProfileAddAsksOnlyForWhatIsMissing(t *testing.T) {
	profiles(t)
	answering(t, "secret\n")

	out, err := run(t, "config", "profile", "add", "-n", "dev", "-b", "https://dev.example.com")

	require.NoError(t, err)
	assert.NotContains(t, out, "Profile name:")
	assert.NotContains(t, out, "Base URL")
	assert.Contains(t, out, "API access token:")
	list, _ := config.Profiles()
	assert.Equal(t, []config.Profile{{Name: "dev", BaseURL: "https://dev.example.com", APIAccessToken: "secret"}}, list)
}

// `-t "$STEADYBIT_TOKEN"` with the variable unset.
func TestProfileAddRefusesAnEmptyToken(t *testing.T) {
	profiles(t)

	_, err := run(t, "config", "profile", "add", "-n", "dev", "-t", "")

	assert.EqualError(t, err, "The token given with --token is empty. Is the variable it comes from set?")
	assert.Empty(t, names(t))
}

// A mistyped name would otherwise fall back to the first profile, possibly production.
func TestProfileSelectRefusesAnUnknownName(t *testing.T) {
	profiles(t, "prod", "dev")

	_, err := run(t, "config", "profile", "select", "dve")
	assert.EqualError(t, err, "No profile named dve. Available: prod, dev")

	out, err := run(t, "config", "profile", "select", "dev")
	require.NoError(t, err)
	assert.Equal(t, "Profile dev is now active.\n", out)
	active, _ := config.ActiveProfile()
	assert.Equal(t, "dev", active.Name)
}

func TestProfileRemoveAsksFirst(t *testing.T) {
	profiles(t, "prod", "dev")
	answering(t, "n\ny\n")

	out, err := run(t, "config", "profile", "remove", "dev")
	require.NoError(t, err)
	assert.Contains(t, out, "Remove profile dev?")
	assert.Equal(t, []string{"prod", "dev"}, names(t))

	out, err = run(t, "config", "profile", "remove", "dev")
	require.NoError(t, err)
	assert.Contains(t, out, "Profile dev removed.\n")
	assert.Equal(t, []string{"prod"}, names(t))
}

func TestProfileRemoveSaysWhenItWasTheActiveOne(t *testing.T) {
	profiles(t, "prod", "dev")
	require.NoError(t, config.SetActiveProfile("dev"))

	out, err := run(t, "config", "profile", "remove", "dev", "--yes")

	require.NoError(t, err)
	assert.Equal(t, "Profile dev removed.\nIt was the active profile; choose another with `steadybit config profile select`.\n", out)
}

func TestProfileListSaysWhenThereAreNone(t *testing.T) {
	profiles(t)

	out, err := run(t, "config", "profile", "list")

	require.NoError(t, err)
	assert.Empty(t, out)
}
