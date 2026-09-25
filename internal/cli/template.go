// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/template"
)

const (
	typeHelp   = `The output format ("json" or "yaml"). (default: "json" if the file ends in ".json", "yaml" otherwise.)`
	yesHelp    = "Skip the confirmation prompt. Not necessary when no TTY is attached."
	templateID = "d7e65100-1d20-4980-be87-c351704910b8"
	hubID      = "0194a7d4-0d1f-7b21-9c64-5b6e3c1f2a10"
)

// fileFlags adds the flags every `apply` takes.
func fileFlags(cmd *cobra.Command, files *[]string, recursive *bool, what string) {
	cmd.Flags().StringArrayVarP(files, "file", "f", nil, "The path to the "+what+" file or a directory containing multiple files.")
	cmd.Flags().BoolVarP(recursive, "recursive", "R", false, "Process the directory used in -f, --file recursively.")
	_ = cmd.MarkFlagRequired("file")
	variadic(cmd, "file")
}

// outputFlags adds the flags every `get` takes.
func outputFlags(cmd *cobra.Command, file, datatype *string, what string) {
	cmd.Flags().StringVarP(file, "file", "f", "", "The path to write the "+what+" to.")
	cmd.Flags().StringVarP(datatype, "type", "t", "", typeHelp)
}

func newTemplate() *cobra.Command {
	cmd := &cobra.Command{Use: "template", Short: "Manage the experiment templates to create experiments from."}

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
	variadic(list, "tag", "target-type", "action", "search")

	var g template.GetOptions
	get := &cobra.Command{
		Use:   "get",
		Short: "Get an experiment template. Output is written to file or stdout.",
		Args:  cobra.NoArgs,
		Example: examples(
			"steadybit template get -i "+templateID+" -f template.yml",
			"steadybit template get -i "+templateID+" --placeholders -f values.yml",
		),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return template.Get(ctx, c, g) }),
	}
	get.Flags().StringVarP(&g.ID, "id", "i", "", "The experiment template id.")
	get.Flags().StringVarP(&g.File, "file", "f", "", "The path to write the template to.")
	get.Flags().StringVarP(&g.Type, "type", "t", "", typeHelp)
	get.Flags().BoolVar(&g.Placeholders, "placeholders", false, "Only output the template placeholders, as a file to fill in and pass to --placeholders.")
	_ = get.MarkFlagRequired("id")

	var a template.ApplyOptions
	apply := &cobra.Command{
		Use:     "apply",
		Short:   "Create or update experiment templates from files. A file without an id creates a template, and the new id is written back to it.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit template apply -f template.yml", "steadybit template apply -f ./templates -R"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return template.Apply(ctx, c, a) }),
	}
	fileFlags(apply, &a.Files, &a.Recursive, "template")

	var d template.DeleteOptions
	del := &cobra.Command{
		Use:     "delete",
		Short:   "Delete an experiment template. Service profiles lose it, and the experiments they provided from it are deleted.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit template delete -i " + templateID),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return template.Delete(ctx, c, d) }),
	}
	idFlag(del, &d.ID, "The experiment template id.")
	del.Flags().BoolVar(&d.Yes, "yes", false, yesHelp)

	var im template.ImportOptions
	imp := &cobra.Command{
		Use:     "import",
		Short:   "Import experiment templates from a connected hub.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit template import --hub " + hubID + " --template " + templateID + " --overwrite"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return template.Import(ctx, c, im) }),
	}
	imp.Flags().StringVar(&im.Hub, "hub", "", "The id of the hub, see `steadybit hub list`.")
	imp.Flags().StringArrayVar(&im.Templates, "template", nil, "The ids of the hub's templates to import.")
	imp.Flags().BoolVar(&im.Overwrite, "overwrite", false, "Replace templates that exist already. Without it, the import fails if any does.")
	_ = imp.MarkFlagRequired("hub")
	_ = imp.MarkFlagRequired("template")
	variadic(imp, "template")

	cmd.AddCommand(list, get, apply, del, imp)
	return cmd
}
