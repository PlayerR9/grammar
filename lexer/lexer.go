package lexer

import (
	"fmt"

	gr "github.com/PlayerR9/grammar/grammar"
	luc "github.com/PlayerR9/lib_units/common"
)

// Lexer is an interface that defines the behavior of a lexer.
type Lexer[S gr.TokenTyper] interface {
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
	//   - *gr.Token[S]: The token of the lexer.
	//   - error: An error if the lexer encounters an error while lexing the next token.
	//
	// If the returned token is nil, then it is marked as 'to skip' and, as a result,
	// not added to the list of tokens.
	LexOne() (*gr.Token[S], error)
}

// get_tokens returns the tokens of the lexer.
//
// Parameters:
//   - tokens: The tokens of the lexer.
//
// Returns:
//   - []T: The tokens of the lexer.
func get_tokens[S gr.TokenTyper](tokens []*gr.Token[S]) []*gr.Token[S] {
	eof_tk := &gr.Token[S]{
		Type:      S(0),
		Data:      "",
		At:        -1,
		Lookahead: nil,
	}

	tokens = append(tokens, eof_tk)
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
//   - []*grammar.Token[T]: The tokens of the lexer.
//   - error: An error if the lexer encounters an error while lexing the input stream.
//
// This function always returns at least one token and the last one is
// always the EOF token.
//
// This function is just a convenience function that calls the SetInputStream, Lex, and
// GetTokens methods of the lexer.
func FullLex[S gr.TokenTyper](lexer Lexer[S], data []byte) ([]*gr.Token[S], error) {
	if lexer == nil {
		tokens := get_tokens[S](nil)

		return tokens, luc.NewErrNilParameter("lexer")
	}

	lexer.SetInputStream(data)

	lexer.Reset()

	var tokens []*gr.Token[S]

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
