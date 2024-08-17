package runes

import (
	"bytes"
	"strings"

	gcint "github.com/PlayerR9/go-commons/ints"
	gcch "github.com/PlayerR9/go-commons/runes"
)

// RuneTable is a table of runes.
type RuneTable struct {
	// table is the table of runes.
	table [][]rune
}

// String implements the fmt.Stringer interface.
func (rt RuneTable) String() string {
	lines := make([]string, 0, len(rt.table))

	for _, row := range rt.table {
		lines = append(lines, string(row))
	}

	return strings.Join(lines, "\n")
}

// NewRuneTable creates a new RuneTable.
//
// Returns:
//   - RuneTable: The new RuneTable.
func NewRuneTable() RuneTable {
	return RuneTable{
		table: make([][]rune, 0),
	}
}

// FromBytes initializes the RuneTable from a slice of slice of bytes.
//
// Parameters:
//   - lines: The slice of slice of bytes.
//
// Returns:
//   - error: An error of type *ints.ErrAt if the lines could not be processed.
func (rt *RuneTable) FromBytes(lines [][]byte) error {
	table := make([][]rune, 0, len(lines))

	for i, line := range lines {
		row, err := gcch.BytesToUtf8(line)
		if err != nil {
			return gcint.NewErrAt(i+1, "line", err)
		}

		table = append(table, row)
	}

	rt.table = table

	return nil
}

// FromRunes initializes the RuneTable from a slice of slice of runes.
//
// Parameters:
//   - lines: The slice of slice of runes.
func (rt *RuneTable) FromRunes(lines [][]rune) {
	rt.table = lines
}

// FromStrings initializes the RuneTable from a slice of strings.
//
// Parameters:
//   - lines: The slice of strings.
//
// Returns:
//   - error: An error of type *ints.ErrAt if the lines could not be processed.
func (rt *RuneTable) FromStrings(lines []string) error {
	table := make([][]rune, 0, len(lines))

	for i, line := range lines {
		row, err := gcch.StringToUtf8(line)
		if err != nil {
			return gcint.NewErrAt(i+1, "line", err)
		}

		table = append(table, row)
	}

	rt.table = table

	return nil
}

// RightMostEdge gets the right most edge of the content.
//
// Parameters:
//   - content: The content.
//
// Returns:
//   - int: The right most edge.
func (rt RuneTable) RightMostEdge() int {
	var longest_line int

	for _, row := range rt.table {
		if len(row) > longest_line {
			longest_line = len(row)
		}
	}

	return longest_line
}

// AlignRightEdge aligns the right edge of the table.
//
// Returns:
//   - int: The right most edge.
func (rt *RuneTable) AlignRightEdge() int {
	edge := rt.RightMostEdge()

	for i := 0; i < len(rt.table); i++ {
		curr_row := rt.table[i]

		padding := edge - len(curr_row)

		padding_right := make([]rune, 0, padding)
		for i := 0; i < padding; i++ {
			padding_right = append(padding_right, ' ')
		}

		rt.table[i] = append(curr_row, padding_right...)
	}

	return edge
}

// PrependTopRow prepends a row to the top of the table.
//
// Parameters:
//   - row: The row to prepend.
func (rt *RuneTable) PrependTopRow(row []rune) {
	rt.table = append([][]rune{row}, rt.table...)
}

// AppendBottomRow appends a row to the bottom of the table.
//
// Parameters:
//   - row: The row to append.
func (rt *RuneTable) AppendBottomRow(row []rune) {
	rt.table = append(rt.table, row)
}

// PrefixEachRow prefixes each row with the given prefix.
//
// Parameters:
//   - prefix: The prefix to add to each row.
func (rt *RuneTable) PrefixEachRow(prefix []rune) {
	for i := 0; i < len(rt.table); i++ {
		new_row := append(prefix, rt.table[i]...)
		rt.table[i] = new_row
	}
}

// SuffixEachRow suffixes each row with the given suffix.
//
// Parameters:
//   - suffix: The suffix to add to each row.
func (rt *RuneTable) SuffixEachRow(suffix []rune) {
	for i := 0; i < len(rt.table); i++ {
		new_row := append(rt.table[i], suffix...)
		rt.table[i] = new_row
	}
}

// Byte returns the byte representation of the table.
//
// Returns:
//   - []byte: The byte representation of the table.
func (rt RuneTable) Byte() []byte {
	if len(rt.table) == 0 {
		return []byte{}
	}

	var buffer bytes.Buffer

	buffer.Grow(gcch.JoinSize(rt.table))

	for _, r := range rt.table[0] {
		buffer.WriteRune(r)
	}

	for i := 1; i < len(rt.table); i++ {
		buffer.WriteRune('\n')

		for _, r := range rt.table[i] {
			buffer.WriteRune(r)
		}
	}

	return buffer.Bytes()
}

// Rune returns the rune representation of the table.
//
// Returns:
//   - []rune: The rune representation of the table.
func (rt RuneTable) Rune() []rune {
	return gcch.Join(rt.table, '\n')
}
