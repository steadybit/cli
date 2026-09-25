// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package prompt asks questions on a terminal. Ctrl-C ends the process with 130 and
// without a stack trace, as the TypeScript CLI's prompts did.
package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/steadybit/cli/internal/interrupt"
	"golang.org/x/term"
)

var reader = bufio.NewReader(os.Stdin)

type Validator func(string) error

func NotBlank(value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("a value is required")
	}
	return nil
}

func HTTPURL(value string) error {
	if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
		return fmt.Errorf("please enter an http:// or https:// URL")
	}
	return nil
}

// Input asks until the answer is valid. An empty answer takes the default.
func Input(message, defaultValue string, validate Validator) (string, error) {
	for {
		if defaultValue != "" {
			fmt.Printf("? %s (%s) ", message, defaultValue)
		} else {
			fmt.Printf("? %s ", message)
		}
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		answer := strings.TrimSpace(line)
		if answer == "" {
			answer = defaultValue
		}
		if err := validate(answer); err != nil {
			fmt.Printf("> %s\n", err)
			continue
		}
		return answer, nil
	}
}

// Password reads without echo when stdin is a terminal.
func Password(message string, validate Validator) (string, error) {
	for {
		fmt.Printf("? %s ", message)
		var answer string
		if fd := int(os.Stdin.Fd()); term.IsTerminal(fd) {
			state, err := term.GetState(fd)
			if err != nil {
				return "", err
			}
			// Reading a password turns echo off. An interrupt meanwhile has to turn it back
			// on, or the user is left with a terminal that shows nothing typed.
			pop := interrupt.Push(func(os.Signal) { _ = term.Restore(fd, state) })
			bytes, err := term.ReadPassword(fd)
			_ = term.Restore(fd, state)
			pop()
			fmt.Println()
			if err != nil {
				return "", err
			}
			answer = string(bytes)
		} else {
			line, err := reader.ReadString('\n')
			if err != nil {
				return "", err
			}
			answer = strings.TrimSpace(line)
		}
		if err := validate(answer); err != nil {
			fmt.Printf("> %s\n", err)
			continue
		}
		return answer, nil
	}
}

// Confirm asks a yes/no question. Without a terminal it answers nonInteractive.
func Confirm(message string, defaultYes, nonInteractive bool) (bool, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return nonInteractive, nil
	}
	hint := "y/N"
	if defaultYes {
		hint = "Y/n"
	}
	fmt.Printf("? %s (%s) ", message, hint)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "":
		return defaultYes, nil
	case "y", "yes":
		return true, nil
	default:
		return false, nil
	}
}
