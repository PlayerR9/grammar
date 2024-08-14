package lexing

import (
	"errors"
	"io"

	gcers "github.com/PlayerR9/go-commons/errors"
)

var (
	// Done is the error returned when the lexing process is done without error and before reaching the end of the
	// stream. Readers must return Done itself and not wrap it as callers will test this error using ==.
	Done error
)

func init() {
	Done = errors.New("done")
}

// LexFunc is a function that can be used with RightLex.
//
// Parameters:
//   - scanner: The rune scanner.
//
// Returns:
//   - []rune: The list of runes.
//   - error: The error.
//
// This function must assume scanner is never nil. Moreover, use io.EOF to signify the end of the stream.
// Lastly, the error Done is returned if the lexing process is done before reaching the end of the stream.
//
// Finally, it is suggested to always push back the last rune read if any error that is not io.EOF is returned.
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
//   - any other error returned by lex_f.
func RightLex(scanner io.RuneScanner, lex_f LexFunc) ([]rune, error) {
	if scanner == nil {
		return nil, gcers.NewErrNilParameter("scanner")
	} else if lex_f == nil {
		return nil, gcers.NewErrNilParameter("lex_f")
	}

	var chars []rune

	for {
		curr, err := lex_f(scanner)
		if err == io.EOF || err == Done {
			chars = append(chars, curr...)

			break
		}

		if err != nil {
			return nil, err
		}

		chars = append(chars, curr...)
	}

	return chars, nil
}
