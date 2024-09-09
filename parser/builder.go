package parser

import gr "github.com/PlayerR9/grammar/grammar"

// Builder is a parser builder.
type Builder[T gr.Enumer] struct {
	// table is the table of rules.
	table map[T]ParseFunc[T]
}

// NewBuilder creates a new parser builder.
//
// Returns:
//   - *Builder: The new parser builder. Never returns nil.
func NewBuilder[T gr.Enumer]() *Builder[T] {
	return &Builder[T]{
		table: make(map[T]ParseFunc[T]),
	}
}

// Register registers a rule.
//
// Parameters:
//   - type_: The type of the rule.
//   - fn: The parse function of the rule.
//
// If fn is nil, the rule will not be registered.
// Previously registered rules with the same type will be overwritten.
func (b *Builder[T]) Register(type_ T, fn ParseFunc[T]) {
	if fn == nil {
		return
	}

	b.table[type_] = fn
}

// Build builds a parser.
//
// Returns:
//   - *Parser: The new parser. Never returns nil.
func (b Builder[T]) Build() *Parser[T] {
	table := make(map[T]ParseFunc[T], len(b.table))

	for k, v := range b.table {
		table[k] = v
	}

	return &Parser[T]{
		table: table,
	}
}

// Reset resets the builder.
func (b *Builder[T]) Reset() {
	for k := range b.table {
		b.table[k] = nil
		delete(b.table, k)
	}

	b.table = nil
}
