// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package cli wires the commands. Names, flags, messages and exit codes follow the
// TypeScript CLI, which pipelines depend on.
package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/internal/experiment"
	"github.com/steadybit/cli/internal/output"
	"github.com/steadybit/cli/internal/platform"
)

// Laid out like the TypeScript CLI's help, which pipelines and the e2e suite read.
const usageTemplate = `Usage: {{if .Runnable}}{{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}{{.CommandPath}} [command]{{end}}
{{if .HasAvailableLocalFlags}}
Options:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableSubCommands}}

Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}
`

func examples(lines ...string) string {
	for i, line := range lines {
		lines[i] = "  $ " + line
	}
	return strings.Join(lines, "\n")
}

// withClient runs a command that talks to the platform. A missing access token is
// reported with the setup help before anything is sent.
func withClient(run func(ctx context.Context, c *platform.Client, args []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client, err := platform.New()
		if err != nil {
			return err
		}
		return run(cmd.Context(), client, args)
	}
}

func newRoot() *cobra.Command {
	root := &cobra.Command{
		Use:           "steadybit",
		Short:         "Command-line interface to interact with the Steadybit API",
		Version:       platform.Version,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: examples(
			"steadybit experiment run -f experiment.yml",
			"steadybit schedule list --team ADM",
			"steadybit experiment --help",
		),
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			platform.Verbose, _ = cmd.Flags().GetBool("verbose")
		},
	}
	root.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")
	root.Flags().BoolP("version", "V", false, "output the version number")
	root.SetVersionTemplate("{{.Version}}\n")
	root.AddCommand(newAccessToken(), newAction(), newAdvice(), newAuditLog(), newConfig(), newEnvironment(), newExecution(), newExperiment(), newHub(), newIntegration(), newKillswitch(), newProperty(), newReport(), newSchedule(), newService(), newServiceProfile(), newTarget(), newTeam(), newTemplate(), newUser())
	// Shell completion is new with the Go CLI; it gets examples like every other command.
	root.InitDefaultCompletionCmd()
	for _, cmd := range root.Commands() {
		if cmd.Name() != "completion" {
			continue
		}
		cmd.Example = examples("steadybit completion zsh > \"${fpath[1]}/_steadybit\"")
		shells := map[string]string{
			"bash":       "source <(steadybit completion bash)",
			"zsh":        `steadybit completion zsh > "${fpath[1]}/_steadybit"`,
			"fish":       "steadybit completion fish > ~/.config/fish/completions/steadybit.fish",
			"powershell": "steadybit completion powershell | Out-String | Invoke-Expression",
		}
		for _, shell := range cmd.Commands() {
			if example, ok := shells[shell.Name()]; ok {
				shell.Example = examples(example)
			}
		}
	}
	for _, cmd := range append(root.Commands(), root) {
		setUsage(cmd)
	}
	return root
}

func setUsage(cmd *cobra.Command) {
	cmd.SetUsageTemplate(usageTemplate)
	for _, sub := range cmd.Commands() {
		setUsage(sub)
	}
}

func Execute() int {
	root := newRoot()
	root.SetArgs(expandVariadic(root, os.Args[1:]))
	err := root.ExecuteContext(context.Background())
	if err == nil {
		return 0
	}
	if errors.Is(err, experiment.ErrIncomplete) {
		return 1 // already reported, with what was missing
	}
	if errors.Is(err, platform.ErrNoAccessToken) {
		fmt.Fprintln(os.Stderr, platform.MissingTokenHelp())
	} else {
		fmt.Fprintln(os.Stderr, output.Red(err.Error()))
	}
	return 1
}
