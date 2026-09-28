// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package resource_test

import (
	"encoding/json"
	"testing"

	"github.com/steadybit/cli/v6/internal/output"
	"github.com/steadybit/cli/v6/internal/platformtest"
	"github.com/steadybit/cli/v6/internal/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var items = []json.RawMessage{json.RawMessage(`{"id":"b","name":"Beta","count":2}`), json.RawMessage(`{"id":"a","name":"Alpha","count":1}`)}

func TestListPrintsTheTableForPeople(t *testing.T) {
	out, err := platformtest.Stdout(t, func() error {
		return resource.List(items, "", func() error { print("TABLE"); return nil })
	})
	require.NoError(t, err)
	assert.Empty(t, out) // the table callback ran instead
}

func TestListPrintsItemsInThePlatformsOrder(t *testing.T) {
	out, err := platformtest.Stdout(t, func() error { return resource.List(items, "yaml", nil) })

	require.NoError(t, err)
	assert.Equal(t, "- id: b\n  name: Beta\n  count: 2\n- id: a\n  name: Alpha\n  count: 1\n", out)
}

func TestJQPrintsStringsRawAndValuesAsJSON(t *testing.T) {
	output.JQ = `.[] | select(.count > 1) | .name, {id}`
	t.Cleanup(func() { output.JQ = "" })

	out, err := platformtest.Stdout(t, func() error { return resource.List(items, "", nil) })

	require.NoError(t, err)
	assert.Equal(t, "Beta\n{\n  \"id\": \"b\"\n}\n", out)
}

func TestJQReportsABadExpression(t *testing.T) {
	output.JQ = `.[`
	t.Cleanup(func() { output.JQ = "" })

	_, err := platformtest.Stdout(t, func() error { return resource.List(items, "", nil) })

	assert.ErrorContains(t, err, "Invalid --jq expression")
}
