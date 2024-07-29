package ast

import (
	"fmt"

	gr "github.com/PlayerR9/grammar/grammar"
	luc "github.com/PlayerR9/lib_units/common"
)

// PrintAst stringifies the AST.
//
// Parameters:
//   - root: The root of the AST.
//
// Returns:
//   - string: The AST as a string.
func PrintAst[N NodeTyper](root *Node[N]) string {
	if root == nil {
		return ""
	}

	str, err := gr.PrintTree(root)
	luc.AssertErr(err, "Strings.PrintTree(root)")

	return str
}

// LeftAstFunc is a function that parses the left-recursive AST.
//
// Parameters:
//   - children: The children of the current node.
//
// Returns:
//   - []*Node[N]: The left-recursive AST.
//   - error: An error if the left-recursive AST could not be parsed.
type LeftAstFunc[N NodeTyper, T gr.TokenTyper] func(children []*gr.Token[T]) ([]*Node[N], error)

// LeftRecursive parses the left-recursive AST.
//
// Parameters:
//   - root: The root of the left-recursive AST.
//   - lhs_type: The type of the left-hand side.
//   - f: The function that parses the left-recursive AST.
//
// Returns:
//   - []*Node[N]: The left-recursive AST.
//   - error: An error if the left-recursive AST could not be parsed.
func LeftRecursive[N NodeTyper, T gr.TokenTyper](root *gr.Token[T], lhs_type T, f LeftAstFunc[N, T]) ([]*Node[N], error) {
	luc.AssertNil(root, "root")

	var nodes []*Node[N]

	for root != nil {
		if root.Type != lhs_type {
			return nodes, fmt.Errorf("expected %q, got %q instead", lhs_type.String(), root.Type.String())
		}

		children, ok := root.Data.([]*gr.Token[T])
		if !ok {
			return nodes, fmt.Errorf("expected non-leaf node, got leaf node instead")
		} else if len(children) == 0 {
			return nodes, fmt.Errorf("expected at least 1 child, got 0 children instead")
		}

		last_child := children[len(children)-1]

		if last_child.Type == lhs_type {
			children = children[:len(children)-1]
			root = last_child
		} else {
			root = nil
		}

		sub_nodes, err := f(children)
		if len(sub_nodes) > 0 {
			nodes = append(nodes, sub_nodes...)
		}

		if err != nil {
			return nodes, fmt.Errorf("in %q: %w", root.Type.String(), err)
		}
	}

	return nodes, nil
}

// ToAstFunc is a function that parses the AST.
//
// Parameters:
//   - root: The root of the AST.
//
// Returns:
//   - []*Node[N]: The AST.
//   - error: An error if the AST could not be parsed.
type ToAstFunc[N NodeTyper, T gr.TokenTyper] func(root *gr.Token[T]) ([]*Node[N], error)

// ToAst parses the AST.
//
// Parameters:
//   - root: The root of the AST.
//   - to_ast: The function that parses the AST.
//
// Returns:
//   - []N: The AST.
//   - error: An error if the AST could not be parsed.
//
// Errors:
//   - *common.ErrInvalidParameter: If the root is nil or the to_ast is nil.
//   - error: Any error returned by the to_ast function.
func ToAst[N NodeTyper, T gr.TokenTyper](root *gr.Token[T], to_ast ToAstFunc[N, T]) ([]*Node[N], error) {
	if root == nil {
		return nil, luc.NewErrNilParameter("root")
	} else if to_ast == nil {
		return nil, luc.NewErrNilParameter("to_ast")
	}

	nodes, err := to_ast(root)
	if err != nil {
		return nodes, err
	}

	return nodes, nil
}

// ExtractData extracts the data from a token.
//
// Parameters:
//   - node: The token to extract the data from.
//
// Returns:
//   - string: The data of the token.
//   - error: An error if the data is not of type string or if the token is nil.
func ExtractData[T gr.TokenTyper](node *gr.Token[T]) (string, error) {
	if node == nil {
		return "", luc.NewErrNilParameter("node")
	}

	data, ok := node.Data.(string)
	if !ok {
		return "", fmt.Errorf("expected string, got %T instead", node.Data)
	}

	return data, nil
}

// ExtractChildren extracts the children from a token.
//
// Parameters:
//   - node: The token to extract the children from.
//
// Returns:
//   - []*gr.Token[T]: The children of the token.
//   - error: An error if the children is not of type []*gr.Token[T] or if the token is nil.
func ExtractChildren[T gr.TokenTyper](node *gr.Token[T]) ([]*gr.Token[T], error) {
	if node == nil {
		return nil, luc.NewErrNilParameter("node")
	}

	children, ok := node.Data.([]*gr.Token[T])
	if !ok {
		return nil, fmt.Errorf("expected []*Token, got %T instead", node.Data)
	}

	return children, nil
}
