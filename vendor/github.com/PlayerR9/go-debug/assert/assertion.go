package assert

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// Assertion is the struct that is used to perform assertions.
type Assertion[T cmp.Ordered] struct {
	// name is the name of the value.
	name string

	// value is the value to assert.
	value T

	// cond is the condition to assert.
	cond Conditioner[T]

	// negative is true if the assertion should be negated.
	negative bool
}

// AssertThat is a constructor for the Assertion struct.
//
// Parameters:
//   - name: the name of the parameter (or variable) to assert.
//   - val: the value of the parameter (or variable) to assert.
//
// Returns:
//   - *Assertion: the assertion. Never returns nil.
//
// A normal construction is a chain of AssertThat() function followed by
// the conditions and the action to perform.
//
// Example:
//
//	foo := "foo"
//	AssertThat("foo", foo).Not().In("bar", "fooo", "baz").Panic()
//	// Does not panic since foo is not in ["bar", "fooo", "baz"]
func AssertThat[T cmp.Ordered](name string, val T) *Assertion[T] {
	return &Assertion[T]{
		name:  name,
		value: val,
	}
}

// Not negates the assertion. Useful for checking the negation of an assertion.
//
// However, if the positive check is more expensive than its negative counterpart,
// it is suggested to create the negative assertion rather than negating the positive one.
//
// Furthermore, if more than one Not() function is called on the same assertion,
// then if the count of the Not() functions is odd, the assertion will be negated. Otherwise,
// the assertion will be positive.
//
// For example, doing .Not().Not().Not() is the same as .Not().
//
// Returns:
//   - *Assertion: the assertion for chaining. Never returns nil.
func (a *Assertion[T]) Not() *Assertion[T] {
	a.negative = !a.negative
	return a
}

// Equal is the assertion for checking if the value is equal to another value.
//
// If any other condition is specified, the furthest condition overwrites any
// other condition.
//
// Parameters:
//   - b: the other value to compare with.
//
// Returns:
//   - *Assertion: the assertion for chaining. Never returns nil.
func (a *Assertion[T]) Equal(b T) *Assertion[T] {
	a.cond = &EqualCond[T]{other: b}

	return a
}

// GreaterThan is the assertion for checking if the value is greater than another value.
//
// If any other condition is specified, the furthest condition overwrites any
// other condition.
//
// Parameters:
//   - b: the other value to compare with.
//
// Returns:
//   - *Assertion: the assertion for chaining. Never returns nil.
func (a *Assertion[T]) GreaterThan(b T) *Assertion[T] {
	a.cond = &GreaterThanCond[T]{other: b}
	return a
}

// LessThan is the assertion for checking if the value is less than another value.
//
// If any other condition is specified, the furthest condition overwrites any
// other condition.
//
// Parameters:
//   - b: the other value to compare with.
//
// Returns:
//   - *Assertion: the assertion for chaining. Never returns nil.
func (a *Assertion[T]) LessThan(b T) *Assertion[T] {
	a.cond = &LessThanCond[T]{other: b}

	return a
}

// GreaterOrEqualThan is the assertion for checking if the value is greater or equal than another value.
//
// If any other condition is specified, the furthest condition overwrites any
// other condition.
//
// Parameters:
//   - b: the other value to compare with.
//
// Returns:
//   - *Assertion: the assertion for chaining. Never returns nil.
func (a *Assertion[T]) GreaterOrEqualThan(b T) *Assertion[T] {
	a.cond = &GreaterOrEqualThanCond[T]{other: b}

	return a
}

// LessOrEqualThan is the assertion for checking if the value is less or equal than another value.
//
// If any other condition is specified, the furthest condition overwrites any
// other condition.
//
// Parameters:
//   - b: the other value to compare with.
//
// Returns:
//   - *Assertion: the assertion for chaining. Never returns nil.
func (a *Assertion[T]) LessOrEqualThan(b T) *Assertion[T] {
	a.cond = &LessOrEqualThanCond[T]{other: b}

	return a
}

// InRange is the assertion for checking if the value is in a range.
//
// If any other condition is specified, the furthest condition overwrites any
// other condition.
//
// Parameters:
//   - min: the minimum value of the range.
//   - max: the maximum value of the range.
//
// Returns:
//   - *Assertion: the assertion for chaining. Never returns nil.
//
// If min is greater than max, the min and max values will be swapped. Moreover, if
// min is equal to max, the assertion will be equal to the EqualCond[T] with the min value.
func (a *Assertion[T]) InRange(min, max T) *Assertion[T] {
	if min > max {
		min, max = max, min
	}

	if min == max {
		a.cond = &EqualCond[T]{other: min}
	} else {
		a.cond = &InRangeCond[T]{min: min, max: max}
	}

	return a
}

// Zero is the assertion for checking if the value is the zero value for its type.
//
// If any other condition is specified, the furthest condition overwrites any
// other condition.
//
// Returns:
//   - *Assertion: the assertion for chaining. Never returns nil.
func (a *Assertion[T]) Zero() *Assertion[T] {
	a.cond = &ZeroCond[T]{}

	return a
}

// In is the assertion for checking if the value is in a list of values.
//
// If any other condition is specified, the furthest condition overwrites any
// other condition.
//
// Parameters:
//   - values: the list of values to check against.
//
// Returns:
//   - *Assertion: the assertion for chaining. Never returns nil.
//
// The list is sorted in ascending order and duplicates are removed. As a special case,
// if only one value is provided, the assertion will be equal to the EqualCond[T] with
// that value.
func (a *Assertion[T]) In(values ...T) *Assertion[T] {
	if len(values) > 2 {
		sorted := make([]T, 0, len(values))

		for _, val := range values {
			pos, ok := slices.BinarySearch(sorted, val)
			if !ok {
				sorted = slices.Insert(sorted, pos, val)
			}
		}

		values = sorted[:len(sorted):len(sorted)]
	}

	switch len(values) {
	case 0:
		a.cond = &InCond[T]{values: []T{}}
	case 1:
		a.cond = &EqualCond[T]{other: values[0]}
	default:
		a.cond = &InCond[T]{values: values}
	}

	return a
}

// Satisfies is the assertion for checking custom conditions.
//
// If any other condition is specified, the furthest condition overwrites any
// other condition. However, if cond is nil, this function will be a no-op.
//
// Parameters:
//   - cond: the condition to check.
//
// Returns:
//   - *Assertion: the assertion for chaining. Never returns nil.
func (a *Assertion[T]) Satisfies(cond Conditioner[T]) *Assertion[T] {
	if cond == nil {
		return a
	}

	a.cond = cond

	return a
}

// Applies is the same as Satisfies but without needing to provide a custom definition
// that implements Conditioner. Best used for checks that are only done once.
//
// If any other condition is specified, the furthest condition overwrites any
// other condition.
//
// Parameters:
//   - msg: the message of the condition.
//   - cond: the condition to check.
//
// Returns:
//   - *Assertion: the assertion for chaining. Never returns nil.
func (a *Assertion[T]) Applies(msg func() string, cond func(value T) bool) *Assertion[T] {
	a.cond = &GenericCond[T]{message: msg, verify: cond}

	return a
}

// Panic will panic if the condition is not met.
//
// The error message is "expected <name> to <message>; got <value> instead" where
// <name> is the name of the assertion, <message> is the message of the condition
// and <value> is the value of the assertion. Finally, this error message is used
// within the *ErrAssertionFailed error.
func (a *Assertion[T]) Panic() {
	if a.cond == nil {
		return
	}

	ok := a.cond.Verify(a.value)
	if ok != a.negative {
		return
	}

	if a.cond.Verify(a.value) {
		return
	}

	var builder strings.Builder

	builder.WriteString("expected ")
	builder.WriteString(strconv.Quote(a.name))

	if a.negative {
		builder.WriteString(" to not ")
	} else {
		builder.WriteString(" to ")
	}

	builder.WriteString(a.cond.Message())
	builder.WriteString("; got ")
	builder.WriteString(strconv.Quote(fmt.Sprintf("%v", a.value)))
	builder.WriteString(" instead")

	panic(NewErrAssertionFailed(builder.String()))
}

// PanicWithMessage is the same as Panic but with a custom error message.
// This error message is overrides the default error message of the Assertion.
//
// Of course, the message is still used within the *ErrAssertionFailed error.
func (a *Assertion[T]) PanicWithMessage(msg string) {
	ok := a.cond.Verify(a.value)
	if ok != a.negative {
		return
	}

	panic(NewErrAssertionFailed(msg))
}

// Error same as Panic but returns the *ErrAssertionFailed error instead of a panic.
//
// The error message is "expected <name> to <message>; got <value> instead" where
// <name> is the name of the assertion, <message> is the message of the condition
// and <value> is the value of the assertion.
//
// Returns:
//   - *ErrAssertionFailed: the error. Nil iff the condition is met.
func (a *Assertion[T]) Error() *ErrAssertionFailed {
	ok := a.cond.Verify(a.value)
	if ok != a.negative {
		return nil
	}

	var builder strings.Builder

	builder.WriteString("expected ")
	builder.WriteString(strconv.Quote(a.name))

	if a.negative {
		builder.WriteString(" to not")
	} else {
		builder.WriteString(" to")
	}

	return NewErrAssertionFailed(fmt.Sprintf("%s %s; got %v instead", builder.String(), a.cond.Message(), a.value))
}

// ErrorWithMessage is the same as PanicWithMessage but returns the *ErrAssertionFailed error instead of a panic.
// This error message is overrides the default error message of the Assertion.
//
// Of course, the message is still used within the *ErrAssertionFailed error.
//
// Returns:
//   - *ErrAssertionFailed: the error. Nil iff the condition is met.
func (a *Assertion[T]) ErrorWithMessage(msg string) error {
	if a.cond.Verify(a.value) {
		return nil
	}

	return NewErrAssertionFailed(msg)
}

// Check checks whether the condition is met.
//
// Returns:
//   - bool: true if the condition is met. false otherwise.
func (a *Assertion[T]) Check() bool {
	ok := a.cond.Verify(a.value)
	return ok != a.negative
}
