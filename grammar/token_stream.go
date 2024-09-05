package grammar

import (
	"io"

	internal "github.com/PlayerR9/grammar/grammar/internal"
)

// TokenReader is an interface for reading tokens from a token stream.
type TokenReader[T internal.TokenTyper] interface {
	// ReadToken reads the next token from the token stream.
	//
	// Returns:
	//   - *Token[T]: The next token.
	//   - error: An error of type io.EOF if there are no more tokens.
	ReadToken() (*Token[T], error)
}

// TokenStream is a token stream.
type TokenStream[T internal.TokenTyper] struct {
	// tokens is the token stream.
	tokens []*Token[T]
}

// ReadToken implements the TokenReader interface.
func (ts *TokenStream[T]) ReadToken() (*Token[T], error) {
	if len(ts.tokens) == 0 {
		return nil, io.EOF
	}

	tk := ts.tokens[0]
	ts.tokens = ts.tokens[1:]

	return tk, nil
}

// NewTokenStream creates a new token stream.
//
// Parameters:
//   - tokens: The tokens in the token stream.
//
// Returns:
//   - *TokenStream[T]: The new token stream. Never returns nil.
func NewTokenStream[T internal.TokenTyper](tokens []*Token[T]) *TokenStream[T] {
	if tokens == nil {
		tokens = make([]*Token[T], 0)
	}

	return &TokenStream[T]{
		tokens: tokens,
	}
}
