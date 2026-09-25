// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package target implements the `target` commands.
package target

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/steadybit/cli/api"
	"github.com/steadybit/cli/internal/jsyaml"
	"github.com/steadybit/cli/internal/output"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/resource"
	"github.com/steadybit/cli/internal/table"
)

func notFoundOr(err error, environment, format string) error {
	if platform.IsStatus(err, http.StatusNotFound) {
		return fmt.Errorf("Environment %s not found.", environment)
	}
	return platform.Failed(err, format, environment)
}

func optional(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// printAs writes values as JSON or YAML when a type is given, and returns false otherwise.
func printAs(values []any, datatype string) (bool, error) {
	if !resource.Machine(datatype) {
		return false, nil
	}
	if output.JQ != "" {
		return true, resource.PrintJSONValue([]byte(jsyaml.CompactJSON(values)), datatype)
	}
	resolved, err := output.ResolveDatatype(datatype, "")
	if err != nil {
		return true, err
	}
	if resolved == output.JSON {
		fmt.Println(jsyaml.JSON(values))
	} else {
		fmt.Println(jsyaml.Dump(values))
	}
	return true, nil
}

type QueryOptions struct {
	Environment, TargetType, Query string
	Attributes                     []string
	// At most this many targets; 0 is all of them.
	Limit int
	Type  string
}

// The most targets the platform returns at once.
const maxPageSize = 1000

func Query(ctx context.Context, c *platform.Client, o QueryOptions) error {
	if o.Limit < 0 {
		return errors.New("--limit cannot be negative.")
	}
	params := api.GetTargetsParams{Environment: o.Environment, TargetType: optional(o.TargetType), Query: optional(o.Query)}
	if len(o.Attributes) > 0 {
		params.Attribute = &o.Attributes
	}
	targets := []any{}
	for {
		size := int32(maxPageSize)
		if o.Limit > 0 {
			size = int32(min(o.Limit-len(targets), maxPageSize))
		}
		params.Size = &size
		body, _, err := platform.Read(c.GetTargets(ctx, &params))
		if err != nil {
			return notFoundOr(err, o.Environment, "Failed to get the targets of environment %s")
		}
		value, err := output.ParseValue(body)
		if err != nil {
			return err
		}
		slice, _ := value.(*jsyaml.Map)
		items, _ := slice.Get("items")
		list, _ := items.([]any)
		targets = append(targets, list...)
		hasNext, _ := slice.Get("hasNext")
		next, _ := slice.Get("nextCursor")
		cursor, _ := next.(string)
		if hasNext != true || cursor == "" || len(list) == 0 || (o.Limit > 0 && len(targets) >= o.Limit) {
			break
		}
		params.Cursor = &cursor
	}
	if printed, err := printAs(targets, o.Type); printed {
		return err
	}
	if len(targets) == 0 {
		fmt.Println("No targets found.")
		return nil
	}
	columns := []table.Column{
		{Name: "name", Title: "Name", Alignment: table.Left},
		{Name: "type", Title: "Type", Alignment: table.Left},
	}
	// The attributes asked for become columns; they are what the query was about.
	for _, a := range o.Attributes {
		columns = append(columns, table.Column{Name: "@" + a, Title: a, Alignment: table.Left})
	}
	t := table.New(columns...)
	for _, item := range targets {
		target, _ := item.(*jsyaml.Map)
		name, _ := target.Get("name")
		kind, _ := target.Get("type")
		cells := [][2]string{table.Cell("name", name), table.Cell("type", kind)}
		for _, a := range o.Attributes {
			cells = append(cells, table.Cell("@"+a, strings.Join(attribute(target, a), ", ")))
		}
		t.AddRow(table.Default, cells...)
	}
	t.Print()
	return nil
}

// attribute returns a target's values of a key, which it may have more than once.
func attribute(target *jsyaml.Map, key string) []string {
	attributes, _ := target.Get("attributes")
	list, _ := attributes.([]any)
	var values []string
	for _, a := range list {
		m, _ := a.(*jsyaml.Map)
		if k, _ := m.Get("key"); k == key {
			v, _ := m.Get("value")
			values = append(values, fmt.Sprint(v))
		}
	}
	return values
}

type AttributeOptions struct {
	Environment, TargetType, Action string
	// The attribute whose values to list; empty lists the keys.
	Key  string
	Type string
}

func (o AttributeOptions) check() error {
	if o.TargetType == "" && o.Action == "" {
		return errors.New("Either --target-type or --action must be specified.")
	}
	return nil
}

func AttributeKeys(ctx context.Context, c *platform.Client, o AttributeOptions) error {
	if err := o.check(); err != nil {
		return err
	}
	keys, err := platform.AllPages[string](func(page, size int32) (*http.Response, error) {
		return c.GetTargetAttributeKeys(ctx, &api.GetTargetAttributeKeysParams{
			Environment: o.Environment, TargetType: optional(o.TargetType), ActionId: optional(o.Action), Page: &page, Size: &size,
		})
	})
	if err != nil {
		return notFoundOr(err, o.Environment, "Failed to get the attribute keys of environment %s")
	}
	return printStrings(keys, "Key", "No attribute keys found.", o.Type)
}

func AttributeValues(ctx context.Context, c *platform.Client, o AttributeOptions) error {
	if err := o.check(); err != nil {
		return err
	}
	values, err := platform.AllPages[string](func(page, size int32) (*http.Response, error) {
		return c.GetTargetAttributeValues(ctx, &api.GetTargetAttributeValuesParams{
			Environment: o.Environment, TargetType: optional(o.TargetType), ActionId: optional(o.Action), AttributeKey: o.Key, Page: &page, Size: &size,
		})
	})
	if err != nil {
		return notFoundOr(err, o.Environment, "Failed to get the attribute values of environment %s")
	}
	return printStrings(values, "Value", "No values found for attribute "+o.Key+".", o.Type)
}

func printStrings(values []string, title, none, datatype string) error {
	list := make([]any, len(values))
	for i, v := range values {
		list[i] = v
	}
	if printed, err := printAs(list, datatype); printed {
		return err
	}
	if len(values) == 0 {
		fmt.Println(none)
		return nil
	}
	t := table.New(table.Column{Name: "value", Title: title, Alignment: table.Left})
	for _, v := range values {
		t.AddRow(table.Default, table.Cell("value", v))
	}
	t.Print()
	return nil
}
