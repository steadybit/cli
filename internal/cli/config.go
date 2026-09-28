// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/v6/internal/config"
	"github.com/steadybit/cli/v6/internal/output"
	"github.com/steadybit/cli/v6/internal/prompt"
)

func newConfig() *cobra.Command {
	cmd := &cobra.Command{Use: "config", Short: "Show/modify the CLI configuration and authentication profiles."}
	profile := &cobra.Command{Use: "profile", Short: "Configure authentication profiles."}
	profile.AddCommand(newProfileAdd(), newProfileList("list", "List all configured profiles."),
		newProfileList("ls", "Alias for list."), newProfileSelect(), newProfileRemove())
	cmd.AddCommand(profile, &cobra.Command{
		Use:     "show",
		Short:   "Show the active CLI configuration. Warning: Prints secrets!",
		Example: examples("steadybit config show"),
		RunE: func(*cobra.Command, []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			fmt.Printf("{\n  \"apiAccessToken\": %q,\n  \"baseUrl\": %q\n}\n", cfg.APIAccessToken, cfg.BaseURL)
			return nil
		},
	})
	return cmd
}

const startHelp = `Configuration profiles enable you to use the CLI without repeatedly providing
passwords or having to remember environment variables. Configuration profiles
are stored in ~/.steadybit`

func newProfileAdd() *cobra.Command {
	var p config.Profile
	cmd := &cobra.Command{
		Use:     "add",
		Short:   "Configure a new profile (interactively or via options).",
		Args:    cobra.NoArgs,
		Example: examples("steadybit config profile add", `steadybit config profile add -n prod -t "$STEADYBIT_TOKEN"`),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if p.Name == "" || p.APIAccessToken == "" {
				var err error
				fmt.Println(startHelp)
				fmt.Println()
				if p.Name, err = prompt.Input("Profile name:", "", prompt.NotBlank); err != nil {
					return err
				}
				if p.BaseURL, err = prompt.Input("Base URL of the Steadybit server:", config.DefaultBaseURL, prompt.HTTPURL); err != nil {
					return err
				}
				fmt.Printf("\nThe CLI will need an API access token of %s to communicate with\nthe Steadybit servers. You can generate one through the following URL:\n\n          %s/settings/api-tokens\n\n",
					output.Bold("type team"), strings.TrimSuffix(p.BaseURL, "/"))
				if p.APIAccessToken, err = prompt.Password("API access token:", prompt.NotBlank); err != nil {
					return err
				}
			}
			if err := config.AddProfile(p); err != nil {
				return err
			}
			fmt.Printf("\n%s You can now start using the CLI. For example, you could start\nto run your first experiment via:\n\n                   %s\n",
				output.Green("Done!"), output.Bold("steadybit experiment run -k <your-key>"))
			return nil
		},
	}
	cmd.Flags().StringVarP(&p.Name, "name", "n", "", "Name of the profile")
	cmd.Flags().StringVarP(&p.BaseURL, "baseUrl", "b", config.DefaultBaseURL, "Base URL to be used")
	cmd.Flags().StringVarP(&p.APIAccessToken, "token", "t", "", "Team API token")
	return cmd
}

func newProfileList(use, short string) *cobra.Command {
	return &cobra.Command{
		Use:     use,
		Short:   short,
		Args:    cobra.NoArgs,
		Example: examples("steadybit config profile " + use),
		RunE: func(*cobra.Command, []string) error {
			profiles, err := config.Profiles()
			if err != nil {
				return err
			}
			active, err := config.ActiveProfile()
			if err != nil {
				return err
			}
			for _, p := range profiles {
				if active != nil && p.Name == active.Name {
					fmt.Printf("* %s\n", output.Green(p.Name))
				} else {
					fmt.Printf("  %s\n", p.Name)
				}
			}
			return nil
		},
	}
}

func chooseProfile(message string, args []string) (string, error) {
	if len(args) == 1 {
		return args[0], nil
	}
	profiles, err := config.Profiles()
	if err != nil {
		return "", err
	}
	if len(profiles) == 0 {
		return "", fmt.Errorf("no profiles configured")
	}
	for i, p := range profiles {
		fmt.Printf("  %d) %s\n", i+1, p.Name)
	}
	answer, err := prompt.Input(message, "", func(v string) error {
		for i, p := range profiles {
			if v == p.Name || v == fmt.Sprint(i+1) {
				return nil
			}
		}
		return fmt.Errorf("choose one of the profiles above")
	})
	if err != nil {
		return "", err
	}
	for i, p := range profiles {
		if answer == fmt.Sprint(i+1) {
			return p.Name, nil
		}
	}
	return answer, nil
}

func newProfileSelect() *cobra.Command {
	return &cobra.Command{
		Use:     "select [name]",
		Short:   "Interactively change the currently active profile.",
		Args:    cobra.MaximumNArgs(1),
		Example: examples("steadybit config profile select", "steadybit config profile select prod"),
		RunE: func(_ *cobra.Command, args []string) error {
			name, err := chooseProfile("Profile to activate:", args)
			if err != nil {
				return err
			}
			return config.SetActiveProfile(name)
		},
	}
}

func newProfileRemove() *cobra.Command {
	return &cobra.Command{
		Use:     "remove [name]",
		Short:   "Interactively remove an existing profile.",
		Args:    cobra.MaximumNArgs(1),
		Example: examples("steadybit config profile remove", "steadybit config profile remove old"),
		RunE: func(_ *cobra.Command, args []string) error {
			name, err := chooseProfile("Profile to remove:", args)
			if err != nil {
				return err
			}
			return config.RemoveProfile(name)
		},
	}
}
