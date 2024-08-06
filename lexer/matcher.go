package lexer

import (
	"slices"
	"strings"
	"unicode/utf8"

	gr "github.com/PlayerR9/grammar/grammar"

	gcch "github.com/PlayerR9/go-commons/runes"
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
}

func NewMatcher[S gr.TokenTyper]() *Matcher[S] {
	return &Matcher[S]{
		rules: make([]*MatchRule[S], 0),
	}
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

func (m *Matcher[S]) UniqueCharsAt(at int) []rune {
	var chars []rune

	for _, rule := range m.rules {
		c, ok := rule.CharAt(at)
		if !ok {
			continue
		}

		pos, ok := slices.BinarySearch(chars, c)
		if !ok {
			chars = slices.Insert(chars, pos, c)
		}
	}

	return chars
}

type matched_data struct {
	data   string
	weight int
	reason error
}

func (m *Matcher[S]) match(lexer *Lexer[S], chars []rune) (*matched_data, error) {
	if len(chars) == 0 {
		data := &matched_data{
			data:   "",
			weight: 0,
			reason: nil,
		}

		return data, nil
	}

	var prev *rune
	var builder strings.Builder

	for i, char := range chars {
		c, err := lexer.Next()
		if IsExhausted(err) {
			data := &matched_data{
				data:   builder.String(),
				weight: i + 1,
				reason: NewErrUnexpectedRune(prev, nil, char),
			}

			return data, NewErrUnexpectedRune(prev, nil, char)
		} else if err != nil {
			return nil, err
		}

		if c != char {
			data := &matched_data{
				data:   builder.String(),
				weight: i + 1,
				reason: NewErrUnexpectedRune(prev, &c, char),
			}

			return data, NewErrUnexpectedRune(prev, &c, char)
		}

		builder.WriteRune(c)
		prev = &c
	}

	data := &matched_data{
		data:   builder.String(),
		weight: len(chars),
		reason: nil,
	}

	return data, nil
}

func (m *Matcher[S]) filter(at int) bool {
	var top int

	for i := 0; i < len(m.indices); i++ {
		rule := m.rules[m.indices[i]]

		c, ok := rule.CharAt(at)
		if !ok {
			continue
		}

	}

	if top == 0 {
		return true
	} else {
		m.indices = m.indices[:top]
	}

	return false
}

func (m *Matcher[S]) Match(lexer *Lexer[S]) error {
	c, err := lexer.Next()
	if err != nil {
		return err
	}

	m.indices = m.indices[:0]

	for i, rule := range m.rules {
		char, _ := rule.CharAt(0)

		if char == c {
			m.indices = append(m.indices, i)
		}
	}

	switch len(m.indices) {
	case 0:
		expecteds := m.UniqueCharsAt(0)

		return NewErrUnexpectedRune(nil, nil, expecteds...)
	case 1:
	default:
	}
}
