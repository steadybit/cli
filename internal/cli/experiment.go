// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/internal/experiment"
	"github.com/steadybit/cli/internal/jsyaml"
	"github.com/steadybit/cli/internal/platform"
)

func newExperiment() *cobra.Command {
	cmd := &cobra.Command{Use: "experiment", Short: "Check and run experiments."}
	cmd.AddCommand(newExperimentRun(), newExperimentGet(), newExperimentApply(), newExperimentDump())
	return cmd
}

// keyValues is a repeatable KEY=VALUE flag, kept in the order given. Only the first `=`
// separates, so a value may contain one, as a URL with a query string does.
type keyValues struct{ values *jsyaml.Map }

func newKeyValues() *keyValues { return &keyValues{values: jsyaml.NewMap()} }

func (k *keyValues) String() string { return "" }
func (k *keyValues) Type() string   { return "KEY=VALUE" }
func (k *keyValues) Set(value string) error {
	i := strings.Index(value, "=")
	if i <= 0 {
		return fmt.Errorf("'%s' is not in the form KEY=VALUE.", value)
	}
	k.values.Set(value[:i], value[i+1:])
	return nil
}

// addTemplateFlags adds what `run` and `apply` share for creating from a template.
func addTemplateFlags(cmd *cobra.Command, o *experiment.TemplateOptions) {
	placeholders, vars := newKeyValues(), newKeyValues()
	o.Placeholder, o.Variable = placeholders.values, vars.values
	f := cmd.Flags()
	f.StringVar(&o.Template, "template", "", "Create the experiment from the experiment template with this id.")
	f.StringVar(&o.Team, "team", "", "With --template: the key of the team owning the experiment.")
	f.StringVar(&o.Environment, "environment", "", "With --template: the environment the experiment runs in.")
	f.StringVar(&o.ExternalID, "external-id", "", "With --template: an identifier of your own. Using the same one again updates the experiment it created before.")
	f.VarP(placeholders, "placeholder", "p", "With --template: a placeholder value. Repeat for more.")
	f.StringVar(&o.PlaceholdersFile, "placeholders", "", "With --template: a YAML/JSON file mapping placeholder keys to values. -p overrides entries.")
	f.Var(vars, "variable", "With --template: an experiment variable to add to the experiment. Repeat for more.")
	var noReset bool
	f.BoolVar(&noReset, "no-reset-properties", false, "With --template: keep the properties of an existing experiment instead of resetting them to the template.")
	// Read once the flags are parsed; the negated flag keeps the default of resetting.
	cmd.PreRun = func(*cobra.Command, []string) { o.ResetProperties = !noReset }
	cmd.MarkFlagsMutuallyExclusive("template", "file")
}

func newExperimentRun() *cobra.Command {
	var o experiment.RunOptions
	var noWait bool
	executionVariables := newKeyValues()
	o.ExecutionVariable = executionVariables.values
	cmd := &cobra.Command{
		Use:     "run",
		Aliases: []string{"exec"},
		Short:   "Executes an experiment run. If a file is specified the experiment is saved before execution.",
		Args:    cobra.NoArgs,
		Example: examples(
			"steadybit experiment run -k ADM-1",
			"steadybit experiment run -f experiment.yml --no-wait",
			"steadybit experiment run -f ./experiments -R --yes --timeout 30m --report steadybit.xml",
			"steadybit experiment run --template d7e65100-1d20-4980-be87-c351704910b8 --team ADM -p CLUSTER=prod",
		),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			o.Wait = !noWait
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
	f.DurationVar(&o.Timeout, "timeout", 0, `With waiting: cancel the run and fail when it has not ended after this long, e.g. "15m".`)
	f.BoolVar(&o.KeepRunningOnInterrupt, "keep-running-on-interrupt", false, "With waiting: leave the run going when the CLI is interrupted, instead of cancelling it.")
	f.BoolVar(&o.ShowSteps, "show-steps", false, "With waiting: print each step's state as it changes.")
	f.StringVar(&o.Report, "report", "", `With waiting: write a JUnit report of the runs to this file, or JSON if it ends in ".json".`)
	f.Var(executionVariables, "execution-variable", "With --template: a variable for this run only, overriding experiment and environment variables. Repeat for more.")
	addTemplateFlags(cmd, &o.TemplateOptions)
	cmd.MarkFlagsMutuallyExclusive("key", "file")
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
	var t experiment.TemplateOptions
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Upload an experiment to Steadybit. If a key is provided, an update is performed. Otherwise, the externalId from the file is used to create or update the experiment. With --template, the experiment is created from an experiment template instead of a file.",
		Args:  cobra.NoArgs,
		Example: examples(
			"steadybit experiment apply -f experiment.yml",
			"steadybit experiment apply -f ./experiments -R",
			"steadybit experiment apply --template d7e65100-1d20-4980-be87-c351704910b8 --team ADM --external-id shop-latency -p CLUSTER=prod",
			"steadybit experiment apply --template d7e65100-1d20-4980-be87-c351704910b8 -k ADM-12 --placeholders values.yml",
		),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			if t.Template != "" {
				return experiment.ApplyTemplate(ctx, c, o.Key, t)
			}
			if len(o.Files) == 0 {
				return errors.New("Either --file or --template must be specified.")
			}
			return experiment.Apply(ctx, c, o)
		}),
	}
	cmd.Flags().StringVarP(&o.Key, "key", "k", "", "The experiment key.")
	cmd.Flags().StringArrayVarP(&o.Files, "file", "f", nil, "The path to the experiment file or a directory containing multiple files.")
	cmd.Flags().BoolVarP(&o.Recursive, "recursive", "R", false, "Process the directory used in -f, --file recursively.")
	addTemplateFlags(cmd, &t)
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
