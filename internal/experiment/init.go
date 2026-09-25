// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package experiment

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/steadybit/cli/api"
	"github.com/steadybit/cli/internal/jsyaml"
	"github.com/steadybit/cli/internal/output"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/prompt"
	"golang.org/x/term"
)

type InitOptions struct {
	Template    string
	Team        string
	Environment string
	File        string
}

// Interactive reports whether questions can be asked. Tests replace it.
var Interactive = func() bool { return term.IsTerminal(int(os.Stdin.Fd())) }

type templateDetails struct {
	Title        string `json:"templateTitle"`
	Placeholders []struct {
		Key         string `json:"key"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"placeholders"`
}

// Init walks through creating an experiment from a template: which template, its
// placeholders, the team and environment. It creates the experiment and writes it to a
// file, ready for `run -f` and for Git.
func Init(ctx context.Context, c *platform.Client, o InitOptions) error {
	if !Interactive() {
		return errors.New("`experiment init` asks questions and needs a terminal. In scripts, use `experiment apply --template`.")
	}
	var err error
	if o.Template == "" {
		if o.Template, err = chooseTemplate(ctx, c); err != nil {
			return err
		}
	}
	var id openapi_types.UUID
	if err := id.UnmarshalText([]byte(o.Template)); err != nil {
		return fmt.Errorf("Experiment template %s not found.", o.Template)
	}
	var template templateDetails
	resp, err := c.GetExperimentTemplate(ctx, id)
	if _, err := platform.Decode(resp, err, &template); err != nil {
		if platform.IsStatus(err, http.StatusNotFound) {
			return fmt.Errorf("Experiment template %s not found.", o.Template)
		}
		return platform.Failed(err, "Failed to get experiment template %s", o.Template)
	}

	fmt.Printf("\n%s\n\n", output.Bold(template.Title))
	placeholders := jsyaml.NewMap()
	for _, p := range template.Placeholders {
		if p.Description != "" {
			fmt.Println(strings.TrimSpace(p.Description))
		}
		value, err := prompt.Input(fmt.Sprintf("%s (%s):", p.Name, p.Key), "", prompt.NotBlank)
		if err != nil {
			return err
		}
		placeholders.Set(p.Key, value)
	}
	if o.Team == "" {
		if o.Team, err = prompt.Input("Team key:", "", prompt.NotBlank); err != nil {
			return err
		}
	}
	if o.Environment == "" {
		if o.Environment, err = prompt.Input("Environment:", "Global", prompt.NotBlank); err != nil {
			return err
		}
	}
	if o.File == "" {
		suggested := slug(template.Title) + ".yml"
		if o.File, err = prompt.Input("Write it to:", suggested, prompt.NotBlank); err != nil {
			return err
		}
	}

	request, err := createRequest(TemplateOptions{Template: o.Template, Team: o.Team, Environment: o.Environment, Placeholder: placeholders})
	if err != nil {
		return err
	}
	reset := true
	_, created, err := platform.Read(c.CreateExperimentByTemplate(ctx, id, &api.CreateExperimentByTemplateParams{ResetProperties: &reset}, request))
	if err != nil {
		return platform.Failed(err, "Failed to create the experiment from template %s", o.Template)
	}
	key := keyFromLocation(created)
	if err := Get(ctx, c, GetOptions{Key: key, File: o.File}); err != nil {
		return err
	}
	fmt.Printf("\nExperiment %s created. Run it with:\n\n    steadybit experiment run -f %s\n", key, o.File)
	return nil
}

func chooseTemplate(ctx context.Context, c *platform.Client) (string, error) {
	search, err := prompt.Input("Search templates:", "", func(string) error { return nil })
	if err != nil {
		return "", err
	}
	params := &api.GetExperimentTemplatesParams{}
	if search != "" {
		params.FreeTextPhrases = &[]string{search}
	}
	var list struct {
		Templates []struct {
			ID    string `json:"id"`
			Title string `json:"templateTitle"`
		} `json:"templates"`
	}
	resp, err := c.GetExperimentTemplates(ctx, params)
	if _, err := platform.Decode(resp, err, &list); err != nil {
		return "", platform.Failed(err, "Failed to get the experiment templates")
	}
	if len(list.Templates) == 0 {
		return "", fmt.Errorf("No experiment templates match '%s'.", search)
	}
	const shown = 20
	for i, t := range list.Templates {
		if i == shown {
			fmt.Printf("  … and %d more; search more precisely to see them.\n", len(list.Templates)-shown)
			break
		}
		fmt.Printf("  %2d) %s\n", i+1, t.Title)
	}
	choice, err := prompt.Input("Template:", "1", func(v string) error {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > min(len(list.Templates), shown) {
			return fmt.Errorf("choose one of the numbers above")
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	n, _ := strconv.Atoi(choice)
	return list.Templates[n-1].ID, nil
}

var nonWord = regexp.MustCompile(`[^a-z0-9]+`)

func slug(s string) string {
	s = strings.Trim(nonWord.ReplaceAllString(strings.ToLower(s), "-"), "-")
	if s == "" {
		return "experiment"
	}
	return s
}
