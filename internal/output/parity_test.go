// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package output

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Each case is a value and what the TypeScript CLI rendered for it, recorded from js-yaml
// and JSON.stringify by internal/jsyaml/testdata/generate.mjs.
func TestRendersExactlyLikeTheTypeScriptCLI(t *testing.T) {
	content, err := os.ReadFile("../jsyaml/testdata/cases.json")
	require.NoError(t, err)
	var cases []struct {
		Value   json.RawMessage `json:"value"`
		YAML    string          `json:"yaml"`
		JSON    string          `json:"json"`
		Compact string          `json:"compact"`
	}
	require.NoError(t, json.Unmarshal(content, &cases))

	for _, c := range cases {
		t.Run(string(c.Value), func(t *testing.T) {
			doc, err := ParseDocument(c.Value)
			require.NoError(t, err)
			yamlOut, _ := doc.Render(YAML)
			assert.Equal(t, c.YAML, string(yamlOut), "yaml")
			jsonOut, _ := doc.Render(JSON)
			assert.Equal(t, c.JSON+"\n", string(jsonOut), "json")
			assert.Equal(t, c.Compact, string(doc.RenderFile(JSON)), "compact json")
		})
	}
}

// Set JSYAML_CORPUS to a directory recorded by a corpus script: <key>.raw.json with the
// platform's response next to <key>.ts.yaml, <key>.ts.json and <key>.ts.compact.json.
func TestRendersARecordedCorpusExactly(t *testing.T) {
	dir := os.Getenv("JSYAML_CORPUS")
	if dir == "" {
		t.Skip("JSYAML_CORPUS not set")
	}
	raws, err := filepath.Glob(filepath.Join(dir, "*.raw.json"))
	require.NoError(t, err)
	require.NotEmpty(t, raws)
	failures := 0
	for _, raw := range raws {
		base := strings.TrimSuffix(raw, ".raw.json")
		content, err := os.ReadFile(raw)
		require.NoError(t, err)
		doc, err := ParseDocument(content)
		require.NoError(t, err)
		doc.Delete("version")
		for suffix, render := range map[string]func() string{
			".ts.yaml":         func() string { out, _ := doc.Render(YAML); return string(out) },
			".ts.json":         func() string { out, _ := doc.Render(JSON); return strings.TrimSuffix(string(out), "\n") },
			".ts.compact.json": func() string { return string(doc.RenderFile(JSON)) },
		} {
			expected, err := os.ReadFile(base + suffix)
			require.NoError(t, err)
			if actual := render(); actual != string(expected) {
				failures++
				if failures <= 5 {
					assert.Equal(t, string(expected), actual, filepath.Base(base)+suffix)
				}
			}
		}
	}
	t.Logf("%d documents, %d renderings differ", len(raws), failures)
	assert.Zero(t, failures)
}
