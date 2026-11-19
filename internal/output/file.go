// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package output

import (
	"bytes"
	"slices"

	"github.com/steadybit/cli/v6/internal/jsyaml"
	"go.yaml.in/yaml/v4"
)

// File is a file of an object, as kept in Git: its text and the object it holds. A change
// to the object goes into the text and changes as little else of it as it can.
type File struct {
	bom      bool
	text     []byte
	datatype Datatype
	doc      *Document
}

// ParseFile reads a file's content as JSON or YAML; JSON is tried first, as it was.
func ParseFile(content []byte) (*File, error) {
	doc, err := ParseDocument(content)
	if err != nil {
		return nil, err
	}
	datatype := YAML
	if IsJSON(content) {
		datatype = JSON
	}
	text, hasBOM := bytes.CutPrefix(content, bom)
	return &File{bom: hasBOM, text: text, datatype: datatype, doc: doc}, nil
}

// Document is a copy of the object the file holds: changing it changes neither the file
// nor its text.
func (f *File) Document() *Document {
	return NewDocument(jsyaml.Clone(f.doc.value).(*jsyaml.Map))
}

// Bytes is the content of the file, with the changes made to it.
func (f *File) Bytes() []byte {
	if f.bom {
		return append(append([]byte{}, bom...), f.text...)
	}
	return f.text
}

// Set sets key to value in the object, as reading the file back gives it: a key the object
// holds gets the value where it is, a key it does not hold goes first. YAML keeps the rest
// of the file as far as the change allows (see setInDocument); JSON is written out again in the
// style it was in.
func (f *File) Set(key, value string) error {
	if f.datatype == YAML {
		edited, err := setInDocument(f.text, f.doc.value, key, value)
		if err != nil {
			return err
		}
		if edited == nil {
			edited = f.doc.RenderFile(YAML)
		}
		f.text = edited
		return nil
	}
	if _, holds := f.doc.value.Get(key); holds {
		f.doc.value.Set(key, value)
	} else {
		f.doc.value.SetFirst(key, value)
	}
	f.text = renderJSONLike(f.text, f.doc)
	return nil
}

// renderJSONLike writes the object as JSON in the style of text: compact or indented, with
// or without a final newline.
func renderJSONLike(text []byte, doc *Document) []byte {
	var rendered string
	if bytes.Contains(bytes.TrimSpace(text), []byte("\n")) {
		rendered = jsyaml.JSON(doc.value)
	} else {
		rendered = jsyaml.CompactJSON(doc.value)
	}
	if bytes.HasSuffix(text, []byte("\n")) {
		rendered += "\n"
	}
	return []byte(rendered)
}

// setInDocument sets the key in YAML text holding one document, and in the object it
// holds, as reading the text back gives it, as the document's syntax tree says it can be
// with the least change.
//
// A new key of a block mapping is a line of its own at the start of the document: after a
// `---` marker, comments and directives before it staying above, and at the mapping's
// indentation; the rest of the text stays byte for byte. A key the mapping holds gets its
// value in the tree, in place, as does a new key of a flow mapping `{...}` or of a mapping
// with an anchor or tag, and the document is written out again from the tree: comments,
// anchors, quoting and long lines stay, blank lines go, and indentation becomes two spaces.
//
// It returns nil where the tree cannot change without changing more, for the object to be
// written out again from its values: where an alias may refer to the value, and where the
// tree would drop a further document.
func setInDocument(text []byte, object *jsyaml.Map, key, value string) ([]byte, error) {
	var root yaml.Node
	if err := yaml.Unmarshal(text, &root); err != nil {
		return nil, err
	}
	mapping := root.Content[0]
	held := heldValue(mapping, key)
	if held != nil {
		object.Set(key, value)
	} else {
		object.SetFirst(key, value)
	}
	if held != nil && hasAnchor(held) {
		return nil, nil
	}
	lines := lineStarts(text)
	marker := markerLine(text, lines, root.Line, mapping.Line)
	if held == nil && mapping.Style&(yaml.FlowStyle|yaml.TaggedStyle) == 0 && mapping.Anchor == "" {
		at := 0
		if marker > 0 {
			at = lines[marker+1]
		}
		return slices.Concat(text[:at], keyLine(mapping.Column, key, value), text[at:]), nil
	}
	// The tree is written out from the line after a marker of its own, or from the marker
	// the mapping starts on, which the dumper does not write.
	at, prefix := 0, []byte(nil)
	switch {
	case marker > 0 && marker < mapping.Line:
		at = lines[marker+1]
	case marker > 0:
		at, prefix = lines[marker], []byte("--- ")
	}
	var tree yaml.Node
	if err := yaml.Load(text[at:], &tree); err != nil {
		return nil, nil
	}
	setInTree(tree.Content[0], key, value)
	out, err := yaml.Dump(&tree, yaml.WithV3Defaults(), yaml.WithIndent(2), yaml.WithLineWidth(-1))
	if err != nil {
		return nil, err
	}
	return slices.Concat(text[:at], prefix, out), nil
}

// heldValue is the value node of the key as reading the mapping gives it: of its last
// occurrence, with alias keys resolved. Keys merged in with `<<` are not the mapping's own.
func heldValue(mapping *yaml.Node, key string) *yaml.Node {
	var held *yaml.Node
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		k := mapping.Content[i]
		for k.Kind == yaml.AliasNode {
			k = k.Alias
		}
		if k.Kind == yaml.ScalarNode && k.Value == key && k.ShortTag() != "!!merge" {
			held = mapping.Content[i+1]
		}
	}
	return held
}

// hasAnchor tells whether an alias may refer to the node or anything in it.
func hasAnchor(node *yaml.Node) bool {
	return node.Anchor != "" || slices.ContainsFunc(node.Content, hasAnchor)
}

// lineStarts are the offsets the lines of text start at, by 1-based line number, counting
// line breaks as YAML does: "\r\n", "\r" and "\n". The entry after the last line is the end
// of the text.
func lineStarts(text []byte) []int {
	starts := []int{0, 0}
	for i := 0; i < len(text); i++ {
		switch text[i] {
		case '\r':
			if i+1 < len(text) && text[i+1] == '\n' {
				i++
			}
		case '\n':
		default:
			continue
		}
		if i+1 < len(text) {
			starts = append(starts, i+1)
		}
	}
	return append(starts, len(text))
}

// markerLine is the line of the `---` marker opening the document, between the line the
// document starts on and the line its mapping starts on, or 0 for a document without one.
// A document starts with its directives; one without a marker starts with its content.
func markerLine(text []byte, lines []int, from, to int) int {
	for line := min(to, len(lines)-2); line >= from && line > 0; line-- {
		rest := text[lines[line]:lines[line+1]]
		if bytes.HasPrefix(rest, []byte("---")) && (len(rest) == 3 || bytes.IndexByte([]byte(" \t\r\n"), rest[3]) >= 0) {
			return line
		}
	}
	return 0
}

// keyLine is `key: value` as a line of a block mapping at column, quoted as js-yaml quotes
// a value.
func keyLine(column int, key, value string) []byte {
	line := jsyaml.NewMap()
	line.Set(key, value)
	return append(bytes.Repeat([]byte(" "), column-1), jsyaml.Dump(line)...)
}

// setInTree gives a key the mapping holds the value in place, keeping its comments, and
// puts a key it does not hold first.
func setInTree(mapping *yaml.Node, key, value string) {
	if held := heldValue(mapping, key); held != nil {
		held.SetString(value)
		held.Style &^= yaml.SingleQuotedStyle | yaml.DoubleQuotedStyle
		return
	}
	k, v := &yaml.Node{}, &yaml.Node{}
	k.SetString(key)
	v.SetString(value)
	mapping.Content = append([]*yaml.Node{k, v}, mapping.Content...)
}
