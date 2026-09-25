// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/internal/advice"
	"github.com/steadybit/cli/internal/platform"
)

func newAdvice() *cobra.Command {
	cmd := &cobra.Command{Use: "advice", Short: "Show/verify advice status."}
	var o advice.Options
	validate := &cobra.Command{
		Use:   "validate-status",
		Short: "Validates the status of one or multiple advice for a given environment and an optional query.",
		Args:  cobra.NoArgs,
		Example: examples(
			"steadybit advice validate-status -e Global",
			`steadybit advice validate-status -e Global -q "k8s.cluster-name=dev-demo and k8s.namespace=steadybit-demo"`,
		),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return advice.ValidateStatus(ctx, c, o)
		}),
	}
	validate.Flags().StringVarP(&o.Environment, "environment", "e", "", "The environment name.")
	validate.Flags().StringVarP(&o.Status, "status", "s", "Implemented", "The expected status of the advice.")
	validate.Flags().StringVarP(&o.Query, "query", "q", "", "(optional) A target query to filter advice by targets.")
	_ = validate.MarkFlagRequired("environment")
	cmd.AddCommand(validate)
	return cmd
}
