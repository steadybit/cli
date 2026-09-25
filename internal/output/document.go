// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"go.yaml.in/yaml/v3"
)

// Document is a JSON or YAML document that keeps its field order. Files written by
// `get` are kept in Git, so the order the platform returns has to survive: a Go map
// would sort the keys and turn every upgrade of the CLI into a diff of every file.
type Document struct {
	node *yaml.Node
}

// ParseDocument reads JSON or YAML; JSON is valid YAML, so one parser handles both.
// Anchors and merge keys (`<<:`) are resolved when the document is turned into JSON.
func ParseDocument(content []byte) (*Document, error) {
	var node yaml.Node
	if err := yaml.Unmarshal(content, &node); err != nil {
		return nil, err
	}
	if node.Kind != yaml.DocumentNode || len(node.Content) != 1 || node.Content[0].Kind != yaml.MappingNode {
		return nil, fmt.Errorf("expected an object")
	}
	return &Document{node: node.Content[0]}, nil
}

func (d *Document) Get(key string) (string, bool) {
	for i := 0; i+1 < len(d.node.Content); i += 2 {
		if d.node.Content[i].Value == key {
			return d.node.Content[i+1].Value, true
		}
	}
	return "", false
}

func (d *Document) Delete(key string) {
	for i := 0; i+1 < len(d.node.Content); i += 2 {
		if d.node.Content[i].Value == key {
			d.node.Content = append(d.node.Content[:i], d.node.Content[i+2:]...)
			return
		}
	}
}

// SetFirst sets a string field, moving it to the top, where `key` and `id` belong.
func (d *Document) SetFirst(key, value string) {
	d.Delete(key)
	d.node.Content = append([]*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: "!!str", Value: key},
		{Kind: yaml.ScalarNode, Tag: "!!str", Value: value},
	}, d.node.Content...)
}

func (d *Document) Render(datatype Datatype) ([]byte, error) {
	if datatype == JSON {
		var buf bytes.Buffer
		if err := writeJSON(&buf, d.node, ""); err != nil {
			return nil, err
		}
		buf.WriteByte('\n')
		return buf.Bytes(), nil
	}
	restyle(d.node)
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(d.node); err != nil {
		return nil, err
	}
	return buf.Bytes(), encoder.Close()
}

// MarshalJSON lets a document be sent as a request body.
func (d *Document) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	err := writeJSON(&buf, d.node, "")
	return buf.Bytes(), err
}

// restyle drops the flow style and double quotes that parsing JSON leaves on every node,
// and quotes a string the way js-yaml did: single quotes when it would otherwise read as
// something else, double quotes only when it needs escapes.
func restyle(node *yaml.Node) {
	node.Style = 0
	// JavaScript has one number type, so the platform's 1.0 was always written as 1.
	if node.Kind == yaml.ScalarNode && node.ShortTag() == "!!float" {
		if f, err := strconv.ParseFloat(node.Value, 64); err == nil {
			node.Value = strconv.FormatFloat(f, 'f', -1, 64)
			if f == float64(int64(f)) {
				node.Tag = "!!int"
			}
		}
	}
	if node.Kind == yaml.ScalarNode && node.Tag == "!!str" && needsQuotes(node.Value) {
		if strings.ContainsAny(node.Value, "\n\t\\") || !strconv.IsPrint(firstUnprintable(node.Value)) {
			node.Style = yaml.DoubleQuotedStyle
		} else {
			node.Style = yaml.SingleQuotedStyle
		}
	}
	if node.Kind == yaml.ScalarNode && node.Tag == "!!str" && strings.Contains(node.Value, "\n") {
		node.Style = yaml.LiteralStyle
	}
	for _, child := range node.Content {
		restyle(child)
	}
}

func firstUnprintable(s string) rune {
	for _, r := range s {
		if !strconv.IsPrint(r) {
			return r
		}
	}
	return 'a'
}

// needsQuotes reports whether a plain scalar would not read back as this string.
func needsQuotes(value string) bool {
	plain := &yaml.Node{Kind: yaml.ScalarNode, Value: value}
	out, err := yaml.Marshal(plain)
	if err != nil {
		return true
	}
	rendered := strings.TrimSuffix(string(out), "\n")
	if rendered != value {
		return true
	}
	var resolved yaml.Node
	if err := yaml.Unmarshal(out, &resolved); err != nil || len(resolved.Content) != 1 {
		return true
	}
	return resolved.Content[0].Tag != "!!str" || isTimestamp(value)
}

// js-yaml's default schema reads timestamps as dates, so it quoted them; so do we.
func isTimestamp(value string) bool {
	return len(value) >= 10 && value[4] == '-' && value[7] == '-' && strings.IndexFunc(value[:4], func(r rune) bool { return r < '0' || r > '9' }) < 0
}

// writeJSON renders a node as indented JSON in its own field order, resolving aliases
// and merge keys on the way.
func writeJSON(buf *bytes.Buffer, node *yaml.Node, indent string) error {
	switch node.Kind {
	case yaml.AliasNode:
		return writeJSON(buf, node.Alias, indent)
	case yaml.MappingNode:
		pairs := mergedPairs(node)
		if len(pairs) == 0 {
			buf.WriteString("{}")
			return nil
		}
		buf.WriteString("{\n")
		for i, pair := range pairs {
			key := jsonString(pair[0].Value)
			buf.WriteString(indent + "  ")
			buf.Write(key)
			buf.WriteString(": ")
			if err := writeJSON(buf, pair[1], indent+"  "); err != nil {
				return err
			}
			if i < len(pairs)-1 {
				buf.WriteByte(',')
			}
			buf.WriteByte('\n')
		}
		buf.WriteString(indent + "}")
	case yaml.SequenceNode:
		if len(node.Content) == 0 {
			buf.WriteString("[]")
			return nil
		}
		buf.WriteString("[\n")
		for i, item := range node.Content {
			buf.WriteString(indent + "  ")
			if err := writeJSON(buf, item, indent+"  "); err != nil {
				return err
			}
			if i < len(node.Content)-1 {
				buf.WriteByte(',')
			}
			buf.WriteByte('\n')
		}
		buf.WriteString(indent + "]")
	case yaml.ScalarNode:
		switch node.ShortTag() {
		case "!!null":
			buf.WriteString("null")
		case "!!bool", "!!int", "!!float":
			var v any
			if err := node.Decode(&v); err != nil {
				return err
			}
			b, _ := json.Marshal(v)
			buf.Write(b)
		default:
			buf.Write(jsonString(node.Value))
		}
	default:
		return fmt.Errorf("unsupported YAML node")
	}
	return nil
}

// jsonString encodes like JSON.stringify: `&`, `<` and `>` stay as they are.
func jsonString(value string) []byte {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(value)
	return bytes.TrimSuffix(buf.Bytes(), []byte("\n"))
}

// mergedPairs returns a mapping's key/value pairs with `<<` merge keys expanded; keys
// written in the mapping itself win over merged ones, as YAML specifies.
func mergedPairs(node *yaml.Node) [][2]*yaml.Node {
	var own, merged [][2]*yaml.Node
	seen := map[string]bool{}
	for i := 0; i+1 < len(node.Content); i += 2 {
		key, value := node.Content[i], node.Content[i+1]
		if key.Value == "<<" && key.Tag == "!!merge" {
			sources := []*yaml.Node{value}
			if resolve(value).Kind == yaml.SequenceNode {
				sources = resolve(value).Content
			}
			for _, source := range sources {
				merged = append(merged, mergedPairs(resolve(source))...)
			}
			continue
		}
		own = append(own, [2]*yaml.Node{key, value})
		seen[key.Value] = true
	}
	for _, pair := range merged {
		if !seen[pair[0].Value] {
			own = append(own, pair)
			seen[pair[0].Value] = true
		}
	}
	return own
}

func resolve(node *yaml.Node) *yaml.Node {
	for node.Kind == yaml.AliasNode {
		node = node.Alias
	}
	return node
}
