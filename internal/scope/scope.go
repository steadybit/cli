// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package scope keeps the targets an environment or a step selects in one form per file.
//
// The platform stores one scope and reports it twice: as `predicate`, and, when it is a
// query, the same query again as `query`. Sent both, it takes `query` and ignores
// `predicate`. A file holding both would hide an edit to `predicate` behind the
// unchanged `query`, so files hold `query` when the platform reports one, and
// `predicate`, which is deprecated, only for a scope that has no query.
package scope

import (
	"fmt"

	"github.com/steadybit/cli/v6/internal/jsyaml"
)

// Holder is a map holding a scope, and where in its document it sits.
type Holder struct {
	// Path is how messages name it: "" for the document itself, or
	// "lanes[0].steps[1].radius".
	Path string
	Map  *jsyaml.Map
}

// Environment finds the scope of an environment: the document itself.
func Environment(doc *jsyaml.Map) []Holder {
	return []Holder{{Map: doc}}
}

// Steps finds the scopes of an experiment or a template: each step's radius.
func Steps(doc *jsyaml.Map) []Holder {
	var holders []Holder
	for i, lane := range list(doc, "lanes") {
		laneMap, _ := lane.(*jsyaml.Map)
		for j, step := range list(laneMap, "steps") {
			stepMap, _ := step.(*jsyaml.Map)
			if stepMap == nil {
				continue
			}
			if radius, ok := stepMap.Get("radius"); ok {
				if radiusMap, ok := radius.(*jsyaml.Map); ok {
					holders = append(holders, Holder{Path: fmt.Sprintf("lanes[%d].steps[%d].radius", i, j), Map: radiusMap})
				}
			}
		}
	}
	return holders
}

func list(m *jsyaml.Map, key string) []any {
	if m == nil {
		return nil
	}
	v, _ := m.Get(key)
	l, _ := v.([]any)
	return l
}

// HasQuery tells whether a scope holds a query, which the platform takes over any
// predicate beside it. An empty query is one too: it selects every target.
func HasQuery(m *jsyaml.Map) bool {
	v, ok := m.Get("query")
	return ok && v != nil
}

// ForFile leaves one form of each scope in a document the platform returned: `query`
// when it holds one, `predicate` otherwise.
func ForFile(doc *jsyaml.Map, holders func(*jsyaml.Map) []Holder) {
	for _, h := range holders(doc) {
		if HasQuery(h.Map) {
			h.Map.Delete("predicate")
		} else {
			h.Map.Delete("query")
		}
	}
}

// Comparable makes the scopes of a file and of the platform's version comparable by the
// platform's rule: a file holding `query` is compared by its query alone; one without is
// compared by its predicate, whatever query the platform derived from it. It returns a
// hint for each scope the file still selects with a predicate while the platform holds
// a query, naming the query to use instead.
func Comparable(local, remote *jsyaml.Map, holders func(*jsyaml.Map) []Holder) []string {
	remotes := map[string]*jsyaml.Map{}
	for _, h := range holders(remote) {
		remotes[h.Path] = h.Map
	}
	var hints []string
	for _, l := range holders(local) {
		r := remotes[l.Path]
		if HasQuery(l.Map) {
			l.Map.Delete("predicate")
			if r != nil {
				r.Delete("predicate")
			}
			continue
		}
		l.Map.Delete("query")
		if r == nil {
			continue
		}
		query, _ := r.Get("query")
		r.Delete("query")
		lp, has := l.Map.Get("predicate")
		rp, _ := r.Get("predicate")
		if q, ok := query.(string); ok && has && jsyaml.CompactJSON(lp) != jsyaml.CompactJSON(rp) {
			hints = append(hints, hint(l.Path, q))
		}
	}
	return hints
}

func hint(path, query string) string {
	where := "predicate"
	if path != "" {
		where = path + ".predicate"
	}
	return fmt.Sprintf("%s is deprecated; the platform holds this scope as query: %s. "+
		"Unless the predicate was edited, replace it with that query.", where, jsyaml.CompactJSON(query))
}
