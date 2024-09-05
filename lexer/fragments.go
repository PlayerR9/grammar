package lexer

import (
	"fmt"
	"io"
	"strings"

	gcers "github.com/PlayerR9/go-commons/errors"
	dbg "github.com/PlayerR9/go-debug/assert"
	internal "github.com/PlayerR9/grammar/internal"
)

// FragNewline lexes a newline if it is found.
//
// Parameters:
//   - l: The lexer.
//
// Returns:
//   - string: The newline.
//   - error: An error if the newline is not found or if the newline is not a newline.
//
// Errors:
//   - NotFound: If the newline is not found.
//   - *gcers.ErrInvalidParameter: If the lexer is nil.
//   - error: any other error.
func FragNewline[T internal.TokenTyper](l *Lexer[T]) (string, error) {
	if l == nil {
		return "", gcers.NewErrNilParameter("l")
	}

	c, err := l.NextRune()
	if err == io.EOF {
		return "", NotFound
	} else if err != nil {
		return "", err
	}

	switch c {
	case '\n':
		return "\n", nil
	case '\r':
		next_c, err := l.NextRune()
		if err != nil {
			err := l.RefuseRune()
			dbg.AssertErr(err, "l.RefuseRune()")
		}

		if err == io.EOF {
			return "", fmt.Errorf("expected %q after %q, got nothing instead", '\n', '\r')
		} else if err != nil {
			return "", err
		}

		if next_c != '\n' {
			err := l.RefuseRune()
			dbg.AssertErr(err, "l.RefuseRune()")

			return "", fmt.Errorf("expected %q after %q, got %q instead", '\n', '\r', next_c)
		}

		return "\n", nil
	default:
		err := l.RefuseRune()
		dbg.AssertErr(err, "l.RefuseRune()")

		return "", NotFound
	}
}

// LexGroup is a helper function for lexing a group of characters that satisfy a given predicate according
// to the following rule:
//
//	group+
//
// Parameters:
//   - l: The lexer.
//   - is_func: The predicate function.
//
// Returns:
//   - string: The group of characters.
//   - error: An error if any.
//
// Errors:
//   - NotFound: When the group cannot be found because it doesn't exist in the text.
//   - errors.ErrInvalidParameter: If is_func or l is nil.
//   - any other error returned by l.NextRune() function.
func LexGroup[T internal.TokenTyper](l *Lexer[T], is_func func(c rune) bool) (string, error) {
	if is_func == nil {
		return "", gcers.NewErrNilParameter("is_func")
	} else if l == nil {
		return "", gcers.NewErrNilParameter("l")
	}

	c, err := l.NextRune()
	if err == io.EOF {
		return "", NotFound
	} else if err != nil {
		return "", err
	}

	if !is_func(c) {
		err := l.RefuseRune()
		dbg.AssertErr(err, "l.RefuseRune()")

		return "", NotFound
	}

	var builder strings.Builder
	builder.WriteRune(c)

	for {
		c, err := l.NextRune()
		if err == io.EOF {
			break
		} else if err != nil {
			tmp := l.RefuseRune()
			dbg.AssertErr(tmp, "l.RefuseRune()")

			return "", err
		}

		if !is_func(c) {
			err := l.RefuseRune()
			dbg.AssertErr(err, "l.RefuseRune()")

			break
		}

		builder.WriteRune(c)
	}

	return builder.String(), nil
}
