package ast

import (
	"errors"
	"fmt"

	lus "github.com/PlayerR9/lib_units/slices"
)

// Result is the result of the AST.
type Result[N NodeTyper] struct {
	// nodes is the nodes of the result.
	nodes []*Node[N]

	// err is the error of the result.
	err error
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

// SetError sets the error of the result.
// Does nothing if the error is nil.
//
// Parameters:
//   - err: The error to set.
func (a *Result[N]) SetError(err error) {
	if err != nil {
		a.err = err
	}
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
func (a *Result[N]) AppendChildren(children []*Node[N]) {
	children = lus.FilterNilValues(children)

	if len(children) == 0 {
		return
	}

	if len(a.nodes) == 0 {
		if a.err == nil {
			a.err = errors.New("no node to append children to")
		}
	} else if len(a.nodes) > 1 {
		if a.err == nil {
			a.err = errors.New("cannot append children to multiple nodes")
		}
	} else {
		a.nodes[0].AppendChildren(children)
	}
}

// Apply applies the result.
//
// Returns:
//   - []*Node[N]: The nodes of the result.
//   - error: The error of the result.
func (a *Result[N]) Apply() ([]*Node[N], error) {
	return a.nodes, a.err
}

// IsError returns true if the result is an error.
//
// Returns:
//   - bool: True if the result is an error. False otherwise.
func (a *Result[N]) IsError() bool {
	return a.err != nil
}

// DoFunc does something with the result.
//
// Parameters:
//   - f: The function to do something with the result.
//   - prev: The previous result of the function.
//
// Returns:
//   - any: The result of the function.
//
// This function does nothing if f is nil or an error is set.
func (a *Result[N]) DoFunc(f DoFunc[N], prev any) any {
	if f == nil || a.err != nil {
		return nil
	}

	return f(a, prev)
}

// Exec executes a set of functions on the result.
//
// Parameters:
//   - fs: The functions to execute.
//
// Returns:
//   - []*Node[N]: The nodes of the result.
//   - error: The error of the result.
//
// This function does nothing ignores the functions that are nil.
func (a *Result[N]) Exec(fs []DoFunc[N]) ([]*Node[N], error) {
	var top int

	for i := 0; i < len(fs); i++ {
		if fs[i] != nil {
			fs[top] = fs[i]
			top++
		}
	}

	fs = fs[:top]

	if len(fs) == 0 {
		return a.Apply()
	}

	var prev any

	for _, f := range fs {
		prev = f(a, prev)
		if a.IsError() {
			prev = nil

			break
		}
	}

	if prev != nil {
		panic(fmt.Sprintf("Last function returned (%v) instead of nil. Maybe you forgot to specify a function?", prev))
	}

	return a.Apply()
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
