// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package experiment

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/steadybit/cli/api"
	"github.com/steadybit/cli/internal/jsyaml"
	"github.com/steadybit/cli/internal/output"
	"github.com/steadybit/cli/internal/platform"
)

type TemplateOptions struct {
	Template          string
	Team              string
	Environment       string
	ExternalID        string
	Placeholder       *jsyaml.Map
	PlaceholdersFile  string
	Variable          *jsyaml.Map
	ResetProperties   bool
	ExecutionVariable *jsyaml.Map
}

// ResolvePlaceholders reads the placeholders file, a map of key to value or the
// platform's list of {key, value}, and applies -p values on top, so that a pipeline can
// keep shared values in a file and override one per stage.
func ResolvePlaceholders(o TemplateOptions) ([]api.ExperimentTemplatePlaceholderValueAO, error) {
	values := jsyaml.NewMap()
	if o.PlaceholdersFile != "" {
		content, err := readAny(o.PlaceholdersFile)
		if err != nil {
			return nil, err
		}
		invalid := fmt.Errorf("Placeholders file '%s' must be a map of key to value or a list of {key, value} entries.", o.PlaceholdersFile)
		switch v := content.(type) {
		case []any:
			for _, entry := range v {
				m, ok := entry.(*jsyaml.Map)
				key, hasKey := m.Get("key")
				value, hasValue := m.Get("value")
				if !ok || !hasKey || !hasValue {
					return nil, invalid
				}
				keyString, ok := key.(string)
				if !ok {
					return nil, invalid
				}
				values.Set(keyString, value)
			}
		case *jsyaml.Map:
			for _, k := range v.Keys() {
				value, _ := v.Get(k)
				values.Set(k, value)
			}
		default:
			return nil, invalid
		}
	}
	if o.Placeholder != nil {
		for _, k := range o.Placeholder.Keys() {
			value, _ := o.Placeholder.Get(k)
			values.Set(k, value)
		}
	}
	result := make([]api.ExperimentTemplatePlaceholderValueAO, 0, values.Len())
	for _, k := range values.Keys() {
		value, _ := values.Get(k)
		result = append(result, api.ExperimentTemplatePlaceholderValueAO{Key: k, Value: plain(value)})
	}
	return result, nil
}

// plain turns a document value into something encoding/json writes the same way.
func plain(value any) any {
	switch v := value.(type) {
	case *jsyaml.Map:
		return rawJSON(jsyaml.CompactJSON(v))
	case []any:
		return rawJSON(jsyaml.CompactJSON(v))
	case float64:
		return rawJSON(jsyaml.CompactJSON(v))
	case jsyaml.Timestamp:
		return v.ISO()
	}
	return value
}

type rawJSON string

func (r rawJSON) MarshalJSON() ([]byte, error) { return []byte(r), nil }

func readAny(file string) (any, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("Failed to read placeholders file at path '%s': %s", file, pathCause(err))
	}
	value, err := output.ParseValue(content)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse placeholders file at path '%s' as YAML/JSON: %s", file, err)
	}
	return value, nil
}

func pathCause(err error) string {
	var pathErr *fs.PathError
	if errors.As(err, &pathErr) {
		return pathErr.Err.Error()
	}
	return err.Error()
}

// variables turns KEY=VALUE flags into the constant-string form of a variable.
func variables(values *jsyaml.Map) *map[string]api.VariableExpressionAO {
	if values == nil || values.Len() == 0 {
		return nil
	}
	result := map[string]api.VariableExpressionAO{}
	for _, k := range values.Keys() {
		value, _ := values.Get(k)
		var v api.VariableExpressionAO
		_ = v.FromVariableExpressionAO0(fmt.Sprint(value))
		result[k] = v
	}
	return &result
}

func templateID(id string) (openapi_types.UUID, error) {
	var uuid openapi_types.UUID
	if err := uuid.UnmarshalText([]byte(id)); err != nil {
		return uuid, fmt.Errorf("Experiment template %s not found.", id)
	}
	return uuid, nil
}

func createRequest(o TemplateOptions) (api.CreateExperimentFromTemplateAO, error) {
	if o.Team == "" {
		return api.CreateExperimentFromTemplateAO{}, errors.New("--team is required to create an experiment from a template.")
	}
	placeholders, err := ResolvePlaceholders(o)
	if err != nil {
		return api.CreateExperimentFromTemplateAO{}, err
	}
	request := api.CreateExperimentFromTemplateAO{Team: o.Team, Placeholders: &placeholders, ExperimentVariables: variables(o.Variable)}
	if o.Environment != "" {
		request.Environment = &o.Environment
	}
	if o.ExternalID != "" {
		request.ExternalId = &o.ExternalID
	}
	return request, nil
}

// ApplyTemplate creates an experiment from a template, or updates the one with key.
func ApplyTemplate(ctx context.Context, c *platform.Client, key string, o TemplateOptions) error {
	id, err := templateID(o.Template)
	if err != nil {
		return err
	}
	if key != "" {
		// An update only re-renders the experiment; it keeps its team and environment,
		// so accepting those here would silently do nothing.
		var ignored []string
		for _, f := range []struct {
			flag string
			set  bool
		}{{"--team", o.Team != ""}, {"--environment", o.Environment != ""}, {"--external-id", o.ExternalID != ""}, {"--variable", o.Variable != nil && o.Variable.Len() > 0}} {
			if f.set {
				ignored = append(ignored, f.flag)
			}
		}
		if len(ignored) > 0 {
			return fmt.Errorf("Updating experiment %s from a template only takes placeholders; remove %s.", key, strings.Join(ignored, ", "))
		}
		placeholders, err := ResolvePlaceholders(o)
		if err != nil {
			return err
		}
		_, _, err = platform.Read(c.UpdateExperimentByTemplate(ctx, id, key,
			&api.UpdateExperimentByTemplateParams{ResetProperties: &o.ResetProperties},
			api.UpdateExperimentFromTemplateAO{Placeholders: &placeholders}))
		if platform.IsStatus(err, http.StatusNotFound) {
			return fmt.Errorf("Experiment template %s or experiment %s not found.", o.Template, key)
		}
		if err != nil {
			return platform.Failed(err, "Failed to update experiment %s from template %s", key, o.Template)
		}
		fmt.Printf("Experiment %s updated from template %s.\n", key, o.Template)
		return nil
	}

	request, err := createRequest(o)
	if err != nil {
		return err
	}
	_, resp, err := platform.Read(c.CreateExperimentByTemplate(ctx, id,
		&api.CreateExperimentByTemplateParams{ResetProperties: &o.ResetProperties}, request))
	if platform.IsStatus(err, http.StatusNotFound) {
		return fmt.Errorf("Experiment template %s not found.", o.Template)
	}
	if err != nil {
		return platform.Failed(err, "Failed to create the experiment from template %s", o.Template)
	}
	verb := "updated"
	if resp.StatusCode == http.StatusCreated {
		verb = "created"
	}
	fmt.Printf("Experiment %s %s from template %s.\n", keyFromLocation(resp), verb, o.Template)
	return nil
}
