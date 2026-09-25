// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/property"
	"github.com/steadybit/cli/internal/resource"
)

const associationID = "0190d7b2-9e8f-7c6d-b5a4-3f2e1d0c9b8a"

func propertyKeyFlag(cmd *cobra.Command, key *string) {
	cmd.Flags().StringVarP(key, "key", "k", "", "The property key.")
	_ = cmd.MarkFlagRequired("key")
}

func newProperty() *cobra.Command {
	cmd := &cobra.Command{Use: "property", Short: "Manage the properties experiments and services carry."}
	cmd.AddCommand(newPropertyDefinition(), newPropertyAssociation())
	return cmd
}

func newPropertyDefinition() *cobra.Command {
	cmd := &cobra.Command{Use: "definition", Short: "Manage property definitions: a property's key, label and type."}

	var listType string
	list := &cobra.Command{
		Use:     "list",
		Short:   "List property definitions.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit property definition list"),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return property.ListDefinitions(ctx, c, listType)
		}),
	}
	list.Flags().StringVarP(&listType, "type", "t", "", resource.ListTypeHelp)

	var g property.GetDefinitionOptions
	get := &cobra.Command{
		Use:     "get",
		Short:   "Get a property definition. Output is written to file or stdout.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit property definition get -k RESULT_COLOR -f result-color.yml"),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return property.GetDefinition(ctx, c, g)
		}),
	}
	propertyKeyFlag(get, &g.Key)
	outputFlags(get, &g.File, &g.Type, "property definition")

	var a property.ApplyDefinitionOptions
	apply := &cobra.Command{
		Use:     "apply",
		Short:   "Create or update property definitions from files. The key names the definition to update.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit property definition apply -f result-color.yml", "steadybit property definition apply -f ./properties -R"),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return property.ApplyDefinitions(ctx, c, a)
		}),
	}
	fileFlags(apply, &a.Files, &a.Recursive, "property definition")
	apply.Flags().BoolVar(&a.DeleteValues, "delete-values", false, "Allow removing enum values still in use, deleting them where they are used.")

	var d property.DeleteDefinitionOptions
	del := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a property definition. Only one without associations can be deleted, unless --delete-associations is given.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit property definition delete -k RESULT_COLOR"),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return property.DeleteDefinition(ctx, c, d)
		}),
	}
	propertyKeyFlag(del, &d.Key)
	del.Flags().BoolVar(&d.Associations, "delete-associations", false, "Also delete the associations of the property.")
	del.Flags().BoolVar(&d.Yes, "yes", false, yesHelp)

	cmd.AddCommand(list, get, apply, del)
	return cmd
}

func newPropertyAssociation() *cobra.Command {
	cmd := &cobra.Command{Use: "association", Short: "Manage property associations: which experiments or services carry a property."}

	var l property.ListAssociationOptions
	list := &cobra.Command{
		Use:     "list",
		Short:   "List property associations.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit property association list", "steadybit property association list --key RESULT_COLOR --type EXPERIMENT"),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return property.ListAssociations(ctx, c, l)
		}),
	}
	list.Flags().StringVar(&l.Key, "key", "", "Only list associations of this property.")
	list.Flags().StringVar(&l.Experiment, "experiment", "", "Only list associations given to this experiment, by key. Those for all experiments are not listed.")
	list.Flags().StringVar(&l.Service, "service", "", "Only list associations given to this service, by id. Those for all services are not listed.")
	list.Flags().StringVar(&l.Type, "type", "", `Only list "EXPERIMENT" or "SERVICE" associations.`)
	list.Flags().StringVar(&l.Output, "output", "", resource.ListTypeHelp)

	var g property.GetAssociationOptions
	get := &cobra.Command{
		Use:     "get",
		Short:   "Get a property association. Output is written to file or stdout.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit property association get -i " + associationID + " -f association.yml"),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return property.GetAssociation(ctx, c, g)
		}),
	}
	idFlag(get, &g.ID, "The property association id.")
	outputFlags(get, &g.File, &g.Type, "property association")

	var a property.ApplyAssociationOptions
	apply := &cobra.Command{
		Use:     "apply",
		Short:   "Create or update property associations from files. A file without an id creates one, and the new id is written back to it.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit property association apply -f association.yml"),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return property.ApplyAssociations(ctx, c, a)
		}),
	}
	fileFlags(apply, &a.Files, &a.Recursive, "property association")

	var d property.DeleteAssociationOptions
	del := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a property association. Only one whose values are not used can be deleted, unless --delete-values is given.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit property association delete -i " + associationID),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return property.DeleteAssociation(ctx, c, d)
		}),
	}
	idFlag(del, &d.ID, "The property association id.")
	del.Flags().BoolVar(&d.DeleteValues, "delete-values", false, "Also delete the values experiments and schedules have for it.")
	del.Flags().BoolVar(&d.Yes, "yes", false, yesHelp)

	cmd.AddCommand(list, get, apply, del)
	return cmd
}
