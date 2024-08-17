package matcher

import (
	"errors"
	"io"
	"slices"
	"unicode/utf8"

	gcers "github.com/PlayerR9/go-commons/errors"
	gcch "github.com/PlayerR9/go-commons/runes"

	gcslc "github.com/PlayerR9/go-commons/slices"
	gcstr "github.com/PlayerR9/go-commons/strings"
)

var (
	// NoMatch is the error that occurs when the matcher does not match any rule. Readers
	// must return this error as is and not wrap it as callers are expected to check for
	// this error using ==.
	NoMatch error
)

func init() {
	NoMatch = errors.New("no match")
}

// Matcher is the matcher of the grammar.
type Matcher[T RuleTyper] struct {
	// rules are the rules to match.
	rules []MatchRule[T]

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

	// matches are the matches of the matcher.
	matches []Matched[T]
}

// GetWords returns the words of the matcher.
//
// Returns:
//   - []string: The words of the matcher.
func (m Matcher[T]) GetWords() []string {
	var words []string

	for _, rule := range m.rules {
		words = append(words, string(rule.chars))
	}

	return words
}

// GetRuleNames returns the names of the rules of the matcher.
//
// Returns:
//   - []string: The names of the rules of the matcher.
func (m Matcher[T]) GetRuleNames() []string {
	var names []string

	for _, rule := range m.rules {
		word := rule.symbol.String()

		pos, ok := slices.BinarySearch(names, word)
		if !ok {
			names = slices.Insert(names, pos, word)
		}
	}

	return names
}

// HasSkipped checks whether the matcher has skipped any characters.
//
// Returns:
//   - bool: True if the matcher has skipped any characters, false otherwise.
func (m Matcher[T]) HasSkipped() bool {
	for _, rule := range m.rules {
		if rule.should_skip {
			return true
		}
	}

	return false
}

// IsEmpty checks whether the matcher has at least one rule.
//
// Returns:
//   - bool: True if matcher is empty, false otherwise.
func (m Matcher[T]) IsEmpty() bool {
	return len(m.rules) == 0
}

// find_index finds the index of the rule to match.
//
// Parameters:
//   - chars: The characters to match.
//
// Returns:
//   - int: The index of the rule to match. -1 if the rule to match is not found.
func (m Matcher[T]) find_index(chars []rune) int {
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
func (m *Matcher[T]) AddToMatch(symbol T, word string) error {
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

	rule := MatchRule[T]{
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

// AddToSkipRule adds a rule to skip.
//
// Parameters:
//   - words: The words to skip.
//
// Returns:
//   - error: An error if the rule to skip is invalid.
func (m *Matcher[T]) AddToSkipRule(words ...string) error {
	words = gcstr.FilterNonEmpty(words)
	if len(words) == 0 {
		return nil
	}

	for _, word := range words {
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

		rule := MatchRule[T]{
			symbol:      T(0),
			chars:       chars,
			should_skip: true,
		}

		idx := m.find_index(chars)
		if idx == -1 {
			m.rules = append(m.rules, rule)
		} else {
			m.rules[idx] = rule
		}
	}

	return nil
}

// match_first matches the first character of the matcher.
//
// Parameters:
//   - scanner: The scanner to match.
//
// Returns:
//   - bool: True if the error is not critical, false otherwise.
//   - error: An error if the first character does not match.
func (m *Matcher[T]) match_first(scanner io.RuneScanner) (bool, error) {
	m.indices = m.indices[:0]
	m.prev = nil
	m.got = nil
	m.at = 0
	m.chars = m.chars[:0]
	m.matches = m.matches[:0]

	char, _, err := scanner.ReadRune()
	if err == io.EOF {
		return true, nil
	} else if err != nil {
		return false, err
	}

	m.got = &char

	for i, rule := range m.rules {
		c, _ := rule.CharAt(m.at)

		if char == c {
			m.indices = append(m.indices, i)
		}
	}

	if len(m.indices) == 0 {
		err := scanner.UnreadRune()
		if err != nil {
			return false, err
		}

		return true, NoMatch
	}

	m.prev = &char
	m.at++
	m.chars = append(m.chars, char)

	return true, nil
}

// filter filters the rules to match.
//
// Parameters:
//   - scanner: The scanner to filter.
//
// Returns:
//   - bool: True if the scanner is exhausted, false otherwise.
//   - error: An error if the scanner is exhausted or invalid.
func (m *Matcher[T]) filter(scanner io.RuneScanner) (bool, error) {
	char, _, err := scanner.ReadRune()
	if err == io.EOF {
		return true, nil
	} else if err != nil {
		return false, err
	}

	m.got = &char

	fn := func(idx int) bool {
		rule := m.rules[idx]

		c, ok := rule.CharAt(m.at)
		if ok && c == char {
			return true
		}

		if !ok {
			tmp := new_matched(rule.symbol, m.chars, rule.should_skip)
			m.matches = append(m.matches, tmp)
		}

		return false
	}

	tmp, ok := gcslc.SFSeparateEarly(m.indices, fn)
	if !ok {
		err := scanner.UnreadRune()
		if err != nil {
			return false, err
		}

		return true, nil
	}

	m.indices = tmp
	m.chars = append(m.chars, char)
	m.prev = &char
	m.at++

	return false, nil
}

// make_error makes an error.
//
// Returns:
//   - error: An error if the next characters do not match.
func (m Matcher[T]) make_error() error {
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

	return gcstr.NewErrUnexpectedRune(m.prev, m.got, chars...)
}

// Match matches the next characters of the matcher.
//
// Parameters:
//   - scanner: The scanner to match.
//
// Returns:
//   - bool: True if the error is not critical, false otherwise.
//   - error: An error if the next characters do not match.
//
// A non-critical error is an error that occurs when the matcher cannot match a word
// due to it not being in the dictionary. Because of that, they can be ignored.
//
// However, critical errors are errors that are external to the dictionary and prevent
// the matching to continue.
func (m *Matcher[T]) Match(scanner io.RuneScanner) (bool, error) {
	if scanner == nil {
		return false, gcers.NewErrNilParameter("scanner")
	}

	not_critical, err := m.match_first(scanner)
	if err != nil {
		if not_critical {
			return true, m.make_error()
		}

		return false, err
	}

	var is_done bool

	for !is_done && len(m.indices) > 0 {
		is_done, err = m.filter(scanner)
		if err != nil {
			return false, err
		}
	}

	if len(m.matches) == 0 {
		return true, m.make_error()
	}

	return true, nil
}

// GetMatches returns the matches of the matcher.
//
// Returns:
//   - []Matched[T]: The matches of the matcher. Nil if no matches were found.
func (m Matcher[T]) GetMatches() []Matched[T] {
	if len(m.matches) == 0 {
		return nil
	}

	matches := make([]Matched[T], 0, len(m.matches))
	copy(matches, m.matches)

	return matches
}
