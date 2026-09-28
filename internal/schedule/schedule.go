// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package schedule implements the `schedule` commands.
package schedule

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/steadybit/cli/v6/api"
	"github.com/steadybit/cli/v6/internal/jsyaml"
	"github.com/steadybit/cli/v6/internal/output"
	"github.com/steadybit/cli/v6/internal/platform"
	"github.com/steadybit/cli/v6/internal/resource"
	"github.com/steadybit/cli/v6/internal/table"
)

// What the platform reports about the last edit and the next run cannot be sent back,
// so it is left out of files, keeping `get` followed by `apply` a round trip.
var readOnly = []string{"editedBy", "lastUpdated", "nextExecution"}

func notFoundOr(err error, id, format string) error {
	if platform.IsStatus(err, http.StatusNotFound) {
		return fmt.Errorf("Experiment schedule %s not found.", id)
	}
	return platform.Failed(err, format, id)
}

func body(m *jsyaml.Map) io.Reader { return bytes.NewReader([]byte(jsyaml.CompactJSON(m))) }

func optional(values []string) *[]string {
	if len(values) == 0 {
		return nil
	}
	return &values
}

type ListOptions struct {
	Teams, Experiments []string
	Type               string
}

func List(ctx context.Context, c *platform.Client, o ListOptions) error {
	var schedules []struct {
		ID            string  `json:"id"`
		ExperimentKey string  `json:"experimentKey"`
		Cron          *string `json:"cron"`
		StartAt       *string `json:"startAt"`
		Timezone      *string `json:"timezone"`
		Enabled       *bool   `json:"enabled"`
		AllowParallel *bool   `json:"allowParallel"`
	}
	resp, err := c.GetAllSchedulesV2(ctx, &api.GetAllSchedulesV2Params{Team: optional(o.Teams), Experiment: optional(o.Experiments)})
	var raw []json.RawMessage
	if _, err := platform.Decode(resp, err, &raw); err != nil {
		return platform.Failed(err, "Failed to get the experiment schedules")
	}
	if resource.Machine(o.Type) {
		return resource.List(raw, o.Type, nil)
	}
	if err := resource.DecodeEach(raw, &schedules); err != nil {
		return err
	}
	if len(schedules) == 0 {
		fmt.Println("No experiment schedules found.")
		return nil
	}
	t := table.New(
		table.Column{Name: "id", Title: "Id", Alignment: table.Left},
		table.Column{Name: "experiment", Title: "Experiment", Alignment: table.Left},
		table.Column{Name: "when", Title: "When", Alignment: table.Left},
		table.Column{Name: "enabled", Title: "Enabled", Alignment: table.Left},
		table.Column{Name: "allowParallel", Title: "Parallel", Alignment: table.Left},
	)
	boolOr := func(b *bool) string {
		if b == nil {
			return "true"
		}
		return fmt.Sprint(*b)
	}
	for _, s := range schedules {
		when := ""
		switch {
		case s.Cron != nil && *s.Cron != "":
			when = *s.Cron
			if s.Timezone != nil && *s.Timezone != "" {
				when += " (" + *s.Timezone + ")"
			}
		case s.StartAt != nil:
			when = *s.StartAt
		}
		t.AddRow(table.Default, table.Cell("id", s.ID), table.Cell("experiment", s.ExperimentKey), table.Cell("when", when),
			table.Cell("enabled", boolOr(s.Enabled)), table.Cell("allowParallel", boolOr(s.AllowParallel)))
	}
	t.Print()
	return nil
}

type GetOptions struct {
	ID, File, Type string
}

func Get(ctx context.Context, c *platform.Client, o GetOptions) error {
	doc, _, err := platform.ReadDocument(c.GetSchedules(ctx, o.ID))
	if err != nil {
		return notFoundOr(err, o.ID, "Failed to get experiment schedule %s")
	}
	if err := resource.Output(resource.Strip(doc, readOnly...), o.File, o.Type); err != nil {
		return err
	}
	if o.File != "" {
		fmt.Printf("Experiment schedule %s written to %s.\n", o.ID, o.File)
	}
	return nil
}

type upserted struct {
	ID            string `json:"id"`
	ExperimentKey string `json:"experimentKey"`
}

type ApplyOptions struct {
	Files     []string
	Recursive bool
}

func Apply(ctx context.Context, c *platform.Client, o ApplyOptions) error {
	return resource.ApplyFiles(o.Files, o.Recursive, "schedule", func(file string, doc *output.Document) (resource.Applied, error) {
		if key, _ := doc.Get("experimentKey"); key == "" {
			return resource.Applied{}, fmt.Errorf("Schedule file '%s' does not name an experimentKey.", file)
		}
		schedule, created, err := save(ctx, c, resource.Strip(doc, readOnly...).Value())
		if err != nil {
			return resource.Applied{}, err
		}
		fmt.Printf("Experiment schedule %s for %s %s.\n", schedule.ID, schedule.ExperimentKey, resource.CreatedOrUpdated(created))
		return resource.Applied{ID: schedule.ID, Created: created}, nil
	})
}

func save(ctx context.Context, c *platform.Client, schedule *jsyaml.Map) (upserted, bool, error) {
	resp, err := c.UpsertScheduleWithBody(ctx, "application/json", body(schedule))
	var result upserted
	resp, err = platform.Decode(resp, err, &result)
	if err != nil {
		key, _ := schedule.Get("experimentKey")
		return result, false, platform.Failed(err, "Failed to save the experiment schedule for %v", key)
	}
	return result, resp.StatusCode == http.StatusCreated, nil
}

type Fields struct {
	Cron, StartAt, Timezone string
	AllowParallel           *bool
	Variables               *jsyaml.Map
}

func (f Fields) check(requireOne bool) error {
	if f.Cron != "" && f.StartAt != "" {
		return errors.New("--cron and --start-at cannot be combined.")
	}
	if requireOne && f.Cron == "" && f.StartAt == "" {
		return errors.New("Either --cron or --start-at must be specified.")
	}
	return nil
}

// into sets the fields that were given, leaving the others out of the request.
func (f Fields) into(m *jsyaml.Map) {
	for _, field := range [][2]string{{"cron", f.Cron}, {"startAt", f.StartAt}, {"timezone", f.Timezone}} {
		if field[1] != "" {
			m.Set(field[0], field[1])
		}
	}
	if f.AllowParallel != nil {
		m.Set("allowParallel", *f.AllowParallel)
	}
	if f.Variables != nil && f.Variables.Len() > 0 {
		m.Set("variables", f.Variables)
	}
}

type CreateOptions struct {
	Fields
	Experiment string
	Disabled   bool
}

func Create(ctx context.Context, c *platform.Client, o CreateOptions) error {
	if err := o.check(true); err != nil {
		return err
	}
	m := jsyaml.NewMap()
	m.Set("experimentKey", o.Experiment)
	o.into(m)
	m.Set("enabled", !o.Disabled)
	schedule, _, err := save(ctx, c, m)
	if err != nil {
		return err
	}
	fmt.Printf("Experiment schedule %s for %s created.\n", schedule.ID, schedule.ExperimentKey)
	return nil
}

type UpdateOptions struct {
	Fields
	ID string
}

func Update(ctx context.Context, c *platform.Client, o UpdateOptions) error {
	if err := o.check(false); err != nil {
		return err
	}
	m := jsyaml.NewMap()
	o.into(m)
	if m.Len() == 0 {
		return errors.New("Nothing to update. Pass at least one of the options, see --help.")
	}
	return patchAndReport(ctx, c, o.ID, m, "updated")
}

func SetEnabled(ctx context.Context, c *platform.Client, id string, enabled bool) error {
	m := jsyaml.NewMap()
	m.Set("enabled", enabled)
	outcome := "disabled"
	if enabled {
		outcome = "enabled"
	}
	return patchAndReport(ctx, c, id, m, outcome)
}

func patchAndReport(ctx context.Context, c *platform.Client, id string, m *jsyaml.Map, outcome string) error {
	resp, err := c.PatchScheduleWithBody(ctx, id, "application/json", body(m))
	var result upserted
	if _, err := platform.Decode(resp, err, &result); err != nil {
		return notFoundOr(err, id, "Failed to update experiment schedule %s")
	}
	fmt.Printf("Experiment schedule %s for %s %s.\n", id, result.ExperimentKey, outcome)
	return nil
}

func Delete(ctx context.Context, c *platform.Client, id string) error {
	if _, _, err := platform.Read(c.RemoveExperimentScheduleById(ctx, id)); err != nil {
		return notFoundOr(err, id, "Failed to delete experiment schedule %s")
	}
	fmt.Printf("Experiment schedule %s deleted.\n", id)
	return nil
}
