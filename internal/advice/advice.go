// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package advice implements `advice validate-status`.
package advice

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/steadybit/cli/api"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/table"
)

type Options struct {
	Environment string
	Query       string
	Status      string
}

type item struct {
	Target struct {
		Reference string `json:"reference"`
	} `json:"target"`
	Advice struct {
		Label  string `json:"label"`
		Status string `json:"status"`
	} `json:"advice"`
}

type page struct {
	TotalItems int    `json:"totalItems"`
	NextOffset *int64 `json:"nextOffset"`
	Items      []item `json:"items"`
}

var separators = regexp.MustCompile(`[\s_-]+`)

// The platform reports IMPLEMENTED while --status defaults to Implemented; case and the
// separator are ignored, so `action needed` matches ACTION_NEEDED too.
func sameStatus(reported, expected string) bool {
	normalise := func(s string) string { return separators.ReplaceAllString(strings.ToLower(strings.TrimSpace(s)), "_") }
	return normalise(reported) == normalise(expected)
}

func fetchAll(ctx context.Context, c *platform.Client, o Options) ([]item, error) {
	var all []item
	offset := int64(0)
	for {
		request := api.GetAdviceApiRequestAO{EnvironmentName: o.Environment, Offset: &offset}
		if o.Query != "" {
			request.Query = &o.Query
		}
		resp, err := c.GetTargetAdviceSummary(ctx, request)
		var p page
		if _, err := platform.Decode(resp, err, &p); err != nil {
			return nil, platform.Failed(err, "Failed to fetch advice status. HTTP request failed.")
		}
		if len(p.Items) > 0 {
			all = append(all, p.Items...)
			fmt.Printf("Fetched %d of %d matching advice.\n", len(all), p.TotalItems)
		} else {
			fmt.Println("No matching advice.")
		}
		if p.NextOffset == nil || *p.NextOffset <= 0 {
			return all, nil
		}
		offset = *p.NextOffset
	}
}

func ValidateStatus(ctx context.Context, c *platform.Client, o Options) error {
	all, err := fetchAll(ctx, c, o)
	if err != nil || len(all) == 0 {
		return err
	}
	errors := 0
	t := table.New()
	for _, a := range all {
		color := table.Green
		if !sameStatus(a.Advice.Status, o.Status) {
			errors++
			color = table.Red
		}
		t.AddRow(color, table.Cell("target", a.Target.Reference), table.Cell("advice", a.Advice.Label), table.Cell("status", a.Advice.Status))
	}
	t.Print()
	if errors > 0 {
		return fmt.Errorf("%d of %d advice did not match the expected status.", errors, len(all))
	}
	return nil
}
