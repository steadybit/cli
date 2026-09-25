// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandsSpaceSeparatedValuesAsCommanderDid(t *testing.T) {
	root := newRoot()

	assert.Equal(t,
		[]string{"experiment", "dump", "-d", "x", "--team", "A", "--team", "B", "--team", "C", "-t", "json"},
		expandVariadic(root, []string{"experiment", "dump", "-d", "x", "--team", "A", "B", "C", "-t", "json"}))
	assert.Equal(t,
		[]string{"experiment", "apply", "-f", "a.yml", "-f", "b.yml", "-R"},
		expandVariadic(root, []string{"experiment", "apply", "-f", "a.yml", "b.yml", "-R"}))
	// `get -f` takes a single file, so a following word is left alone.
	assert.Equal(t,
		[]string{"experiment", "get", "-k", "ADM-1", "-f", "x.yml"},
		expandVariadic(root, []string{"experiment", "get", "-k", "ADM-1", "-f", "x.yml"}))
}
