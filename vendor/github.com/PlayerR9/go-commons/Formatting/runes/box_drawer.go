package runes

import (
	"errors"

	gcch "github.com/PlayerR9/go-commons/runes"
)

var (
	// DefaultBoxStyle is the default box style.
	DefaultBoxStyle *BoxStyle
)

func init() {
	DefaultBoxStyle = &BoxStyle{
		LineType: BtNormal,
		IsHeavy:  false,
		Padding:  [4]int{1, 1, 1, 1},
	}
}

// BoxBorderType is the type of the box border.
type BoxBorderType int

const (
	// BtNormal is the normal box border type.
	BtNormal BoxBorderType = iota

	// BtTriple is the triple box border type.
	BtTriple

	// BtQuadruple is the quadruple box border type.
	BtQuadruple

	// BtDouble is the double box border type.
	BtDouble

	// BtRounded is like BtNormal but with rounded corners.
	BtRounded
)

// BoxStyle is the style of the box.
type BoxStyle struct {
	// LineType is the type of the line.
	LineType BoxBorderType

	// IsHeavy is whether the line is heavy or not.
	// Only applicable to BtNormal, BtTriple, and BtQuadruple.
	IsHeavy bool

	// Padding is the padding of the box.
	// [Top, Right, Bottom, Left]
	Padding [4]int
}

// NewBoxStyle creates a new box style.
//
// Negative padding are set to 0.
//
// Parameters:
//   - line_type: The line type.
//   - is_heavy: Whether the line is heavy or not.
//   - padding: The padding of the box. [Top, Right, Bottom, Left]
//
// Returns:
//   - *BoxStyle: The new box style.
func NewBoxStyle(line_type BoxBorderType, is_heavy bool, padding [4]int) *BoxStyle {
	for i := 0; i < 4; i++ {
		if padding[i] < 0 {
			padding[i] = 0
		}
	}

	bs := &BoxStyle{
		LineType: line_type,
		IsHeavy:  is_heavy,
		Padding:  padding,
	}

	return bs
}

// Corners gets the corners of the box.
//
// Returns:
//   - [4]rune: The corners. [TopLeft, TopRight, BottomLeft, BottomRight]
func (bs *BoxStyle) Corners() [4]rune {
	var corners [4]rune

	if bs.IsHeavy {
		corners = [4]rune{'┏', '┓', '┗', '┛'}
	} else {
		corners = [4]rune{'┌', '┐', '└', '┘'}
	}

	return corners
}

// TopBorder gets the top border of the box.
//
// It also applies to the bottom border as they are the same.
//
// Returns:
//   - string: The top border.
func (bs *BoxStyle) TopBorder() rune {
	var tb_border rune

	switch bs.LineType {
	case BtNormal:
		if bs.IsHeavy {
			tb_border = '━'
		} else {
			tb_border = '─'
		}
	case BtTriple:
		if bs.IsHeavy {
			tb_border = '┅'
		} else {
			tb_border = '┄'
		}
	case BtQuadruple:
		if bs.IsHeavy {
			tb_border = '┉'
		} else {
			tb_border = '┅'
		}
	case BtDouble:
		tb_border = '═'
	case BtRounded:
		tb_border = '─'
	}

	return tb_border
}

// SideBorder gets the side border of the box.
//
// It also applies to the left border as they are the same.
//
// Returns:
//   - string: The side border.
func (bs *BoxStyle) SideBorder() rune {
	var side_border rune

	switch bs.LineType {
	case BtNormal:
		if bs.IsHeavy {
			side_border = '┃'
		} else {
			side_border = '│'
		}
	case BtTriple:
		if bs.IsHeavy {
			side_border = '┇'
		} else {
			side_border = '┆'
		}
	case BtQuadruple:
		if bs.IsHeavy {
			side_border = '┋'
		} else {
			side_border = '┆'
		}
	case BtDouble:
		side_border = '║'
	case BtRounded:
		side_border = '│'
	}

	return side_border
}

// make_side_padding is a helper function to make side padding.
//
// Parameters:
//   - width: The width of the padding.
//
// Returns:
//   - []rune: The side padding.
func make_side_padding(width int) []rune {
	// dbg.AssertParam("width", width >= 0, luc.NewErrGTE(0))

	side_padding := make([]rune, 0, width)
	for i := 0; i < width; i++ {
		side_padding = append(side_padding, ' ')
	}

	return side_padding
}

// make_tb_border is a helper function to make a top or bottom border.
//
// Parameters:
//   - width: The width of the border.
//   - border: The border character.
//   - left_corner: The left corner character.
//   - right_corner: The right corner character.
//
// Returns:
//   - []rune: The top or bottom border.
//
// Assertions:
//   - width >= 0
//   - border != 0
//   - left_corner != 0
//   - right_corner != 0
func make_tb_border(width int, border, left_corner, right_corner rune) []rune {
	// dbg.AssertParam("width", width >= 0, luc.NewErrGTE(0))
	// dbg.AssertParam("border", border != 0, errors.New("border cannot be \\0"))
	// dbg.AssertParam("left_corner", left_corner != 0, errors.New("left_corner cannot be \\0"))
	// dbg.AssertParam("right_corner", right_corner != 0, errors.New("right_corner cannot be \\0"))

	row := make([]rune, 0, width+2)

	row = append(row, left_corner)
	for i := 0; i < width; i++ {
		row = append(row, border)
	}

	row = append(row, right_corner)

	return row
}

// Apply draws a box around a content that is specified in a table.
//
// Format: If the content is [['H', 'e', 'l', 'l', 'o'], ['W', 'o', 'r', 'l', 'd']], the box will be:
//
//	┏━━━━━━━┓
//	┃ Hello ┃
//	┃ World ┃
//	┗━━━━━━━┛
//
// Parameters:
//   - table: The table that contains the content to be drawn.
//
// Returns:
//   - error: An error if the content could not be processed.
//
// Behaviors:
//   - If the box style is nil, the default box style will be used.
//
// Each string of the content represents a row in the box.
func (bs *BoxStyle) Apply(table *RuneTable) error {
	if table == nil {
		return errors.New("table cannot be nil")
	}

	for i := 0; i < 4; i++ {
		if bs.Padding[i] < 0 {
			bs.Padding[i] = 0
		}
	}

	side_border := bs.SideBorder()
	left_padding := make_side_padding(bs.Padding[3])
	right_padding := make_side_padding(bs.Padding[1])
	tbb_char := bs.TopBorder()
	corners := bs.Corners()
	prefix := append([]rune{side_border}, left_padding...)
	suffix := append(right_padding, side_border)

	right_edge := table.AlignRightEdge()

	total_width := right_edge + bs.Padding[1] + bs.Padding[3]
	empty_row := gcch.Repeat(' ', right_edge)

	top_border := make_tb_border(total_width, tbb_char, corners[0], corners[1])
	bottom_border := make_tb_border(total_width, tbb_char, corners[2], corners[3])

	for i := 0; i < bs.Padding[0]; i++ {
		table.PrependTopRow(empty_row)
	}

	for i := 0; i < bs.Padding[2]; i++ {
		table.AppendBottomRow(empty_row)
	}

	table.PrefixEachRow(prefix)
	table.SuffixEachRow(suffix)
	table.PrependTopRow(top_border)
	table.AppendBottomRow(bottom_border)

	return nil
}
