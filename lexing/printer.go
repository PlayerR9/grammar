package lexing

import (
	"unicode"

	gcby "github.com/PlayerR9/go-commons/bytes"
)

// make_arrow is a helper function that creates an arrow pointing to the faulty token.
//
// Parameters:
//   - faulty_line: The faulty line.
//   - faulty_point: The faulty point.
//
// Returns:
//   - []byte: The arrow data.
func make_arrow(faulty_line []byte, faulty_point int) []byte {
	// luc.AssertParam("faulty_point", faulty_point >= 0 && faulty_point < len(faulty_line), luc.NewErrOutOfBounds(faulty_point, 0, len(faulty_line)))

	arrow_data := make([]byte, 0, faulty_point)

	for i := 0; i < faulty_point; i++ {
		if faulty_line[i] == '\t' {
			arrow_data = append(arrow_data, '\t')
		} else {
			arrow_data = append(arrow_data, ' ')
		}
	}

	for i := faulty_point; i < len(faulty_line); i++ {
		if unicode.IsSpace(rune(faulty_line[i])) {
			break
		}

		arrow_data = append(arrow_data, '^')
	}

	arrow_data = append(arrow_data, '\n')

	return arrow_data
}

// PrintSyntaxError is a helper function that prints the syntax error.
//
// Parameters:
//   - data: The data of the faulty line.
//   - start_pos: The start position of the faulty token.
//   - delta: The end position of the faulty token. Calculated as start_pos + delta.
//
// Returns:
//   - []byte: The syntax error data.
func PrintSyntaxError(data []byte, start_pos, delta int) []byte {
	if len(data) == 0 {
		return nil
	}

	if start_pos < 0 {
		start_pos = len(data) + start_pos
	}

	end_pos := start_pos + delta

	if end_pos >= len(data) {
		end_pos = len(data)
	}

	var before, faulty_line, after []byte

	before_idx := gcby.ReverseSearch(data, start_pos, []byte("\n"))
	after_idx := gcby.ForwardSearch(data, start_pos, []byte("\n"))

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

	fault_point := start_pos + end_pos - len(before) - 1

	arrow_data := make_arrow(faulty_line, fault_point)

	var full_data []byte

	full_data = append(full_data, before...)
	full_data = append(full_data, '\n')
	full_data = append(full_data, faulty_line...)
	full_data = append(full_data, '\n')
	full_data = append(full_data, arrow_data...)
	full_data = append(full_data, after...)

	return full_data
}
