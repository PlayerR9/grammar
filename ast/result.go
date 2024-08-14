package ast

import (
	"errors"

	gcslc "github.com/PlayerR9/go-commons/slices"
)

// Result is the result of the AST.
type Result[N Noder] struct {
	// nodes is the nodes of the result.
	nodes []N
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
func (a Result[N]) Apply() []N {
	return a.nodes
}
