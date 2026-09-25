// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package execution

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/steadybit/cli/internal/jsyaml"
	"github.com/steadybit/cli/internal/output"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/table"
	"golang.org/x/term"
)

type WatchOptions struct {
	ID       int64
	Key      string
	Interval time.Duration
}

var ended = map[string]bool{"FAILED": true, "ERRORED": true, "CANCELED": true, "COMPLETED": true}

// Latest is the most recent run of an experiment.
func Latest(ctx context.Context, c *platform.Client, key string) (int64, error) {
	var list struct {
		Executions []struct {
			ID int64 `json:"id"`
		} `json:"executions"`
	}
	resp, err := c.GetExperimentExecutions3(ctx, key, nil)
	if _, err := platform.Decode(resp, err, &list); err != nil {
		return 0, platform.Failed(err, "Failed to get the runs of experiment %s", key)
	}
	var latest int64
	for _, e := range list.Executions {
		latest = max(latest, e.ID)
	}
	if latest == 0 {
		return 0, fmt.Errorf("Experiment %s has not run yet.", key)
	}
	return latest, nil
}

// Watch shows a run as it progresses, until it ends: redrawn in place on a terminal,
// as a line per change otherwise. It only watches; stopping it leaves the run alone.
func Watch(ctx context.Context, c *platform.Client, o WatchOptions) error {
	id := o.ID
	if o.Key != "" {
		latest, err := Latest(ctx, c, o.Key)
		if err != nil {
			return err
		}
		id = latest
	}
	interval := o.Interval
	if interval <= 0 {
		interval = 2 * time.Second
	}
	live := term.IsTerminal(int(os.Stdout.Fd()))
	drawn := 0
	previous := map[string]string{}
	for {
		doc, err := Fetch(ctx, c, id)
		if err != nil {
			return err
		}
		run := doc.Value()
		state := str(run, "state")
		if live {
			frame := render(run, time.Now())
			// Back to the top of the previous frame, and clear it, before drawing.
			if drawn > 0 {
				fmt.Printf("\x1b[%dA\x1b[J", drawn)
			}
			fmt.Print(frame)
			drawn = strings.Count(frame, "\n")
		} else {
			changes(os.Stdout, run, previous)
		}
		if ended[state] {
			if state != "COMPLETED" {
				reason := str(run, "reason")
				if reason != "" {
					reason = ", reason: " + reason
				}
				return fmt.Errorf("Experiment %s (#%d) %s%s", str(run, "key"), id, strings.ToLower(state), reason)
			}
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(interval):
		}
	}
}

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339Nano, s)
	return t
}

func elapsed(started, ended string, now time.Time) string {
	from := parseTime(started)
	if from.IsZero() {
		return ""
	}
	to := parseTime(ended)
	if to.IsZero() {
		to = now
	}
	return to.Sub(from).Round(time.Second).String()
}

func stepName(step *jsyaml.Map) string {
	if label := str(step, "customLabel"); label != "" {
		return label
	}
	if action := str(step, "actionId"); action != "" {
		return action
	}
	if strings.EqualFold(str(step, "stepType"), "wait") {
		if params, ok := step.Get("parameters"); ok {
			if m, ok := params.(*jsyaml.Map); ok {
				return "wait " + str(m, "duration")
			}
		}
	}
	return strings.ToLower(str(step, "stepType"))
}

// targets summarises a step's target executions: how many ended, of how many.
func targets(step *jsyaml.Map) string {
	all := list(step, "targetExecutions")
	if len(all) == 0 {
		return ""
	}
	done := 0
	for _, t := range all {
		if m, ok := t.(*jsyaml.Map); ok && ended[str(m, "state")] {
			done++
		}
	}
	return fmt.Sprintf("%d/%d", done, len(all))
}

func colorState(state string) string {
	lower := strings.ToLower(state)
	switch state {
	case "COMPLETED":
		return output.Green(lower)
	case "FAILED", "ERRORED":
		return output.Red(lower)
	case "RUNNING":
		return output.Bold(lower)
	}
	return lower
}

func render(run *jsyaml.Map, now time.Time) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s %s · run #%s · %s · %s\n", output.Bold(str(run, "key")), str(run, "name"),
		number(run, "id"), colorState(str(run, "state")), elapsed(str(run, "started"), str(run, "ended"), now))
	t := table.New(
		table.Column{Name: "n", Title: "#"},
		table.Column{Name: "step", Title: "Step", Alignment: table.Left},
		table.Column{Name: "state", Title: "State", Alignment: table.Left},
		table.Column{Name: "targets", Title: "Targets"},
		table.Column{Name: "time", Title: "Time"},
	)
	for i, s := range list(run, "steps") {
		step, _ := s.(*jsyaml.Map)
		t.AddRow(table.Default, table.Cell("n", i+1), table.Cell("step", stepName(step)), table.Cell("state", colorState(str(step, "state"))),
			table.Cell("targets", targets(step)), table.Cell("time", elapsed(str(step, "started"), str(step, "ended"), now)))
	}
	b.WriteString(t.Render())
	b.WriteString("\n")
	return b.String()
}

func number(m *jsyaml.Map, key string) string {
	v, _ := m.Get(key)
	if f, ok := v.(float64); ok {
		return jsyaml.NumberString(f)
	}
	return fmt.Sprint(v)
}

// changes prints what moved since the last poll, one line each, for logs and pipes.
func changes(w io.Writer, run *jsyaml.Map, previous map[string]string) {
	if state := str(run, "state"); previous["run"] != state {
		previous["run"] = state
		fmt.Fprintf(w, "Experiment %s run #%s: %s\n", str(run, "key"), number(run, "id"), strings.ToLower(state))
	}
	steps := list(run, "steps")
	for i, s := range steps {
		step, _ := s.(*jsyaml.Map)
		key := fmt.Sprint("step", i)
		now := str(step, "state") + " " + targets(step)
		if previous[key] != now {
			previous[key] = now
			line := fmt.Sprintf("  step %d/%d %s: %s", i+1, len(steps), stepName(step), strings.ToLower(str(step, "state")))
			if t := targets(step); t != "" {
				line += fmt.Sprintf(" (targets %s)", t)
			}
			fmt.Fprintln(w, line)
		}
	}
}
