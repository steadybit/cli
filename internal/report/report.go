// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package report implements the `report` commands: time series over the tenant, printed as
// the platform sends them.
package report

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/steadybit/cli/api"
	"github.com/steadybit/cli/internal/jsyaml"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/resource"
)

// Filter says which filters a report takes beyond the time range.
type Filter int

const (
	Teams Filter = 1 << iota
	Environments
	Services
	ServiceProperties
	Categories
)

type Kind struct {
	Name, Short string
	Filters     Filter
	// The groupings the report offers; none means it cannot be grouped.
	GroupBy []string
	send    func(ctx context.Context, c *platform.Client, groupBy string, body io.Reader) (*http.Response, error)
}

func (k Kind) Takes(f Filter) bool { return k.Filters&f != 0 }

const contentType = "application/json"

var serviceFilters = Teams | Environments | Services | ServiceProperties | Categories

var Kinds = []Kind{
	{Name: "users", Short: "Users in the tenant over time.", send: func(ctx context.Context, c *platform.Client, _ string, body io.Reader) (*http.Response, error) {
		return c.GetUserCountsWithBody(ctx, contentType, body)
	}},
	{Name: "teams", Short: "Teams in the tenant over time.", send: func(ctx context.Context, c *platform.Client, _ string, body io.Reader) (*http.Response, error) {
		return c.GetTeamCountsWithBody(ctx, contentType, body)
	}},
	{Name: "environments", Short: "Environments in the tenant over time.", send: func(ctx context.Context, c *platform.Client, _ string, body io.Reader) (*http.Response, error) {
		return c.GetEnvironmentCountsWithBody(ctx, contentType, body)
	}},
	{
		Name: "experiments-executed", Short: "Experiment runs over time.", Filters: Teams | Environments | Services,
		GroupBy: []string{"NONE", "STATE", "TRIGGER", "ACTION", "ISSUES_FIXED", "ISSUES_DISCOVERED"},
		send: func(ctx context.Context, c *platform.Client, groupBy string, body io.Reader) (*http.Response, error) {
			params := &api.GetExperimentExecutionsParams{}
			if groupBy != "" {
				g := api.GetExperimentExecutionsParamsGroupBy(groupBy)
				params.GroupBy = &g
			}
			return c.GetExperimentExecutionsWithBody(ctx, params, contentType, body)
		},
	},
	{
		Name: "experiments-created", Short: "Experiments created over time.", Filters: Teams | Environments,
		GroupBy: []string{"NONE", "CREATED_VIA", "ORIGIN"},
		send: func(ctx context.Context, c *platform.Client, groupBy string, body io.Reader) (*http.Response, error) {
			params := &api.GetExperimentCreationsParams{}
			if groupBy != "" {
				g := api.GetExperimentCreationsParamsGroupBy(groupBy)
				params.GroupBy = &g
			}
			return c.GetExperimentCreationsWithBody(ctx, params, contentType, body)
		},
	},
	{Name: "services-distribution", Short: "Services per risk level (low, medium, high) over time.", Filters: serviceFilters,
		send: func(ctx context.Context, c *platform.Client, _ string, body io.Reader) (*http.Response, error) {
			return c.GetRiskDistributionWithBody(ctx, contentType, body)
		}},
	{Name: "services-by-category", Short: "The average risk of services per category over time.", Filters: serviceFilters,
		send: func(ctx context.Context, c *platform.Client, _ string, body io.Reader) (*http.Response, error) {
			return c.GetRiskByCategoryWithBody(ctx, contentType, body)
		}},
	{Name: "services-average", Short: "The average risk of services, 0 to 100, over time.", Filters: serviceFilters,
		send: func(ctx context.Context, c *platform.Client, _ string, body io.Reader) (*http.Response, error) {
			return c.GetAverageRiskWithBody(ctx, contentType, body)
		}},
}

type Options struct {
	From, To, Rollup, GroupBy string
	// Teams by key and environments by name, as users know them; the platform takes ids.
	Teams, Environments, Services, Categories []string
	ServiceProperties                         []string
	File, Type                                string
}

// Today is when a report ends by default. A variable, so that tests can fix it.
var Today = func() time.Time { return time.Now().UTC() }

func date(flag, value string) (string, error) {
	if _, err := time.Parse(time.DateOnly, value); err != nil {
		return "", fmt.Errorf("--%s '%s' is not a date like 2026-09-01.", flag, value)
	}
	return value, nil
}

func Run(ctx context.Context, c *platform.Client, k Kind, o Options) error {
	body := jsyaml.NewMap()
	to := Today().Format(time.DateOnly)
	if o.To != "" {
		var err error
		if to, err = date("to", o.To); err != nil {
			return err
		}
	}
	// 30 days back, unless told otherwise.
	end, _ := time.Parse(time.DateOnly, to)
	from := end.AddDate(0, 0, -30).Format(time.DateOnly)
	if o.From != "" {
		var err error
		if from, err = date("from", o.From); err != nil {
			return err
		}
	}
	body.Set("from", from)
	body.Set("to", to)
	if o.Rollup != "" {
		rollup := strings.ToUpper(o.Rollup)
		if rollup != "DAILY" && rollup != "MONTHLY" {
			return fmt.Errorf("--rollup must be DAILY or MONTHLY, not '%s'.", o.Rollup)
		}
		body.Set("rollup", rollup)
	}
	groupBy := strings.ToUpper(o.GroupBy)
	if groupBy != "" && !contains(k.GroupBy, groupBy) {
		return fmt.Errorf("--group-by must be one of %s, not '%s'.", strings.Join(k.GroupBy, ", "), o.GroupBy)
	}

	if len(o.Teams) > 0 {
		ids, err := teamIDs(ctx, c, o.Teams)
		if err != nil {
			return err
		}
		body.Set("teamIds", ids)
	}
	if len(o.Environments) > 0 {
		ids, err := environmentIDs(ctx, c, o.Environments)
		if err != nil {
			return err
		}
		body.Set("environmentIds", ids)
	}
	if len(o.Services) > 0 {
		body.Set("serviceIds", list(o.Services))
	}
	if len(o.ServiceProperties) > 0 {
		properties := jsyaml.NewMap()
		for _, pair := range o.ServiceProperties {
			key, value, ok := strings.Cut(pair, "=")
			if !ok || key == "" {
				return fmt.Errorf("'%s' is not in the form KEY=VALUE.", pair)
			}
			values, _ := properties.Get(key)
			existing, _ := values.([]any)
			properties.Set(key, append(existing, value))
		}
		body.Set("serviceProperties", properties)
	}
	if len(o.Categories) > 0 {
		body.Set("categoryKeys", list(o.Categories))
	}

	doc, _, err := platform.ReadDocument(k.send(ctx, c, groupBy, resource.Body(body)))
	if err != nil {
		return platform.Failed(err, "Failed to get the %s report", k.Name)
	}
	if err := resource.Output(doc, o.File, o.Type); err != nil {
		return err
	}
	if o.File != "" {
		fmt.Printf("Report %s written to %s.\n", k.Name, o.File)
	}
	return nil
}

func contains(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func list(values []string) []any {
	out := make([]any, len(values))
	for i, v := range values {
		out[i] = v
	}
	return out
}

func teamIDs(ctx context.Context, c *platform.Client, keys []string) ([]any, error) {
	ids := make([]any, len(keys))
	for i, key := range keys {
		var team struct{ ID string }
		resp, err := c.GetTeam(ctx, key)
		if _, err := platform.Decode(resp, err, &team); err != nil {
			if platform.IsStatus(err, http.StatusNotFound) {
				return nil, fmt.Errorf("Team %s not found.", key)
			}
			return nil, platform.Failed(err, "Failed to get team %s", key)
		}
		ids[i] = team.ID
	}
	return ids, nil
}

func environmentIDs(ctx context.Context, c *platform.Client, names []string) ([]any, error) {
	var summaries struct {
		Environments []struct{ ID, Name string } `json:"environments"`
	}
	resp, err := c.GetEnvironments(ctx, &api.GetEnvironmentsParams{})
	if _, err := platform.Decode(resp, err, &summaries); err != nil {
		return nil, platform.Failed(err, "Failed to get the environments")
	}
	ids := make([]any, len(names))
	for i, name := range names {
		for _, e := range summaries.Environments {
			if e.Name == name {
				ids[i] = e.ID
			}
		}
		if ids[i] == nil {
			return nil, fmt.Errorf("Environment %s not found.", name)
		}
	}
	return ids, nil
}
