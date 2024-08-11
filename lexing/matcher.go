package lexing

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"unicode/utf8"

	gcers "github.com/PlayerR9/go-commons/errors"
	gcch "github.com/PlayerR9/go-commons/runes"
	gr "github.com/PlayerR9/grammar/grammar"

	gcslc "github.com/PlayerR9/go-commons/slices"
)

// MatchRule is a rule to match.
type MatchRule[S gr.TokenTyper] struct {
	// symbol is the symbol to match.
	symbol S

	// chars are the characters to match.
	chars []rune
}

// CharAt returns the character at the given index.
//
// Returns:
//   - rune: The character at the given index.
//   - bool: True if the index is valid, false otherwise.
func (r *MatchRule[S]) CharAt(at int) (rune, bool) {
	if at < 0 || at >= len(r.chars) {
		return 0, false
	}

	return r.chars[at], true
}

// Matcher is the matcher of the grammar.
type Matcher[S gr.TokenTyper] struct {
	// rules are the rules to match.
	rules []*MatchRule[S]

	// indices are the indices of the rules to match.
	indices []int

	// at is the position of the matcher in the input stream.
	at int

	// prev is the previous rune of the matcher.
	prev *rune

	// prev_size is the size of the previous rune of the matcher.
	got *rune

	// chars are the characters extracted from the input stream.
	chars []rune
}

// NewMatcher creates a new matcher.
//
// Returns:
//   - *Matcher: The new matcher. Never returns nil.
func NewMatcher[S gr.TokenTyper]() *Matcher[S] {
	return new(Matcher[S])
}

// find_index finds the index of the rule to match.
//
// Parameters:
//   - chars: The characters to match.
//
// Returns:
//   - int: The index of the rule to match. -1 if the rule to match is not found.
func (m *Matcher[S]) find_index(chars []rune) int {
	for i, rule := range m.rules {
		if len(rule.chars) != len(chars) {
			continue
		}

		if slices.Equal(rule.chars, chars) {
			return i
		}
	}

	return -1
}

// AddToMatch adds a rule to match.
//
// Parameters:
//   - symbol: The symbol to match.
//   - word: The word to match.
//
// Returns:
//   - error: An error if the rule to match is invalid.
func (m *Matcher[S]) AddToMatch(symbol S, word string) error {
	if word == "" {
		return nil
	}

	var chars []rune

	for at := 0; len(word) > 0; at++ {
		c, size := utf8.DecodeRuneInString(word)
		if c == utf8.RuneError {
			return gcch.NewErrInvalidUTF8Encoding(at)
		}

		chars = append(chars, c)
		at += size
		word = word[size:]
	}

	rule := &MatchRule[S]{
		symbol: symbol,
		chars:  chars,
	}

	idx := m.find_index(chars)
	if idx == -1 {
		m.rules = append(m.rules, rule)
	} else {
		m.rules[idx] = rule
	}

	return nil
}

func (m *Matcher[S]) match_first(scanner io.RuneScanner) error {
	c, _, err := scanner.ReadRune()
	if err != nil {
		return err
	}

	m.indices = m.indices[:0]
	m.prev = nil
	m.got = &c
	m.at = 0
	m.chars = m.chars[:0]

	for i, rule := range m.rules {
		char, _ := rule.CharAt(m.at)

		if char == c {
			m.indices = append(m.indices, i)
		}
	}

	if len(m.indices) == 0 {
		_ = scanner.UnreadRune()

		return m.make_error()
	}

	m.prev = &c
	m.at++

	m.chars = append(m.chars, c)

	return nil
}

// filter filters the rules to match.
//
// Parameters:
//   - scanner: The scanner to filter.
//
// Returns:
//   - bool: True if the scanner is exhausted, false otherwise.
//   - error: An error if the scanner is exhausted or invalid.
func (m *Matcher[S]) filter(scanner io.RuneScanner) (bool, error) {
	if scanner == nil {
		return true, gcers.NewErrNilParameter("scanner")
	}

	char, _, err := scanner.ReadRune()
	if err == io.EOF {
		return true, nil
	} else if err != nil {
		return false, err
	}

	m.got = &char

	f := func(idx int) bool {
		rule := m.rules[idx]

		c, ok := rule.CharAt(m.at)
		return ok && c == char
	}

	tmp, ok := gcslc.SFSeparateEarly(m.indices, f)
	if !ok {
		// No valid matches exist.
		_ = scanner.UnreadRune()

		tmp, ok := m.filter_size(m.indices)
		if ok {
			m.indices = tmp
		}

		return true, nil
	}

	m.indices = tmp

	m.prev = &char
	m.at++
	m.chars = append(m.chars, char)

	return false, nil
}

// make_error makes an error.
//
// Returns:
//   - error: An error if the next characters do not match.
func (m *Matcher[S]) make_error() error {
	var chars []rune

	for _, rule := range m.rules {
		c, ok := rule.CharAt(m.at)
		if !ok {
			continue
		}

		pos, ok := slices.BinarySearch(chars, c)
		if !ok {
			chars = slices.Insert(chars, pos, c)
		}
	}

	return NewErrUnexpectedRune(m.prev, m.got, chars...)
}

// Match matches the next characters of the matcher.
//
// Parameters:
//   - scanner: The scanner to match.
//
// Returns:
//   - S: The matched symbol.
//   - error: An error if the next characters do not match.
func (m *Matcher[S]) Match(scanner io.RuneScanner) (Matched[S], error) {
	if scanner == nil {
		return Matched[S]{
			symbol: nil,
			chars:  m.chars,
		}, gcers.NewErrNilParameter("scanner")
	}

	err := m.match_first(scanner)
	if err != nil {
		return Matched[S]{
			symbol: nil,
			chars:  m.chars,
		}, err
	}

	for {
		is_done, err := m.filter(scanner)
		if err != nil {
			return Matched[S]{
				symbol: nil,
				chars:  m.chars,
			}, err
		}

		if is_done {
			break
		}
	}

	if len(m.indices) == 0 {
		return Matched[S]{
			symbol: nil,
			chars:  m.chars,
		}, m.make_error()
	}

	if len(m.indices) > 1 {
		words := make([]string, 0, len(m.indices))

		for _, idx := range m.indices {
			rule := m.rules[idx]

			words = append(words, string(rule.chars))
		}

		return Matched[S]{
			symbol: nil,
			chars:  m.chars,
		}, fmt.Errorf("ambiguous match: %s", strings.Join(words, ", "))
	}

	tmp, ok := m.filter_size(m.indices)
	if !ok {
		return Matched[S]{
			symbol: nil,
			chars:  m.chars,
		}, m.make_error()
	}

	m.indices = tmp

	rule := m.rules[m.indices[0]]

	return Matched[S]{
		symbol: &rule.symbol,
		chars:  m.chars,
	}, nil
}

func (m *Matcher[S]) filter_size(indices []int) ([]int, bool) {
	f := func(idx int) bool {
		rule := m.rules[idx]
		return len(rule.chars) == m.at
	}

	return gcslc.SFSeparateEarly(indices, f)
}

type Matched[S gr.TokenTyper] struct {
	symbol *S
	chars  []rune
}

func (m *Matched[S]) GetMatch() (S, string) {
	if m.symbol == nil {
		return S(0), ""
	}

	symbol := *m.symbol

	return symbol, string(m.chars)
}

func (m *Matched[S]) GetChars() []rune {
	return m.chars
}

func (m *Matched[S]) IsValidMatch() bool {
	return m.symbol != nil
}

var (
	// Done is the error returned when the lexing process is done without error and before reaching the end of the
	// stream. Readers must return Done itself and not wrap it as callers will test this error using ==.
	Done error
)

func init() {
	Done = errors.New("done")
}

// LexFunc is a function that can be used with RightLex.
//
// Parameters:
//   - scanner: The rune scanner.
//
// Returns:
//   - []rune: The list of runes.
//   - error: The error.
//
// This function must assume scanner is never nil. Moreover, use io.EOF to signify the end of the stream.
// Lastly, the error Done is returned if the lexing process is done before reaching the end of the stream.
//
// Finally, it is suggested to always push back the last rune read if any error that is not io.EOF is returned.
type LexFunc func(scanner io.RuneScanner) ([]rune, error)

// RightLex reads the content of the stream and returns the list of runes according to the given function.
//
// Parameters:
//   - scanner: The rune scanner.
//   - lex_f: The lexing function.
//
// Returns:
//   - []rune: The list of runes.
//   - error: The error.
//
// Errors:
//   - *errors.ErrInvalidParameter: When scanner or lex_f is nil.
//   - any other error returned by lex_f.
func RightLex(scanner io.RuneScanner, lex_f LexFunc) ([]rune, error) {
	if scanner == nil {
		return nil, gcers.NewErrNilParameter("scanner")
	} else if lex_f == nil {
		return nil, gcers.NewErrNilParameter("lex_f")
	}

	var chars []rune

	for {
		curr, err := lex_f(scanner)
		if err == io.EOF || err == Done {
			chars = append(chars, curr...)

			break
		}

		if err != nil {
			return nil, err
		}

		chars = append(chars, curr...)
	}

	return chars, nil
}

func (m *Matcher[S]) GetWords() []string {
	var words []string

	for _, rule := range m.rules {
		words = append(words, string(rule.chars))
	}

	return words
}
