// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package service implements the `service` commands.
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/steadybit/cli/api"
	"github.com/steadybit/cli/internal/experiment"
	"github.com/steadybit/cli/internal/jsyaml"
	"github.com/steadybit/cli/internal/output"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/resource"
	"github.com/steadybit/cli/internal/table"
)

// Who created and edited a service cannot be sent back. The version is dropped as
// `experiment get` drops it: kept in a file, it turns every apply after an edit in the UI
// into a conflict.
var readOnly = []string{"created", "createdBy", "edited", "editedBy", "version"}

var ErrNotFound = errors.New("not found")

func uuid(id string) (openapi_types.UUID, error) {
	var u openapi_types.UUID
	if err := u.UnmarshalText([]byte(id)); err != nil {
		return u, fmt.Errorf("Service %s not found.", id)
	}
	return u, nil
}

func notFoundOr(err error, id, format string, args ...any) error {
	if platform.IsStatus(err, http.StatusNotFound) {
		return fmt.Errorf("Service %s not found.", id)
	}
	return platform.Failed(err, format, append(args, id)...)
}

func optional(values []string) *[]string {
	if len(values) == 0 {
		return nil
	}
	return &values
}

type ListOptions struct {
	Teams, Environments, Experiments []string
	Type                             string
}

func List(ctx context.Context, c *platform.Client, o ListOptions) error {
	type summary struct {
		ID, Name, Team, Environment string
	}
	raw, err := platform.AllPagesRaw(func(page, size int32) (*http.Response, error) {
		return c.GetServiceList(ctx, &api.GetServiceListParams{
			TeamKey: optional(o.Teams), EnvironmentName: optional(o.Environments), ExperimentKey: optional(o.Experiments),
			Page: api.PageRequestAO{Page: &page, Size: &size},
		})
	})
	if err != nil {
		return platform.Failed(err, "Failed to get the services")
	}
	if resource.Machine(o.Type) {
		return resource.List(raw, o.Type, nil)
	}
	var services []summary
	if err := resource.DecodeEach(raw, &services); err != nil {
		return err
	}
	if len(services) == 0 {
		fmt.Println("No services found.")
		return nil
	}
	t := table.New(
		table.Column{Name: "id", Title: "Id", Alignment: table.Left},
		table.Column{Name: "name", Title: "Name", Alignment: table.Left},
		table.Column{Name: "team", Title: "Team", Alignment: table.Left},
		table.Column{Name: "environment", Title: "Environment", Alignment: table.Left},
	)
	for _, s := range services {
		t.AddRow(table.Default, table.Cell("id", s.ID), table.Cell("name", s.Name), table.Cell("team", s.Team), table.Cell("environment", s.Environment))
	}
	t.Print()
	return nil
}

type GetOptions struct {
	ID, File, Type string
}

func Get(ctx context.Context, c *platform.Client, o GetOptions) error {
	id, err := uuid(o.ID)
	if err != nil {
		return err
	}
	doc, _, err := platform.ReadDocument(c.GetService(ctx, id))
	if err != nil {
		return notFoundOr(err, o.ID, "Failed to get service %s")
	}
	if err := resource.Output(resource.Strip(doc, readOnly...), o.File, o.Type); err != nil {
		return err
	}
	if o.File != "" {
		fmt.Printf("Service %s written to %s.\n", o.ID, o.File)
	}
	return nil
}

type ApplyOptions struct {
	Files             []string
	Recursive         bool
	DeleteExperiments bool
}

// RefusedForProvidedExperiments recognises the platform refusing a change that would
// orphan provided experiments. It names its query parameter, so the user is pointed at
// the CLI flag instead.
func RefusedForProvidedExperiments(err error) bool {
	var apiErr *platform.APIError
	if !errors.As(err, &apiErr) || apiErr.Status != http.StatusUnprocessableEntity {
		return false
	}
	var problem struct {
		Violations []struct {
			Message string `json:"message"`
		} `json:"violations"`
	}
	_ = json.Unmarshal(apiErr.Body, &problem)
	for _, v := range problem.Violations {
		if strings.Contains(v.Message, "deleteExperiments") {
			return true
		}
	}
	return false
}

func Apply(ctx context.Context, c *platform.Client, o ApplyOptions) error {
	return resource.ApplyFiles(o.Files, o.Recursive, "service", func(file string, doc *output.Document) (resource.Applied, error) {
		name, _ := doc.Get("name")
		if name == "" {
			return resource.Applied{}, fmt.Errorf("Service file '%s' does not name the service.", file)
		}
		var saved struct{ ID, Name string }
		resp, err := c.UpsertServiceWithBody(ctx, &api.UpsertServiceParams{DeleteExperiments: &o.DeleteExperiments}, "application/json",
			bytes.NewReader([]byte(jsyaml.CompactJSON(resource.Strip(doc, readOnly...).Value()))))
		resp, err = platform.Decode(resp, err, &saved)
		if err != nil {
			if !o.DeleteExperiments && RefusedForProvidedExperiments(err) {
				return resource.Applied{}, fmt.Errorf("Service %s was not saved: the change would remove provided experiments. Pass --delete-experiments to delete them.", name)
			}
			return resource.Applied{}, platform.Failed(err, "Failed to save service %s", name)
		}
		created := resp.StatusCode == http.StatusCreated
		fmt.Printf("Service %s (%s) %s.\n", saved.Name, saved.ID, resource.CreatedOrUpdated(created))
		return resource.Applied{ID: saved.ID, Created: created}, nil
	})
}

func Delete(ctx context.Context, c *platform.Client, idText string) error {
	id, err := uuid(idText)
	if err != nil {
		return err
	}
	if _, _, err := platform.Read(c.DeleteService(ctx, id)); err != nil {
		return notFoundOr(err, idText, "Failed to delete service %s")
	}
	fmt.Printf("Service %s deleted.\n", idText)
	return nil
}

type RiskOptions struct {
	ID        string
	Type      string
	FailAbove *int
}

func Risk(ctx context.Context, c *platform.Client, o RiskOptions) error {
	id, err := uuid(o.ID)
	if err != nil {
		return err
	}
	doc, _, err := platform.ReadDocument(c.GetRisk(ctx, id))
	if platform.IsStatus(err, http.StatusNotFound) {
		return fmt.Errorf("Service %s not found, or its risk has not been calculated yet.", o.ID)
	}
	if err != nil {
		return platform.Failed(err, "Failed to get the risk of service %s", o.ID)
	}
	value := doc.Value()
	riskValue, hasRisk := value.Get("risk")
	risk, _ := riskValue.(float64)
	riskText := "unknown"
	if n, ok := riskValue.(float64); hasRisk && ok {
		riskText = jsyaml.NumberString(n)
	}

	if resource.Machine(o.Type) {
		if err := resource.Output(doc, "", o.Type); err != nil {
			return err
		}
	} else {
		calculated, _ := doc.Get("lastCalculated")
		if calculated == "" {
			calculated = "never"
		}
		fmt.Printf("Risk of service %s: %s (calculated %s)\n", o.ID, riskText, calculated)
		if categories, ok := value.Get("categoryRisks"); ok {
			if m, ok := categories.(*jsyaml.Map); ok && m.Len() > 0 {
				t := table.New(table.Column{Name: "category", Title: "Category", Alignment: table.Left},
					table.Column{Name: "total", Title: "Total"}, table.Column{Name: "experiment", Title: "Experiments"}, table.Column{Name: "advice", Title: "Advice"})
				for _, category := range m.Keys() {
					r, _ := m.Get(category)
					rm, _ := r.(*jsyaml.Map)
					t.AddRow(table.Default, table.Cell("category", category), table.Cell("total", number(rm, "total")),
						table.Cell("experiment", number(rm, "experiment")), table.Cell("advice", number(rm, "advice")))
				}
				t.Print()
			}
		}
		if experiments, ok := value.Get("experimentRisks"); ok {
			if list, ok := experiments.([]any); ok && len(list) > 0 {
				t := table.New(table.Column{Name: "experimentKey", Title: "Experiment", Alignment: table.Left}, table.Column{Name: "risk", Title: "Risk"})
				for _, e := range list {
					em, _ := e.(*jsyaml.Map)
					key, _ := em.Get("experimentKey")
					t.AddRow(table.Default, table.Cell("experimentKey", key), table.Cell("risk", number(em, "risk")))
				}
				t.Print()
			}
		}
	}

	// Lets a pipeline stop the rollout of a service whose risk is too high.
	if o.FailAbove != nil && (!hasRisk || riskText == "unknown" || risk > float64(*o.FailAbove)) {
		return fmt.Errorf("Risk of service %s is %s, above the accepted %d.", o.ID, riskText, *o.FailAbove)
	}
	return nil
}

func number(m *jsyaml.Map, key string) any {
	if m == nil {
		return nil
	}
	v, _ := m.Get(key)
	if n, ok := v.(float64); ok {
		return jsyaml.NumberString(n)
	}
	return v
}

type ExperimentListOptions struct {
	ID                string
	Categories, Types []string
	Type              string
}

func ListExperiments(ctx context.Context, c *platform.Client, o ExperimentListOptions) error {
	id, err := uuid(o.ID)
	if err != nil {
		return err
	}
	var types *[]api.GetServiceExperimentsParamsType
	if len(o.Types) > 0 {
		list := make([]api.GetServiceExperimentsParamsType, len(o.Types))
		for i, t := range o.Types {
			list[i] = api.GetServiceExperimentsParamsType(strings.ToUpper(t))
		}
		types = &list
	}
	type entry struct {
		ExperimentKey   *string `json:"experimentKey"`
		TemplateID      *string `json:"templateId"`
		Category        string  `json:"category"`
		AssociationType string  `json:"associationType"`
	}
	raw, err := platform.AllPagesRaw(func(page, size int32) (*http.Response, error) {
		return c.GetServiceExperiments(ctx, id, &api.GetServiceExperimentsParams{Category: optional(o.Categories), Type: types, Page: api.PageRequestAO{Page: &page, Size: &size}})
	})
	if err != nil {
		return notFoundOr(err, o.ID, "Failed to get the experiments of service %s")
	}
	if resource.Machine(o.Type) {
		return resource.List(raw, o.Type, nil)
	}
	var experiments []entry
	if err := resource.DecodeEach(raw, &experiments); err != nil {
		return err
	}
	if len(experiments) == 0 {
		if len(o.Categories) > 0 || len(o.Types) > 0 {
			fmt.Printf("Service %s has no matching experiments.\n", o.ID)
		} else {
			fmt.Printf("Service %s has no experiments.\n", o.ID)
		}
		return nil
	}
	t := table.New(
		table.Column{Name: "category", Title: "Category", Alignment: table.Left},
		table.Column{Name: "associationType", Title: "Type", Alignment: table.Left},
		// A provided experiment not created yet has no key, only its template.
		table.Column{Name: "experimentKey", Title: "Experiment", Alignment: table.Left},
		table.Column{Name: "templateId", Title: "Template", Alignment: table.Left},
	)
	for _, e := range experiments {
		key, template := "(not created)", ""
		if e.ExperimentKey != nil {
			key = *e.ExperimentKey
		}
		if e.TemplateID != nil {
			template = *e.TemplateID
		}
		t.AddRow(table.Default, table.Cell("category", e.Category), table.Cell("associationType", e.AssociationType), table.Cell("experimentKey", key), table.Cell("templateId", template))
	}
	t.Print()
	return nil
}

func Link(ctx context.Context, c *platform.Client, idText, experimentKey, category string) error {
	id, err := uuid(idText)
	if err != nil {
		return err
	}
	if _, _, err := platform.Read(c.LinkCustomExperiment(ctx, id, api.LinkCustomExperimentRequestAO{ExperimentKey: experimentKey, Category: category})); err != nil {
		return notFoundOr(err, idText, "Failed to link experiment %s to service %s", experimentKey)
	}
	fmt.Printf("Experiment %s linked to service %s in category %s.\n", experimentKey, idText, category)
	return nil
}

func Unlink(ctx context.Context, c *platform.Client, idText, experimentKey string) error {
	id, err := uuid(idText)
	if err != nil {
		return err
	}
	if _, _, err := platform.Read(c.UnlinkCustomExperiment(ctx, id, &api.UnlinkCustomExperimentParams{ExperimentKey: experimentKey})); err != nil {
		return notFoundOr(err, idText, "Failed to unlink experiment %s from service %s", experimentKey)
	}
	fmt.Printf("Experiment %s unlinked from service %s.\n", experimentKey, idText)
	return nil
}

type ProvideOptions struct {
	ID         string
	Experiment string
	experiment.TemplateOptions
}

func Provide(ctx context.Context, c *platform.Client, o ProvideOptions) error {
	id, err := uuid(o.ID)
	if err != nil {
		return err
	}
	placeholders, err := experiment.ResolvePlaceholders(o.TemplateOptions)
	if err != nil {
		return err
	}
	var templateID openapi_types.UUID
	if err := templateID.UnmarshalText([]byte(o.Template)); err != nil {
		return fmt.Errorf("Service %s or experiment template %s not found.", o.ID, o.Template)
	}
	request := api.UpsertProvidedExperimentRequestAO{TemplateId: templateID, Placeholders: &placeholders}
	if o.Experiment != "" {
		request.ExperimentKey = &o.Experiment
	}
	_, resp, err := platform.Read(c.UpsertProvidedExperiment(ctx, id, &api.UpsertProvidedExperimentParams{ResetProperties: &o.ResetProperties}, request))
	if platform.IsStatus(err, http.StatusNotFound) {
		return fmt.Errorf("Service %s or experiment template %s not found.", o.ID, o.Template)
	}
	if err != nil {
		return platform.Failed(err, "Failed to save the provided experiment of service %s", o.ID)
	}
	key := o.Experiment
	if location := resp.Header.Get("Location"); location != "" {
		key = location[strings.LastIndex(location, "/")+1:]
	}
	fmt.Printf("Provided experiment %s of service %s %s from template %s.\n", key, o.ID, resource.CreatedOrUpdated(resp.StatusCode == http.StatusCreated), o.Template)
	return nil
}

type VariableGetOptions struct {
	ID, Type string
}

func GetVariables(ctx context.Context, c *platform.Client, o VariableGetOptions) error {
	id, err := uuid(o.ID)
	if err != nil {
		return err
	}
	doc, _, err := platform.ReadDocument(c.GetServiceVariables(ctx, id))
	if err != nil {
		return notFoundOr(err, o.ID, "Failed to get the variables of service %s")
	}
	return resource.Output(doc, "", o.Type)
}

type VariableSetOptions struct {
	ID      string
	File    string
	Replace bool
}

// SetVariables merges KEY=VALUE arguments over a file's variables. With --replace, the
// result is all there is.
func SetVariables(ctx context.Context, c *platform.Client, pairs []string, o VariableSetOptions) error {
	variables, err := resource.Variables(pairs, o.File, o.Replace)
	if err != nil {
		return err
	}
	id, err := uuid(o.ID)
	if err != nil {
		return err
	}
	body := resource.Body(variables)
	if o.Replace {
		_, _, err = platform.Read(c.SetServiceVariablesWithBody(ctx, id, "application/json", body))
	} else {
		_, _, err = platform.Read(c.MergeServiceVariablesWithBody(ctx, id, "application/json", body))
	}
	if err != nil {
		return notFoundOr(err, o.ID, "Failed to update the variables of service %s")
	}
	fmt.Printf("%d variable(s) of service %s %s.\n", variables.Len(), o.ID, resource.VariablesOutcome(o.Replace))
	return nil
}
