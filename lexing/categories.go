package lexing

import (
	"errors"
	"fmt"
	"io"

	dbg "github.com/PlayerR9/go-debug/assert"
)

var (
	// NoMatch is the error that is returned when there is no match. Readers must return NoMatch
	// itself and not wrap it as callers will test this error using ==.
	NoMatch error
)

func init() {
	NoMatch = errors.New("no match")
}

// LexMode is the lexing mode of a category.
type LexMode int

const (
	// LexOne is the lexing mode that lexes one character at a time. No more, no less.
	LexOne LexMode = iota

	// LexMany is the lexing mode that lexes many characters at a time; including
	// no characters at all.
	LexMany

	// LexOptional is the lexing mode that lexes one character at a time. However, no
	// character is also accepted.
	LexOptional

	// MustLexMany is the lexing mode that lexes many characters at a time. However,
	// there must be at least one matching character.
	MustLexMany
)

// LexCategory reads the content of the stream and returns the list of runes according to the given function.
//
// Parameters:
//   - scanner: The rune scanner.
//   - is: The function that checks if a character is part of the category.
//   - mode: The lexing mode.
//
// Returns:
//   - []rune: The list of runes.
//   - error: The error.
//
// Errors:
//   - NoMatch: When a match is not found, regardless of whether the end of the stream is reached or not.
//   - io.EOF: When the end of the stream is reached but a match was found.
//   - any other error that is not NoMatch or io.EOF.
func LexCategory(scanner io.RuneScanner, is func(rune) bool, mode LexMode) ([]rune, error) {
	if is == nil {
		return nil, NoMatch
	}

	var chars []rune

	switch mode {
	case LexOne:
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
			err = scanner.UnreadRune()
			dbg.AssertErr(err, "scanner.UnreadRune()")

			return nil, NoMatch
		}

		chars = []rune{c}
	case LexMany:
		if scanner == nil {
			return nil, nil
		}

		for {
			c, _, err := scanner.ReadRune()
			if err == io.EOF {
				break
			} else if err != nil {
				return chars, err
			}

			if !is(c) {
				err = scanner.UnreadRune()
				dbg.AssertErr(err, "scanner.UnreadRune()")

				break
			}

			chars = append(chars, c)
		}
	case LexOptional:
		if scanner == nil {
			return nil, nil
		}

		c, _, err := scanner.ReadRune()
		if err == io.EOF {
			return nil, nil
		} else if err != nil {
			return nil, err
		}

		if !is(c) {
			err = scanner.UnreadRune()
			dbg.AssertErr(err, "scanner.UnreadRune()")
		} else {
			chars = append(chars, c)
		}
	case MustLexMany:
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
				err = scanner.UnreadRune()
				dbg.AssertErr(err, "scanner.UnreadRune()")

				break
			}

			chars = append(chars, c)
		}

		if len(chars) == 0 {
			return nil, NoMatch
		}
	default:
		return nil, fmt.Errorf("invalid mode: %d", mode)
	}

	return chars, nil
}

// MakeCategoryLexer creates a new LexFunc that reads the content of the stream and returns the list of runes
// according to the given function.
//
// Parameters:
//   - is: The function that checks if a character is part of the category.
//   - mode: The lexing mode.
//
// Returns:
//   - LexFunc: The new lexing function. Nil if is is nil.
func MakeCategoryLexer(is func(c rune) bool, mode LexMode) LexFunc {
	if is == nil {
		return nil
	}

	var f LexFunc

	switch mode {
	case LexOne:
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
				err = scanner.UnreadRune()
				dbg.AssertErr(err, "scanner.UnreadRune()")

				return nil, NoMatch
			}

			return []rune{c}, nil
		}
	case LexMany:
		f = func(scanner io.RuneScanner) ([]rune, error) {
			if scanner == nil {
				return nil, nil
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
					err = scanner.UnreadRune()
					dbg.AssertErr(err, "scanner.UnreadRune()")

					break
				}

				chars = append(chars, c)
			}

			return chars, nil
		}
	case LexOptional:
		f = func(scanner io.RuneScanner) ([]rune, error) {
			if scanner == nil {
				return nil, nil
			}

			c, _, err := scanner.ReadRune()
			if err != nil {
				if err != io.EOF {
					return nil, err
				}

				return nil, nil
			}

			if !is(c) {
				err = scanner.UnreadRune()
				dbg.AssertErr(err, "scanner.UnreadRune()")

				return nil, nil
			}

			return []rune{c}, nil
		}
	case MustLexMany:
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
					err = scanner.UnreadRune()
					dbg.AssertErr(err, "scanner.UnreadRune()")

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
