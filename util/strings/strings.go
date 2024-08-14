package strings

import (
	"strconv"
	"strings"
)

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

// QuoteStrings is a function that quotes a slice of strings in-place.
//
// Parameters:
//   - values: The values to quote.
func QuoteStrings(values []string) {
	if len(values) == 0 {
		return
	}

	for i := 0; i < len(values); i++ {
		values[i] = strconv.Quote(values[i])
	}
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

// SliceOfRunes is a function that returns a slice of strings
// from a slice of runes.
//
// Parameters:
//   - values: The values to convert to a slice of strings.
//
// Returns:
//   - []string: The slice of strings.
func SliceOfRunes(values []rune) []string {
	if len(values) == 0 {
		return nil
	}

	elems := make([]string, 0, len(values))

	for _, value := range values {
		elems = append(elems, string(value))
	}

	return elems
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
