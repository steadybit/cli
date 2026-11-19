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
	doc.Value().SetFirst("key", "ADM-1")
	out, _ := doc.Render(YAML)
	assert.Equal(t, "key: ADM-1\nname: x\n", string(out))
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

// yaml.Unmarshal would read the first of several documents and drop the others, so
// applying such a file would have saved one object and reported success.
func TestReadsEveryYAMLDocument(t *testing.T) {
	docs, err := ParseDocuments([]byte("---\nname: a\n---\n# b follows\nname: b\n---\n---\nname: c\n...\n"))
	require.NoError(t, err)
	var names []string
	for _, doc := range docs {
		name, _ := doc.Get("name")
		names = append(names, name)
	}
	assert.Equal(t, []string{"a", "b", "c"}, names)
}

func TestNamesTheDocumentThatIsNoObject(t *testing.T) {
	_, err := ParseDocuments([]byte("name: a\n---\n- b\n"))
	assert.EqualError(t, err, "document 2: expected an object")
}

func TestOneObjectIsExpectedWhereSeveralCannotBeTaken(t *testing.T) {
	in := []byte("name: a\n---\nname: b\n")
	_, err := ParseDocument(in)
	assert.EqualError(t, err, "found 2 documents where one object is expected")
	_, err = ParseValue(in)
	assert.EqualError(t, err, "found 2 documents where one object is expected")
}

func TestReadsOneYAMLDocumentBetweenMarkers(t *testing.T) {
	for name, in := range map[string]string{
		"leading marker":         "---\nname: a\n",
		"trailing marker":        "name: a\n---\n",
		"comment after a marker": "name: a\n---\n# nothing more\n",
		"document end marker":    "name: a\n...\n",
		"markers on either side": "---\nname: a\n---\n",
	} {
		t.Run(name, func(t *testing.T) {
			doc, err := ParseDocument([]byte(in))
			require.NoError(t, err)
			name, _ := doc.Get("name")
			assert.Equal(t, "a", name)
		})
	}
}
