// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/v6/internal/killswitch"
	"github.com/steadybit/cli/v6/internal/platform"
)

func newKillswitch() *cobra.Command {
	cmd := &cobra.Command{Use: "killswitch", Short: "Stop all experiments of the tenant at once, and keep them stopped."}

	var s killswitch.StatusOptions
	status := &cobra.Command{
		Use:     "status",
		Short:   "Show whether the kill switch is active.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit killswitch status", "steadybit killswitch status -t json"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return killswitch.Status(ctx, c, s) }),
	}
	status.Flags().StringVarP(&s.Type, "type", "t", "", `Print the status as "json" or "yaml".`)

	var yes bool
	activate := &cobra.Command{
		Use:     "activate",
		Short:   "Activate the kill switch: stop every running experiment of every team, and let none run until it is deactivated.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit killswitch activate", "steadybit killswitch activate --yes"),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return killswitch.Activate(ctx, c, yes)
		}),
	}
	activate.Flags().BoolVar(&yes, "yes", false, yesHelp)

	deactivate := &cobra.Command{
		Use:     "deactivate",
		Short:   "Deactivate the kill switch, so that experiments can run again. Those it stopped are not restarted.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit killswitch deactivate"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return killswitch.Deactivate(ctx, c) }),
	}

	cmd.AddCommand(status, activate, deactivate)
	return cmd
}
