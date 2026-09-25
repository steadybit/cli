// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/resource"
	"github.com/steadybit/cli/internal/team"
)

func teamKeyFlag(cmd *cobra.Command, key *string) {
	cmd.Flags().StringVarP(key, "key", "k", "", "The team key.")
	_ = cmd.MarkFlagRequired("key")
}

func newTeam() *cobra.Command {
	cmd := &cobra.Command{Use: "team", Short: "Manage teams, their members and environments."}

	var l team.ListOptions
	list := &cobra.Command{
		Use:     "list",
		Short:   "List teams.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit team list", "steadybit team list --search jane@example.com"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return team.List(ctx, c, l) }),
	}
	list.Flags().StringVar(&l.Search, "search", "", "Only list teams whose name or key, or a member's name or email, matches.")
	list.Flags().StringVarP(&l.Type, "type", "t", "", resource.ListTypeHelp)

	var g team.GetOptions
	get := &cobra.Command{
		Use:     "get",
		Short:   "Get a team. Output is written to file or stdout.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit team get -k ADM -f team.yml"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return team.Get(ctx, c, g) }),
	}
	teamKeyFlag(get, &g.Key)
	outputFlags(get, &g.File, &g.Type, "team")

	var a team.ApplyOptions
	apply := &cobra.Command{
		Use:     "apply",
		Short:   "Create or update teams from files. The team key names the team to update.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit team apply -f team.yml", "steadybit team apply -f ./teams -R"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return team.Apply(ctx, c, a) }),
	}
	fileFlags(apply, &a.Files, &a.Recursive, "team")

	var d team.DeleteOptions
	del := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a team, with --purge-experiments, which the platform requires. Nothing may be running in it.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit team delete -k OPS --purge-experiments", "steadybit team delete -k OPS --purge-experiments --yes"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return team.Delete(ctx, c, d) }),
	}
	teamKeyFlag(del, &d.Key)
	del.Flags().BoolVar(&d.Experiments, "purge-experiments", false, "Also delete the team's experiments and their runs.")
	del.Flags().BoolVar(&d.Yes, "yes", false, yesHelp)

	cmd.AddCommand(list, get, apply, del, newTeamMember(), newTeamEnvironment())
	return cmd
}

func newTeamMember() *cobra.Command {
	cmd := &cobra.Command{Use: "member", Short: "Manage the members of a team."}

	var key, listType string
	list := &cobra.Command{
		Use:     "list",
		Short:   "List the members of a team.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit team member list -k ADM"),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return team.ListMembers(ctx, c, key, listType)
		}),
	}
	teamKeyFlag(list, &key)
	list.Flags().StringVarP(&listType, "type", "t", "", resource.ListTypeHelp)

	change := func(use, short string, withRole, confirm bool, example string, run func(context.Context, *platform.Client, team.MemberOptions) error) *cobra.Command {
		var o team.MemberOptions
		c := &cobra.Command{
			Use:     use,
			Short:   short,
			Args:    cobra.NoArgs,
			Example: examples(example),
			RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return run(ctx, c, o) }),
		}
		teamKeyFlag(c, &o.Key)
		c.Flags().StringArrayVar(&o.Usernames, "username", nil, "Members by username.")
		c.Flags().StringArrayVar(&o.Emails, "email", nil, "Members by email.")
		variadic(c, "username", "email")
		if withRole {
			c.Flags().StringVar(&o.Role, "role", "MEMBER", `The role of the given members, "MEMBER" or "OWNER".`)
			c.Flags().BoolVar(&o.Validate, "validate", false, "Fail on users the platform does not know, instead of skipping them.")
		}
		if confirm {
			c.Flags().BoolVar(&o.Yes, "yes", false, yesHelp)
		}
		return c
	}
	cmd.AddCommand(list,
		change("add", "Add members to a team, or change the role of existing ones.", true, false,
			"steadybit team member add -k ADM --email jane@example.com joe@example.com --role OWNER", team.AddMembers),
		change("remove", "Remove members from a team. They keep their account.", false, true,
			"steadybit team member remove -k ADM --email jane@example.com", team.RemoveMembers),
		change("set", "Make the given users the only members of a team, removing everyone else.", true, true,
			"steadybit team member set -k ADM --email jane@example.com --role OWNER --yes", team.SetMembers),
	)
	return cmd
}

func newTeamEnvironment() *cobra.Command {
	cmd := &cobra.Command{Use: "environment", Short: "Manage the environments a team may use."}

	var key, envListType string
	list := &cobra.Command{
		Use:     "list",
		Short:   "List the environments of a team.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit team environment list -k ADM"),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return team.ListEnvironments(ctx, c, key, envListType)
		}),
	}
	teamKeyFlag(list, &key)
	list.Flags().StringVarP(&envListType, "type", "t", "", resource.ListTypeHelp)

	change := func(use, short string, validate, confirm bool, example string, run func(context.Context, *platform.Client, team.EnvironmentOptions) error) *cobra.Command {
		var o team.EnvironmentOptions
		c := &cobra.Command{
			Use:     use,
			Short:   short,
			Args:    cobra.NoArgs,
			Example: examples(example),
			RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return run(ctx, c, o) }),
		}
		teamKeyFlag(c, &o.Key)
		c.Flags().StringArrayVar(&o.Environments, "environment", nil, "Environments by name.")
		variadic(c, "environment")
		if validate {
			c.Flags().BoolVar(&o.Validate, "validate", false, "Fail on environments the platform does not know, instead of skipping them.")
		}
		if confirm {
			c.Flags().BoolVar(&o.Yes, "yes", false, yesHelp)
		}
		return c
	}
	cmd.AddCommand(list,
		change("add", "Allow a team to use environments.", true, false,
			`steadybit team environment add -k ADM --environment Global "Online Shop"`, team.AddEnvironments),
		change("remove", "Stop a team from using environments.", false, false,
			"steadybit team environment remove -k ADM --environment Global", team.RemoveEnvironments),
		change("set", "Make the given environments the only ones a team may use.", true, true,
			"steadybit team environment set -k ADM --environment Global --yes", team.SetEnvironments),
	)
	return cmd
}
