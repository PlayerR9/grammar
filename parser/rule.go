package parser

import (
	gr "github.com/PlayerR9/grammar/grammar"
	luc "github.com/PlayerR9/lib_units/common"
)

// Rule is a struct that represents a rule of type S.
type Rule[S gr.TokenTyper] struct {
	// lhs is the left-hand side of the rule.
	lhs S

	// rhss is the right-hand side of the rule.
	rhss []S
}

// Iterator implements the Ruler interface.
func (r *Rule[S]) Iterator() luc.Iterater[S] {
	return luc.NewSimpleIterator(r.rhss)
}

// NewRule creates a new rule.
//
// Parameters:
//   - lhs: The left-hand side of the rule.
//   - rhss: The right-hand side of the rule.
//
// Returns:
//   - *Rule[S]: The new rule.
//
// Returns nil if the rhss is empty.
func NewRule[S gr.TokenTyper](lhs S, rhss []S) *Rule[S] {
	if len(rhss) == 0 {
		return nil
	}

	return &Rule[S]{
		lhs:  lhs,
		rhss: rhss,
	}
}
