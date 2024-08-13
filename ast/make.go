package ast

import (
	"errors"
	"fmt"

	gcers "github.com/PlayerR9/go-commons/errors"
	gcslc "github.com/PlayerR9/go-commons/slices"
	gr "github.com/PlayerR9/grammar/grammar"
)

// Make is the constructor for the AST.
type Make[N Noder, T gr.TokenTyper] struct {
	// ast_map is the map of the AST.
	ast_map map[T][]DoFunc[N]
}

// NewMake creates a new Make.
//
// Returns:
//   - *Make[N, T]: The new Make.
func NewMake[N Noder, T gr.TokenTyper]() *Make[N, T] {
	return &Make[N, T]{
		ast_map: make(map[T][]DoFunc[N]),
	}
}

// AddEntry adds an entry to the AST. Nil steps are ignored.
//
// Parameters:
//   - t: The type of the entry.
//   - steps: The steps of the entry.
//
// Returns:
//   - error: An error if no steps were provided or if the entry already exists.
func (m *Make[N, T]) AddEntry(t T, steps []DoFunc[N]) error {
	f := func(do DoFunc[N]) bool {
		return do != nil
	}

	steps = gcslc.SliceFilter(steps, f)
	if len(steps) == 0 {
		return errors.New("no steps provided")
	}

	if m.ast_map == nil {
		m.ast_map = make(map[T][]DoFunc[N])
	}

	_, ok := m.ast_map[t]
	if ok {
		return fmt.Errorf("entry with type %q already exists", t.String())
	}

	m.ast_map[t] = steps

	return nil
}

// Apply creates the AST given the root.
//
// Parameters:
//   - tree: The root of the AST.
//
// Returns:
//   - []N: The AST.
//   - error: An error if the AST could not be created.
func (m *Make[N, T]) Apply(root *gr.Token[T]) ([]N, error) {
	if root == nil {
		return nil, gcers.NewErrNilParameter("tree")
	}

	steps, ok := m.ast_map[root.Type]
	if !ok {
		return nil, fmt.Errorf("unexpected token type: %q", root.Type.String())
	}

	res := NewResult[N]()

	var prev any = root
	var err error

	for _, step := range steps {
		prev, err = step(res, prev)
		if err != nil {
			nodes := res.Apply()

			return nodes, fmt.Errorf("in %q: %w", root.Type.String(), err)
		}
	}

	if prev != nil {
		panic(fmt.Sprintf("Last function returned (%v) instead of nil. Did you forget to specify a function?", prev))
	}

	nodes := res.Apply()

	return nodes, nil
}

// Apply creates the AST given the token.
//
// Parameters:
//   - token: The token to create the AST from.
//
// Returns:
//   - []N: The AST.
//   - error: An error if the AST could not be created.
func (m *Make[N, T]) ApplyToken(token *gr.Token[T]) ([]N, error) {
	if token == nil {
		return nil, gcers.NewErrNilParameter("tree")
	}

	steps, ok := m.ast_map[token.Type]
	if !ok {
		return nil, fmt.Errorf("unexpected token type: %q", token.Type.String())
	}

	res := NewResult[N]()

	var prev any = token
	var err error

	for _, step := range steps {
		prev, err = step(res, prev)
		if err != nil {
			nodes := res.Apply()

			return nodes, fmt.Errorf("in %q: %w", token.Type.String(), err)
		}
	}

	if prev != nil {
		panic(fmt.Sprintf("Last function returned (%v) instead of nil. Did you forget to specify a function?", prev))
	}

	nodes := res.Apply()

	return nodes, nil
}

// CheckTokenType checks if the token is of the expected type.
//
// Parameters:
//   - tk: The token to check.
//   - tk_type: The expected type of the token.
//
// Returns:
//   - error: An error of type *ErrInvalidType if the token is not of the expected type.
func CheckTokenType[T gr.TokenTyper](tk *gr.Token[T], tk_type T) error {
	if tk == nil {
		return NewErrInvalidType(tk_type, nil)
	} else if tk.Type != tk_type {
		return NewErrInvalidType(tk_type, &tk.Type)
	}

	return nil
}
