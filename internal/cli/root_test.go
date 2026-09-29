// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Completion never waits for the release check, wherever global flags put it.
func TestCompletionSkipsTheUpdateCheck(t *testing.T) {
	assert.True(t, completing([]string{"completion", "zsh"}))
	assert.True(t, completing([]string{"--profile", "prod", "completion", "zsh"}))
	assert.True(t, completing([]string{"-v", "__complete", "experiment", "get", "-k", ""}))
	assert.True(t, completing([]string{"__completeNoDesc", "team"}))
	assert.False(t, completing([]string{"experiment", "run", "-k", "ADM-1"}))
}
