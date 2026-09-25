// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/user"
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

	cmd.AddCommand(invite)
	return cmd
}
