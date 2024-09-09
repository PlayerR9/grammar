package lexer

import (
	"fmt"
	"io"

	gr "github.com/PlayerR9/grammar/PREV/grammar"

	internal "github.com/PlayerR9/grammar/PREV/internal"
)

// Builder is the builder of the lexer.
type Builder[T internal.TokenTyper] struct {
	// table is the table of the lexer.
	table map[rune]*build_rule[T]

	// def_case is the default case of the lexer.
	def_case LexOnceFunc[T]
}

// Register registers a new rule for the given character.
//
// Parameters:
//   - char: The character to register the rule for.
//   - type_: The type of the token.
//   - fn: The function to call when the character is encountered.
//
// If fn is nil, Register does nothing.
// When multiple rules are registered for the same character, the last rule is used.
func (b *Builder[T]) Register(char rune, type_ T, fn LexFunc[T]) {
	if fn == nil {
		return
	}

	if b.table == nil {
		b.table = make(map[rune]*build_rule[T])
	}

	rule := new_build_rule(type_, false, fn)

	b.table[char] = rule
}

// RegisterSkip registers a new rule for the given character as a skip rule.
//
// Parameters:
//   - char: The character to register the rule for.
//   - fn: The function to call when the character is encountered.
//
// If fn is nil, RegisterSkip does nothing.
// When multiple rules are registered for the same character, the last rule is used.
func (b *Builder[T]) RegisterSkip(char rune, fn LexFunc[T]) {
	if fn == nil {
		return
	}

	if b.table == nil {
		b.table = make(map[rune]*build_rule[T])
	}

	rule := new_skip_build_rule(fn)

	b.table[char] = rule
}

// SetDefaultCase sets the default case of the lexer.
//
// Parameters:
//   - fn: The default case of the lexer.
//
// If fn is nil, the default case is removed.
func (b *Builder[T]) SetDefaultCase(fn LexOnceFunc[T]) {
	b.def_case = fn
}

// Build builds the lexer.
//
// Returns:
//   - *Lexer[T]: The lexer. Never returns nil.
func (b Builder[T]) Build() *Lexer[T] {
	var fn LexOnceFunc[T]

	if b.table == nil {
		fn = func(lexer *ActiveLexer[T]) ([]*gr.Token[T], error) {
			// dbg.AssertNotNil(lexer, "l")

			char, ok := lexer.PeekRune()
			if !ok {
				return nil, io.EOF
			}

			return nil, fmt.Errorf("unknown character: %q", string(char))
		}
	} else if b.def_case == nil {
		fn = func(lexer *ActiveLexer[T]) ([]*gr.Token[T], error) {
			// dbg.AssertNotNil(lexer, "l")

			char, ok := lexer.PeekRune()
			if !ok {
				return nil, io.EOF
			}

			rule, ok := b.table[char]
			if !ok {
				return nil, fmt.Errorf("unknown character: %q", string(char))
			}

			tk, err := rule.apply(lexer)
			if err != nil {
				return nil, err
			}

			return []*gr.Token[T]{tk}, nil
		}
	} else {
		def_case := b.def_case

		fn = func(lexer *ActiveLexer[T]) ([]*gr.Token[T], error) {
			// dbg.AssertNotNil(lexer, "l")

			char, ok := lexer.PeekRune()
			if !ok {
				return nil, io.EOF
			}

			rule, ok := b.table[char]
			if !ok {
				tk, err := def_case(lexer)
				if err != nil {
					return nil, err
				}

				return tk, nil
			}

			tk, err := rule.apply(lexer)
			if err != nil {
				return nil, err
			}

			return []*gr.Token[T]{tk}, nil
		}
	}

	return &Lexer[T]{
		fn: fn,
	}
}

// Reset resets the lexer builder.
//
// This function resets the table and the default case function of the lexer builder.
func (b *Builder[T]) Reset() {
	b.table = nil
	b.def_case = nil
}
