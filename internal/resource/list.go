// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package resource

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/steadybit/cli/internal/jsyaml"
	"github.com/steadybit/cli/internal/output"
	"github.com/steadybit/cli/internal/platform"
)

// ListTypeHelp describes the -t flag of listings.
const ListTypeHelp = `Print the platform's items as "json" or "yaml" instead of a table.`

// Machine reports whether a listing prints its items rather than a table: with -t, or
// when --jq filters them.
func Machine(explicitType string) bool {
	return explicitType != "" || output.JQ != ""
}

// List prints items, the platform's objects as it sent them, as JSON or YAML, or hands
// over to table when a person is reading.
func List(items []json.RawMessage, explicitType string, table func() error) error {
	if !Machine(explicitType) {
		return table()
	}
	parts := make([]string, len(items))
	for i, item := range items {
		parts[i] = string(item)
	}
	return PrintJSONValue([]byte("["+strings.Join(parts, ",")+"]"), explicitType)
}

// PrintJSONValue prints any JSON value as -t and --jq ask: through the jq expression,
// or as YAML, or as indented JSON.
func PrintJSONValue(raw []byte, explicitType string) error {
	value, err := output.ParseValue(raw)
	if err != nil {
		return err
	}
	if output.JQ != "" {
		return output.ApplyJQ(os.Stdout, jsyaml.CompactJSON(value), output.JQ)
	}
	if explicitType == "yaml" {
		fmt.Print(jsyaml.Dump(value))
		return nil
	}
	if explicitType != "" && explicitType != "json" {
		return fmt.Errorf("Unsupported output format '%s'. Use \"json\" or \"yaml\".", explicitType)
	}
	fmt.Println(jsyaml.JSON(value))
	return nil
}

// DecodeEach decodes raw items into typed ones, for the table.
func DecodeEach[T any](raw []json.RawMessage, into *[]T) error {
	for _, item := range raw {
		var v T
		if err := json.Unmarshal(item, &v); err != nil {
			return err
		}
		*into = append(*into, v)
	}
	return nil
}

// DecodeListed reads a listing response once: into typed, for the table, and as the raw
// items under field ("" for a response that is the array itself), for -t and --jq.
func DecodeListed(resp *http.Response, err error, field string, typed any) ([]json.RawMessage, error) {
	body, _, err := platform.Read(resp, err)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, typed); err != nil {
		return nil, err
	}
	var raw []json.RawMessage
	if field == "" {
		err = json.Unmarshal(body, &raw)
	} else {
		var wrapper map[string]json.RawMessage
		if err = json.Unmarshal(body, &wrapper); err == nil && wrapper[field] != nil {
			err = json.Unmarshal(wrapper[field], &raw)
		}
	}
	return raw, err
}
