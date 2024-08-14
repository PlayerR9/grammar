package matcher

// Matched is the matched result.
type Matched[T RuleTyper] struct {
	// symbol is the matched symbol.
	symbol *T

	// chars are the matched characters.
	chars []rune

	// should_skip is true if the rule should be skipped.
	should_skip bool
}

// new_err_matched creates a new matched with an error.
//
// Parameters:
//   - chars: The matched characters.
//   - should_skip: True if the rule should be skipped.
//
// Returns:
//   - Matched: The new matched with an error.
func new_err_matched[T RuleTyper](chars []rune, should_skip bool) Matched[T] {
	return Matched[T]{
		symbol:      nil,
		chars:       chars,
		should_skip: should_skip,
	}
}

// new_matched creates a new matched.
//
// Parameters:
//   - symbol: The matched symbol.
//   - chars: The matched characters.
//   - should_skip: True if the rule should be skipped.
//
// Returns:
//   - Matched: The new matched.
func new_matched[T RuleTyper](symbol T, chars []rune, should_skip bool) Matched[T] {
	return Matched[T]{
		symbol:      &symbol,
		chars:       chars,
		should_skip: should_skip,
	}
}

// GetMatch returns the matched symbol and the matched characters.
//
// Returns:
//   - T: The matched symbol.
//   - string: The matched characters.
func (m Matched[T]) GetMatch() (T, string) {
	if m.symbol == nil {
		return T(0), ""
	}

	symbol := *m.symbol

	return symbol, string(m.chars)
}

// GetChars returns the matched characters.
//
// Returns:
//   - []rune: The matched characters.
func (m Matched[T]) GetChars() []rune {
	return m.chars
}

// IsValidMatch returns true if the matched symbol is not nil.
//
// Returns:
//   - bool: True if the matched symbol is not nil, false otherwise.
func (m Matched[T]) IsValidMatch() bool {
	return m.symbol != nil
}

// IsShouldSkip returns true if the rule should be skipped.
//
// Returns:
//   - bool: True if the rule should be skipped, false otherwise.
func (m Matched[T]) IsShouldSkip() bool {
	return m.should_skip
}
