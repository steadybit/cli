// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package gitops

import (
	"context"
	"errors"
	"fmt"

	"github.com/steadybit/cli/v6/internal/experiment"
	"github.com/steadybit/cli/v6/internal/output"
	"github.com/steadybit/cli/v6/internal/platform"
	"github.com/steadybit/cli/v6/internal/resource"
)

// ErrDifferent ends `diff` with exit status 2 when files and platform disagree, so a
// pipeline can tell drift (2) from a failure (1).
var ErrDifferent = errors.New("files differ from the platform")

func compareAll(ctx context.Context, c *platform.Client, k Kind, paths []string, recursive bool) ([]Result, error) {
	files, err := experiment.ResolveFiles(paths, recursive)
	if err != nil {
		return nil, err
	}
	results := make([]Result, 0, len(files))
	for _, file := range files {
		doc, _, err := resource.Read(file, k.Name)
		if err != nil {
			return nil, err
		}
		r, err := Compare(ctx, c, k, file, doc)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

// DiffFiles prints how each file differs from the platform: a unified diff for changes,
// a line for files that would create something new.
func DiffFiles(ctx context.Context, c *platform.Client, k Kind, paths []string, recursive bool) error {
	results, err := compareAll(ctx, c, k, paths, recursive)
	if err != nil {
		return err
	}
	different := 0
	for _, r := range results {
		switch r.State {
		case Changed:
			different++
			fmt.Print(colorDiff(r.Diff))
		case New:
			different++
			fmt.Println(r.Describe(k))
		}
	}
	if different == 0 {
		fmt.Printf("%d %s file(s) match the platform.\n", len(results), k.Name)
		return nil
	}
	fmt.Printf("%d of %d %s file(s) differ from the platform.\n", different, len(results), k.Name)
	return ErrDifferent
}

// DryRun reports what applying the files would do, without applying anything.
func DryRun(ctx context.Context, c *platform.Client, k Kind, paths []string, recursive bool) error {
	results, err := compareAll(ctx, c, k, paths, recursive)
	if err != nil {
		return err
	}
	for _, r := range results {
		fmt.Println(r.Describe(k))
	}
	return nil
}

func colorDiff(diff string) string {
	if !output.ColorsEnabled() {
		return diff
	}
	var out []byte
	for _, line := range splitKeep(diff) {
		switch {
		case len(line) > 3 && (line[:3] == "---" || line[:3] == "+++"):
			out = append(out, output.Bold(line)...)
		case len(line) > 0 && line[0] == '+':
			out = append(out, output.Green(line)...)
		case len(line) > 0 && line[0] == '-':
			out = append(out, output.Red(line)...)
		default:
			out = append(out, line...)
		}
	}
	return string(out)
}

// splitKeep splits after each newline, keeping it, so colours never swallow one.
func splitKeep(s string) []string {
	var lines []string
	for len(s) > 0 {
		i := 0
		for i < len(s) && s[i] != '\n' {
			i++
		}
		if i < len(s) {
			lines = append(lines, s[:i]+"\n")
			s = s[i+1:]
		} else {
			lines = append(lines, s)
			s = ""
		}
	}
	return lines
}
