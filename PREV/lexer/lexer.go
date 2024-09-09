package lexer

import (
	"fmt"
	"iter"

	gcbk "github.com/PlayerR9/go-commons/backup"
	gcch "github.com/PlayerR9/go-commons/runes"
	internal "github.com/PlayerR9/grammar/PREV/internal"
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

// Lexer is the lexer of the grammar.
type Lexer[T internal.TokenTyper] struct {
	// data is the scanner of the lexer.
	data []rune

	// fn is the function that lexes the next token of the lexer.
	fn LexOnceFunc[T]
}

// SetInputStream sets the input stream of the lexer.
//
// Parameters:
//   - data: The input stream of the lexer.
//
// Returns:
//   - error: An error if any.
func (l *Lexer[T]) SetInputStream(data []byte) error {
	runes, err := gcch.BytesToUtf8(data)
	if err != nil {
		return err
	}

	l.data = runes

	return nil
}

// Lex lexes tokens in the input stream.
//
// Returns:
//   - error: An error if the lexer encounters an error.
func (l *Lexer[T]) Lex() iter.Seq[*ActiveLexer[T]] {
	return gcbk.Execute(func() *ActiveLexer[T] {
		return &ActiveLexer[T]{
			global: l,
		}
	})
}

// RuneAt returns the rune at the given position.
//
// Parameters:
//   - pos: The position of the rune.
//
// Returns:
//   - rune: The rune.
//   - bool: True if the rune exists. False otherwise.
func (l Lexer[T]) RuneAt(pos int) (rune, bool) {
	if pos < 0 || pos >= len(l.data) {
		return 0, false
	}

	return l.data[pos], true
}
