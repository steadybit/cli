// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package table

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Recorded from console-table-printer with colours off.
func TestRendersLikeConsoleTablePrinter(t *testing.T) {
	tbl := New(Column{Name: "a", Title: "A", Alignment: Left}, Column{Name: "n", Title: "Num"})
	tbl.AddRow(Default, Cell("a", "x"), Cell("n", 5))
	tbl.AddRow(Default, Cell("a", "longer ünï 😀"), Cell("n", 12345))
	tbl.AddRow(Default, Cell("a", ""), Cell("n", nil))

	assert.Equal(t, "┌───────────────┬───────┐\n│ A             │   Num │\n├───────────────┼───────┤\n│ x             │     5 │\n│ longer ünï 😀 │ 12345 │\n│               │       │\n└───────────────┴───────┘", tbl.Render())
}

func TestTakesColumnsFromTheRowsWhenNoneAreDeclared(t *testing.T) {
	tbl := New()
	tbl.AddRow(Red, Cell("target", "a"), Cell("advice", "b"))

	assert.Equal(t, "┌────────┬────────┐\n│ target │ advice │\n├────────┼────────┤\n│      a │      b │\n└────────┴────────┘", tbl.Render())
}

func TestColouredCellsKeepTheColumnsAligned(t *testing.T) {
	tbl := New(Column{Name: "s", Title: "State", Alignment: Left})
	tbl.AddRow(Default, Cell("s", "\x1b[32mcompleted\x1b[0m"))

	assert.Equal(t, "┌───────────┐\n│ State     │\n├───────────┤\n│ \x1b[32mcompleted\x1b[0m │\n└───────────┘", tbl.Render())
}
