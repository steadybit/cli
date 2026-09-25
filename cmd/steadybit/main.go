// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package main

import (
	"os"

	"github.com/steadybit/cli/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
