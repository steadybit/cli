// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"context"
	"strings"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/internal/integration"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/resource"
)

const integrationID = "0190d7b2-7d3e-7a4b-8c5d-6e7f8a9b0c1d"

func newIntegration() *cobra.Command {
	cmd := &cobra.Command{Use: "integration", Short: "Manage webhook, Slack and preflight integrations. Changing them needs an admin access token."}
	for _, kind := range integration.Kinds {
		cmd.AddCommand(newIntegrationKind(kind))
	}
	return cmd
}

func newIntegrationKind(k integration.Kind) *cobra.Command {
	plural := k.Plural
	// Slack keeps its capital in the middle of a sentence.
	lower := k.Title
	if k.Name != "slack" {
		lower = strings.ToLower(k.Title[:1]) + k.Title[1:]
	}
	prefix := "steadybit integration " + k.Name
	cmd := &cobra.Command{Use: k.Name, Short: "Manage " + plural + "."}

	var listType string
	list := &cobra.Command{
		Use:     "list",
		Short:   "List " + plural + ".",
		Args:    cobra.NoArgs,
		Example: examples(prefix + " list"),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return integration.List(ctx, c, k, listType)
		}),
	}
	list.Flags().StringVarP(&listType, "type", "t", "", resource.ListTypeHelp)

	var g integration.GetOptions
	get := &cobra.Command{
		Use:     "get",
		Short:   "Get a " + lower + ". Output is written to file or stdout.",
		Args:    cobra.NoArgs,
		Example: examples(prefix + " get -i " + integrationID + " -f " + k.Name + ".yml"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return integration.Get(ctx, c, k, g) }),
	}
	idFlag(get, &g.ID, "The "+lower+" id.")
	outputFlags(get, &g.File, &g.Type, lower)

	var a integration.ApplyOptions
	apply := &cobra.Command{
		Use:     "apply",
		Short:   "Create or update " + plural + " from files. A file without an id creates one, and the new id is written back to it.",
		Args:    cobra.NoArgs,
		Example: examples(prefix+" apply -f "+k.Name+".yml", prefix+" apply -f ./integrations -R"),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return integration.Apply(ctx, c, k, a)
		}),
	}
	fileFlags(apply, &a.Files, &a.Recursive, lower)

	var d integration.DeleteOptions
	del := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a " + lower + ".",
		Args:    cobra.NoArgs,
		Example: examples(prefix + " delete -i " + integrationID),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return integration.Delete(ctx, c, k, d)
		}),
	}
	idFlag(del, &d.ID, "The "+lower+" id.")
	del.Flags().BoolVar(&d.Yes, "yes", false, yesHelp)

	cmd.AddCommand(list, get, apply, del)
	return cmd
}
