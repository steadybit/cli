// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/v6/internal/platform"
	"github.com/steadybit/cli/v6/internal/resource"
	"github.com/steadybit/cli/v6/internal/user"
)

func newUser() *cobra.Command {
	cmd := &cobra.Command{Use: "user", Short: "Manage the users of the tenant."}

	var o user.InviteOptions
	invite := &cobra.Command{
		Use:     "invite",
		Short:   "Invite users by email. Each receives an email with a link to join. Needs an admin access token.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit user invite --email jane@example.com joe@example.com --team ADM"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return user.Invite(ctx, c, o) }),
	}
	invite.Flags().StringArrayVar(&o.Emails, "email", nil, "The email addresses to invite.")
	invite.Flags().StringVar(&o.Role, "role", "", `The role in the tenant, "USER" or "ADMIN". (default: the platform's default role)`)
	invite.Flags().StringVar(&o.Team, "team", "", "The key of a team to add them to.")
	_ = invite.MarkFlagRequired("email")
	variadic(invite, "email")

	var l user.ListOptions
	list := &cobra.Command{
		Use:     "list",
		Short:   "List the users of the tenant. Needs an admin access token.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit user list", "steadybit user list --query jane", `steadybit user list --jq '.[] | select(.managedBy == "SCIM") | .username'`),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return user.List(ctx, c, l) }),
	}
	list.Flags().StringVar(&l.Query, "query", "", "Only list users whose username, name or email contains this text.")
	list.Flags().StringVarP(&l.Type, "type", "t", "", resource.ListTypeHelp)

	var d user.DeleteOptions
	del := &cobra.Command{
		Use:     "delete",
		Short:   "Remove a user from the tenant. Users managed by SCIM or an LDAP sync are removed there instead. Needs an admin access token.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit user delete --username 13a52737-b318-4048-a79d-4789d645bc31"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return user.Delete(ctx, c, d) }),
	}
	del.Flags().StringVarP(&d.Username, "username", "u", "", "The username, as user list shows it.")
	del.Flags().BoolVar(&d.Yes, "yes", false, yesHelp)
	_ = del.MarkFlagRequired("username")

	cmd.AddCommand(invite, list, del)
	return cmd
}
