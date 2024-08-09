// Code generated by go generate; do not edit.
package ast

import (
	"io"
	"strconv"
	"strings"	
)

// NodeIterator is a pull-based iterator that iterates over the children of a Node.
type NodeIterator[N NodeTyper] struct {
	first, current *Node[N]
}

// Consume implements the Iterater interface.
func (iter *NodeIterator[N]) Consume() (Noder, error) {
	n := iter.current

	if n == nil {
		return nil, io.EOF
	}

	iter.current = n.NextSibling

	return n, nil
}

// Restart implements the Iterater interface.
func (iter *NodeIterator[N]) Restart() {
	iter.current = iter.first
}

// Node is a node in a ast.
type Node[N NodeTyper] struct {
	Parent, FirstChild, NextSibling, LastChild, PrevSibling *Node[N]

	Type N
	Data string
	Pos int
}

// IsLeaf implements the Noder interface.
func (tn *Node[N]) IsLeaf() bool {
	return tn.FirstChild == nil
}

// AddChild implements the Noder interface.
func (tn *Node[N]) AddChild(target Noder) {
	if target == nil {
		return
	}

	tmp, ok := target.(*Node[N])
	if !ok {
		return
	}
	
	tmp.NextSibling = nil
	tmp.PrevSibling = nil

	last_child := tn.LastChild

	if last_child == nil {
		tn.FirstChild = tmp
	} else {
		last_child.NextSibling = tmp
		tmp.PrevSibling = last_child
	}

	tmp.Parent = tn
	tn.LastChild = tmp
}

// AddChildren implements the Noder interface.
func (tn *Node[N]) AddChildren(children []Noder) {
	if len(children) == 0 {
		return
	}
	
	var valid_children []*Node[N]

	for _, child := range children {
		if child == nil {
			continue
		}

		c, ok := child.(*Node[N])
		if !ok {
			continue
		}

		valid_children = append(valid_children, c)
	}

	if len(valid_children) == 0 {
		return
	}

	// Deal with the first child
	first_child := valid_children[0]

	first_child.NextSibling = nil
	first_child.PrevSibling = nil

	last_child := tn.LastChild

	if last_child == nil {
		tn.FirstChild = first_child
	} else {
		last_child.NextSibling = first_child
		first_child.PrevSibling = last_child
	}

	first_child.Parent = tn
	tn.LastChild = first_child

	// Deal with the rest of the children
	for i := 1; i < len(valid_children); i++ {
		child := valid_children[i]

		child.NextSibling = nil
		child.PrevSibling = nil

		last_child := tn.LastChild
		last_child.NextSibling = child
		child.PrevSibling = last_child

		child.Parent = tn
		tn.LastChild = child
	}
}

// Iterator implements the Noder interface.
func (tn *Node[N]) Iterator() Iterater {
	return &NodeIterator[N]{
		first:   tn.FirstChild,
		current: tn.FirstChild,
	}
}

// String implements the Noder interface.
func (tn *Node[N]) String() string {
	var builder strings.Builder

	builder.WriteString(strconv.Itoa(tn.Pos))
	builder.WriteString(":Node[")
	builder.WriteString(tn.Type.String())

	if tn.Data != "" {
		builder.WriteString(" (")
		builder.WriteString(strconv.Quote(tn.Data))
		builder.WriteRune(')')
	}

	builder.WriteRune(']')

	return builder.String()
}

// NewNode creates a new node with the given data.
//
// Parameters:
//   - n_type: The type of the node.
//   - data: The data of the node.
//   - pos: The position of the node in the source code.
//
// Returns:
//   - *Node[N]: A pointer to the newly created node. It is never nil.
func NewNode[N NodeTyper](n_type N, data string, pos int) *Node[N] {
	return &Node[N]{
		Type: n_type,
		Data: data,
		Pos:  pos,
	}
}