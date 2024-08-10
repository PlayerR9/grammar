package errors

import (
	"bytes"
	"strings"
	"unicode"
	"unicode/utf8"

	gcby "github.com/PlayerR9/go-commons/bytes"
	gcint "github.com/PlayerR9/go-commons/ints"
	"github.com/PlayerR9/grammar/lexing"
)

type PrintOption func(s *PrintSettings)

func WithLimitPrevLines(prev_lines int) PrintOption {
	if prev_lines < 0 {
		prev_lines = -1
	}

	return func(s *PrintSettings) {
		s.prev_lines = prev_lines
	}
}

func WithLimitNextLines(next_lines int) PrintOption {
	if next_lines < 0 {
		next_lines = -1
	}

	return func(s *PrintSettings) {
		s.next_lines = next_lines
	}
}

func WithDelta(delta int) PrintOption {
	if delta < 0 {
		delta = -1
	} else if delta == 0 {
		delta = 1
	}

	return func(s *PrintSettings) {
		s.delta = delta
	}
}

func WithFixedTabSize(tab_size int) PrintOption {
	if tab_size < 0 {
		tab_size = -1
	} else if tab_size == 0 {
		tab_size = 3
	}

	return func(s *PrintSettings) {
		s.tab_size = tab_size
	}
}

type PrintSettings struct {
	prev_lines int
	next_lines int
	delta      int
	tab_size   int
}

// make_arrow is a helper function that creates an arrow pointing to the faulty token.
//
// Parameters:
//   - faulty_line: The faulty line.
//   - start_pos: The start position of the arrow.
//
// Returns:
//   - []byte: The arrow data.
func (s *PrintSettings) make_arrow(faulty_line []byte, start_pos int) []byte {
	// luc.AssertParam("faulty_point", faulty_point >= 0 && faulty_point < len(faulty_line), luc.NewErrOutOfBounds(faulty_point, 0, len(faulty_line)))

	arrow_data := make([]byte, 0, start_pos)

	var first_tab []byte

	if s.tab_size > 0 {
		first_tab = bytes.Repeat([]byte{' '}, s.tab_size)
	} else {
		first_tab = []byte{'\t'}
	}

	for i := 0; i < start_pos; i++ {
		if faulty_line[i] == '\t' {
			arrow_data = append(arrow_data, first_tab...)
		} else {
			arrow_data = append(arrow_data, ' ')
		}
	}

	if s.delta < 0 {
		faulty_line = faulty_line[start_pos:]

		for len(faulty_line) > 0 {
			r, size := utf8.DecodeRune(faulty_line)
			faulty_line = faulty_line[size:]

			if r == utf8.RuneError {
				break
			}

			if unicode.IsSpace(r) {
				break
			}

			arrow_data = append(arrow_data, '^')
		}
	} else {
		var second_tab []byte

		if s.tab_size > 0 {
			second_tab = bytes.Repeat([]byte{'~'}, s.tab_size)
		} else {
			second_tab = []byte{'\t'}
		}

		for i := start_pos; i < start_pos+s.delta; i++ {
			if faulty_line[i] != '\t' {
				arrow_data = append(arrow_data, '^')
			} else {
				arrow_data = append(arrow_data, second_tab...)
			}
		}
	}

	return arrow_data
}

// PrintSyntaxError is a helper function that prints the syntax error.
//
// Parameters:
//   - data: The data of the faulty line.
//   - start_pos: The start position of the faulty token.
//   - opts: The print options.
//
// Returns:
//   - []byte: The syntax error data.
func PrintSyntaxError(data []byte, start_pos int, opts ...PrintOption) []byte {
	if len(data) == 0 {
		return nil
	}

	s := PrintSettings{
		prev_lines: -1,
		next_lines: -1,
		delta:      -1,
		tab_size:   -1,
	}

	for _, opt := range opts {
		opt(&s)
	}

	if start_pos < 0 {
		start_pos = len(data) + start_pos
	} else if start_pos >= len(data) {
		start_pos = len(data) - 1
	}

	if s.delta != -1 && start_pos+s.delta >= len(data) {
		s.delta = len(data) - start_pos
	}

	var before, faulty_line, after []byte

	before_idx := gcby.ReverseSearch(data, start_pos, Newline)
	after_idx := gcby.ForwardSearch(data, start_pos, Newline)

	if before_idx == -1 {
		if after_idx == -1 {
			faulty_line = data
		} else {
			faulty_line = data[:after_idx]
			after = data[after_idx+1:]
		}
	} else {
		if after_idx == -1 {
			before = data[:before_idx]
			faulty_line = data[before_idx+1:]
		} else {
			before = data[:before_idx]
			faulty_line = data[before_idx+1 : after_idx]
			after = data[after_idx+1:]
		}
	}

	arrow_data := s.make_arrow(faulty_line, start_pos-len(before))

	full_data := make([]byte, 0, len(data)+1+len(arrow_data))

	if len(before) > 0 {
		before = LimitReverseLines(before, s.prev_lines)
		full_data = append(full_data, before...)
		full_data = append(full_data, '\n')
	}

	full_data = append(full_data, faulty_line...)
	full_data = append(full_data, '\n')
	full_data = append(full_data, arrow_data...)

	if len(after) > 0 {
		full_data = append(full_data, '\n')
		after = LimitLines(after, s.next_lines)
		full_data = append(full_data, after...)
	}

	return full_data
}

// DetermineCoords is a helper function that determines the coordinates of the given position.
//
// Parameters:
//   - data: The data read from the input stream.
//   - pos: The position of the faulty token.
//
// Returns:
//   - int: The x coordinate of the faulty token.
//   - int: The y coordinate of the faulty token.
func DetermineCoords(data []byte, pos int) (int, int) {
	if len(data) == 0 {
		return 0, 0
	}

	if pos < 0 {
		pos = len(data) + pos
	} else if pos >= len(data) {
		pos = len(data) - 1
	}

	var x int
	var y int

	for i := 0; i < pos; i++ {
		if data[i] == '\n' {
			x = 0
			y++
		} else {
			x++
		}
	}

	return x, y
}

func DisplayError(data []byte, err error, opts ...PrintOption) string {
	if err == nil {
		return ""
	}

	var builder strings.Builder

	switch err := err.(type) {
	case *lexing.ErrLexing:
		x, y := DetermineCoords(data, err.StartPos)

		builder.WriteString("Lexing error at the ")
		builder.WriteString(gcint.GetOrdinalSuffix(x))
		builder.WriteString(" character of the ")
		builder.WriteString(gcint.GetOrdinalSuffix(y))
		builder.WriteString(" line:")
		builder.WriteRune('\n')
		builder.WriteRune('\t')
		builder.WriteString(err.Error())
		builder.WriteRune('\n')
		builder.WriteRune('\n')

		opts = append(opts, WithDelta(err.Delta))

		builder.Write(PrintSyntaxError(data, err.StartPos, opts...))
		builder.WriteRune('\n')

		suggestion := err.Suggestion
		if suggestion != "" {
			builder.WriteRune('\n')
			builder.WriteString("Hint: ")
			builder.WriteString(suggestion)
		}
	default:
		builder.WriteString("Error: ")
		builder.WriteString(err.Error())
	}

	return builder.String()
}
