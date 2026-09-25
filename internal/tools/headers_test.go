// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package tools

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var header = regexp.MustCompile(`^(#!.*\n)?// SPDX-License-Identifier: MIT\n// SPDX-FileCopyrightText: \d{4} Steadybit GmbH\n`)

// Every hand-written source file starts with the SPDX header; generated ones are exempt.
func TestSourceFilesCarryTheSPDXHeader(t *testing.T) {
	root := "../.."
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && (d.Name() == "node_modules" || d.Name() == "dist" || strings.HasPrefix(d.Name(), ".")) && path != root {
			return filepath.SkipDir
		}
		if d.IsDir() || !(strings.HasSuffix(path, ".go") || strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".mjs")) || strings.HasSuffix(path, ".gen.go") {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !header.Match(content) {
			t.Errorf("%s does not start with the SPDX header", path)
		}
		return nil
	})
}
