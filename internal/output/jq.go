// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package output

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/itchyny/gojq"
)

// JQ is the --jq expression, applied to whatever JSON a command prints. Set once from the
// root command's flags.
var JQ string

// ApplyJQ filters JSON through the --jq expression and writes each result on its own
// line: strings raw, everything else as indented JSON, as `gh --jq` does. jq itself is
// not needed, the expression runs in-process.
func ApplyJQ(w io.Writer, jsonText string, expression string) error {
	query, err := gojq.Parse(expression)
	if err != nil {
		return fmt.Errorf("Invalid --jq expression: %w", err)
	}
	code, err := gojq.Compile(query, gojq.WithEnvironLoader(func() []string { return nil }))
	if err != nil {
		return fmt.Errorf("Invalid --jq expression: %w", err)
	}
	var input any
	decoder := json.NewDecoder(strings.NewReader(jsonText))
	decoder.UseNumber()
	if err := decoder.Decode(&input); err != nil {
		return err
	}
	iter := code.Run(normalizeNumbers(input))
	for {
		v, ok := iter.Next()
		if !ok {
			return nil
		}
		if err, isErr := v.(error); isErr {
			var halt *gojq.HaltError
			if errors.As(err, &halt) && halt.Value() == nil {
				return nil
			}
			return fmt.Errorf("--jq: %w", err)
		}
		if s, isString := v.(string); isString {
			fmt.Fprintln(w, s)
			continue
		}
		out, err := gojq.Marshal(v)
		if err != nil {
			return err
		}
		var pretty bytes.Buffer
		if err := json.Indent(&pretty, out, "", "  "); err != nil {
			return err
		}
		fmt.Fprintln(w, pretty.String())
	}
}

// gojq works on float64 and big numbers, not json.Number.
func normalizeNumbers(v any) any {
	switch x := v.(type) {
	case json.Number:
		if i, err := x.Int64(); err == nil {
			return int(i)
		}
		f, _ := x.Float64()
		return f
	case map[string]any:
		for k, item := range x {
			x[k] = normalizeNumbers(item)
		}
		return x
	case []any:
		for i, item := range x {
			x[i] = normalizeNumbers(item)
		}
		return x
	}
	return v
}
