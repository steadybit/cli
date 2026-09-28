// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

//go:build !windows

package output

import "os"

// Terminals elsewhere interpret escape codes as they are.
func enableEscapeCodes(*os.File) bool { return true }
