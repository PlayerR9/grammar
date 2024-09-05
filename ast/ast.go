package ast

import (
	"fmt"
	"iter"

	gcers "github.com/PlayerR9/go-commons/errors"
	uttr "github.com/PlayerR9/go-commons/tree"
	gr "github.com/PlayerR9/grammar/grammar"
	internal "github.com/PlayerR9/grammar/internal"
)

// ToAstFunc is a function that converts a token to an AST node.
//
// Parameters:
//   - tk: The token. Assume tk is not nil.
//
// Returns:
//   - N: The AST node.
//   - error: An error if the function failed.
type ToAstFunc[T internal.TokenTyper, N interface {
	Child() iter.Seq[N]
	BackwardChild() iter.Seq[N]

	uttr.TreeNoder
}] func(tk *gr.Token[T]) (N, error)

// AstBuilder is an AST builder.
type AstBuilder[T internal.TokenTyper, N interface {
	Child() iter.Seq[N]
	BackwardChild() iter.Seq[N]

	uttr.TreeNoder
}] struct {
	// table is the table of the AST builder.
	table map[T]ToAstFunc[T, N]
}

// NewAstBuilder creates a new AST builder.
//
// Returns:
//   - *AstBuilder[T, N]: The new AST builder. Never returns nil.
func NewAstBuilder[T internal.TokenTyper, N interface {
	Child() iter.Seq[N]
	BackwardChild() iter.Seq[N]

	uttr.TreeNoder
}]() *AstBuilder[T, N] {
	return &AstBuilder[T, N]{
		table: make(map[T]ToAstFunc[T, N]),
	}
}

// Register registers a function to convert a token to an AST node.
//
// Parameters:
//   - type_: The type of the token.
//   - fn: The function to convert a token to an AST node.
//
// It ignores the function if it is nil and when multiple types are registered
// the previous one will be overwritten.
func (b *AstBuilder[T, N]) Register(type_ T, fn ToAstFunc[T, N]) {
	if fn == nil {
		return
	}

	b.table[type_] = fn
}

// Build builds an AST from a token. This function is an helper function that is
// used within the registered functions.
//
// Parameters:
//   - root: The root of the parse tree.
//
// Returns:
//   - N: The AST node.
//   - error: An error if the function failed.
func (b *AstBuilder[T, N]) Build(root *gr.Token[T]) (N, error) {
	if root == nil {
		return *new(N), gcers.NewErrNilParameter("root")
	}

	fn, ok := b.table[root.Type]
	if !ok {
		return *new(N), NewErrIn(root.Type, fmt.Errorf("unknown token type: %q", root.Type.String()))
	}

	node, err := fn(root)
	if err != nil {
		return node, NewErrIn(root.Type, err)
	}

	return node, nil
}

// Make creates an AST from a tree.
//
// Parameters:
//   - tree: The tree to create the AST from.
//
// Returns:
//   - *tree.Tree[N]: The AST.
//   - error: An error if the function failed.
func (b *AstBuilder[T, N]) Make(tree *uttr.Tree[*gr.Token[T]]) (*uttr.Tree[N], error) {
	if tree == nil {
		return nil, gcers.NewErrNilParameter("tree")
	}

	root := tree.Root()

	node, err := b.Build(root)
	if err != nil {
		return nil, NewErrIn(root.Type, err)
	}

	final_tree := uttr.NewTree(node)

	return final_tree, nil
}

// LhsAst is a helper function for building ASTs with a given LHS according
// to the rule LHS -> RHS LHS?
//
// Parameters:
//   - root: The root token.
//   - lhs: The LHS token type.
//   - f: The function that builds the AST.
//
// Returns:
//   - []N: The extracted nodes.
//   - error: An error if any.
func LhsAst[T internal.TokenTyper, N interface {
	Child() iter.Seq[N]
	BackwardChild() iter.Seq[N]

	uttr.TreeNoder
}](root *gr.Token[T], lhs T, f func(children []*gr.Token[T]) ([]N, error)) ([]N, error) {
	if root == nil {
		return nil, gcers.NewErrNilParameter("root")
	} else if root.Type != lhs {
		return nil, NewErrIn(lhs, fmt.Errorf("expected token type %q, got %q instead", lhs.String(), root.Type.String()))
	}

	var nodes []N

	for root != nil {
		children := root.Children()
		if len(children) == 0 {
			return nil, NewErrIn(lhs, fmt.Errorf("expected at least 1 child, got 0 instead"))
		}

		var sub_children []*gr.Token[T]

		if children[len(children)-1].Type == lhs {
			root = children[len(children)-1]
			sub_children = children[: len(children)-1 : len(children)-1]
		} else {
			root = nil
			sub_children = children
		}

		sub_nodes, err := f(sub_children)
		if err != nil {
			return nil, NewErrIn(lhs, err)
		}

		nodes = append(nodes, sub_nodes...)
	}

	return nodes, nil
}
