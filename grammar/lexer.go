package grammar

import (
	"fmt"

	luc "github.com/PlayerR9/lib_units/common"
)

// Lexer is an interface that defines the behavior of a lexer.
type Lexer[T TokenTyper] interface {
	// SetInputStream sets the input stream of the lexer.
	//
	// Parameters:
	//   - data: The input stream of the lexer.
	SetInputStream(data []byte)

	// Reset resets the lexer.
	//
	// This utility function allows to reset the information contained in the lexer
	// so that it can be used multiple times.
	Reset()

	// IsDone checks if the lexer is done.
	//
	// Returns:
	//   - bool: True if the lexer is done, false otherwise.
	IsDone() bool

	// LexOne lexes the next token of the lexer.
	//
	// Returns:
	//   - *Token[T]: The token of the lexer.
	//   - error: An error if the lexer encounters an error while lexing the next token.
	//
	// If the token lexed is marked as 'to skip', then the return value will be nil, nil instead.
	LexOne() (*Token[T], error)
}

// get_tokens returns the tokens of the lexer.
//
// Parameters:
//   - tokens: The tokens of the lexer.
//
// Returns:
//   - []*Token[T]: The tokens of the lexer.
func get_tokens[T TokenTyper](tokens []*Token[T]) []*Token[T] {
	eof_tok, err := NewToken(T(0), "", -1, nil)
	luc.AssertErr(err, "NewToken(%s, %q, %d, nil)", T(0).String(), "", -1)

	tokens = append(tokens, eof_tok)
	if len(tokens) == 1 {
		return tokens
	}

	prev := tokens[0]

	for _, next := range tokens[1:] {
		prev.Lookahead = next
		prev = next
	}

	return tokens
}

// FullLex lexes the input stream of the lexer and returns the tokens.
//
// Parameters:
//   - lexer: The lexer.
//   - data: The input stream of the lexer.
//
// Returns:
//   - []*Token[T]: The tokens of the lexer.
//   - error: An error if the lexer encounters an error while lexing the input stream.
//
// This function always returns at least one token and the last one is
// always the EOF token.
//
// This function is just a convenience function that calls the SetInputStream, Lex, and
// GetTokens methods of the lexer.
func FullLex[T TokenTyper](lexer Lexer[T], data []byte) ([]*Token[T], error) {
	if lexer == nil {
		tokens := get_tokens[T](nil)

		return tokens, luc.NewErrNilParameter("lexer")
	}

	lexer.SetInputStream(data)

	lexer.Reset()

	var tokens []*Token[T]

	for !lexer.IsDone() {
		tk, err := lexer.LexOne()
		if err != nil {
			tokens = get_tokens(tokens)
			return tokens, fmt.Errorf("error while lexing: %w", err)
		}

		if tk != nil {
			tokens = append(tokens, tk)
		}
	}

	tokens = get_tokens(tokens)

	return tokens, nil
}
