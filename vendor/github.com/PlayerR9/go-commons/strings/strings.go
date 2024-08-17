package strings

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	gcers "github.com/PlayerR9/go-commons/errors"
)

// LimitReverseLines is a function that limits the lines of the data in reverse order.
//
// Parameters:
//   - data: The data to limit.
//   - limit: The limit of the lines.
//
// Returns:
//   - []byte: The limited data.
func LimitReverseLines(data string, limit int) string {
	if len(data) == 0 {
		return ""
	}

	lines := strings.Split(data, "\n")

	if limit == -1 || limit > len(lines) {
		limit = len(lines)
	}

	start_idx := len(lines) - limit

	lines = lines[start_idx:]

	return strings.Join(lines, "\n")
}

// LimitLines is a function that limits the lines of the data.
//
// Parameters:
//   - data: The data to limit.
//   - limit: The limit of the lines.
//
// Returns:
//   - string: The limited data.
func LimitLines(data string, limit int) string {
	if len(data) == 0 {
		return ""
	}

	lines := strings.Split(data, "\n")

	if limit == -1 || limit > len(lines) {
		limit = len(lines)
	}

	lines = lines[:limit]

	return strings.Join(lines, "\n")
}

// GoStringOf returns a string representation of the element.
//
// Parameters:
//   - elem: The element to get the string representation of.
//
// Returns:
//   - string: The string representation of the element.
//
// Behaviors:
//   - If the element is nil, the function returns "nil".
//   - If the element implements the fmt.GoStringer interface, the function
//     returns the result of the GoString method.
//   - If the element implements the fmt.Stringer interface, the function
//     returns the result of the String method.
//   - If the element is a string, the function returns the string enclosed in
//     double quotes.
//   - If the element is an error, the function returns the error message
//     enclosed in double quotes.
//   - Otherwise, the function returns the result of the %#v format specifier.
func GoStringOf(elem any) string {
	if elem == nil {
		return "nil"
	}

	switch elem := elem.(type) {
	case fmt.GoStringer:
		return elem.GoString()
	case fmt.Stringer:
		return strconv.Quote(elem.String())
	case string:
		return strconv.Quote(elem)
	case error:
		return strconv.Quote(elem.Error())
	default:
		return fmt.Sprintf("%#v", elem)
	}
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
func filter_equals(indices []int, data []string, other string, offset int) []int {
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
func IndicesOf(data []string, sep []string, exclude_sep bool) []int {
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
//   - *luc.ErrInvalidParameter: If the closingToken is an empty string.
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
func FindContentIndexes(op_token, cl_token string, tokens []string) (result [2]int, err error) {
	result[0] = -1
	result[1] = -1

	if cl_token == "" {
		err = gcers.NewErrInvalidParameter("cl_token", gcers.NewErrEmpty(cl_token))
		return
	}

	op_tok_idx := slices.Index(tokens, op_token)
	if op_tok_idx < 0 {
		err = NewErrTokenNotFound(op_token, true)
		return
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
		return
	}

	if balance < 0 {
		err = NewErrNeverOpened(op_token, cl_token)
		return
	} else if balance != 1 || cl_token != "\n" {
		err = NewErrTokenNotFound(cl_token, false)
		return
	}

	result[1] = len(tokens)
	return
}

// AndString is a function that returns a string representation of a slice
// of strings. Empty strings are ignored.
//
// Parameters:
//   - values: The values to convert to a string.
//
// Returns:
//   - string: The string representation of the values.
func AndString(values []string) string {
	values = TrimEmpty(values)
	if len(values) == 0 {
		return ""
	}

	if len(values) == 1 {
		return values[0]
	}

	var builder strings.Builder

	if len(values) > 2 {
		builder.WriteString(strings.Join(values[:len(values)-1], ", "))
		builder.WriteRune(',')
	} else {
		builder.WriteString(values[0])
	}

	builder.WriteString(" and ")
	builder.WriteString(values[len(values)-1])

	return builder.String()
}

// EitherOrString is a function that returns a string representation of a slice
// of strings. Empty strings are ignored.
//
// Parameters:
//   - values: The values to convert to a string.
//
// Returns:
//   - string: The string representation.
//
// Example:
//
//	EitherOrString([]string{"a", "b", "c"}, false) // "a, b or c"
func EitherOrString(values []string) string {
	values = TrimEmpty(values)

	if len(values) == 0 {
		return ""
	}

	if len(values) == 1 {
		return values[0]
	}

	var builder strings.Builder

	builder.WriteString("either ")

	if len(values) > 2 {
		builder.WriteString(strings.Join(values[:len(values)-1], ", "))
		builder.WriteRune(',')
	} else {
		builder.WriteString(values[0])
	}

	builder.WriteString(" or ")
	builder.WriteString(values[len(values)-1])

	return builder.String()
}

// OrString is a function that returns a string representation of a slice of
// strings. Empty strings are ignored.
//
// Parameters:
//   - values: The values to convert to a string.
//   - is_negative: True if the string should use "nor" instead of "or", false
//     otherwise.
//
// Returns:
//   - string: The string representation.
//
// Example:
//
//	OrString([]string{"a", "b", "c"}, true) // "a, b, nor c"
func OrString(values []string, is_negative bool) string {
	values = TrimEmpty(values)
	if len(values) == 0 {
		return ""
	}

	if len(values) == 1 {
		return values[0]
	}

	var sep string

	if is_negative {
		sep = " nor "
	} else {
		sep = " or "
	}

	var builder strings.Builder

	if len(values) > 2 {
		builder.WriteString(strings.Join(values[:len(values)-1], ", "))
		builder.WriteRune(',')
	} else {
		builder.WriteString(values[0])
	}

	builder.WriteString(sep)
	builder.WriteString(values[len(values)-1])

	return builder.String()
}

// QuoteInt returns a quoted string of an integer prefixed and suffixed with
// square brackets.
//
// Parameters:
//   - value: The integer to quote.
//
// Returns:
//   - string: The quoted integer.
func QuoteInt(value int) string {
	var builder strings.Builder

	builder.WriteRune('[')
	builder.WriteString(strconv.Itoa(value))
	builder.WriteRune(']')

	return builder.String()
}

// TrimEmpty removes empty strings from a slice of strings.
// Empty spaces at the beginning and end of the strings are also removed from
// the strings.
//
// Parameters:
//   - values: The slice of strings to trim.
//
// Returns:
//   - []string: The slice of strings with empty strings removed.
func TrimEmpty(values []string) []string {
	if len(values) == 0 {
		return values
	}

	var top int

	for i := 0; i < len(values); i++ {
		current_value := values[i]

		str := strings.TrimSpace(current_value)
		if str != "" {
			values[top] = str
			top++
		}
	}

	values = values[:top]

	return values
}

// FixTabSize fixes the tab size by replacing it with a specified string iff
// the tab size is greater than 0. The replacement string is repeated for the
// specified number of times.
//
// Parameters:
//   - size: The size of the tab.
//   - rep: The replacement string.
//
// Returns:
//   - string: The fixed tab size.
func FixTabSize(size int, rep string) string {
	if size <= 0 {
		return "\t"
	}

	return strings.Repeat(rep, size)
}

// FilterNonEmpty removes empty strings from a slice of strings.
//
// Parameters:
//   - values: The slice of strings to trim.
//
// Returns:
//   - []string: The slice of strings with empty strings removed.
func FilterNonEmpty(values []string) []string {
	if len(values) == 0 {
		return nil
	}

	var top int

	for i := 0; i < len(values); i++ {
		if values[i] != "" {
			values[top] = values[i]
			top++
		}
	}

	return values[:top:top]
}

const (
	// Ellipsis is an ellipsis string.
	Ellipsis string = "..."

	// EllipsisLen is the length of the ellipsis string.
	EllipsisLen int = len(Ellipsis)
)

// AdaptToScreenWidth is a function that adapts a slice of strings to a
// specified screen width.
//
// Parameters:
//   - elems: The slice of strings to adapt.
//   - width: The screen width to adapt to.
//   - sep: The separator string to use.
//
// Returns:
//   - string: The adapted string.
//   - int: The number of elements that were not adapted.
func AdaptToScreenWidth(elems []string, width int, sep string) (string, int) {
	if width <= 0 || len(elems) == 0 {
		return "", len(elems)
	}

	if len(elems) == 1 {
		if len(elems[0]) <= width {
			return elems[0], 0
		}

		return "", 1
	}

	var res_string string
	cut := len(elems)

	sizes := make([]int, 0, len(elems))

	for i := 0; i < len(elems); i++ {
		sizes = append(sizes, len(elems[i]))
	}

	if sep == "" {
		total := EllipsisLen

		var idx int

		for idx < len(sizes) && sizes[idx]+total+sizes[len(sizes)-1] <= width {
			total += sizes[idx]
			idx++
		}

		if idx == 0 {
			right_idx := len(elems) - 1

			for right_idx >= 0 && sizes[right_idx]+total <= width {
				total += sizes[right_idx]
				right_idx--
			}

			if right_idx == len(elems)-1 {
				return "", len(elems)
			}

			elems = append([]string{Ellipsis}, elems[:right_idx+1]...)

			res_string = strings.Join(elems, "")
		} else if idx == len(sizes) {
			res_string = strings.Join(elems, "")
		} else {
			elems = append(elems[:idx], Ellipsis, elems[len(elems)-1])
			res_string = strings.Join(elems, "")
		}
	} else {
		sep_len := len(sep)

		total := EllipsisLen + sep_len

		var idx int

		for idx < len(sizes) && sizes[idx]+total+sizes[len(sizes)-1]+sep_len <= width {
			total += sizes[idx]
			idx++
		}

		if idx == 0 {
			right_idx := len(elems) - 1

			for right_idx >= 0 && sizes[right_idx]+total+sep_len <= width {
				total += sizes[right_idx]
				right_idx--
			}

			if right_idx == len(elems)-1 {
				return "", len(elems)
			}

			elems = append([]string{Ellipsis}, elems[:right_idx+1]...)

			res_string = strings.Join(elems, sep)
		} else if idx == len(sizes) {
			res_string = strings.Join(elems, sep)
		} else {
			elems = append(elems[:idx], Ellipsis, elems[len(elems)-1])
			res_string = strings.Join(elems, sep)
		}
	}

	return res_string, cut - len(elems) + 1
}
