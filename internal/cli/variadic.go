// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const variadicAnnotation = "steadybit_variadic"

// variadic marks flags that take several space-separated values, as commander's
// `<values...>` options did: `--team ADM WEBHOOK` and `-f a.yml b.yml`. pflag only takes
// a repeated flag or a comma-separated list, and pipelines written for the TypeScript CLI
// use the space-separated form.
func variadic(cmd *cobra.Command, names ...string) {
	for _, name := range names {
		_ = cmd.Flags().SetAnnotation(name, variadicAnnotation, []string{"true"})
	}
}

// expandVariadic rewrites `--flag a b` into `--flag a --flag b` for the command the
// arguments address. Like commander, it takes values until the next one starting with '-'.
func expandVariadic(root *cobra.Command, args []string) []string {
	cmd, _, err := root.Find(args)
	if err != nil || cmd == nil {
		return args
	}
	isVariadic := func(arg string) (string, bool) {
		var flag *pflag.Flag
		switch {
		case strings.HasPrefix(arg, "--"):
			flag = cmd.Flags().Lookup(strings.TrimPrefix(arg, "--"))
		case strings.HasPrefix(arg, "-") && len(arg) == 2:
			flag = cmd.Flags().ShorthandLookup(arg[1:])
		}
		if flag == nil || flag.Annotations[variadicAnnotation] == nil {
			return "", false
		}
		return arg, true
	}

	var out []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		out = append(out, arg)
		if arg == "--" {
			return append(out, args[i+1:]...)
		}
		name, ok := isVariadic(arg)
		if !ok || i+1 >= len(args) {
			continue
		}
		out = append(out, args[i+1])
		i++
		for i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
			out = append(out, name, args[i+1])
			i++
		}
	}
	return out
}
