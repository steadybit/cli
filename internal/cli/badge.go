// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/v6/internal/badge"
	"github.com/steadybit/cli/v6/internal/platform"
)

func newExperimentBadge() *cobra.Command {
	var o badge.Options
	cmd := &cobra.Command{
		Use:   "badge",
		Short: "Print a status badge of an experiment to paste into a README: its latest run's state, linking to the platform.",
		Long: `Print a status badge of an experiment to paste into a README: its latest run's state, linking to the platform.

With --tag, such as an incident id, the badge shows the latest run of the experiment
having the tag, or, while there is none, invites to create one with the tag.

The badge URL carries no access token: anyone who knows the tenant key can load it,
and it shows the experiment key and the state of its latest run. The link opens the
platform, which asks to log in. Finding the tenant key needs an admin access token;
with any other, pass it with --tenant. With an admin access token, --tenant must be
the token's own tenant.

-t prints the image URL, the link and every snippet at once, so it does not combine
with --format.`,
		Args: cobra.NoArgs,
		Example: examples(
			"steadybit experiment badge -k ADM-1",
			"steadybit experiment badge -k ADM-1 --format html --scale 2",
			`steadybit experiment badge --tag INCIDENT-100 --create-caption "Create experiment for incident 100" --tenant demo`,
		),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return badge.Print(ctx, c, o) }),
	}
	f := cmd.Flags()
	f.StringVarP(&o.Key, "key", "k", "", "The experiment key.")
	f.StringVar(&o.Tag, "tag", "", "Instead of --key: the tag of the experiments the badge is for.")
	f.StringVar(&o.CreateCaption, "create-caption", "", "With --tag: the caption shown while no experiment has the tag. (default: the platform's, \"Create experiment\")")
	f.StringVar(&o.Tenant, "tenant", "", "The tenant key, the tenant= of a platform URL. (default: read from the license)")
	f.IntVar(&o.Scale, "scale", 0, "Scale the badge image by this factor. (default: the platform's, 1)")
	f.StringVar(&o.Format, "format", "", `Print the badge as "markdown", "html", or only the image "url". Not with -t or --jq. (default: markdown)`)
	f.StringVarP(&o.Type, "type", "t", "", `Print the image URL, link and snippets as "json" or "yaml" instead. Not with --format.`)
	cmd.MarkFlagsMutuallyExclusive("key", "tag")
	return cmd
}
