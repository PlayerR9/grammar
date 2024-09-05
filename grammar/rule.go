package grammar

import (
	"iter"
	"slices"

	utst "github.com/PlayerR9/go-commons/cmp"
	gcers "github.com/PlayerR9/go-commons/errors"
	internal "github.com/PlayerR9/grammar/grammar/internal"
)

// Rule is a grammar rule.
type Rule[T internal.TokenTyper] struct {
	// lhs is the left-hand side of the rule.
	lhs T

	// rhss is the right-hand side of the rule.
	rhss []T
}

// NewRule creates a new rule with the given left-hand side and right-hand side.
//
// Parameters:
//   - lhs: The left-hand side of the rule.
//   - rhss: The right-hand side of the rule.
//
// Returns:
//   - *Rule[T]: The new rule.
//   - error: An error of type *errors.ErrInvalidParameter if 'rhss' is empty.
func NewRule[T internal.TokenTyper](lhs T, rhss []T) (*Rule[T], error) {
	if len(rhss) == 0 {
		return nil, gcers.NewErrInvalidParameter("rhss", gcers.NewErrEmpty(rhss))
	}

	return &Rule[T]{
		lhs:  lhs,
		rhss: rhss,
	}, nil
}

// Equals checks if the given rule is equal to the current rule.
//
// Parameters:
//   - other: The other rule.
//
// Returns:
//   - bool: True if the given rule is equal to the current rule, otherwise false.
//
// If 'other' is nil, it always returns false.
func (r *Rule[T]) Equals(other *Rule[T]) bool {
	return other != nil && r.lhs == other.lhs && slices.Equal(r.rhss, other.rhss)
}

// Size returns the amount of right-hand sides of the rule.
//
// Returns:
//   - int: The amount of right-hand sides of the rule.
//
// This function never returns 0.
func (r Rule[T]) Size() int {
	return len(r.rhss)
}

// IndicesOf returns a slice of indices of the given right-hand side in the rule.
//
// Parameters:
//   - rhs: The right-hand side of the rule.
//
// Returns:
//   - []int: The slice of indices.
func (r Rule[T]) IndicesOf(rhs T) []int {
	var indices []int

	for i, tmp := range r.rhss {
		if tmp == rhs {
			indices = append(indices, i)
		}
	}

	return indices
}

// RhsAt returns the right-hand side at the given index.
//
// Parameters:
//   - idx: The index of the right-hand side.
//
// Returns:
//   - T: The right-hand side.
//   - bool: True if the right-hand side exists, otherwise false.
func (r Rule[T]) RhsAt(idx int) (T, bool) {
	if idx < 0 || idx >= len(r.rhss) {
		return T(0), false
	}

	return r.rhss[idx], true
}

// Symbols returns the set of symbols in the rule.
//
// Returns:
//   - *utst.Set[T]: The set of symbols. Never returns nil.
func (r Rule[T]) Symbols() *utst.Set[T] {
	symbols := utst.NewSet[T]()

	symbols.Add(r.lhs)

	for _, rhs := range r.rhss {
		symbols.Add(rhs)
	}

	return symbols
}

// Rhs returns an iterator over the right-hand side of the rule.
//
// Returns:
//   - iter.Seq[T]: The iterator. Never returns nil.
func (r Rule[T]) Rhs() iter.Seq[T] {
	fn := func(yield func(T) bool) {
		for _, rhs := range r.rhss {
			if !yield(rhs) {
				return
			}
		}
	}

	return fn
}

// Backwards returns an iterator over the right-hand side of the rule in reverse order.
//
// Returns:
//   - iter.Seq[T]: The iterator. Never returns nil.
func (r Rule[T]) Backwards() iter.Seq[T] {
	rhss := make([]T, len(r.rhss))
	copy(rhss, r.rhss)

	slices.Reverse(rhss)

	fn := func(yield func(T) bool) {
		for i := 0; i < len(rhss); i++ {
			if !yield(rhss[i]) {
				return
			}
		}
	}

	return fn
}

// Lhs returns the left-hand side of the rule.
//
// Returns:
//   - T: The left-hand side.
func (r Rule[T]) Lhs() T {
	return r.lhs
}

// ExtractRhsAt returns a slice of the right-hand side at the given index.
//
// Parameters:
//   - at: The index of the right-hand side.
//
// Returns:
//   - []T: The slice of the right-hand side.
func ExtractRhsAt[T internal.TokenTyper](rules []*Rule[T], at int) []T {
	if len(rules) == 0 || at < 0 {
		return nil
	}

	var rhss []T

	for _, r := range rules {
		if r == nil {
			continue
		}

		rhs, ok := r.RhsAt(at)
		if !ok {
			continue
		}

		pos, ok := slices.BinarySearch(rhss, rhs)
		if !ok {
			rhss = slices.Insert(rhss, pos, rhs)
		}
	}

	return rhss
}
