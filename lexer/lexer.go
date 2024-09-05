package lexer

import (
	"fmt"
	"io"

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

// LexOnceFunc is the function that lexes the next token of the lexer.
//
// Parameters:
//   - lexer: The lexer. Assume that lexer is not nil.
//
// Returns:
//   - *grammar.Token: The next token of the lexer.
//   - error: An error if the lexer encounters an error while lexing the next token.
type LexOnceFunc[T internal.TokenTyper] func(lexer *Lexer[T]) (*gr.Token[T], error)

// Lexer is the lexer of the grammar.
type Lexer[T internal.TokenTyper] struct {
	// scanner is the scanner of the lexer.
	scanner io.RuneScanner

	// tokens is the tokens of the lexer.
	tokens []*gr.Token[T]

	// fn is the function that lexes the next token of the lexer.
	fn LexOnceFunc[T]
}

// NewLexer creates a new lexer.
//
// Parameters:
//   - fn: The function that lexes the next token of the lexer.
//
// Returns:
//   - *Lexer[T]: The new lexer.
//   - error: An error if the lexing function is nil.
func NewLexer[T internal.TokenTyper](fn LexOnceFunc[T]) (*Lexer[T], error) {
	if fn == nil {
		return nil, gcers.NewErrNilParameter("fn")
	}

	return &Lexer[T]{
		scanner: nil,
		tokens:  nil,
		fn:      fn,
	}, nil
}

// SetInputStream sets the input stream of the lexer.
//
// Parameters:
//   - data: The input stream of the lexer.
func (l *Lexer[T]) SetInputStream(data []byte) {
	var stream gcch.CharStream

	stream.Init(data)

	l.scanner = &stream
}

// Lex lexes tokens in the input stream.
//
// Returns:
//   - error: An error if the lexer encounters an error.
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

// Tokens returns the tokens of the lexer.
//
// Returns:
//   - []*grammar.Token: The tokens of the lexer.
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

// PeekRune returns the next rune in the input stream without consuming it.
//
// Returns:
//   - rune: The next rune in the input stream.
//   - error: An error if any.
func (l *Lexer[T]) PeekRune() (rune, error) {
	char, _, err := l.scanner.ReadRune()
	if err != nil {
		return '\000', err
	}

	err = l.scanner.UnreadRune()
	dbg.AssertErr(err, "l.scanner.UnreadRune()")

	return char, nil
}

// NextRune returns the next rune in the input stream.
//
// Returns:
//   - rune: The next rune in the input stream.
//   - error: An error if any.
func (l *Lexer[T]) NextRune() (rune, error) {
	char, _, err := l.scanner.ReadRune()
	if err != nil {
		return '\000', err
	}

	return char, nil
}

// RefuseRune rejects the last read rune in the input stream.
//
// Returns:
//   - error: An error if any.
func (l *Lexer[T]) RefuseRune() error {
	err := l.scanner.UnreadRune()
	return err
}
