// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/v6/internal/gitops"
	"github.com/steadybit/cli/v6/internal/platform"
	"github.com/steadybit/cli/v6/internal/resource"
	"github.com/steadybit/cli/v6/internal/schedule"
)

const scheduleID = "01951394-727f-76a0-8675-c7519ebd0ff5"

func scheduleIDFlag(cmd *cobra.Command, id *string) {
	cmd.Flags().StringVarP(id, "id", "i", "", "The experiment schedule id.")
	_ = cmd.MarkFlagRequired("id")
}

// scheduleFields adds the flags `create` and `update` share. --allow-parallel and
// --no-allow-parallel stay unset unless one is given, so an update leaves it alone.
func scheduleFields(cmd *cobra.Command, f *schedule.Fields) {
	vars := newKeyValues()
	f.Variables = vars.values
	var allow, noAllow bool
	flags := cmd.Flags()
	flags.StringVar(&f.Cron, "cron", "", `Run repeatedly on this Quartz cron expression (seconds first), e.g. "0 0 9 ? * MON-FRI".`)
	flags.StringVar(&f.StartAt, "start-at", "", "Run once at this ISO 8601 time, e.g. 2026-10-01T09:00:00Z.")
	flags.StringVar(&f.Timezone, "timezone", "", "The timezone of the cron expression, e.g. Europe/Berlin.")
	flags.BoolVar(&allow, "allow-parallel", false, "Run even when another experiment is running.")
	flags.BoolVar(&noAllow, "no-allow-parallel", false, "Skip the run when another experiment is running.")
	flags.Var(vars, "variable", "A variable for the scheduled runs, overriding experiment and environment variables. Repeat for more.")
	cmd.MarkFlagsMutuallyExclusive("allow-parallel", "no-allow-parallel")
	cmd.PreRun = func(*cobra.Command, []string) {
		switch {
		case allow:
			f.AllowParallel = &allow
		case noAllow:
			no := false
			f.AllowParallel = &no
		}
	}
}

func newSchedule() *cobra.Command {
	cmd := &cobra.Command{Use: "schedule", Short: "Schedule experiments."}

	var l schedule.ListOptions
	list := &cobra.Command{
		Use:     "list",
		Short:   "List experiment schedules.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit schedule list", "steadybit schedule list --team ADM --experiment ADM-1 ADM-2", "steadybit schedule list --jq '.[] | select(.enabled) | .experimentKey'"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return schedule.List(ctx, c, l) }),
	}
	list.Flags().StringArrayVar(&l.Teams, "team", nil, "Only list schedules of these teams, by team key.")
	list.Flags().StringArrayVar(&l.Experiments, "experiment", nil, "Only list schedules of these experiments, by experiment key.")
	list.Flags().StringVarP(&l.Type, "type", "t", "", resource.ListTypeHelp)
	variadic(list, "team", "experiment")

	var g schedule.GetOptions
	get := &cobra.Command{
		Use:     "get",
		Short:   "Get an experiment schedule. Output is written to file or stdout.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit schedule get -i " + scheduleID + " -f schedule.yml"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return schedule.Get(ctx, c, g) }),
	}
	scheduleIDFlag(get, &g.ID)
	get.Flags().StringVarP(&g.File, "file", "f", "", "The path to write the schedule to.")
	get.Flags().StringVarP(&g.Type, "type", "t", "", typeHelp)

	var a schedule.ApplyOptions
	apply := &cobra.Command{
		Use:     "apply",
		Short:   "Create or update experiment schedules from files. A file without an id creates a schedule, and the new id is written back to it.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit schedule apply -f schedule.yml", "steadybit schedule apply -f ./schedules -R"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return schedule.Apply(ctx, c, a) }),
	}
	apply.Flags().StringArrayVarP(&a.Files, "file", "f", nil, "The path to the schedule file or a directory containing multiple files.")
	apply.Flags().BoolVarP(&a.Recursive, "recursive", "R", false, "Process the directory used in -f, --file recursively.")
	_ = apply.MarkFlagRequired("file")
	variadic(apply, "file")
	dryRun(apply, gitops.Schedule, &a.Files, &a.Recursive)

	var cr schedule.CreateOptions
	create := &cobra.Command{
		Use:   "create",
		Short: "Schedule an experiment, either repeatedly with --cron or once with --start-at.",
		Args:  cobra.NoArgs,
		Example: examples(
			`steadybit schedule create -k ADM-1 --cron "0 0 9 ? * MON-FRI" --timezone Europe/Berlin`,
			"steadybit schedule create -k ADM-1 --start-at 2026-10-01T09:00:00Z --no-allow-parallel",
		),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return schedule.Create(ctx, c, cr) }),
	}
	create.Flags().StringVarP(&cr.Experiment, "experiment", "k", "", "The key of the experiment to schedule.")
	create.Flags().BoolVar(&cr.Disabled, "disabled", false, "Create the schedule disabled.")
	_ = create.MarkFlagRequired("experiment")
	scheduleFields(create, &cr.Fields)

	var u schedule.UpdateOptions
	update := &cobra.Command{
		Use:     "update",
		Short:   "Change an experiment schedule. Only the given fields are changed.",
		Args:    cobra.NoArgs,
		Example: examples(`steadybit schedule update -i ` + scheduleID + ` --cron "0 30 8 ? * *"`),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return schedule.Update(ctx, c, u) }),
	}
	scheduleIDFlag(update, &u.ID)
	scheduleFields(update, &u.Fields)

	idCommand := func(use, short string, run func(ctx context.Context, c *platform.Client, id string) error) *cobra.Command {
		var id string
		c := &cobra.Command{
			Use:     use,
			Short:   short,
			Args:    cobra.NoArgs,
			Example: examples("steadybit schedule " + use + " -i " + scheduleID),
			RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return run(ctx, c, id) }),
		}
		scheduleIDFlag(c, &id)
		return c
	}
	cmd.AddCommand(list, get, apply, newDiff(gitops.Schedule, "schedule", "schedule.yml"), create, update,
		idCommand("enable", "Enable an experiment schedule.", func(ctx context.Context, c *platform.Client, id string) error {
			return schedule.SetEnabled(ctx, c, id, true)
		}),
		idCommand("disable", "Disable an experiment schedule without deleting it.", func(ctx context.Context, c *platform.Client, id string) error {
			return schedule.SetEnabled(ctx, c, id, false)
		}),
		idCommand("delete", "Delete an experiment schedule.", schedule.Delete),
	)
	return cmd
}
