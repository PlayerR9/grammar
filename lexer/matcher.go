package lexer

import (
	"fmt"
	"io"
	"slices"
	"strings"
	"unicode/utf8"

	gcers "github.com/PlayerR9/go-commons/errors"
	gcch "github.com/PlayerR9/go-commons/runes"
	gr "github.com/PlayerR9/grammar/grammar"
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
//   - symbol: The symbol to match.
//
// Returns:
//   - int: The index of the rule to match.
//   - bool: True if the rule to match is found, false otherwise.
func (m *Matcher[S]) find_index(symbol S) (int, bool) {
	cmp_func := func(e *MatchRule[S], symbol S) int {
		return int(symbol) - int(e.symbol)
	}

	return slices.BinarySearchFunc(m.rules, symbol, cmp_func)
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

	pos, ok := m.find_index(symbol)
	if !ok {
		m.rules = slices.Insert(m.rules, pos, rule)
	} else {
		m.rules[pos] = rule
	}

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

	var top int

	for i := 0; i < len(m.indices); i++ {
		rule := m.rules[m.indices[i]]

		c, ok := rule.CharAt(m.at)
		if !ok {
			continue
		}

		if c == char {
			m.indices[top] = m.indices[i]
			top++
		}
	}

	if top == 0 {
		// No valid matches exist.
		_ = scanner.UnreadRune()

		return true, nil
	}

	m.indices = m.indices[:top]
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

	rule := m.rules[m.indices[0]]

	if len(rule.chars) != m.at {
		return Matched[S]{
			symbol: nil,
			chars:  m.chars,
		}, m.make_error()
	}

	return Matched[S]{
		symbol: &rule.symbol,
		chars:  m.chars,
	}, nil
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
