// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// Every command shows an example in its help: an option list says what can be passed,
// not which combination does what a pipeline author came for.
func TestEveryCommandHasAnExample(t *testing.T) {
	var walk func(*cobra.Command)
	walk = func(cmd *cobra.Command) {
		if cmd.Name() == "help" {
			return
		}
		if !cmd.HasSubCommands() {
			assert.NotEmpty(t, cmd.Example, cmd.CommandPath())
		}
		for _, sub := range cmd.Commands() {
			walk(sub)
		}
	}
	walk(newRoot())
}
