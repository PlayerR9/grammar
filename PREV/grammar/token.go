package grammar

import (
	"iter"
	"strconv"
	"strings"

	gcers "github.com/PlayerR9/go-commons/errors"
	gcslc "github.com/PlayerR9/go-commons/slices"
	internal "github.com/PlayerR9/grammar/PREV/internal"
)

// Token is a token in the token stream.
type Token[T internal.TokenTyper] struct {
	// Parent is the parent of the token.
	Parent *Token[T]

	// FirstChild is the first child of the token.
	FirstChild *Token[T]

	// NextSibling is the next sibling of the token.
	NextSibling *Token[T]

	// LastChild is the last child of the token.
	LastChild *Token[T]

	// PrevSibling is the previous sibling of the token.
	PrevSibling *Token[T]

	// Type is the type of the token.
	Type T

	// Data is the data of the token.
	Data string

	// Lookahead is the lookahead token.
	Lookahead *Token[T]
}

func (t *Token[T]) Cleanup() []*Token[T] {
	panic("implement me")
}

func (t *Token[T]) IsSingleton() bool {
	panic("implement me")
}

func (t *Token[T]) LinkChildren(children []*Token[T]) {
	panic("implement me")
}

// String implements the pkg.Type interface.
func (t *Token[T]) String() string {
	var builder strings.Builder

	builder.WriteString("Token[T][")
	builder.WriteString(t.Type.String())

	if t.Data != "" {
		builder.WriteString(" (")
		builder.WriteString(strconv.Quote(t.Data))
		builder.WriteRune(')')
	}

	builder.WriteRune(']')

	return builder.String()
}

// DeepCopy returns a copy of the token.
//
// Returns:
//   - sdpkg.Type: The copy of the token. Never returns nil.
//
// However, pointers are not copied.
func (tk *Token[T]) Copy() *Token[T] {
	if tk == nil {
		return nil
	}

	return &Token[T]{
		Type: tk.Type,
		Data: tk.Data,
	}
}

// NewToken creates a new token with the given type and data.
//
// Parameters:
//   - type_: The type of the token.
//   - data: The data of the token.
//   - lookahead: The lookahead token.
//
// Returns:
//   - *Token[T]: A pointer to the new token. Never returns nil.
func NewToken[T internal.TokenTyper](type_ T, data string, lookahead *Token[T]) *Token[T] {
	return &Token[T]{
		Type:      type_,
		Data:      data,
		Lookahead: lookahead,
	}
}

// AddChildren is a convenience function to add multiple children to the node at once.
// It is more efficient than adding them one by one. Therefore, the behaviors are the
// same as the behaviors of the Token.AddChild function.
//
// Parameters:
//   - children: The children to add.
func (tk *Token[T]) AddChildren(children []*Token[T]) {
	if len(children) == 0 {
		return
	}

	children = gcslc.FilterNilValues(children)
	if len(children) == 0 {
		return
	}

	// Deal with the first child
	first_child := children[0]

	first_child.NextSibling = nil
	first_child.PrevSibling = nil

	last_child := tk.LastChild

	if last_child == nil {
		tk.FirstChild = first_child
	} else {
		last_child.NextSibling = first_child
		first_child.PrevSibling = last_child
	}

	first_child.Parent = tk
	tk.LastChild = first_child

	// Deal with the rest of the children
	for _, child := range children[1:] {
		child.NextSibling = nil
		child.PrevSibling = nil

		last_child := tk.LastChild
		last_child.NextSibling = child
		child.PrevSibling = last_child

		child.Parent = tk
		tk.LastChild = child
	}
}

// Children returns the children of the token.
//
// Returns:
//   - []*Token[T]: The children of the token. Never returns nil.
func (tk Token[T]) Children() []*Token[T] {
	if tk.FirstChild == nil {
		return nil
	}

	var children []*Token[T]

	for child := tk.FirstChild; child != nil; child = child.NextSibling {
		children = append(children, child)
	}

	return children
}

// Child returns an iterator over the children of the token.
//
// Returns:
//   - iter.Seq[*Token[T]]: An iterator over the children of the token.
func (tk *Token[T]) Child() iter.Seq[*Token[T]] {
	fn := func(yield func(*Token[T]) bool) {
		for c := tk.FirstChild; c != nil; c = c.NextSibling {
			if !yield(c) {
				return
			}
		}
	}

	return fn
}

// BackwardChild returns an iterator over the children of the token in reverse order.
//
// Returns:
//   - iter.Seq[*Token[T]]: An iterator over the children of the token in reverse order.
func (tk *Token[T]) BackwardChild() iter.Seq[*Token[T]] {
	fn := func(yield func(*Token[T]) bool) {
		for c := tk.LastChild; c != nil; c = c.PrevSibling {
			if !yield(c) {
				return
			}
		}
	}

	return fn
}

// IsLeaf checks if the token is a leaf.
//
// Returns:
//   - bool: True if the token is a leaf, false otherwise.
func (tk *Token[T]) IsLeaf() bool {
	return tk.FirstChild == nil
}

// CheckTokenAt checks if the token at the given index is of the given type.
//
// Parameters:
//   - tokens: The tokens to check.
//   - idx: The index of the token to check.
//   - type_: The type of the token to check.
//
// Returns:
//   - error: An error if the token at the given index is not of the given type.
//
// Errors:
//   - *errors.ErrInvalidParameter: If 'idx' is less than 0.
//   - *ErrUnexpectedToken: If the token at the given index is not of the given type.
func CheckTokenAt[T internal.TokenTyper](tokens []Token[T], idx int, type_ T) error {
	if idx < 0 {
		return gcers.NewErrInvalidParameter("idx", gcers.NewErrGTE(0))
	}

	var prev *T

	if idx >= len(tokens) {
		if idx > 0 && idx < len(tokens) {
			prev = &tokens[idx-1].Type
		}

		return NewErrUnexpectedToken(prev, nil, type_)
	} else if tokens[idx].Type != type_ {
		if idx > 0 {
			prev = &tokens[idx-1].Type
		}

		return NewErrUnexpectedToken(prev, &tokens[idx].Type, type_)
	}

	return nil
}

// CheckToken checks if the token is of the given type.
//
// Parameters:
//   - token: The token to check.
//   - prev: The previous token.
//   - type_: The type of the token to check.
//
// Returns:
//   - error: An error if the token is not of the given type.
//
// Errors:
//   - *ErrUnexpectedToken: If the token is not of the given type.
func CheckToken[T internal.TokenTyper](token, prev *Token[T], type_ T) error {
	if token == nil {
		var pt *T

		if prev != nil {
			pt = &prev.Type
		}

		return NewErrUnexpectedToken(pt, nil, type_)
	}

	if token.Type != type_ {
		var pt *T

		if prev != nil {
			pt = &prev.Type
		}

		return NewErrUnexpectedToken(pt, &token.Type, type_)
	}

	return nil
}
