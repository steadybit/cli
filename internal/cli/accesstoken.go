// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/internal/accesstoken"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/resource"
)

const (
	accessTokenID = "Xy12AbCd"
	tokenTypeHelp = `Print only the id and the token as "json" or "yaml", for a script to read.`
)

func newAccessToken() *cobra.Command {
	cmd := &cobra.Command{Use: "access-token", Short: "Manage the API access tokens of the tenant."}

	var l accesstoken.ListOptions
	var expired bool
	list := &cobra.Command{
		Use:     "list",
		Short:   "List access tokens. The tokens themselves are never shown.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit access-token list", "steadybit access-token list --type TEAM --team ADM OPS --expired=false"),
		PreRun: func(cmd *cobra.Command, _ []string) {
			if cmd.Flags().Changed("expired") {
				l.Expired = &expired
			}
		},
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return accesstoken.List(ctx, c, l) }),
	}
	list.Flags().StringVar(&l.Name, "name", "", "Only list tokens with this name.")
	list.Flags().StringVar(&l.CreatedBy, "created-by", "", "Only list tokens created by this user.")
	list.Flags().StringVar(&l.Type, "type", "", `Only list tokens of this type, "ADMIN", "TEAM" or "WILDCARD".`)
	list.Flags().StringVar(&l.Output, "output", "", resource.ListTypeHelp)
	list.Flags().StringArrayVar(&l.Teams, "team", nil, "Only list tokens of these teams, by team key.")
	list.Flags().BoolVar(&expired, "expired", false, "Only list expired tokens, or with --expired=false those still valid.")
	variadic(list, "team")

	var cr accesstoken.CreateOptions
	create := &cobra.Command{
		Use:   "create",
		Short: "Create an access token. The token is printed once and cannot be shown again.",
		Args:  cobra.NoArgs,
		Example: examples(
			"steadybit access-token create --name ci --type TEAM --team ADM --expires-at 2026-12-31",
			"steadybit access-token create --name ci --type TEAM --team ADM -t json",
		),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return accesstoken.Create(ctx, c, cr) }),
	}
	create.Flags().StringVar(&cr.Name, "name", "", "The name of the token.")
	create.Flags().StringVar(&cr.Type, "type", "", `The type of the token: "ADMIN", "TEAM" for the teams given with --team, or "WILDCARD" for every team.`)
	create.Flags().StringArrayVar(&cr.Teams, "team", nil, "The teams of a TEAM token, by team key.")
	create.Flags().StringVar(&cr.ExpiresAt, "expires-at", "", "When the token expires, a date or an RFC 3339 time. (default: never)")
	create.Flags().StringVarP(&cr.Output, "output", "t", "", tokenTypeHelp)
	_ = create.MarkFlagRequired("name")
	_ = create.MarkFlagRequired("type")
	variadic(create, "team")

	var r accesstoken.RecreateOptions
	recreate := &cobra.Command{
		Use:     "recreate",
		Short:   "Replace an access token with a new one of the same name, type and teams. The old token stops working.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit access-token recreate -i " + accessTokenID + " --expires-at 2027-06-30"),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return accesstoken.Recreate(ctx, c, r)
		}),
	}
	idFlag(recreate, &r.ID, "The access token id.")
	recreate.Flags().StringVar(&r.ExpiresAt, "expires-at", "", "When the new token expires, a date or an RFC 3339 time. (default: never)")
	recreate.Flags().StringVarP(&r.Output, "output", "t", "", tokenTypeHelp)
	recreate.Flags().BoolVar(&r.Yes, "yes", false, yesHelp)

	var d accesstoken.DeleteOptions
	del := &cobra.Command{
		Use:     "delete",
		Short:   "Delete an access token. Everything using it loses access.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit access-token delete -i " + accessTokenID),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return accesstoken.Delete(ctx, c, d) }),
	}
	idFlag(del, &d.ID, "The access token id.")
	del.Flags().BoolVar(&d.Yes, "yes", false, yesHelp)

	cmd.AddCommand(list, create, recreate, del)
	return cmd
}
