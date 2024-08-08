package lexer

import (
	"fmt"
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

	// matcher is the matcher of the lexer.
	matcher *Matcher[S]
}

// NewLexer creates a new lexer.
//
// Parameters:
//   - lex_one_func: The function that lexes the next token of the lexer.
//
// Returns:
//   - *Lexer: The new lexer.
//
// It returns nil if the lex_one_func is nil.
func NewLexer[S gr.TokenTyper](lex_one_func LexOneFunc[S], matcher *Matcher[S]) *Lexer[S] {
	if lex_one_func == nil {
		return nil
	}

	return &Lexer[S]{
		lex_one: lex_one_func,
		matcher: matcher,
	}
}

// Reset resets the lexer.
//
// This utility function allows to reset the information contained in the lexer
// so that it can be used multiple times.
func (l *Lexer[S]) Reset() {
	l.tokens = l.tokens[:0]
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
func (l *Lexer[S]) FullLex(data []byte) ([]*gr.Token[S], error) {
	l.Init(data)

	l.Reset()

	var tokens []*gr.Token[S]

	for {
		tk, err := l.lex_one(l)
		if err == io.EOF {
			break
		}

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

/* // MatchChars matches the next characters of the lexer.
//
// Parameters:
//   - lexer: The lexer.
//   - chars: The characters to match.
//
// Returns:
//   - string: The matched characters.
//   - error: An error if the next characters do not match.
func (l *Lexer[S]) Match() (string, error) {
	if l.matcher == nil {
		return "", nil
	}

	match, err := l.matcher.Match(l)
	if err != nil {
		return "", err
	}


	if len(chars) == 0 {
		return "", nil
	}

	var prev *rune
	var builder strings.Builder

	for _, char := range chars {
		c, err := l.Next()
		if IsExhausted(err) {
			return builder.String(), NewErrUnexpectedRune(prev, nil, char)
		} else if err != nil {
			return builder.String(), err
		}

		if c != char {
			return builder.String(), NewErrUnexpectedRune(prev, &c, char)
		}

		builder.WriteRune(c)
		prev = &c
	}

	return builder.String(), nil
} */

/* // MatchRegex matches the next regex of the lexer.
//
// Parameters:
//   - regex: The regex to match.
//
// Returns:
//   - string: The matched regex.
//   - bool: True if the next regex matches, false otherwise.
func (l *Lexer[S]) MatchRegex(regex *regexp.Regexp) (string, bool) {
	if regex == nil {
		return "", false
	}

	l.input_stream.

	match := regex.Find(l.input_stream)

	if len(match) == 0 {
		return "", false
	}

	l.input_stream = l.input_stream[len(match):]
	l.at += len(match)

	return string(match), true
} */

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
