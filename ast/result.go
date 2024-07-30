package ast

import (
	"errors"

	luc "github.com/PlayerR9/lib_units/common"
	lus "github.com/PlayerR9/lib_units/slices"
)

// Result is the result of the AST.
type Result[N NodeTyper] struct {
	// nodes is the nodes of the result.
	nodes []*Node[N]
}

// NewResult creates a new AstResult.
//
// Returns:
//   - *AstResult[N]: The new AstResult. Never returns nil.
func NewResult[N NodeTyper]() *Result[N] {
	return &Result[N]{}
}

// MakeNode creates a new node and adds it to the result; replacing any existing nodes.
//
// Parameters:
//   - t: The type of the node.
//   - data: The data of the node.
func (a *Result[N]) MakeNode(t N, data string) {
	n := NewNode(t, data)

	a.nodes = []*Node[N]{n}
}

// SetNodes sets the nodes of the result. It ignores the nodes that are nil.
//
// Parameters:
//   - nodes: The nodes to set.
func (a *Result[N]) SetNodes(nodes []*Node[N]) {
	nodes = lus.FilterNilValues(nodes)
	if len(nodes) > 0 {
		a.nodes = nodes
	}
}

// AppendNodes appends the nodes of the result. It ignores the nodes that are nil.
//
// Parameters:
//   - nodes: The nodes to append.
func (a *Result[N]) AppendNodes(nodes []*Node[N]) {
	nodes = lus.FilterNilValues(nodes)

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
func (a *Result[N]) AppendChildren(children []*Node[N]) error {
	children = lus.FilterNilValues(children)

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
//   - []*Node[N]: The nodes of the result.
func (a *Result[N]) Apply() []*Node[N] {
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
		return nil, luc.NewErrNilParameter("f")
	}

	res, err := f(a, prev)
	if err != nil {
		return res, err
	}

	return res, nil
}

// TransformNodes transforms the nodes of the result.
//
// Parameters:
//   - new_type: The new type of the nodes.
//   - new_data: The new data of the nodes.
func (a *Result[N]) TransformNodes(new_type N, new_data string) {
	if len(a.nodes) == 0 {
		return
	}

	for _, node := range a.nodes {
		node.Type = new_type
		node.Data = new_data
	}
}
