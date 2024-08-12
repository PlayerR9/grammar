package ast

import (
	"strconv"
	"strings"

	gr "github.com/PlayerR9/grammar/grammar"
)

// ErrInvalidType is an error that occurs when an invalid type is encountered.
type ErrInvalidType[T gr.TokenTyper] struct {
	// Expected is the expected type.
	Expected T

	// Got is the actual type.
	Got *T
}

// Error implements the error interface.
//
// Message: "expected <expected>, got <actual> instead"
func (e ErrInvalidType[T]) Error() string {
	var builder strings.Builder

	builder.WriteString("expected ")
	builder.WriteString(strconv.Quote(e.Expected.String()))
	builder.WriteString(", got ")

	if e.Got == nil {
		builder.WriteString("nothing")
	} else {
		builder.WriteString(strconv.Quote((*e.Got).String()))
	}

	builder.WriteString(" instead")

	return builder.String()
}

// NewErrInvalidType creates a new error that occurs when an invalid type
// is encountered.
//
// Parameters:
//   - expected: The expected type.
//   - got: The actual type.
//
// Returns:
//   - error: An error if the type is invalid. Never returns nil.
func NewErrInvalidType[T gr.TokenTyper](expected T, got *T) *ErrInvalidType[T] {
	return &ErrInvalidType[T]{
		Expected: expected,
		Got:      got,
	}
}

// ErrInRule is an error that occurs when an error is encountered in a rule.
type ErrInRule[T gr.TokenTyper] struct {
	// Lhs is the left-hand side of the rule.
	Lhs T

	// Reason is the reason for the error.
	Reason error
}

// Error implements the error interface.
//
// Format:
//
//	"in rule <lhs>: <reason>" or "an error occurred in the rule <lhs>"
func (e ErrInRule[T]) Error() string {
	var builder strings.Builder

	if e.Reason == nil {
		builder.WriteString("an error occurred in the rule ")
		builder.WriteString(strconv.Quote(e.Lhs.String()))
	} else {
		builder.WriteString("in rule ")
		builder.WriteString(strconv.Quote(e.Lhs.String()))
		builder.WriteString(": ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

// Unwrap returns the reason of the error.
//
// Returns:
//   - error: The reason of the error.
func (e ErrInRule[T]) Unwrap() error {
	return e.Reason
}

// NewErrInRule creates a new error that occurs when an error is encountered
// in a rule.
//
// Parameters:
//   - lhs: The left-hand side of the rule.
//   - reason: The reason for the error.
//
// Returns:
//   - error: An error if the rule is invalid. Never returns nil.
func NewErrInRule[T gr.TokenTyper](lhs T, reason error) *ErrInRule[T] {
	return &ErrInRule[T]{
		Lhs:    lhs,
		Reason: reason,
	}
}

// ChangeReason changes the reason of the error.
//
// Parameters:
//   - reason: The new reason of the error.
func (e *ErrInRule[T]) ChangeReason(reason error) {
	e.Reason = reason
}
