// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"context"
	"strings"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/report"
)

func newReport() *cobra.Command {
	cmd := &cobra.Command{Use: "report", Short: "Print reports on the tenant as time series. Output is written to file or stdout."}
	for _, kind := range report.Kinds {
		cmd.AddCommand(newReportKind(kind))
	}
	return cmd
}

func newReportKind(k report.Kind) *cobra.Command {
	var o report.Options
	example := []string{"steadybit report " + k.Name + " --from 2026-01-01 --to 2026-06-30 --rollup MONTHLY"}
	cmd := &cobra.Command{
		Use:   k.Name,
		Short: k.Short,
		Args:  cobra.NoArgs,
		RunE:  withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return report.Run(ctx, c, k, o) }),
	}
	flags := cmd.Flags()
	flags.StringVar(&o.From, "from", "", "The first day, e.g. 2026-09-01. (default: 30 days before --to)")
	flags.StringVar(&o.To, "to", "", "The last day, e.g. 2026-09-30. (default: today)")
	flags.StringVar(&o.Rollup, "rollup", "", `Sum up per "DAILY" or "MONTHLY" bucket.`)
	outputFlags(cmd, &o.File, &o.Type, "report")
	if len(k.GroupBy) > 0 {
		flags.StringVar(&o.GroupBy, "group-by", "", "Split the series by one of "+strings.Join(k.GroupBy, ", ")+".")
		example = append(example, "steadybit report "+k.Name+" --group-by "+k.GroupBy[1]+" --team ADM -t json")
	}
	var variadics []string
	if k.Takes(report.Teams) {
		flags.StringArrayVar(&o.Teams, "team", nil, "Only count these teams, by team key.")
		variadics = append(variadics, "team")
	}
	if k.Takes(report.Environments) {
		flags.StringArrayVar(&o.Environments, "environment", nil, "Only count these environments, by name.")
		variadics = append(variadics, "environment")
	}
	if k.Takes(report.Services) {
		flags.StringArrayVar(&o.Services, "service", nil, "Only count these services, by id.")
		variadics = append(variadics, "service")
	}
	if k.Takes(report.ServiceProperties) {
		flags.StringArrayVar(&o.ServiceProperties, "service-property", nil, "Only count services with this property, as KEY=VALUE. Repeat for more.")
	}
	if k.Takes(report.Categories) {
		flags.StringArrayVar(&o.Categories, "category", nil, "Only count these risk categories.")
		variadics = append(variadics, "category")
	}
	variadic(cmd, variadics...)
	cmd.Example = examples(example...)
	return cmd
}
