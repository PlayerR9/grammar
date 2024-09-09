package grammar

import (
	gcers "github.com/PlayerR9/go-commons/errors"
)

// Enumer is an interface for all token types. The 0th value is reserved for the EOF token.
type Enumer interface {
	~int

	// String returns the literal name of the token type.
	// This is used for debugging and error messages.
	//
	// Returns:
	// 	- string: The literal name of the token type.
	String() string
}

// Token represents a token in the grammar.
type Token[T Enumer] struct {
	// Type is the type of the token.
	Type T

	// Data is the value of the token.
	Data string

	// Pos is the position of the token in the input stream.
	Pos int

	// Lookahead is the next token in the input stream.
	Lookahead *Token[T]

	// Children are the children of the token.
	Children []*Token[T]
}

// NewTerminalToken creates a new terminal token with the given type, data, and lookahead.
//
// Parameters:
//   - type_: The type of the token.
//   - data: The value of the token.
//
// Returns:
//   - *Token: The new token. Never returns nil.
func NewTerminalToken[T Enumer](type_ T, data string) *Token[T] {
	return &Token[T]{
		Type:      type_,
		Data:      data,
		Lookahead: nil,
		Children:  nil,
	}
}

// NewToken creates a new non-terminal token with the given type, data, and children.
//
// Keep in mind that the last children must be the furthest in the input stream.
//
// Parameters:
//   - type_: The type of the token.
//   - data: The value of the token.
//   - children: The children of the token.
//
// Returns:
//   - *Token: The new token.
//   - error: An error of type *errors.ErrInvalidParameter if there is an empty list of children.
func NewToken[T Enumer](type_ T, data string, children []*Token[T]) (*Token[T], error) {
	if len(children) == 0 {
		return nil, gcers.NewErrInvalidParameter("children", gcers.NewErrEmpty(children))
	}

	return &Token[T]{
		Type:      type_,
		Data:      data,
		Lookahead: children[len(children)-1].Lookahead,
		Children:  children,
		Pos:       children[0].Pos,
	}, nil
}
