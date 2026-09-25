// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/internal/action"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/resource"
	"github.com/steadybit/cli/internal/target"
)

const listTypeHelp = `Print the list as "json" or "yaml" instead of a table.`

func environmentNameFlag(cmd *cobra.Command, environment *string) {
	cmd.Flags().StringVarP(environment, "environment", "e", "", "The environment name.")
	_ = cmd.MarkFlagRequired("environment")
}

func newTarget() *cobra.Command {
	cmd := &cobra.Command{Use: "target", Short: "Find the targets of an environment and their attributes."}

	var q target.QueryOptions
	query := &cobra.Command{
		Use:   "query",
		Short: "List the targets of an environment, optionally of one type and matching a query.",
		Args:  cobra.NoArgs,
		Example: examples(
			"steadybit target query -e Global --target-type com.steadybit.extension_container.container",
			`steadybit target query -e Global -q 'k8s.namespace="shop"' --attribute k8s.deployment k8s.pod.name --limit 0 -t json`,
		),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return target.Query(ctx, c, q) }),
	}
	environmentNameFlag(query, &q.Environment)
	query.Flags().StringVar(&q.TargetType, "target-type", "", "Only list targets of this type.")
	query.Flags().StringVarP(&q.Query, "query", "q", "", "Only list targets matching this target query.")
	query.Flags().StringArrayVar(&q.Attributes, "attribute", nil, "Only fetch these attributes, shown as columns. (default: all, not shown)")
	query.Flags().IntVar(&q.Limit, "limit", 100, "List at most this many targets; 0 lists all.")
	query.Flags().StringVarP(&q.Type, "type", "t", "", listTypeHelp)
	variadic(query, "attribute")

	attribute := &cobra.Command{Use: "attribute", Short: "List the attribute keys and values of targets."}
	attributeFlags := func(c *cobra.Command, o *target.AttributeOptions) {
		environmentNameFlag(c, &o.Environment)
		c.Flags().StringVar(&o.TargetType, "target-type", "", "The target type. Either this or --action is required.")
		c.Flags().StringVar(&o.Action, "action", "", "The action whose target selection to use, for actions with an extended one.")
		c.Flags().StringVarP(&o.Type, "type", "t", "", listTypeHelp)
	}
	var k target.AttributeOptions
	keys := &cobra.Command{
		Use:     "keys",
		Short:   "List the attribute keys the targets of a type have.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit target attribute keys -e Global --target-type com.steadybit.extension_container.container"),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return target.AttributeKeys(ctx, c, k)
		}),
	}
	attributeFlags(keys, &k)
	var v target.AttributeOptions
	values := &cobra.Command{
		Use:     "values",
		Short:   "List the values an attribute has on the targets of a type.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit target attribute values -e Global --target-type com.steadybit.extension_container.container --key k8s.namespace"),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return target.AttributeValues(ctx, c, v)
		}),
	}
	attributeFlags(values, &v)
	values.Flags().StringVarP(&v.Key, "key", "k", "", "The attribute key.")
	_ = values.MarkFlagRequired("key")
	attribute.AddCommand(keys, values)

	cmd.AddCommand(query, attribute)
	return cmd
}

func newAction() *cobra.Command {
	cmd := &cobra.Command{Use: "action", Short: "Find the actions experiments can use."}

	var l action.ListOptions
	list := &cobra.Command{
		Use:     "list",
		Short:   "List the actions the extensions provide.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit action list", "steadybit action list --kind ATTACK CHECK"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return action.List(ctx, c, l) }),
	}
	list.Flags().StringArrayVar(&l.Kinds, "kind", nil, `Only list actions of these kinds: "ATTACK", "CHECK", "LOAD_TEST", "OTHER" or "BASIC".`)
	list.Flags().StringVarP(&l.Type, "type", "t", "", resource.ListTypeHelp)
	variadic(list, "kind")

	var g action.GetOptions
	get := &cobra.Command{
		Use:     "get",
		Short:   "Get an action with its parameters. Output is written to file or stdout.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit action get -i com.steadybit.extension_container.stress_cpu"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return action.Get(ctx, c, g) }),
	}
	idFlag(get, &g.ID, "The action id.")
	outputFlags(get, &g.File, &g.Type, "action")

	cmd.AddCommand(list, get)
	return cmd
}
