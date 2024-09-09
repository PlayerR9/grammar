package parser

import (
	gcers "github.com/PlayerR9/go-commons/errors"
	gr "github.com/PlayerR9/grammar/grammar"
)

// Actioner is an interface for actions.
type Actioner interface {
}

// ShiftAct is a shift action.
type ShiftAct struct {
}

// NewShiftAct creates a new shift action.
//
// Returns:
//   - *ShiftAct: The new shift action. Never returns nil.
func NewShiftAct() *ShiftAct {
	return &ShiftAct{}
}

// ReduceAct is a reduce action.
type ReduceAct[T gr.Enumer] struct {
	// rule is the rule that is being reduced.
	rule *Rule[T]
}

// NewReduceAct creates a new reduce action.
//
// Parameters:
//   - rule: The rule that is being reduced.
//
// Returns:
//   - *ReduceAct: The new reduce action.
//   - error: An error if rule is nil.
func NewReduceAct[T gr.Enumer](rule *Rule[T]) (*ReduceAct[T], error) {
	if rule == nil {
		return nil, gcers.NewErrNilParameter("rule")
	}

	return &ReduceAct[T]{
		rule: rule,
	}, nil
}

// Rule returns the rule that is being reduced.
//
// Returns:
//   - *Rule: The rule that is being reduced. Never returns nil.
func (a ReduceAct[T]) Rule() *Rule[T] {
	return a.rule
}

// AcceptAct is an accept action.
type AcceptAct[T gr.Enumer] struct {
	// rule is the rule that is being accepted.
	rule *Rule[T]
}

// NewAcceptAct creates a new accept action.
//
// Parameters:
//   - rule: The rule that is being accepted.
//
// Returns:
//   - *AcceptAct: The new accept action.
//   - error: An error if rule is nil.
func NewAcceptAct[T gr.Enumer](rule *Rule[T]) (*AcceptAct[T], error) {
	if rule == nil {
		return nil, gcers.NewErrNilParameter("rule")
	}

	return &AcceptAct[T]{
		rule: rule,
	}, nil
}

// Rule returns the rule that is being accepted.
//
// Returns:
//   - *Rule: The rule that is being accepted. Never returns nil.
func (a AcceptAct[T]) Rule() *Rule[T] {
	return a.rule
}
