// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/steadybit/cli/v6/api"
	"github.com/steadybit/cli/v6/internal/config"
	"github.com/steadybit/cli/v6/internal/platform"
)

// Shell completion offers what exists on the platform: experiment keys, team keys and
// the ids of templates, schedules, services and profiles, each with its name. It must
// never print an error or keep the shell waiting, so any failure offers nothing.

type completer func(ctx context.Context, c *platform.Client, toComplete string) ([]string, cobra.ShellCompDirective)

func remote(complete completer) cobra.CompletionFunc {
	return func(cmd *cobra.Command, _ []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		c, err := platform.New()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		values, directive := complete(ctx, c, toComplete)
		return values, directive | cobra.ShellCompDirectiveNoFileComp
	}
}

func matching(values []string, prefix string) []string {
	var out []string
	for _, v := range values {
		if strings.HasPrefix(strings.ToLower(v), strings.ToLower(prefix)) {
			out = append(out, v)
		}
	}
	return out
}

func completeTeams(ctx context.Context, c *platform.Client, prefix string) ([]string, cobra.ShellCompDirective) {
	var teams struct {
		Teams []struct{ Key, Name string } `json:"teams"`
	}
	all := false
	resp, err := c.GetTeams(ctx, &api.GetTeamsParams{OnlyAccessible: &all})
	if _, err := platform.Decode(resp, err, &teams); err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var values []string
	for _, t := range teams.Teams {
		values = append(values, t.Key+"\t"+t.Name)
	}
	return matching(values, prefix), cobra.ShellCompDirectiveDefault
}

// Keys are TEAM-number, so the team is completed first and then only that team's
// experiments are listed: a tenant can hold thousands.
func completeExperimentKeys(ctx context.Context, c *platform.Client, prefix string) ([]string, cobra.ShellCompDirective) {
	team, _, hasTeam := strings.Cut(prefix, "-")
	if !hasTeam {
		teams, directive := completeTeams(ctx, c, prefix)
		for i, t := range teams {
			key, name, _ := strings.Cut(t, "\t")
			teams[i] = key + "-\t" + name
		}
		return teams, directive | cobra.ShellCompDirectiveNoSpace
	}
	var list struct {
		Experiments []struct{ Key, Name string } `json:"experiments"`
	}
	teams := []string{strings.ToUpper(team)}
	resp, err := c.GetExperiments(ctx, &api.GetExperimentsParams{Team: &teams})
	if _, err := platform.Decode(resp, err, &list); err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var values []string
	for _, e := range list.Experiments {
		values = append(values, e.Key+"\t"+e.Name)
	}
	return matching(values, prefix), cobra.ShellCompDirectiveDefault
}

func completeTemplates(ctx context.Context, c *platform.Client, prefix string) ([]string, cobra.ShellCompDirective) {
	var list struct {
		Templates []struct {
			ID    string `json:"id"`
			Title string `json:"templateTitle"`
		} `json:"templates"`
	}
	resp, err := c.GetExperimentTemplates(ctx, nil)
	if _, err := platform.Decode(resp, err, &list); err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var values []string
	for _, t := range list.Templates {
		values = append(values, t.ID+"\t"+t.Title)
	}
	return matching(values, prefix), cobra.ShellCompDirectiveDefault
}

func completeSchedules(ctx context.Context, c *platform.Client, prefix string) ([]string, cobra.ShellCompDirective) {
	var schedules []struct {
		ID            string  `json:"id"`
		ExperimentKey string  `json:"experimentKey"`
		Cron          *string `json:"cron"`
		StartAt       *string `json:"startAt"`
	}
	resp, err := c.GetAllSchedulesV2(ctx, nil)
	if _, err := platform.Decode(resp, err, &schedules); err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var values []string
	for _, s := range schedules {
		when := ""
		if s.Cron != nil {
			when = *s.Cron
		} else if s.StartAt != nil {
			when = *s.StartAt
		}
		values = append(values, fmt.Sprintf("%s\t%s %s", s.ID, s.ExperimentKey, when))
	}
	return matching(values, prefix), cobra.ShellCompDirectiveDefault
}

func completePaged(fetch func(ctx context.Context, c *platform.Client, page, size int32) (*http.Response, error)) completer {
	return func(ctx context.Context, c *platform.Client, prefix string) ([]string, cobra.ShellCompDirective) {
		items, err := platform.AllPagesRaw(func(page, size int32) (*http.Response, error) { return fetch(ctx, c, page, size) })
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		var values []string
		for _, raw := range items {
			var item struct{ ID, Name string }
			if json.Unmarshal(raw, &item) == nil {
				values = append(values, item.ID+"\t"+item.Name)
			}
		}
		return matching(values, prefix), cobra.ShellCompDirectiveDefault
	}
}

func completeEnvironments(ctx context.Context, c *platform.Client, prefix string) ([]string, cobra.ShellCompDirective) {
	var list struct {
		Environments []struct{ ID, Name string } `json:"environments"`
	}
	resp, err := c.GetEnvironments(ctx, nil)
	if _, err := platform.Decode(resp, err, &list); err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var values []string
	for _, e := range list.Environments {
		values = append(values, e.ID+"\t"+e.Name)
	}
	return matching(values, prefix), cobra.ShellCompDirectiveDefault
}

var completeServices = completePaged(func(ctx context.Context, c *platform.Client, page, size int32) (*http.Response, error) {
	return c.GetServiceList(ctx, &api.GetServiceListParams{Page: api.PageRequestAO{Page: &page, Size: &size}})
})

var completeProfiles = completePaged(func(ctx context.Context, c *platform.Client, page, size int32) (*http.Response, error) {
	return c.GetProfiles(ctx, &api.GetProfilesParams{Page: api.PageRequestAO{Page: &page, Size: &size}})
})

// Profile names come from the local configuration and need no platform.
func completeProfileNames(*cobra.Command, []string, string) ([]cobra.Completion, cobra.ShellCompDirective) {
	profiles, err := config.Profiles()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var names []string
	for _, p := range profiles {
		names = append(names, p.Name+"\t"+p.BaseURL)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

// registerCompletions attaches completion to flags by what they name, across all
// commands, so a new command gets it without doing anything.
func registerCompletions(root *cobra.Command) {
	_ = root.RegisterFlagCompletionFunc("profile", completeProfileNames)
	var walk func(cmd *cobra.Command)
	walk = func(cmd *cobra.Command) {
		path := cmd.CommandPath()
		group := ""
		if parts := strings.Fields(path); len(parts) > 1 {
			group = parts[1]
		}
		register := func(flag string, complete completer) {
			if cmd.Flags().Lookup(flag) != nil {
				_ = cmd.RegisterFlagCompletionFunc(flag, remote(complete))
			}
		}
		register("team", completeTeams)
		register("template", completeTemplates)
		switch group {
		case "experiment", "execution":
			register("key", completeExperimentKeys)
		case "schedule":
			register("id", completeSchedules)
			register("experiment", completeExperimentKeys)
		case "service":
			register("id", completeServices)
			register("experiment", completeExperimentKeys)
		case "service-profile":
			register("id", completeProfiles)
		case "template":
			register("id", completeTemplates)
		case "team":
			register("key", completeTeams)
		case "environment":
			register("id", completeEnvironments)
		}
		for _, sub := range cmd.Commands() {
			walk(sub)
		}
	}
	walk(root)
}
