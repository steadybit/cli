// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package jsyaml

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf16"
)

// A port of js-yaml 5's dump with the options the TypeScript CLI used: indent 2, line
// width 80, single quotes preferred, and a schema of the YAML core tags plus merge keys
// and timestamps. Only what JSON-shaped values need is ported; there are no tags,
// anchors or flow collections beyond the empty [] and {}.

const (
	indentStep = 2
	lineWidth  = 80
)

type style int

const (
	plain style = iota
	singleQuoted
	doubleQuoted
	literal
	folded
)

type dumper struct {
	openEnded bool
}

// Dump renders a value as js-yaml's dump did, trailing newline included.
func Dump(value any) string {
	d := &dumper{}
	out := d.node(0, value, false, true, true) + "\n"
	if d.openEnded {
		out += "...\n"
	}
	return out
}

func nextLine(level int) string {
	return "\n" + strings.Repeat(" ", indentStep*level)
}

func (d *dumper) node(level int, value any, isKey, block, compact bool) string {
	switch v := value.(type) {
	case *Map:
		if block && v.Len() > 0 {
			return d.blockMapping(level, v, compact)
		}
		d.openEnded = false
		return "{}"
	case []any:
		if block && len(v) > 0 {
			return d.blockSequence(level, v, compact)
		}
		d.openEnded = false
		return "[]"
	case string:
		return d.scalar(level, v, "str", isKey, !block)
	case nil:
		return d.scalar(level, "null", "null", isKey, !block)
	case bool:
		return d.scalar(level, strconv.FormatBool(v), "bool", isKey, !block)
	case float64:
		text, tag := representNumber(v)
		return d.scalar(level, text, tag, isKey, !block)
	case Timestamp:
		return d.scalar(level, v.ISO(), "timestamp", isKey, !block)
	default:
		return d.scalar(level, fmt.Sprint(v), "str", isKey, !block)
	}
}

// representNumber picks js-yaml's int or float tag and text for a JavaScript number.
func representNumber(f float64) (string, string) {
	s := NumberString(f)
	if f == math.Trunc(f) && !math.IsInf(f, 0) && !(f == 0 && math.Signbit(f)) && !strings.Contains(s, "e") {
		return s, "int"
	}
	switch {
	case math.IsNaN(f):
		return ".nan", "float"
	case math.IsInf(f, 1):
		return ".inf", "float"
	case math.IsInf(f, -1):
		return "-.inf", "float"
	case f == 0 && math.Signbit(f):
		return "-0.0", "float"
	}
	if leadingExponent.MatchString(s) {
		s = strings.Replace(s, "e", ".e", 1)
	}
	return s, "float"
}

var leadingExponent = regexp.MustCompile(`^[-+]?[0-9]+e`)

func (d *dumper) blockSequence(level int, items []any, compact bool) string {
	var result strings.Builder
	for _, item := range items {
		text := d.node(level+1, item, false, true, true)
		if !compact || result.Len() > 0 {
			result.WriteString(nextLine(level))
		}
		if text == "" || text[0] == '\n' {
			result.WriteString("-")
		} else {
			result.WriteString("- ")
		}
		result.WriteString(text)
	}
	return result.String()
}

func (d *dumper) blockMapping(level int, m *Map, compact bool) string {
	var result strings.Builder
	for _, key := range m.Keys() {
		var pair strings.Builder
		if !compact || result.Len() > 0 {
			pair.WriteString(nextLine(level))
		}
		keyText := d.node(level+1, key, true, true, true)
		explicit := strings.Contains(key, "\n") || length(keyText) > 1024
		if explicit {
			if keyText != "" && keyText[0] == '\n' {
				pair.WriteString("?")
			} else {
				pair.WriteString("? ")
			}
		}
		pair.WriteString(keyText)
		if explicit {
			pair.WriteString(nextLine(level))
		}
		valueText := d.node(level+1, m.values[key], false, true, explicit)
		if valueText == "" || valueText[0] == '\n' {
			pair.WriteString(":")
		} else {
			pair.WriteString(": ")
		}
		pair.WriteString(valueText)
		result.WriteString(pair.String())
	}
	return result.String()
}

type layout struct {
	value                                       string
	tag                                         string
	isKey, flowOnly                             bool
	shiftOfParent, shiftOfContent, shiftOfFirst int
	allowPlain, allowSingle, allowBlock         bool
}

func (d *dumper) scalar(level int, value, tag string, isKey, flowOnly bool) string {
	l := layout{value: value, tag: tag, isKey: isKey, flowOnly: flowOnly}
	if level == 0 {
		l.shiftOfParent = -1
	} else {
		l.shiftOfParent = indentStep * (level - 1)
		l.shiftOfFirst = indentStep * level
	}
	l.shiftOfContent = indentStep * max(1, level)
	l.allowPlain = canUsePlain(l)
	l.allowSingle = canUseSingleQuoted(l)
	l.allowBlock = canUseBlock(l)

	s := chooseStyle(l)
	d.openEnded = (s == literal || s == folded) && (value == "\n" || strings.HasSuffix(value, "\n\n"))
	return render(l, s)
}

func (l layout) allowed(s style) bool {
	switch s {
	case plain:
		return l.allowPlain
	case singleQuoted:
		return l.allowSingle
	case literal, folded:
		return l.allowBlock
	}
	return true
}

var (
	invisibles     = regexp.MustCompile(`[\t\x{7F}-\x{A0}\x{2028}\x{2029}\x{FEFF}\x{FFFE}\x{FFFF}]`)
	onlyWhitespace = regexp.MustCompile(`^[\t\n\v\f\r \x{A0}\x{1680}\x{2000}-\x{200A}\x{2028}\x{2029}\x{202F}\x{205F}\x{3000}\x{FEFF}]+$`)
	foldableSpace  = regexp.MustCompile(` [^ \t]`)
)

// chooseStyle applies js-yaml's default scalar style rules, in their order.
func chooseStyle(l layout) style {
	s := plain
	if invisibles.MatchString(l.value) || onlyWhitespace.MatchString(l.value) {
		s = doubleQuoted
	}
	if s == plain && !l.isKey {
		multiline := strings.Contains(l.value, "\n")
		if !l.allowBlock {
			if multiline {
				s = doubleQuoted
			}
		} else {
			available := max(min(lineWidth, 40), lineWidth-l.shiftOfContent)
			fold := false
			for _, line := range strings.Split(l.value, "\n") {
				if length(line) > available && !strings.HasPrefix(line, " ") && foldableSpace.MatchString(line) {
					fold = true
				}
			}
			if fold {
				s = folded
			} else if multiline {
				s = literal
			}
		}
	}
	if s == plain && !l.allowPlain {
		if l.allowSingle {
			s = singleQuoted
		} else {
			s = doubleQuoted
		}
	}
	if !l.allowed(s) {
		s = doubleQuoted
	}
	return s
}

func render(l layout, s style) string {
	switch s {
	case plain:
		return encodeFlowBreaks(l.value, l.shiftOfContent)
	case singleQuoted:
		return "'" + strings.ReplaceAll(encodeFlowBreaks(l.value, l.shiftOfContent), "'", "''") + "'"
	case literal:
		return "|" + blockHeader(l.value, l.shiftOfParent, l.shiftOfContent) +
			dropEndingNewline(indentString(l.value, l.shiftOfContent))
	case folded:
		available := max(min(lineWidth, 40), lineWidth-l.shiftOfContent)
		return ">" + blockHeader(l.value, l.shiftOfParent, l.shiftOfContent) +
			dropEndingNewline(indentString(foldBlockScalar(l.value, available), l.shiftOfContent))
	default:
		return `"` + escapeDoubleQuoted(l.value) + `"`
	}
}

var lineBreakRun = regexp.MustCompile(`(\n+)([^\n]*)`)

func encodeFlowBreaks(s string, shift int) string {
	first := strings.Index(s, "\n")
	if first == -1 {
		return s
	}
	pad := strings.Repeat(" ", shift)
	var b strings.Builder
	b.WriteString(s[:first])
	for _, m := range lineBreakRun.FindAllStringSubmatch(s[first:], -1) {
		b.WriteString(strings.Repeat("\n", len(m[1])+1) + pad + m[2])
	}
	return b.String()
}

func indentString(s string, spaces int) string {
	indent := strings.Repeat(" ", spaces)
	var b strings.Builder
	for len(s) > 0 {
		var line string
		if i := strings.Index(s, "\n"); i == -1 {
			line, s = s, ""
		} else {
			line, s = s[:i+1], s[i+1:]
		}
		if line != "" && line != "\n" {
			b.WriteString(indent)
		}
		b.WriteString(line)
	}
	return b.String()
}

var leadingSpaceAfterBreaks = regexp.MustCompile(`^\n* `)

func blockHeader(s string, shiftOfParent, shiftOfContent int) string {
	indicator := ""
	if leadingSpaceAfterBreaks.MatchString(s) {
		indicator = strconv.Itoa(shiftOfContent - shiftOfParent)
	}
	clip := strings.HasSuffix(s, "\n")
	chomp := "-"
	if clip {
		chomp = ""
		if strings.HasSuffix(s, "\n\n") || s == "\n" {
			chomp = "+"
		}
	}
	return indicator + chomp + "\n"
}

func dropEndingNewline(s string) string {
	return strings.TrimSuffix(s, "\n")
}

func isMoreIndented(c uint16) bool { return c == ' ' || c == '\t' }

// foldLine works in UTF-16 code units, which is what JavaScript's string indices and
// lengths count, so that lines with astral characters break where js-yaml broke them.
func foldLine(line []uint16, width int) []uint16 {
	if len(line) == 0 || isMoreIndented(line[0]) {
		return line
	}
	var result []uint16
	start, curr, next := 0, 0, 0
	for i := 0; i+1 < len(line); i++ {
		if line[i] != ' ' || line[i+1] == ' ' || line[i+1] == '\t' {
			continue
		}
		next = i
		if next-start > width {
			end := next
			if curr > start {
				end = curr
			}
			result = append(append(result, '\n'), line[start:end]...)
			start = end + 1
		}
		curr = next
	}
	result = append(result, '\n')
	if len(line)-start > width && curr > start {
		result = append(append(append(result, line[start:curr]...), '\n'), line[curr+1:]...)
	} else {
		result = append(result, line[start:]...)
	}
	return result[1:]
}

func foldBlockScalar(s string, width int) string {
	units := utf16.Encode([]rune(s))
	firstBreak := indexOf(units, '\n', 0)
	if firstBreak == -1 {
		firstBreak = len(units)
	}
	result := foldLine(units[:firstBreak], width)
	prevMoreIndented := len(units) > 0 && (units[0] == '\n' || isMoreIndented(units[0]))
	for pos := firstBreak; pos < len(units); {
		breaksEnd := pos
		for breaksEnd < len(units) && units[breaksEnd] == '\n' {
			breaksEnd++
		}
		lineEnd := indexOf(units, '\n', breaksEnd)
		if lineEnd == -1 {
			lineEnd = len(units)
		}
		prefix, line := units[pos:breaksEnd], units[breaksEnd:lineEnd]
		moreIndented := len(line) > 0 && isMoreIndented(line[0])
		result = append(result, prefix...)
		if !prevMoreIndented && !moreIndented && len(line) > 0 {
			result = append(result, '\n')
		}
		result = append(result, foldLine(line, width)...)
		prevMoreIndented = moreIndented
		pos = lineEnd
	}
	return string(utf16.Decode(result))
}

func indexOf(units []uint16, c uint16, from int) int {
	for i := from; i < len(units); i++ {
		if units[i] == c {
			return i
		}
	}
	return -1
}

func escapeDoubleQuoted(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case 0:
			b.WriteString(`\0`)
		case 7:
			b.WriteString(`\a`)
		case '\b':
			b.WriteString(`\b`)
		case '\t':
			b.WriteString(`\t`)
		case '\n':
			b.WriteString(`\n`)
		case '\v':
			b.WriteString(`\v`)
		case '\f':
			b.WriteString(`\f`)
		case '\r':
			b.WriteString(`\r`)
		case 0x1B:
			b.WriteString(`\e`)
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		case 0x85:
			b.WriteString(`\N`)
		case 0xA0:
			b.WriteString(`\_`)
		case 0x2028:
			b.WriteString(`\L`)
		case 0x2029:
			b.WriteString(`\P`)
		default:
			switch {
			case r < 0x20 || (r >= 0x7F && r <= 0xA0):
				fmt.Fprintf(&b, `\x%02X`, r)
			case r == 0xFEFF || r == 0xFFFE || r == 0xFFFF:
				fmt.Fprintf(&b, `\u%04X`, r)
			default:
				b.WriteRune(r)
			}
		}
	}
	return b.String()
}
