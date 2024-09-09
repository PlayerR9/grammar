package lexer

import (
	"fmt"
	"io"

	gcch "github.com/PlayerR9/go-commons/runes"
	gr "github.com/PlayerR9/grammar/grammar"
)

// Lexer is a lexer.
type Lexer[T gr.Enumer] struct {
	// chars is the characters left in the input stream.
	chars []rune

	// prev_pos is the previous position in the input stream.
	prev_pos int

	// curr_pos is the current position in the input stream.
	curr_pos int

	// tokens is the list of tokens lexed so far.
	tokens []*gr.Token[T]

	// table is the table of lexing functions.
	table map[rune]LexFunc[T]

	// def_fn is the default lexing function.
	def_fn LexFunc[T]
}

// NextRune advances the lexer to the next rune in the input stream.
//
// Returns:
//   - rune: The next rune in the input stream.
//   - bool: True if the next rune is in the input stream, false otherwise.
func (l *Lexer[T]) NextRune() (rune, bool) {
	if len(l.chars) == 0 {
		return 0, false
	}

	r := l.chars[0]
	l.chars = l.chars[1:]

	l.curr_pos++

	return r, true
}

// PeekRune returns the next rune in the input stream without consuming it.
//
// Returns:
//   - rune: The next rune in the input stream.
//   - bool: True if the next rune is in the input stream, false otherwise.
func (l Lexer[T]) PeekRune() (rune, bool) {
	if len(l.chars) == 0 {
		return 0, false
	}

	return l.chars[0], true
}

// lex_one is a helper function that lexes a single token.
//
// Returns:
//   - *Token: The token that was lexed.
//   - error: An error if the token could not be lexed.
//
// Nil tokens are ignored.
func (l *Lexer[T]) lex_one(char rune) (*gr.Token[T], error) {
	fn, ok := l.table[char]
	if ok {
		tk, err := fn(l)
		if err != nil {
			return nil, err
		}

		return tk, nil
	}

	if l.def_fn == nil {
		return nil, fmt.Errorf("unexpected character %q", char)
	}

	tk, err := l.def_fn(l)
	if err != nil {
		return nil, err
	}

	return tk, nil
}

// Tokens is a function that returns the list of tokens. The last token
// is guaranteed to be an EOF token.
//
// Parameters:
//   - tokens: The list of tokens.
//
// Returns:
//   - []*Token: The list of tokens with an EOF token added to the end.
func (l *Lexer[T]) Tokens() []*gr.Token[T] {
	tk_eof := gr.NewTerminalToken(T(0), "")
	tk_eof.Pos = -1

	tokens := append(l.tokens, tk_eof)

	for i := 0; i < len(tokens)-1; i++ {
		tokens[i].Lookahead = tokens[i+1]
	}

	return tokens
}

// SetInputStream sets the input stream for the lexer.
//
// Parameters:
//   - data: The input stream to set.
//
// Returns:
//   - error: An error if the input stream could not be set.
func (l *Lexer[T]) SetInputStream(data []byte) error {
	chars, err := gcch.BytesToUtf8(data)
	if err != nil {
		return err
	}

	l.chars = chars

	return nil
}

// Lex lexes the input stream and returns a list of tokens.
//
// Parameters:
//   - data: The input stream to lex.
//
// Returns:
//   - error: An error if the input stream could not be lexed.
func (l *Lexer[T]) Lex() error {
	if l.chars == nil {
		l.tokens = make([]*gr.Token[T], 0)
	} else {
		l.tokens = l.tokens[:0]
	}

	for len(l.chars) > 0 {
		char := l.chars[0]

		tk, err := l.lex_one(char)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if tk != nil {
			tk.Pos = l.prev_pos
			l.tokens = append(l.tokens, tk)
		}

		l.prev_pos = l.curr_pos
	}

	return nil
}
