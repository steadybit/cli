// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"context"
	"errors"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/internal/gitops"
	"github.com/steadybit/cli/internal/platform"
)

// newDiff is the `diff` command of a kind of file: `experiment diff`, `schedule diff`...
func newDiff(k gitops.Kind, group, example string) *cobra.Command {
	var files []string
	var recursive bool
	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Show how " + k.Name + " files differ from the platform. Exits with 2 when they do, so a pipeline can detect drift.",
		Args:  cobra.NoArgs,
		Example: examples(
			"steadybit "+group+" diff -f "+example,
			"steadybit "+group+" diff -f ./"+group+"s -R || echo drift",
		),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return gitops.DiffFiles(ctx, c, k, files, recursive)
		}),
	}
	cmd.Flags().StringArrayVarP(&files, "file", "f", nil, "The path to the file or a directory containing multiple files.")
	cmd.Flags().BoolVarP(&recursive, "recursive", "R", false, "Process the directory used in -f, --file recursively.")
	_ = cmd.MarkFlagRequired("file")
	variadic(cmd, "file")
	return cmd
}

// dryRun adds --dry-run to an apply command: report what would change, change nothing.
func dryRun(cmd *cobra.Command, k gitops.Kind, files *[]string, recursive *bool) {
	var enabled bool
	cmd.Flags().BoolVar(&enabled, "dry-run", false, "Report what applying the files would create or update, without changing anything.")
	run := cmd.RunE
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if !enabled {
			return run(cmd, args)
		}
		if len(*files) == 0 {
			return errors.New("--dry-run compares files with the platform; pass them with -f.")
		}
		return withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return gitops.DryRun(ctx, c, k, *files, *recursive)
		})(cmd, args)
	}
}
