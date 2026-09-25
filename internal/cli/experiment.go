// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/internal/experiment"
	"github.com/steadybit/cli/internal/platform"
)

func newExperiment() *cobra.Command {
	cmd := &cobra.Command{Use: "experiment", Short: "Check and run experiments."}
	cmd.AddCommand(newExperimentRun(), newExperimentGet(), newExperimentApply(), newExperimentDump())
	return cmd
}

// keyValues is a repeatable KEY=VALUE flag. Only the first `=` separates.
type keyValues map[string]string

func (k *keyValues) String() string { return "" }
func (k *keyValues) Type() string   { return "KEY=VALUE" }
func (k *keyValues) Set(value string) error {
	i := strings.Index(value, "=")
	if i <= 0 {
		return fmt.Errorf("'%s' is not in the form KEY=VALUE", value)
	}
	if *k == nil {
		*k = keyValues{}
	}
	(*k)[value[:i]] = value[i+1:]
	return nil
}

func newExperimentRun() *cobra.Command {
	var o experiment.RunOptions
	var noWait bool
	placeholders := keyValues{}
	cmd := &cobra.Command{
		Use:     "run",
		Aliases: []string{"exec"},
		Short:   "Executes an experiment run. If a file is specified the experiment is saved before execution.",
		Args:    cobra.NoArgs,
		Example: examples(
			"steadybit experiment run -k ADM-1",
			"steadybit experiment run -f experiment.yml --no-wait",
			"steadybit experiment run --template d7e65100-1d20-4980-be87-c351704910b8 --team ADM -p CLUSTER=prod",
		),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			o.Wait = !noWait
			o.Placeholder = placeholders
			return experiment.Run(ctx, c, o)
		}),
	}
	f := cmd.Flags()
	f.StringVarP(&o.Key, "key", "k", "", "The experiment key.")
	f.StringArrayVarP(&o.Files, "file", "f", nil, "The path to the experiment file or a directory containing multiple files.")
	f.BoolVarP(&o.Recursive, "recursive", "R", false, "Process the directory used in -f, --file recursively.")
	f.BoolVar(&noWait, "no-wait", false, "Do not wait for experiment run to finish.")
	f.BoolVar(&o.Yes, "yes", false, "Skip the prompt asking for experiment run confirmation. Not necessary when no TTY is attached.")
	f.BoolVar(&o.AllowParallel, "allowParallel", false, "Skip the prompt warning about another experiment running and allow always parallel execution.")
	f.IntVar(&o.Retries, "retries", 0, "Number of retries when the experiment fails validation (e.g., missing targets). 0 means no retry.")
	f.IntVar(&o.RetryInterval, "retryInterval", 10, "Interval in seconds between retries.")
	f.StringVar(&o.Template, "template", "", "Create the experiment from the experiment template with this id.")
	f.StringVar(&o.Team, "team", "", "With --template: the key of the team owning the experiment.")
	f.StringVar(&o.Environment, "environment", "", "With --template: the environment the experiment runs in.")
	f.StringVar(&o.ExternalID, "external-id", "", "With --template: an identifier of your own; reusing it updates the experiment.")
	f.VarP(&placeholders, "placeholder", "p", "With --template: a placeholder value. Repeat for more.")
	cmd.MarkFlagsMutuallyExclusive("key", "file")
	cmd.MarkFlagsMutuallyExclusive("template", "file")
	variadic(cmd, "file")
	return cmd
}

func newExperimentGet() *cobra.Command {
	var o experiment.GetOptions
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get an experiment from Steadybit. Output is written to file or stdout.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit experiment get -k ADM-1", "steadybit experiment get -k ADM-1 -f experiment.json"),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return experiment.Get(ctx, c, o)
		}),
	}
	cmd.Flags().StringVarP(&o.Key, "key", "k", "", "The experiment key.")
	cmd.Flags().StringVarP(&o.File, "file", "f", "", "The path to the experiment file.")
	cmd.Flags().StringVarP(&o.Type, "type", "t", "", `The output format of the experiment ("json" or "yaml"). (default: "json" if the file ends in ".json", "yaml" otherwise.)`)
	_ = cmd.MarkFlagRequired("key")
	return cmd
}

func newExperimentApply() *cobra.Command {
	var o experiment.ApplyOptions
	cmd := &cobra.Command{
		Use:     "apply",
		Short:   "Upload an experiment to Steadybit. If a key is provided, an update is performed. Otherwise, the externalId from the file is used to create or update the experiment.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit experiment apply -f experiment.yml", "steadybit experiment apply -f ./experiments -R"),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return experiment.Apply(ctx, c, o)
		}),
	}
	cmd.Flags().StringVarP(&o.Key, "key", "k", "", "The experiment key.")
	cmd.Flags().StringArrayVarP(&o.Files, "file", "f", nil, "The path to the experiment file or a directory containing multiple files.")
	cmd.Flags().BoolVarP(&o.Recursive, "recursive", "R", false, "Process the directory used in -f, --file recursively.")
	_ = cmd.MarkFlagRequired("file")
	variadic(cmd, "file")
	return cmd
}

func newExperimentDump() *cobra.Command {
	var o experiment.DumpOptions
	cmd := &cobra.Command{
		Use:   "dump",
		Short: "Dump all experiments and executions from all teams in Steadybit.",
		Args:  cobra.NoArgs,
		Example: examples(
			"steadybit experiment dump -d ./dump",
			"steadybit experiment dump -d ./dump -t json --team ADM WEBHOOK",
		),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return experiment.Dump(ctx, c, o)
		}),
	}
	cmd.Flags().StringVarP(&o.Directory, "directory", "d", ".", "The path to dump all the experiments to")
	cmd.Flags().StringVarP(&o.Type, "type", "t", "yaml", `The output format of the experiment ("json" or "yaml").`)
	cmd.Flags().StringArrayVar(&o.Teams, "team", nil, "Only dump the given teams, by team key. Defaults to every accessible team.")
	variadic(cmd, "team")
	return cmd
}
