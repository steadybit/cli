// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package table prints tables laid out as console-table-printer did for the TypeScript
// CLI: box drawing, one space of padding, right-aligned unless a column says otherwise.
package table

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/steadybit/cli/internal/output"
)

type Alignment int

const (
	Right Alignment = iota
	Left
)

type Column struct {
	Name      string
	Title     string
	Alignment Alignment
}

type Color int

const (
	Default Color = iota
	Red
	Green
)

type Table struct {
	columns []Column
	rows    []row
}

type row struct {
	cells map[string]string
	color Color
}

func New(columns ...Column) *Table { return &Table{columns: columns} }

// AddRow adds a row. Without declared columns, they are taken from the row's keys in
// order, as console-table-printer did.
func (t *Table) AddRow(color Color, cells ...[2]string) {
	values := map[string]string{}
	for _, cell := range cells {
		values[cell[0]] = cell[1]
		if !t.hasColumn(cell[0]) {
			t.columns = append(t.columns, Column{Name: cell[0]})
		}
	}
	t.rows = append(t.rows, row{cells: values, color: color})
}

func (t *Table) hasColumn(name string) bool {
	for _, c := range t.columns {
		if c.Name == name {
			return true
		}
	}
	return false
}

func Cell(name string, value any) [2]string {
	if value == nil {
		return [2]string{name, ""}
	}
	return [2]string{name, fmt.Sprint(value)}
}

func (c Column) title() string {
	if c.Title != "" {
		return c.Title
	}
	return c.Name
}

func pad(s string, width int, a Alignment) string {
	gap := strings.Repeat(" ", width-runewidth.StringWidth(s))
	if a == Left {
		return s + gap
	}
	return gap + s
}

func colored(s string, c Color) string {
	switch c {
	case Red:
		return output.Red(s)
	case Green:
		return output.Green(s)
	}
	return s
}

func (t *Table) Render() string {
	widths := make([]int, len(t.columns))
	for i, c := range t.columns {
		widths[i] = runewidth.StringWidth(c.title())
		for _, r := range t.rows {
			widths[i] = max(widths[i], runewidth.StringWidth(r.cells[c.Name]))
		}
	}
	line := func(left, middle, right string) string {
		parts := make([]string, len(widths))
		for i, w := range widths {
			parts[i] = strings.Repeat("─", w+2)
		}
		return left + strings.Join(parts, middle) + right
	}
	var b strings.Builder
	b.WriteString(line("┌", "┬", "┐") + "\n│")
	for i, c := range t.columns {
		b.WriteString(" " + output.Bold(pad(c.title(), widths[i], c.Alignment)) + " │")
	}
	b.WriteString("\n" + line("├", "┼", "┤") + "\n")
	for _, r := range t.rows {
		b.WriteString("│")
		for i, c := range t.columns {
			b.WriteString(" " + colored(pad(r.cells[c.Name], widths[i], c.Alignment), r.color) + " │")
		}
		b.WriteString("\n")
	}
	b.WriteString(line("└", "┴", "┘"))
	return b.String()
}

func (t *Table) Print() { fmt.Println(t.Render()) }
