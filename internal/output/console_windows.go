// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package output

import (
	"os"

	"golang.org/x/sys/windows"
)

// Colours on stderr are gated on stdout, so both have to interpret escape codes.
func enableEscapeCodes(stdout *os.File) bool {
	ok := enable(stdout)
	enable(os.Stderr)
	return ok
}

func enable(f *os.File) bool {
	handle := windows.Handle(f.Fd())
	var mode uint32
	if err := windows.GetConsoleMode(handle, &mode); err != nil {
		return false
	}
	return windows.SetConsoleMode(handle, mode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING) == nil
}
