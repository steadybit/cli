// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package experiment implements `experiment get`, `apply` and `run`.
//
// Experiment designs pass through as documents rather than generated structs: decoding
// a file into typed Go values and encoding it again would drop fields the spec does not
// know yet and rewrite zero values, silently changing files kept in Git. The generated
// client still types every path, parameter and the smaller request bodies.
package experiment

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/steadybit/cli/v6/api"
	"github.com/steadybit/cli/v6/internal/interrupt"
	"github.com/steadybit/cli/v6/internal/jsyaml"
	"github.com/steadybit/cli/v6/internal/output"
	"github.com/steadybit/cli/v6/internal/platform"
	"github.com/steadybit/cli/v6/internal/prompt"
)

type Document = *output.Document

const anotherExperimentRunning = "https://steadybit.com/problems/another-experiment-running-exception"

// jsonBody is an experiment as it is sent. A version in the file, as the UI's download
// writes one, is left out: the platform rejects a stale one with 409, and the file, not
// what was edited since it was downloaded, is what an apply means to save.
func jsonBody(document Document) (io.Reader, error) {
	sent := output.NewDocument(jsyaml.Clone(document.Value()).(*jsyaml.Map))
	sent.Delete("version")
	b, err := json.Marshal(sent)
	return bytes.NewReader(b), err
}

func Fetch(ctx context.Context, c *platform.Client, key string) (Document, error) {
	body, _, err := platform.Read(c.GetExperiment(ctx, key))
	if platform.IsStatus(err, http.StatusNotFound) {
		return nil, fmt.Errorf("Experiment %s not found.", key)
	}
	if err != nil {
		return nil, platform.Failed(err, "Failed to get the experiment. HTTP request failed.")
	}
	document, err := output.ParseDocument(body)
	if err != nil {
		return nil, err
	}
	// Left out of files: kept, it would make every apply after an edit in the UI a 409.
	document.Delete("version")
	return document, nil
}

type GetOptions struct {
	Key, File, Type string
}

func Get(ctx context.Context, c *platform.Client, o GetOptions) error {
	document, err := Fetch(ctx, c, o.Key)
	if err != nil {
		return err
	}
	if output.JQ != "" && o.File == "" {
		return output.ApplyJQ(os.Stdout, jsyaml.CompactJSON(document.Value()), output.JQ)
	}
	datatype, err := output.ResolveDatatype(o.Type, o.File)
	if err != nil {
		return err
	}
	rendered, err := document.Render(datatype)
	if err != nil {
		return err
	}
	if o.File == "" {
		// As console.log printed it: YAML already ends in a newline and gets another.
		fmt.Print(string(rendered))
		if datatype == output.YAML {
			fmt.Println()
		}
		return nil
	}
	if err := os.WriteFile(o.File, document.RenderFile(datatype), 0o644); err != nil {
		return err
	}
	fmt.Printf("Experiment %s written to %s.\n", o.Key, o.File)
	return nil
}

// Delete removes an experiment.
func Delete(ctx context.Context, c *platform.Client, key string) error {
	_, _, err := platform.Read(c.DeleteExperiment(ctx, key))
	if platform.IsStatus(err, http.StatusNotFound) {
		return fmt.Errorf("Experiment %s not found.", key)
	}
	if err != nil {
		return platform.Failed(err, "Failed to delete the experiment. HTTP request failed.")
	}
	fmt.Printf("Experiment %s deleted.\n", key)
	return nil
}

// ResolveFiles expands directories into their YAML files, recursively on request.
func ResolveFiles(paths []string, recursive bool) ([]string, error) {
	var files []string
	for _, path := range paths {
		info, err := os.Stat(path)
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("File or directory '%s' not found.", path)
		}
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			files = append(files, path)
			continue
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}
		var dirs []string
		for _, entry := range entries {
			name := strings.ToLower(entry.Name())
			switch {
			case entry.IsDir():
				dirs = append(dirs, filepath.Join(path, entry.Name()))
			case strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml"):
				files = append(files, filepath.Join(path, entry.Name()))
			}
		}
		if recursive && len(dirs) > 0 {
			nested, err := ResolveFiles(dirs, recursive)
			if err != nil {
				return nil, err
			}
			files = append(files, nested...)
		}
	}
	return files, nil
}

func load(file string) (Document, output.Datatype, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, "", fmt.Errorf("Failed to read experiment file at path '%s': %w", file, err)
	}
	document, err := output.ParseDocument(content)
	if err != nil {
		return nil, "", fmt.Errorf("Failed to parse experiment file at path '%s' as YAML/JSON: %w", file, err)
	}
	datatype := output.YAML
	if output.IsJSON(content) {
		datatype = output.JSON
	}
	return document, datatype, nil
}

// writeBack puts the key the platform assigned at the top of the file, so that the next
// apply updates this experiment instead of creating another one.
func writeBack(file string, document Document, datatype output.Datatype, key string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	if datatype == output.YAML {
		if rendered, ok := output.WithFieldFirst(content, document, "key", key); ok {
			return os.WriteFile(file, rendered, 0o644)
		}
	}
	document.SetFirst("key", key)
	return os.WriteFile(file, document.RenderFile(datatype), 0o644)
}

func keyOf(document Document) string {
	key, _ := document.Get("key")
	return key
}

func keyFromLocation(resp *http.Response) string {
	location := resp.Header.Get("Location")
	return location[strings.LastIndex(location, "/")+1:]
}

func update(ctx context.Context, c *platform.Client, key string, document Document) error {
	body, err := jsonBody(document)
	if err != nil {
		return err
	}
	_, _, err = platform.Read(c.UpdateExperimentWithBody(ctx, key, "application/json", body))
	if platform.IsStatus(err, http.StatusNotFound) {
		return fmt.Errorf("Experiment %s not found.", key)
	}
	if err != nil {
		return platform.Failed(err, "Failed to save the experiment. HTTP request failed.")
	}
	return nil
}

type ApplyOptions struct {
	Key       string
	Files     []string
	Recursive bool
}

func Apply(ctx context.Context, c *platform.Client, o ApplyOptions) error {
	files, err := ResolveFiles(o.Files, o.Recursive)
	if err != nil {
		return err
	}
	if o.Key != "" && len(files) > 1 {
		return errors.New("If --key is specified, at most one --file can be specified.")
	}
	for _, file := range files {
		document, datatype, err := load(file)
		if err != nil {
			return err
		}
		key := o.Key
		if key == "" {
			key = keyOf(document)
		}
		if key != "" {
			if err := update(ctx, c, key, document); err != nil {
				return err
			}
			fmt.Printf("Experiment %s updated.\n", key)
			continue
		}
		body, err := jsonBody(document)
		if err != nil {
			return err
		}
		_, resp, err := platform.Read(c.CreateOrUpdateExperimentWithBody(ctx, "application/json", body))
		if err != nil {
			return platform.Failed(err, "Failed to save the experiment. HTTP request failed.")
		}
		key = keyFromLocation(resp)
		if resp.StatusCode == http.StatusCreated {
			if err := writeBack(file, document, datatype, key); err != nil {
				return err
			}
			fmt.Printf("Experiment %s created.\n", key)
		} else {
			fmt.Printf("Experiment %s updated.\n", key)
		}
	}
	return nil
}

type RunOptions struct {
	Key           string
	Files         []string
	Recursive     bool
	Yes, Wait     bool
	AllowParallel bool
	Retries       int
	RetryInterval int
	// Parallel is how many runs go at once; 0 or 1 runs them one after another.
	Parallel int
	// BusyRetries tries a run again, BusyRetryInterval apart, when another experiment is
	// running and running in parallel is not allowed, instead of asking or failing.
	BusyRetries       int
	BusyRetryInterval time.Duration
	// ExpectationRetries runs an experiment again, ExpectationRetryInterval apart, when
	// its run did not end as expected.
	ExpectationRetries       int
	ExpectationRetryInterval time.Duration
	WaitOptions
	// Report is a file to write a JUnit (or, for .json, JSON) report of the runs to.
	Report string

	TemplateOptions
}

type started struct {
	Key, APILocation, UILocation string
}

func Run(ctx context.Context, c *platform.Client, o RunOptions) error {
	if !o.Yes {
		ok, err := prompt.Confirm("Are you sure you want to run the experiment?", false, true)
		if err != nil {
			return err
		}
		if !ok {
			os.Exit(0)
		}
	}

	// Each run names its experiment in the question about running in parallel, as the
	// TypeScript CLI did.
	runs := []runner{}
	switch {
	case o.Template != "" && o.Key != "":
		return errors.New("--key cannot be combined with --template. Use `experiment apply --template -k` to update it.")
	case o.Template != "":
		runs = append(runs, runner{"this one", func(parallel, persist bool) (started, error) { return runTemplate(ctx, c, o, parallel, persist) }, false})
	case len(o.Files) > 0:
		files, err := ResolveFiles(o.Files, o.Recursive)
		if err != nil {
			return err
		}
		if o.Key != "" && len(files) > 1 {
			return errors.New("If --key is specified, at most one --file can be specified.")
		}
		for _, file := range files {
			runs = append(runs, runner{fileExperimentName(file), func(parallel, persist bool) (started, error) { return runFile(ctx, c, o, file, parallel, persist) }, false})
		}
	case o.Key != "":
		if o.ExternalID != "" {
			return errors.New("--external-id finds the experiment to run; leave out --key.")
		}
		runs = append(runs, runner{o.Key, func(parallel, persist bool) (started, error) { return runKey(ctx, c, o.Key, parallel, persist) }, true})
	case o.ExternalID != "":
		// Without --template, the external id names an experiment that exists, as the
		// run-experiment action's externalId does.
		key, err := keyByExternalID(ctx, c, o.ExternalID)
		if err != nil {
			return err
		}
		runs = append(runs, runner{key, func(parallel, persist bool) (started, error) { return runKey(ctx, c, key, parallel, persist) }, true})
	default:
		return errors.New("Either --key, --file or --template must be specified.")
	}

	if o.ExpectReason != "" && o.ExpectState == "" {
		o.ExpectState = "COMPLETED"
	}
	if o.ExpectState != "" {
		o.ExpectState = strings.ToUpper(o.ExpectState)
		if !slices.Contains(runStates, o.ExpectState) {
			return fmt.Errorf("--expect-state must be one of %s, not '%s'.", strings.Join(runStates, ", "), o.ExpectState)
		}
	}
	if !o.Wait && (o.ExpectState != "" || o.ExpectationRetries > 0) {
		return errors.New("--expect-state, --expect-reason and --expectation-retries need to wait for the run; remove --no-wait.")
	}
	if o.Retries < 0 || o.BusyRetries < 0 || o.ExpectationRetries < 0 {
		return errors.New("--retries, --busy-retries and --expectation-retries cannot be negative.")
	}

	if o.Parallel < 0 {
		return errors.New("--parallel cannot be negative.")
	}
	if o.Report != "" && !o.Wait {
		return errors.New("--report needs to wait for the runs to end; remove --no-wait.")
	}
	// Without waiting, a slot frees as soon as its run started: --parallel would limit how
	// fast runs start, not how many run at once.
	if o.Parallel > 1 && !o.Wait {
		return errors.New("--parallel needs to wait for the runs; remove --no-wait.")
	}
	o.WaitOptions.Steps = o.Report != "" || os.Getenv("GITHUB_STEP_SUMMARY") != ""

	var finished []*RunResult
	// The report and summary cover the runs that were started: one after the other, up to
	// and including the first one that failed, which ends the command; with --parallel,
	// every run, a failed one not stopping the others.
	report := func() error {
		if o.Report != "" {
			if err := WriteReport(o.Report, finished); err != nil {
				return fmt.Errorf("Failed to write the report to %s: %w", o.Report, err)
			}
		}
		return WriteGitHubSummary(finished)
	}
	// A single run has nothing to run in parallel with, and telling the platform it runs in
	// parallel on purpose would skip the check for another experiment running.
	if o.Parallel > 1 && len(runs) > 1 {
		return runConcurrently(ctx, c, o, runs, report, &finished)
	}
	for _, r := range runs {
		done, err := runUntilExpected(ctx, c, o, r, "", func(result started) { printStarted(result, r.keyLast) })
		if done != nil {
			finished = append(finished, done)
		}
		if err != nil {
			return errors.Join(err, report())
		}
	}
	return report()
}

// canceledForAnother is a run the platform accepted and then canceled, because another
// experiment was running.
func canceledForAnother(run *RunResult) bool {
	return run.State == "CANCELED" && strings.Contains(run.Reason, "another experiment was running")
}

// runUntilExpected starts one run and waits for it. It runs the experiment again when
// --expectation-retries asks for it and the run did not end as expected, and, with
// --busy-retries, when the platform canceled it because another experiment was running.
// It returns the last run, for the report, or nil when nothing started or --no-wait
// left a run going.
func runUntilExpected(ctx context.Context, c *platform.Client, o RunOptions, r runner, prefix string, onStart func(started)) (*RunResult, error) {
	busy := 0
	for attempt := 0; ; attempt++ {
		result, err := withRetries(o, r.what, prefix, r.run)
		if err != nil {
			return nil, err
		}
		onStart(result)
		if result.APILocation == "" {
			return nil, nil
		}
		if !o.Wait {
			// Only a run that failed the check is reported: one still running has no result yet.
			run, err := checkStarted(ctx, c, result.APILocation)
			if err != nil {
				run.UILocation, run.APILocation = result.UILocation, result.APILocation
				return run, err
			}
			return nil, nil
		}
		waitOptions := o.WaitOptions
		if prefix != "" {
			waitOptions.Prefix = "[" + result.Key + "] "
		}
		done, err := wait(ctx, c, result.APILocation, waitOptions)
		if done != nil {
			done.UILocation, done.APILocation = result.UILocation, result.APILocation
		}
		if err == nil || done == nil || !errors.Is(err, ErrUnexpected) || interrupt.Interrupted() {
			return done, err
		}
		if !o.AllowParallel && canceledForAnother(done) && busy < o.BusyRetries {
			busy++
			fmt.Printf("%sAnother experiment is running, trying again in %s (%d/%d).\n", prefix, o.BusyRetryInterval, busy, o.BusyRetries)
			time.Sleep(o.BusyRetryInterval)
			attempt--
			continue
		}
		if attempt >= o.ExpectationRetries {
			return done, err
		}
		fmt.Printf("%sExperiment run %d did not end as expected (attempt %d/%d). Running it again in %s.\n", prefix, done.ID, attempt+1, o.ExpectationRetries+1, o.ExpectationRetryInterval)
		time.Sleep(o.ExpectationRetryInterval)
	}
}

// A run by key alone printed its locations before the key, one from a file after.
func printStarted(result started, keyLast bool) {
	if !keyLast {
		fmt.Println("Executing experiment:", result.Key)
	}
	fmt.Println("Experiment run API:", result.APILocation)
	fmt.Println("Experiment run UI:", result.UILocation)
	if keyLast {
		fmt.Println("Executing experiment:", result.Key)
	}
}

// runConcurrently starts up to o.Parallel runs at a time. The platform is told they run in
// parallel on purpose, or it cancels all but the first. A failed run does not stop the
// others: every run is waited for, reported, and the command fails if any did.
func runConcurrently(ctx context.Context, c *platform.Client, o RunOptions, runs []runner, report func() error, finished *[]*RunResult) error {
	o.AllowParallel = true
	var printing sync.Mutex
	results := make([]*RunResult, len(runs))
	errs := make([]error, len(runs))
	// A run that has no result from the platform is still reported, so that the report and
	// summary name every run that failed the command.
	unreported := func(i int, state string, err error) {
		results[i] = &RunResult{Key: runs[i].what, State: state, Reason: err.Error()}
		errs[i] = fmt.Errorf("%s: %w", runs[i].what, err)
	}
	slots := make(chan struct{}, o.Parallel)
	var group sync.WaitGroup
	for i, r := range runs {
		group.Add(1)
		slots <- struct{}{}
		go func() {
			defer group.Done()
			defer func() { <-slots }()
			// Interrupted while waiting for a slot: the handlers that cancel the runs have
			// already been collected, and would miss this one.
			if interrupt.Interrupted() {
				unreported(i, "CANCELED", errors.New("not started, the command was interrupted"))
				return
			}
			done, err := runUntilExpected(ctx, c, o, r, "["+r.what+"] ", func(result started) {
				printing.Lock()
				defer printing.Unlock()
				printStarted(result, r.keyLast)
			})
			if done == nil {
				if err != nil {
					unreported(i, "ERRORED", err)
				}
				return
			}
			results[i] = done
			errs[i] = err
		}()
	}
	group.Wait()
	for _, run := range results {
		if run != nil {
			*finished = append(*finished, run)
		}
	}
	return errors.Join(append(errs, report())...)
}

// runner is one run to start: how the question about running in parallel names it, how
// to start it, and where its key goes in what is printed.
type runner struct {
	what string
	// run starts it; persist keeps a run the platform refused, so that its errors can be
	// seen in the platform.
	run     func(parallel, persist bool) (started, error)
	keyLast bool
}

// withRetries retries validation errors, which clear up once targets appear, and offers
// a parallel run when another experiment is already running, or, with --busy-retries,
// waits for it instead. Only the last attempt at a run with validation errors is kept on
// the platform, so that it shows what was wrong without the attempts before it.
// The prefix starts its messages, telling runs apart when several start at once.
func withRetries(o RunOptions, what, prefix string, run func(parallel, persist bool) (started, error)) (started, error) {
	parallel := o.AllowParallel
	busy := 0
	for attempt := 0; ; attempt++ {
		result, err := run(parallel, attempt >= o.Retries)
		if err == nil {
			return result, nil
		}
		var apiErr *platform.APIError
		if !errors.As(err, &apiErr) {
			return result, err
		}
		if !parallel && apiErr.ProblemType() == anotherExperimentRunning {
			if busy < o.BusyRetries {
				busy++
				fmt.Printf("%sAnother experiment is running, trying again in %s (%d/%d).\n", prefix, o.BusyRetryInterval, busy, o.BusyRetries)
				time.Sleep(o.BusyRetryInterval)
				attempt--
				continue
			}
			// Waiting was asked for; starting in parallel after all is what it rules out.
			if o.BusyRetries > 0 {
				return result, platform.Failed(err, "Failed to execute experiment, another one was still running after %d tries", o.BusyRetries)
			}
			ok := o.Yes
			if !ok {
				// Its own error: the platform's is what a "no" reports.
				var promptErr error
				if ok, promptErr = prompt.Confirm("There is already an experiment running. Do you want to start "+what+" in parallel?", false, false); promptErr != nil {
					return result, promptErr
				}
			}
			if ok {
				parallel = true
				attempt--
				continue
			}
		}
		if apiErr.Status == http.StatusUnprocessableEntity && attempt < o.Retries {
			fmt.Printf("%sExperiment has validation errors (attempt %d/%d). Retrying in %ds...\n", prefix, attempt+1, o.Retries+1, o.RetryInterval)
			time.Sleep(time.Duration(o.RetryInterval) * time.Second)
			continue
		}
		return result, platform.Failed(err, "Failed to execute experiment")
	}
}

func decodeStarted(body []byte, resp *http.Response, fallbackKey string) (started, error) {
	var r api.ExecuteExperimentResponseAO
	if len(body) > 0 {
		if err := json.Unmarshal(body, &r); err != nil {
			return started{}, err
		}
	}
	s := started{Key: r.Key, APILocation: resp.Header.Get("Location"), UILocation: r.UiLocation}
	if s.Key == "" {
		s.Key = fallbackKey
	}
	if s.APILocation == "" {
		s.APILocation = r.ApiLocation
	}
	return s, nil
}

func runKey(ctx context.Context, c *platform.Client, key string, parallel, persist bool) (started, error) {
	body, resp, err := platform.Read(c.ExecuteExperimentWithBody(ctx, key,
		&api.ExecuteExperimentParams{AllowParallel: &parallel, ForcePersist: &persist}, "application/json", nil))
	if err != nil {
		return started{}, err
	}
	return decodeStarted(body, resp, key)
}

func runFile(ctx context.Context, c *platform.Client, o RunOptions, file string, parallel, persist bool) (started, error) {
	document, datatype, err := load(file)
	if err != nil {
		return started{}, err
	}
	key := o.Key
	if key == "" {
		key = keyOf(document)
	}
	if key != "" {
		if err := update(ctx, c, key, document); err != nil {
			return started{}, err
		}
		return runKey(ctx, c, key, parallel, persist)
	}
	reqBody, err := jsonBody(document)
	if err != nil {
		return started{}, err
	}
	body, resp, err := platform.Read(c.SaveAndRunWithBody(ctx,
		&api.SaveAndRunParams{AllowParallel: &parallel, ForcePersist: &persist}, "application/json", reqBody))
	if err != nil {
		return started{}, err
	}
	result, err := decodeStarted(body, resp, "")
	if err != nil {
		return started{}, err
	}
	return result, writeBack(file, document, datatype, result.Key)
}

func runTemplate(ctx context.Context, c *platform.Client, o RunOptions, parallel, persist bool) (started, error) {
	id, err := templateID(o.Template)
	if err != nil {
		return started{}, err
	}
	create, err := createRequest(o.TemplateOptions)
	if err != nil {
		return started{}, err
	}
	request := api.CreateAndRunExperimentFromTemplateAO{
		Team: create.Team, Environment: create.Environment, ExternalId: create.ExternalId,
		Placeholders: create.Placeholders, ExperimentVariables: create.ExperimentVariables,
		ExecutionVariables: variables(o.ExecutionVariable),
	}
	body, resp, err := platform.Read(c.SaveAndRunFromTemplate(ctx, id,
		&api.SaveAndRunFromTemplateParams{ResetProperties: &o.ResetProperties, AllowParallel: &parallel, ForcePersist: &persist}, request))
	if platform.IsStatus(err, http.StatusNotFound) {
		return started{}, fmt.Errorf("Experiment template %s not found.", o.Template)
	}
	if err != nil {
		return started{}, err
	}
	return decodeStarted(body, resp, "")
}

// PollInterval is how often --wait asks for the state of a run. Tests shorten it.
var PollInterval = 5 * time.Second

// runStates are the states a run goes through, which --expect-state can name.
var runStates = []string{"CREATED", "PREPARED", "RUNNING", "FAILED", "CANCELED", "COMPLETED", "ERRORED"}

// keyByExternalID finds the experiment with an external id; exactly one must have it.
func keyByExternalID(ctx context.Context, c *platform.Client, externalID string) (string, error) {
	var list struct {
		Experiments []struct {
			Key string `json:"key"`
		} `json:"experiments"`
	}
	ids := []string{externalID}
	resp, err := c.GetExperiments(ctx, &api.GetExperimentsParams{ExternalId: &ids})
	if _, err := platform.Decode(resp, err, &list); err != nil {
		return "", platform.Failed(err, "Failed to find the experiment with external id %s", externalID)
	}
	switch len(list.Experiments) {
	case 0:
		return "", fmt.Errorf("No experiment has the external id '%s'.", externalID)
	case 1:
	default:
		return "", fmt.Errorf("%d experiments have the external id '%s'; run one of them with --key.", len(list.Experiments), externalID)
	}
	return list.Experiments[0].Key, nil
}

var terminal = map[string]bool{"FAILED": true, "ERRORED": true, "CANCELED": true, "COMPLETED": true}

// WaitOptions shape what `run --wait` does besides waiting.
type WaitOptions struct {
	// Timeout cancels the run once it has taken this long; zero waits indefinitely.
	Timeout time.Duration
	// KeepRunningOnInterrupt leaves the run going when the CLI is interrupted, as the
	// TypeScript CLI did. By default it is cancelled: an aborted pipeline should not
	// leave an attack running on its own.
	KeepRunningOnInterrupt bool
	// ShowSteps prints each step's state as it changes.
	ShowSteps bool
	// Prefix starts each line about the run, telling runs apart when several go at once.
	Prefix string
	// ExpectState passes the run once it reaches this state, which need not be an end
	// such as RUNNING, and fails it when it ends in another. Empty expects COMPLETED.
	ExpectState string
	// ExpectReason also requires the run's reason to be exactly this.
	ExpectReason string
	// Steps asks the platform for the steps of the run, which reports need.
	Steps bool
}

// settle waits a little for a cancelled run to end, and returns it as it last was.
func settle(ctx context.Context, c *platform.Client, path string, last *RunResult) *RunResult {
	for range 10 {
		time.Sleep(PollInterval)
		body, _, err := platform.Read(c.Get(ctx, path))
		if err != nil {
			return last
		}
		run, err := parseRun(body)
		if err != nil {
			return last
		}
		last = run
		if terminal[run.State] {
			return run
		}
	}
	return last
}

// runIDFromLocation reads the id at the end of a run's location, or 0 if there is none.
func runIDFromLocation(location string) int64 {
	u, err := url.Parse(location)
	if err != nil {
		return 0
	}
	id, _ := strconv.ParseInt(path.Base(u.Path), 10, 64)
	return id
}

// fileExperimentName is how the question about running in parallel names the experiment
// of a file: its key, or its name.
func fileExperimentName(file string) string {
	document, _, err := load(file)
	if err != nil {
		return "the experiment"
	}
	for _, field := range []string{"key", "name"} {
		if v, ok := document.Get(field); ok && v != "" {
			return v
		}
	}
	return "the experiment"
}

// ErrTimedOut is returned when --timeout cancelled the run.
var ErrTimedOut = errors.New("timed out")

// wait polls the run until it ends. A run that did not complete is an error, which is
// what lets a pipeline fail on it. The finished run is returned for reports.
func wait(ctx context.Context, c *platform.Client, location string, o WaitOptions) (*RunResult, error) {
	path := runPath(location)
	if o.Steps || o.ShowSteps {
		separator := "?"
		if strings.Contains(path, "?") {
			separator = "&"
		}
		path += separator + "fields=steps"
	}

	// Known from the location before the first poll, so an interrupt right after the
	// start still cancels the run. Read by the signal handler, hence atomic.
	var runID atomic.Int64
	runID.Store(runIDFromLocation(location))
	cancel := func(why string) {
		id := runID.Load()
		if id == 0 {
			return
		}
		fmt.Fprintf(os.Stderr, "%s%s, canceling experiment run %d.\n", o.Prefix, why, id)
		cancelCtx, done := context.WithTimeout(context.Background(), 30*time.Second)
		defer done()
		if _, _, err := platform.Read(c.CancelExperimentExecution(cancelCtx, id)); err != nil {
			fmt.Fprintf(os.Stderr, "%sFailed to cancel experiment run %d: %s\n", o.Prefix, id, err)
		}
	}
	if !o.KeepRunningOnInterrupt {
		// Once: the handler and the check below can both see the interrupt.
		var once sync.Once
		interrupted := func() { once.Do(func() { cancel("Interrupted") }) }
		pop := interrupt.Push(func(os.Signal) { interrupted() })
		defer pop()
		// Interrupted while the run was starting, the handlers ran before this one was
		// pushed: the run is canceled here, or it would be left running.
		if interrupt.Interrupted() {
			interrupted()
		}
	}

	var deadline time.Time
	if o.Timeout > 0 {
		deadline = time.Now().Add(o.Timeout)
	}
	shown := map[string]string{}
	for {
		time.Sleep(PollInterval)
		body, _, err := platform.Read(c.Get(ctx, path))
		if err != nil {
			return nil, platform.Failed(err, "Failed to get experiment run ")
		}
		run, err := parseRun(body)
		if err != nil {
			return nil, err
		}
		runID.Store(run.ID)
		fmt.Printf("%sCurrent run state: %s\n", o.Prefix, strings.ToLower(run.State))
		if o.ShowSteps {
			for i, step := range run.Steps {
				id := fmt.Sprint(i)
				if shown[id] != step.State {
					shown[id] = step.State
					fmt.Printf("%s  step %d/%d %s: %s\n", o.Prefix, i+1, len(run.Steps), step.Name, strings.ToLower(step.State))
				}
			}
		}
		if o.ExpectState != "" && run.State == o.ExpectState {
			if o.ExpectReason != "" && run.Reason != o.ExpectReason {
				return run, unexpected(fmt.Sprintf("Experiment %s (#%d) %s with reason %q, but the reason %q was expected", run.Key, run.ID, strings.ToLower(run.State), run.Reason, o.ExpectReason))
			}
			return run, nil
		}
		if !terminal[run.State] {
			if !deadline.IsZero() && time.Now().After(deadline) {
				cancel(fmt.Sprintf("Experiment run %d did not end within %s", run.ID, o.Timeout))
				// Reported once the platform has stopped it, so the report shows what was
				// cut short rather than a run that seems to be still going.
				run = settle(ctx, c, path, run)
				if run.Reason == "" {
					run.Reason = fmt.Sprintf("did not end within %s", o.Timeout)
				}
				return run, fmt.Errorf("Experiment %s (#%d) did not end within %s and was canceled: %w", run.Key, run.ID, o.Timeout, ErrTimedOut)
			}
			continue
		}
		if o.ExpectState != "" {
			return run, unexpected(notCompleted(run).Error() + ", but " + strings.ToLower(o.ExpectState) + " was expected")
		}
		if run.State != "COMPLETED" {
			return run, unexpected(notCompleted(run).Error())
		}
		return run, nil
	}
}

// ErrUnexpected marks a run that did not end as expected, by --expect-state and
// --expect-reason or by completing, which --expectation-retries runs again.
var ErrUnexpected = errors.New("the run did not end as expected")

// unexpected is such a run's error: its message says how it ended, and it is ErrUnexpected.
type unexpected string

func (e unexpected) Error() string      { return string(e) }
func (unexpected) Is(target error) bool { return target == ErrUnexpected }

func notCompleted(run *RunResult) error {
	reason := ""
	if run.Reason != "" {
		reason = ", reason: " + run.Reason
	}
	return fmt.Errorf("Experiment %s (#%d) %s%s", run.Key, run.ID, strings.ToLower(run.State), reason)
}

// StartCheckDelay is how long --no-wait gives a run before first looking at it. The
// platform accepts a run and may cancel it moments later, when its validation finds
// another experiment running; unchecked, a pipeline would pass on a run that never ran.
var StartCheckDelay = 2 * time.Second

// StartCheckTimeout bounds the whole check, requests and the client's back-off included:
// --no-wait promises not to wait for the run, so a slow platform only earns a warning.
var StartCheckTimeout = 15 * time.Second

// checkStarted polls the run until the platform is done validating it, and fails when the
// platform ended it before it ran: canceled or errored. A run that failed is the
// experiment's result, which --no-wait does not wait for. The check is a courtesy: when
// the platform cannot be asked in time, the run was still started, so it only warns.
func checkStarted(ctx context.Context, c *platform.Client, location string) (*RunResult, error) {
	ctx, done := context.WithTimeout(ctx, StartCheckTimeout)
	defer done()
	path := runPath(location)
	var last *RunResult
	warn := func(err error) (*RunResult, error) {
		if last != nil && errors.Is(ctx.Err(), context.DeadlineExceeded) {
			err = fmt.Errorf("the run was still %s after %s", strings.ToLower(last.State), StartCheckTimeout)
		}
		fmt.Fprintln(os.Stderr, platform.Failed(err, "Could not check that the experiment run started"))
		return nil, nil
	}
	for delay := StartCheckDelay; ; delay = PollInterval {
		if err := sleep(ctx, delay); err != nil {
			return warn(err)
		}
		body, _, err := platform.Read(c.Get(ctx, path))
		if err != nil {
			return warn(err)
		}
		run, err := parseRun(body)
		if err != nil {
			return warn(err)
		}
		switch run.State {
		case "CREATED", "REQUESTED", "PREPARED":
			last = run
			continue
		case "CANCELED", "ERRORED":
			return run, notCompleted(run)
		}
		return run, nil
	}
}

// sleep waits for d, or until ctx ends.
func sleep(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// runPath is a run's location relative to the API, which the client sends to its host.
func runPath(location string) string {
	if i := strings.Index(location, "/api/"); i >= 0 {
		return location[i:]
	}
	return location
}
