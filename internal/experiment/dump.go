// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package experiment

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/steadybit/cli/v6/api"
	"github.com/steadybit/cli/v6/internal/jsyaml"
	"github.com/steadybit/cli/v6/internal/output"
	"github.com/steadybit/cli/v6/internal/platform"
)

// A dump walks every experiment of every team and every execution of every experiment.
// Both levels are bounded so the request volume stays predictable instead of scaling
// with the size of the tenant.
const (
	experimentConcurrency = 4
	executionConcurrency  = 16
	largeDumpExperiments  = 100
)

type DumpOptions struct {
	Directory string
	Type      string
	Teams     []string
}

type team struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type listed struct {
	Key string `json:"key"`
}

// ErrIncomplete makes the command exit non-zero after everything that could be fetched
// has been written, so that a pipeline does not mistake a partial dump for a full one.
var ErrIncomplete = fmt.Errorf("incomplete dump")

func Dump(ctx context.Context, c *platform.Client, o DumpOptions) error {
	datatype := output.YAML
	if o.Type == "json" {
		datatype = output.JSON
	}
	if err := os.MkdirAll(o.Directory, 0o755); err != nil {
		return err
	}

	var teamList struct {
		Teams []team `json:"teams"`
	}
	onlyAccessible := false
	if err := getJSON(c.GetTeams(ctx, &api.GetTeamsParams{OnlyAccessible: &onlyAccessible}))(&teamList); err != nil {
		return fmt.Errorf("Failed to get teams: %w", err)
	}
	teams, err := selectTeams(teamList.Teams, o.Teams)
	if err != nil {
		return err
	}

	// The lists are fetched up front, which costs nothing extra because each team needs
	// one anyway, so the size of the walk is known before it starts.
	plural := "teams"
	if len(teams) == 1 {
		plural = "team"
	}
	fmt.Printf("Listing experiments for %d %s", len(teams), plural)
	lists := map[string][]listed{}
	total := 0
	for _, t := range teams {
		var list struct {
			Experiments []listed `json:"experiments"`
		}
		key := []string{t.Key}
		if err := getJSON(c.GetExperiments(ctx, &api.GetExperimentsParams{Team: &key}))(&list); err != nil {
			fmt.Println()
			return fmt.Errorf("Failed to get the experiments. HTTP request failed: %w", err)
		}
		lists[t.Key] = list.Experiments
		total += len(list.Experiments)
		fmt.Print(".")
	}
	fmt.Println()
	if total > largeDumpExperiments {
		minutes := int(math.Ceil(platform.Limiter().DurationFor(total * 2).Minutes()))
		fmt.Fprintf(os.Stderr, "Dumping %d experiments. Requests are paced to the platform's rate limit, so this takes at least %d minutes, longer with executions.\n\n", total, minutes)
	}

	var experiments, executions, failedExperiments, failedExecutions int
	for _, t := range teams {
		fmt.Printf("Fetching experiments for team %s (%s)... ", t.Name, t.Key)
		d := dumpTeam(ctx, c, lists[t.Key], o.Directory, datatype)
		experiments += d.experiments
		executions += d.executions
		failedExperiments += d.failedExperiments
		failedExecutions += d.failedExecutions
		failed := ""
		if n := d.failedExperiments + d.failedExecutions; n > 0 {
			failed = fmt.Sprintf(", failed: %d", n)
		}
		fmt.Printf("experiments: %d, executions: %d%s\n", d.experiments, d.executions, failed)
		// Only once the progress line is finished, so the two streams stay readable
		// when redirected to different places.
		for _, problem := range d.problems {
			fmt.Fprintf(os.Stderr, "  %s\n", problem)
		}
	}
	fmt.Printf("Written %d experiments with %d executions\n", experiments, executions)
	if failedExperiments > 0 || failedExecutions > 0 {
		fmt.Fprintf(os.Stderr, "Incomplete: %d experiments and %d executions could not be dumped\n", failedExperiments, failedExecutions)
		return ErrIncomplete
	}
	return nil
}

// getJSON reads a response into target once the request has been checked.
func getJSON(resp *http.Response, err error) func(target any) error {
	return func(target any) error {
		body, _, err := platform.Read(resp, err)
		if err != nil {
			return err
		}
		return json.Unmarshal(body, target)
	}
}

// Keys are matched case-insensitively, as they are shown and typed. An unknown one
// aborts: a dump that quietly covers less than asked is indistinguishable from a full one.
func selectTeams(teams []team, keys []string) ([]team, error) {
	if len(keys) == 0 {
		return teams, nil
	}
	wanted := map[string]bool{}
	var order []string
	for _, k := range keys {
		upper := strings.ToUpper(k)
		if !wanted[upper] {
			order = append(order, upper)
		}
		wanted[upper] = true
	}
	var selected []team
	found := map[string]bool{}
	for _, t := range teams {
		if wanted[strings.ToUpper(t.Key)] {
			selected = append(selected, t)
			found[strings.ToUpper(t.Key)] = true
		}
	}
	var missing []string
	for _, k := range order {
		if !found[k] {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		available := make([]string, 0, len(teams))
		for _, t := range teams {
			available = append(available, t.Key)
		}
		sort.Strings(available)
		return nil, fmt.Errorf("No accessible team with key %s. Available: %s", strings.Join(missing, ", "), strings.Join(available, ", "))
	}
	return selected, nil
}

type teamDump struct {
	experiments, executions, failedExperiments, failedExecutions int
	problems                                                     []string
}

type experimentDump struct {
	executions, failedExecutions int
	failed                       bool
	problem                      string
}

func dumpTeam(ctx context.Context, c *platform.Client, list []listed, dir string, datatype output.Datatype) teamDump {
	results := make([]experimentDump, len(list))
	forEach(len(list), experimentConcurrency, func(i int) {
		results[i] = dumpExperiment(ctx, c, list[i].Key, filepath.Join(dir, list[i].Key), datatype)
	})
	var d teamDump
	for _, r := range results {
		if r.failed {
			d.failedExperiments++
		} else {
			d.experiments++
		}
		d.executions += r.executions
		d.failedExecutions += r.failedExecutions
		if r.problem != "" {
			d.problems = append(d.problems, r.problem)
		}
	}
	return d
}

// A single unlucky request must not discard the whole walk: failures are counted and
// reported instead of ending the command.
func dumpExperiment(ctx context.Context, c *platform.Client, key, dir string, datatype output.Datatype) experimentDump {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return experimentDump{failed: true, problem: fmt.Sprintf("%s: %s", key, err)}
	}
	var (
		wg         sync.WaitGroup
		document   Document
		designErr  error
		executions struct {
			Executions []struct {
				ID int64 `json:"id"`
			} `json:"executions"`
		}
		listErr error
	)
	wg.Add(2)
	go func() { defer wg.Done(); document, designErr = Fetch(ctx, c, key) }()
	go func() {
		defer wg.Done()
		listErr = getJSON(c.GetExperimentExecutions3(ctx, key, nil))(&executions)
	}()
	wg.Wait()
	if designErr != nil {
		return experimentDump{failed: true, problem: fmt.Sprintf("%s: %s", key, designErr)}
	}
	if listErr != nil {
		return experimentDump{failed: true, problem: fmt.Sprintf("%s: Failed to get the executions. HTTP request failed: %s", key, listErr)}
	}
	removeDeprecatedFields(document.Value())
	if err := os.WriteFile(filepath.Join(dir, "experiment."+string(datatype)), document.RenderFile(datatype), 0o644); err != nil {
		return experimentDump{failed: true, problem: fmt.Sprintf("%s: %s", key, err)}
	}

	ok := make([]bool, len(executions.Executions))
	forEach(len(ok), executionConcurrency, func(i int) {
		id := executions.Executions[i].ID
		body, _, err := platform.Read(c.GetExperimentExecution(ctx, id, nil))
		if err != nil {
			return
		}
		execution, err := output.ParseDocument(body)
		if err != nil {
			return
		}
		ok[i] = os.WriteFile(filepath.Join(dir, fmt.Sprintf("execution-%d.%s", id, datatype)), execution.RenderFile(datatype), 0o644) == nil
	})
	r := experimentDump{}
	for _, written := range ok {
		if written {
			r.executions++
		} else {
			r.failedExecutions++
		}
	}
	if r.failedExecutions > 0 {
		r.problem = fmt.Sprintf("%s: %d of %d executions could not be fetched", key, r.failedExecutions, len(ok))
	}
	return r
}

// The query and list of a step's radius are deprecated and left out of dumps.
func removeDeprecatedFields(experiment *jsyaml.Map) {
	lanes, _ := experiment.Get("lanes")
	laneList, _ := lanes.([]any)
	for _, lane := range laneList {
		laneMap, _ := lane.(*jsyaml.Map)
		if laneMap == nil {
			continue
		}
		steps, _ := laneMap.Get("steps")
		stepList, _ := steps.([]any)
		for _, step := range stepList {
			stepMap, _ := step.(*jsyaml.Map)
			if stepMap == nil {
				continue
			}
			if radius, ok := stepMap.Get("radius"); ok {
				if radiusMap, ok := radius.(*jsyaml.Map); ok {
					radiusMap.Delete("query")
					radiusMap.Delete("list")
				}
			}
		}
	}
}

// forEach runs fn for 0..n-1 with at most limit calls in flight.
func forEach(n, limit int, fn func(int)) {
	slots := make(chan struct{}, max(1, limit))
	var wg sync.WaitGroup
	for i := range n {
		wg.Add(1)
		slots <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-slots }()
			fn(i)
		}()
	}
	wg.Wait()
}
