// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package experiment_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/steadybit/cli/internal/experiment"
	"github.com/steadybit/cli/internal/jsyaml"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/platformtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	platform.RetryUnit = time.Millisecond
	experiment.PollInterval = time.Millisecond
}

var ctx = context.Background()

const design = `{"key":"TST-1","version":3,"name":"Verify TTR","team":"TST","environment":"Global","lanes":[{"steps":[{"type":"wait","parameters":{"duration":"10s"}}]}]}`

func TestGetPrintsYAMLWithoutTheVersion(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/experiments/TST-1", platformtest.Reply{Body: design})

	out, err := platformtest.Stdout(t, func() error { return experiment.Get(ctx, p.Client, experiment.GetOptions{Key: "TST-1"}) })

	require.NoError(t, err)
	assert.Equal(t, "key: TST-1\nname: Verify TTR\nteam: TST\nenvironment: Global\nlanes:\n  - steps:\n      - type: wait\n        parameters:\n          duration: 10s\n\n", out)
}

func TestGetWritesCompactJSONFiles(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/experiments/TST-1", platformtest.Reply{Body: design})
	file := filepath.Join(t.TempDir(), "experiment.json")

	out, err := platformtest.Stdout(t, func() error {
		return experiment.Get(ctx, p.Client, experiment.GetOptions{Key: "TST-1", File: file})
	})

	require.NoError(t, err)
	assert.Equal(t, "Experiment TST-1 written to "+file+".\n", out)
	content, _ := os.ReadFile(file)
	assert.Equal(t, `{"key":"TST-1","name":"Verify TTR","team":"TST","environment":"Global","lanes":[{"steps":[{"type":"wait","parameters":{"duration":"10s"}}]}]}`, string(content))
}

func TestGetReportsAMissingExperiment(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/experiments/TST-9", platformtest.Reply{Status: http.StatusNotFound})

	err := experiment.Get(ctx, p.Client, experiment.GetOptions{Key: "TST-9"})

	assert.EqualError(t, err, "Experiment TST-9 not found.")
}

func TestApplyCreatesAndPrependsTheKeyKeepingTheFile(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments", platformtest.Reply{Status: http.StatusCreated, Headers: map[string]string{"Location": p.URL + "/api/experiments/NEW-1"}})
	file := filepath.Join(t.TempDir(), "experiment.yml")
	original := "# kept\nname: new\nlanes:\n  - steps:\n      - &w\n        type: wait\n      - <<: *w\n"
	require.NoError(t, os.WriteFile(file, []byte(original), 0o644))

	out, err := platformtest.Stdout(t, func() error { return experiment.Apply(ctx, p.Client, experiment.ApplyOptions{Files: []string{file}}) })

	require.NoError(t, err)
	assert.Equal(t, "Experiment NEW-1 created.\n", out)
	content, _ := os.ReadFile(file)
	assert.Equal(t, "key: NEW-1\n"+original, string(content))
	// Anchors and merge keys are resolved in what is sent.
	assert.Equal(t, map[string]any{"name": "new", "lanes": []any{map[string]any{"steps": []any{
		map[string]any{"type": "wait"}, map[string]any{"type": "wait"},
	}}}}, p.Requests("POST /api/experiments")[0].JSON(t))
}

func TestApplyUpdatesByTheKeyInTheFile(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/TST-1", platformtest.Reply{})
	file := filepath.Join(t.TempDir(), "experiment.json")
	require.NoError(t, os.WriteFile(file, []byte(design), 0o644))

	out, err := platformtest.Stdout(t, func() error { return experiment.Apply(ctx, p.Client, experiment.ApplyOptions{Files: []string{file}}) })

	require.NoError(t, err)
	assert.Equal(t, "Experiment TST-1 updated.\n", out)
}

func TestApplyRefusesAKeyWithSeveralFiles(t *testing.T) {
	dir := t.TempDir()
	for _, f := range []string{"a.yml", "b.yml"} {
		require.NoError(t, os.WriteFile(filepath.Join(dir, f), []byte("name: x\n"), 0o644))
	}

	err := experiment.Apply(ctx, nil, experiment.ApplyOptions{Key: "TST-1", Files: []string{dir}})

	assert.EqualError(t, err, "If --key is specified, at most one --file can be specified.")
}

func started(p *platformtest.Platform, key string, run int) platformtest.Reply {
	return platformtest.Reply{Status: http.StatusCreated, JSON: map[string]any{
		"key": key, "executionId": run,
		"apiLocation": p.URL + "/api/experiments/executions/1",
		"uiLocation":  "https://ui/" + key,
	}, Headers: map[string]string{"Location": "https://elsewhere.example.com/api/experiments/executions/1"}}
}

func TestRunByKeyAndWaitPollsTheConfiguredHost(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/TST-1/execute", started(p, "TST-1", 1))
	var polls atomic.Int32
	p.Handle("GET /api/experiments/executions/1", func(platformtest.Request) platformtest.Reply {
		if polls.Add(1) < 2 {
			return platformtest.Reply{JSON: map[string]any{"id": 1, "key": "TST-1", "state": "RUNNING"}}
		}
		return platformtest.Reply{JSON: map[string]any{"id": 1, "key": "TST-1", "state": "COMPLETED"}}
	})

	out, err := platformtest.Stdout(t, func() error {
		return experiment.Run(ctx, p.Client, experiment.RunOptions{Key: "TST-1", Yes: true, Wait: true})
	})

	require.NoError(t, err)
	assert.Equal(t, "Executing experiment: TST-1\nExperiment run API: https://elsewhere.example.com/api/experiments/executions/1\nExperiment run UI: https://ui/TST-1\nCurrent run state: running\nCurrent run state: completed\n", out)
	query := p.Requests("POST /api/experiments/TST-1/execute")[0].Query
	assert.Equal(t, []string{"false"}, query["allowParallel"])
	assert.Equal(t, []string{"true"}, query["forcePersist"])
}

func TestAFailedRunFailsTheCommand(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/TST-1/execute", started(p, "TST-1", 1))
	p.Reply("GET /api/experiments/executions/1", platformtest.Reply{JSON: map[string]any{"id": 1, "key": "TST-1", "state": "FAILED", "reason": "hypothesis violated"}})

	_, err := platformtest.Stdout(t, func() error {
		return experiment.Run(ctx, p.Client, experiment.RunOptions{Key: "TST-1", Yes: true, Wait: true})
	})

	assert.EqualError(t, err, "Experiment TST-1 (#1) failed, reason: hypothesis violated")
}

func TestRunRetriesValidationErrorsWithoutPersistingThem(t *testing.T) {
	p := platformtest.New(t)
	var calls atomic.Int32
	p.Handle("POST /api/experiments/TST-1/execute", func(platformtest.Request) platformtest.Reply {
		if calls.Add(1) <= 2 {
			return platformtest.Reply{Status: http.StatusUnprocessableEntity, JSON: map[string]any{"title": "no targets"}}
		}
		return started(p, "TST-1", 1)
	})

	out, err := platformtest.Stdout(t, func() error {
		return experiment.Run(ctx, p.Client, experiment.RunOptions{Key: "TST-1", Yes: true, Retries: 2})
	})

	require.NoError(t, err)
	assert.Contains(t, out, "Experiment has validation errors (attempt 1/3). Retrying in 0s...")
	for _, r := range p.Requests("POST /api/experiments/TST-1/execute") {
		assert.Equal(t, []string{"false"}, r.Query["forcePersist"])
	}
}

func TestRunGivesUpAfterTheLastRetry(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/TST-1/execute", platformtest.Reply{Status: http.StatusUnprocessableEntity, JSON: map[string]any{"title": "no targets"}})

	_, err := platformtest.Stdout(t, func() error {
		return experiment.Run(ctx, p.Client, experiment.RunOptions{Key: "TST-1", Yes: true, Retries: 1})
	})

	assert.ErrorContains(t, err, "Failed to execute experiment: Steadybit API at POST")
	assert.Len(t, p.Requests("POST /api/experiments/TST-1/execute"), 2)
}

func TestRunsInParallelWithYesWhenAnotherIsRunning(t *testing.T) {
	p := platformtest.New(t)
	p.Handle("POST /api/experiments/TST-1/execute", func(r platformtest.Request) platformtest.Reply {
		if r.Query["allowParallel"][0] == "true" {
			return started(p, "TST-1", 1)
		}
		return platformtest.Reply{Status: http.StatusConflict, JSON: map[string]any{"type": "https://steadybit.com/problems/another-experiment-running-exception"}}
	})

	_, err := platformtest.Stdout(t, func() error {
		return experiment.Run(ctx, p.Client, experiment.RunOptions{Key: "TST-1", Yes: true})
	})

	require.NoError(t, err)
	assert.Len(t, p.Requests("POST /api/experiments/TST-1/execute"), 2)
}

func TestRunByFileWithoutKeyUpsertsAndWritesTheKey(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/execute", started(p, "NEW-1", 1))
	file := filepath.Join(t.TempDir(), "experiment.yml")
	require.NoError(t, os.WriteFile(file, []byte("name: new\n"), 0o644))

	_, err := platformtest.Stdout(t, func() error {
		return experiment.Run(ctx, p.Client, experiment.RunOptions{Files: []string{file}, Yes: true})
	})

	require.NoError(t, err)
	content, _ := os.ReadFile(file)
	assert.Equal(t, "key: NEW-1\nname: new\n", string(content))
}

func TestRunNeedsSomethingToRun(t *testing.T) {
	err := experiment.Run(ctx, nil, experiment.RunOptions{Yes: true})

	assert.EqualError(t, err, "Either --key, --file or --template must be specified.")
}

const templateID = "d7e65100-1d20-4980-be87-c351704910b8"

func templateOptions() experiment.TemplateOptions {
	return experiment.TemplateOptions{Template: templateID, Team: "ADM", Placeholder: jsyaml.NewMap(), Variable: jsyaml.NewMap(), ExecutionVariable: jsyaml.NewMap(), ResetProperties: true}
}

func TestRunFromATemplate(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/templates/"+templateID+"/experiment-execute", started(p, "ADM-12", 7))
	o := experiment.RunOptions{Yes: true, TemplateOptions: templateOptions()}
	o.Placeholder.Set("CLUSTER", "prod")
	o.ExecutionVariable.Set("region", "eu")

	_, err := platformtest.Stdout(t, func() error { return experiment.Run(ctx, p.Client, o) })

	require.NoError(t, err)
	r := p.Requests("POST /api/experiments/templates/" + templateID + "/experiment-execute")[0]
	assert.Equal(t, map[string]any{"team": "ADM", "placeholders": []any{map[string]any{"key": "CLUSTER", "value": "prod"}}, "executionVariables": map[string]any{"region": "eu"}}, r.JSON(t))
	assert.Equal(t, []string{"true"}, r.Query["resetProperties"])
}

func TestApplyFromATemplateMergesAPlaceholdersFile(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/templates/"+templateID+"/experiment-create", platformtest.Reply{Status: http.StatusCreated, Headers: map[string]string{"Location": p.URL + "/api/experiments/ADM-12"}})
	file := filepath.Join(t.TempDir(), "values.yml")
	require.NoError(t, os.WriteFile(file, []byte("CLUSTER: dev\nREPLICAS: 3\n"), 0o644))
	o := templateOptions()
	o.PlaceholdersFile = file
	o.Placeholder.Set("CLUSTER", "prod")
	o.ResetProperties = false

	out, err := platformtest.Stdout(t, func() error { return experiment.ApplyTemplate(ctx, p.Client, "", o) })

	require.NoError(t, err)
	assert.Equal(t, "Experiment ADM-12 created from template "+templateID+".\n", out)
	r := p.Requests("POST /api/experiments/templates/" + templateID + "/experiment-create")[0]
	assert.Equal(t, []any{map[string]any{"key": "CLUSTER", "value": "prod"}, map[string]any{"key": "REPLICAS", "value": float64(3)}}, r.JSON(t).(map[string]any)["placeholders"])
	assert.Equal(t, []string{"false"}, r.Query["resetProperties"])
}

func TestApplyFromATemplateRefusesWhatAnUpdateIgnores(t *testing.T) {
	o := templateOptions()
	o.Environment = "Global"

	err := experiment.ApplyTemplate(ctx, nil, "ADM-12", o)

	assert.EqualError(t, err, "Updating experiment ADM-12 from a template only takes placeholders; remove --team, --environment.")
}

func TestApplyFromATemplateRejectsAnInvalidPlaceholdersFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "values.yml")
	require.NoError(t, os.WriteFile(file, []byte("- just\n- strings\n"), 0o644))
	o := templateOptions()
	o.PlaceholdersFile = file

	err := experiment.ApplyTemplate(ctx, nil, "", o)

	assert.ErrorContains(t, err, "must be a map of key to value or a list of {key, value} entries.")
}

func TestDumpWritesEveryTeamAndCountsWhatFailed(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/teams", platformtest.Reply{JSON: map[string]any{"teams": []any{map[string]any{"key": "TST", "name": "Test"}}}})
	p.Reply("GET /api/experiments", platformtest.Reply{JSON: map[string]any{"experiments": []any{map[string]any{"key": "TST-1"}}}})
	p.Reply("GET /api/experiments/TST-1", platformtest.Reply{Body: `{"key":"TST-1","lanes":[{"steps":[{"type":"action","radius":{"query":"x","list":[],"percentage":50}}]}]}`})
	p.Reply("GET /api/experiments/TST-1/executions", platformtest.Reply{JSON: map[string]any{"executions": []any{map[string]any{"id": 1}, map[string]any{"id": 2}}}})
	p.Reply("GET /api/experiments/executions/1", platformtest.Reply{JSON: map[string]any{"id": 1}})
	p.Reply("GET /api/experiments/executions/2", platformtest.Reply{Status: http.StatusInternalServerError})
	dir := t.TempDir()

	out, err := platformtest.Stdout(t, func() error { return experiment.Dump(ctx, p.Client, experiment.DumpOptions{Directory: dir}) })

	assert.ErrorIs(t, err, experiment.ErrIncomplete)
	assert.Equal(t, "Listing experiments for 1 team.\nFetching experiments for team Test (TST)... experiments: 1, executions: 1, failed: 1\nWritten 1 experiments with 1 executions\n", out)
	design, _ := os.ReadFile(filepath.Join(dir, "TST-1", "experiment.yaml"))
	assert.Equal(t, "key: TST-1\nlanes:\n  - steps:\n      - type: action\n        radius:\n          percentage: 50\n", string(design))
	_, err = os.Stat(filepath.Join(dir, "TST-1", "execution-1.yaml"))
	assert.NoError(t, err)
}

func TestDumpRefusesAnUnknownTeam(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/teams", platformtest.Reply{JSON: map[string]any{"teams": []any{map[string]any{"key": "B"}, map[string]any{"key": "A"}}}})

	err := experiment.Dump(ctx, p.Client, experiment.DumpOptions{Directory: t.TempDir(), Teams: []string{"a", "nope"}})

	assert.EqualError(t, err, "No accessible team with key NOPE. Available: A, B")
}

func TestDelete(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("DELETE /api/experiments/TST-1", platformtest.Reply{})
	p.Reply("DELETE /api/experiments/TST-9", platformtest.Reply{Status: http.StatusNotFound})

	out, err := platformtest.Stdout(t, func() error { return experiment.Delete(ctx, p.Client, "TST-1") })

	require.NoError(t, err)
	assert.Equal(t, "Experiment TST-1 deleted.\n", out)
	assert.EqualError(t, experiment.Delete(ctx, p.Client, "TST-9"), "Experiment TST-9 not found.")
}
