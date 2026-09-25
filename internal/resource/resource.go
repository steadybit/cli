// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package resource holds what the commands managing schedules, services, profiles,
// templates and runs share: writing a document to a file or stdout, reading one back,
// and applying files with the new id written into them.
package resource

import (
	"fmt"
	"os"
	"strings"

	"github.com/steadybit/cli/internal/experiment"
	"github.com/steadybit/cli/internal/jsyaml"
	"github.com/steadybit/cli/internal/output"
)

// Output writes to the file when one is given and to stdout otherwise, as JSON
// indented by two, or YAML.
func Output(doc *output.Document, file, explicitType string) error {
	// --jq filters what would be printed; a file is written as asked.
	if output.JQ != "" && file == "" {
		return output.ApplyJQ(os.Stdout, jsyaml.CompactJSON(doc.Value()), output.JQ)
	}
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
		if existingID == nil || existingID == "" {
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
