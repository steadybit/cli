// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package main

import (
	"os"

	"github.com/steadybit/cli/v6/internal/cli"
	// Imported for its signal handling, which every command needs from the start.
	_ "github.com/steadybit/cli/v6/internal/interrupt"
)

func main() {
	os.Exit(cli.Execute())
}
