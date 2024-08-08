package lexer

import (
	"unicode"

	gcby "github.com/PlayerR9/go-commons/bytes"
	gr "github.com/PlayerR9/grammar/grammar"
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
//   - tokens: The tokens of the faulty line.
//
// Returns:
//   - []byte: The syntax error data.
func PrintSyntaxError[T gr.TokenTyper](data []byte, tokens []*gr.Token[T]) []byte {
	if len(tokens) < 2 || len(data) == 0 {
		return nil
	}

	last_token := tokens[len(tokens)-2]

	idx := last_token.At

	var before, faulty_line, after []byte

	before_idx := gcby.ReverseSearch(data, idx, []byte("\n"))
	after_idx := gcby.ForwardSearch(data, idx, []byte("\n"))

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

	fault_point := idx + last_token.Size() - len(before) - 1

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
