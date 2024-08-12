package ast

import (
	"fmt"
)

// NodeTyper is an interface that defines the behavior of a node type.
type NodeTyper interface {
	~int

	fmt.Stringer
}

// Iterater is an interface that defines the behavior of an iterator.
type Iterater interface {
	// Consume is a method that consumes the next node in the iterator.
	//
	// Returns:
	//   - Noder: The next node in the iterator.
	//   - error: An error of type io.EOF if there are no more nodes in the iterator.
	//
	// The returned node is never nil; unless an error is returned.
	Consume() (Noder, error)

	// Restart is a method that restarts the iterator.
	Restart()
}

// Noder is an interface that defines the behavior of a node.
type Noder interface {
	// IsLeaf is a method that checks whether the node is a leaf.
	//
	// Returns:
	//   - bool: True if the node is a leaf, false otherwise.
	IsLeaf() bool

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

	// Iterator is a method that returns an iterator that iterates over the direct
	// children of the node.
	//
	// Returns:
	//   - Iterater: The iterator. Never returns nil.
	Iterator() Iterater

	// ReverseIterator is a method that returns an iterator that iterates over the
	// direct children of the node in reverse order.
	//
	// Returns:
	//   - Iterater: The iterator. Never returns nil.
	ReverseIterator() Iterater

	fmt.Stringer
}
