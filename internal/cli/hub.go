// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/internal/hub"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/resource"
)

func newHub() *cobra.Command {
	cmd := &cobra.Command{Use: "hub", Short: "Manage the hubs experiment templates are imported from."}

	var listType string
	list := &cobra.Command{
		Use:     "list",
		Short:   "List hubs.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit hub list"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return hub.List(ctx, c, listType) }),
	}
	list.Flags().StringVarP(&listType, "type", "t", "", resource.ListTypeHelp)

	var g hub.GetOptions
	get := &cobra.Command{
		Use:     "get",
		Short:   "Get a hub. Output is written to file or stdout.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit hub get -i " + hubID + " -f hub.yml"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return hub.Get(ctx, c, g) }),
	}
	idFlag(get, &g.ID, "The hub id.")
	outputFlags(get, &g.File, &g.Type, "hub")

	var a hub.ApplyOptions
	apply := &cobra.Command{
		Use:     "apply",
		Short:   "Create or update hubs from files. A file without an id creates a hub, and the new id is written back to it.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit hub apply -f hub.yml --synchronize"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return hub.Apply(ctx, c, a) }),
	}
	fileFlags(apply, &a.Files, &a.Recursive, "hub")
	apply.Flags().BoolVar(&a.Synchronize, "synchronize", false, "Fetch the hub's templates from its repository, waiting until it is done.")

	var d hub.DeleteOptions
	del := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a hub. The templates imported from it are kept unless --delete-imported-templates is given.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit hub delete -i " + hubID),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return hub.Delete(ctx, c, d) }),
	}
	idFlag(del, &d.ID, "The hub id.")
	del.Flags().BoolVar(&d.Templates, "delete-imported-templates", false, "Also delete the templates imported from the hub.")
	del.Flags().BoolVar(&d.Yes, "yes", false, yesHelp)

	var resyncID string
	resync := &cobra.Command{
		Use:     "resync",
		Short:   "Fetch a hub's templates from its repository again, waiting until it is done.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit hub resync -i " + hubID),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return hub.Resync(ctx, c, resyncID) }),
	}
	idFlag(resync, &resyncID, "The hub id.")

	cmd.AddCommand(list, get, apply, del, resync)
	return cmd
}
