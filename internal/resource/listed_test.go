// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package resource_test

import (
	"context"
	"testing"

	"github.com/steadybit/cli/internal/platformtest"
	"github.com/steadybit/cli/internal/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeListedKeepsTheRawItemsBesideTheTypedOnes(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/teams", platformtest.Reply{Body: `{"teams":[{"key":"ADM","name":"Admins","extra":{"kept":true}}]}`})

	var typed struct {
		Teams []struct{ Key string } `json:"teams"`
	}
	resp, err := p.Client.GetTeams(context.Background(), nil)
	items, err := resource.DecodeListed(resp, err, "teams", &typed)

	require.NoError(t, err)
	assert.Equal(t, "ADM", typed.Teams[0].Key)
	assert.JSONEq(t, `{"key":"ADM","name":"Admins","extra":{"kept":true}}`, string(items[0]))
}
