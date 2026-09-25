// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package output formats documents and gates colour on stdout being a terminal, as the
// TypeScript CLI did: its output is routinely parsed by GitOps pipelines.
package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"go.yaml.in/yaml/v3"
	"golang.org/x/term"
)

var colorsEnabled = os.Getenv("NO_COLOR") == "" &&
	(os.Getenv("FORCE_COLOR") != "" || term.IsTerminal(int(os.Stdout.Fd())))

func style(code, s string) string {
	if !colorsEnabled {
		return s
	}
	return "\x1b[" + code + "m" + s + "\x1b[0m"
}

func Bold(s string) string  { return style("1", s) }
func Red(s string) string   { return style("31", s) }
func Green(s string) string { return style("32", s) }

type Datatype string

const (
	JSON Datatype = "json"
	YAML Datatype = "yaml"
)

// ResolveDatatype: an explicit type first, then the file's extension, then YAML.
func ResolveDatatype(explicit, file string) (Datatype, error) {
	switch explicit {
	case "json", "yaml":
		return Datatype(explicit), nil
	case "":
		if strings.HasSuffix(strings.ToLower(file), ".json") {
			return JSON, nil
		}
		return YAML, nil
	default:
		return "", fmt.Errorf("unsupported output format '%s'. Use \"json\" or \"yaml\"", explicit)
	}
}

// Format renders a document. JSON is decoded with UseNumber first, so that large
// numbers survive the trip through Go values unchanged.
func Format(document any, datatype Datatype) ([]byte, error) {
	if datatype == JSON {
		return json.MarshalIndent(document, "", "  ")
	}
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(document); err != nil {
		return nil, err
	}
	return buf.Bytes(), encoder.Close()
}

// Parse reads JSON or YAML. YAML anchors and merge keys (`<<:`) are resolved, which is
// what the TypeScript CLI's schema was configured to do for experiment files.
func Parse(content []byte) (map[string]any, Datatype, error) {
	var document map[string]any
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.UseNumber()
	if err := decoder.Decode(&document); err == nil {
		return document, JSON, nil
	}
	var node yaml.Node
	if err := yaml.Unmarshal(content, &node); err != nil {
		return nil, "", err
	}
	if err := node.Decode(&document); err != nil {
		return nil, "", err
	}
	return normalize(document).(map[string]any), YAML, nil
}

// yaml.v3 decodes nested maps as map[string]any already, but keeps integers as int,
// which is fine; this only exists to turn map[any]any from older documents into JSON-able maps.
func normalize(value any) any {
	switch v := value.(type) {
	case map[string]any:
		for k, item := range v {
			v[k] = normalize(item)
		}
		return v
	case map[any]any:
		m := make(map[string]any, len(v))
		for k, item := range v {
			m[fmt.Sprint(k)] = normalize(item)
		}
		return m
	case []any:
		for i, item := range v {
			v[i] = normalize(item)
		}
		return v
	default:
		return v
	}
}
