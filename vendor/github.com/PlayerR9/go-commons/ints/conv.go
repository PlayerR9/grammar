package ints

import (
	"math"
)

const (
	// Unary is the unary base number.
	Unary int = 1

	// Binary is the binary base number.
	Binary int = 2

	// Ternary is the ternary base number.
	Ternary int = 3

	// Quaternary is the quaternary base number.
	Quaternary int = 4

	// Octal is the octal base number.
	Octal int = 8

	// Decimal is the decimal base number.
	Decimal int = 10

	// Duodecimal is the duodecimal base number.
	Duodecimal int = 12

	// Hexadecimal is the hexadecimal base number.
	Hexadecimal int = 16
)

// DecToBase converts a decimal number to a base number. The converted sequence
// of digits are in Least Significant Digit (LSD) order. This means that the
// digit at index 0 is the least significant.
//
// Parameters:
//   - number: The decimal number to convert.
//   - base: The base number to convert to.
//
// Returns:
//   - []int: The converted base number.
//   - bool: True if base is greater than 0. False otherwise.
//
// If base is 1, the returned slice will have a length equal to the number itself and
// it will be filled with 0s.
func DecToBase(number, base int) ([]int, bool) {
	if base <= 0 {
		return nil, false
	}

	if number < 0 {
		number *= -1
	}

	if base == 1 {
		return make([]int, number), true
	} else if number < base {
		return []int{number}, true
	}

	log_base := math.Log(float64(base))
	size := int(math.Log(float64(number))/log_base + 1)

	digits := make([]int, 0, size)

	for number > 0 {
		digits = append(digits, number%base)
		number /= base
	}

	return digits, true
}

// BaseToDec converts a base number to a decimal number.
//
// Parameters:
//   - digits: The base number to convert. Must be in Least Significant Digit (LSD) order.
//   - base: The base number to convert from.
//
// Returns:
//   - int: The converted decimal number.
//   - bool: True if base is greater than 0. False otherwise.
//
// Because this does not perform any digit out of bounds checks, the caller must guarantee that
// all digits are in the range [0, base-1] as, otherwise, the result will be incorrect.
//
// If the base is 1, the result will always be the length of the digits.
func BaseToDec(digits []int, base int) (int, bool) {
	if base <= 0 {
		return 0, false
	} else if base == 1 {
		return len(digits), true
	} else if len(digits) == 0 {
		return 0, true
	}

	var result float64

	for i, digit := range digits {
		result += math.Pow(float64(base), float64(i)) * float64(digit)
	}

	return int(result), true
}

// CheckDigits checks that all digits are in the range [0, base-1] and whether
// the base is greater than 0.
//
// Parameters:
//   - digits: The digits to check.
//   - base: The base number to check.
//
// Returns:
//   - error: An error if either the base or the digits are invalid.
//
// Errors:
//   - *ErrInvalidBase: If the base is not greater than 0.
//   - *ErrAt: If any digit is not in the range [0, base-1].
func CheckDigits(digits []int, base int) error {
	if base <= 0 {
		return NewErrInvalidBase("base")
	} else if len(digits) == 0 || base == 1 {
		return nil
	}

	for i, digit := range digits {
		if digit < 0 || digit >= base {
			return NewErrInvalidDigit(i, digit, base)
		}
	}

	return nil
}

// BaseToBase converts a base number to another base number.
//
// Parameters:
//   - digits: The base number to convert. Must be in Least Significant Digit (LSD) order.
//   - base_from: The base number to convert from.
//   - base_to: The base number to convert to.
//
// Returns:
//   - []int: The converted base number.
//   - error: An error of type *errors.ErrInvalidParameter if either base_from or base_to
//     is not greater than 0.
//
// If both base_from and base_to are the same, the slice will be the same as the input.
func BaseToBase(digits []int, base_from, base_to int) ([]int, error) {
	if base_from == base_to {
		return digits, nil
	}

	if base_from <= 0 {
		return nil, NewErrInvalidBase("base_from")
	} else if base_to <= 0 {
		return nil, NewErrInvalidBase("base_to")
	}

	res, _ := BaseToDec(digits, base_from)
	new_digits, _ := DecToBase(res, base_to)

	return new_digits, nil
}

// GetDigits returns the digits of a number.
//
// Parameters:
//   - number: The number to get the digits of.
//
// Returns:
//   - []int: The digits of the number.
func GetDigits(number int) []int {
	digits, _ := DecToBase(number, Decimal)
	return digits
}
