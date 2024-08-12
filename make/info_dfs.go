package make

import (
	"io"

	dbg "github.com/PlayerR9/go-debug/assert"
	ast "github.com/PlayerR9/grammar/ast"
)

// InfoDFSDoFunc is a function that is called for each node.
//
// Parameters:
//   - node: The node.
//   - data: The data.
//
// Returns:
//   - error: An error if the DFS could not be applied.
type InfoDFSDoFunc[N ast.Noder, I interface{ Copy() I }] func(node N, data I) error

// InitFunc is a function that initializes the data.
//
// Returns:
//   - I: The data.
type InitFunc[N ast.Noder, I interface{ Copy() I }] func() I

// InfoDFS is a simple depth-first search.
type InfoDFS[N ast.Noder, I interface{ Copy() I }] struct {
	// do_func is the function that is called for each node.
	do_func InfoDFSDoFunc[N, I]

	// init is the function that initializes the data.
	init InitFunc[N, I]
}

// NewInfoDFS creates a new InfoDFS.
//
// Parameters:
//   - f: The function that is called for each node.
//
// Returns:
//   - *InfoDFS[N, I]: The new InfoDFS.
//
// If f is nil, infoDFS is returned as nil.
// If init is nil, the default init function is used which returns the zero value of I.
func NewInfoDFS[N ast.Noder, I interface{ Copy() I }](f InfoDFSDoFunc[N, I], init InitFunc[N, I]) *InfoDFS[N, I] {
	if f == nil {
		return nil
	}

	if init == nil {
		init = func() I { return *new(I) }
	}

	return &InfoDFS[N, I]{
		do_func: f,
		init:    init,
	}
}

// Apply applies the InfoDFS.
//
// Parameters:
//   - root: The root of the AST.
//
// Returns:
//   - I: The data.
//   - error: An error if the InfoDFS could not be applied.
func (s InfoDFS[N, I]) Apply(root N) (I, error) {
	dbg.Assert(s.do_func != nil, "do_func should not be nil")

	type stack_elem struct {
		node N
		data I
	}

	data := s.init()

	stack := []stack_elem{{root, data}}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		err := s.do_func(top.node, top.data)
		if err != nil {
			return data, err
		}

		if top.node.IsLeaf() {
			continue
		}

		iter := top.node.Iterator()
		dbg.Assert(iter != nil, "iterator expected to be non-nil")

		for {
			value, err := iter.Consume()
			if err == io.EOF {
				break
			}

			dbg.AssertErr(err, "iter.Consume()")

			child := dbg.AssertConv[N](value, "value")
			stack = append(stack, stack_elem{child, top.data.Copy()})
		}
	}

	return data, nil
}
