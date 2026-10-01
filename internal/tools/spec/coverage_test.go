// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package main

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const root = "../../.."

// Every operation of the committed spec is called by a command, deprecated, or named in
// notCovered with a reason. A spec refresh that brings a new endpoint fails here until
// someone decides what the CLI does with it.
func TestEveryOperationIsCoveredOrListed(t *testing.T) {
	raw, err := os.ReadFile(root + "/" + specFile)
	require.NoError(t, err)

	report, err := checkCoverage(raw, root)
	require.NoError(t, err)

	for _, op := range report.Uncovered {
		t.Errorf("%s %s (%s) is neither called by a command nor listed in notCovered", op.Method, op.Path, op.ID)
	}
	for _, name := range report.Stale {
		t.Errorf("notCovered lists %s, which is called by a command or no longer in the spec", name)
	}
}

// The check matches operations to the client's methods by name, so the names it derives
// have to be the ones oapi-codegen generated.
func TestOperationIDsMapToGeneratedMethods(t *testing.T) {
	raw, err := os.ReadFile(root + "/" + specFile)
	require.NoError(t, err)
	generated, err := os.ReadFile(root + "/api/platform.gen.go")
	require.NoError(t, err)
	methods := map[string]bool{}
	for _, m := range regexp.MustCompile(`(?m)^func \(c \*Client\) (\w+)\(ctx`).FindAllSubmatch(generated, -1) {
		methods[string(m[1])] = true
	}

	ops, err := operations(raw)
	require.NoError(t, err)
	require.NotEmpty(t, ops)
	for _, op := range ops {
		assert.True(t, methods[methodName(op.ID)], "%s: no client method %s", op.ID, methodName(op.ID))
	}
}

func TestMethodNames(t *testing.T) {
	assert.Equal(t, "GetAccessTokens1", methodName("getAccessTokens_1"))
	assert.Equal(t, "GetTeamMembers", methodName("getTeamMembers"))
	assert.Equal(t, "Health", methodName("health"))
}

// A new endpoint shows as uncovered; once a command calls it, its notCovered line is stale.
func TestNewAndStaleOperations(t *testing.T) {
	spec := []byte(`{"paths": {
		"/api/health": {"get": {"operationId": "health"}},
		"/api/new": {"get": {"operationId": "getNew", "summary": "Something new"}, "parameters": []},
		"/api/old": {"get": {"operationId": "getOld", "deprecated": true}}
	}}`)
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(dir+"/internal/thing", 0o755))
	require.NoError(t, os.MkdirAll(dir+"/cmd", 0o755))
	require.NoError(t, os.WriteFile(dir+"/internal/thing/thing.go", []byte("package thing\n\nfunc f(c interface{ Health() }) { c.Health() }\n"), 0o644))

	report, err := checkCoverage(spec, dir)
	require.NoError(t, err)

	require.Len(t, report.Uncovered, 1)
	assert.Equal(t, "getNew", report.Uncovered[0].ID)
	assert.Equal(t, 1, report.Called)
	assert.Contains(t, report.Stale, "Health", "listed, but called")
	assert.Contains(t, report.Stale, "ConnectionCheck", "listed, but not in this spec")
}
