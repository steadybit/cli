// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package gitops compares files kept in Git with what the platform holds, for `diff`
// and `apply --dry-run`, and moves whole projects between the two.
package gitops

import (
	"math"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/steadybit/cli/internal/jsyaml"
)

// Comparable projects the platform's version of a document onto the file's: every field
// the file sets, and only those fields the file leaves out that hold something. The
// platform fills in defaults (false, empty lists and maps) that a file written by hand
// omits; reporting those would bury every real difference.
func Comparable(local, remote any) any {
	lm, lok := local.(*jsyaml.Map)
	rm, rok := remote.(*jsyaml.Map)
	if lok && rok {
		out := jsyaml.NewMap()
		for _, k := range lm.Keys() {
			lv, _ := lm.Get(k)
			if rv, ok := rm.Get(k); ok {
				out.Set(k, Comparable(lv, rv))
			} else {
				out.Set(k, missing{})
			}
		}
		for _, k := range rm.Keys() {
			if _, ok := lm.Get(k); ok {
				continue
			}
			if rv, _ := rm.Get(k); !isEmpty(rv) {
				out.Set(k, rv)
			}
		}
		return out
	}
	ll, lok := local.([]any)
	rl, rok := remote.([]any)
	if lok && rok {
		out := make([]any, len(rl))
		for i := range rl {
			if i < len(ll) {
				out[i] = Comparable(ll[i], rl[i])
			} else {
				out[i] = rl[i]
			}
		}
		return out
	}
	return remote
}

// missing marks a field the file sets and the platform does not have.
type missing struct{}

// isEmpty is what the platform writes when nothing was set: nothing, false, zero, an
// empty string, list or map.
func isEmpty(v any) bool {
	switch x := v.(type) {
	case nil:
		return true
	case bool:
		return !x
	case float64:
		return x == 0 && !math.Signbit(x)
	case string:
		return x == ""
	case []any:
		return len(x) == 0
	case *jsyaml.Map:
		for _, k := range x.Keys() {
			if value, _ := x.Get(k); !isEmpty(value) {
				return false
			}
		}
		return true
	}
	return false
}

// withoutMissing drops the markers, so a field the platform lacks shows as removed.
func withoutMissing(v any) any {
	switch x := v.(type) {
	case *jsyaml.Map:
		out := jsyaml.NewMap()
		for _, k := range x.Keys() {
			value, _ := x.Get(k)
			if _, gone := value.(missing); !gone {
				out.Set(k, withoutMissing(value))
			}
		}
		return out
	case []any:
		out := make([]any, len(x))
		for i, item := range x {
			out[i] = withoutMissing(item)
		}
		return out
	}
	return v
}

// Diff is a unified diff from the platform's version to the file's, empty when they
// agree. Both sides are rendered as YAML the way `get` writes it.
func Diff(file string, local, remote *jsyaml.Map) (string, error) {
	projected := withoutMissing(Comparable(local, remote)).(*jsyaml.Map)
	before, after := jsyaml.Dump(projected), jsyaml.Dump(local)
	if before == after {
		return "", nil
	}
	return difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        lines(before),
		B:        lines(after),
		FromFile: "platform",
		ToFile:   file,
		Context:  3,
	})
}

// lines splits a YAML rendering after each newline. difflib's own splitting adds an empty
// last line, which shows up as a stray line of context.
func lines(text string) []string {
	split := strings.SplitAfter(text, "\n")
	if len(split) > 0 && split[len(split)-1] == "" {
		split = split[:len(split)-1]
	}
	return split
}

// Changes counts the changed lines of a diff, for dry runs.
func Changes(diff string) int {
	n := 0
	for _, line := range strings.Split(diff, "\n") {
		if (strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-")) && !strings.HasPrefix(line, "+++") && !strings.HasPrefix(line, "---") {
			n++
		}
	}
	return n
}
