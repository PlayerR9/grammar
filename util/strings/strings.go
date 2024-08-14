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
