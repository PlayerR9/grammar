package traversing

import (
	"io"

	dbg "github.com/PlayerR9/go-debug/assert"
	ast "github.com/PlayerR9/grammar/ast"
)

// SimpleDFSDoFunc is a function that is called for each node.
//
// Parameters:
//   - node: The node.
//
// Returns:
//   - error: An error if the DFS could not be applied.
type SimpleDFSDoFunc[N ast.Noder] func(node N) error

// SimpleDFS is a simple depth-first search.
type SimpleDFS[N ast.Noder] struct {
	// do_func is the function that is called for each node.
	do_func SimpleDFSDoFunc[N]
}

// NewSimpleDFS creates a new SimpleDFS.
//
// Parameters:
//   - f: The function that is called for each node.
//
// Returns:
//   - SimpleDFS[N]: The new SimpleDFS.
//
// If f is nil, simpleDFS is returned as nil.
func NewSimpleDFS[N ast.Noder](f SimpleDFSDoFunc[N]) SimpleDFS[N] {
	if f == nil {
		return SimpleDFS[N]{}
	}

	return SimpleDFS[N]{do_func: f}
}

// SetDoFunc sets the function that is called for each node.
//
// Parameters:
//   - f: The function that is called for each node.
func (s *SimpleDFS[N]) SetDoFunc(f SimpleDFSDoFunc[N]) {
	s.do_func = f
}

// Apply applies the SimpleDFS. Does nothing if the do_func is nil.
//
// Parameters:
//   - root: The root of the AST.
//
// Returns:
//   - error: An error if the SimpleDFS could not be applied.
func (s SimpleDFS[N]) Apply(root N) error {
	if s.do_func == nil {
		return nil
	}

	stack := []N{root}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		err := s.do_func(top)
		if err != nil {
			return err
		}

		if top.IsLeaf() {
			continue
		}

		iter := top.ReverseIterator()
		dbg.AssertNotNil(iter, "iter")

		for {
			value, err := iter.Consume()
			if err == io.EOF {
				break
			}

			dbg.AssertErr(err, "iter.Consume()")

			child := dbg.AssertConv[N](value, "value")
			stack = append(stack, child)
		}
	}

	return nil
}
