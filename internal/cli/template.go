// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/resource"
	"github.com/steadybit/cli/internal/template"
)

const typeHelp = `The output format ("json" or "yaml"). (default: "json" if the file ends in ".json", "yaml" otherwise.)`

func newTemplate() *cobra.Command {
	cmd := &cobra.Command{Use: "template", Short: "Find experiment templates to create experiments from."}

	var l template.ListOptions
	list := &cobra.Command{
		Use:   "list",
		Short: "List experiment templates. Filters of the same kind match any of the given values.",
		Args:  cobra.NoArgs,
		Example: examples(
			"steadybit template list",
			"steadybit template list --search kubernetes --action com.steadybit.extension_host.stress-cpu",
		),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return template.List(ctx, c, l) }),
	}
	list.Flags().StringArrayVar(&l.Tags, "tag", nil, "Only list templates with one of these tags.")
	list.Flags().StringArrayVar(&l.TargetTypes, "target-type", nil, "Only list templates targeting one of these target types.")
	list.Flags().StringArrayVar(&l.Actions, "action", nil, "Only list templates using one of these actions.")
	list.Flags().StringArrayVar(&l.Search, "search", nil, "Only list templates whose title or description match.")
	list.Flags().StringVarP(&l.Type, "type", "t", "", resource.ListTypeHelp)
	variadic(list, "tag", "target-type", "action", "search")

	var g template.GetOptions
	get := &cobra.Command{
		Use:   "get",
		Short: "Get an experiment template. Output is written to file or stdout.",
		Args:  cobra.NoArgs,
		Example: examples(
			"steadybit template get -i d7e65100-1d20-4980-be87-c351704910b8",
			"steadybit template get -i d7e65100-1d20-4980-be87-c351704910b8 --placeholders -f values.yml",
		),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return template.Get(ctx, c, g) }),
	}
	get.Flags().StringVarP(&g.ID, "id", "i", "", "The experiment template id.")
	get.Flags().StringVarP(&g.File, "file", "f", "", "The path to write the template to.")
	get.Flags().StringVarP(&g.Type, "type", "t", "", typeHelp)
	get.Flags().BoolVar(&g.Placeholders, "placeholders", false, "Only output the template placeholders, as a file to fill in and pass to --placeholders.")
	_ = get.MarkFlagRequired("id")

	cmd.AddCommand(list, get)
	return cmd
}
