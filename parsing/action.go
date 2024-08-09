package parsing

import (
	"fmt"

	gcers "github.com/PlayerR9/go-commons/errors"
	gr "github.com/PlayerR9/grammar/grammar"
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
type ReduceAction[S gr.TokenTyper] struct {
	// rule is the rule to reduce.
	rule *Rule[S]
}

// String implements the Actioner interface.
func (r *ReduceAction[S]) String() string {
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
func NewReduceAction[S gr.TokenTyper](rule *Rule[S]) (*ReduceAction[S], error) {
	if rule == nil {
		return nil, gcers.NewErrNilParameter("rule")
	}

	return &ReduceAction[S]{
		rule: rule,
	}, nil
}

// AcceptAction is the accept action.
type AcceptAction[S gr.TokenTyper] struct {
	// rule is the rule to accept.
	rule *Rule[S]
}

// String implements the Actioner interface.
func (a *AcceptAction[S]) String() string {
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
func NewAcceptAction[S gr.TokenTyper](rule *Rule[S]) (*AcceptAction[S], error) {
	if rule == nil {
		return nil, gcers.NewErrNilParameter("rule")
	}

	return &AcceptAction[S]{
		rule: rule,
	}, nil
}
