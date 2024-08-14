package ast

import (
	"fmt"

	gcers "github.com/PlayerR9/go-commons/errors"
	gr "github.com/PlayerR9/grammar/grammar"
)

// DoFunc is a function that does something with the AST.
//
// Parameters:
//   - a: The result of the AST.
//   - root: The root of the parse tree.
//
// Returns:
//   - error: An error if the function failed.
type DoFunc[N Noder, T gr.TokenTyper] func(a *Result[N], root *gr.Token[T]) error

// Make is the constructor for the AST.
type Make[N Noder, T gr.TokenTyper] struct {
	// ast_map is the map of the AST.
	ast_map map[T]DoFunc[N, T]
}

// AddEntry adds an entry to the AST. Nil steps are ignored.
//
// Parameters:
//   - t: The type of the entry.
//   - steps: The steps of the entry.
//
// Returns:
//   - error: An error if no steps were provided or if the entry already exists.
func (m *Make[N, T]) AddEntry(t T, step DoFunc[N, T]) error {
	if step == nil {
		return gcers.NewErrNilParameter("step")
	}

	if m.ast_map == nil {
		m.ast_map = make(map[T]DoFunc[N, T])
	}

	_, ok := m.ast_map[t]
	if ok {
		return fmt.Errorf("entry with type %q already exists", t.String())
	}

	m.ast_map[t] = step

	return nil
}

// Apply creates the AST given a token (most often the root).
//
// Parameters:
//   - token: The token to create the AST from.
//
// Returns:
//   - []N: The AST.
//   - error: An error if the AST could not be created.
func (m Make[N, T]) Apply(token *gr.Token[T]) ([]N, error) {
	if token == nil {
		return nil, gcers.NewErrNilParameter("tree")
	}

	step, ok := m.ast_map[token.Type]
	if !ok {
		return nil, fmt.Errorf("unexpected token type: %q", token.Type.String())
	}

	var res Result[N]

	err := step(&res, token)
	nodes := res.Apply()

	if err != nil {
		return nodes, NewErrInRule(token.Type, err)
	}

	return nodes, nil
}
