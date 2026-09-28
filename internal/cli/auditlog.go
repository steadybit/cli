// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/v6/internal/auditlog"
	"github.com/steadybit/cli/v6/internal/platform"
)

func newAuditLog() *cobra.Command {
	var o auditlog.Options
	cmd := &cobra.Command{
		Use:     "audit-log",
		Short:   "Show the audit log: who changed what, and when. Needs an admin access token.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit audit-log", "steadybit audit-log --from 2026-09-01 --to 2026-09-08T12:00:00Z -t json"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return auditlog.Show(ctx, c, o) }),
	}
	cmd.Flags().StringVar(&o.From, "from", "", "The earliest time, a date or an RFC 3339 time. (default: 7 days before --to)")
	cmd.Flags().StringVar(&o.To, "to", "", "The latest time, a date or an RFC 3339 time. (default: 7 days after --from, or now)")
	cmd.Flags().StringVarP(&o.Type, "type", "t", "", `Print the entries as "json" or "yaml" instead of a table.`)
	return cmd
}
