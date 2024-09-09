package lexer

import (
	dbg "github.com/PlayerR9/go-debug/assert"
	gr "github.com/PlayerR9/grammar/PREV/grammar"
	internal "github.com/PlayerR9/grammar/PREV/internal"
)

// build_rule is the rule of the lexer.
type build_rule[T internal.TokenTyper] struct {
	// type_ is the type of the token.
	type_ T

	// fn is the function that is called when the rule is applied.
	fn LexFunc[T]

	// is_skip is true if the rule is a skip rule. False otherwise.
	is_skip bool
}

// new_build_rule creates a new build rule.
//
// Parameters:
//   - type_: The type of the token.
//   - is_skip: True if the rule is a skip rule. False otherwise.
//   - fn: The function that is called when the rule is applied.
//
// Returns:
//   - *build_rule[T]: The new build rule. Never returns nil.
func new_build_rule[T internal.TokenTyper](type_ T, is_skip bool, fn LexFunc[T]) *build_rule[T] {
	dbg.AssertNotNil(fn, "fn")

	return &build_rule[T]{
		type_:   type_,
		fn:      fn,
		is_skip: is_skip,
	}
}

// new_skip_build_rule creates a new build rule that is a skip rule.
//
// Parameters:
//   - fn: The function that is called when the rule is applied.
//
// Returns:
//   - *build_rule[T]: The new build rule. Never returns nil.
func new_skip_build_rule[T internal.TokenTyper](fn LexFunc[T]) *build_rule[T] {
	dbg.AssertNotNil(fn, "fn")

	return &build_rule[T]{
		type_:   T(0),
		fn:      fn,
		is_skip: true,
	}
}

// apply applies the build rule to the lexer.
//
// Parameters:
//   - lexer: The lexer.
//
// Returns:
//   - *gr.Token[T]: The new token if the rule is not a skip rule. Nil otherwise.
//   - error: An error if the rule failed.
func (r build_rule[T]) apply(lexer *ActiveLexer[T]) (*gr.Token[T], error) {
	dbg.AssertNotNil(lexer, "lexer")

	str, err := r.fn(lexer)
	if err != nil {
		return nil, err
	}

	if r.is_skip {
		return nil, nil
	}

	tk := gr.NewToken(r.type_, str, nil)

	return tk, nil
}
