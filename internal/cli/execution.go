// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/v6/internal/execution"
	"github.com/steadybit/cli/v6/internal/platform"
	"github.com/steadybit/cli/v6/internal/resource"
)

func runID(cmd *cobra.Command, id *int64) {
	cmd.Flags().Int64VarP(id, "id", "i", 0, "The experiment run id.")
	_ = cmd.MarkFlagRequired("id")
}

func newExecution() *cobra.Command {
	cmd := &cobra.Command{Use: "execution", Short: "Search, inspect, cancel and annotate experiment runs, and download their artifacts."}

	var g execution.GetOptions
	get := &cobra.Command{
		Use:     "get",
		Short:   "Get an experiment run, including its steps and target executions. Output is written to file or stdout.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit execution get -i 1234", "steadybit execution get -i 1234 -t json | jq .state"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return execution.Get(ctx, c, g) }),
	}
	runID(get, &g.ID)
	get.Flags().StringVarP(&g.File, "file", "f", "", "The path to write the experiment run to.")
	get.Flags().StringVarP(&g.Type, "type", "t", "", typeHelp)

	var cancelID int64
	cancel := &cobra.Command{
		Use:     "cancel",
		Short:   "Cancel a running experiment run. The run stops as soon as its agents have been told.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit execution cancel -i 1234"),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return execution.Cancel(ctx, c, cancelID)
		}),
	}
	runID(cancel, &cancelID)

	property := &cobra.Command{Use: "property", Short: "Change the properties of an experiment run."}
	var s execution.PropertyOptions
	set := &cobra.Command{
		Use:   "set",
		Short: "Set the value of a property of an experiment run. Only properties editable in a run can be changed. Several --value set a list property.",
		Args:  cobra.NoArgs,
		Example: examples(
			`steadybit execution property set -i 1234 -k approvedBy --value "Jane Doe"`,
			"steadybit execution property set -i 1234 -k tickets --value SHOP-1 SHOP-2",
			"steadybit execution property set -i 1234 -k score --value 7 --json",
		),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return execution.SetProperty(ctx, c, s)
		}),
	}
	runID(set, &s.ID)
	set.Flags().StringVarP(&s.Key, "key", "k", "", "The property key.")
	set.Flags().StringArrayVar(&s.Values, "value", nil, "The value to set.")
	set.Flags().BoolVar(&s.JSON, "json", false, "Parse each value as JSON, to send a number or an object.")
	_ = set.MarkFlagRequired("key")
	_ = set.MarkFlagRequired("value")
	variadic(set, "value")

	var a execution.PropertyOptions
	var addValue string
	add := &cobra.Command{
		Use:     "add",
		Short:   "Add a value to a list property of an experiment run.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit execution property add -i 1234 -k tickets --value SHOP-3"),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			a.Values = []string{addValue}
			return execution.AddProperty(ctx, c, a)
		}),
	}
	runID(add, &a.ID)
	add.Flags().StringVarP(&a.Key, "key", "k", "", "The property key.")
	add.Flags().StringVar(&addValue, "value", "", "The value to add.")
	add.Flags().BoolVar(&a.JSON, "json", false, "Parse the value as JSON, to send a number or an object.")
	_ = add.MarkFlagRequired("key")
	_ = add.MarkFlagRequired("value")
	property.AddCommand(set, add)

	artifact := &cobra.Command{Use: "artifact", Short: "List and download the artifacts of an experiment run."}
	var listID int64
	var listType string
	list := &cobra.Command{
		Use:     "list",
		Short:   "List the artifacts that the actions of an experiment run attached.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit execution artifact list -i 1234", "steadybit execution artifact list -i 1234 --jq '.[].artifactId'"),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return execution.ListArtifacts(ctx, c, listID, listType)
		}),
	}
	runID(list, &listID)
	list.Flags().StringVarP(&listType, "type", "t", "", resource.ListTypeHelp)
	var d execution.DownloadOptions
	download := &cobra.Command{
		Use:   "download",
		Short: "Download the artifacts of an experiment run into <directory>/<target execution>/<artifact>. Without filters, all of them are downloaded.",
		Args:  cobra.NoArgs,
		Example: examples(
			"steadybit execution artifact download -i 1234 -d ./artifacts",
			"steadybit execution artifact download -i 1234 -a jmeter-report.zip -o report.zip",
		),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return execution.Download(ctx, c, d) }),
	}
	runID(download, &d.ID)
	download.Flags().StringVarP(&d.Artifact, "artifact", "a", "", "Only download artifacts with this id, usually the file name.")
	download.Flags().StringVar(&d.TargetExecution, "target-execution", "", "Only download artifacts of this target execution.")
	download.Flags().StringVarP(&d.Directory, "directory", "d", ".", "The directory to download into.")
	download.Flags().StringVarP(&d.Output, "output", "o", "", "Write the artifact to this file instead. Requires exactly one match.")
	download.MarkFlagsMutuallyExclusive("output", "directory")
	artifact.AddCommand(list, download)

	var w execution.WatchOptions
	watch := &cobra.Command{
		Use:   "watch",
		Short: "Follow an experiment run live until it ends: its steps, their targets and timings. Stopping the watch leaves the run alone.",
		Args:  cobra.NoArgs,
		Example: examples(
			"steadybit execution watch -i 1234",
			"steadybit execution watch -k ADM-1",
		),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			if (w.ID == 0) == (w.Key == "") {
				return errors.New("Pass either --id or --key.")
			}
			return execution.Watch(ctx, c, w)
		}),
	}
	watch.Flags().Int64VarP(&w.ID, "id", "i", 0, "The experiment run id.")
	watch.Flags().StringVarP(&w.Key, "key", "k", "", "Watch the latest run of this experiment instead.")
	watch.Flags().DurationVar(&w.Interval, "interval", 2*time.Second, "How often to refresh.")

	var l execution.ListOptions
	runs := &cobra.Command{
		Use:   "list",
		Short: "Search the runs of all experiments, most recent first. Filters of the same kind match any of the given values.",
		Args:  cobra.NoArgs,
		Example: examples(
			"steadybit execution list --team ADM",
			"steadybit execution list --team ADM --state FAILED ERRORED --from 2026-09-28 --fail-on-match",
			"steadybit execution list -k ADM-1 --limit 0 --jq '.[] | select(.state == \"COMPLETED\") | .id'",
		),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return execution.List(ctx, c, l) }),
	}
	runs.Flags().StringArrayVar(&l.Teams, "team", nil, "Only list runs of experiments of these teams, by team key.")
	runs.Flags().StringArrayVarP(&l.Experiments, "key", "k", nil, "Only list runs of these experiments, by experiment key.")
	runs.Flags().StringArrayVar(&l.States, "state", nil, "Only list runs in these states: "+strings.Join(execution.States, ", ")+".")
	runs.Flags().StringArrayVar(&l.Environments, "environment", nil, "Only list runs in these environments, by name.")
	runs.Flags().StringArrayVar(&l.Services, "service", nil, "Only list runs of experiments linked to these services, by name.")
	runs.Flags().StringVar(&l.Name, "name", "", "Only list runs of experiments whose name or key contains this text.")
	runs.Flags().StringVar(&l.From, "from", "", "Only list runs created at or after this time, a date or an RFC 3339 time. A date means the start of that day, UTC.")
	runs.Flags().StringVar(&l.To, "to", "", "Only list runs created at or before this time, a date or an RFC 3339 time. A date means the end of that day, UTC.")
	runs.Flags().IntVar(&l.Limit, "limit", 50, "The most runs to list, 0 for all of them.")
	runs.Flags().BoolVar(&l.FailOnMatch, "fail-on-match", false, "Exit with a non-zero status when any run matches, to stop a pipeline.")
	runs.Flags().StringVarP(&l.Type, "type", "t", "", resource.ListTypeHelp)
	variadic(runs, "team", "key", "state", "environment", "service")
	_ = runs.RegisterFlagCompletionFunc("state", cobra.FixedCompletions(execution.States, cobra.ShellCompDirectiveNoFileComp))

	cmd.AddCommand(runs, get, cancel, property, artifact, watch)
	return cmd
}
