package lexing

import (
	"errors"
	"io"
	"unicode"
)

var (
	// NoMatch is the error that is returned when there is no match. Readers must return NoMatch
	// itself and not wrap it as callers will test this error using ==.
	NoMatch error
)

func init() {
	NoMatch = errors.New("no match")
}

// LexCategory reads the content of the stream and returns the list of runes according to the given function.
//
// Parameters:
//   - scanner: The rune scanner.
//   - is: The function that checks if a character is part of the category.
//   - lex_one: Whether to lex one or more times.
//
// Returns:
//   - []rune: The list of runes.
//   - error: The error.
//
// Errors:
//   - NoMatch: When a match is not found, regardless of whether the end of the stream is reached or not.
//   - io.EOF: When the end of the stream is reached but a match was found.
//   - any other error that is not NoMatch or io.EOF.
func LexCategory(scanner io.RuneScanner, is func(rune) bool, lex_one bool) ([]rune, error) {
	if is == nil {
		return nil, NoMatch
	}

	var chars []rune

	if lex_one {
		if scanner == nil {
			return nil, NoMatch
		}

		c, _, err := scanner.ReadRune()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}

			return nil, NoMatch
		}

		if !is(c) {
			_ = scanner.UnreadRune()
			// dbg.AssertErr(err, "scanner.UnreadRune()")

			return nil, NoMatch
		}

		chars = []rune{c}
	} else {
		if scanner == nil {
			return nil, NoMatch
		}

		for {
			c, _, err := scanner.ReadRune()
			if err == io.EOF {
				break
			} else if err != nil {
				return chars, err
			}

			if !is(c) {
				_ = scanner.UnreadRune()
				// dbg.AssertErr(err, "scanner.UnreadRune()")

				break
			}

			chars = append(chars, c)
		}

		if len(chars) == 0 {
			return nil, NoMatch
		}
	}

	return chars, nil
}

// MakeCategoryLexer creates a new LexFunc that reads the content of the stream and returns the list of runes
// according to the given function.
//
// Parameters:
//   - is: The function that checks if a character is part of the category.
//   - lex_one: Whether to lex one or more times.
//
// Returns:
//   - LexFunc: The new lexing function. Nil if is is nil.
func MakeCategoryLexer(is func(c rune) bool, lex_one bool) LexFunc {
	if is == nil {
		return nil
	}

	var f LexFunc

	if lex_one {
		f = func(scanner io.RuneScanner) ([]rune, error) {
			if scanner == nil {
				return nil, NoMatch
			}

			c, _, err := scanner.ReadRune()
			if err != nil {
				if err != io.EOF {
					return nil, err
				}

				return nil, NoMatch
			}

			if !is(c) {
				_ = scanner.UnreadRune()
				// dbg.AssertErr(err, "scanner.UnreadRune()")

				return nil, NoMatch
			}

			return []rune{c}, nil
		}
	} else {
		f = func(scanner io.RuneScanner) ([]rune, error) {
			if scanner == nil {
				return nil, NoMatch
			}

			var chars []rune

			for {
				c, _, err := scanner.ReadRune()
				if err == io.EOF {
					break
				} else if err != nil {
					return chars, err
				}

				if !is(c) {
					_ = scanner.UnreadRune()
					// dbg.AssertErr(err, "scanner.UnreadRune()")

					break
				}

				chars = append(chars, c)
			}

			if len(chars) == 0 {
				return nil, NoMatch
			}

			return chars, nil
		}
	}

	return f
}

var (
	// [0-9]+
	CatDigits LexFunc

	// [A-Z]+
	CatUppercases LexFunc

	// [a-z]+
	CatLowercases LexFunc
)

func init() {
	// [0-9]+
	CatDigits = MakeCategoryLexer(unicode.IsDigit, false)

	// [A-Z]+
	CatUppercases = MakeCategoryLexer(unicode.IsUpper, false)

	// [a-z]+
	CatLowercases = MakeCategoryLexer(unicode.IsLower, false)
}
