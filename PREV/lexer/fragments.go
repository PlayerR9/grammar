package lexer

import (
	"fmt"
	"strings"

	gcers "github.com/PlayerR9/go-commons/errors"
	dbg "github.com/PlayerR9/go-debug/assert"
	internal "github.com/PlayerR9/grammar/PREV/internal"
)

// LexGroup is a helper function for lexing a group of characters that satisfy a given predicate according
// to the following rule:
//
//	group+
//
// Parameters:
//   - lexer: The lexer.
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
func LexGroup[T internal.TokenTyper](lexer *ActiveLexer[T], is_func func(c rune) bool) (string, error) {
	if lexer == nil {
		return "", gcers.NewErrNilParameter("lexer")
	} else if is_func == nil {
		return "", gcers.NewErrNilParameter("is_func")
	}

	c, ok := lexer.NextRune()
	if !ok {
		return "", NotFound
	}

	if !is_func(c) {
		ok := lexer.RefuseRune()
		dbg.AssertOk(ok, "lexer.RefuseRune()")

		return "", NotFound
	}

	var builder strings.Builder
	builder.WriteRune(c)

	for {
		c, ok := lexer.NextRune()
		if !ok {
			break
		}

		if !is_func(c) {
			ok := lexer.RefuseRune()
			dbg.AssertOk(ok, "lexer.RefuseRune()")

			break
		}

		builder.WriteRune(c)
	}

	return builder.String(), nil
}

// LexAtLeastOne is a convenience function for lexing text according to the lexing function
// such that at least one character is found and continues until the end of the text (or the valid input)
// is reached.
//
// Parameters:
//   - lexer: The lexer.
//   - lex_func: The lexing function.
//
// Returns:
//   - string: Text that was lexed.
//   - error: An error if any.
//
// Errors:
//   - NotFound: When the group cannot be found because it doesn't exist in the text.
//   - errors.ErrInvalidParameter: If lexer or lex_func is nil.
//   - any other error returned by lex_func.
func LexAtLeastOne[T internal.TokenTyper](lexer *ActiveLexer[T], lex_func LexFunc[T]) (string, error) {
	if lexer == nil {
		return "", gcers.NewErrNilParameter("lexer")
	} else if lex_func == nil {
		return "", gcers.NewErrNilParameter("lex_func")
	}

	str, err := lex_func(lexer)
	if err != nil {
		return str, err
	}

	var builder strings.Builder
	builder.WriteString(str)

	for {
		str, err := lex_func(lexer)
		if err == NotFound {
			break
		} else if err != nil {
			return builder.String(), err
		}

		builder.WriteString(str)
	}

	return builder.String(), nil
}

// FragLiteral lexes a literal if it is found.
//
// Parameters:
//   - lexer: The lexer.
//   - chars: The literal.
//
// Returns:
//   - string: The literal.
//   - error: An error if the literal is not found or if the literal is not a literal.
//
// Errors:
//   - NotFound: If the literal is not found.
//   - *gcers.ErrInvalidParameter: If the lexer is nil or the chars is empty.
//   - error: any other error.
func FragLiteral[T internal.TokenTyper](lexer *ActiveLexer[T], chars []rune) (string, error) {
	if lexer == nil {
		return "", gcers.NewErrNilParameter("lexer")
	} else if len(chars) == 0 {
		return "", gcers.NewErrInvalidParameter("chars", gcers.NewErrEmpty(chars))
	}

	char, ok := lexer.NextRune()
	if !ok {
		return "", NotFound
	}

	if char != chars[0] {
		ok := lexer.RefuseRune()
		dbg.AssertOk(ok, "lexer.RefuseRune()")

		return "", NotFound
	}

	var builder strings.Builder
	builder.WriteRune(char)

	for i := 1; i < len(chars); i++ {
		c, ok := lexer.NextRune()
		if !ok {
			return builder.String(), fmt.Errorf("expected %q after %q, got nothing instead", chars[i], builder.String())
		} else if c != chars[i] {
			return builder.String(), fmt.Errorf("expected %q after %q, got %q instead", chars[i], builder.String(), c)
		}

		builder.WriteRune(c)
	}

	return builder.String(), nil
}

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
func FragNewline[T internal.TokenTyper](lexer *ActiveLexer[T]) (string, error) {
	if lexer == nil {
		return "", gcers.NewErrNilParameter("l")
	}

	c, ok := lexer.NextRune()
	if !ok {
		return "", NotFound
	}

	switch c {
	case '\n':
		return "\n", nil
	case '\r':
		next_c, ok := lexer.NextRune()
		if !ok {
			return "", fmt.Errorf("expected %q after %q, got nothing instead", '\n', '\r')
		}

		if next_c != '\n' {
			ok := lexer.RefuseRune()
			dbg.AssertOk(ok, "lexer.RefuseRune()")

			return "", fmt.Errorf("expected %q after %q, got %q instead", '\n', '\r', next_c)
		}

		return "\n", nil
	default:
		ok := lexer.RefuseRune()
		dbg.AssertOk(ok, "lexer.RefuseRune()")

		return "", NotFound
	}
}
