// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package jsyaml writes YAML and JSON exactly as the TypeScript CLI did with js-yaml's
// dump and JSON.stringify. Users keep the files `get` writes in Git; any difference in
// formatting would show up as a change to every one of them after upgrading.
package jsyaml

import (
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Map is a JavaScript object: string keys in JavaScript's iteration order.
type Map struct {
	keys   []string
	values map[string]any
}

func NewMap() *Map { return &Map{values: map[string]any{}} }

// Set keeps a key's first position when it is set again, as assigning to a JavaScript
// object property does.
func (m *Map) Set(key string, value any) {
	if _, exists := m.values[key]; !exists {
		m.keys = append(m.keys, key)
	}
	m.values[key] = value
}

func (m *Map) Get(key string) (any, bool) {
	v, ok := m.values[key]
	return v, ok
}

func (m *Map) Delete(key string) {
	if _, ok := m.values[key]; !ok {
		return
	}
	delete(m.values, key)
	for i, k := range m.keys {
		if k == key {
			m.keys = append(m.keys[:i], m.keys[i+1:]...)
			return
		}
	}
}

// SetFirst puts a key at the start, as the spread `{ key, ...experiment }` did.
func (m *Map) SetFirst(key string, value any) {
	m.Delete(key)
	m.keys = append([]string{key}, m.keys...)
	m.values[key] = value
}

func (m *Map) Len() int { return len(m.keys) }

// Keys in JavaScript's order: keys that are array indices first, ascending, then the
// rest in insertion order. JSON.parse and js-yaml both see objects this way.
func (m *Map) Keys() []string {
	var indices, others []string
	for _, k := range m.keys {
		if isArrayIndex(k) {
			indices = append(indices, k)
		} else {
			others = append(others, k)
		}
	}
	sort.SliceStable(indices, func(i, j int) bool {
		a, _ := strconv.ParseUint(indices[i], 10, 64)
		b, _ := strconv.ParseUint(indices[j], 10, 64)
		return a < b
	})
	return append(indices, others...)
}

func isArrayIndex(k string) bool {
	if k == "" || len(k) > 10 || (len(k) > 1 && k[0] == '0') {
		return false
	}
	n, err := strconv.ParseUint(k, 10, 64)
	return err == nil && n < math.MaxUint32
}

// Timestamp is a YAML timestamp read from a file, a JavaScript Date in the TypeScript CLI.
type Timestamp time.Time

func (t Timestamp) ISO() string {
	return time.Time(t).UTC().Format("2006-01-02T15:04:05.000Z")
}

// NumberString formats a float64 as JavaScript's Number.prototype.toString does.
func NumberString(f float64) string {
	switch {
	case math.IsNaN(f):
		return "NaN"
	case math.IsInf(f, 1):
		return "Infinity"
	case math.IsInf(f, -1):
		return "-Infinity"
	case f == 0:
		return "0"
	}
	sign := ""
	if f < 0 {
		sign, f = "-", -f
	}
	// Shortest round-trip digits, the same ones JavaScript picks.
	e := strconv.FormatFloat(f, 'e', -1, 64)
	mantissa, exponent, _ := strings.Cut(e, "e")
	digits := strings.Replace(mantissa, ".", "", 1)
	exp, _ := strconv.Atoi(exponent)
	k, n := len(digits), exp+1
	switch {
	case k <= n && n <= 21:
		return sign + digits + strings.Repeat("0", n-k)
	case 0 < n && n <= 21:
		return sign + digits[:n] + "." + digits[n:]
	case -6 < n && n <= 0:
		return sign + "0." + strings.Repeat("0", -n) + digits
	}
	expSign := "+"
	if n-1 < 0 {
		expSign = "-"
	}
	abs := n - 1
	if abs < 0 {
		abs = -abs
	}
	if k == 1 {
		return sign + digits + "e" + expSign + strconv.Itoa(abs)
	}
	return sign + digits[:1] + "." + digits[1:] + "e" + expSign + strconv.Itoa(abs)
}

// Clone copies a value deeply, so that one copy can be trimmed for a request while the
// other is written back to its file intact.
func Clone(value any) any {
	switch v := value.(type) {
	case *Map:
		c := NewMap()
		for _, k := range v.keys {
			c.keys = append(c.keys, k)
			c.values[k] = Clone(v.values[k])
		}
		return c
	case []any:
		c := make([]any, len(v))
		for i, item := range v {
			c[i] = Clone(item)
		}
		return c
	default:
		return v
	}
}
