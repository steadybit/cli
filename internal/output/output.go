// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package output formats documents and gates colour on stdout being a terminal, as the
// TypeScript CLI did: its output is routinely parsed by GitOps pipelines.
package output

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"
)

// stdoutTerminal is whether stdout is a terminal that interprets escape codes. The classic
// Windows console only does once asked to, as Node asked on start.
var stdoutTerminal = term.IsTerminal(int(os.Stdout.Fd())) && enableEscapeCodes(os.Stdout)

var colorsEnabled = os.Getenv("NO_COLOR") == "" && (os.Getenv("FORCE_COLOR") != "" || stdoutTerminal)

// Live reports whether output can be redrawn in place, as `execution watch` does.
func Live() bool { return stdoutTerminal }

// ColorsEnabled reports whether output is coloured: only on a terminal, unless
// NO_COLOR or FORCE_COLOR say otherwise.
func ColorsEnabled() bool { return colorsEnabled }

func style(code, s string) string {
	if !colorsEnabled {
		return s
	}
	return "\x1b[" + code + "m" + s + "\x1b[0m"
}

func Bold(s string) string  { return style("1", s) }
func Red(s string) string   { return style("31", s) }
func Green(s string) string { return style("32", s) }

type Datatype string

const (
	JSON Datatype = "json"
	YAML Datatype = "yaml"
)

// ResolveDatatype: an explicit type first, then the file's extension, then YAML.
func ResolveDatatype(explicit, file string) (Datatype, error) {
	switch explicit {
	case "json", "yaml":
		return Datatype(explicit), nil
	case "":
		if strings.HasSuffix(strings.ToLower(file), ".json") {
			return JSON, nil
		}
		return YAML, nil
	default:
		return "", fmt.Errorf("unsupported output format '%s'. Use \"json\" or \"yaml\"", explicit)
	}
}

// PathSegment reduces an id from the platform to one path segment, so that one
// containing "../" cannot write outside the chosen directory; Base alone leaves "..".
func PathSegment(id string) string {
	segment := filepath.Base(id)
	if segment == "" || segment == "." || segment == ".." || segment == string(filepath.Separator) {
		return "_"
	}
	return segment
}
