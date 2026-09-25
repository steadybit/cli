// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package experiment

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"time"
)

// RunResult is a finished (or abandoned) experiment run, as reports describe it.
type RunResult struct {
	ID         int64     `json:"id"`
	Key        string    `json:"key"`
	Name       string    `json:"name"`
	State      string    `json:"state"`
	Reason     string    `json:"reason,omitempty"`
	Started    time.Time `json:"started"`
	Ended      time.Time `json:"ended"`
	UILocation string    `json:"uiLocation,omitempty"`
	Steps      []Step    `json:"steps"`
}

type Step struct {
	Name    string    `json:"name"`
	State   string    `json:"state"`
	Reason  string    `json:"reason,omitempty"`
	Started time.Time `json:"started"`
	Ended   time.Time `json:"ended"`
}

func (r RunResult) Duration() time.Duration { return duration(r.Started, r.Ended) }
func (s Step) Duration() time.Duration      { return duration(s.Started, s.Ended) }

func duration(from, to time.Time) time.Duration {
	if from.IsZero() || to.IsZero() || to.Before(from) {
		return 0
	}
	return to.Sub(from)
}

func parseRun(body []byte) (*RunResult, error) {
	var raw struct {
		ID      int64     `json:"id"`
		Key     string    `json:"key"`
		Name    string    `json:"name"`
		State   string    `json:"state"`
		Reason  string    `json:"reason"`
		Started time.Time `json:"started"`
		Ended   time.Time `json:"ended"`
		Steps   []struct {
			StepType    string         `json:"stepType"`
			ActionID    string         `json:"actionId"`
			CustomLabel string         `json:"customLabel"`
			State       string         `json:"state"`
			Reason      string         `json:"reason"`
			Started     time.Time      `json:"started"`
			Ended       time.Time      `json:"ended"`
			Parameters  map[string]any `json:"parameters"`
		} `json:"steps"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	run := &RunResult{ID: raw.ID, Key: raw.Key, Name: raw.Name, State: raw.State, Reason: raw.Reason, Started: raw.Started, Ended: raw.Ended}
	for _, s := range raw.Steps {
		name := s.CustomLabel
		switch {
		case name != "":
		case s.ActionID != "":
			name = s.ActionID
		case strings.EqualFold(s.StepType, "wait") && s.Parameters["duration"] != nil:
			name = fmt.Sprintf("wait %v", s.Parameters["duration"])
		default:
			name = strings.ToLower(s.StepType)
		}
		run.Steps = append(run.Steps, Step{Name: name, State: s.State, Reason: s.Reason, Started: s.Started, Ended: s.Ended})
	}
	return run, nil
}

// WriteReport writes the runs as JUnit XML, which CI systems show as test results, or
// as JSON, chosen by the file's extension.
func WriteReport(file string, runs []*RunResult) error {
	var content []byte
	var err error
	if strings.HasSuffix(strings.ToLower(file), ".json") {
		content, err = json.MarshalIndent(runs, "", "  ")
	} else {
		content, err = junit(runs)
	}
	if err != nil {
		return err
	}
	return os.WriteFile(file, append(content, '\n'), 0o644)
}

type junitSuites struct {
	XMLName  xml.Name     `xml:"testsuites"`
	Name     string       `xml:"name,attr"`
	Tests    int          `xml:"tests,attr"`
	Failures int          `xml:"failures,attr"`
	Errors   int          `xml:"errors,attr"`
	Time     string       `xml:"time,attr"`
	Suites   []junitSuite `xml:"testsuite"`
}

type junitSuite struct {
	Name       string          `xml:"name,attr"`
	Tests      int             `xml:"tests,attr"`
	Failures   int             `xml:"failures,attr"`
	Errors     int             `xml:"errors,attr"`
	Skipped    int             `xml:"skipped,attr"`
	Time       string          `xml:"time,attr"`
	Timestamp  string          `xml:"timestamp,attr,omitempty"`
	Properties []junitProperty `xml:"properties>property,omitempty"`
	Cases      []junitCase     `xml:"testcase"`
}

type junitProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type junitCase struct {
	ClassName string        `xml:"classname,attr"`
	Name      string        `xml:"name,attr"`
	Time      string        `xml:"time,attr"`
	Failure   *junitProblem `xml:"failure,omitempty"`
	Error     *junitProblem `xml:"error,omitempty"`
	Skipped   *struct{}     `xml:"skipped,omitempty"`
}

type junitProblem struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Text    string `xml:",chardata"`
}

func seconds(d time.Duration) string { return fmt.Sprintf("%.3f", d.Seconds()) }

// A failed step is a failure (the hypothesis did not hold), an errored one an error
// (the run could not do what it should); steps that never ran are skipped.
func junitCaseFor(className, name, state, reason string, d time.Duration) junitCase {
	c := junitCase{ClassName: className, Name: name, Time: seconds(d)}
	message := strings.ToLower(state)
	if reason != "" {
		message += ": " + reason
	}
	switch state {
	case "FAILED":
		c.Failure = &junitProblem{Message: message, Type: state, Text: reason}
	case "ERRORED":
		c.Error = &junitProblem{Message: message, Type: state, Text: reason}
	case "CANCELED", "SKIPPED", "CREATED", "PREPARED", "", "COMPLETED":
		if state != "COMPLETED" {
			c.Skipped = &struct{}{}
		}
	default:
		// Still going when the run was cut short, by a timeout for instance.
		c.Error = &junitProblem{Message: "did not end: " + message, Type: state, Text: reason}
	}
	return c
}

func junit(runs []*RunResult) ([]byte, error) {
	suites := junitSuites{Name: "steadybit"}
	var total time.Duration
	for _, run := range runs {
		suite := junitSuite{Name: strings.TrimSpace(run.Key + " " + run.Name), Time: seconds(run.Duration())}
		if !run.Started.IsZero() {
			suite.Timestamp = run.Started.UTC().Format(time.RFC3339)
		}
		suite.Properties = append(suite.Properties, junitProperty{Name: "executionId", Value: fmt.Sprint(run.ID)})
		if run.UILocation != "" {
			suite.Properties = append(suite.Properties, junitProperty{Name: "uiLocation", Value: run.UILocation})
		}
		// A run that has no steps to report still has to show up, as one case.
		if len(run.Steps) == 0 {
			suite.Cases = append(suite.Cases, junitCaseFor(run.Key, run.Key, run.State, run.Reason, run.Duration()))
		}
		for i, step := range run.Steps {
			suite.Cases = append(suite.Cases, junitCaseFor(run.Key, fmt.Sprintf("%d. %s", i+1, step.Name), step.State, step.Reason, step.Duration()))
		}
		for _, c := range suite.Cases {
			suite.Tests++
			switch {
			case c.Failure != nil:
				suite.Failures++
			case c.Error != nil:
				suite.Errors++
			case c.Skipped != nil:
				suite.Skipped++
			}
		}
		// A run that ended badly without any step to blame, a canceled or timed-out run
		// for instance, still fails its suite.
		if run.State != "COMPLETED" && suite.Failures == 0 && suite.Errors == 0 {
			message := strings.ToLower(run.State)
			if run.Reason != "" {
				message += ": " + run.Reason
			}
			problem := &junitProblem{Message: message, Type: run.State, Text: run.Reason}
			c := junitCase{ClassName: run.Key, Name: "run", Time: seconds(run.Duration())}
			if run.State == "FAILED" {
				c.Failure = problem
				suite.Failures++
			} else {
				c.Error = problem
				suite.Errors++
			}
			suite.Cases = append(suite.Cases, c)
			suite.Tests++
		}
		suites.Tests += suite.Tests
		suites.Failures += suite.Failures
		suites.Errors += suite.Errors
		total += run.Duration()
		suites.Suites = append(suites.Suites, suite)
	}
	suites.Time = seconds(total)
	out, err := xml.MarshalIndent(suites, "", "  ")
	return append([]byte(xml.Header), out...), err
}

// WriteGitHubSummary appends a Markdown summary of the runs to the job summary when the
// CLI runs in GitHub Actions, which names the file in GITHUB_STEP_SUMMARY.
func WriteGitHubSummary(runs []*RunResult) error {
	file := os.Getenv("GITHUB_STEP_SUMMARY")
	if file == "" || len(runs) == 0 {
		return nil
	}
	var b strings.Builder
	for _, run := range runs {
		icon := "✅"
		if run.State != "COMPLETED" {
			icon = "❌"
		}
		title := run.Key
		if run.Name != "" {
			title += " · " + run.Name
		}
		fmt.Fprintf(&b, "### %s Steadybit experiment %s\n\n", icon, title)
		link := fmt.Sprintf("#%d", run.ID)
		if run.UILocation != "" {
			link = fmt.Sprintf("[#%d](%s)", run.ID, run.UILocation)
		}
		fmt.Fprintf(&b, "Run %s %s after %s", link, strings.ToLower(run.State), run.Duration().Round(time.Second))
		if run.Reason != "" {
			fmt.Fprintf(&b, ": %s", run.Reason)
		}
		b.WriteString("\n\n")
		if len(run.Steps) > 0 {
			b.WriteString("| # | Step | State | Duration |\n|---|---|---|---|\n")
			for i, step := range run.Steps {
				state := strings.ToLower(step.State)
				if step.Reason != "" {
					state += ": " + step.Reason
				}
				fmt.Fprintf(&b, "| %d | %s | %s | %s |\n", i+1, escapeTable(step.Name), escapeTable(state), step.Duration().Round(time.Second))
			}
			b.WriteString("\n")
		}
	}
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(b.String())
	return err
}

func escapeTable(s string) string {
	return strings.NewReplacer("|", `\|`, "\n", " ").Replace(s)
}
