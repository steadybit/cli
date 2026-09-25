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

func TestKeepsThePlatformsFieldOrder(t *testing.T) {
	assert.Equal(t, "name: x\nteam: ADM\nactive: true\n", render(t, `{"name":"x","team":"ADM","active":true}`, YAML))
}

func TestQuotesLikeJsYaml(t *testing.T) {
	out := render(t, `{"a":"*","b":"2026-02-25T06:35:52Z","c":"true","d":"plain","e":"line1\nline2"}`, YAML)
	assert.Equal(t, "a: '*'\nb: '2026-02-25T06:35:52Z'\nc: 'true'\nd: plain\ne: |-\n  line1\n  line2\n", out)
}

func TestWritesWholeFloatsAsJavaScriptDid(t *testing.T) {
	assert.Equal(t, "a: 1\nb: 1.5\n", render(t, `{"a":1.0,"b":1.5}`, YAML))
	assert.Equal(t, "{\n  \"a\": 1,\n  \"b\": 1.5\n}\n", render(t, `{"a":1.0,"b":1.5}`, JSON))
}

func TestIndentsSequencesUnderTheirKey(t *testing.T) {
	assert.Equal(t, "tags:\n  - a\n  - b\n", render(t, `{"tags":["a","b"]}`, YAML))
}

func TestResolvesMergeKeysWhenSending(t *testing.T) {
	doc, err := ParseDocument([]byte("base: &b\n  x: 1\n  y: 2\nderived:\n  <<: *b\n  y: 3\n"))
	require.NoError(t, err)
	out, err := doc.MarshalJSON()
	require.NoError(t, err)
	assert.JSONEq(t, `{"base":{"x":1,"y":2},"derived":{"y":3,"x":1}}`, string(out))
}

func TestDoesNotEscapeHTMLInJSON(t *testing.T) {
	assert.Equal(t, "{\n  \"u\": \"a?b=1&c=<d>\"\n}\n", render(t, `{"u":"a?b=1&c=<d>"}`, JSON))
}

func TestSetFirstMovesTheKeyToTheTop(t *testing.T) {
	doc, err := ParseDocument([]byte(`{"name":"x","key":"old"}`))
	require.NoError(t, err)
	doc.SetFirst("key", "ADM-1")
	out, err := doc.Render(YAML)
	require.NoError(t, err)
	assert.Equal(t, "key: ADM-1\nname: x\n", string(out))
}
