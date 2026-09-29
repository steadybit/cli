// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/v6/internal/license"
	"github.com/steadybit/cli/v6/internal/platform"
)

func newLicense() *cobra.Command {
	cmd := &cobra.Command{Use: "license", Short: "Show the license of the tenant and what of it is used. Needs an admin access token."}

	var s license.ShowOptions
	show := &cobra.Command{
		Use:     "show",
		Short:   "Show the license, when it expires, and how much of each limit is used.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit license show", "steadybit license show --jq '.expires'"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return license.Show(ctx, c, s) }),
	}
	show.Flags().StringVarP(&s.Type, "type", "t", "", `Print the license summary as "json" or "yaml" instead of text.`)

	var r license.ReportOptions
	report := &cobra.Command{
		Use:     "report",
		Short:   "Download the license usage report, a zip archive with the usage of each license period.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit license report", "steadybit license report -o usage.zip"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return license.Report(ctx, c, r) }),
	}
	report.Flags().StringVarP(&r.Output, "output", "o", "", "Write the report to this file. (default: the name the platform gives it, in the current directory)")

	cmd.AddCommand(show, report)
	return cmd
}
