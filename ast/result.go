package ast

import (
	"errors"

	gcers "github.com/PlayerR9/go-commons/errors"
	gcint "github.com/PlayerR9/go-commons/ints"
	gcslc "github.com/PlayerR9/go-commons/slices"
)

// Result is the result of the AST.
type Result[N Noder] struct {
	// nodes is the nodes of the result.
	nodes []N
}

// NewResult creates a new AstResult.
//
// Returns:
//   - *AstResult[N]: The new AstResult. Never returns nil.
func NewResult[N Noder]() *Result[N] {
	return &Result[N]{}
}

// SetNode sets the node of the result. It replaces any existing node.
//
// Parameters:
//   - node: The node to set.
func (a *Result[N]) SetNode(node N) {
	a.nodes = []N{node}
}

// SetNodes sets the nodes of the result. It ignores the nodes that are nil.
//
// Parameters:
//   - nodes: The nodes to set.
func (a *Result[N]) SetNodes(nodes []N) {
	if len(nodes) > 0 {
		a.nodes = nodes
	}
}

// AppendNodes appends the nodes of the result. It ignores the nodes that are nil.
//
// Parameters:
//   - nodes: The nodes to append.
func (a *Result[N]) AppendNodes(nodes []N) {
	if len(nodes) > 0 {
		a.nodes = append(a.nodes, nodes...)
	}
}

// AppendChildren appends the children to the node. It ignores the children that are nil.
//
// Parameters:
//   - children: The children to append.
//
// Returns:
//   - error: An error if the result is an error.
func (a *Result[N]) AppendChildren(children []Noder) error {
	children = gcslc.SliceFilter(children, filter_non_nil_noders)
	if len(children) == 0 {
		return nil
	}

	if len(a.nodes) == 0 {
		return errors.New("no node to append children to")
	} else if len(a.nodes) > 1 {
		return errors.New("cannot append children to multiple nodes")
	}

	a.nodes[0].AddChildren(children)

	return nil
}

// Apply applies the result.
//
// Returns:
//   - []N: The nodes of the result.
func (a *Result[N]) Apply() []N {
	return a.nodes
}

// DoFunc does something with the result.
//
// Parameters:
//   - f: The function to do something with the result.
//   - prev: The previous result of the function.
//
// Returns:
//   - any: The result of the function.
//   - error: An error if the function failed.
//
// Errors:
//   - *common.ErrInvalidParameter: If the f is nil.
//   - error: Any error returned by the f function.
func (a *Result[N]) DoFunc(f DoFunc[N], prev any) (any, error) {
	if f == nil {
		return nil, gcers.NewErrNilParameter("f")
	}

	res, err := f(a, prev)
	if err != nil {
		return res, err
	}

	return res, nil
}

// DoForEach does something with the nodes of the result.
//
// Parameters:
//   - f: The function to do something with the nodes of the result.
//
// Returns:
//   - error: An error if the function failed.
//
// Errors:
//   - *common.ErrAt: With the error of the 'f' function.
func (a *Result[N]) DoForEach(f func(N) error) error {
	if len(a.nodes) == 0 || f == nil {
		return nil
	}

	for i, node := range a.nodes {
		err := f(node)
		if err != nil {
			return gcint.NewErrAt(i+1, "node", err)
		}
	}

	return nil
}
