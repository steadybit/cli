// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Profiles written by the TypeScript CLI must keep working after the upgrade.
func TestReadsProfilesWrittenByTheTypeScriptCLI(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	require.NoError(t, os.MkdirAll(filepath.Join(home, ".steadybit"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(home, ".steadybit", "profiles.json"), []byte(`[
  {"name": "prod", "apiAccessToken": "p", "baseUrl": "https://platform.steadybit.com"},
  {"name": "dev", "apiAccessToken": "d", "baseUrl": "https://platform.dev.steadybit.com/"}
]`), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(home, ".steadybit", "activeProfile"), []byte("dev\n"), 0o644))
	os.Unsetenv("STEADYBIT_TOKEN")
	os.Unsetenv("STEADYBIT_URL")

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, Configuration{APIAccessToken: "d", BaseURL: "https://platform.dev.steadybit.com"}, cfg)
}

func TestEnvironmentWinsAndAnEmptyTokenCounts(t *testing.T) {
	setHome(t)
	require.NoError(t, AddProfile(Profile{Name: "p", APIAccessToken: "from-profile"}))
	t.Setenv("STEADYBIT_TOKEN", "")
	t.Setenv("STEADYBIT_URL", "http://localhost:8080")

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, "", cfg.APIAccessToken)
	assert.Equal(t, "http://localhost:8080", cfg.BaseURL)
}

func TestFallsBackToTheFirstProfile(t *testing.T) {
	setHome(t)
	require.NoError(t, AddProfile(Profile{Name: "a", APIAccessToken: "1"}))
	require.NoError(t, AddProfile(Profile{Name: "b", APIAccessToken: "2"}))

	active, err := ActiveProfile()

	require.NoError(t, err)
	assert.Equal(t, "a", active.Name)
}

// Go reads USERPROFILE for the home directory on Windows and HOME elsewhere.
func setHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
}
