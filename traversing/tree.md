package traversing

import (
	"slices"
	"strings"

	gcers "github.com/PlayerR9/go-commons/errors"
	gr "github.com/PlayerR9/grammar/grammar"
)

// TokenTree is a token tree.
type TokenTree[S gr.TokenTyper] struct {
	// root is the root of the tree.
	root *gr.Token[S]
}

// String implements the fmt.Stringer interface.
func (t *TokenTree[S]) String() string {
	if t.root == nil {
		return ""
	}

	p := &token_printer[S]{
		lines: make([]string, 0),
		seen:  make(map[*Token[S]]bool),
	}

	se := &stack_element[S]{
		indent:     "",
		node:       t.root,
		same_level: false,
		is_last:    true,
	}

	stack := []*stack_element[S]{se}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		sub := p.trav(top)
		if len(sub) == 0 {
			continue
		}

		slices.Reverse(sub)

		stack = append(stack, sub...)
	}

	return strings.Join(p.lines, "\n")
}

// NewTokenTree creates a new token tree.
//
// Parameters:
//   - root: The root of the tree.
//
// Returns:
//   - *TokenTree: The new token tree. Never returns nil.
//   - error: An error of type *common.ErrInvalidParameter if the root is nil.
func NewTokenTree[S TokenTyper](root *Token[S]) (*TokenTree[S], error) {
	if root == nil {
		return nil, gcers.NewErrNilParameter("root")
	}

	return &TokenTree[S]{
		root: root,
	}, nil
}

// Root returns the root of the tree.
//
// Returns:
//   - *Token: The root of the tree. Never returns nil.
func (t *TokenTree[S]) Root() *Token[S] {
	return t.root
}
