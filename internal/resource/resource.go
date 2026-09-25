// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package resource holds what the commands managing schedules, services, profiles,
// templates and runs share: writing a document to a file or stdout, reading one back,
// and applying files with the new id written into them.
package resource

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/steadybit/cli/internal/experiment"
	"github.com/steadybit/cli/internal/jsyaml"
	"github.com/steadybit/cli/internal/output"
	"github.com/steadybit/cli/internal/prompt"
)

// Output writes to the file when one is given and to stdout otherwise, as JSON
// indented by two, or YAML.
func Output(doc *output.Document, file, explicitType string) error {
	datatype, err := output.ResolveDatatype(explicitType, file)
	if err != nil {
		return err
	}
	rendered := format(doc, datatype)
	if file == "" {
		fmt.Println(rendered)
		return nil
	}
	return os.WriteFile(file, []byte(rendered), 0o644)
}

// OutputValue is Output for a value that is not an object, such as a map of variables.
func OutputValue(value *jsyaml.Map, file, explicitType string) error {
	return Output(output.NewDocument(value), file, explicitType)
}

func format(doc *output.Document, datatype output.Datatype) string {
	if datatype == output.JSON {
		return jsyaml.JSON(doc.Value())
	}
	return jsyaml.Dump(doc.Value())
}

// Read loads a JSON or YAML file; JSON is tried first, as it was.
func Read(file, what string) (*output.Document, output.Datatype, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, "", fmt.Errorf("Failed to read %s file at path '%s': %s", what, file, cause(err))
	}
	doc, err := output.ParseDocument(content)
	if err != nil {
		return nil, "", fmt.Errorf("Failed to parse %s file at path '%s' as YAML/JSON: %s", what, file, err)
	}
	datatype := output.YAML
	if isJSON(content) {
		datatype = output.JSON
	}
	return doc, datatype, nil
}

func isJSON(content []byte) bool {
	trimmed := strings.TrimSpace(string(content))
	return strings.HasPrefix(trimmed, "{") && output.IsJSON(content)
}

func cause(err error) string {
	if pathErr, ok := err.(*os.PathError); ok {
		return pathErr.Err.Error()
	}
	return err.Error()
}

// Strip removes fields the platform reports but does not accept back, so that a file
// written by `get` can be applied again unchanged.
func Strip(doc *output.Document, fields ...string) *output.Document {
	for _, f := range fields {
		doc.Delete(f)
	}
	return doc
}

type Applied struct {
	ID      string
	Created bool
}

// ApplyFiles upserts every file. A file without an id gets the new one written into it,
// first, so the next apply updates what this one created instead of creating another.
func ApplyFiles(paths []string, recursive bool, what string, upsert func(file string, doc *output.Document) (Applied, error)) error {
	files, err := experiment.ResolveFiles(paths, recursive)
	if err != nil {
		return err
	}
	for _, file := range files {
		doc, datatype, err := Read(file, what)
		if err != nil {
			return err
		}
		existingID, _ := doc.Value().Get("id")
		result, err := upsert(file, output.NewDocument(jsyaml.Clone(doc.Value()).(*jsyaml.Map)))
		if err != nil {
			return err
		}
		// Resources named by a key, like teams, have no id to write back.
		if (existingID == nil || existingID == "") && result.ID != "" {
			doc.Value().SetFirst("id", result.ID)
			if err := os.WriteFile(file, []byte(format(doc, datatype)), 0o644); err != nil {
				return err
			}
		}
	}
	return nil
}

func CreatedOrUpdated(created bool) string {
	if created {
		return "created"
	}
	return "updated"
}

// Body sends a document or value as the JSON the platform expects.
func Body(value any) io.Reader { return bytes.NewReader([]byte(jsyaml.CompactJSON(value))) }

// Optional leaves an empty filter out of the request.
func Optional(values []string) *[]string {
	if len(values) == 0 {
		return nil
	}
	return &values
}

// UUID parses an id; a malformed one cannot name anything, so callers report it as not found.
func UUID(id string) (openapi_types.UUID, bool) {
	var u openapi_types.UUID
	return u, u.UnmarshalText([]byte(id)) == nil
}

// Confirmed asks before something that cannot be undone, unless --yes was given.
// Without a terminal, as in a pipeline, it goes ahead, as `experiment run` does.
func Confirmed(yes bool, question string) (bool, error) {
	if yes {
		return true, nil
	}
	ok, err := prompt.Confirm(question, false, true)
	if err == nil && !ok {
		fmt.Println("Aborted.")
	}
	return ok, err
}

// Variables merges KEY=VALUE arguments, always strings, over a file's variables, which
// may be lists or select expressions. Nothing given is only allowed when replacing, where
// it removes every variable.
func Variables(pairs []string, file string, replace bool) (*jsyaml.Map, error) {
	given := jsyaml.NewMap()
	for _, pair := range pairs {
		i := strings.Index(pair, "=")
		if i <= 0 {
			return nil, fmt.Errorf("'%s' is not in the form KEY=VALUE.", pair)
		}
		given.Set(pair[:i], pair[i+1:])
	}
	variables := jsyaml.NewMap()
	if file != "" {
		doc, _, err := Read(file, "variables")
		if err != nil {
			return nil, fmt.Errorf("Variables file '%s' must be a map of variable names to values.", file)
		}
		variables = doc.Value()
	}
	for _, k := range given.Keys() {
		v, _ := given.Get(k)
		variables.Set(k, v)
	}
	if variables.Len() == 0 && !replace {
		return nil, errors.New("No variables given. Pass KEY=VALUE arguments or --file.")
	}
	return variables, nil
}

func VariablesOutcome(replace bool) string {
	if replace {
		return "set, all others removed"
	}
	return "set"
}

// Time takes a date, meaning its start in UTC, or a full RFC 3339 time.
func Time(flag, value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	for _, layout := range []string{time.RFC3339, time.DateOnly} {
		if t, err := time.Parse(layout, value); err == nil {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("--%s '%s' is neither a date like 2026-09-01 nor a time like 2026-09-01T12:00:00Z.", flag, value)
}
