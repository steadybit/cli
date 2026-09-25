// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package jsyaml

import (
	"fmt"
	"math"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

// JSON renders like JSON.stringify(value, undefined, 2).
func JSON(value any) string {
	var b strings.Builder
	writeJSON(&b, value, "  ", "")
	return b.String()
}

// CompactJSON renders like JSON.stringify(value).
func CompactJSON(value any) string {
	var b strings.Builder
	writeJSON(&b, value, "", "")
	return b.String()
}

func writeJSON(b *strings.Builder, value any, step, indent string) {
	newline, space := "", ""
	if step != "" {
		newline, space = "\n", " "
	}
	switch v := value.(type) {
	case nil:
		b.WriteString("null")
	case bool:
		if v {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			b.WriteString("null")
		} else {
			b.WriteString(NumberString(v))
		}
	case string:
		b.WriteString(quoteJSON(v))
	case Timestamp:
		b.WriteString(quoteJSON(v.ISO()))
	case []any:
		if len(v) == 0 {
			b.WriteString("[]")
			return
		}
		b.WriteString("[" + newline)
		for i, item := range v {
			b.WriteString(indent + step)
			writeJSON(b, item, step, indent+step)
			if i < len(v)-1 {
				b.WriteString(",")
			}
			b.WriteString(newline)
		}
		b.WriteString(indent + "]")
	case *Map:
		if v.Len() == 0 {
			b.WriteString("{}")
			return
		}
		b.WriteString("{" + newline)
		keys := v.Keys()
		for i, k := range keys {
			b.WriteString(indent + step + quoteJSON(k) + ":" + space)
			writeJSON(b, v.values[k], step, indent+step)
			if i < len(keys)-1 {
				b.WriteString(",")
			}
			b.WriteString(newline)
		}
		b.WriteString(indent + "}")
	default:
		b.WriteString(quoteJSON(fmt.Sprint(v)))
	}
}

// quoteJSON escapes as JSON.stringify does: only quotes, backslashes and control
// characters. Unlike Go's encoder it leaves <, >, & and U+2028/U+2029 as they are.
func quoteJSON(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		s = s[size:]
		switch {
		case r == '"':
			b.WriteString(`\"`)
		case r == '\\':
			b.WriteString(`\\`)
		case r == '\b':
			b.WriteString(`\b`)
		case r == '\f':
			b.WriteString(`\f`)
		case r == '\n':
			b.WriteString(`\n`)
		case r == '\r':
			b.WriteString(`\r`)
		case r == '\t':
			b.WriteString(`\t`)
		case r < 0x20:
			fmt.Fprintf(&b, `\u%04x`, r)
		case r == utf8.RuneError && size == 1:
			b.WriteString(`�`)
		default:
			b.WriteRune(r)
		}
	}
	b.WriteByte('"')
	return b.String()
}

// length is a string's length in JavaScript, in UTF-16 code units.
func length(s string) int {
	n := 0
	for _, r := range s {
		n += utf16.RuneLen(r)
	}
	return n
}
