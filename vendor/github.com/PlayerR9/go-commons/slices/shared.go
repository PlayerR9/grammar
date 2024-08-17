package slices

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

// NeverOpened is the error returned when a closing token is found without
// a corresponding opening token. Readers must return NeverOpened without wrapping it
// as callers should be able to check this error by using ==.
var NeverOpened error

func init() {
	NeverOpened = errors.New("closing token is found without a corresponding opening token")
}

// filter_equals returns the indices of the other in the data.
//
// Parameters:
//   - indices: The indices.
//   - data: The data.
//   - other: The other value.
//   - offset: The offset to start the search from.
//
// Returns:
//   - []int: The indices.
func filter_equals[T comparable](indices []int, data []T, other T, offset int) []int {
	var top int

	for i := 0; i < len(indices); i++ {
		idx := indices[i]

		if data[idx+offset] == other {
			indices[top] = idx
			top++
		}
	}

	indices = indices[:top]

	return indices
}

// Indices returns the indices of the separator in the data.
//
// Parameters:
//   - data: The data.
//   - sep: The separator.
//   - exclude_sep: Whether the separator is inclusive. If set to true, the indices will point to the character right after the
//     separator. Otherwise, the indices will point to the character right before the separator.
//
// Returns:
//   - []int: The indices.
func IndicesOf[T comparable](data []T, sep []T, exclude_sep bool) []int {
	if len(data) == 0 || len(sep) == 0 {
		return nil
	}

	var indices []int

	for i := 0; i < len(data)-len(sep); i++ {
		if data[i] == sep[0] {
			indices = append(indices, i)
		}
	}

	if len(indices) == 0 {
		return nil
	}

	for i := 1; i < len(sep); i++ {
		other := sep[i]

		indices = filter_equals(indices, data, other, i)

		if len(indices) == 0 {
			return nil
		}
	}

	if exclude_sep {
		for i := 0; i < len(indices); i++ {
			indices[i] += len(sep)
		}
	}

	return indices
}

// FindContentIndexes searches for the positions of opening and closing
// tokens in a slice of strings.
//
// Parameters:
//   - op_token: The string that marks the beginning of the content.
//   - cl_token: The string that marks the end of the content.
//   - tokens: The slice of strings in which to search for the tokens.
//
// Returns:
//   - result: An array of two integers representing the start and end indexes
//     of the content.
//   - err: Any error that occurred while searching for the tokens.
//
// Errors:
//   - *ErrTokenNotFound: If the opening or closing token is not found in the
//     content.
//   - *ErrNeverOpened: If the closing token is found without any
//     corresponding opening token.
//
// Behaviors:
//   - The first index of the content is inclusive, while the second index is
//     exclusive.
//   - This function returns a partial result when errors occur. ([-1, -1] if
//     errors occur before finding the opening token, [index, 0] if the opening
//     token is found but the closing token is not found.
func FindContentIndexes[T comparable](op_token, cl_token T, tokens []T) ([2]int, error) {
	var result [2]int

	result[0] = -1
	result[1] = -1

	op_tok_idx := slices.Index(tokens, op_token)
	if op_tok_idx < 0 {
		return result, NewErrTokenNotFound(true)
	} else {
		result[0] = op_tok_idx + 1
	}

	balance := 1
	cl_tok_idx := -1

	for i := result[0]; i < len(tokens) && cl_tok_idx == -1; i++ {
		curr_tok := tokens[i]

		if curr_tok == cl_token {
			balance--

			if balance == 0 {
				cl_tok_idx = i
			}
		} else if curr_tok == op_token {
			balance++
		}
	}

	if cl_tok_idx != -1 {
		result[1] = cl_tok_idx + 1
		return result, nil
	}

	if balance < 0 {
		return result, NeverOpened
	} else if balance != 1 {
		return result, NewErrTokenNotFound(false)
	}

	result[1] = len(tokens)
	return result, nil
}

// AndString is a function that returns a string representation of a slice
// of strings.
//
// Parameters:
//   - values: The values to convert to a string.
//   - quote: Whether to quote the values.
//
// Returns:
//   - string: The string representation of the values.
func AndString[T fmt.Stringer](values []T, quote bool) string {
	if len(values) == 0 {
		return ""
	}

	var elems []string

	if quote {
		for i := 0; i < len(values); i++ {
			str := values[i].String()
			str = strings.TrimSpace(str)
			if str == "" {
				continue
			}

			elems = append(elems, strconv.Quote(str))
		}
	} else {
		for i := 0; i < len(values); i++ {
			str := values[i].String()
			str = strings.TrimSpace(str)

			if str == "" {
				continue
			}

			elems = append(elems, str)
		}
	}

	if len(elems) == 0 {
		return ""
	} else if len(elems) == 1 {
		return elems[0]
	}

	var builder strings.Builder

	if len(elems) > 2 {
		builder.WriteString(strings.Join(elems[:len(elems)-1], ", "))
		builder.WriteRune(',')
	} else {
		builder.WriteString(elems[0])
	}

	builder.WriteString(" and ")
	builder.WriteString(elems[len(elems)-1])

	return builder.String()
}

// EitherOrString is a function that returns a string representation of a slice
// of elements.
//
// Parameters:
//   - values: The elements to convert to a string.
//   - quote: True if the elements should be quoted, false otherwise.
//
// Returns:
//   - string: The string representation.
func EitherOrString[T fmt.Stringer](values []T, quote bool) string {
	if len(values) == 0 {
		return ""
	}

	var elems []string

	if quote {
		for _, v := range values {
			str := v.String()
			str = strings.TrimSpace(str)
			if str == "" {
				continue
			}

			elems = append(elems, strconv.Quote(str))
		}
	} else {
		for _, v := range values {
			str := v.String()
			str = strings.TrimSpace(str)
			if str == "" {
				continue
			}

			elems = append(elems, str)
		}
	}

	if len(elems) == 0 {
		return ""
	} else if len(elems) == 1 {
		return elems[0]
	}

	var builder strings.Builder

	builder.WriteString("either ")

	if len(values) > 2 {
		builder.WriteString(strings.Join(elems[:len(elems)-1], ", "))
		builder.WriteRune(',')
	} else {
		builder.WriteString(elems[0])
	}

	builder.WriteString(" or ")
	builder.WriteString(elems[len(elems)-1])

	return builder.String()
}

// OrString is a function that returns a string representation of a slice of
// elements.
//
// Parameters:
//   - values: The elements to convert to a string.
//   - quote: True if the elements should be quoted, false otherwise.
//   - is_negative: True if the string should use "nor" instead of "or", false
//     otherwise.
//
// Returns:
//   - string: The string representation.
func OrString[T fmt.Stringer](values []T, quote, is_negative bool) string {
	if len(values) == 0 {
		return ""
	}

	var elems []string

	if quote {
		for _, v := range values {
			str := v.String()
			str = strings.TrimSpace(str)

			if str == "" {
				continue
			}

			elems = append(elems, strconv.Quote(str))
		}
	} else {
		for _, v := range values {
			str := v.String()
			str = strings.TrimSpace(str)
			if str == "" {
				continue
			}

			elems = append(elems, str)
		}
	}

	if len(elems) == 0 {
		return ""
	} else if len(elems) == 1 {
		return elems[0]
	}

	var sep string

	if is_negative {
		sep = " nor "
	} else {
		sep = " or "
	}

	var builder strings.Builder

	if len(values) > 2 {
		builder.WriteString(strings.Join(elems[:len(elems)-1], ", "))
		builder.WriteRune(',')
	} else {
		builder.WriteString(elems[0])
	}

	builder.WriteString(sep)
	builder.WriteString(elems[len(elems)-1])

	return builder.String()
}
