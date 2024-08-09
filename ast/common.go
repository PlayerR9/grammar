package ast

import (
	"fmt"

	gcers "github.com/PlayerR9/go-commons/errors"
	gr "github.com/PlayerR9/grammar/grammar"
)

// PrintAst stringifies the AST.
//
// Parameters:
//   - root: The root of the AST.
//
// Returns:
//   - string: The AST as a string.
func PrintAst[N Noder](root N) string {
	str, _ := PrintTree(root)
	// luc.AssertErr(err, "PrintTree(root)")

	return str
}

// LeftAstFunc is a function that parses the left-recursive AST.
//
// Parameters:
//   - children: The children of the current node.
//
// Returns:
//   - []N: The left-recursive AST.
//   - error: An error if the left-recursive AST could not be parsed.
type LeftAstFunc[N Noder, T gr.TokenTyper] func(children []*gr.Token[T]) ([]N, error)

// LeftRecursive parses the left-recursive AST.
//
// Parameters:
//   - root: The root of the left-recursive AST.
//   - lhs_type: The type of the left-hand side.
//   - f: The function that parses the left-recursive AST.
//
// Returns:
//   - []N: The left-recursive AST.
//   - error: An error if the left-recursive AST could not be parsed.
func LeftRecursive[N Noder, T gr.TokenTyper](root *gr.Token[T], lhs_type T, f LeftAstFunc[N, T]) ([]N, error) {
	// luc.AssertNil(root, "root")

	var nodes []N

	for root != nil {
		if root.Type != lhs_type {
			return nodes, fmt.Errorf("expected %q, got %q instead", lhs_type.String(), root.Type.String())
		}

		children, err := ExtractChildren(root)
		if err != nil {
			return nodes, err
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
type ToAstFunc[N Noder, T gr.TokenTyper] func(root *gr.Token[T]) ([]N, error)

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
func ToAst[N Noder, T gr.TokenTyper](root *gr.Token[T], to_ast ToAstFunc[N, T]) ([]N, error) {
	if root == nil {
		return nil, gcers.NewErrNilParameter("root")
	} else if to_ast == nil {
		return nil, gcers.NewErrNilParameter("to_ast")
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
		return "", gcers.NewErrNilParameter("node")
	}

	return node.Data, nil
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
		return nil, gcers.NewErrNilParameter("node")
	}

	var children []*gr.Token[T]

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, c)
	}

	return children, nil
}
