package ast

import (
	"fmt"
)

// NodeTyper is an interface that defines the behavior of a node type.
type NodeTyper interface {
	~int

	fmt.Stringer
}

// Noder is an interface that defines the behavior of a node.
type Noder interface {
	// AddChild is a method that adds a child to the node. If the child is nil or not
	// of the correct type, nothing happens.
	//
	// Parameters:
	//   - child: The child to add.
	AddChild(child Noder)

	// AddChildren is a convenience function to add multiple children to the node at once.
	// It is more efficient than adding them one by one. Therefore, the behaviors are the
	// same as the behaviors of the Noder.AddChild function.
	//
	// Parameters:
	//   - children: The children to add.
	AddChildren(children []Noder)

	// IsLeaf is a method that checks if the node is a leaf.
	//
	// Returns:
	//   - bool: True if the node is a leaf, false otherwise.
	IsLeaf() bool

	// String is a method that returns a string representation of the node.
	//
	// Returns:
	//   - string: The string representation of the node.
	String() string
}
