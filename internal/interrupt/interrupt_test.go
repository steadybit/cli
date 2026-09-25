// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package interrupt

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunsHandlersMostRecentFirstUntilPopped(t *testing.T) {
	var calls []string
	popFirst := Push(func(os.Signal) { calls = append(calls, "first") })
	popSecond := Push(func(os.Signal) { calls = append(calls, "second") })

	RunHandlers(os.Interrupt)
	popSecond()
	RunHandlers(os.Interrupt)
	popFirst()
	RunHandlers(os.Interrupt)

	assert.Equal(t, []string{"second", "first", "first"}, calls)
}
