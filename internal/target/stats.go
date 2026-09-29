// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package target

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/steadybit/cli/v6/api"
	"github.com/steadybit/cli/v6/internal/jsyaml"
	"github.com/steadybit/cli/v6/internal/output"
	"github.com/steadybit/cli/v6/internal/platform"
	"github.com/steadybit/cli/v6/internal/resource"
	"github.com/steadybit/cli/v6/internal/table"
)

type StatsOptions struct {
	Query string
	Type  string
}

// Stats prints how many targets of each type the platform knows. The platform counts
// over the whole tenant: it takes a query but no environment.
func Stats(ctx context.Context, c *platform.Client, o StatsOptions) error {
	if _, err := output.ResolveDatatype(o.Type, ""); err != nil {
		return err
	}
	var resp *http.Response
	var err error
	if o.Query == "" {
		resp, err = c.GetTargetsStats(ctx)
	} else {
		resp, err = c.GetTargetsStats1(ctx, api.TargetStatsRequest{Query: &o.Query})
	}
	body, _, err := platform.Read(resp, err)
	if err != nil {
		return platform.Failed(err, "Failed to get the target statistics")
	}
	if resource.Machine(o.Type) {
		return resource.PrintJSONValue(body, o.Type)
	}
	value, err := output.ParseValue(body)
	if err != nil {
		return err
	}
	counts, _ := value.(*jsyaml.Map)
	if counts == nil || counts.Len() == 0 {
		fmt.Println("No targets found.")
		return nil
	}
	// The platform sends the types in no particular order.
	types := counts.Keys()
	sort.Strings(types)
	t := table.New(
		table.Column{Name: "type", Title: "Target type", Alignment: table.Left},
		table.Column{Name: "count", Title: "Targets", Alignment: table.Right},
	)
	for _, kind := range types {
		count, _ := counts.Get(kind)
		t.AddRow(table.Default, table.Cell("type", kind), table.Cell("count", count))
	}
	t.Print()
	return nil
}
