// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package cli

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Every command of the TypeScript CLI, which pipelines call by name. Removing one is a
// breaking change; this list is what the Go CLI promised to keep.
var typeScriptCommands = []string{
	"advice validate-status",
	"config show", "config profile add", "config profile list", "config profile ls", "config profile remove", "config profile select",
	"experiment run", "experiment exec", "experiment get", "experiment apply", "experiment delete", "experiment dump",
	"template list", "template get",
	"execution get", "execution cancel", "execution property set", "execution property add", "execution artifact list", "execution artifact download",
	"schedule list", "schedule get", "schedule apply", "schedule create", "schedule update", "schedule enable", "schedule disable", "schedule delete",
	"service list", "service get", "service apply", "service delete", "service risk",
	"service experiment list", "service experiment provide", "service experiment link", "service experiment unlink",
	"service variable get", "service variable set",
	"service-profile list", "service-profile get", "service-profile apply", "service-profile delete",
}

func TestKeepsEveryCommandOfTheTypeScriptCLI(t *testing.T) {
	root := newRoot()
	for _, path := range typeScriptCommands {
		cmd, rest, err := root.Find(strings.Fields(path))
		if assert.NoError(t, err, path) {
			assert.Empty(t, rest, path)
			assert.Contains(t, append(cmd.Aliases, cmd.Name()), strings.Fields(path)[len(strings.Fields(path))-1], path)
		}
	}
}
