package parsing

import (
	"strings"

	gr "github.com/PlayerR9/grammar/grammar"
)

// Rule is a struct that represents a rule of type S.
type Rule[S gr.TokenTyper] struct {
	// lhs is the left-hand side of the rule.
	lhs S

	// rhss is the right-hand side of the rule.
	rhss []S
}

// String implements the fmt.Stringer interface.
//
// Format:
//
//	RHS(n) RHS(n-1) ... RHS(1) -> LHS .
func (r *Rule[S]) String() string {
	var values []string

	for _, rhs := range r.rhss {
		values = append(values, rhs.GoString())
	}

	values = append(values, "->")
	values = append(values, r.lhs.GoString())
	values = append(values, ".")

	return strings.Join(values, " ")
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

// GetLhs returns the left-hand side of the rule.
//
// Returns:
//   - S: The left-hand side of the rule.
func (r *Rule[S]) GetLhs() S {
	return r.lhs
}

// GetIndicesOfRhs returns the ocurrence indices of the rhs in the rule.
//
// Parameters:
//   - rhs: The right-hand side to search.
//
// Returns:
//   - []int: The indices of the rhs in the rule.
func (r *Rule[S]) GetIndicesOfRhs(rhs S) []int {
	var indices []int

	for i := 0; i < len(r.rhss); i++ {
		if r.rhss[i] == rhs {
			indices = append(indices, i)
		}
	}

	return indices
}

// GetRhss returns the right-hand sides of the rule.
//
// Returns:
//   - []S: The right-hand sides of the rule.
func (r *Rule[S]) GetRhss() []S {
	return r.rhss
}

// Size returns the number of right-hand sides of the rule.
//
// Returns:
//   - int: The "size" of the rule.
func (r *Rule[S]) Size() int {
	return len(r.rhss)
}
