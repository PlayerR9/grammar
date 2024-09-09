package parser

import (
	"iter"

	gcers "github.com/PlayerR9/go-commons/errors"
	gr "github.com/PlayerR9/grammar/grammar"
)

// Rule represents a rule in the grammar.
type Rule[T gr.Enumer] struct {
	// lhs is the left hand side of the rule.
	lhs T

	// rhss is the right hand side of the rule.
	rhss []T
}

// NewRule creates a new rule.
//
// Parameters:
//   - lhs: The left hand side of the rule.
//   - rhss: The right hand side of the rule.
//
// Returns:
//   - *Rule: The new rule.
//   - error: An error if rhss is empty.
func NewRule[T gr.Enumer](lhs T, rhss ...T) (*Rule[T], error) {
	if len(rhss) == 0 {
		return nil, gcers.NewErrInvalidParameter("rhss", gcers.NewErrEmpty(rhss))
	}

	return &Rule[T]{
		lhs:  lhs,
		rhss: rhss,
	}, nil
}

// BackwardRhs returns the right hand side of the rule in reverse order.
//
// Returns:
//   - iter.Seq[T]: The right hand side of the rule in reverse order.
func (r Rule[T]) BackwardRhs() iter.Seq[T] {
	fn := func(yield func(T) bool) {
		for i := len(r.rhss) - 1; i >= 0; i-- {
			if !yield(r.rhss[i]) {
				break
			}
		}
	}

	return fn
}

// Lhs returns the left hand side of the rule.
//
// Returns:
//   - T: The left hand side of the rule.
func (r Rule[T]) Lhs() T {
	return r.lhs
}
