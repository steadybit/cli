// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package scope_test

import (
	"testing"

	"github.com/steadybit/cli/v6/internal/jsyaml"
	"github.com/steadybit/cli/v6/internal/output"
	"github.com/steadybit/cli/v6/internal/scope"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func doc(t *testing.T, text string) *jsyaml.Map {
	t.Helper()
	d, err := output.ParseDocument([]byte(text))
	require.NoError(t, err)
	return d.Value()
}

func TestForFileKeepsTheQueryWhenThePlatformReportsOne(t *testing.T) {
	env := doc(t, `{"name":"Prod","predicate":{"query":"a=\"b\""},"query":"a=\"b\""}`)

	scope.ForFile(env, scope.Environment)

	assert.Equal(t, `{"name":"Prod","query":"a=\"b\""}`, jsyaml.CompactJSON(env))
}

func TestForFileKeepsThePredicateOfAScopeWithoutQuery(t *testing.T) {
	env := doc(t, `{"name":"Prod","predicate":{"operator":"AND","predicates":[]},"query":null}`)

	scope.ForFile(env, scope.Environment)

	assert.Equal(t, `{"name":"Prod","predicate":{"operator":"AND","predicates":[]}}`, jsyaml.CompactJSON(env))
}

// An empty query selects every target, as Global's does, and is a query all the same.
func TestForFileKeepsAnEmptyQuery(t *testing.T) {
	env := doc(t, `{"name":"Global","predicate":{"query":""},"query":""}`)

	scope.ForFile(env, scope.Environment)

	assert.Equal(t, `{"name":"Global","query":""}`, jsyaml.CompactJSON(env))
}

func TestForFileTreatsEveryStepRadius(t *testing.T) {
	experiment := doc(t, `{"lanes":[{"steps":[{"type":"wait"},{"type":"action","radius":{"predicate":{"query":"x"},"query":"x","percentage":50}}]},`+
		`{"steps":[{"type":"action","radius":{"predicate":{"operator":"AND","predicates":[]}}}]}]}`)

	scope.ForFile(experiment, scope.Steps)

	assert.Equal(t, `{"lanes":[{"steps":[{"type":"wait"},{"type":"action","radius":{"query":"x","percentage":50}}]},`+
		`{"steps":[{"type":"action","radius":{"predicate":{"operator":"AND","predicates":[]}}}]}]}`, jsyaml.CompactJSON(experiment))
}

// The platform takes a query over a predicate, so a file holding a query is compared by
// it alone, even when it still holds a stale predicate beside it.
func TestComparableComparesAFileWithAQueryByItsQuery(t *testing.T) {
	local := doc(t, `{"query":"x","predicate":{"operator":"AND","predicates":[]}}`)
	remote := doc(t, `{"predicate":{"query":"x"},"query":"x"}`)

	hints := scope.Comparable(local, remote, scope.Environment)

	assert.Empty(t, hints)
	assert.Equal(t, `{"query":"x"}`, jsyaml.CompactJSON(local))
	assert.Equal(t, `{"query":"x"}`, jsyaml.CompactJSON(remote))
}

// The query the platform derives from a predicate is no difference to a file without one.
func TestComparableComparesAFileWithoutQueryByItsPredicate(t *testing.T) {
	local := doc(t, `{"predicate":{"query":"x"}}`)
	remote := doc(t, `{"predicate":{"query":"x"},"query":"x"}`)

	hints := scope.Comparable(local, remote, scope.Environment)

	assert.Empty(t, hints)
	assert.Equal(t, `{"predicate":{"query":"x"}}`, jsyaml.CompactJSON(remote))
}

// Once the platform holds a predicate tree as a query, a file keeping the tree learns
// which query to use instead.
func TestComparableHintsAtTheQueryForAPredicateTree(t *testing.T) {
	local := doc(t, `{"lanes":[{"steps":[{"radius":{"predicate":{"operator":"AND","predicates":[{"key":"k","operator":"EQUALS","values":["v"]}]}}}]}]}`)
	remote := doc(t, `{"lanes":[{"steps":[{"radius":{"predicate":{"query":"k=\"v\""},"query":"k=\"v\""}}]}]}`)

	hints := scope.Comparable(local, remote, scope.Steps)

	assert.Equal(t, []string{`lanes[0].steps[0].radius.predicate is deprecated; the platform holds this scope as query: "k=\"v\"". ` +
		`Unless the predicate was edited, replace it with that query.`}, hints)
	assert.Equal(t, `{"lanes":[{"steps":[{"radius":{"predicate":{"query":"k=\"v\""}}}]}]}`, jsyaml.CompactJSON(remote))
}

// A tree the platform keeps as a tree, having no query that selects the same targets,
// is compared as it is, without a hint.
func TestComparableLeavesATreeThePlatformKeeps(t *testing.T) {
	local := doc(t, `{"predicate":{"operator":"AND","predicates":[]}}`)
	remote := doc(t, `{"predicate":{"operator":"AND","predicates":[]},"query":null}`)

	assert.Empty(t, scope.Comparable(local, remote, scope.Environment))
	assert.Equal(t, `{"predicate":{"operator":"AND","predicates":[]}}`, jsyaml.CompactJSON(remote))
}
