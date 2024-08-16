package ast

import (
	gcers "github.com/PlayerR9/go-commons/errors"
	dbg "github.com/PlayerR9/go-debug/assert"
)

// TravData is a container for the data associated with the node before the node is visited.
type TravData[T any] struct {
	// Node is the node.
	Node T

	// Data is the data associated with the node before the node is visited.
	Data interface {
		Reset()
		Apply(node T) ([]TravData[T], error)
	}
}

// Apply applies the Traverser on the root node.
//
// Parameters:
//   - root: The root of the AST.
//
// Returns:
//   - error: An error if the Traverser could not be applied.
func Apply[T any](trav interface {
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
	Apply(node T) ([]TravData[T], error)
}, root T) error {
	if trav == nil {
		return gcers.NewErrNilParameter("trav")
	}

	trav.Reset()

	/* if root == nil {
		return nil
	} */

	stack := []TravData[T]{{root, trav}}

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
