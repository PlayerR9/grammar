package grammar

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	luc "github.com/PlayerR9/lib_units/common"
)

// TokenTyper is an interface that defines the behavior of a token type.
//
// Value of 0 is reserved for the EOF token.
type TokenTyper interface {
	~int

	fmt.Stringer
}

// Token is a struct that represents a generic token of type T.
type Token[T TokenTyper] struct {
	// Type is the type of the token.
	Type T

	// Data is the data of the token.
	// Can only be a string or []*Token[T].
	//
	// For nil values, an empty string should be used.
	Data any

	// Lookahead is the lookahead token.
	Lookahead *Token[T]

	// At is the position of the token in the input. It is the rune position of the token in
	// the input string.
	At int
}

// IsLeaf implements the Strings.Noder interface.
func (t *Token[T]) IsLeaf() bool {
	_, ok := t.Data.(string)
	return ok
}

// Iterator implements the Strings.Noder interface.
func (t *Token[T]) Iterator() luc.Iterater[Noder] {
	children, ok := t.Data.([]*Token[T])
	if !ok {
		return nil
	}

	nodes := make([]Noder, 0, len(children))
	for _, child := range children {
		nodes = append(nodes, child)
	}

	return luc.NewSimpleIterator(nodes)
}

// String implements the Strings.Noder interface.
//
// Format:
//
//	"Token[T][{{ .Type }} ({{ .Data }})] : {{ .At }}]"
func (t *Token[T]) String() string {
	var builder strings.Builder

	builder.WriteString("Token[T][")
	builder.WriteString(t.Type.String())

	data, ok := t.Data.(string)
	if ok && data != "" {
		builder.WriteString(" (")
		builder.WriteString(strconv.Quote(data))
		builder.WriteRune(')')
	}

	builder.WriteString(" : ")
	builder.WriteString(strconv.Itoa(t.At))
	builder.WriteRune(']')

	return builder.String()
}

// NewToken creates a new token of type T.
//
// Parameters:
//   - t: The type of the token.
//   - d: The data of the token.
//   - at: The position of the token in the input.
//   - lookahead: The lookahead token.
//
// Returns:
//   - *Token[T]: A pointer to the newly created token.
//   - error: An error of type *common.ErrInvalidParameter if the data is nil or not
//     of type string or []*Token[T].
func NewToken[T TokenTyper](t T, d any, at int, lookahead *Token[T]) (*Token[T], error) {
	if d == nil {
		return nil, luc.NewErrNilParameter("d")
	}

	switch d := d.(type) {
	case string, []*Token[T]:
		return &Token[T]{
			Type:      t,
			Data:      d,
			Lookahead: lookahead,
			At:        at,
		}, nil
	default:
		return nil, luc.NewErrInvalidParameter("d", fmt.Errorf("expected string or []*Token[T], got %T instead", d))
	}
}

// Size returns the number of runes in the token's data.
//
// Returns:
//   - int: The number of runes in the token's data.
func (t *Token[T]) Size() int {
	switch data := t.Data.(type) {
	case string:
		return utf8.RuneCountInString(data)
	case []*Token[T]:
		var size int
		for _, token := range data {
			size += token.Size()
		}
		return size
	default:
		return 0
	}
}
