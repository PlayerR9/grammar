package matcher

// RuleTyper is a rule type.
type RuleTyper interface {
	~int

	// String returns the name of the rule type.
	//
	// Returns:
	//   - string: The name of the rule type.
	String() string
}

// MatchRule is a rule to match.
type MatchRule[T RuleTyper] struct {
	// symbol is the symbol to match.
	symbol T

	// chars are the characters to match.
	chars []rune

	// should_skip is true if the rule should be skipped.
	should_skip bool
}

// CharAt returns the character at the given index.
//
// Returns:
//   - rune: The character at the given index.
//   - bool: True if the index is valid, false otherwise.
func (r MatchRule[T]) CharAt(at int) (rune, bool) {
	if at < 0 || at >= len(r.chars) {
		return 0, false
	}

	return r.chars[at], true
}
