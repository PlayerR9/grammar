package grammar

import (
	"fmt"
	"io"
	"strings"

	gcers "github.com/PlayerR9/go-commons/errors"
	gcch "github.com/PlayerR9/go-commons/runes"
	dbg "github.com/PlayerR9/go-debug/assert"
	gr "github.com/PlayerR9/grammar/grammar"
	internal "github.com/PlayerR9/grammar/internal"
)

var (
	// NotFound is an error that is returned when text cannot be lexed yet it is not due to an error
	// within the text. Readers must return this error as is and not wrap it as callers check for this
	// error with ==.
	NotFound error
)

func init() {
	NotFound = fmt.Errorf("not found")
}

// LexOneFunc is the function that lexes the next token of the lexer.
//
// Parameters:
//   - lexer: The lexer. Assume that lexer is not nil.
//
// Returns:
//   - *grammar.Token: The next token of the lexer.
//   - error: An error if the lexer encounters an error while lexing the next token.
type LexOneFunc[T internal.TokenTyper] func(l *Lexer[T]) (*gr.Token[T], error)

// Lexer is the lexer of the grammar.
type Lexer[T internal.TokenTyper] struct {
	// scanner is the scanner of the lexer.
	scanner io.RuneScanner

	// tokens is the tokens of the lexer.
	tokens []*gr.Token[T]

	// fn is the function that lexes the next token of the lexer.
	fn LexOneFunc[T]
}

// NewLexer creates a new lexer.
//
// Parameters:
//   - fn: The function that lexes the next token of the lexer.
//
// Returns:
//   - *Lexer[T]: The new lexer.
//   - error: An error if the lexing function is nil.
func NewLexer[T internal.TokenTyper](fn LexOneFunc[T]) (*Lexer[T], error) {
	if fn == nil {
		return nil, gcers.NewErrNilParameter("fn")
	}

	return &Lexer[T]{
		scanner: nil,
		tokens:  nil,
		fn:      fn,
	}, nil
}

func (l *Lexer[T]) SetInputStream(data []byte) {
	var stream gcch.CharStream

	stream.Init(data)

	l.scanner = &stream
}

func (l *Lexer[T]) Lex() error {
	// Clear previous tokens
	if len(l.tokens) > 0 {
		l.tokens = l.tokens[:0]
	}

	for {
		tk, err := l.fn(l)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if tk != nil {
			l.tokens = append(l.tokens, tk)
		}
	}

	return nil
}

func (l Lexer[T]) Tokens() []*gr.Token[T] {
	tokens := make([]*gr.Token[T], len(l.tokens), len(l.tokens)+1)
	copy(tokens, l.tokens)

	eof := gr.NewToken(T(0), "", nil)
	tokens = append(tokens, eof)

	for i := 0; i < len(tokens)-1; i++ {
		tokens[i].Lookahead = tokens[i+1]
	}

	return tokens
}

func (l *Lexer[T]) NextRune() (rune, error) {
	char, _, err := l.scanner.ReadRune()
	if err != nil {
		return '\000', err
	}

	return char, nil
}

func (l *Lexer[T]) RefuseRune() error {
	err := l.scanner.UnreadRune()
	return err
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
