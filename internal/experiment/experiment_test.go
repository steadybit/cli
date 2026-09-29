// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package experiment_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/steadybit/cli/v6/internal/experiment"
	"github.com/steadybit/cli/v6/internal/interrupt"
	"github.com/steadybit/cli/v6/internal/jsyaml"
	"github.com/steadybit/cli/v6/internal/output"
	"github.com/steadybit/cli/v6/internal/platform"
	"github.com/steadybit/cli/v6/internal/platformtest"
	"github.com/steadybit/cli/v6/internal/prompt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	platform.RetryUnit = time.Millisecond
	experiment.PollInterval = time.Millisecond
	experiment.StartCheckDelay = time.Millisecond
}

// stillRunning answers the look --no-wait takes at the run it started.
func stillRunning(p *platformtest.Platform) {
	p.Reply("GET /api/experiments/executions/1", platformtest.Reply{JSON: map[string]any{"id": 1, "key": "TST-1", "state": "RUNNING"}})
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

// Each of these came out as a different experiment when the key was simply put in front.
func TestApplyWritesTheKeyIntoFilesALineInFrontWouldBreak(t *testing.T) {
	for name, tc := range map[string]struct{ file, original, written string }{
		"document marker": {"e.yml", "---\nname: new\n", "---\nkey: NEW-1\nname: new\n"},
		"empty key":       {"e.yml", "key: ''\nname: new\n", "key: NEW-1\nname: new\n"},
		"flow mapping":    {"e.yml", "{name: new}\n", "key: NEW-1\nname: new\n"},
		"json with a bom": {"e.json", "\ufeff{\"name\":\"new\"}", `{"key":"NEW-1","name":"new"}`},
	} {
		t.Run(name, func(t *testing.T) {
			p := platformtest.New(t)
			p.Reply("POST /api/experiments", platformtest.Reply{Status: http.StatusCreated, Headers: map[string]string{"Location": p.URL + "/api/experiments/NEW-1"}})
			p.Reply("POST /api/experiments/NEW-1", platformtest.Reply{})
			file := filepath.Join(t.TempDir(), tc.file)
			require.NoError(t, os.WriteFile(file, []byte(tc.original), 0o644))
			apply := func() error { return experiment.Apply(ctx, p.Client, experiment.ApplyOptions{Files: []string{file}}) }

			_, err := platformtest.Stdout(t, apply)
			require.NoError(t, err)
			content, _ := os.ReadFile(file)
			assert.Equal(t, tc.written, string(content))

			// The next apply updates what the first one created, with the whole design.
			_, err = platformtest.Stdout(t, apply)
			require.NoError(t, err)
			assert.Len(t, p.Requests("POST /api/experiments"), 1)
			assert.Equal(t, map[string]any{"key": "NEW-1", "name": "new"}, p.Requests("POST /api/experiments/NEW-1")[0].JSON(t))
		})
	}
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
	// By key, the TypeScript CLI printed the locations before the key; a file run the other way round.
	assert.Equal(t, "Experiment run API: https://elsewhere.example.com/api/experiments/executions/1\nExperiment run UI: https://ui/TST-1\nExecuting experiment: TST-1\nCurrent run state: running\nCurrent run state: completed\n", out)
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

// Only the last attempt is kept on the platform, so that it shows what was wrong without
// the attempts before it, as the run-experiment action does.
func TestRunRetriesValidationErrorsKeepingOnlyTheLast(t *testing.T) {
	p := platformtest.New(t)
	stillRunning(p)
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
	var persisted []string
	for _, r := range p.Requests("POST /api/experiments/TST-1/execute") {
		persisted = append(persisted, r.Query["forcePersist"][0])
	}
	assert.Equal(t, []string{"false", "false", "true"}, persisted)
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
	stillRunning(p)
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
	stillRunning(p)
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
	stillRunning(p)
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

func finished(state, reason string) platformtest.Reply {
	return platformtest.Reply{JSON: map[string]any{
		"id": 1, "key": "TST-1", "name": "Verify TTR", "state": state, "reason": reason,
		"started": "2026-09-25T10:00:00Z", "ended": "2026-09-25T10:00:42Z",
		"steps": []any{
			map[string]any{"stepType": "wait", "state": "COMPLETED", "parameters": map[string]any{"duration": "10s"}, "started": "2026-09-25T10:00:00Z", "ended": "2026-09-25T10:00:10Z"},
			map[string]any{"stepType": "action", "actionId": "com.steadybit.extension_http.check", "state": state, "reason": reason, "started": "2026-09-25T10:00:10Z", "ended": "2026-09-25T10:00:42Z"},
		},
	}}
}

func TestRunWritesAJUnitReport(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/TST-1/execute", started(p, "TST-1", 1))
	p.Reply("GET /api/experiments/executions/1", finished("FAILED", "HTTP check failed"))
	report := filepath.Join(t.TempDir(), "report.xml")

	_, err := platformtest.Stdout(t, func() error {
		return experiment.Run(ctx, p.Client, experiment.RunOptions{Key: "TST-1", Yes: true, Wait: true, Report: report})
	})

	assert.EqualError(t, err, "Experiment TST-1 (#1) failed, reason: HTTP check failed")
	content, _ := os.ReadFile(report)
	assert.Equal(t, `<?xml version="1.0" encoding="UTF-8"?>
<testsuites name="steadybit" tests="2" failures="1" errors="0" time="42.000">
  <testsuite name="TST-1 Verify TTR" tests="2" failures="1" errors="0" skipped="0" time="42.000" timestamp="2026-09-25T10:00:00Z">
    <properties>
      <property name="executionId" value="1"></property>
      <property name="uiLocation" value="https://ui/TST-1"></property>
    </properties>
    <testcase classname="TST-1" name="1. wait 10s" time="10.000"></testcase>
    <testcase classname="TST-1" name="2. com.steadybit.extension_http.check" time="32.000">
      <failure message="failed: HTTP check failed" type="FAILED">HTTP check failed</failure>
    </testcase>
  </testsuite>
</testsuites>
`, string(content))
	// Only a report asks the platform for the steps.
	assert.Equal(t, []string{"steps"}, p.Requests("GET /api/experiments/executions/1")[0].Query["fields"])
}

func TestRunWritesAJSONReportAndAGitHubSummary(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/TST-1/execute", started(p, "TST-1", 1))
	p.Reply("GET /api/experiments/executions/1", finished("COMPLETED", ""))
	dir := t.TempDir()
	summary := filepath.Join(dir, "summary.md")
	t.Setenv("GITHUB_STEP_SUMMARY", summary)

	_, err := platformtest.Stdout(t, func() error {
		return experiment.Run(ctx, p.Client, experiment.RunOptions{Key: "TST-1", Yes: true, Wait: true, Report: filepath.Join(dir, "r.json")})
	})

	require.NoError(t, err)
	report, _ := os.ReadFile(filepath.Join(dir, "r.json"))
	assert.Contains(t, string(report), `"state": "COMPLETED"`)
	content, _ := os.ReadFile(summary)
	assert.Equal(t, "### ✅ Steadybit experiment TST-1 · Verify TTR\n\nRun [#1](https://ui/TST-1) completed after 42s\n\n| # | Step | State | Duration |\n|---|---|---|---|\n| 1 | wait 10s | completed | 10s |\n| 2 | com.steadybit.extension_http.check | completed | 32s |\n\n", string(content))
}

func TestWaitingWithoutAReportDoesNotAskForSteps(t *testing.T) {
	p := platformtest.New(t)
	t.Setenv("GITHUB_STEP_SUMMARY", "")
	p.Reply("POST /api/experiments/TST-1/execute", started(p, "TST-1", 1))
	p.Reply("GET /api/experiments/executions/1", finished("COMPLETED", ""))

	_, err := platformtest.Stdout(t, func() error {
		return experiment.Run(ctx, p.Client, experiment.RunOptions{Key: "TST-1", Yes: true, Wait: true})
	})

	require.NoError(t, err)
	assert.Nil(t, p.Requests("GET /api/experiments/executions/1")[0].Query["fields"])
}

func TestATimeoutCancelsTheRun(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/TST-1/execute", started(p, "TST-1", 1))
	p.Reply("GET /api/experiments/executions/1", platformtest.Reply{JSON: map[string]any{"id": 1, "key": "TST-1", "state": "RUNNING"}})
	p.Reply("POST /api/experiments/executions/1/cancel", platformtest.Reply{Status: http.StatusAccepted})

	_, err := platformtest.Stdout(t, func() error {
		return experiment.Run(ctx, p.Client, experiment.RunOptions{Key: "TST-1", Yes: true, Wait: true, WaitOptions: experiment.WaitOptions{Timeout: time.Nanosecond}})
	})

	assert.ErrorIs(t, err, experiment.ErrTimedOut)
	assert.ErrorContains(t, err, "Experiment TST-1 (#1) did not end within 1ns and was canceled")
	assert.Len(t, p.Requests("POST /api/experiments/executions/1/cancel"), 1)
}

func TestShowStepsPrintsEachChange(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/TST-1/execute", started(p, "TST-1", 1))
	p.Reply("GET /api/experiments/executions/1", finished("COMPLETED", ""))

	out, err := platformtest.Stdout(t, func() error {
		return experiment.Run(ctx, p.Client, experiment.RunOptions{Key: "TST-1", Yes: true, Wait: true, WaitOptions: experiment.WaitOptions{ShowSteps: true}})
	})

	require.NoError(t, err)
	assert.Contains(t, out, "Current run state: completed\n  step 1/2 wait 10s: completed\n  step 2/2 com.steadybit.extension_http.check: completed\n")
}

func TestAReportNeedsWaiting(t *testing.T) {
	err := experiment.Run(ctx, nil, experiment.RunOptions{Key: "TST-1", Yes: true, Report: "r.xml"})

	assert.EqualError(t, err, "--report needs to wait for the runs to end; remove --no-wait.")
}

// An aborted pipeline must not leave the attack it started running on its own.
func TestAnInterruptCancelsTheRunItWaitsFor(t *testing.T) {
	t.Cleanup(interrupt.Reset)
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/TST-1/execute", started(p, "TST-1", 1))
	var canceled, interrupted atomic.Bool
	p.Handle("GET /api/experiments/executions/1", func(platformtest.Request) platformtest.Reply {
		// Before the first poll has been answered: the run id comes from where it started.
		if !interrupted.Swap(true) {
			interrupt.RunHandlers(os.Interrupt)
		}
		if canceled.Load() {
			return platformtest.Reply{JSON: map[string]any{"id": 1, "key": "TST-1", "state": "CANCELED"}}
		}
		return platformtest.Reply{JSON: map[string]any{"id": 1, "key": "TST-1", "state": "RUNNING"}}
	})
	p.Handle("POST /api/experiments/executions/1/cancel", func(platformtest.Request) platformtest.Reply {
		canceled.Store(true)
		return platformtest.Reply{Status: http.StatusAccepted}
	})

	_, err := platformtest.Stdout(t, func() error {
		return experiment.Run(ctx, p.Client, experiment.RunOptions{Key: "TST-1", Yes: true, Wait: true})
	})

	assert.EqualError(t, err, "Experiment TST-1 (#1) canceled")
	assert.Len(t, p.Requests("POST /api/experiments/executions/1/cancel"), 1)
}

func TestKeepRunningOnInterruptLeavesTheRunAlone(t *testing.T) {
	t.Cleanup(interrupt.Reset)
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/TST-1/execute", started(p, "TST-1", 1))
	var polls atomic.Int32
	p.Handle("GET /api/experiments/executions/1", func(platformtest.Request) platformtest.Reply {
		if polls.Add(1) == 2 {
			interrupt.RunHandlers(os.Interrupt)
		}
		state := "RUNNING"
		if polls.Load() >= 3 {
			state = "COMPLETED"
		}
		return platformtest.Reply{JSON: map[string]any{"id": 1, "key": "TST-1", "state": state}}
	})

	_, err := platformtest.Stdout(t, func() error {
		return experiment.Run(ctx, p.Client, experiment.RunOptions{Key: "TST-1", Yes: true, Wait: true, WaitOptions: experiment.WaitOptions{KeepRunningOnInterrupt: true}})
	})

	require.NoError(t, err)
	assert.Empty(t, p.Requests("POST /api/experiments/executions/1/cancel"))
}

func TestJQFiltersWhatGetPrints(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/experiments/TST-1", platformtest.Reply{Body: design})
	output.JQ = ".lanes[0].steps[0].parameters.duration"
	t.Cleanup(func() { output.JQ = "" })

	out, err := platformtest.Stdout(t, func() error { return experiment.Get(ctx, p.Client, experiment.GetOptions{Key: "TST-1"}) })

	require.NoError(t, err)
	assert.Equal(t, "10s\n", out)
}

func TestInitAsksForThePlaceholdersAndWritesTheExperiment(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/experiments/templates", platformtest.Reply{JSON: map[string]any{"templates": []any{
		map[string]any{"id": templateID, "templateTitle": "Checkout survives latency"},
	}}})
	p.Reply("GET /api/experiments/templates/"+templateID, platformtest.Reply{JSON: map[string]any{
		"templateTitle": "Checkout survives latency",
		"placeholders":  []any{map[string]any{"key": "DELAY", "name": "Delay", "description": "How slow the network gets."}},
	}})
	p.Reply("POST /api/experiments/templates/"+templateID+"/experiment-create", platformtest.Reply{Status: http.StatusCreated, Headers: map[string]string{"Location": p.URL + "/api/experiments/ADM-7"}})
	p.Reply("GET /api/experiments/ADM-7", platformtest.Reply{Body: `{"key":"ADM-7","name":"Checkout survives latency","team":"ADM"}`})
	experiment.Interactive = func() bool { return true }
	t.Cleanup(func() { experiment.Interactive = func() bool { return false } })
	file := filepath.Join(t.TempDir(), "checkout.yml")
	// search, template number, placeholder, team, environment (default), file
	prompt.UseInput(strings.NewReader("checkout\n1\n500ms\nADM\n\n" + file + "\n"))

	out, err := platformtest.Stdout(t, func() error { return experiment.Init(ctx, p.Client, experiment.InitOptions{}) })

	require.NoError(t, err)
	assert.Contains(t, out, "How slow the network gets.\n? Delay (DELAY): ")
	assert.Contains(t, out, "Experiment ADM-7 created. Run it with:\n\n    steadybit experiment run -f "+file+"\n")
	assert.Equal(t, map[string]any{"team": "ADM", "environment": "Global", "placeholders": []any{map[string]any{"key": "DELAY", "value": "500ms"}}},
		p.Requests("POST /api/experiments/templates/" + templateID + "/experiment-create")[0].JSON(t))
	content, _ := os.ReadFile(file)
	assert.Equal(t, "key: ADM-7\nname: Checkout survives latency\nteam: ADM\n", string(content))
}

func TestInitNeedsATerminal(t *testing.T) {
	experiment.Interactive = func() bool { return false }

	err := experiment.Init(ctx, nil, experiment.InitOptions{})

	assert.ErrorContains(t, err, "needs a terminal. In scripts, use `experiment apply --template`.")
}

// The report shows the run once the platform has stopped it, not the moment it was cut.
func TestATimedOutRunIsReportedAsCanceled(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/TST-1/execute", started(p, "TST-1", 1))
	var canceled atomic.Bool
	p.Handle("GET /api/experiments/executions/1", func(platformtest.Request) platformtest.Reply {
		state := "RUNNING"
		if canceled.Load() {
			state = "CANCELED"
		}
		return platformtest.Reply{JSON: map[string]any{"id": 1, "key": "TST-1", "state": state,
			"steps": []any{map[string]any{"stepType": "WAIT", "state": state}}}}
	})
	p.Handle("POST /api/experiments/executions/1/cancel", func(platformtest.Request) platformtest.Reply {
		canceled.Store(true)
		return platformtest.Reply{Status: http.StatusAccepted}
	})
	report := filepath.Join(t.TempDir(), "r.xml")

	_, err := platformtest.Stdout(t, func() error {
		return experiment.Run(ctx, p.Client, experiment.RunOptions{Key: "TST-1", Yes: true, Wait: true, Report: report, WaitOptions: experiment.WaitOptions{Timeout: time.Nanosecond}})
	})

	assert.ErrorIs(t, err, experiment.ErrTimedOut)
	content, _ := os.ReadFile(report)
	assert.Contains(t, string(content), `<testcase classname="TST-1" name="1. wait" time="0.000">`)
	assert.Contains(t, string(content), `<testcase classname="TST-1" name="run" time="0.000">
      <error message="canceled: did not end within 1ns" type="CANCELED">did not end within 1ns</error>`)
	assert.Contains(t, string(content), `errors="1" skipped="1"`)
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

// The question names the experiment, and --yes answers it, as in the TypeScript CLI.
func TestAsksBeforeRunningInParallel(t *testing.T) {
	for name, tc := range map[string]struct {
		yes      bool
		answer   string
		parallel bool
	}{
		"answered yes": {answer: "y\n", parallel: true},
		"answered no":  {answer: "\n"},
		"--yes":        {yes: true, parallel: true},
	} {
		t.Run(name, func(t *testing.T) {
			p := platformtest.New(t)
			stillRunning(p)
			p.Handle("POST /api/experiments/TST-1/execute", func(r platformtest.Request) platformtest.Reply {
				if q := r.Query["allowParallel"]; len(q) == 1 && q[0] == "true" {
					return started(p, "TST-1", 1)
				}
				return platformtest.Reply{Status: http.StatusConflict, Body: `{"type":"https://steadybit.com/problems/another-experiment-running-exception"}`}
			})
			original := prompt.Interactive
			prompt.Interactive = func() bool { return true }
			t.Cleanup(func() { prompt.Interactive = original })
			prompt.UseInput(strings.NewReader("y\n" + tc.answer))

			out, err := platformtest.Stdout(t, func() error {
				return experiment.Run(ctx, p.Client, experiment.RunOptions{Key: "TST-1", Yes: tc.yes})
			})

			if tc.parallel {
				require.NoError(t, err)
				assert.Contains(t, out, "Executing experiment: TST-1")
			} else {
				assert.Error(t, err)
			}
			assert.Equal(t, !tc.yes, strings.Contains(out, "Do you want to start TST-1 in parallel?"), out)
		})
	}
}

// With --key, a file updates that experiment before it runs, whatever key the file has.
func TestRunWithAKeyUpdatesItFromTheFile(t *testing.T) {
	p := platformtest.New(t)
	stillRunning(p)
	p.Reply("POST /api/experiments/TST-7", platformtest.Reply{})
	p.Reply("POST /api/experiments/TST-7/execute", started(p, "TST-7", 1))
	file := filepath.Join(t.TempDir(), "e.yml")
	require.NoError(t, os.WriteFile(file, []byte("name: from file\n"), 0o644))

	out, err := platformtest.Stdout(t, func() error {
		return experiment.Run(ctx, p.Client, experiment.RunOptions{Key: "TST-7", Files: []string{file}, Yes: true})
	})

	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(out, "Executing experiment: TST-7\n"), out)
	assert.Equal(t, map[string]any{"name": "from file"}, p.Requests("POST /api/experiments/TST-7")[0].JSON(t))
	content, _ := os.ReadFile(file)
	assert.Equal(t, "name: from file\n", string(content))
}

// A file downloaded from the UI carries a version. It is not sent: the platform would
// reject it with 409 once the experiment was edited, and the file is what is meant.
func TestAVersionInTheFileIsNotSent(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments", platformtest.Reply{Status: http.StatusCreated, Headers: map[string]string{"Location": p.URL + "/api/experiments/NEW-1"}})
	p.Reply("POST /api/experiments/NEW-1", platformtest.Reply{})
	file := filepath.Join(t.TempDir(), "e.yml")
	original := "# from the UI\nname: new\nversion: 3\n"
	require.NoError(t, os.WriteFile(file, []byte(original), 0o644))
	apply := func() error { return experiment.Apply(ctx, p.Client, experiment.ApplyOptions{Files: []string{file}}) }

	_, err := platformtest.Stdout(t, apply)
	require.NoError(t, err)
	_, err = platformtest.Stdout(t, apply)
	require.NoError(t, err)

	assert.Equal(t, map[string]any{"name": "new"}, p.Requests("POST /api/experiments")[0].JSON(t))
	assert.Equal(t, map[string]any{"key": "NEW-1", "name": "new"}, p.Requests("POST /api/experiments/NEW-1")[0].JSON(t))
	content, _ := os.ReadFile(file)
	assert.Equal(t, "key: NEW-1\n"+original, string(content))
}

// The platform accepts a run and may cancel or error it right after; --no-wait must not
// pass then. A run that failed ran, and its result is what --no-wait does not wait for.
func TestNoWaitChecksThatTheRunStarted(t *testing.T) {
	for state, want := range map[string]string{
		"RUNNING":   "",
		"COMPLETED": "",
		"FAILED":    "",
		"CANCELED":  "Experiment TST-1 (#1) canceled, reason: The run was started via CLI, but another experiment was running in parallel.",
		"ERRORED":   "Experiment TST-1 (#1) errored, reason: The run was started via CLI, but another experiment was running in parallel.",
	} {
		t.Run(state, func(t *testing.T) {
			p := platformtest.New(t)
			p.Reply("POST /api/experiments/TST-1/execute", started(p, "TST-1", 1))
			p.Reply("GET /api/experiments/executions/1", platformtest.Reply{JSON: map[string]any{"id": 1, "key": "TST-1", "state": state,
				"reason": "The run was started via CLI, but another experiment was running in parallel."}})

			_, err := platformtest.Stdout(t, func() error {
				return experiment.Run(ctx, p.Client, experiment.RunOptions{Key: "TST-1", Yes: true})
			})

			if want == "" {
				require.NoError(t, err)
			} else {
				assert.EqualError(t, err, want)
			}
		})
	}
}

// Validation can take longer than the first look: the check keeps looking while the run is
// still being prepared, and the failed run is in the summary of the job it failed.
func TestNoWaitKeepsCheckingWhileTheRunIsPrepared(t *testing.T) {
	p := platformtest.New(t)
	summary := filepath.Join(t.TempDir(), "summary.md")
	t.Setenv("GITHUB_STEP_SUMMARY", summary)
	p.Reply("POST /api/experiments/TST-1/execute", started(p, "TST-1", 1))
	var polls atomic.Int32
	p.Handle("GET /api/experiments/executions/1", func(platformtest.Request) platformtest.Reply {
		state := []string{"CREATED", "PREPARED", "CANCELED"}[min(polls.Add(1)-1, 2)]
		return platformtest.Reply{JSON: map[string]any{"id": 1, "key": "TST-1", "state": state, "reason": "another experiment was running"}}
	})

	_, err := platformtest.Stdout(t, func() error {
		return experiment.Run(ctx, p.Client, experiment.RunOptions{Key: "TST-1", Yes: true})
	})

	assert.EqualError(t, err, "Experiment TST-1 (#1) canceled, reason: another experiment was running")
	assert.Equal(t, int32(3), polls.Load())
	content, _ := os.ReadFile(summary)
	assert.Contains(t, string(content), "Run [#1](https://ui/TST-1) canceled")
}

// The run was started either way; only being unable to look at it is no reason to fail.
func TestNoWaitOnlyWarnsWhenTheRunCannotBeChecked(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/TST-1/execute", started(p, "TST-1", 1))
	p.Reply("GET /api/experiments/executions/1", platformtest.Reply{Status: http.StatusBadGateway, JSON: map[string]any{"title": "Bad Gateway"}})

	var out string
	stderr, err := platformtest.Stderr(t, func() error {
		var err error
		out, err = platformtest.Stdout(t, func() error {
			return experiment.Run(ctx, p.Client, experiment.RunOptions{Key: "TST-1", Yes: true})
		})
		return err
	})

	require.NoError(t, err)
	assert.Contains(t, out, "Executing experiment: TST-1\n")
	assert.Contains(t, stderr, "Could not check that the experiment run started: ")
	assert.Contains(t, stderr, `"title": "Bad Gateway"`)
}

// --no-wait does not wait for a run the platform takes long to validate: past the
// deadline, it warns and passes.
func TestNoWaitGivesUpCheckingAfterTheTimeout(t *testing.T) {
	timeout := experiment.StartCheckTimeout
	experiment.StartCheckTimeout = 50 * time.Millisecond
	t.Cleanup(func() { experiment.StartCheckTimeout = timeout })
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/TST-1/execute", started(p, "TST-1", 1))
	p.Reply("GET /api/experiments/executions/1", platformtest.Reply{JSON: map[string]any{"id": 1, "key": "TST-1", "state": "CREATED"}})

	stderr, err := platformtest.Stderr(t, func() error {
		_, err := platformtest.Stdout(t, func() error {
			return experiment.Run(ctx, p.Client, experiment.RunOptions{Key: "TST-1", Yes: true})
		})
		return err
	})

	require.NoError(t, err)
	assert.Equal(t, "Could not check that the experiment run started: the run was still created after 50ms\n", stderr)
}

// Three experiment files, two at a time: the platform is told they run in parallel on
// purpose, a failed run does not stop the others, and the report has all three in order.
func TestRunsFilesInParallel(t *testing.T) {
	p := platformtest.New(t)
	var mu sync.Mutex
	inFlight, most := 0, 0
	polls := map[string]int{}
	for id, n := range map[int]string{1: "1", 2: "2", 3: "3"} {
		key := "TST-" + n
		p.Reply("POST /api/experiments/"+key, platformtest.Reply{})
		p.Handle("POST /api/experiments/"+key+"/execute", func(r platformtest.Request) platformtest.Reply {
			assert.Equal(t, []string{"true"}, r.Query["allowParallel"], key)
			mu.Lock()
			inFlight++
			most = max(most, inFlight)
			mu.Unlock()
			return platformtest.Reply{Status: http.StatusCreated, JSON: map[string]any{"key": key, "executionId": id,
				"apiLocation": p.URL + "/api/experiments/executions/" + n, "uiLocation": "https://ui/" + key}}
		})
		p.Handle("GET /api/experiments/executions/"+n, func(platformtest.Request) platformtest.Reply {
			mu.Lock()
			defer mu.Unlock()
			polls[key]++
			state := "RUNNING"
			if polls[key] > 20 {
				state = "COMPLETED"
				if key == "TST-2" {
					state = "FAILED"
				}
				inFlight--
			}
			return platformtest.Reply{JSON: map[string]any{"id": id, "key": key, "name": key, "state": state}}
		})
	}
	dir := t.TempDir()
	var files []string
	for _, n := range []string{"1", "2", "3"} {
		file := filepath.Join(dir, "e"+n+".yml")
		require.NoError(t, os.WriteFile(file, []byte("key: TST-"+n+"\nname: TST-"+n+"\n"), 0o644))
		files = append(files, file)
	}
	report := filepath.Join(dir, "report.json")

	out, err := platformtest.Stdout(t, func() error {
		return experiment.Run(ctx, p.Client, experiment.RunOptions{Files: files, Yes: true, Wait: true, Parallel: 2, Report: report})
	})

	assert.EqualError(t, err, "Experiment TST-2 (#2) failed")
	assert.Equal(t, 2, most, "at most two runs at a time")
	assert.Contains(t, out, "[TST-1] Current run state: completed\n")
	assert.Contains(t, out, "[TST-3] Current run state: completed\n")
	content, err := os.ReadFile(report)
	require.NoError(t, err)
	var runs []map[string]any
	require.NoError(t, json.Unmarshal(content, &runs))
	var keys []any
	for _, run := range runs {
		keys = append(keys, run["key"])
	}
	assert.Equal(t, []any{"TST-1", "TST-2", "TST-3"}, keys)
}

// Interrupted, every run still going is canceled, not only the last one started.
func TestAnInterruptCancelsEveryParallelRun(t *testing.T) {
	t.Cleanup(interrupt.Reset)
	p := platformtest.New(t)
	var mu sync.Mutex
	canceled := map[int]bool{}
	var polled atomic.Int32
	var once sync.Once
	dir := t.TempDir()
	var files []string
	for _, id := range []int{1, 2} {
		key := fmt.Sprintf("TST-%d", id)
		p.Reply("POST /api/experiments/"+key, platformtest.Reply{})
		p.Reply("POST /api/experiments/"+key+"/execute", platformtest.Reply{Status: http.StatusCreated, JSON: map[string]any{"key": key, "executionId": id,
			"apiLocation": fmt.Sprintf("%s/api/experiments/executions/%d", p.URL, id), "uiLocation": "https://ui/" + key}})
		p.Handle(fmt.Sprintf("GET /api/experiments/executions/%d", id), func(platformtest.Request) platformtest.Reply {
			// Once both runs have been polled, both waits have their handler in place.
			if polled.Add(1) >= 2 {
				once.Do(func() { go interrupt.RunHandlers(os.Interrupt) })
			}
			mu.Lock()
			defer mu.Unlock()
			state := "RUNNING"
			if canceled[id] {
				state = "CANCELED"
			}
			return platformtest.Reply{JSON: map[string]any{"id": id, "key": key, "state": state}}
		})
		p.Handle(fmt.Sprintf("POST /api/experiments/executions/%d/cancel", id), func(platformtest.Request) platformtest.Reply {
			mu.Lock()
			defer mu.Unlock()
			canceled[id] = true
			return platformtest.Reply{Status: http.StatusAccepted}
		})
		file := filepath.Join(dir, key+".yml")
		require.NoError(t, os.WriteFile(file, []byte("key: "+key+"\n"), 0o644))
		files = append(files, file)
	}

	_, err := platformtest.Stdout(t, func() error {
		return experiment.Run(ctx, p.Client, experiment.RunOptions{Files: files, Yes: true, Wait: true, Parallel: 2})
	})

	assert.Error(t, err)
	assert.Len(t, p.Requests("POST /api/experiments/executions/1/cancel"), 1)
	assert.Len(t, p.Requests("POST /api/experiments/executions/2/cancel"), 1)
}

// Interrupted while two runs go, the third waiting for a slot is not started: the
// handlers that cancel runs were already collected and would leave it running.
func TestAnInterruptStartsNoFurtherParallelRun(t *testing.T) {
	t.Cleanup(interrupt.Reset)
	p := platformtest.New(t)
	var mu sync.Mutex
	canceled, polls := map[int]bool{}, map[int]int{}
	var polled atomic.Int32
	var once sync.Once
	dir := t.TempDir()
	var files []string
	for _, id := range []int{1, 2, 3} {
		key := fmt.Sprintf("TST-%d", id)
		p.Reply("POST /api/experiments/"+key, platformtest.Reply{})
		p.Reply("POST /api/experiments/"+key+"/execute", platformtest.Reply{Status: http.StatusCreated, JSON: map[string]any{"key": key, "executionId": id,
			"apiLocation": fmt.Sprintf("%s/api/experiments/executions/%d", p.URL, id), "uiLocation": "https://ui/" + key}})
		p.Handle(fmt.Sprintf("GET /api/experiments/executions/%d", id), func(platformtest.Request) platformtest.Reply {
			if polled.Add(1) >= 2 {
				once.Do(func() { go interrupt.RunHandlers(os.Interrupt) })
			}
			mu.Lock()
			defer mu.Unlock()
			polls[id]++
			// A run nobody cancels ends on its own, so that a regression fails, not hangs.
			state := "RUNNING"
			if canceled[id] {
				state = "CANCELED"
			} else if polls[id] > 50 {
				state = "COMPLETED"
			}
			return platformtest.Reply{JSON: map[string]any{"id": id, "key": key, "state": state}}
		})
		p.Handle(fmt.Sprintf("POST /api/experiments/executions/%d/cancel", id), func(platformtest.Request) platformtest.Reply {
			mu.Lock()
			defer mu.Unlock()
			canceled[id] = true
			return platformtest.Reply{Status: http.StatusAccepted}
		})
		file := filepath.Join(dir, key+".yml")
		require.NoError(t, os.WriteFile(file, []byte("key: "+key+"\n"), 0o644))
		files = append(files, file)
	}
	report := filepath.Join(dir, "report.json")

	_, err := platformtest.Stdout(t, func() error {
		return experiment.Run(ctx, p.Client, experiment.RunOptions{Files: files, Yes: true, Wait: true, Parallel: 2, Report: report})
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "TST-3: not started, the command was interrupted")
	assert.Empty(t, p.Requests("POST /api/experiments/TST-3/execute"))
	assert.Len(t, p.Requests("POST /api/experiments/executions/1/cancel"), 1)
	assert.Len(t, p.Requests("POST /api/experiments/executions/2/cancel"), 1)
	content, err := os.ReadFile(report)
	require.NoError(t, err)
	var runs []map[string]any
	require.NoError(t, json.Unmarshal(content, &runs))
	var states []any
	for _, run := range runs {
		states = append(states, run["key"].(string)+" "+run["state"].(string))
	}
	assert.Equal(t, []any{"TST-1 CANCELED", "TST-2 CANCELED", "TST-3 CANCELED"}, states)
}

// A run whose start the interrupt overtook has no handler yet when the handlers run; it
// is canceled as soon as its wait begins.
func TestARunStartedDuringAnInterruptIsCanceled(t *testing.T) {
	t.Cleanup(interrupt.Reset)
	p := platformtest.New(t)
	var canceled atomic.Bool
	var polls atomic.Int32
	p.Handle("POST /api/experiments/TST-1/execute", func(platformtest.Request) platformtest.Reply {
		interrupt.RunHandlers(os.Interrupt)
		return started(p, "TST-1", 1)
	})
	p.Handle("GET /api/experiments/executions/1", func(platformtest.Request) platformtest.Reply {
		// Left alone, the run ends on its own, so that a regression fails, not hangs.
		state := "RUNNING"
		if canceled.Load() {
			state = "CANCELED"
		} else if polls.Add(1) > 50 {
			state = "COMPLETED"
		}
		return platformtest.Reply{JSON: map[string]any{"id": 1, "key": "TST-1", "state": state}}
	})
	p.Handle("POST /api/experiments/executions/1/cancel", func(platformtest.Request) platformtest.Reply {
		canceled.Store(true)
		return platformtest.Reply{Status: http.StatusAccepted}
	})

	_, err := platformtest.Stdout(t, func() error {
		return experiment.Run(ctx, p.Client, experiment.RunOptions{Key: "TST-1", Yes: true, Wait: true})
	})

	assert.EqualError(t, err, "Experiment TST-1 (#1) canceled")
	assert.Len(t, p.Requests("POST /api/experiments/executions/1/cancel"), 1)
}

func TestParallelNeedsWaiting(t *testing.T) {
	err := experiment.Run(ctx, nil, experiment.RunOptions{Key: "TST-1", Yes: true, Parallel: 2})

	assert.EqualError(t, err, "--parallel needs to wait for the runs; remove --no-wait.")
}

// A single run has nothing to run in parallel with: it keeps the platform's check for
// another experiment running.
func TestParallelWithASingleRunDoesNotAllowParallelRuns(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/TST-1/execute", started(p, "TST-1", 1))
	p.Reply("GET /api/experiments/executions/1", finished("COMPLETED", ""))

	_, err := platformtest.Stdout(t, func() error {
		return experiment.Run(ctx, p.Client, experiment.RunOptions{Key: "TST-1", Yes: true, Wait: true, Parallel: 2})
	})

	require.NoError(t, err)
	assert.Equal(t, []string{"false"}, p.Requests("POST /api/experiments/TST-1/execute")[0].Query["allowParallel"])
}

// A run that could not start is named in the error, its retries are told apart from the
// other runs', and it is in the report, in file order.
func TestAParallelRunThatCannotStartIsReported(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/TST-1", platformtest.Reply{})
	p.Reply("POST /api/experiments/TST-1/execute", platformtest.Reply{Status: http.StatusUnprocessableEntity, JSON: map[string]any{"title": "No targets"}})
	p.Reply("POST /api/experiments/TST-2", platformtest.Reply{})
	p.Reply("POST /api/experiments/TST-2/execute", platformtest.Reply{Status: http.StatusCreated, JSON: map[string]any{"key": "TST-2", "executionId": 2,
		"apiLocation": p.URL + "/api/experiments/executions/2", "uiLocation": "https://ui/TST-2"}})
	p.Reply("GET /api/experiments/executions/2", platformtest.Reply{JSON: map[string]any{"id": 2, "key": "TST-2", "state": "COMPLETED"}})
	dir := t.TempDir()
	var files []string
	for _, key := range []string{"TST-1", "TST-2"} {
		file := filepath.Join(dir, key+".yml")
		require.NoError(t, os.WriteFile(file, []byte("key: "+key+"\n"), 0o644))
		files = append(files, file)
	}
	report := filepath.Join(dir, "report.json")

	out, err := platformtest.Stdout(t, func() error {
		return experiment.Run(ctx, p.Client, experiment.RunOptions{Files: files, Yes: true, Wait: true, Parallel: 2, Retries: 1, Report: report})
	})

	require.Error(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "TST-1: Failed to execute experiment: "), err.Error())
	assert.Contains(t, out, "[TST-1] Experiment has validation errors (attempt 1/2). Retrying in 0s...\n")
	content, err := os.ReadFile(report)
	require.NoError(t, err)
	var runs []map[string]any
	require.NoError(t, json.Unmarshal(content, &runs))
	require.Len(t, runs, 2)
	assert.Equal(t, "TST-1", runs[0]["key"])
	assert.Equal(t, "ERRORED", runs[0]["state"])
	assert.Contains(t, runs[0]["reason"], "Failed to execute experiment")
	assert.Equal(t, "TST-2", runs[1]["key"])
	assert.Equal(t, "COMPLETED", runs[1]["state"])
}

// run is one poll of a run, for the expectation tests.
func runState(state, reason string) platformtest.Reply {
	return platformtest.Reply{JSON: map[string]any{"id": 1, "key": "TST-1", "state": state, "reason": reason}}
}

func TestExpectedStatesAndReasons(t *testing.T) {
	for name, tc := range map[string]struct {
		states     []string
		reason     string
		expect     experiment.WaitOptions
		err        string
		lastPolled int
	}{
		"a failure that was expected passes":       {states: []string{"RUNNING", "FAILED"}, reason: "Check failure.", expect: experiment.WaitOptions{ExpectState: "FAILED"}},
		"RUNNING passes before the run ends":       {states: []string{"CREATED", "RUNNING", "COMPLETED"}, expect: experiment.WaitOptions{ExpectState: "RUNNING"}, lastPolled: 2},
		"another end fails, naming both":           {states: []string{"COMPLETED"}, expect: experiment.WaitOptions{ExpectState: "FAILED"}, err: "Experiment TST-1 (#1) completed, but failed was expected"},
		"the reason has to match exactly":          {states: []string{"FAILED"}, reason: "Check failure.", expect: experiment.WaitOptions{ExpectState: "FAILED", ExpectReason: "Timeout."}, err: `Experiment TST-1 (#1) failed with reason "Check failure.", but the reason "Timeout." was expected`},
		"without an expectation, as it always was": {states: []string{"FAILED"}, reason: "Check failure.", err: "Experiment TST-1 (#1) failed, reason: Check failure."},
	} {
		t.Run(name, func(t *testing.T) {
			p := platformtest.New(t)
			p.Reply("POST /api/experiments/TST-1/execute", started(p, "TST-1", 1))
			var polls atomic.Int32
			p.Handle("GET /api/experiments/executions/1", func(platformtest.Request) platformtest.Reply {
				i := int(polls.Add(1)) - 1
				if i >= len(tc.states) {
					i = len(tc.states) - 1
				}
				return runState(tc.states[i], tc.reason)
			})

			_, err := platformtest.Stdout(t, func() error {
				return experiment.Run(ctx, p.Client, experiment.RunOptions{Key: "TST-1", Yes: true, Wait: true, WaitOptions: tc.expect})
			})

			if tc.err == "" {
				require.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.err)
			}
			if tc.lastPolled > 0 {
				assert.EqualValues(t, tc.lastPolled, polls.Load(), "stops polling once the state is reached")
			}
		})
	}
}

// The platform refuses a run while another one goes; --busy-retries waits instead of
// starting in parallel, which --yes would otherwise do.
func TestBusyRetriesWaitInsteadOfRunningInParallel(t *testing.T) {
	busy := platformtest.Reply{Status: http.StatusUnprocessableEntity, Body: `{"type":"https://steadybit.com/problems/another-experiment-running-exception","title":"Another experiment is running"}`}
	for name, tc := range map[string]struct {
		refusals int
		err      string
	}{
		"it starts once the other has ended": {refusals: 2},
		"it gives up after its tries":        {refusals: 5, err: "Failed to execute experiment, another one was still running after 3 tries"},
	} {
		t.Run(name, func(t *testing.T) {
			p := platformtest.New(t)
			var calls atomic.Int32
			p.Handle("POST /api/experiments/TST-1/execute", func(r platformtest.Request) platformtest.Reply {
				assert.Equal(t, []string{"false"}, r.Query["allowParallel"])
				if int(calls.Add(1)) <= tc.refusals {
					return busy
				}
				return started(p, "TST-1", 1)
			})
			p.Reply("GET /api/experiments/executions/1", runState("COMPLETED", ""))

			out, err := platformtest.Stdout(t, func() error {
				return experiment.Run(ctx, p.Client, experiment.RunOptions{Key: "TST-1", Yes: true, Wait: true, BusyRetries: 3})
			})

			if tc.err == "" {
				require.NoError(t, err)
				assert.Contains(t, out, "Another experiment is running, trying again in 0s (2/3).")
			} else {
				assert.ErrorContains(t, err, tc.err)
				assert.EqualValues(t, 4, calls.Load(), "the first try and three more")
			}
		})
	}
}

// The platform may also accept a run and cancel it right away for the same reason.
func TestBusyRetriesAlsoCoverARunCanceledForAnother(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/TST-1/execute", started(p, "TST-1", 1))
	var polls atomic.Int32
	p.Handle("GET /api/experiments/executions/1", func(platformtest.Request) platformtest.Reply {
		if polls.Add(1) == 1 {
			return runState("CANCELED", "The run was started via CLI, but another experiment was running in parallel.")
		}
		return runState("COMPLETED", "")
	})

	_, err := platformtest.Stdout(t, func() error {
		return experiment.Run(ctx, p.Client, experiment.RunOptions{Key: "TST-1", Yes: true, Wait: true, BusyRetries: 1})
	})

	require.NoError(t, err)
	assert.Len(t, p.Requests("POST /api/experiments/TST-1/execute"), 2)
}

func TestExpectationRetriesRunTheExperimentAgain(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/TST-1/execute", started(p, "TST-1", 1))
	var polls atomic.Int32
	p.Handle("GET /api/experiments/executions/1", func(platformtest.Request) platformtest.Reply {
		if polls.Add(1) <= 2 {
			return runState("FAILED", "flaky")
		}
		return runState("COMPLETED", "")
	})
	report := filepath.Join(t.TempDir(), "run.json")

	out, err := platformtest.Stdout(t, func() error {
		return experiment.Run(ctx, p.Client, experiment.RunOptions{Key: "TST-1", Yes: true, Wait: true, ExpectationRetries: 2, Report: report})
	})

	require.NoError(t, err)
	assert.Len(t, p.Requests("POST /api/experiments/TST-1/execute"), 3)
	assert.Contains(t, out, "Experiment run 1 did not end as expected (attempt 1/3). Running it again in 0s.")
	// The report holds the last run, with the API location run-experiment's output gives.
	content, err := os.ReadFile(report)
	require.NoError(t, err)
	var runs []map[string]any
	require.NoError(t, json.Unmarshal(content, &runs))
	require.Len(t, runs, 1)
	assert.Equal(t, "COMPLETED", runs[0]["state"])
	assert.Equal(t, "https://elsewhere.example.com/api/experiments/executions/1", runs[0]["apiLocation"])
}

func TestRunByExternalID(t *testing.T) {
	p := platformtest.New(t)
	p.Handle("GET /api/experiments", func(r platformtest.Request) platformtest.Reply {
		switch r.Query["externalId"][0] {
		case "shop":
			return platformtest.Reply{JSON: map[string]any{"experiments": []any{map[string]any{"key": "TST-1"}}}}
		case "twice":
			return platformtest.Reply{JSON: map[string]any{"experiments": []any{map[string]any{"key": "TST-1"}, map[string]any{"key": "TST-2"}}}}
		}
		return platformtest.Reply{JSON: map[string]any{"experiments": []any{}}}
	})
	p.Reply("POST /api/experiments/TST-1/execute", started(p, "TST-1", 1))
	p.Reply("GET /api/experiments/executions/1", runState("COMPLETED", ""))
	run := func(o experiment.RunOptions) error {
		o.Yes, o.Wait = true, true
		_, err := platformtest.Stdout(t, func() error { return experiment.Run(ctx, p.Client, o) })
		return err
	}

	require.NoError(t, run(experiment.RunOptions{TemplateOptions: experiment.TemplateOptions{ExternalID: "shop"}}))
	assert.EqualError(t, run(experiment.RunOptions{TemplateOptions: experiment.TemplateOptions{ExternalID: "gone"}}), "No experiment has the external id 'gone'.")
	assert.EqualError(t, run(experiment.RunOptions{TemplateOptions: experiment.TemplateOptions{ExternalID: "twice"}}), "2 experiments have the external id 'twice'; run one of them with --key.")
	assert.EqualError(t, run(experiment.RunOptions{Key: "TST-1", TemplateOptions: experiment.TemplateOptions{ExternalID: "shop"}}), "--external-id finds the experiment to run; leave out --key.")
}

func TestExpectationsNeedWaitingAndARealState(t *testing.T) {
	err := experiment.Run(ctx, nil, experiment.RunOptions{Key: "TST-1", Yes: true, Wait: true, WaitOptions: experiment.WaitOptions{ExpectState: "done"}})
	assert.EqualError(t, err, "--expect-state must be one of CREATED, PREPARED, RUNNING, FAILED, CANCELED, COMPLETED, ERRORED, not 'DONE'.")
	err = experiment.Run(ctx, nil, experiment.RunOptions{Key: "TST-1", Yes: true, WaitOptions: experiment.WaitOptions{ExpectState: "FAILED"}})
	assert.EqualError(t, err, "--expect-state, --expect-reason and --expectation-retries need to wait for the run; remove --no-wait.")
}
