package runes

import (
	"slices"
	"unicode"
	"unicode/utf8"
)

// split_size returns the number of lines and the maximum line length.
//
// Parameters:
//   - data: The data to split.
//   - sep: The separator to use.
//
// Returns:
//   - int: The number of lines.
//   - int: The maximum line length.
func split_size(data []rune, sep rune) (int, int) {
	var count int
	var max int
	var current int

	for _, c := range data {
		if c == sep {
			count++

			if current > max {
				max = current
			}

			current = 0
		} else {
			current++
		}
	}

	if current != 0 {
		count++

		if current > max {
			max = current
		}
	}

	return count, max
}

// Split is a function that splits the data into lines. Returns nil if the data is empty.
//
// Parameters:
//   - data: The data to split.
//   - sep: The separator to use.
//
// Returns:
//   - [][]rune: The lines.
func Split(data []rune, sep rune) [][]rune {
	if len(data) == 0 {
		return nil
	}

	count, max := split_size(data, sep)

	lines := make([][]rune, 0, count)
	current_line := make([]rune, 0, max)

	for i := 0; i < len(data); i++ {
		if data[i] != sep {
			current_line = append(current_line, data[i])

			continue
		}

		lines = append(lines, current_line[:len(current_line):len(current_line)])

		current_line = make([]rune, 0, max)
	}

	if len(current_line) > 0 {
		lines = append(lines, current_line)
	}

	return lines
}

// JoinSize returns the number of runes in the data.
//
// Parameters:
//   - data: The data to join.
//
// Returns:
//   - int: The number of runes.
func JoinSize(data [][]rune) int {
	if len(data) == 0 {
		return 0
	}

	var size int

	for _, line := range data {
		size += len(line)
	}

	size += len(data) - 1

	return size
}

// Join is a function that joins the data. Returns nil if the data is empty.
//
// Parameters:
//   - data: The data to join.
//   - sep: The separator to use.
//
// Returns:
//   - []rune: The joined data.
func Join(data [][]rune, sep rune) []rune {
	if len(data) == 0 {
		return nil
	}

	size := JoinSize(data)

	result := make([]rune, 0, size)

	result = append(result, data[0]...)

	for _, line := range data[1:] {
		result = append(result, sep)
		result = append(result, line...)
	}

	return result
}

// Repeat is a function that repeats the character.
//
// Parameters:
//   - char: The character to repeat.
//   - count: The number of times to repeat the character.
//
// Returns:
//   - []rune: The repeated character. Returns nil if count is less than 0.
func Repeat(char rune, count int) []rune {
	if count < 0 {
		return nil
	} else if count == 0 {
		return []rune{}
	}

	chars := make([]rune, 0, count)

	for i := 0; i < count; i++ {
		chars = append(chars, char)
	}

	return chars
}

// LimitReverseLines is a function that limits the lines of the data in reverse order.
//
// Parameters:
//   - data: The data to limit.
//   - limit: The limit of the lines.
//
// Returns:
//   - []byte: The limited data.
func LimitReverseLines(data []rune, limit int) []rune {
	if len(data) == 0 {
		return nil
	}

	lines := Split(data, '\n')

	if limit == -1 || limit > len(lines) {
		limit = len(lines)
	}

	start_idx := len(lines) - limit

	lines = lines[start_idx:]

	return Join(lines, '\n')
}

// LimitLines is a function that limits the lines of the data.
//
// Parameters:
//   - data: The data to limit.
//   - limit: The limit of the lines.
//
// Returns:
//   - []byte: The limited data.
func LimitLines(data []rune, limit int) []rune {
	if len(data) == 0 {
		return nil
	}

	lines := Split(data, '\n')

	if limit == -1 || limit > len(lines) {
		limit = len(lines)
	}

	lines = lines[:limit]

	return Join(lines, '\n')
}

// ToInt converts a rune to an integer if possible. Conversion is case-insensitive and
// values from 0-9 and a-z are converted to 0-35.
//
// Parameters:
//   - char: The rune to convert.
//
// Returns:
//   - int: The converted integer.
//   - bool: True if the conversion was successful. False otherwise.
//
// Example:
//
//	digit, ok := ToInt('A')
//	if !ok {
//		panic("Could not convert 'A' to an integer")
//	}
//
//	fmt.Println(digit) // 10
func ToInt(char rune) (int, bool) {
	ok := unicode.IsDigit(char)
	if ok {
		return int(char - '0'), true
	}

	ok = unicode.IsLetter(char)
	if !ok {
		return 0, false
	}

	char = unicode.ToLower(char)

	return int(char - 'a' + 10), true
}

// FromInt converts an integer to a rune if possible. Conversion is case-insensitive and
// values from 0-9 and a-z are converted to 0-35.
//
// Parameters:
//   - digit: The integer to convert.
//
// Returns:
//   - rune: The converted rune.
//   - bool: True if the conversion was successful. False otherwise.
//
// Example:
//
//	char, ok := FromInt(10)
//	if !ok {
//		panic("Could not convert 10 to a rune")
//	}
//
//	fmt.Println(char) // 'A'
func FromInt(digit int) (rune, bool) {
	if digit < 0 || digit > 35 {
		return 0, false
	}

	if digit < 10 {
		return rune(digit + '0'), true
	}

	return rune(digit - 10 + 'a'), true
}

// BytesToUtf8 is a function that converts bytes to runes.
//
// Parameters:
//   - data: The bytes to convert.
//
// Returns:
//   - []rune: The runes.
//   - error: An error of type *ErrInvalidUTF8Encoding if the bytes are not
//     valid UTF-8.
//
// This function also converts '\r\n' to '\n'. Plus, whenever an error occurs, it returns the runes
// decoded so far and the index of the error rune.
func BytesToUtf8(data []byte) ([]rune, error) {
	if len(data) == 0 {
		return nil, nil
	}

	var chars []rune
	var i int

	for len(data) > 0 {
		c, size := utf8.DecodeRune(data)
		if c == utf8.RuneError {
			return chars, NewErrInvalidUTF8Encoding(i)
		}

		data = data[size:]
		i += size

		if c != '\r' {
			chars = append(chars, c)
			continue
		}

		if len(data) == 0 {
			return chars, NewErrInvalidUTF8Encoding(i)
		}

		c, size = utf8.DecodeRune(data)
		if c == utf8.RuneError {
			return chars, NewErrInvalidUTF8Encoding(i)
		}

		data = data[size:]
		i += size

		if c != '\n' {
			return chars, NewErrInvalidUTF8Encoding(i)
		}

		chars = append(chars, '\n')
	}

	return chars, nil
}

// StringToUtf8 converts a string to a slice of runes.
//
// Parameters:
//   - str: The string to convert.
//
// Returns:
//   - runes: The slice of runes.
//   - error: An error of type *ErrInvalidUTF8Encoding if the string is not
//     valid UTF-8.
//
// Behaviors:
//   - An empty string returns a nil slice with no errors.
//   - The function stops at the first invalid UTF-8 encoding; returning an
//     error and the runes found up to that point.
//   - The function converts '\r\n' to '\n'.
func StringToUtf8(str string) ([]rune, error) {
	if str == "" {
		return nil, nil
	}

	var chars []rune
	var i int

	for len(str) > 0 {
		c, size := utf8.DecodeRuneInString(str)
		if c == utf8.RuneError {
			return chars, NewErrInvalidUTF8Encoding(i)
		}

		str = str[size:]
		i += size

		if c != '\r' {
			chars = append(chars, c)
			continue
		}

		if len(str) == 0 {
			return chars, NewErrInvalidUTF8Encoding(i)
		}

		c, size = utf8.DecodeRuneInString(str)
		if c == utf8.RuneError {
			return chars, NewErrInvalidUTF8Encoding(i)
		}

		str = str[size:]
		i += size

		if c != '\n' {
			return chars, NewErrInvalidUTF8Encoding(i)
		}

		chars = append(chars, '\n')
	}

	return chars, nil
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
func IndicesOf(data []rune, sep rune, exclude_sep bool) []int {
	if len(data) == 0 {
		return nil
	}

	var indices []int

	for i := 0; i < len(data); i++ {
		if data[i] == sep {
			indices = append(indices, i)
		}
	}

	if len(indices) == 0 {
		return nil
	}

	if exclude_sep {
		for i := 0; i < len(indices); i++ {
			indices[i] += 1
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
//   - *luc.ErrInvalidParameter: If the openingToken or closingToken is an
//     empty string.
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
func FindContentIndexes(op_token, cl_token rune, tokens []rune) (result [2]int, err error) {
	result[0] = -1
	result[1] = -1

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
	} else if balance != 1 || cl_token != '\n' {
		err = NewErrTokenNotFound(cl_token, false)
		return
	}

	result[1] = len(tokens)
	return
}

// FixTabSize fixes the tab size by replacing it with a specified rune iff
// the tab size is greater than 0. The replacement rune is repeated for the
// specified number of times.
//
// Parameters:
//   - size: The size of the tab.
//   - rep: The replacement rune.
//
// Returns:
//   - []rune: The fixed tab size.
func FixTabSize(size int, rep rune) []rune {
	if size <= 0 {
		return []rune{'\t'}
	}

	return Repeat(rep, size)
}

// FilterNonEmpty removes zero runes from a slice of runes.
//
// Parameters:
//   - values: The slice of runes to trim.
//
// Returns:
//   - []rune: The slice of runes with zero runes removed.
func FilterNonEmpty(values []rune) []rune {
	if len(values) == 0 {
		return nil
	}

	var top int

	for i := 0; i < len(values); i++ {
		if values[i] != 0 && values[i] != '\000' {
			values[top] = values[i]
			top++
		}
	}

	return values[:top:top]
}
