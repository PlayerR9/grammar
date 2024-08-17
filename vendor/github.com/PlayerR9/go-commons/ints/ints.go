package ints

import (
	"strconv"
	"strings"
)

// GetOrdinalSuffix returns the ordinal suffix for a given integer.
//
// Parameters:
//   - number: The integer for which to get the ordinal suffix. Negative
//     numbers are treated as positive.
//
// Returns:
//   - string: The ordinal suffix for the number.
//
// Example:
//   - GetOrdinalSuffix(1) returns "1st"
//   - GetOrdinalSuffix(2) returns "2nd"
func GetOrdinalSuffix(number int) string {
	var builder strings.Builder

	builder.WriteString(strconv.Itoa(number))

	if number < 0 {
		number = -number
	}

	lastTwoDigits := number % 100
	lastDigit := lastTwoDigits % 10

	if lastTwoDigits >= 11 && lastTwoDigits <= 13 {
		builder.WriteString("th")
	} else {
		switch lastDigit {
		case 1:
			builder.WriteString("st")
		case 2:
			builder.WriteString("nd")
		case 3:
			builder.WriteString("rd")
		default:
			builder.WriteString("th")
		}
	}

	return builder.String()
}
