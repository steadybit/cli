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
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/steadybit/cli/api"
	"github.com/steadybit/cli/internal/interrupt"
	"github.com/steadybit/cli/internal/output"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/prompt"
)

type Document = *output.Document

const anotherExperimentRunning = "https://steadybit.com/problems/another-experiment-running-exception"

func jsonBody(document any) (io.Reader, error) {
	b, err := json.Marshal(document)
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
	// Removed because it makes files awkward to reapply; the API will drop it too.
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
	if json.Valid(content) {
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
	var rendered []byte
	if datatype == output.YAML {
		// Prepending keeps the rest of the file byte for byte, comments and anchors included.
		rendered = append([]byte("key: "+key+"\n"), content...)
	} else {
		document.SetFirst("key", key)
		rendered = document.RenderFile(datatype)
	}
	return os.WriteFile(file, rendered, 0o644)
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

	persist := o.Retries == 0
	runs := []func(parallel bool) (started, error){}
	switch {
	case o.Template != "" && o.Key != "":
		return errors.New("--key cannot be combined with --template. Use `experiment apply --template -k` to update it.")
	case o.Template != "":
		runs = append(runs, func(parallel bool) (started, error) { return runTemplate(ctx, c, o, parallel, persist) })
	case len(o.Files) > 0:
		files, err := ResolveFiles(o.Files, o.Recursive)
		if err != nil {
			return err
		}
		if o.Key != "" && len(files) > 1 {
			return errors.New("If --key is specified, at most one --file can be specified.")
		}
		for _, file := range files {
			runs = append(runs, func(parallel bool) (started, error) { return runFile(ctx, c, o, file, parallel, persist) })
		}
	case o.Key != "":
		runs = append(runs, func(parallel bool) (started, error) { return runKey(ctx, c, o.Key, parallel, persist) })
	default:
		return errors.New("Either --key, --file or --template must be specified.")
	}

	if o.Report != "" && !o.Wait {
		return errors.New("--report needs to wait for the runs to end; remove --no-wait.")
	}
	o.WaitOptions.Steps = o.Report != "" || os.Getenv("GITHUB_STEP_SUMMARY") != ""

	var finished []*RunResult
	// The report and summary cover the runs up to and including the first one that
	// failed, which ends the command.
	report := func() error {
		if o.Report != "" {
			if err := WriteReport(o.Report, finished); err != nil {
				return fmt.Errorf("Failed to write the report to %s: %w", o.Report, err)
			}
		}
		return WriteGitHubSummary(finished)
	}
	for _, run := range runs {
		result, err := withRetries(o, run)
		if err != nil {
			return errors.Join(err, report())
		}
		fmt.Println("Executing experiment:", result.Key)
		fmt.Println("Experiment run API:", result.APILocation)
		fmt.Println("Experiment run UI:", result.UILocation)
		if o.Wait && result.APILocation != "" {
			done, err := wait(ctx, c, result.APILocation, o.WaitOptions)
			if done != nil {
				done.UILocation = result.UILocation
				finished = append(finished, done)
			}
			if err != nil {
				return errors.Join(err, report())
			}
		}
	}
	return report()
}

// withRetries retries validation errors, which clear up once targets appear, and offers
// a parallel run when another experiment is already running.
func withRetries(o RunOptions, run func(parallel bool) (started, error)) (started, error) {
	parallel := o.AllowParallel
	for attempt := 0; ; attempt++ {
		result, err := run(parallel)
		if err == nil {
			return result, nil
		}
		var apiErr *platform.APIError
		if !errors.As(err, &apiErr) {
			return result, err
		}
		if apiErr.Status == http.StatusUnprocessableEntity && attempt < o.Retries {
			fmt.Printf("Experiment has validation errors (attempt %d/%d). Retrying in %ds...\n", attempt+1, o.Retries+1, o.RetryInterval)
			time.Sleep(time.Duration(o.RetryInterval) * time.Second)
			continue
		}
		if !parallel && apiErr.ProblemType() == anotherExperimentRunning {
			ok := o.Yes
			if !ok {
				if ok, err = prompt.Confirm("There is already an experiment running. Do you want to start it in parallel?", false, false); err != nil {
					return result, err
				}
			}
			if ok {
				parallel = true
				attempt--
				continue
			}
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
	// Steps asks the platform for the steps of the run, which reports need.
	Steps bool
}

// ErrTimedOut is returned when --timeout cancelled the run.
var ErrTimedOut = errors.New("timed out")

// wait polls the run until it ends. A run that did not complete is an error, which is
// what lets a pipeline fail on it. The finished run is returned for reports.
func wait(ctx context.Context, c *platform.Client, location string, o WaitOptions) (*RunResult, error) {
	path := location
	if i := strings.Index(location, "/api/"); i >= 0 {
		path = location[i:]
	}
	if o.Steps || o.ShowSteps {
		separator := "?"
		if strings.Contains(path, "?") {
			separator = "&"
		}
		path += separator + "fields=steps"
	}

	var runID int64
	cancel := func(why string) {
		if runID == 0 {
			return
		}
		fmt.Fprintf(os.Stderr, "%s, canceling experiment run %d.\n", why, runID)
		cancelCtx, done := context.WithTimeout(context.Background(), 30*time.Second)
		defer done()
		if _, _, err := platform.Read(c.CancelExperimentExecution(cancelCtx, runID)); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to cancel experiment run %d: %s\n", runID, err)
		}
	}
	if !o.KeepRunningOnInterrupt {
		pop := interrupt.Push(func(os.Signal) { cancel("Interrupted") })
		defer pop()
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
		runID = run.ID
		fmt.Println("Current run state:", strings.ToLower(run.State))
		if o.ShowSteps {
			for i, step := range run.Steps {
				id := fmt.Sprint(i)
				if shown[id] != step.State {
					shown[id] = step.State
					fmt.Printf("  step %d/%d %s: %s\n", i+1, len(run.Steps), step.Name, strings.ToLower(step.State))
				}
			}
		}
		if !terminal[run.State] {
			if !deadline.IsZero() && time.Now().After(deadline) {
				cancel(fmt.Sprintf("Experiment run %d did not end within %s", run.ID, o.Timeout))
				return run, fmt.Errorf("Experiment %s (#%d) did not end within %s and was canceled: %w", run.Key, run.ID, o.Timeout, ErrTimedOut)
			}
			continue
		}
		if run.State != "COMPLETED" {
			reason := ""
			if run.Reason != "" {
				reason = ", reason: " + run.Reason
			}
			return run, fmt.Errorf("Experiment %s (#%d) %s%s", run.Key, run.ID, strings.ToLower(run.State), reason)
		}
		return run, nil
	}
}
