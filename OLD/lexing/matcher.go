package lexing

import (
	"io"

	gcers "github.com/PlayerR9/go-commons/errors"
)

// LexFunc is a function that can be used with RightLex.
//
// Parameters:
//   - scanner: The rune scanner. Assumed to never be nil.
//
// Returns:
//   - []rune: The list of runes.
//   - error: The error.
//
// Errors:
//   - NoMatch: When a match is not found, regardless of whether the end of the stream is reached or not.
//   - io.EOF: When the end of the stream is reached but a match was found.
//   - any other error that is not NoMatch or io.EOF.
//
// Notes: Always push back the last rune read if any error that is not io.EOF is returned. That's because
// this function always reads the next rune from the stream and, if not properly pushed back, it will
// skip the next rune.
type LexFunc func(scanner io.RuneScanner) ([]rune, error)

// RightLex reads the content of the stream and returns the list of runes according to the given function.
//
// Parameters:
//   - scanner: The rune scanner.
//   - lex_f: The lexing function.
//
// Returns:
//   - []rune: The list of runes.
//   - error: The error.
//
// Errors:
//   - *errors.ErrInvalidParameter: When scanner or lex_f is nil.
//   - NoMatch: When a match is not found, regardless of whether the end of the stream is reached or not.
//   - any other error that is returned by lex_f that is not io.EOF.
func RightLex(scanner io.RuneScanner, lex_f LexFunc) ([]rune, error) {
	if scanner == nil {
		return nil, gcers.NewErrNilParameter("scanner")
	} else if lex_f == nil {
		return nil, gcers.NewErrNilParameter("lex_f")
	}

	// First match
	chars, err := lex_f(scanner)
	if err != nil {
		if err != io.EOF {
			return chars, err
		}

		return chars, nil
	}

	// Subsequent matches
	var tmp []rune

	for err == nil {
		tmp, err = lex_f(scanner)
		chars = append(chars, tmp...)

		if err != nil {
			if err == io.EOF || err == NoMatch {
				return chars, nil
			}

			return chars, err
		}
	}

	if err != io.EOF && err != NoMatch {
		return chars, err
	}

	return chars, nil
}
