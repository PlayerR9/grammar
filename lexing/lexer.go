package lexing

import (
	"errors"
	"io"
	"unicode/utf8"

	gccdm "github.com/PlayerR9/go-commons/CustomData/matcher"
	gcch "github.com/PlayerR9/go-commons/runes"
	gcstr "github.com/PlayerR9/go-commons/strings"
	dbg "github.com/PlayerR9/go-debug/assert"
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
	matcher gccdm.Matcher[S]

	// table is the lavenshtein table of the lexer.
	table *gccdm.LavenshteinTable

	// skipped is the number of skipped characters.
	skipped int
}

// WithLexFunc sets the function that lexes the next token of the lexer.
//
// Parameters:
//   - lex_one: The function that lexes the next token of the lexer.
//
// Use this to specify custom lexing functions that are not matched by the keyword matcher.
func (l *Lexer[S]) WithLexFunc(lex_one LexOneFunc[S]) {
	l.lex_one = lex_one
}

// Reset resets the lexer.
//
// This utility function allows to reset the information contained in the lexer
// so that it can be used multiple times.
//
// However, the internal table is not resetted for efficiency reasons. As such, calling
// functions such as AddToSkipRule() and AddToMatchRule() won't update the table after
// the first call; leading to inconsistencies.
//
// Make sure to prepare everything before calling this or the Lex function.
func (l *Lexer[S]) Reset() {
	gr.CleanTokens(l.tokens)
	l.tokens = l.tokens[:0]

	l.Err = nil
	l.skipped = 0

	if l.table == nil {
		var table gccdm.LavenshteinTable

		err := table.AddWords(l.matcher.GetWords())
		dbg.AssertErr(err, "table.AddWords(l.matcher.GetWords())")

		l.table = &table
	}
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
func (l Lexer[S]) make_error(reason error) *ErrLexing {
	if reason == nil || reason == io.EOF {
		return nil
	}

	var pos int

	if len(l.tokens) < 2 {
		pos = 0
	} else {
		last_tk := l.tokens[len(l.tokens)-2]

		pos = last_tk.At + len(last_tk.Data)
	}

	return NewErrLexing(pos+l.skipped, -1, reason)
}

// get_tokens returns the tokens of the lexer.
//
// Parameters:
//   - tokens: The tokens of the lexer.
//
// Returns:
//   - []T: The tokens of the lexer.
func (lexer Lexer[S]) get_tokens() []*gr.Token[S] {
	eof_tk := &gr.Token[S]{
		Type:      S(0),
		Data:      "",
		At:        -1,
		Lookahead: nil,
	}

	if len(lexer.tokens) == 0 {
		return []*gr.Token[S]{eof_tk}
	}

	tokens := make([]*gr.Token[S], len(lexer.tokens)+1)
	copy(tokens, lexer.tokens)

	tokens[len(tokens)-1] = eof_tk

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
//   - data: The input stream of the lexer.
//
// Returns:
//   - []*gr.Token[S]: The tokens of the lexer that were lexed so far.
//   - error: An error of type *ErrLexing if the lexing failed.
func (lexer *Lexer[S]) FullLex(data []byte) ([]*gr.Token[S], error) {
	lexer.Init(data)

	lexer.Reset()

	has_matcher := !lexer.matcher.IsEmpty()
	has_lexer := lexer.lex_one != nil

	if !has_matcher && !has_lexer {
		return lexer.get_tokens(), errors.New("no lexing function or matcher provided")
	}

	if has_matcher && has_lexer {
		for !lexer.IsExhausted() {
			at := lexer.Pos()

			match, _ := lexer.matcher.Match(lexer)

			if match.IsValidMatch() {
				if match.IsShouldSkip() {
					lexer.skip(match.GetChars())
				} else {
					symbol, data := match.GetMatch()

					tk := gr.NewToken(symbol, data, at, nil)

					lexer.tokens = append(lexer.tokens, tk)

					lexer.skipped = 0
				}
			} else {
				tmp, err := lexer.lex_one(lexer)
				if err != nil {
					lexer.Err = lexer.make_error(err)

					str, err := lexer.table.Closest(match.GetChars(), 2) // Magic number.
					if err == nil {
						lexer.Err.SetSuggestion("Did you mean '" + str + "'?")
					} else {
						words := lexer.matcher.GetRuleNames()
						gcstr.QuoteStrings(words)

						if lexer.matcher.HasSkipped() {
							words = append(words, "any other skipped character")
						}

						lexer.Err.SetSuggestion("Did you mean " + gcstr.OrString(words, false) + "?")
					}

					return lexer.get_tokens(), lexer.Err
				}

				if tmp != nil {
					lexer.tokens = append(lexer.tokens, tmp)
					lexer.skipped = 0
				}
			}
		}
	} else if has_matcher {
		for !lexer.IsExhausted() {
			tmp, err := lexer.lex_one(lexer)
			if err != nil {
				return lexer.get_tokens(), lexer.make_error(err)
			}

			if tmp != nil {
				lexer.tokens = append(lexer.tokens, tmp)
				lexer.skipped = 0
			}
		}
	} else {
		for !lexer.IsExhausted() {
			tmp, err := lexer.lex_one(lexer)
			if err != nil {
				lexer.Err = lexer.make_error(err)

				return lexer.get_tokens(), lexer.Err
			}

			if tmp != nil {
				lexer.tokens = append(lexer.tokens, tmp)
				lexer.skipped = 0
			}
		}
	}

	return lexer.get_tokens(), nil
}

// skip skips the characters of the lexer.
//
// Parameters:
//   - chars: The characters to skip.
func (lexer *Lexer[S]) skip(chars []rune) {
	for _, c := range chars {
		lexer.skipped += utf8.RuneLen(c)
	}
}

// AddToMatch is a method that adds a new match to the lexer.
//
// Parameters:
//   - symbol: The symbol of the match.
//   - word: The word of the match.
//
// Returns:
//   - error: An error if the word cannot be added to the lexer.
func (lexer *Lexer[S]) AddToMatch(symbol S, word string) error {
	err := lexer.matcher.AddToMatch(symbol, word)
	if err != nil {
		return err
	}

	return nil
}

// AddToSkipRule is a method that adds a new skip rule to the lexer.
//
// Parameters:
//   - words: The words of the skip rule.
//
// Returns:
//   - error: An error if the word cannot be added to the lexer.
func (lexer *Lexer[S]) AddToSkipRule(words ...string) error {
	err := lexer.matcher.AddToSkipRule(words...)
	if err != nil {
		return err
	}

	return nil
}
