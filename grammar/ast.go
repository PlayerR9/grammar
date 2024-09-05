package grammar

import (
	"fmt"
	"iter"

	internal "github.com/PlayerR9/grammar/grammar/internal"

	gcers "github.com/PlayerR9/go-commons/errors"
	uttr "github.com/PlayerR9/go-commons/tree"
)

type ToAstFunc[T internal.TokenTyper, N interface {
	Child() iter.Seq[N]
	BackwardChild() iter.Seq[N]

	uttr.TreeNoder
}] func(tk *Token[T]) (N, error)

type AstBuilder[T internal.TokenTyper, N interface {
	Child() iter.Seq[N]
	BackwardChild() iter.Seq[N]

	uttr.TreeNoder
}] struct {
	table map[T]ToAstFunc[T, N]
}

func NewAstBuilder[T internal.TokenTyper, N interface {
	Child() iter.Seq[N]
	BackwardChild() iter.Seq[N]

	uttr.TreeNoder
}]() *AstBuilder[T, N] {
	return &AstBuilder[T, N]{
		table: make(map[T]ToAstFunc[T, N]),
	}
}

func (b *AstBuilder[T, N]) Register(type_ T, fn ToAstFunc[T, N]) {
	if fn == nil {
		return
	}

	b.table[type_] = fn
}

func (b *AstBuilder[T, N]) Build(root *Token[T]) (N, error) {
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

func (b *AstBuilder[T, N]) Make(tree *uttr.Tree[*Token[T]]) (*uttr.Tree[N], error) {
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
//   - []N: The AST.
//   - error: An error if any.
//
// Panics if the function fails.
func LhsAst[T internal.TokenTyper, N interface {
	Child() iter.Seq[N]
	BackwardChild() iter.Seq[N]

	uttr.TreeNoder
}](root *Token[T], lhs T, f func(children []*Token[T]) ([]N, error)) ([]N, error) {
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

		var sub_children []*Token[T]

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
