// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package auditlog implements the `audit-log` command.
package auditlog

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/steadybit/cli/api"
	"github.com/steadybit/cli/internal/jsyaml"
	"github.com/steadybit/cli/internal/output"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/resource"
	"github.com/steadybit/cli/internal/table"
)

type Options struct {
	From, To string
	Type     string
}

// Show prints the entries of the time range; the platform defaults to the last 7 days.
// The endpoint is not paged: it returns the whole range at once.
func Show(ctx context.Context, c *platform.Client, o Options) error {
	from, err := resource.Time("from", o.From)
	if err != nil {
		return err
	}
	to, err := resource.Time("to", o.To)
	if err != nil {
		return err
	}
	body, _, err := platform.Read(c.Find(ctx, &api.FindParams{From: from, To: to}))
	if platform.IsStatus(err, http.StatusForbidden) {
		return errors.New("The audit log needs an admin access token.")
	}
	if err != nil {
		return platform.Failed(err, "Failed to get the audit log")
	}
	value, err := output.ParseValue(body)
	if err != nil {
		return err
	}
	entries, _ := value.([]any)
	if o.Type != "" {
		datatype, err := output.ResolveDatatype(o.Type, "")
		if err != nil {
			return err
		}
		if datatype == output.JSON {
			fmt.Println(jsyaml.JSON(entries))
		} else {
			fmt.Println(jsyaml.Dump(entries))
		}
		return nil
	}
	if len(entries) == 0 {
		fmt.Println("No audit log entries found.")
		return nil
	}
	t := table.New(
		table.Column{Name: "time", Title: "Time", Alignment: table.Left},
		table.Column{Name: "event", Title: "Event", Alignment: table.Left},
		table.Column{Name: "by", Title: "By", Alignment: table.Left},
		table.Column{Name: "team", Title: "Team", Alignment: table.Left},
		table.Column{Name: "environment", Title: "Environment", Alignment: table.Left},
	)
	for _, e := range entries {
		entry, _ := e.(*jsyaml.Map)
		t.AddRow(table.Default, table.Cell("time", field(entry, "eventTime")), table.Cell("event", field(entry, "eventName")),
			table.Cell("by", principal(entry)), table.Cell("team", field(entry, "team", "key")), table.Cell("environment", field(entry, "environment", "name")))
	}
	t.Print()
	return nil
}

// principal names who did it: a user by name, an access token by its name, a batch job as such.
func principal(entry *jsyaml.Map) string {
	switch field(entry, "principal", "principalType") {
	case "USER":
		if name := field(entry, "principal", "name"); name != "" {
			return name
		}
		return field(entry, "principal", "username")
	case "ACCESS_TOKEN":
		return "token " + field(entry, "principal", "name")
	case "BATCH_JOB":
		return "batch job"
	}
	return ""
}

func field(m *jsyaml.Map, path ...string) string {
	var value any = m
	for _, key := range path {
		next, ok := value.(*jsyaml.Map)
		if !ok {
			return ""
		}
		value, _ = next.Get(key)
	}
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}
