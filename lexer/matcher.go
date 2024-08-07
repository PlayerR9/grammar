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

type MatchRule[S gr.TokenTyper] struct {
	// symbol is the symbol to match.
	symbol S

	// chars are the characters to match.
	chars []rune
}

func (r *MatchRule[S]) CharAt(at int) (rune, bool) {
	if at < 0 || at >= len(r.chars) {
		return 0, false
	}

	return r.chars[at], true
}

type Matcher[S gr.TokenTyper] struct {
	rules   []*MatchRule[S]
	indices []int
	at      int
	prev    *rune
	got     *rune
}

func NewMatcher[S gr.TokenTyper]() *Matcher[S] {
	return new(Matcher[S])
}

func (m *Matcher[S]) find_index(symbol S) (int, bool) {
	cmp_func := func(e *MatchRule[S], symbol S) int {
		return int(symbol) - int(e.symbol)
	}

	return slices.BinarySearchFunc(m.rules, symbol, cmp_func)
}

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

func (m *Matcher[S]) filter(scanner io.RuneScanner) (bool, error) {
	if scanner == nil {
		return true, gcers.NewErrNilParameter("scanner")
	}

	char, _, err := scanner.ReadRune()
	if IsExhausted(err) {
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

	return false, nil
}

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

func (m *Matcher[S]) Match(scanner io.RuneScanner) (S, error) {
	if scanner == nil {
		return S(0), gcers.NewErrNilParameter("scanner")
	}

	c, _, err := scanner.ReadRune()
	if err != nil {
		return S(0), err
	}

	m.indices = m.indices[:0]
	m.prev = nil
	m.got = &c
	m.at = 0

	for i, rule := range m.rules {
		char, _ := rule.CharAt(m.at)

		if char == c {
			m.indices = append(m.indices, i)
		}
	}

	if len(m.indices) == 0 {
		return S(0), m.make_error()
	}

	m.prev = &c

	at := 1

	for {
		is_done, err := m.filter(scanner)
		if err != nil {
			return S(0), err
		}

		if is_done {
			break
		}

		at++
	}

	if len(m.indices) == 0 {
		return S(0), m.make_error()
	}

	if len(m.indices) > 1 {
		words := make([]string, 0, len(m.indices))

		for _, idx := range m.indices {
			rule := m.rules[idx]

			words = append(words, string(rule.chars))
		}

		return S(0), fmt.Errorf("ambiguous match: %s", strings.Join(words, ", "))
	}

	rule := m.rules[m.indices[0]]

	if len(rule.chars) != m.at {
		return S(0), m.make_error()
	}

	return rule.symbol, nil
}
