// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package output

import (
	"strings"
	"testing"

	"github.com/steadybit/cli/v6/internal/jsyaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// set sets the key and checks that the file reads back as the object Set holds.
func set(t *testing.T, in, key, value string) string {
	t.Helper()
	f, err := ParseFile([]byte(in))
	require.NoError(t, err)
	require.NoError(t, f.Set(key, value))
	read, err := ParseDocument(f.Bytes())
	require.NoError(t, err, string(f.Bytes()))
	assert.Equal(t, jsyaml.CompactJSON(f.Document().Value()), jsyaml.CompactJSON(read.Value()), "what the file reads back as")
	return string(f.Bytes())
}

// A new key of a block mapping is a line of its own at the start of the document; the
// rest of the file stays byte for byte.
func TestSetPutsANewKeyOfABlockMappingOnALineOfItsOwn(t *testing.T) {
	for name, tc := range map[string]struct{ in, out string }{
		"plain":               {"# kept\nname: new\n", "key: NEW-1\n# kept\nname: new\n"},
		"blank lines":         {"name: new\n\nlanes:\n- a\nnested:\n    deep: 1\n", "key: NEW-1\nname: new\n\nlanes:\n- a\nnested:\n    deep: 1\n"},
		"document marker":     {"# head\n---\nname: new\n", "# head\n---\nkey: NEW-1\nname: new\n"},
		"marker with comment": {"--- # c\nname: new\n", "--- # c\nkey: NEW-1\nname: new\n"},
		"comment after it":    {"---\n# d\nname: new\n", "---\nkey: NEW-1\n# d\nname: new\n"},
		"directive":           {"%YAML 1.1\n---\nname: new\n", "%YAML 1.1\n---\nkey: NEW-1\nname: new\n"},
		"indented":            {"  name: new\n  q: 1\n", "  key: NEW-1\n  name: new\n  q: 1\n"},
		"complex first key":   {"? name\n: new\nother: x\n", "key: NEW-1\n? name\n: new\nother: x\n"},
		"byte order mark":     {"\ufeffname: new\n", "\ufeffkey: NEW-1\nname: new\n"},
		"crlf":                {"---\r\nname: new\r\n", "---\r\nkey: NEW-1\nname: new\r\n"},
		"cr":                  {"# c\r---\rname: new\r", "# c\r---\rkey: NEW-1\nname: new\r"},
		"merged in only":      {"b: &b {key: old}\nm:\n  <<: *b\n", "key: NEW-1\nb: &b {key: old}\nm:\n  <<: *b\n"},
	} {
		t.Run(name, func(t *testing.T) { assert.Equal(t, tc.out, set(t, tc.in, "key", "NEW-1")) })
	}
}

// A value that would read as another type is quoted, so that it stays a string.
func TestSetQuotesAValueThatWouldReadAsANumber(t *testing.T) {
	assert.Equal(t, "key: '123'\nname: new\n", set(t, "name: new\n", "key", "123"))
	assert.Equal(t, "{key: \"123\", name: new}\n", set(t, "{name: new}\n", "key", "123"))
}

// A key the object holds gets the value where it is, in the tree, which keeps comments,
// anchors, quoting and long lines; blank lines go.
func TestSetGivesAKeyTheObjectHoldsTheValueWhereItIs(t *testing.T) {
	long := strings.Repeat("word ", 30) + "end"
	for name, tc := range map[string]struct{ in, out string }{
		"empty":         {"# head\nkey: '' # was empty\n\n# about the name\nname: new\nbase: &b {x: 1}\nuse: *b\n", "# head\nkey: NEW-1 # was empty\n# about the name\nname: new\nbase: &b {x: 1}\nuse: *b\n"},
		"further down":  {"name: new # name\nkey: ''\n", "name: new # name\nkey: NEW-1\n"},
		"null":          {"name: new\nkey:\n", "name: new\nkey: NEW-1\n"},
		"comment below": {"name: new\nkey: \"\"\n# end of file\n", "name: new\nkey: NEW-1\n# end of file\n"},
		"marker":        {"# head\n---\nkey: ''\n", "# head\n---\nkey: NEW-1\n"},
		"merge keys":    {"b: &b {x: 1}\nm:\n  <<: *b\n  y: 2\nkey: ''\n", "b: &b {x: 1}\nm:\n  <<: *b\n  y: 2\nkey: NEW-1\n"},
		"long line":     {"key: ''\nname: " + long + "\n", "key: NEW-1\nname: " + long + "\n"},
		"alias key":     {"a: &k key\n*k : ''\n", "a: &k key\n*k : NEW-1\n"},
		"key twice":     {"key: ''\nname: x\nkey: ''\n", "key: ''\nname: x\nkey: NEW-1\n"},
	} {
		t.Run(name, func(t *testing.T) { assert.Equal(t, tc.out, set(t, tc.in, "key", "NEW-1")) })
	}
}

// A new key of a flow mapping, or of a mapping with an anchor or a tag, goes first in the
// tree: a line in front would not belong to the mapping.
func TestSetPutsANewKeyFirstInTheTreeOfAFlowOrMarkedMapping(t *testing.T) {
	for name, tc := range map[string]struct{ in, out string }{
		"flow":                  {"# head\n{name: new}  # flow\n", "# head\n{key: NEW-1, name: new} # flow\n"},
		"flow on the marker":    {"--- {name: new}\n", "--- {key: NEW-1, name: new}\n"},
		"comment before marker": {"# c\n--- {name: new}\n", "# c\n--- {key: NEW-1, name: new}\n"},
		"anchored":              {"&a\nname: new\n", "&a\nkey: NEW-1\nname: new\n"},
		"anchored on marker":    {"# c\n--- &a\nname: new\n", "# c\n--- &a\nkey: NEW-1\nname: new\n"},
		"tagged":                {"!!map\nname: new\n", "!!map\nkey: NEW-1\nname: new\n"},
	} {
		t.Run(name, func(t *testing.T) { assert.Equal(t, tc.out, set(t, tc.in, "key", "NEW-1")) })
	}
}

// Where the tree cannot change without changing more than the key, the object is written
// out again from its values.
func TestSetWritesTheObjectOutWhereTheTreeWouldChangeMore(t *testing.T) {
	for name, tc := range map[string]struct{ in, out string }{
		"anchored value":    {"key: &k ''\nname: x\nother: *k\n", "key: NEW-1\nname: x\nother: ''\n"},
		"anchor inside it":  {"key:\n  a: &x 1\nother: *x\n", "key: NEW-1\nother: 1\n"},
		"a second document": {"key: ''\nname: x\n---\nother: y\n", "key: NEW-1\nname: x\n"},
	} {
		t.Run(name, func(t *testing.T) { assert.Equal(t, tc.out, set(t, tc.in, "key", "NEW-1")) })
	}
}

func TestSetKeepsTheStyleOfJSON(t *testing.T) {
	for name, tc := range map[string]struct{ in, out string }{
		"compact":         {`{"name":"new"}`, `{"key":"NEW-1","name":"new"}`},
		"indented":        {"{\n  \"name\": \"new\"\n}", "{\n  \"key\": \"NEW-1\",\n  \"name\": \"new\"\n}"},
		"final newline":   {"{\"name\":\"new\"}\n", "{\"key\":\"NEW-1\",\"name\":\"new\"}\n"},
		"byte order mark": {"\ufeff{\"name\":\"new\"}", "\ufeff{\"key\":\"NEW-1\",\"name\":\"new\"}"},
		"held key":        {`{"name":"new","key":""}`, `{"name":"new","key":"NEW-1"}`},
	} {
		t.Run(name, func(t *testing.T) { assert.Equal(t, tc.out, set(t, tc.in, "key", "NEW-1")) })
	}
}

func TestDocumentIsACopy(t *testing.T) {
	f, err := ParseFile([]byte("name: new\n"))
	require.NoError(t, err)
	f.Document().Value().SetFirst("key", "changed")
	assert.Equal(t, "name: new\n", string(f.Bytes()))
	_, has := f.Document().Get("key")
	assert.False(t, has)
}
