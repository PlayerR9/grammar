package grammar

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/PlayerR9/lib_units/common"
)

// TokenIterator is a pull-based iterator that iterates
// over the children of a Token.
type TokenIterator[S TokenTyper] struct {
	parent, current *Token[S]
}

// Consume implements the common.Iterater interface.
//
// The only error type that can be returned by this function is the *common.ErrExhaustedIter type.
//
// Moreover, the return value is always of type *Token[S] and never nil; unless the iterator
// has reached the end of the branch.
func (iter *TokenIterator[S]) Consume() (*Token[S], error) {
	if iter.current == nil {
		return nil, common.NewErrExhaustedIter()
	}

	node := iter.current
	iter.current = iter.current.NextSibling

	return node, nil
}

// Restart implements the common.Iterater interface.
func (iter *TokenIterator[S]) Restart() {
	iter.current = iter.parent.FirstChild
}

// Token is a node in a tree.
type Token[S TokenTyper] struct {
	Parent, FirstChild, NextSibling, LastChild, PrevSibling *Token[S]

	Type      S
	Data      string
	At        int
	Lookahead *Token[S]
}

// IsLeaf is a method that checks whether the node is a leaf.
//
// Returns:
//   - bool: True if the node is a leaf, false otherwise.
func (tk *Token[S]) IsLeaf() bool {
	return tk.FirstChild == nil
}

// Iterator implements the common.Iterable interface.
//
// This function returns an iterator that iterates over the direct children of the node.
// Implemented as a pull-based iterator, this function never returns nil and any of the
// values is guaranteed to be a non-nil node of type Token[S].
func (tk *Token[S]) Iterator() common.Iterater[*Token[S]] {
	return &TokenIterator[S]{
		parent:  tk,
		current: tk.FirstChild,
	}
}

// String implements the fmt.Stringer interface.
func (tk *Token[S]) String() string {
	var builder strings.Builder

	builder.WriteString("Token[")

	builder.WriteString(tk.Type.String())

	if tk.Data != "" {
		builder.WriteString(" (")
		builder.WriteString(strconv.Quote(tk.Data))
		builder.WriteRune(')')
	}

	builder.WriteString(" : ")
	builder.WriteString(strconv.Itoa(tk.At))
	builder.WriteRune(']')

	return builder.String()
}

// Size returns the number of runes in the token's data.
//
// Returns:
//   - int: The number of runes in the token's data.
func (t *Token[S]) Size() int {
	if t.Data != "" {
		return utf8.RuneCountInString(t.Data)
	}

	var size int

	for c := t.FirstChild; c != nil; c = c.NextSibling {
		size += c.Size()
	}

	return size
}

// SetLookahead sets the lookahead token.
//
// Parameters:
//   - lookahead: The lookahead token.
func (t *Token[S]) SetLookahead(lookahead *Token[S]) {
	if lookahead == nil {
		t.Lookahead = nil

		return
	}

	t.Lookahead = common.AssertConv[*Token[S]](lookahead, "lookahead")
}

// GetLookahead returns the lookahead token.
//
// Returns:
//   - *Token[S]: The lookahead token.
func (t *Token[S]) GetLookahead() *Token[S] {
	return t.Lookahead
}

// GetType returns the type of the token.
//
// Returns:
//   - TokenTyper: The type of the token.
func (t *Token[S]) GetType() S {
	return t.Type
}

// NewToken creates a new node with the given data.
//
// Parameters:
//   - t_type: The type of the node.
//   - data: The data of the node.
//   - at: The position of the node in the source code.
//   - lookahead: The lookahead of the node.
//
// Returns:
//   - *Token[S]: A pointer to the newly created node. It is
//     never nil.
func NewToken[S TokenTyper](t_type S, data string, at int, lookahead *Token[S]) *Token[S] {
	return &Token[S]{
		Type:      t_type,
		Data:      data,
		At:        at,
		Lookahead: lookahead,
	}
}

// AddChildren is a convenience function to add multiple children to the node at once.
// It is more efficient than adding them one by one. Therefore, the behaviors are the
// same as the behaviors of the Token.AddChild function.
//
// Parameters:
//   - children: The children to add.
func (tk *Token[S]) AddChildren(children []*Token[S]) {
	if len(children) == 0 {
		return
	}

	var top int

	for i := 0; i < len(children); i++ {
		child := children[i]

		if child != nil {
			children[top] = child
			top++
		}
	}

	children = children[:top]
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
	for i := 1; i < len(children); i++ {
		child := children[i]

		child.NextSibling = nil
		child.PrevSibling = nil

		last_child := tk.LastChild
		last_child.NextSibling = child
		child.PrevSibling = last_child

		child.Parent = tk
		tk.LastChild = child
	}
}
