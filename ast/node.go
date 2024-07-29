package ast

import (
	"fmt"
	"strconv"
	"strings"

	gr "github.com/PlayerR9/grammar/grammar"
	luc "github.com/PlayerR9/lib_units/common"
	lus "github.com/PlayerR9/lib_units/slices"
)

// NodeTyper is an interface that defines the behavior of a node type.
type NodeTyper interface {
	~int

	fmt.Stringer
}

// Node is a node in the AST.
type Node[N NodeTyper] struct {
	// Parent is the parent of the node.
	Parent *Node[N]

	// Children is the children of the node.
	Children []*Node[N]

	// Type is the type of the node.
	Type N

	// Data is the data of the node.
	Data string
}

// IsLeaf implements the grammar.Noder interface.
func (n *Node[N]) IsLeaf() bool {
	return len(n.Children) == 0
}

// Iterator implements the grammar.Noder interface.
func (n *Node[N]) Iterator() luc.Iterater[gr.Noder] {
	if len(n.Children) == 0 {
		return nil
	}

	nodes := make([]gr.Noder, 0, len(n.Children))
	for _, child := range n.Children {

		nodes = append(nodes, child)
	}

	return luc.NewSimpleIterator(nodes)
}

// String implements the grammar.Noder interface.
func (n *Node[N]) String() string {
	var builder strings.Builder

	builder.WriteString("Node[")
	builder.WriteString(n.Type.String())

	if n.Data != "" {
		builder.WriteString(" (")
		builder.WriteString(strconv.Quote(n.Data))
		builder.WriteRune(')')
	}

	builder.WriteRune(']')

	return builder.String()
}

// NewNode creates a new node.
//
// Parameters:
//   - t: The type of the node.
//   - data: The data of the node.
//
// Returns:
//   - *Node[N]: The new node. Never returns nil.
func NewNode[N NodeTyper](t N, data string) *Node[N] {
	return &Node[N]{
		Type: t,
		Data: data,
	}
}

// AppendChildren implements the *Node[N] interface.
func (n *Node[N]) AppendChildren(children []*Node[N]) {
	children = lus.FilterNilValues(children)
	if len(children) == 0 {
		return
	}

	for _, child := range children {
		child.Parent = n
	}

	n.Children = append(n.Children, children...)
}

// SetChildren sets the children of the node. Nil children are ignored.
func (n *Node[N]) SetChildren(children []*Node[N]) {
	children = lus.FilterNilValues(children)
	if len(children) == 0 {
		return
	}

	for _, child := range children {
		child.Parent = n
	}

	n.Children = children
}
