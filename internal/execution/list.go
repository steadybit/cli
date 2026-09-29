// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package execution

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/steadybit/cli/v6/api"
	"github.com/steadybit/cli/v6/internal/platform"
	"github.com/steadybit/cli/v6/internal/resource"
	"github.com/steadybit/cli/v6/internal/table"
)

// States are the ones the platform filters by; REQUESTED is a run not yet created, which
// the search does not know.
var States = []string{"CREATED", "PREPARED", "RUNNING", "FAILED", "CANCELED", "COMPLETED", "ERRORED"}

type ListOptions struct {
	Name                                               string
	Teams, Experiments, Environments, Services, States []string
	From, To                                           string
	Limit                                              int
	FailOnMatch                                        bool
	Type                                               string
}

type summary struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	ExperimentKey    string `json:"experimentKey"`
	State            string `json:"state"`
	Started          string `json:"started"`
	Ended            string `json:"ended"`
	Scheduled        bool   `json:"scheduled"`
	CreatedBy        string `json:"createdBy"`
	CreatedByDetails *struct {
		Name     string `json:"name"`
		Username string `json:"username"`
	} `json:"createdByDetails"`
}

func (o ListOptions) request() (api.ExperimentExecutionsRequestAO, error) {
	from, err := resource.Time("from", o.From)
	if err != nil {
		return api.ExperimentExecutionsRequestAO{}, err
	}
	to, err := resource.Time("to", o.To)
	if err != nil {
		return api.ExperimentExecutionsRequestAO{}, err
	}
	// A date given to --to means the whole day: 00:00 would leave it out. The platform
	// includes a run created at createdTo, so the next day's 00:00 would take its first
	// runs too; the last instant of the day is exact.
	if _, err := time.Parse(time.DateOnly, o.To); err == nil {
		end := to.AddDate(0, 0, 1).Add(-time.Nanosecond)
		to = &end
	}
	var states []string
	for _, s := range o.States {
		state := strings.ToUpper(s)
		if !slices.Contains(States, state) {
			return api.ExperimentExecutionsRequestAO{}, fmt.Errorf("--state must be one of %s, not '%s'.", strings.Join(States, ", "), s)
		}
		states = append(states, state)
	}
	r := api.ExperimentExecutionsRequestAO{
		TeamKeys: resource.Optional(o.Teams), ExperimentKeys: resource.Optional(o.Experiments),
		Environments: resource.Optional(o.Environments), Services: resource.Optional(o.Services),
		States: resource.Optional(states), CreatedFrom: from, CreatedTo: to,
	}
	if o.Name != "" {
		r.Name = &o.Name
	}
	return r, nil
}

// search pages through the runs, newest first, until limit of them (0 for all). The
// page size stays the same throughout, as the platform counts pages in it. The request
// takes no sort: the platform sends the most recently requested run first.
func search(ctx context.Context, c *platform.Client, r api.ExperimentExecutionsRequestAO, limit int) ([]json.RawMessage, int64, error) {
	size := platform.PageSize
	if limit > 0 && limit < int(size) {
		size = int32(limit)
	}
	var items []json.RawMessage
	var total int64
	// A run created during the walk pushes the others one place down, so the next page
	// starts with the last of this one again.
	seen := map[int64]bool{}
	page := int32(0)
	for {
		r.Page, r.Size = &page, &size
		var body struct {
			Items      []json.RawMessage `json:"items"`
			TotalItems int64             `json:"totalItems"`
			NextPage   *int32            `json:"nextPage"`
		}
		resp, err := c.GetExperimentExecutions2(ctx, r)
		if _, err := platform.Decode(resp, err, &body); err != nil {
			return nil, 0, err
		}
		for _, item := range body.Items {
			var run struct {
				ID int64 `json:"id"`
			}
			if err := json.Unmarshal(item, &run); err != nil {
				return nil, 0, err
			}
			if !seen[run.ID] {
				seen[run.ID] = true
				items = append(items, item)
			}
		}
		total = body.TotalItems
		if limit > 0 && len(items) >= limit {
			return items[:limit], total, nil
		}
		// An empty page, or a next page that does not move on, ends it too, rather than
		// asking for the same pages forever.
		if body.NextPage == nil || *body.NextPage <= page || len(body.Items) == 0 {
			return items, total, nil
		}
		page = *body.NextPage
	}
}

// List searches the runs of all experiments. With FailOnMatch it lets a pipeline stop
// when any run matches, such as a failed one of its team since the last deployment.
func List(ctx context.Context, c *platform.Client, o ListOptions) error {
	r, err := o.request()
	if err != nil {
		return err
	}
	if o.Limit < 0 {
		return fmt.Errorf("--limit must be 0 or more, not %d.", o.Limit)
	}
	raw, total, err := search(ctx, c, r, o.Limit)
	if err != nil {
		return platform.Failed(err, "Failed to get the experiment runs")
	}
	if err := resource.List(raw, o.Type, func() error { return printRuns(raw, total) }); err != nil {
		return err
	}
	// JSON and YAML go to a script, which the note must not break, yet it must learn
	// that the limit cut the list.
	if resource.Machine(o.Type) && total > int64(len(raw)) {
		fmt.Fprintln(os.Stderr, cutNote(len(raw), total))
	}
	if o.FailOnMatch && len(raw) > 0 {
		if n := max(total, int64(len(raw))); n > 1 {
			return fmt.Errorf("%d experiment runs match.", n)
		}
		return errors.New("1 experiment run matches.")
	}
	return nil
}

func printRuns(raw []json.RawMessage, total int64) error {
	var runs []summary
	if err := resource.DecodeEach(raw, &runs); err != nil {
		return err
	}
	if len(runs) == 0 {
		fmt.Println("No experiment runs found.")
		return nil
	}
	t := table.New(
		table.Column{Name: "id", Title: "Id"},
		table.Column{Name: "experiment", Title: "Experiment", Alignment: table.Left},
		table.Column{Name: "name", Title: "Name", Alignment: table.Left},
		table.Column{Name: "state", Title: "State", Alignment: table.Left},
		table.Column{Name: "started", Title: "Started", Alignment: table.Left},
		table.Column{Name: "duration", Title: "Duration"},
		table.Column{Name: "trigger", Title: "Trigger", Alignment: table.Left},
	)
	now := time.Now()
	for _, r := range runs {
		t.AddRow(table.Default, table.Cell("id", r.ID), table.Cell("experiment", r.ExperimentKey), table.Cell("name", r.Name),
			table.Cell("state", colorState(r.State)), table.Cell("started", timestamp(r.Started)),
			table.Cell("duration", elapsed(r.Started, r.Ended, now)), table.Cell("trigger", r.trigger()))
	}
	t.Print()
	if total > int64(len(runs)) {
		fmt.Println(cutNote(len(runs), total))
	}
	return nil
}

func cutNote(shown int, total int64) string {
	return fmt.Sprintf("Showing the %d most recent of %d experiment runs. Raise --limit, or 0 for all.", shown, total)
}

// trigger names who started a run: a schedule, or the user or access token.
func (s summary) trigger() string {
	if s.Scheduled {
		return "schedule"
	}
	if d := s.CreatedByDetails; d != nil && d.Name != "" {
		return d.Name
	}
	return s.CreatedBy
}

// timestamp drops the platform's fractions of a second, which only widen the table.
func timestamp(s string) string {
	t := parseTime(s)
	if t.IsZero() {
		return s
	}
	return t.UTC().Format(time.RFC3339)
}
