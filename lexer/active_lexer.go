package lexer

import (
	gr "github.com/PlayerR9/grammar/grammar"
	internal "github.com/PlayerR9/grammar/internal"
)

// LexFunc is the lexing function.
//
// Parameters:
//   - lexer: The lexer.
//
// Returns:
//   - string: The group of characters.
//   - error: An error if any.
//
// Errors:
//   - NotFound: When the group cannot be found because it doesn't exist in the text.
//   - errors.ErrInvalidParameter: If lexer is nil.
//   - any other error returned by lexer.NextRune() function or the text is invalid.
type LexFunc[T internal.TokenTyper] func(lexer *ActiveLexer[T]) (string, error)

// LexOnceFunc is the function that lexes the next token of the lexer.
//
// Parameters:
//   - lexer: The lexer. Assume that lexer is not nil.
//
// Returns:
//   - []*grammar.Token: The next tokens of the lexer.
//   - error: An error if the lexer encounters an error while lexing the next token.
type LexOnceFunc[T internal.TokenTyper] func(lexer *ActiveLexer[T]) ([]*gr.Token[T], error)

// ActiveLexer is the lexer of the grammar.
type ActiveLexer[T internal.TokenTyper] struct {
	// global is the global of the lexer.
	global *Lexer[T]

	// pos is the position of the lexer.
	pos int

	// tokens is the tokens of the lexer.
	tokens []*gr.Token[T]

	// err is the error of the lexer.
	err error
}

// HasError checks whether the lexer has an error.
//
// Returns:
//   - bool: True if the lexer has an error, false otherwise.
func (al ActiveLexer[T]) HasError() bool {
	return al.err != nil
}

// Error returns the error of the lexer.
//
// Returns:
//   - error: The error of the lexer.
func (al ActiveLexer[T]) Error() error {
	return al.err
}

// NextEvents returns the next events of the lexer.
//
// Returns:
//   - []*grammar.Token: The next events of the lexer.
func (al *ActiveLexer[T]) NextEvents() []*gr.Token[T] {
	tks, err := al.global.fn(al)
	if err != nil {
		al.err = err

		return nil
	}

	return tks
}

// WalkOne walks one token.
//
// Parameters:
//   - tk: The token.
//
// Returns:
//   - bool: True if the token is the same as T(0). False otherwise.
func (al *ActiveLexer[T]) WalkOne(tk *gr.Token[T]) bool {
	if tk == nil {
		return false
	}

	al.tokens = append(al.tokens, tk)

	return tk.Type == T(0)
}

/* // Lex lexes tokens in the input stream.
//
// Returns:
//   - error: An error if the lexer encounters an error.
func (l *ActiveLexer[T]) Lex() error {
	for {

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
} */

// Tokens returns the tokens of the lexer.
//
// Returns:
//   - []*grammar.Token: The tokens of the lexer.
func (l ActiveLexer[T]) Tokens() []*gr.Token[T] {
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
//   - bool: True if the next rune exists. False otherwise.
func (l *ActiveLexer[T]) PeekRune() (rune, bool) {
	char, ok := l.global.RuneAt(l.pos)
	if !ok {
		return 0, false
	}

	return char, true
}

// NextRune returns the next rune in the input stream.
//
// Returns:
//   - rune: The next rune in the input stream.
//   - bool: True if the next rune exists. False otherwise.
func (l *ActiveLexer[T]) NextRune() (rune, bool) {
	char, ok := l.global.RuneAt(l.pos)
	if !ok {
		return 0, false
	}

	l.pos += 1

	return char, true
}

// RefuseRune rejects the last read rune in the input stream.
//
// Returns:
//   - bool: True if the last read rune is rejected. False otherwise.
func (l *ActiveLexer[T]) RefuseRune() bool {
	if l.pos == 0 {
		return false
	}

	l.pos -= 1

	return true
}
