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

func TestWithFieldFirstKeepsTheFileWhereALineInFrontDoes(t *testing.T) {
	for name, tc := range map[string]struct{ in, out string }{
		"plain":           {"# kept\nname: new\n", "key: NEW-1\n# kept\nname: new\n"},
		"document marker": {"# head\n---\nname: new\n", "# head\n---\nkey: NEW-1\nname: new\n"},
		"byte order mark": {"\ufeffname: new\n", "\ufeffkey: NEW-1\nname: new\n"},
		"crlf":            {"---\r\nname: new\r\n", "---\r\nkey: NEW-1\nname: new\r\n"},
	} {
		t.Run(name, func(t *testing.T) {
			doc, err := ParseDocument([]byte(tc.in))
			require.NoError(t, err)
			out, ok := WithFieldFirst([]byte(tc.in), doc, "key", "NEW-1")
			require.True(t, ok)
			assert.Equal(t, tc.out, string(out))
		})
	}
}

// Each of these would have come out as a different document, or not parse at all.
func TestWithFieldFirstRefusesWhatALineInFrontWouldBreak(t *testing.T) {
	for name, in := range map[string]string{
		"flow mapping":      "{name: new, team: ADM}\n",
		"empty field":       "key: ''\nname: new\n",
		"null field":        "name: new\nkey:\n",
		"marker with value": "--- {name: new}\n",
	} {
		t.Run(name, func(t *testing.T) {
			doc, err := ParseDocument([]byte(in))
			require.NoError(t, err)
			_, ok := WithFieldFirst([]byte(in), doc, "key", "NEW-1")
			assert.False(t, ok)
		})
	}
}

func TestReadsJSONAfterAByteOrderMark(t *testing.T) {
	content := []byte("\ufeff{\"name\":\"new\"}")
	assert.True(t, IsJSON(content))
	doc, err := ParseDocument(content)
	require.NoError(t, err)
	assert.Equal(t, []string{"name"}, doc.Value().Keys())
}

func TestNeverLetsAnIdStepOutOfADirectory(t *testing.T) {
	for _, id := range []string{"..", ".", "", "../..", "a/../..", "/"} {
		assert.Equal(t, "_", PathSegment(id), id)
	}
	assert.Equal(t, "evil", PathSegment("../../evil"))
	assert.Equal(t, "report.zip", PathSegment("report.zip"))
}
