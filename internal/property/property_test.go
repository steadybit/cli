// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package property_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/steadybit/cli/v6/internal/platformtest"
	"github.com/steadybit/cli/v6/internal/property"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

const id = "0190d7b2-9e8f-7c6d-b5a4-3f2e1d0c9b8a"

func TestListDefinitionsWalksEveryPage(t *testing.T) {
	p := platformtest.New(t)
	p.Handle("GET /api/properties/definitions", func(r platformtest.Request) platformtest.Reply {
		if r.Query["page"][0] == "0" {
			return platformtest.Reply{JSON: map[string]any{"items": []any{map[string]any{"key": "COLOR", "label": "Color", "dataType": "ENUM",
				"enumValues": []any{"a", "b", "c", "d", "e", "f", "g"}}}, "nextPage": 1}}
		}
		return platformtest.Reply{JSON: map[string]any{"items": []any{map[string]any{"key": "OWNER", "label": "Owner", "dataType": "STRING"}}}}
	})

	out, err := platformtest.Stdout(t, func() error { return property.ListDefinitions(ctx, p.Client, "") })

	require.NoError(t, err)
	assert.Contains(t, out, "│ COLOR │ Color │ ENUM   │ a, b, c, d, e and 2 more │")
	assert.Contains(t, out, "│ OWNER │ Owner │ STRING │                          │")
	assert.Equal(t, []string{"100"}, p.Requests("GET /api/properties/definitions")[0].Query["size"])
}

func TestDefinitionRoundTripByKey(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/properties/definitions/COLOR", platformtest.Reply{Body: `{"key":"COLOR","label":"Color","dataType":"ENUM","enumValues":["red"],"version":3}`})
	p.Reply("POST /api/properties/definitions", platformtest.Reply{Body: `{"key":"COLOR"}`})
	file := filepath.Join(t.TempDir(), "color.yml")

	out, err := platformtest.Stdout(t, func() error {
		if err := property.GetDefinition(ctx, p.Client, property.GetDefinitionOptions{Key: "COLOR", File: file}); err != nil {
			return err
		}
		return property.ApplyDefinitions(ctx, p.Client, property.ApplyDefinitionOptions{Files: []string{file}, DeleteValues: true})
	})

	require.NoError(t, err)
	assert.Contains(t, out, "Property definition COLOR updated.")
	content, _ := os.ReadFile(file)
	assert.Equal(t, "key: COLOR\nlabel: Color\ndataType: ENUM\nenumValues:\n  - red\n", string(content))
	sent := p.Requests("POST /api/properties/definitions")[0]
	assert.Equal(t, []string{"true"}, sent.Query["deleteValues"])
	assert.Equal(t, map[string]any{"key": "COLOR", "label": "Color", "dataType": "ENUM", "enumValues": []any{"red"}}, sent.JSON(t))
}

func TestDeleteDefinition(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("DELETE /api/properties/definitions/COLOR", platformtest.Reply{})
	p.Reply("DELETE /api/properties/definitions/NOPE", platformtest.Reply{Status: http.StatusNotFound})

	out, err := platformtest.Stdout(t, func() error {
		return property.DeleteDefinition(ctx, p.Client, property.DeleteDefinitionOptions{Key: "COLOR", Associations: true, Yes: true})
	})

	require.NoError(t, err)
	assert.Equal(t, "Property definition COLOR deleted.\n", out)
	assert.Equal(t, []string{"true"}, p.Requests("DELETE /api/properties/definitions/COLOR")[0].Query["deleteAssociations"])
	assert.EqualError(t, property.DeleteDefinition(ctx, p.Client, property.DeleteDefinitionOptions{Key: "NOPE", Yes: true}), "Property definition NOPE not found.")
}

func TestListAssociationsWithFilters(t *testing.T) {
	p := platformtest.New(t)
	p.Handle("GET /api/properties/associations", func(r platformtest.Request) platformtest.Reply {
		if r.Query["page"][0] == "0" {
			return platformtest.Reply{JSON: map[string]any{"items": []any{
				map[string]any{"id": id, "key": "COLOR", "associationType": "EXPERIMENT", "required": true},
			}, "nextPage": 1}}
		}
		return platformtest.Reply{JSON: map[string]any{"items": []any{
			map[string]any{"id": "b", "key": "COLOR", "associationType": "EXPERIMENT", "experimentKey": "ADM-1", "editableInExecution": true},
		}}}
	})

	out, err := platformtest.Stdout(t, func() error {
		return property.ListAssociations(ctx, p.Client, property.ListAssociationOptions{Key: "COLOR", Type: "experiment"})
	})

	require.NoError(t, err)
	assert.Contains(t, out, "│ "+id+" │ COLOR │ EXPERIMENT │ all   │ true     │ false            │")
	assert.Contains(t, out, "│ ADM-1 │ false    │ true             │")
	q := p.Requests("GET /api/properties/associations")[0].Query
	assert.Equal(t, []string{"COLOR"}, q["key"])
	assert.Equal(t, []string{"EXPERIMENT"}, q["associationTypeAO"])
	assert.NotContains(t, q, "size")
	assert.EqualError(t, property.ListAssociations(ctx, p.Client, property.ListAssociationOptions{Type: "team"}), "--type must be EXPERIMENT or SERVICE, not 'team'.")
}

func TestAssociationGetApplyDelete(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/properties/associations/"+id, platformtest.Reply{Body: `{"key":"COLOR","associationType":"EXPERIMENT","id":"` + id + `","version":1}`})
	p.Reply("POST /api/properties/associations", platformtest.Reply{Status: http.StatusCreated, JSON: map[string]any{"id": id, "key": "COLOR"}})
	p.Reply("DELETE /api/properties/associations/"+id, platformtest.Reply{})
	dir := t.TempDir()
	got := filepath.Join(dir, "got.yml")
	created := filepath.Join(dir, "new.yml")
	require.NoError(t, os.WriteFile(created, []byte("key: COLOR\nassociationType: SERVICE\n"), 0o644))

	out, err := platformtest.Stdout(t, func() error {
		if err := property.GetAssociation(ctx, p.Client, property.GetAssociationOptions{ID: id, File: got}); err != nil {
			return err
		}
		if err := property.ApplyAssociations(ctx, p.Client, property.ApplyAssociationOptions{Files: []string{created}}); err != nil {
			return err
		}
		return property.DeleteAssociation(ctx, p.Client, property.DeleteAssociationOptions{ID: id, Yes: true})
	})

	require.NoError(t, err)
	assert.Contains(t, out, "Property association "+id+" of COLOR created.\nProperty association "+id+" deleted.\n")
	content, _ := os.ReadFile(got)
	assert.Equal(t, "key: COLOR\nassociationType: EXPERIMENT\nid: "+id+"\n", string(content))
	content, _ = os.ReadFile(created)
	assert.Equal(t, "id: "+id+"\nkey: COLOR\nassociationType: SERVICE\n", string(content))
	assert.Equal(t, []string{"false"}, p.Requests("DELETE /api/properties/associations/" + id)[0].Query["deleteValues"])
	assert.EqualError(t, property.GetAssociation(ctx, p.Client, property.GetAssociationOptions{ID: "x"}), "Property association x not found.")
}
