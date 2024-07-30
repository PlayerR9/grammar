package parser

import (
	"fmt"

	gr "github.com/PlayerR9/grammar/grammar"
	luc "github.com/PlayerR9/lib_units/common"
)

// Actioner is an interface that defines the behavior of an action.
type Actioner interface {
	fmt.Stringer
}

// ShiftAction is the shift action.
type ShiftAction struct {
}

// String implements the Actioner interface.
func (s *ShiftAction) String() string {
	return "Shift"
}

// NewShiftAction creates a new shift action.
//
// Returns:
//   - *ShiftAction: The new shift action. Never returns nil.
func NewShiftAction() *ShiftAction {
	return &ShiftAction{}
}

// ReduceAction is the reduce action.
type ReduceAction[T gr.TokenTyper] struct {
	// rule is the rule to reduce.
	rule *Rule[T]
}

// String implements the Actioner interface.
func (r *ReduceAction[T]) String() string {
	return "Reduce"
}

// NewReduceAction creates a new reduce action.
//
// Parameters:
//   - rule: The rule to reduce.
//
// Returns:
//   - *ReduceAction: The new reduce action.
//   - error: An error of type *common.ErrInvalidParameter if the rule is nil.
func NewReduceAction[T gr.TokenTyper](rule *Rule[T]) (*ReduceAction[T], error) {
	if rule == nil {
		return nil, luc.NewErrNilParameter("rule")
	}

	return &ReduceAction[T]{
		rule: rule,
	}, nil
}

// AcceptAction is the accept action.
type AcceptAction[T gr.TokenTyper] struct {
	// rule is the rule to accept.
	rule *Rule[T]
}

// String implements the Actioner interface.
func (a *AcceptAction[T]) String() string {
	return "Accept"
}

// NewAcceptAction creates a new accept action.
//
// Parameters:
//   - rule: The rule to accept.
//
// Returns:
//   - *AcceptAction: The new accept action.
//   - error: An error of type *common.ErrInvalidParameter if the rule is nil.
func NewAcceptAction[T gr.TokenTyper](rule *Rule[T]) (*AcceptAction[T], error) {
	if rule == nil {
		return nil, luc.NewErrNilParameter("rule")
	}

	return &AcceptAction[T]{
		rule: rule,
	}, nil
}
