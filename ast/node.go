package ast

import (
	"fmt"
	"io"

	lls "github.com/PlayerR9/listlike/stack"
)

// NodeTyper is an interface that defines the behavior of a node type.
type NodeTyper interface {
	~int

	fmt.Stringer
}

type NoderIterater interface {
	Consume() (Noder, error)
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
	//   - NoderIterater: The iterator.
	Iterator() NoderIterater

	fmt.Stringer
}

// DFSDoFunc is a function that is called for each node.
//
// Parameters:
//   - node: The node.
//   - data: The data.
//
// Returns:
//   - error: An error if the DFS could not be applied.
type DFSDoFunc[N Noder, I any] func(node N, data I) error

// InitFunc is a function that initializes the data.
//
// Returns:
//   - I: The data.
type InitFunc[N Noder, I any] func() I

// SimpleDFS is a simple depth-first search.
type SimpleDFS[N Noder, I any] struct {
	// do_func is the function that is called for each node.
	do_func DFSDoFunc[N, I]

	// init is the function that initializes the data.
	init InitFunc[N, I]
}

// NewSimpleDFS creates a new SimpleDFS.
//
// Parameters:
//   - f: The function that is called for each node.
//
// Returns:
//   - *SimpleDFS[N, I]: The new SimpleDFS.
//
// If f is nil, simpleDFS is returned as nil.
// If init is nil, the default init function is used which returns the zero value of I.
func NewSimpleDFS[N Noder, I any](f DFSDoFunc[N, I], init InitFunc[N, I]) *SimpleDFS[N, I] {
	if f == nil {
		return nil
	}

	if init == nil {
		init = func() I { return *new(I) }
	}

	return &SimpleDFS[N, I]{
		do_func: f,
		init:    init,
	}
}

// Apply applies the SimpleDFS.
//
// Parameters:
//   - root: The root of the AST.
//
// Returns:
//   - I: The data.
//   - error: An error if the SimpleDFS could not be applied.
func (s *SimpleDFS[N, I]) Apply(root N) (I, error) {
	stack := lls.NewLinkedStack[N]()
	stack.Push(root)

	data := s.init()

	for {
		top, ok := stack.Pop()
		if !ok {
			break
		}

		err := s.do_func(top, data)
		if err != nil {
			return data, err
		}

		if top.IsLeaf() {
			continue
		}

		iter := top.Iterator()
		// luc.Assert(iter != nil, "iterator expected to be non-nil")

		for {
			value, err := iter.Consume()
			if err == io.EOF {
				break
			}

			// luc.AssertErr(err, "iter.Consume()")
			child := value.(N)
			// child := luc.AssertConv[N](value, "value")

			stack.Push(child)
		}
	}

	return data, nil
}
