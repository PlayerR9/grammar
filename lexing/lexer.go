package lexing

import (
	"io"

	gcch "github.com/PlayerR9/go-commons/runes"
	gr "github.com/PlayerR9/grammar/grammar"
)

// LexOneFunc is the function that lexes the next token of the lexer.
//
// Parameters:
//   - lexer: The lexer. Assume that lexer is not nil.
//
// Returns:
//   - *Token: The next token of the lexer.
//   - error: An error if the lexer encounters an error while lexing the next token.
//
// It must return io.EOF if the lexer has reached the end of the input stream.
type LexOneFunc[S gr.TokenTyper] func(lexer *Lexer[S]) (*gr.Token[S], error)

// Lexer is the lexer of the grammar.
type Lexer[S gr.TokenTyper] struct {
	// input_stream is the input stream of the lexer.
	gcch.CharStream

	// tokens is the tokens of the lexer.
	tokens []*gr.Token[S]

	// lex_one is the function that lexes the next token of the lexer.
	lex_one LexOneFunc[S]

	// Err is the error reason of the lexer.
	Err *ErrLexing

	// matcher is the matcher of the lexer.
	matcher *Matcher[S]

	// table is the lavenshtein table of the lexer.
	table *LavenshteinTable
}

// NewLexer creates a new lexer.
//
// Parameters:
//   - lex_one_func: The function that lexes the next token of the lexer.
//   - matcher: The matcher of the lexer. If matcher is nil, it won't be used.
//
// Returns:
//   - *Lexer: The new lexer.
//
// It returns nil if the lex_one_func is nil.
func NewLexer[S gr.TokenTyper](lex_one_func LexOneFunc[S], matcher *Matcher[S]) *Lexer[S] {
	if lex_one_func == nil {
		return nil
	}

	if matcher == nil {
		return &Lexer[S]{
			lex_one: lex_one_func,
		}
	} else {
		table, _ := NewLevenshteinTable(matcher.GetWords()...)

		return &Lexer[S]{
			lex_one: lex_one_func,
			matcher: matcher,
			table:   table,
		}
	}
}

// Reset resets the lexer.
//
// This utility function allows to reset the information contained in the lexer
// so that it can be used multiple times.
func (l *Lexer[S]) Reset() {
	gr.CleanTokens(l.tokens)
	l.tokens = l.tokens[:0]

	l.Err = nil
}

// make_error returns the error reason of the lexer.
//
// Parameters:
//   - reason: The error reason of the lexer.
//
// Returns:
//   - *ErrLexing: The error reason of the lexer.
//
// This function returns nil iff the lexer has no error.
func (l *Lexer[S]) make_error(reason error) *ErrLexing {
	if reason == nil || reason == io.EOF {
		return nil
	}

	var pos, delta int

	if len(l.tokens) < 2 {
		pos = 0
		delta = 1
	} else {
		last_tk := l.tokens[len(l.tokens)-2]

		pos = last_tk.At

		if last_tk.Data == "" {
			delta = 1
		} else {
			delta = len(last_tk.Data)
		}
	}

	return NewErrLexing(pos, delta, reason)
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
func (lexer *Lexer[S]) FullLex(data []byte) []*gr.Token[S] {
	lexer.Init(data)

	lexer.Reset()

	var tokens []*gr.Token[S]

	if lexer.matcher == nil {
		for !lexer.IsExhausted() && lexer.Err == nil {
			tmp, err := lexer.lex_one(lexer)
			if err != nil {
				lexer.Err = lexer.make_error(err)
			} else if tmp != nil {
				tokens = append(tokens, tmp)
			}
		}
	} else {
		for !lexer.IsExhausted() && lexer.Err == nil {
			at := lexer.Pos()

			match, _ := lexer.matcher.Match(lexer)

			if match.IsValidMatch() {
				symbol, data := match.GetMatch()

				tk := gr.NewToken(symbol, data, at, nil)
				tokens = append(tokens, tk)
			} else {
				tmp, err := lexer.lex_one(lexer)
				if err != nil {
					lexer.Err = lexer.make_error(err)

					str, err := lexer.table.Closest(match.chars, 2) // Magic number.
					if err == nil {
						lexer.Err.SetSuggestion("Did you mean '" + str + "'?")
					}
				} else if tmp != nil {
					tokens = append(tokens, tmp)
				}
			}
		}

	}

	tokens = get_tokens(tokens)
	return tokens
}
