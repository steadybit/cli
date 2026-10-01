// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"
)

// notCovered names the operations no command calls, and why. Deprecated operations need
// no entry: the CLI uses the latest version of an API only. Anything else the platform
// offers has a command, or a line here; the coverage check fails otherwise, so a new
// endpoint is a decision rather than something nobody noticed.
var notCovered = map[string]string{
	"Health":    "a probe for the platform's own monitoring",
	"Liveness":  "a probe for the platform's own monitoring",
	"Readiness": "a probe for the platform's own monitoring",

	"GetExperimentBadge": "experiment badge prints its URL, which a README loads",
	"GetLinkedBadge":     "experiment badge prints its URL, which a README loads",
	"ForwardToPlatform":  "experiment badge prints its URL, which a README links to",

	"GetLandscapeViews":   "saved views of the Explore landscape, which only the UI shows",
	"GetLandscapeView":    "saved views of the Explore landscape, which only the UI shows",
	"CreateLandscapeView": "saved views of the Explore landscape, which only the UI shows",
	"UpdateLandscapeView": "saved views of the Explore landscape, which only the UI shows",
	"DeleteLandscapeView": "saved views of the Explore landscape, which only the UI shows",

	"ConnectionCheck":           "not covered yet: checking a hub's repository before saving it",
	"GetPreflightActionSummary": "not covered yet: listing the preflight actions extensions offer",
	"UpdateExecutionProperties": "not covered yet: execution property set and add change one property at a time",
}

type operation struct {
	Method, Path, ID, Summary string
	Deprecated                bool
}

var nonAlphanumeric = regexp.MustCompile(`[^A-Za-z0-9]+(.?)`)

// methodName is the name oapi-codegen gives the client method of an operation id:
// getAccessTokens_1 becomes GetAccessTokens1.
func methodName(operationID string) string {
	name := nonAlphanumeric.ReplaceAllStringFunc(operationID, func(m string) string {
		return strings.ToUpper(nonAlphanumeric.FindStringSubmatch(m)[1])
	})
	if name == "" {
		return name
	}
	return strings.ToUpper(name[:1]) + name[1:]
}

// readSpec reads a spec from a file or, for an http(s) URL, from the platform.
func readSpec(source string) ([]byte, error) {
	if !strings.HasPrefix(source, "http://") && !strings.HasPrefix(source, "https://") {
		return os.ReadFile(source)
	}
	client := &http.Client{Timeout: time.Minute}
	resp, err := client.Get(source)
	if err != nil {
		return nil, fmt.Errorf("fetching the platform spec from %s: %w", source, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetching the platform spec from %s failed with status %d", source, resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func operations(raw []byte) ([]operation, error) {
	var spec struct {
		Paths map[string]map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal(raw, &spec); err != nil {
		return nil, fmt.Errorf("not an OpenAPI document: %w", err)
	}
	var ops []operation
	for path, methods := range spec.Paths {
		for method, body := range methods {
			var op struct {
				OperationID string `json:"operationId"`
				Summary     string `json:"summary"`
				Deprecated  bool   `json:"deprecated"`
			}
			// Path items also hold parameters and the like, which are not operations.
			if json.Unmarshal(body, &op) != nil || op.OperationID == "" {
				continue
			}
			ops = append(ops, operation{strings.ToUpper(method), path, op.OperationID, op.Summary, op.Deprecated})
		}
	}
	slices.SortFunc(ops, func(a, b operation) int { return strings.Compare(a.Path+" "+a.Method, b.Path+" "+b.Method) })
	return ops, nil
}

// calledMethods collects the names of the methods the CLI's own code calls: everything
// under cmd and internal but tests, the fake platform and these tools.
func calledMethods(root string) (map[string]bool, error) {
	called := map[string]bool{}
	fset := token.NewFileSet()
	for _, dir := range []string{"cmd", "internal"} {
		err := filepath.WalkDir(filepath.Join(root, dir), func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			rel, _ := filepath.Rel(root, path)
			rel = filepath.ToSlash(rel)
			if d.IsDir() {
				if rel == "internal/tools" || rel == "internal/platformtest" || d.Name() == "testdata" {
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			file, err := parser.ParseFile(fset, path, nil, 0)
			if err != nil {
				return err
			}
			ast.Inspect(file, func(n ast.Node) bool {
				if call, ok := n.(*ast.CallExpr); ok {
					if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
						called[strings.TrimSuffix(sel.Sel.Name, "WithBody")] = true
					}
				}
				return true
			})
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return called, nil
}

type coverageReport struct {
	Total, Called int
	// Uncovered are operations that no command calls and that notCovered does not name.
	Uncovered []operation
	// Stale are entries of notCovered for operations that are gone or now called.
	Stale []string
}

func checkCoverage(raw []byte, root string) (coverageReport, error) {
	ops, err := operations(raw)
	if err != nil {
		return coverageReport{}, err
	}
	called, err := calledMethods(root)
	if err != nil {
		return coverageReport{}, err
	}
	var report coverageReport
	listed := map[string]bool{}
	for _, op := range ops {
		name := methodName(op.ID)
		switch {
		case called[name]:
			report.Called++
		case op.Deprecated:
		case notCovered[name] != "":
			listed[name] = true
		default:
			report.Uncovered = append(report.Uncovered, op)
		}
		report.Total++
	}
	for name := range notCovered {
		if !listed[name] {
			report.Stale = append(report.Stale, name)
		}
	}
	slices.Sort(report.Stale)
	return report, nil
}

// coverage prints the operations of the spec at source that no command covers, as a
// Markdown list an issue can hold, and fails when there are any.
func coverage(source string) error {
	raw, err := readSpec(source)
	if err != nil {
		return err
	}
	report, err := checkCoverage(raw, ".")
	if err != nil {
		return err
	}
	fmt.Printf("%d of %d operations in %s are called by a command.\n", report.Called, report.Total, source)
	for _, op := range report.Uncovered {
		summary := ""
		if op.Summary != "" {
			summary = ": " + op.Summary
		}
		fmt.Printf("- `%s %s` (`%s`)%s\n", op.Method, op.Path, op.ID, summary)
	}
	for _, name := range report.Stale {
		fmt.Printf("- `%s` is listed as not covered in internal/tools/spec/coverage.go, but is called or gone: remove it there.\n", name)
	}
	if len(report.Uncovered) > 0 || len(report.Stale) > 0 {
		return fmt.Errorf("%d operation(s) are neither called by a command nor listed in internal/tools/spec/coverage.go, and %d listed one(s) are stale", len(report.Uncovered), len(report.Stale))
	}
	return nil
}
