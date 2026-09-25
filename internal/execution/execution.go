// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package execution implements the `execution` commands on experiment runs.
package execution

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/steadybit/cli/api"
	"github.com/steadybit/cli/internal/jsyaml"
	"github.com/steadybit/cli/internal/output"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/resource"
	"github.com/steadybit/cli/internal/table"
)

func notFoundOr(err error, id int64, format string, args ...any) error {
	if platform.IsStatus(err, http.StatusNotFound) {
		return fmt.Errorf("Experiment run %d not found.", id)
	}
	return platform.Failed(err, format, append(args, id)...)
}

// Fetch gets a run with its steps. The platform leaves the steps out unless asked, and
// with them every target execution and artifact.
func Fetch(ctx context.Context, c *platform.Client, id int64) (*output.Document, error) {
	fields := "steps"
	doc, _, err := platform.ReadDocument(c.GetExperimentExecution(ctx, id, &api.GetExperimentExecutionParams{Fields: &fields}))
	if err != nil {
		return nil, notFoundOr(err, id, "Failed to get experiment run %d")
	}
	return doc, nil
}

type GetOptions struct {
	ID         int64
	File, Type string
}

func Get(ctx context.Context, c *platform.Client, o GetOptions) error {
	doc, err := Fetch(ctx, c, o.ID)
	if err != nil {
		return err
	}
	if err := resource.Output(doc, o.File, o.Type); err != nil {
		return err
	}
	if o.File != "" {
		fmt.Printf("Experiment run %d written to %s.\n", o.ID, o.File)
	}
	return nil
}

// Cancel asks the platform to stop a run. A 202 means it is being stopped, a 200 that
// there was nothing left to cancel.
func Cancel(ctx context.Context, c *platform.Client, id int64) error {
	_, resp, err := platform.Read(c.CancelExperimentExecution(ctx, id))
	if err != nil {
		return notFoundOr(err, id, "Failed to cancel experiment run %d")
	}
	if resp.StatusCode == http.StatusAccepted {
		fmt.Printf("Experiment run %d is being canceled.\n", id)
	} else {
		fmt.Printf("Experiment run %d has already ended.\n", id)
	}
	return nil
}

type PropertyOptions struct {
	ID     int64
	Key    string
	Values []string
	JSON   bool
}

// Values are sent as strings unless --json asks otherwise. Guessing from the text would
// turn a ticket number such as "0042" into the number 42.
func parseValue(value string, asJSON bool) (json.RawMessage, error) {
	if !asJSON {
		return json.RawMessage(jsyaml.CompactJSON(value)), nil
	}
	parsed, err := output.ParseValue([]byte(value))
	if err != nil || !json.Valid([]byte(value)) {
		return nil, fmt.Errorf("'%s' is not valid JSON: %s", value, jsonError(value))
	}
	return json.RawMessage(jsyaml.CompactJSON(parsed)), nil
}

func jsonError(value string) string {
	var v any
	if err := json.Unmarshal([]byte(value), &v); err != nil {
		return err.Error()
	}
	return "invalid JSON"
}

func SetProperty(ctx context.Context, c *platform.Client, o PropertyOptions) error {
	values := make([]json.RawMessage, len(o.Values))
	for i, v := range o.Values {
		parsed, err := parseValue(v, o.JSON)
		if err != nil {
			return err
		}
		values[i] = parsed
	}
	// Several values set a list property; a single one stays a scalar.
	var body json.RawMessage
	if len(values) == 1 {
		body = values[0]
	} else {
		body, _ = json.Marshal(values)
	}
	_, _, err := platform.Read(c.SetExecutionPropertyValueWithBody(ctx, o.ID, o.Key, "application/json", bytes.NewReader(body)))
	return reportProperty(err, "set", o)
}

func AddProperty(ctx context.Context, c *platform.Client, o PropertyOptions) error {
	if len(o.Values) != 1 {
		return fmt.Errorf("Adding to a list property takes exactly one --value.")
	}
	body, err := parseValue(o.Values[0], o.JSON)
	if err != nil {
		return err
	}
	_, _, err = platform.Read(c.AddExecutionPropertyValueWithBody(ctx, o.ID, o.Key, "application/json", bytes.NewReader(body)))
	return reportProperty(err, "add", o)
}

func reportProperty(err error, operation string, o PropertyOptions) error {
	if err != nil {
		return notFoundOr(err, o.ID, fmt.Sprintf("Failed to %s property %s of experiment run ", operation, o.Key)+"%d")
	}
	fmt.Printf("Property %s of experiment run %d updated.\n", o.Key, o.ID)
	return nil
}

type Artifact struct {
	Step, Target, TargetExecutionID, ArtifactID string
}

func str(m *jsyaml.Map, key string) string {
	if m == nil {
		return ""
	}
	v, _ := m.Get(key)
	s, _ := v.(string)
	return s
}

func list(m *jsyaml.Map, key string) []any {
	if m == nil {
		return nil
	}
	v, _ := m.Get(key)
	l, _ := v.([]any)
	return l
}

// Collect gathers artifacts from the target executions of action steps, and of the
// actions a service validation step runs; the platform offers no listing of its own.
func Collect(run *jsyaml.Map) []Artifact {
	var artifacts []Artifact
	addFrom := func(step string, targets []any) {
		for _, t := range targets {
			target, _ := t.(*jsyaml.Map)
			for _, a := range list(target, "artifacts") {
				if id, ok := a.(string); ok {
					artifacts = append(artifacts, Artifact{Step: step, Target: str(target, "name"), TargetExecutionID: str(target, "id"), ArtifactID: id})
				}
			}
		}
	}
	firstOf := func(values ...string) string {
		for _, v := range values {
			if v != "" {
				return v
			}
		}
		return ""
	}
	for _, s := range list(run, "steps") {
		step, _ := s.(*jsyaml.Map)
		label := firstOf(str(step, "customLabel"), str(step, "actionId"), str(step, "stepType"))
		addFrom(label, list(step, "targetExecutions"))
		for _, v := range list(step, "validations") {
			validation, _ := v.(*jsyaml.Map)
			addFrom(firstOf(str(validation, "customLabel"), str(validation, "actionId"), label), list(validation, "targetExecutions"))
		}
	}
	return artifacts
}

func ListArtifacts(ctx context.Context, c *platform.Client, id int64) error {
	doc, err := Fetch(ctx, c, id)
	if err != nil {
		return err
	}
	artifacts := Collect(doc.Value())
	if len(artifacts) == 0 {
		fmt.Printf("Experiment run %d has no artifacts.\n", id)
		return nil
	}
	t := table.New(
		table.Column{Name: "artifactId", Title: "Artifact", Alignment: table.Left},
		table.Column{Name: "target", Title: "Target", Alignment: table.Left},
		table.Column{Name: "step", Title: "Step", Alignment: table.Left},
		table.Column{Name: "targetExecutionId", Title: "Target execution", Alignment: table.Left},
	)
	for _, a := range artifacts {
		t.AddRow(table.Default, table.Cell("artifactId", a.ArtifactID), table.Cell("target", a.Target), table.Cell("step", a.Step), table.Cell("targetExecutionId", a.TargetExecutionID))
	}
	t.Print()
	return nil
}

type DownloadOptions struct {
	ID                                   int64
	Artifact, TargetExecution, Directory string
	Output                               string
}

// PathSegment reduces an id from the platform to one path segment, so that one
// containing "../" cannot write outside the chosen directory; Base alone leaves "..".
func PathSegment(id string) string {
	segment := filepath.Base(id)
	if segment == "" || segment == "." || segment == ".." || segment == string(filepath.Separator) {
		return "_"
	}
	return segment
}

// Download writes every selected artifact to <directory>/<target execution>/<artifact>:
// two targets of one step usually produce files of the same name.
func Download(ctx context.Context, c *platform.Client, o DownloadOptions) error {
	doc, err := Fetch(ctx, c, o.ID)
	if err != nil {
		return err
	}
	var selected []Artifact
	for _, a := range Collect(doc.Value()) {
		if (o.Artifact == "" || a.ArtifactID == o.Artifact) && (o.TargetExecution == "" || a.TargetExecutionID == o.TargetExecution) {
			selected = append(selected, a)
		}
	}
	if len(selected) == 0 {
		return fmt.Errorf("No matching artifacts found in experiment run %d.", o.ID)
	}
	if o.Output != "" && len(selected) > 1 {
		return fmt.Errorf("%d artifacts match, but --output takes exactly one. Narrow it down with --artifact and --target-execution.", len(selected))
	}
	// Artifacts are reports and log archives, which take far longer than an API response.
	ctx = platform.WithTimeout(ctx, 5*time.Minute)
	for _, a := range selected {
		file := o.Output
		if file == "" {
			file = filepath.Join(o.Directory, PathSegment(a.TargetExecutionID), PathSegment(a.ArtifactID))
		}
		content, _, err := platform.Read(c.GetArtifact(ctx, o.ID, a.TargetExecutionID, a.ArtifactID))
		if platform.IsStatus(err, http.StatusNotFound) {
			return fmt.Errorf("Artifact %s of experiment run %d not found.", a.ArtifactID, o.ID)
		}
		if err != nil {
			return platform.Failed(err, "Failed to download artifact %s of experiment run %d", a.ArtifactID, o.ID)
		}
		if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(file, content, 0o644); err != nil {
			return err
		}
		fmt.Printf("Artifact %s written to %s.\n", a.ArtifactID, file)
	}
	return nil
}
