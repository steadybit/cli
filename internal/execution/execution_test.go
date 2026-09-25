// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package execution_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/steadybit/cli/internal/execution"
	"github.com/steadybit/cli/internal/output"
	"github.com/steadybit/cli/internal/platformtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

const run = `{"id":42,"key":"ADM-1","state":"COMPLETED","steps":[
 {"stepType":"wait"},
 {"stepType":"action","actionId":"com.steadybit.extension_jmeter.run","targetExecutions":[
   {"id":"te-1","name":"host-a","artifacts":["report.zip","log.txt"]},
   {"id":"te-2","name":"host-b","artifacts":["report.zip"]}]},
 {"stepType":"service-validation","customLabel":"shop is healthy","validations":[
   {"stepType":"action","targetExecutions":[{"id":"te-3","name":"check","artifacts":["result.json"]}]}]}]}`

func TestCollectsArtifactsOfActionsAndServiceValidations(t *testing.T) {
	doc, err := output.ParseDocument([]byte(run))
	require.NoError(t, err)

	assert.Equal(t, []execution.Artifact{
		{Step: "com.steadybit.extension_jmeter.run", Target: "host-a", TargetExecutionID: "te-1", ArtifactID: "report.zip"},
		{Step: "com.steadybit.extension_jmeter.run", Target: "host-a", TargetExecutionID: "te-1", ArtifactID: "log.txt"},
		{Step: "com.steadybit.extension_jmeter.run", Target: "host-b", TargetExecutionID: "te-2", ArtifactID: "report.zip"},
		{Step: "shop is healthy", Target: "check", TargetExecutionID: "te-3", ArtifactID: "result.json"},
	}, execution.Collect(doc.Value()))
}

// The platform leaves the steps out unless asked, and with them every artifact.
func TestAsksForTheSteps(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/experiments/executions/42", platformtest.Reply{JSON: map[string]any{"id": 42, "state": "RUNNING"}})

	out, err := platformtest.Stdout(t, func() error { return execution.Get(ctx, p.Client, execution.GetOptions{ID: 42}) })

	require.NoError(t, err)
	assert.Equal(t, "id: 42\nstate: RUNNING\n\n", out)
	assert.Equal(t, []string{"steps"}, p.Requests("GET /api/experiments/executions/42")[0].Query["fields"])
}

func TestCancel(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/executions/42/cancel", platformtest.Reply{Status: http.StatusAccepted})
	p.Reply("POST /api/experiments/executions/43/cancel", platformtest.Reply{Status: http.StatusOK})
	p.Reply("POST /api/experiments/executions/44/cancel", platformtest.Reply{Status: http.StatusNotFound})

	out, err := platformtest.Stdout(t, func() error {
		if err := execution.Cancel(ctx, p.Client, 42); err != nil {
			return err
		}
		return execution.Cancel(ctx, p.Client, 43)
	})
	require.NoError(t, err)
	assert.Equal(t, "Experiment run 42 is being canceled.\nExperiment run 43 has already ended.\n", out)
	assert.EqualError(t, execution.Cancel(ctx, p.Client, 44), "Experiment run 44 not found.")
}

func TestPropertyValues(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/executions/42/properties/*/set", platformtest.Reply{})
	p.Reply("POST /api/experiments/executions/42/properties/*/add", platformtest.Reply{})
	set := func(values []string, asJSON bool) string {
		_, err := platformtest.Stdout(t, func() error {
			return execution.SetProperty(ctx, p.Client, execution.PropertyOptions{ID: 42, Key: "k", Values: values, JSON: asJSON})
		})
		require.NoError(t, err)
		requests := p.Requests("POST /api/experiments/executions/42/properties/k/set")
		return string(requests[len(requests)-1].Body)
	}

	assert.Equal(t, `"0042"`, set([]string{"0042"}, false), "a number-like value stays a string")
	assert.Equal(t, `["a","b"]`, set([]string{"a", "b"}, false))
	assert.Equal(t, `7`, set([]string{"7"}, true))
	assert.Equal(t, `0`, set([]string{"0"}, true), "falsy values are sent")
	assert.Equal(t, `false`, set([]string{"false"}, true))
	assert.Equal(t, `""`, set([]string{""}, false))

	err := execution.SetProperty(ctx, p.Client, execution.PropertyOptions{ID: 42, Key: "k", Values: []string{"seven"}, JSON: true})
	assert.ErrorContains(t, err, "'seven' is not valid JSON")
	assert.EqualError(t, execution.AddProperty(ctx, p.Client, execution.PropertyOptions{ID: 42, Key: "k", Values: []string{"a", "b"}}),
		"Adding to a list property takes exactly one --value.")
}

func TestPropertyErrorsNameTheProperty(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("POST /api/experiments/executions/42/properties/locked/set", platformtest.Reply{Status: 422, JSON: map[string]any{"title": "not editable"}})

	err := execution.SetProperty(ctx, p.Client, execution.PropertyOptions{ID: 42, Key: "locked", Values: []string{"x"}})

	assert.ErrorContains(t, err, "Failed to set property locked of experiment run 42: ")
	assert.ErrorContains(t, err, "not editable")
}

func TestDownloadsEveryArtifactIntoADirectoryPerTarget(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/experiments/executions/42", platformtest.Reply{Body: run})
	p.Handle("GET /api/experiments/executions/42/artifacts/*/*", func(r platformtest.Request) platformtest.Reply {
		parts := strings.Split(r.Path, "/")
		return platformtest.Reply{Body: "content of " + parts[6] + "/" + parts[7]}
	})
	dir := t.TempDir()

	_, err := platformtest.Stdout(t, func() error {
		return execution.Download(ctx, p.Client, execution.DownloadOptions{ID: 42, Directory: dir})
	})

	require.NoError(t, err)
	for _, f := range []string{"te-1/report.zip", "te-1/log.txt", "te-2/report.zip", "te-3/result.json"} {
		content, err := os.ReadFile(filepath.Join(dir, f))
		require.NoError(t, err)
		assert.Equal(t, "content of "+f, string(content))
	}
}

func TestDownloadToASingleFile(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/experiments/executions/42", platformtest.Reply{Body: run})
	p.Reply("GET /api/experiments/executions/42/artifacts/te-1/log.txt", platformtest.Reply{Body: "the log"})
	file := filepath.Join(t.TempDir(), "my.log")

	_, err := platformtest.Stdout(t, func() error {
		return execution.Download(ctx, p.Client, execution.DownloadOptions{ID: 42, Artifact: "log.txt", Output: file})
	})

	require.NoError(t, err)
	content, _ := os.ReadFile(file)
	assert.Equal(t, "the log", string(content))
	err = execution.Download(ctx, p.Client, execution.DownloadOptions{ID: 42, Artifact: "report.zip", Output: file})
	assert.EqualError(t, err, "2 artifacts match, but --output takes exactly one. Narrow it down with --artifact and --target-execution.")
	err = execution.Download(ctx, p.Client, execution.DownloadOptions{ID: 42, Artifact: "nope"})
	assert.EqualError(t, err, "No matching artifacts found in experiment run 42.")
}

func TestNeverLetsAnIdStepOutOfADirectory(t *testing.T) {
	for _, id := range []string{"..", ".", "", "../..", "a/../..", "/"} {
		assert.Equal(t, "_", execution.PathSegment(id), id)
	}
	assert.Equal(t, "evil", execution.PathSegment("../../evil"))
	assert.Equal(t, "report.zip", execution.PathSegment("report.zip"))
}

func TestWatchPrintsChangesUntilTheRunEnds(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/experiments/ADM-1/executions", platformtest.Reply{JSON: map[string]any{"executions": []any{map[string]any{"id": 41}, map[string]any{"id": 42}}}})
	var polls int
	p.Handle("GET /api/experiments/executions/42", func(platformtest.Request) platformtest.Reply {
		polls++
		state, step := "RUNNING", "RUNNING"
		if polls > 1 {
			state, step = "FAILED", "FAILED"
		}
		return platformtest.Reply{JSON: map[string]any{"id": 42, "key": "ADM-1", "state": state, "reason": "check failed", "steps": []any{
			map[string]any{"stepType": "action", "actionId": "http-check", "state": step, "targetExecutions": []any{map[string]any{"state": step}}},
		}}}
	})

	out, err := platformtest.Stdout(t, func() error {
		return execution.Watch(ctx, p.Client, execution.WatchOptions{Key: "ADM-1", Interval: time.Millisecond})
	})

	assert.EqualError(t, err, "Experiment ADM-1 (#42) failed, reason: check failed")
	assert.Equal(t, "Experiment ADM-1 run #42: running\n  step 1/1 http-check: running (targets 0/1)\nExperiment ADM-1 run #42: failed\n  step 1/1 http-check: failed (targets 1/1)\n", out)
	assert.Empty(t, p.Requests("POST /api/experiments/executions/42/cancel"), "watching never cancels")
}
