package traversing

import (
	gcers "github.com/PlayerR9/go-commons/errors"
	dbg "github.com/PlayerR9/go-debug/assert"

	itr "github.com/PlayerR9/go-commons/iterator"
)

// Traverser is an interface that defines the behavior of a traverser.
type Traverser interface {
	// Reset resets the traverser. Used for initialization.
	Reset()

	// Apply applies the traverser on the node.
	//
	// Parameters:
	//   - node: The node. Assumed to be non-nil.
	//
	// Returns:
	//   - []TravData: The children of the node.
	//   - error: An error if the traversal failed.
	//
	// WARNING: Should not be called directly. Use Apply instead.
	Apply(node TreeNoder) ([]TravData, error)
}

// TreeNoder is an interface that defines the behavior of a tree node.
type TreeNoder interface {
	// IsLeaf is a method that checks whether the node is a leaf.
	//
	// Returns:
	//   - bool: True if the node is a leaf, false otherwise.
	IsLeaf() bool

	// String is a method that returns a string representation of the node.
	//
	// Returns:
	//   - string: The string representation.
	String() string

	// Iterator is a method that returns an iterator that iterates over the direct
	// children of the node.
	//
	// Returns:
	//   - Iterater: The iterator. Never returns nil.
	Iterator() itr.Iterable

	// ReverseIterator is a method that returns an iterator that iterates over the
	// direct children of the node in reverse order.
	//
	// Returns:
	//   - iterator.Iterable: The iterator. Never returns nil.
	ReverseIterator() itr.Iterable
}

// TravData is a container for the data associated with the node before the node is visited.
type TravData struct {
	// Node is the node.
	Node TreeNoder

	// Data is the data associated with the node before the node is visited.
	Data Traverser
}

// Apply applies the Traverser on the root node.
//
// Parameters:
//   - root: The root of the AST.
//
// Returns:
//   - error: An error if the Traverser could not be applied.
func Apply(trav Traverser, root TreeNoder) error {
	if trav == nil {
		return gcers.NewErrNilParameter("trav")
	}

	trav.Reset()

	if root == nil {
		return nil
	}

	stack := []TravData{{root, trav}}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		dbg.Assert(top.Data != nil, "data must not be nil")

		children, err := top.Data.Apply(top.Node)
		if err != nil {
			return err
		}

		for i := len(children) - 1; i >= 0; i-- {
			stack = append(stack, children[i])
		}
	}

	return nil
}
