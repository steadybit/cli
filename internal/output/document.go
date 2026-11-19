// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package output

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/steadybit/cli/v6/internal/jsyaml"
	"go.yaml.in/yaml/v4"
)

// Document is a JSON or YAML object as the TypeScript CLI held it: a JavaScript value
// with JavaScript's key order and number semantics. Files written from it are kept in
// Git, so they have to come out exactly as they did before.
type Document struct {
	value *jsyaml.Map
}

// ParseList reads a JSON array, or JSON values one after another as newline-delimited
// JSON sends them. The platform streams the audit log that way, although its spec says
// array, so both are read.
func ParseList(content []byte) ([]any, error) {
	decoder := json.NewDecoder(bytes.NewReader(bytes.TrimPrefix(content, bom)))
	values := []any{}
	for {
		value, err := decodeJSON(decoder)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	if len(values) == 1 {
		if list, ok := values[0].([]any); ok {
			return list, nil
		}
	}
	return values, nil
}

// ParseValue reads any JSON or YAML value, such as a placeholders file that is a list.
func ParseValue(content []byte) (any, error) {
	content = bytes.TrimPrefix(content, bom)
	if json.Valid(content) {
		return decodeJSON(json.NewDecoder(bytes.NewReader(content)))
	}
	nodes, err := yamlDocuments(content)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 0:
		return nil, nil
	case 1:
		return toValue(nodes[0].Content[0])
	}
	return nil, errMoreThanOne(len(nodes))
}

// ParseDocument reads one object as JSON or YAML; JSON is valid YAML, so one parser handles
// both. Anchors, aliases and merge keys (`<<:`) are resolved, as js-yaml's load did.
func ParseDocument(content []byte) (*Document, error) {
	docs, err := ParseDocuments(content)
	if err != nil {
		return nil, err
	}
	if len(docs) > 1 {
		return nil, errMoreThanOne(len(docs))
	}
	return docs[0], nil
}

// ParseDocuments reads the objects of a file: JSON holds one, YAML one per document, so a
// file can keep several objects of a kind apart with `---`. Empty documents, as a trailing
// `---` leaves, hold none.
func ParseDocuments(content []byte) ([]*Document, error) {
	// JSON is read with a JSON decoder: a YAML parser rejects characters JSON allows
	// unescaped, such as DEL, and would fail on a platform response containing one.
	content = bytes.TrimPrefix(content, bom)
	if json.Valid(content) {
		value, err := decodeJSON(json.NewDecoder(bytes.NewReader(content)))
		if err != nil {
			return nil, err
		}
		m, ok := value.(*jsyaml.Map)
		if !ok {
			return nil, fmt.Errorf("expected an object")
		}
		return []*Document{{value: m}}, nil
	}
	_, docs, err := yamlObjects(content)
	return docs, err
}

// yamlObjects decodes the documents of a YAML file into objects, keeping the nodes they
// came from for the positions of their keys.
func yamlObjects(content []byte) ([]*yaml.Node, []*Document, error) {
	nodes, err := yamlDocuments(content)
	if err != nil {
		return nil, nil, err
	}
	if len(nodes) == 0 {
		return nil, nil, fmt.Errorf("expected an object")
	}
	docs := make([]*Document, 0, len(nodes))
	for i, node := range nodes {
		value, err := toValue(node.Content[0])
		if err != nil {
			return nil, nil, inDocument(err, i, len(nodes))
		}
		m, ok := value.(*jsyaml.Map)
		if !ok {
			return nil, nil, inDocument(fmt.Errorf("expected an object"), i, len(nodes))
		}
		docs = append(docs, &Document{value: m})
	}
	return nodes, docs, nil
}

// yamlDocuments decodes every document of a YAML file, leaving out empty ones.
// yaml.Unmarshal would read only the first and drop the rest without a word.
func yamlDocuments(content []byte) ([]*yaml.Node, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(content))
	var nodes []*yaml.Node
	for {
		var node yaml.Node
		err := decoder.Decode(&node)
		if errors.Is(err, io.EOF) {
			return nodes, nil
		}
		if err != nil {
			return nil, err
		}
		if len(node.Content) != 1 || (node.Content[0].Kind == yaml.ScalarNode && node.Content[0].Tag == "!!null") {
			continue
		}
		nodes = append(nodes, &node)
	}
}

func errMoreThanOne(n int) error {
	return fmt.Errorf("found %d documents where one object is expected", n)
}

func inDocument(err error, index, count int) error {
	if count == 1 {
		return err
	}
	return fmt.Errorf("document %d: %w", index+1, err)
}

// Label names one object of a file in messages: the file alone when it holds one.
func Label(file string, index, count int) string {
	if count == 1 {
		return file
	}
	return fmt.Sprintf("%s (document %d)", file, index+1)
}

func NewDocument(m *jsyaml.Map) *Document { return &Document{value: m} }

func (d *Document) Value() *jsyaml.Map { return d.value }

func (d *Document) Get(key string) (string, bool) {
	v, ok := d.value.Get(key)
	s, isString := v.(string)
	return s, ok && isString
}

func (d *Document) Delete(key string) { d.value.Delete(key) }

// Render as `get` printed it: JSON.stringify(value, undefined, 2), or js-yaml's dump.
func (d *Document) Render(datatype Datatype) ([]byte, error) {
	if datatype == JSON {
		return []byte(jsyaml.JSON(d.value) + "\n"), nil
	}
	return []byte(jsyaml.Dump(d.value)), nil
}

// RenderFile as the TypeScript CLI wrote files: compact JSON without a newline, or YAML.
func (d *Document) RenderFile(datatype Datatype) []byte {
	if datatype == JSON {
		return []byte(jsyaml.CompactJSON(d.value))
	}
	return []byte(jsyaml.Dump(d.value))
}

// MarshalJSON sends a document as a request body, as JSON.stringify did.
func (d *Document) MarshalJSON() ([]byte, error) {
	return []byte(jsyaml.CompactJSON(d.value)), nil
}

// decodeJSON reads the next JSON value in order, with numbers as JavaScript sees them.
func decodeJSON(decoder *json.Decoder) (any, error) {
	decoder.UseNumber()
	token, err := decoder.Token()
	if err != nil {
		return nil, err
	}
	switch t := token.(type) {
	case json.Delim:
		switch t {
		case '{':
			m := jsyaml.NewMap()
			for decoder.More() {
				keyToken, err := decoder.Token()
				if err != nil {
					return nil, err
				}
				value, err := decodeJSON(decoder)
				if err != nil {
					return nil, err
				}
				m.Set(keyToken.(string), value)
			}
			_, err := decoder.Token()
			return m, err
		case '[':
			items := []any{}
			for decoder.More() {
				item, err := decodeJSON(decoder)
				if err != nil {
					return nil, err
				}
				items = append(items, item)
			}
			_, err := decoder.Token()
			return items, err
		}
	case json.Number:
		return jsNumber(string(t)), nil
	case string, bool, nil:
		return t, nil
	}
	if errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("unexpected end of JSON")
	}
	return nil, fmt.Errorf("unexpected JSON token %v", token)
}

// jsNumber is Number(text): the nearest double, with the sign of zero kept.
func jsNumber(text string) float64 {
	f, _ := strconv.ParseFloat(text, 64)
	if f == 0 && strings.HasPrefix(text, "-") {
		return math.Copysign(0, -1)
	}
	return f
}

func toValue(node *yaml.Node) (any, error) {
	switch node.Kind {
	case yaml.AliasNode:
		return toValue(node.Alias)
	case yaml.MappingNode:
		m := jsyaml.NewMap()
		for _, pair := range mergedPairs(node) {
			key, err := toValue(pair[0])
			if err != nil {
				return nil, err
			}
			value, err := toValue(pair[1])
			if err != nil {
				return nil, err
			}
			m.Set(keyString(key), value)
		}
		return m, nil
	case yaml.SequenceNode:
		items := make([]any, 0, len(node.Content))
		for _, child := range node.Content {
			item, err := toValue(child)
			if err != nil {
				return nil, err
			}
			items = append(items, item)
		}
		return items, nil
	case yaml.ScalarNode:
		return scalarValue(node)
	}
	return nil, fmt.Errorf("unsupported YAML node at line %d", node.Line)
}

// keyString is the property name JavaScript would use for a key of any type.
func keyString(key any) string {
	switch k := key.(type) {
	case string:
		return k
	case nil:
		return "null"
	case bool:
		return strconv.FormatBool(k)
	case float64:
		return jsyaml.NumberString(k)
	case jsyaml.Timestamp:
		return time.Time(k).UTC().Format("Mon Jan 02 2006 15:04:05 GMT+0000 (Coordinated Universal Time)")
	}
	return fmt.Sprint(key)
}

func scalarValue(node *yaml.Node) (any, error) {
	switch node.ShortTag() {
	case "!!null":
		return nil, nil
	case "!!bool":
		var b bool
		err := node.Decode(&b)
		return b, err
	case "!!int":
		// JavaScript has one number type; big integers lose precision just as they did.
		var i int64
		if err := node.Decode(&i); err == nil {
			return float64(i), nil
		}
		var u uint64
		if err := node.Decode(&u); err == nil {
			return float64(u), nil
		}
		f, err := strconv.ParseFloat(strings.ReplaceAll(node.Value, "_", ""), 64)
		return f, err
	case "!!float":
		var f float64
		if err := node.Decode(&f); err != nil {
			return math.NaN(), err
		}
		return f, nil
	case "!!timestamp":
		var t time.Time
		if err := node.Decode(&t); err != nil {
			return node.Value, nil
		}
		return jsyaml.Timestamp(t), nil
	default:
		return node.Value, nil
	}
}

// mergedPairs returns a mapping's key/value pairs with `<<` merge keys expanded; keys
// written in the mapping itself win over merged ones, as YAML specifies.
func mergedPairs(node *yaml.Node) [][2]*yaml.Node {
	var own, merged [][2]*yaml.Node
	seen := map[string]bool{}
	for i := 0; i+1 < len(node.Content); i += 2 {
		key, value := node.Content[i], node.Content[i+1]
		if key.Value == "<<" && key.ShortTag() == "!!merge" {
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

// IsJSON reports whether content is a JSON document, as JSON.parse would accept it. A
// leading byte order mark, which Windows editors write, is ignored.
func IsJSON(content []byte) bool { return json.Valid(bytes.TrimPrefix(content, bom)) }

var bom = []byte("\ufeff")
