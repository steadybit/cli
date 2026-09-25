// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package action implements the `action` commands.
package action

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/steadybit/cli/api"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/resource"
	"github.com/steadybit/cli/internal/table"
)

type ListOptions struct {
	// Only actions of these kinds, e.g. ATTACK or CHECK; the endpoint cannot filter.
	Kinds []string
}

type summary struct {
	ID, Name, Kind, Category, Technology string
}

// all follows nextPage like platform.AllPages, but this endpoint lists under `actions`.
func all(ctx context.Context, c *platform.Client) ([]summary, error) {
	var actions []summary
	page, size := int32(0), platform.PageSize
	for {
		var body struct {
			Actions  []summary `json:"actions"`
			NextPage *int32    `json:"nextPage"`
		}
		resp, err := c.FindAllActions(ctx, &api.FindAllActionsParams{Page: &page, Size: &size})
		if _, err := platform.Decode(resp, err, &body); err != nil {
			return nil, err
		}
		actions = append(actions, body.Actions...)
		if body.NextPage == nil || *body.NextPage == page {
			return actions, nil
		}
		page = *body.NextPage
	}
}

func List(ctx context.Context, c *platform.Client, o ListOptions) error {
	actions, err := all(ctx, c)
	if err != nil {
		return platform.Failed(err, "Failed to get the actions")
	}
	t := table.New(
		table.Column{Name: "id", Title: "Id", Alignment: table.Left},
		table.Column{Name: "name", Title: "Name", Alignment: table.Left},
		table.Column{Name: "kind", Title: "Kind", Alignment: table.Left},
		table.Column{Name: "category", Title: "Category", Alignment: table.Left},
	)
	rows := 0
	for _, a := range actions {
		if len(o.Kinds) > 0 && !anyEqualFold(o.Kinds, a.Kind) {
			continue
		}
		t.AddRow(table.Default, table.Cell("id", a.ID), table.Cell("name", a.Name), table.Cell("kind", a.Kind), table.Cell("category", a.Category))
		rows++
	}
	if rows == 0 {
		fmt.Println("No actions found.")
		return nil
	}
	t.Print()
	return nil
}

func anyEqualFold(values []string, s string) bool {
	for _, v := range values {
		if strings.EqualFold(v, s) {
			return true
		}
	}
	return false
}

type GetOptions struct {
	ID, File, Type string
}

func Get(ctx context.Context, c *platform.Client, o GetOptions) error {
	doc, _, err := platform.ReadDocument(c.GetAction(ctx, o.ID))
	if platform.IsStatus(err, http.StatusNotFound) {
		return fmt.Errorf("Action %s not found.", o.ID)
	}
	if err != nil {
		return platform.Failed(err, "Failed to get action %s", o.ID)
	}
	if err := resource.Output(doc, o.File, o.Type); err != nil {
		return err
	}
	if o.File != "" {
		fmt.Printf("Action %s written to %s.\n", o.ID, o.File)
	}
	return nil
}
