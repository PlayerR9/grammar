package lexer

import (
	"fmt"
	"unicode/utf8"

	gcch "github.com/PlayerR9/go-commons/runes"
	gr "github.com/PlayerR9/grammar/grammar"
)

// LexFunc is the function that lexers call to lex the input stream.
//
// Parameters:
//   - lexer: The lexer that is lexing the input stream. Assumed to be non-nil.
//
// Returns:
//   - *Token[T]: The token that was lexed.
//   - error: Any error that occurred during lexing.
//
// If the returned token is nil, then it is ignored as a 'skip' rule.
type LexFunc[T gr.Enumer] func(lexer *Lexer[T]) (*gr.Token[T], error)

// Builder is a lexer builder.
type Builder[T gr.Enumer] struct {
	// table is the table of rules.
	table map[rune]LexFunc[T]

	// def_fn is the default function to call for unrecognized tokens.
	// If it is nil, then it is ignored.
	def_fn LexFunc[T]
}

// Register registers a new rule.
//
// Parameters:
//   - first_char: The first character of the rule.
//   - fn: The function to call when the rule is matched.
//
// If fn is nil, then it is ignored.
func (b *Builder[T]) Register(first_char rune, fn LexFunc[T]) {
	if fn == nil {
		return
	}

	if b.table == nil {
		b.table = make(map[rune]LexFunc[T])
	}

	b.table[first_char] = fn
}

// RegisterLiteral registers a new literal rule.
//
// Parameters:
//   - type_: The type of the token.
//   - literal: The literal to match.
//
// Returns:
//   - error: Any error that occurred during registration.
//
// If literal is empty, then it is ignored.
func (b *Builder[T]) RegisterLiteral(type_ T, literal string) error {
	if literal == "" {
		return nil
	}

	chars, err := gcch.StringToUtf8(literal)
	if err != nil {
		return err
	}

	if b.table == nil {
		b.table = make(map[rune]LexFunc[T])
	}

	char := chars[0]

	if len(chars) == 1 {
		b.table[char] = func(lexer *Lexer[T]) (*gr.Token[T], error) {
			_, _ = lexer.NextRune()

			tk := gr.NewTerminalToken(type_, literal)
			return tk, nil
		}
	} else {
		b.table[char] = func(lexer *Lexer[T]) (*gr.Token[T], error) {
			_, _ = lexer.NextRune()

			for i := 1; i < len(chars); i++ {
				exp := chars[i]

				r, ok := lexer.NextRune()
				if !ok {
					return nil, fmt.Errorf("expected %q after %q, got nothing instead", exp, chars[i-1])
				} else if r != exp {
					return nil, fmt.Errorf("expected %q after %q, got %q instead", exp, chars[i-1], r)
				}
			}

			tk := gr.NewTerminalToken(type_, literal)
			return tk, nil
		}
	}

	return nil
}

// RegisterSkip registers a new 'skip' rule.
//
// Parameters:
//   - literal: The literal to match.
//
// Returns:
//   - error: Any error that occurred during registration.
//
// If literal is empty, then it is ignored.
func (b *Builder[T]) RegisterSkip(literal string) error {
	if literal == "" {
		return nil
	}

	var chars []rune

	for len(literal) > 0 {
		c, size := utf8.DecodeRuneInString(literal)
		literal = literal[size:]

		if c == utf8.RuneError {
			return fmt.Errorf("invalid literal %q", literal)
		}

		chars = append(chars, c)
	}

	if b.table == nil {
		b.table = make(map[rune]LexFunc[T])
	}

	char := chars[0]

	if len(chars) == 1 {
		b.table[char] = func(lexer *Lexer[T]) (*gr.Token[T], error) {
			_, _ = lexer.NextRune()
			return nil, nil
		}
	} else {
		b.table[char] = func(lexer *Lexer[T]) (*gr.Token[T], error) {
			_, _ = lexer.NextRune()

			for i := 1; i < len(chars); i++ {
				exp := chars[i]

				r, ok := lexer.NextRune()
				if !ok {
					return nil, fmt.Errorf("expected %q after %q, got nothing instead", exp, chars[i-1])
				} else if r != exp {
					return nil, fmt.Errorf("expected %q after %q, got %q instead", exp, chars[i-1], r)
				}
			}

			return nil, nil
		}
	}

	return nil
}

// RegisterDefault registers a new 'default' rule.
//
// Parameters:
//   - fn: The function to call when the rule is matched.
//
// If fn is nil, then the previous default function is cleared.
func (b *Builder[T]) RegisterDefault(fn LexFunc[T]) {
	b.def_fn = fn
}

// Build builds a new Lexer instance.
//
// Returns:
//   - *Lexer: The new Lexer instance.
func (b Builder[T]) Build() *Lexer[T] {
	table := make(map[rune]LexFunc[T], len(b.table))

	for k, v := range b.table {
		table[k] = v
	}

	fn := b.def_fn

	return &Lexer[T]{
		table:  table,
		def_fn: fn,
	}
}

// Reset resets the builder.
func (b *Builder[T]) Reset() {
	for k := range b.table {
		b.table[k] = nil
		delete(b.table, k)
	}

	b.table = nil
	b.def_fn = nil
}
