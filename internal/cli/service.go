// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/v6/internal/gitops"
	"github.com/steadybit/cli/v6/internal/platform"
	"github.com/steadybit/cli/v6/internal/resource"
	"github.com/steadybit/cli/v6/internal/service"
	"github.com/steadybit/cli/v6/internal/serviceprofile"
)

const (
	serviceID = "019cd80d-a4c9-775b-bdf8-2672a280ce7c"
	profileID = "019eacd7-fb2c-733a-bed5-99a935323db5"
)

func idFlag(cmd *cobra.Command, id *string, help string) {
	cmd.Flags().StringVarP(id, "id", "i", "", help)
	_ = cmd.MarkFlagRequired("id")
}

func newService() *cobra.Command {
	cmd := &cobra.Command{Use: "service", Short: "Manage services, their experiments, variables and risk."}

	var l service.ListOptions
	list := &cobra.Command{
		Use:     "list",
		Short:   "List services. Filters of the same kind match any of the given values.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit service list", "steadybit service list --team ADM --environment Global", "steadybit service list --jq '.[].id'"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return service.List(ctx, c, l) }),
	}
	list.Flags().StringArrayVar(&l.Teams, "team", nil, "Only list services of these teams, by team key.")
	list.Flags().StringArrayVar(&l.Environments, "environment", nil, "Only list services in these environments.")
	list.Flags().StringArrayVar(&l.Experiments, "experiment", nil, "Only list services these experiments are linked to.")
	list.Flags().StringVarP(&l.Type, "type", "t", "", resource.ListTypeHelp)
	variadic(list, "team", "environment", "experiment")

	var g service.GetOptions
	get := &cobra.Command{
		Use:     "get",
		Short:   "Get a service. Output is written to file or stdout.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit service get -i " + serviceID + " -f service.yml"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return service.Get(ctx, c, g) }),
	}
	idFlag(get, &g.ID, "The service id.")
	get.Flags().StringVarP(&g.File, "file", "f", "", "The path to write the service to.")
	get.Flags().StringVarP(&g.Type, "type", "t", "", typeHelp)

	var a service.ApplyOptions
	apply := &cobra.Command{
		Use:     "apply",
		Short:   "Create or update services from files. A file without an id creates a service, and the new id is written back to it.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit service apply -f service.yml", "steadybit service apply -f ./services -R"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return service.Apply(ctx, c, a) }),
	}
	apply.Flags().StringArrayVarP(&a.Files, "file", "f", nil, "The path to the service file or a directory containing multiple files.")
	apply.Flags().BoolVarP(&a.Recursive, "recursive", "R", false, "Process the directory used in -f, --file recursively.")
	apply.Flags().BoolVar(&a.DeleteExperiments, "delete-experiments", false, "When the service profile changes, delete provided experiments whose templates the new profile does not contain. Without it, such a change is refused.")
	_ = apply.MarkFlagRequired("file")
	variadic(apply, "file")
	dryRun(apply, gitops.Service, &a.Files, &a.Recursive)

	var deleteID string
	del := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a service.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit service delete -i " + serviceID),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return service.Delete(ctx, c, deleteID)
		}),
	}
	idFlag(del, &deleteID, "The service id.")

	var r service.RiskOptions
	var failAbove int
	risk := &cobra.Command{
		Use:     "risk",
		Short:   "Show the risk score of a service, overall, per category and per experiment.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit service risk -i "+serviceID, "steadybit service risk -i "+serviceID+" --fail-above 50"),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return service.Risk(ctx, c, r)
		}),
	}
	idFlag(risk, &r.ID, "The service id.")
	risk.Flags().StringVarP(&r.Type, "type", "t", "", `Print the raw risk as "json" or "yaml" instead of tables.`)
	risk.Flags().IntVar(&failAbove, "fail-above", 0, "Exit with a non-zero status when the overall risk is above this score.")
	risk.PreRun = func(cmd *cobra.Command, _ []string) {
		if cmd.Flags().Changed("fail-above") {
			r.FailAbove = &failAbove
		}
	}

	experiments := &cobra.Command{Use: "experiment", Short: "Manage the experiments of a service."}
	var el service.ExperimentListOptions
	elist := &cobra.Command{
		Use:     "list",
		Short:   "List the experiments of a service: those provided by its service profile, created or not, and custom ones linked to it.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit service experiment list -i "+serviceID, "steadybit service experiment list -i "+serviceID+" --type custom"),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return service.ListExperiments(ctx, c, el)
		}),
	}
	idFlag(elist, &el.ID, "The service id.")
	elist.Flags().StringArrayVar(&el.Categories, "category", nil, "Only list experiments in these categories.")
	elist.Flags().StringArrayVar(&el.Types, "type", nil, `Only list "provided" or "custom" experiments.`)
	elist.Flags().StringVar(&el.Type, "output", "", resource.ListTypeHelp)
	variadic(elist, "category", "type")

	var p service.ProvideOptions
	placeholders := newKeyValues()
	p.Placeholder = placeholders.values
	var noReset bool
	provide := &cobra.Command{
		Use:     "provide",
		Short:   "Create or update a provided experiment of a service from one of its service profile's templates.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit service experiment provide -i " + serviceID + " --template d7e65100-1d20-4980-be87-c351704910b8 -p REPLICAS=3"),
		PreRun:  func(*cobra.Command, []string) { p.ResetProperties = !noReset },
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return service.Provide(ctx, c, p) }),
	}
	idFlag(provide, &p.ID, "The service id.")
	provide.Flags().StringVar(&p.Template, "template", "", "The template, which must be part of the service profile.")
	provide.Flags().StringVarP(&p.Experiment, "experiment", "k", "", "Update this existing provided experiment instead of creating one.")
	provide.Flags().VarP(placeholders, "placeholder", "p", "A placeholder value. Repeat for more.")
	provide.Flags().StringVar(&p.PlaceholdersFile, "placeholders", "", "A YAML/JSON file mapping placeholder keys to values. -p overrides entries.")
	provide.Flags().BoolVar(&noReset, "no-reset-properties", false, "Keep the properties of an existing experiment instead of resetting them to the template.")
	_ = provide.MarkFlagRequired("template")

	var lk struct{ id, experiment, category string }
	link := &cobra.Command{
		Use:     "link",
		Short:   "Link an existing experiment to a service as a custom experiment.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit service experiment link -i " + serviceID + " -k ADM-1 --category Redundancy"),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return service.Link(ctx, c, lk.id, lk.experiment, lk.category)
		}),
	}
	idFlag(link, &lk.id, "The service id.")
	link.Flags().StringVarP(&lk.experiment, "experiment", "k", "", "The experiment to link.")
	link.Flags().StringVar(&lk.category, "category", "", "The category to link it in.")
	_ = link.MarkFlagRequired("experiment")
	_ = link.MarkFlagRequired("category")

	var ul struct{ id, experiment string }
	unlink := &cobra.Command{
		Use:     "unlink",
		Short:   "Remove a custom experiment from a service. The experiment itself is kept.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit service experiment unlink -i " + serviceID + " -k ADM-1"),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return service.Unlink(ctx, c, ul.id, ul.experiment)
		}),
	}
	idFlag(unlink, &ul.id, "The service id.")
	unlink.Flags().StringVarP(&ul.experiment, "experiment", "k", "", "The experiment to unlink.")
	_ = unlink.MarkFlagRequired("experiment")
	experiments.AddCommand(elist, provide, link, unlink)

	variable := &cobra.Command{Use: "variable", Short: "Manage the variables of a service."}
	var vg service.VariableGetOptions
	vget := &cobra.Command{
		Use:     "get",
		Short:   "Print the variables of a service.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit service variable get -i " + serviceID),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return service.GetVariables(ctx, c, vg)
		}),
	}
	idFlag(vget, &vg.ID, "The service id.")
	vget.Flags().StringVarP(&vg.Type, "type", "t", "yaml", `The output format ("json" or "yaml").`)
	var vs service.VariableSetOptions
	vset := &cobra.Command{
		Use:   "set [KEY=VALUE...]",
		Short: "Set variables of a service, keeping the others. With --replace, the given variables become the only ones.",
		Example: examples(
			"steadybit service variable set -i "+serviceID+" endpoint=http://shop.internal region=eu",
			"steadybit service variable set -i "+serviceID+" -f variables.yml --replace",
		),
		RunE: withClient(func(ctx context.Context, c *platform.Client, args []string) error {
			return service.SetVariables(ctx, c, args, vs)
		}),
	}
	idFlag(vset, &vs.ID, "The service id.")
	vset.Flags().StringVarP(&vs.File, "file", "f", "", "A YAML/JSON file mapping variable names to values, which may be lists or select expressions.")
	vset.Flags().BoolVar(&vs.Replace, "replace", false, "Remove every variable not given.")
	variable.AddCommand(vget, vset)

	cmd.AddCommand(list, get, apply, newDiff(gitops.Service, "service", "service.yml"), del, risk, experiments, variable)
	return cmd
}

func newServiceProfile() *cobra.Command {
	cmd := &cobra.Command{Use: "service-profile", Short: "Manage the service profiles that provide experiments to services."}

	var l serviceprofile.ListOptions
	list := &cobra.Command{
		Use:     "list",
		Short:   "List service profiles.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit service-profile list", "steadybit service-profile list --origin custom"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return serviceprofile.List(ctx, c, l) }),
	}
	list.Flags().StringVar(&l.Name, "name", "", "Only list profiles whose name contains this.")
	list.Flags().StringArrayVar(&l.Origins, "origin", nil, `Only list "provided" or "custom" profiles.`)
	list.Flags().BoolVar(&l.Default, "default", false, "Only list the default profile.")
	list.Flags().StringVarP(&l.Type, "type", "t", "", resource.ListTypeHelp)
	variadic(list, "origin")

	var g serviceprofile.GetOptions
	get := &cobra.Command{
		Use:     "get",
		Short:   "Get a service profile. Output is written to file or stdout.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit service-profile get -i " + profileID + " -f profile.yml"),
		RunE:    withClient(func(ctx context.Context, c *platform.Client, _ []string) error { return serviceprofile.Get(ctx, c, g) }),
	}
	idFlag(get, &g.ID, "The service profile id.")
	get.Flags().StringVarP(&g.File, "file", "f", "", "The path to write the service profile to.")
	get.Flags().StringVarP(&g.Type, "type", "t", "", typeHelp)

	var a serviceprofile.ApplyOptions
	apply := &cobra.Command{
		Use:     "apply",
		Short:   "Create or update service profiles from files. A file without an id creates a profile, and the new id is written back to it.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit service-profile apply -f profile.yml"),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return serviceprofile.Apply(ctx, c, a)
		}),
	}
	apply.Flags().StringArrayVarP(&a.Files, "file", "f", nil, "The path to the service profile file or a directory containing multiple files.")
	apply.Flags().BoolVarP(&a.Recursive, "recursive", "R", false, "Process the directory used in -f, --file recursively.")
	apply.Flags().BoolVar(&a.DeleteExperiments, "delete-experiments", false, "Delete the provided experiments of services that use templates removed from the profile.")
	_ = apply.MarkFlagRequired("file")
	variadic(apply, "file")
	dryRun(apply, gitops.ServiceProfile, &a.Files, &a.Recursive)

	var deleteID string
	del := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a custom service profile.",
		Args:    cobra.NoArgs,
		Example: examples("steadybit service-profile delete -i " + profileID),
		RunE: withClient(func(ctx context.Context, c *platform.Client, _ []string) error {
			return serviceprofile.Delete(ctx, c, deleteID)
		}),
	}
	idFlag(del, &deleteID, "The service profile id.")

	cmd.AddCommand(list, get, apply, newDiff(gitops.ServiceProfile, "service-profile", "profile.yml"), del)
	return cmd
}
