// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package output

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func render(t *testing.T, input string, datatype Datatype) string {
	t.Helper()
	doc, err := ParseDocument([]byte(input))
	require.NoError(t, err)
	out, err := doc.Render(datatype)
	require.NoError(t, err)
	return string(out)
}

// JSON.parse keeps the sign of -0, and js-yaml wrote it as a float.
func TestKeepsNegativeZero(t *testing.T) {
	assert.Equal(t, "n: -0.0\n", render(t, `{"n":-0}`, YAML))
	assert.Equal(t, "{\n  \"n\": 0\n}\n", render(t, `{"n":-0}`, JSON))
}

// JSON.parse orders integer-like keys first; the fixture cannot show it, since the
// generator's own JSON.stringify already reordered them.
func TestOrdersKeysAsJavaScriptObjectsDo(t *testing.T) {
	assert.Equal(t, "'2': two\n'10': ten\nb: b\n'01': zero-one\n", render(t, `{"b":"b","10":"ten","01":"zero-one","2":"two"}`, YAML))
}

// A YAML parser rejects DEL, which JSON allows unescaped.
func TestReadsJSONThatYAMLWouldReject(t *testing.T) {
	assert.Equal(t, "s: \"a\\x7Fb\"\n", render(t, "{\"s\":\"a\x7fb\"}", YAML))
}

func TestResolvesMergeKeysInYAMLFiles(t *testing.T) {
	doc, err := ParseDocument([]byte("base: &b\n  x: 1\n  y: 2\nderived:\n  <<: *b\n  y: 3\n"))
	require.NoError(t, err)
	out, err := doc.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, `{"base":{"x":1,"y":2},"derived":{"y":3,"x":1}}`, string(out))
}

// js-yaml's load turned timestamps into Dates, which JSON.stringify writes as ISO strings.
func TestSendsYAMLTimestampsAsJavaScriptDates(t *testing.T) {
	doc, err := ParseDocument([]byte("startAt: 2030-06-01T09:00:00Z\n"))
	require.NoError(t, err)
	out, _ := doc.MarshalJSON()
	assert.Equal(t, `{"startAt":"2030-06-01T09:00:00.000Z"}`, string(out))
}

func TestSetFirstMovesTheKeyToTheTop(t *testing.T) {
	doc, err := ParseDocument([]byte(`{"name":"x","key":"old"}`))
	require.NoError(t, err)
	doc.SetFirst("key", "ADM-1")
	out, _ := doc.Render(YAML)
	assert.Equal(t, "key: ADM-1\nname: x\n", string(out))
}
