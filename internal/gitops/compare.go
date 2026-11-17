// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package gitops compares files kept in Git with what the platform holds, for `diff`
// and `apply --dry-run`, and moves whole projects between the two.
package gitops

import (
	"math"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/steadybit/cli/v6/internal/jsyaml"
)

// Comparable projects the platform's version of a document onto the file's: every field
// the file sets, and only those fields the file leaves out that hold something. The
// platform fills in defaults (false, empty lists and maps) that a file written by hand
// omits; reporting those would bury every real difference.
func Comparable(local, remote any) any {
	return comparable(local, remote, nil)
}

// comparable is Comparable with fields reported even when they hold what looks like an
// empty default, because their default is something else.
func comparable(local, remote any, always map[string]bool) any {
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
			if rv, _ := rm.Get(k); always[k] || !isEmpty(rv) {
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
	// An unquoted YAML date is the same moment as the platform's ISO string, however
	// either writes its fractions of a second.
	if lt, ok := local.(jsyaml.Timestamp); ok {
		if rs, ok := remote.(string); ok {
			if rt, err := time.Parse(time.RFC3339Nano, rs); err == nil && rt.Equal(time.Time(lt)) {
				return local
			}
		}
	}
	return remote
}

// inFileOrder puts the platform's sets into the order of the file's, so the order alone
// is never a difference: the items the file lists, in its order, then those only the
// platform holds. paths are the sets, as api.Unordered names them; keys name the fields
// that identify an item of a set of objects, by its path. Items without keys match when
// they agree on every field the file sets.
func inFileOrder(local, remote *jsyaml.Map, paths []string, keys map[string][]string) *jsyaml.Map {
	sorted := append([]string{}, paths...)
	// A set within a set is matched once its own items are lined up.
	sort.SliceStable(sorted, func(i, j int) bool {
		return strings.Count(sorted[i], ".") < strings.Count(sorted[j], ".")
	})
	out := jsyaml.Clone(remote).(*jsyaml.Map)
	for _, p := range sorted {
		alignAt(local, out, strings.Split(p, "."), keys[p])
	}
	return out
}

// alignAt lines up the set at path in remote, going into every item of the lists and,
// for "*", every value of the maps on its way. Lists on the way pair up by position.
func alignAt(local, remote any, path []string, keys []string) {
	switch r := remote.(type) {
	case *jsyaml.Map:
		l, _ := local.(*jsyaml.Map)
		if l == nil {
			return
		}
		fields := []string{path[0]}
		if path[0] == "*" {
			fields = r.Keys()
		}
		for _, field := range fields {
			rv, ok := r.Get(field)
			if !ok {
				continue
			}
			lv, ok := l.Get(field)
			if !ok {
				continue
			}
			if len(path) == 1 {
				ll, lok := lv.([]any)
				rl, rok := rv.([]any)
				if lok && rok {
					r.Set(field, matchItems(ll, rl, keys))
				}
				continue
			}
			alignAt(lv, rv, path[1:], keys)
		}
	case []any:
		l, _ := local.([]any)
		for i := range r {
			if i < len(l) {
				alignAt(l[i], r[i], path, keys)
			}
		}
	}
}

// matchItems orders remote like local. An item matched by one of its keys takes the
// file's spelling of it and loses the other keys the file leaves out: the platform fills
// in a member's username for an email, and its email for a username.
func matchItems(local, remote []any, keys []string) []any {
	out := make([]any, 0, len(remote))
	used := make([]bool, len(remote))
	for _, l := range local {
		for i, r := range remote {
			if used[i] {
				continue
			}
			if key, ok := sameItem(l, r, keys); ok {
				used[i] = true
				out = append(out, withoutOtherKeys(l, r, keys, key))
				break
			}
		}
	}
	for i, r := range remote {
		if !used[i] {
			out = append(out, r)
		}
	}
	return out
}

// sameItem tells whether r is the platform's version of l, and by which key.
func sameItem(l, r any, keys []string) (string, bool) {
	lm, lok := l.(*jsyaml.Map)
	rm, rok := r.(*jsyaml.Map)
	if lok && rok && len(keys) > 0 {
		for _, key := range keys {
			lv, _ := lm.Get(key)
			rv, _ := rm.Get(key)
			ls, _ := lv.(string)
			rs, _ := rv.(string)
			if ls != "" && strings.EqualFold(ls, rs) {
				return key, true
			}
		}
		return "", false
	}
	return "", jsyaml.CompactJSON(withoutMissing(Comparable(l, r))) == jsyaml.CompactJSON(l)
}

func withoutOtherKeys(l, r any, keys []string, matched string) any {
	lm, lok := l.(*jsyaml.Map)
	rm, rok := r.(*jsyaml.Map)
	if !lok || !rok || matched == "" {
		return r
	}
	out := jsyaml.Clone(rm).(*jsyaml.Map)
	// Matched regardless of case, an email is the same however either spells it.
	spelled, _ := lm.Get(matched)
	out.Set(matched, spelled)
	for _, key := range keys {
		if _, set := lm.Get(key); key != matched && !set {
			out.Delete(key)
		}
	}
	return out
}

// sameSet tells whether two lists hold the same items, in any order.
func sameSet(a, b any) bool {
	al, aok := a.([]any)
	bl, bok := b.([]any)
	if !aok || !bok || len(al) != len(bl) {
		return false
	}
	items := func(l []any) []string {
		out := make([]string, len(l))
		for i, item := range l {
			out[i] = jsyaml.CompactJSON(item)
		}
		sort.Strings(out)
		return out
	}
	return slices.Equal(items(al), items(bl))
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
// Fields in always are reported whenever the platform holds them and the file does not.
func Diff(file string, local, remote *jsyaml.Map, always ...string) (string, error) {
	report := map[string]bool{}
	for _, k := range always {
		report[k] = true
	}
	projected := withoutMissing(comparable(local, remote, report)).(*jsyaml.Map)
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
