// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package jsyaml

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func isPrintable(c rune) bool {
	return c == 0x09 || c == 0x0A || c == 0x0D || (c >= 0x20 && c <= 0x7E) || c == 0x85 ||
		(c >= 0xA0 && c <= 0xD7FF) || (c >= 0xE000 && c <= 0xFFFD) || (c >= 0x10000 && c <= 0x10FFFF)
}

func isNbChar(c rune) bool { return isPrintable(c) && c != '\n' && c != '\r' && c != 0xFEFF }
func isNsChar(c rune) bool { return isNbChar(c) && c != ' ' && c != '\t' }
func isWhite(c rune) bool  { return c == ' ' || c == '\t' }

func isIndicator(c rune) bool { return strings.ContainsRune("-?:,[]{}#&*!|>'\"%@`", c) }

// isPlainSafe reports whether s matches YAML's ns-plain-multi-line (or, for keys, the
// one-line form) in block context: the regular expressions js-yaml checks, written as
// a scanner because Go's regexp has no lookahead.
func isPlainSafe(s string, oneLine bool) bool {
	r := []rune(s)
	i := 0
	at := func(j int) (rune, bool) {
		if j < len(r) {
			return r[j], true
		}
		return 0, false
	}
	safeAt := func(j int) bool {
		c, ok := at(j)
		return ok && isNsChar(c)
	}
	// plainChar consumes one ns-plain-char and the '#' characters that may follow it.
	plainChar := func() bool {
		c, ok := at(i)
		if !ok {
			return false
		}
		switch {
		case c == ':':
			if !safeAt(i + 1) {
				return false
			}
		case c == '#' || !isNsChar(c):
			return false
		}
		i++
		for c, ok := at(i); ok && c == '#'; c, ok = at(i) {
			i++
		}
		return true
	}

	c, ok := at(0)
	if !ok {
		return false
	}
	switch {
	case isNsChar(c) && !isIndicator(c):
		i++
	case (c == '?' || c == ':' || c == '-') && safeAt(1):
		i++
	default:
		return false
	}
	for c, ok := at(i); ok && c == '#'; c, ok = at(i) {
		i++
	}
	for i < len(r) {
		j := i
		for j < len(r) && isWhite(r[j]) {
			j++
		}
		if j == len(r) {
			return false // trailing whitespace
		}
		if r[j] == '\n' {
			if j != i || oneLine {
				return false
			}
			for i < len(r) && r[i] == '\n' {
				i++
			}
			if !plainChar() {
				return false
			}
			continue
		}
		i = j
		if !plainChar() {
			return false
		}
	}
	return true
}

var forbiddenFirstLine = regexp.MustCompile(`^(?:---|\.\.\.)(?:$|[ \t\n\r])`)

func canUsePlain(l layout) bool {
	if l.value != "" {
		if !isPlainSafe(l.value, l.isKey) {
			return false
		}
		if l.shiftOfFirst == 0 && forbiddenFirstLine.MatchString(l.value) {
			return false
		}
	}
	resolved := resolveImplicit(l.value)
	if resolved != l.tag {
		return false
	}
	return !(l.value == "=" && resolved == "str")
}

func canUseSingleQuoted(l layout) bool {
	for _, c := range l.value {
		ok := c == 0x09 || (c >= 0x20 && c <= 0xD7FF) || (c >= 0xE000 && c <= 0xFFFF) || c >= 0x10000
		if !ok && !(c == '\n' && !l.isKey) {
			return false
		}
	}
	return !strings.Contains(l.value, " \n") && !strings.Contains(l.value, "\t\n") &&
		!strings.Contains(l.value, "\n ") && !strings.Contains(l.value, "\n\t")
}

var startsWithSpaceAfterBreaks = regexp.MustCompile(`^\n* `)

func canUseBlock(l layout) bool {
	if l.flowOnly {
		return false
	}
	for _, c := range l.value {
		if !isNbChar(c) && c != '\n' {
			return false
		}
	}
	contentIndent := l.shiftOfContent - l.shiftOfParent
	if contentIndent < 1 {
		return false
	}
	return !(contentIndent > 9 && startsWithSpaceAfterBreaks.MatchString(l.value))
}

var (
	coreInt   = regexp.MustCompile(`^(?:0o[0-7]+|0x[0-9a-fA-F]+|[-+]?[0-9]+)$`)
	coreFloat = regexp.MustCompile(`^(?:[-+]?[0-9]+(?:\.[0-9]*)?(?:[eE][-+]?[0-9]+)?|[-+]?\.[0-9]+(?:[eE][-+]?[0-9]+)?|[-+]?\.(?:inf|Inf|INF)|\.(?:nan|NaN|NAN))$`)
	special   = regexp.MustCompile(`^(?:[-+]?\.(?:inf|Inf|INF)|\.(?:nan|NaN|NAN))$`)
	date      = regexp.MustCompile(`^([0-9]{4})-([0-9]{2})-([0-9]{2})$`)
	timestamp = regexp.MustCompile(`^([0-9]{4})-([0-9][0-9]?)-([0-9][0-9]?)(?:[Tt]|[ \t]+)([0-9][0-9]?):([0-9]{2}):([0-9]{2})(?:\.([0-9]*))?(?:[ \t]*(Z|([-+])([0-9][0-9]?)(?::([0-9]{2}))?))?$`)
)

// resolveImplicit is the tag a plain scalar would read back as under the CLI's schema:
// the YAML core tags, merge keys and timestamps.
func resolveImplicit(s string) string {
	switch s {
	case "", "~", "null", "Null", "NULL":
		return "null"
	case "true", "True", "TRUE", "false", "False", "FALSE":
		return "bool"
	case "<<":
		return "merge"
	}
	if coreInt.MatchString(s) {
		return "int"
	}
	if coreFloat.MatchString(s) {
		// JavaScript's parseFloat overflows to Infinity, which only resolves when it was
		// written as .inf; otherwise the text stays a string.
		if special.MatchString(s) {
			return "float"
		}
		if f, _ := strconv.ParseFloat(s, 64); !math.IsInf(f, 0) {
			return "float"
		}
	}
	if isTimestamp(s) {
		return "timestamp"
	}
	return "str"
}

func isTimestamp(s string) bool {
	if m := date.FindStringSubmatch(s); m != nil {
		return validDate(m[1], m[2], m[3])
	}
	m := timestamp.FindStringSubmatch(s)
	if m == nil || !validDate(m[1], m[2], m[3]) {
		return false
	}
	h, _ := strconv.Atoi(m[4])
	mi, _ := strconv.Atoi(m[5])
	sec, _ := strconv.Atoi(m[6])
	if h > 23 || mi > 59 || sec > 59 {
		return false
	}
	if m[9] != "" {
		oh, _ := strconv.Atoi(m[10])
		om, _ := strconv.Atoi(m[11])
		if oh > 23 || om > 59 {
			return false
		}
	}
	return true
}

func validDate(y, m, d string) bool {
	year, _ := strconv.Atoi(y)
	month, _ := strconv.Atoi(m)
	day, _ := strconv.Atoi(d)
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return t.Year() == year && int(t.Month()) == month && t.Day() == day
}
