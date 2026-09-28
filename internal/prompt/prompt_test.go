// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package prompt

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func interactive(t *testing.T, is bool) {
	original := Interactive
	Interactive = func() bool { return is }
	t.Cleanup(func() { Interactive = original })
}

// With output redirected nobody sees the question, so it is not asked, even though
// stdin could still answer it.
func TestConfirmDoesNotAskWhenTheQuestionWouldNotBeSeen(t *testing.T) {
	interactive(t, false)
	UseInput(strings.NewReader("n\n"))

	ok, err := Confirm("Sure?", false, true)

	require.NoError(t, err)
	assert.True(t, ok)
}

func TestConfirmReadsTheAnswer(t *testing.T) {
	interactive(t, true)
	UseInput(strings.NewReader("y\n\n"))

	yes, err := Confirm("Sure?", false, false)
	require.NoError(t, err)
	byDefault, err := Confirm("Sure?", false, true)
	require.NoError(t, err)

	assert.True(t, yes)
	assert.False(t, byDefault)
}
