package assert

import (
	"cmp"
	"fmt"
	"strings"
)

// Conditioner is an interface that describes the behavior of the condition
// for the Assertion struct.
//
// Since all messages are preceded by the "expected <value> to", it is important
// to make sure that the grammar is correct. Thus, a message of "be true" is shown
// as "expected <value> to be true; got <value> instead".
type Conditioner[T cmp.Ordered] interface {
	// Message returns the message that is shown when the condition is not met.
	//
	// Returns:
	//   - string: the message.
	Message() string

	// Verify checks if the condition is met. In other words, whether the value
	// satisfies the condition.
	//
	// Parameters:
	//   - value: the value to check.
	//
	// Returns:
	//   - bool: true if the condition is met. false otherwise.
	Verify(value T) bool
}

// EqualCond is the condition for the Assertion struct and checks if the
// value is equal to another value.
type EqualCond[T cmp.Ordered] struct {
	// other is the other value to compare with.
	other T
}

// Message implements the Conditioner interface.
//
// Message: "have the value <other>"
func (ec *EqualCond[T]) Message() string {
	return fmt.Sprintf("have the value %v", ec.other)
}

// Verify implements the Conditioner interface.
func (ec *EqualCond[T]) Verify(value T) bool {
	return ec.other == value
}

// GreaterThanCond is the condition for the Assertion struct and checks if the
// value is greater than another value.
type GreaterThanCond[T cmp.Ordered] struct {
	// other is the other value to compare with.
	other T
}

// Message implements the Conditioner interface.
//
// Message: "be greater than <other>"
func (gtc *GreaterThanCond[T]) Message() string {
	return fmt.Sprintf("be greater than %v", gtc.other)
}

// Verify implements the Conditioner interface.
func (gtc *GreaterThanCond[T]) Verify(value T) bool {
	return value > gtc.other
}

// LessThanCond is the condition for the Assertion struct and checks if the
// value is less than another value.
type LessThanCond[T cmp.Ordered] struct {
	// other is the other value to compare with.
	other T
}

// Message implements the Conditioner interface.
//
// Message: "be less than <other>"
func (ltc *LessThanCond[T]) Message() string {
	return fmt.Sprintf("be less than %v", ltc.other)
}

// Verify implements the Conditioner interface.
func (ltc *LessThanCond[T]) Verify(value T) bool {
	return value < ltc.other
}

// GreaterOrEqualThanCond is the condition for the Assertion struct and checks if the
// value is greater or equal than another value.
type GreaterOrEqualThanCond[T cmp.Ordered] struct {
	// other is the other value to compare with.
	other T
}

// Message implements the Conditioner interface.
//
// Message: "be greater or equal than <other>"
func (gtec *GreaterOrEqualThanCond[T]) Message() string {
	return fmt.Sprintf("be greater or equal than %v", gtec.other)
}

// Verify implements the Conditioner interface.
func (gtec *GreaterOrEqualThanCond[T]) Verify(value T) bool {
	return value >= gtec.other
}

// LessOrEqualThanCond is the condition for the Assertion struct and checks if the
// value is less or equal than another value.
type LessOrEqualThanCond[T cmp.Ordered] struct {
	// other is the other value to compare with.
	other T
}

// Message implements the Conditioner interface.
//
// Message: "be less or equal than <other>"
func (letc *LessOrEqualThanCond[T]) Message() string {
	return fmt.Sprintf("be less or equal than %v", letc.other)
}

// Verify implements the Conditioner interface.
func (letc *LessOrEqualThanCond[T]) Verify(value T) bool {
	return value <= letc.other
}

// InRangeCond is the condition for the Assertion struct and checks if the
// value is in a range. Both bounds are included in the range.
type InRangeCond[T cmp.Ordered] struct {
	// min is the lower bound of the range.
	min T

	// max is the upper bound of the range.
	max T
}

// Message implements the Conditioner interface.
//
// Message: "be in range [<min> : <max>]"
func (irc *InRangeCond[T]) Message() string {
	return fmt.Sprintf("be in range [%v : %v]", irc.min, irc.max)
}

// Verify implements the Conditioner interface.
func (irc *InRangeCond[T]) Verify(value T) bool {
	return irc.min <= value && irc.max >= value
}

// ZeroCond is the condition for the Assertion struct and checks if the
// value is the zero value for its type. For example, 0 for int, 0.0 for
// float, "" for string, etc.
type ZeroCond[T cmp.Ordered] struct{}

// Message implements the Conditioner interface.
//
// Message: "be its zero value"
func (izc *ZeroCond[T]) Message() string {
	return "be its zero value"
}

// Verify implements the Conditioner interface.
func (izc *ZeroCond[T]) Verify(value T) bool {
	return value == *new(T)
}

// InCond is the condition for the Assertion struct and checks if the
// value is in a list of values.
type InCond[T cmp.Ordered] struct {
	// values is the list of values to check for.
	values []T
}

// Message implements the Conditioner interface.
//
// Message: "be one of {<value1>, <value2>, ...}"
func (iic *InCond[T]) Message() string {
	if len(iic.values) == 0 {
		return "be one of {}"
	}

	var builder strings.Builder

	builder.WriteString("be one of {")
	fmt.Fprintf(&builder, "%v", iic.values[0])

	for i := 1; i < len(iic.values); i++ {
		fmt.Fprintf(&builder, ", %v", iic.values[i])
	}

	builder.WriteRune('}')

	return builder.String()
}

// Verify implements the Conditioner interface.
func (iic *InCond[T]) Verify(value T) bool {
	if len(iic.values) == 0 {
		return false
	}

	for _, v := range iic.values {
		if v == value {
			return true
		}
	}

	return false
}

// GenericCond is the condition for the Assertion struct and is used for
// specifying custom conditions without having to write a custom Conditioner.
type GenericCond[T any] struct {
	// message is the function that returns the message of the condition.
	message func() string

	// verify is the function that returns the verification of the condition.
	verify func(T) bool
}

// Message implements the Conditioner interface.
//
// The returned message is the result of the message function and, if no message function
// was provided, an empty string is returned.
func (gc *GenericCond[T]) Message() string {
	if gc.message == nil {
		return ""
	}

	return gc.message()
}

// Verify implements the Conditioner interface.
//
// The returned verification is the result of the verify function and, if no verify
// function was provided, true is returned.
func (gc *GenericCond[T]) Verify(value T) bool {
	if gc.verify == nil {
		return true
	}

	return gc.verify(value)
}
