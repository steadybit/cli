// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/internal/environment"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/resource"
)

const environmentID = "0190d7b2-1c5e-7f3a-8e4b-2d6f9a1c3e57"

func newEnvironment() *cobra.Command {
	cmd := &cobra.Command{Use: "environment", Short: "Manage environments and their variables."}

	var l environment.ListOptions
	list := &cobra.Command{
		Use:     "list",
		Short:   "List environments.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit environment list", "steadybit environment list --search prod"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return environment.List(ctx, c, l) }),
	}
	list.Flags().StringVar(&l.Search, "search", "", "Only list environments whose name, or the name or key of a team using them, matches.")
	list.Flags().StringVarP(&l.Type, "type", "t", "", resource.ListTypeHelp)

	var g environment.GetOptions
	get := &cobra.Command{
		Use:     "get",
		Short:   "Get an environment. Output is written to file or stdout.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit environment get -i " + environmentID + " -f environment.yml"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return environment.Get(ctx, c, g) }),
	}
	idFlag(get, &g.ID, "The environment id.")
	outputFlags(get, &g.File, &g.Type, "environment")

	var a environment.ApplyOptions
	apply := &cobra.Command{
		Use:     "apply",
		Short:   "Create or update environments from files. A file without an id creates an environment, and the new id is written back to it.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit environment apply -f environment.yml", "steadybit environment apply -f ./environments -R"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return environment.Apply(ctx, c, a) }),
	}
	fileFlags(apply, &a.Files, &a.Recursive, "environment")

	var d environment.DeleteOptions
	del := &cobra.Command{
		Use:     "delete",
		Short:   "Delete an environment.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit environment delete -i " + environmentID),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return environment.Delete(ctx, c, d) }),
	}
	idFlag(del, &d.ID, "The environment id.")
	del.Flags().BoolVar(&d.Yes, "yes", false, yesHelp)

	variable := &cobra.Command{Use: "variable", Short: "Manage the variables of an environment."}
	var vg environment.VariableGetOptions
	vget := &cobra.Command{
		Use:     "get",
		Short:   "Print the variables of an environment.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit environment variable get -i " + environmentID),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return environment.GetVariables(ctx, c, vg)
		}),
	}
	idFlag(vget, &vg.ID, "The environment id.")
	vget.Flags().StringVarP(&vg.Type, "type", "t", "yaml", `The output format ("json" or "yaml").`)
	var vs environment.VariableSetOptions
	vset := &cobra.Command{
		Use:   "set [KEY=VALUE...]",
		Short: "Set variables of an environment, keeping the others. With --replace, the given variables become the only ones.",
		Example: examples(
			"steadybit environment variable set -i "+environmentID+" region=eu cluster=prod",
			"steadybit environment variable set -i "+environmentID+" -f variables.yml --replace",
		),
		RunE: withClient(func(ctx context.Context, c *platform.Client, args []string) error {
			return environment.SetVariables(ctx, c, args, vs)
		}),
	}
	idFlag(vset, &vs.ID, "The environment id.")
	vset.Flags().StringVarP(&vs.File, "file", "f", "", "A YAML/JSON file mapping variable names to values, which may be lists or select expressions.")
	vset.Flags().BoolVar(&vs.Replace, "replace", false, "Remove every variable not given.")
	variable.AddCommand(vget, vset)

	cmd.AddCommand(list, get, apply, del, variable)
	return cmd
}
