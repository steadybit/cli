// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package license

import (
	"testing"

	"github.com/steadybit/cli/v6/internal/table"
	"github.com/stretchr/testify/assert"
)

func TestHighlightsAHardLimitOnceReachedAndASoftOneOnceExceeded(t *testing.T) {
	n := func(v int64) *int64 { return &v }
	for _, c := range []struct {
		name  string
		f     feature
		limit string
		color table.Color
	}{
		{"below a hard limit", feature{Type: "HARD_LIMIT", Usage: n(9), HardLimit: n(10)}, "10", table.Default},
		{"at a hard limit", feature{Type: "HARD_LIMIT", Usage: n(10), HardLimit: n(10)}, "10", table.Red},
		{"over a hard limit", feature{Type: "HARD_LIMIT", Usage: n(11), HardLimit: n(10)}, "10", table.Red},
		{"at a soft limit", feature{Type: "SOFT_LIMIT", Usage: n(10), SoftLimit: n(10)}, "10 (soft)", table.Default},
		{"over a soft limit", feature{Type: "SOFT_LIMIT", Usage: n(11), SoftLimit: n(10)}, "10 (soft)", table.Red},
		// The type decides, not which field is sent.
		{"a soft limit also sending a hard one", feature{Type: "SOFT_LIMIT", Usage: n(5), SoftLimit: n(3), HardLimit: n(10)}, "3 (soft)", table.Red},
		{"a hard limit also sending a soft one", feature{Type: "HARD_LIMIT", Usage: n(5), SoftLimit: n(3), HardLimit: n(10)}, "10", table.Default},
		{"a hard limit without an amount", feature{Type: "HARD_LIMIT", Usage: n(5), SoftLimit: n(3)}, "unlimited", table.Default},
		{"no usage", feature{Type: "HARD_LIMIT", HardLimit: n(0)}, "0", table.Default},
	} {
		t.Run(c.name, func(t *testing.T) {
			limit, color := limitOf(c.f)
			assert.Equal(t, c.limit, limit)
			assert.Equal(t, c.color, color)
		})
	}
}
