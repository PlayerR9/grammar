package grammar

import (
	"strconv"
	"strings"
	"unicode/utf8"

	itr "github.com/PlayerR9/go-commons/iterator"
	gcslc "github.com/PlayerR9/go-commons/slices"
)

// TokenIterator is a pull-based iterator that iterates over the children of a Node.
type TokenIterator[T TokenTyper] struct {
	first, current *Token[T]
}

// Consume implements the ast.Iterater interface.
func (iter *TokenIterator[T]) Apply(fn itr.IteratorFunc) error {
	if iter.current == nil {
		return itr.ErrExausted
	}

	err := fn(iter.current)
	if err != nil {
		return err
	}

	iter.current = iter.current.NextSibling

	return nil
}

// Restart implements the ast.Iterater interface.
func (iter *TokenIterator[T]) Reset() {
	iter.current = iter.first
}

// TokenReverseIterator is a pull-based iterator that iterates over the children of a Node in reverse order.
type TokenReverseIterator[T TokenTyper] struct {
	last, current *Token[T]
}

// Consume implements the Iterater interface.
func (iter *TokenReverseIterator[T]) Apply(fn itr.IteratorFunc) error {
	if iter.current == nil {
		return itr.ErrExausted
	}

	err := fn(iter.current)
	if err != nil {
		return err
	}

	iter.current = iter.current.PrevSibling

	return nil
}

// Restart implements the ast.Iterater interface.
func (iter *TokenReverseIterator[T]) Reset() {
	iter.current = iter.last
}

// Token is a node in a tree.
type Token[S TokenTyper] struct {
	Parent, FirstChild, NextSibling, LastChild, PrevSibling *Token[S]

	Type      S
	Data      string
	At        int
	Lookahead *Token[S]
}

// String implements the fmt.Stringer interface.
func (tk Token[S]) String() string {
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
func (t Token[S]) Size() int {
	if t.Data != "" {
		return utf8.RuneCountInString(t.Data)
	}

	var size int

	for c := t.FirstChild; c != nil; c = c.NextSibling {
		size += c.Size()
	}

	return size
}

// GetType returns the type of the token.
//
// Returns:
//   - TokenTyper: The type of the token.
func (t Token[S]) GetType() S {
	return t.Type
}

// IsLeaf checks if the token is a leaf.
//
// Returns:
//   - bool: True if the token is a leaf, false otherwise.
func (t Token[S]) IsLeaf() bool {
	return t.FirstChild == nil
}

// Iterator implements the ast.Noder interface.
func (t Token[S]) Iterator() itr.Iterable {
	return &TokenIterator[S]{
		first:   t.FirstChild,
		current: t.FirstChild,
	}
}

// ReverseIterator implements the ast.Noder interface.
func (t Token[S]) ReverseIterator() itr.Iterable {
	return &TokenReverseIterator[S]{
		last:    t.LastChild,
		current: t.LastChild,
	}
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

// Cleanup cleans up the token.
//
// Returns:
//   - []*Token[S]: The children of the token.
func (tk *Token[S]) Cleanup() []*Token[S] {
	var prev, next *Token[S]

	if tk.PrevSibling != nil {
		prev = tk.PrevSibling
	}

	if tk.NextSibling != nil {
		next = tk.NextSibling
	}

	if prev != nil {
		prev.NextSibling = next
	}

	if next != nil {
		next.PrevSibling = prev
	}

	var children []*Token[S]

	for c := tk.FirstChild; c != nil; c = c.NextSibling {
		c.Parent = nil

		children = append(children, c)
	}

	tk.NextSibling = nil
	tk.PrevSibling = nil
	tk.Lookahead = nil
	tk.Parent = nil
	tk.FirstChild = nil
	tk.LastChild = nil

	return children
}
