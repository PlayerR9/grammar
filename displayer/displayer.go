package displayer

import (
	"bytes"
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"

	gfch "github.com/PlayerR9/go-commons/Formatting/runes"
	gcby "github.com/PlayerR9/go-commons/bytes"
	gcint "github.com/PlayerR9/go-commons/ints"
	dbg "github.com/PlayerR9/go-debug/assert"
	"github.com/PlayerR9/grammar/lexing"
)

var (
	BoxStyle *gfch.BoxStyle
)

func init() {
	BoxStyle = gfch.NewBoxStyle(gfch.BtNormal, true, [4]int{1, 2, 1, 2})
}

// PrintSettings is the settings of the printer.
type PrintSettings struct {
	// prev_lines is the number of previous lines.
	prev_lines int

	// next_lines is the number of next lines.
	next_lines int

	// delta is the delta.
	delta int

	// tab_size is the tab size.
	tab_size int
}

// make_arrow is a helper function that creates an arrow pointing to the faulty token.
//
// Parameters:
//   - faulty_line: The faulty line.
//   - start_pos: The start position of the arrow.
//
// Returns:
//   - []byte: The arrow data.
func (s *PrintSettings) make_arrow(faulty_line []byte, start_pos int) ([]byte, error) {
	var buffer bytes.Buffer

	buffer.Grow(len(faulty_line))

	first_tab := gcby.FixTabSize(s.tab_size, []byte{' '})

	for i := 0; i < start_pos; i++ {
		if faulty_line[i] == '\t' {
			buffer.Write(first_tab)
		} else {
			buffer.WriteByte(' ')
		}
	}

	if s.delta < 0 {
		faulty_line = faulty_line[start_pos:]

		dbg.Assert(len(faulty_line) > 0, "faulty_line is empty; this should never happen")

		r, size := utf8.DecodeRune(faulty_line)
		faulty_line = faulty_line[size:]

		if r == utf8.RuneError {
			return nil, errors.New("invalid utf8 sequence")
		}

		buffer.WriteRune('^')

		for len(faulty_line) > 0 {
			r, size := utf8.DecodeRune(faulty_line)
			faulty_line = faulty_line[size:]

			if r == utf8.RuneError {
				break
			}

			if unicode.IsSpace(r) {
				break
			}

			buffer.WriteRune('^')
		}
	} else {
		second_tab := gcby.FixTabSize(s.tab_size, []byte{'~'})

		for i := start_pos; i < start_pos+s.delta; i++ {
			if faulty_line[i] != '\t' {
				buffer.WriteByte('^')
			} else {
				buffer.Write(second_tab)
			}
		}
	}

	return buffer.Bytes(), nil
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

	before_idx := gcby.ReverseSearch(data, start_pos, gcby.Newline)
	after_idx := gcby.ForwardSearch(data, start_pos, gcby.Newline)

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
		} else if before_idx == after_idx {
			before = data[:before_idx]
			after = data[after_idx+1:]
		} else {
			before = data[:before_idx]
			faulty_line = data[before_idx+1 : after_idx]
			after = data[after_idx+1:]
		}
	}

	arrow_data, err := s.make_arrow(faulty_line, start_pos-len(before))
	dbg.AssertErr(err, "PrintSettings.make_arrow(%q, %d)", string(faulty_line), start_pos-len(before))

	before = gcby.LimitReverseLines(before, s.prev_lines)
	after = gcby.LimitLines(after, s.next_lines)

	var buffer bytes.Buffer

	buffer.Grow(len(data) + 1 + len(arrow_data))

	if len(before) > 0 {
		buffer.Write(before)
		buffer.WriteByte('\n')
	}

	buffer.Write(faulty_line)
	buffer.WriteByte('\n')
	buffer.Write(arrow_data)

	if len(after) == 0 {
		return buffer.Bytes()
	}

	buffer.WriteByte('\n')
	buffer.Write(after)

	return buffer.Bytes()
}

// PrintBoxedData is a helper function that prints the boxed data.
//
// Parameters:
//   - data: The data of the faulty line.
//   - at: The start position of the faulty token.
//   - opts: The print options.
//
// Returns:
//   - []byte: The boxed data.
func PrintBoxedData(data []byte, at int, opts ...PrintOption) []byte {
	var table gfch.RuneTable

	err := table.FromBytes(bytes.Split(PrintSyntaxError(data, at, opts...), []byte("\n")))
	dbg.AssertErr(err, "table.FromBytes(data)")

	err = BoxStyle.Apply(&table)
	dbg.AssertErr(err, "BoxStyle.Apply(&table)")

	return table.Byte()
}

// DisplayError is a helper function that displays the error.
//
// Parameters:
//   - data: The data read from the input stream.
//   - err: The error.
//   - opts: The print options.
//
// Returns:
//   - string: The error data.
func DisplayError(data []byte, err error, opts ...PrintOption) string {
	if err == nil {
		return ""
	}

	var builder strings.Builder

	switch reason := err.(type) {
	case *lexing.ErrLexing:
		x, y := gcby.DetermineCoords(data, reason.StartPos)

		builder.WriteString("Lexing error at the ")
		builder.WriteString(gcint.GetOrdinalSuffix(x + 1))
		builder.WriteString(" character of the ")
		builder.WriteString(gcint.GetOrdinalSuffix(y + 1))
		builder.WriteString(" line:")
		builder.WriteRune('\n')
		builder.WriteRune('\t')
		builder.WriteString(reason.Reason.Error())
		builder.WriteRune('\n')
		builder.WriteRune('\n')

		opts = append(opts, WithDelta(reason.Delta))

		_, _ = builder.Write(PrintBoxedData(data, reason.StartPos, opts...))
		builder.WriteRune('\n')

		suggestion := reason.Suggestion
		if suggestion != "" {
			builder.WriteRune('\n')
			builder.WriteString("Hint: ")
			builder.WriteString(suggestion)
		}
	case *ErrParsing:
		x, y := gcby.DetermineCoords(data, reason.StartPos)

		builder.WriteString("Parsing error at the ")
		builder.WriteString(gcint.GetOrdinalSuffix(x))
		builder.WriteString(" character of the ")
		builder.WriteString(gcint.GetOrdinalSuffix(y))
		builder.WriteString(" line:")
		builder.WriteRune('\n')
		builder.WriteRune('\t')
		builder.WriteString(reason.Reason.Error())
		builder.WriteRune('\n')
		builder.WriteRune('\n')

		opts = append(opts, WithDelta(reason.Delta))

		_, _ = builder.Write(PrintBoxedData(data, reason.StartPos, opts...))
		builder.WriteRune('\n')

		suggestion := reason.Suggestion
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
