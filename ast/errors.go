package ast

import (
	"strings"

	gcers "github.com/PlayerR9/go-commons/errors"
	internal "github.com/PlayerR9/grammar/internal"
)

// ErrIn is an error that occurs when an error is encountered in a rule.
type ErrIn[T internal.TokenTyper] struct {
	// Type is the type in which the error occurred.
	Type T

	// Reason is the reason for the error.
	Reason error
}

// Error implements the error interface.
//
// Message: "in <type>: <reason>"
func (e ErrIn[T]) Error() string {
	var builder strings.Builder

	builder.WriteString("in ")
	builder.WriteString(e.Type.String())
	builder.WriteString(": ")
	builder.WriteString(gcers.Error(e.Reason))

	return builder.String()
}

// Unwrap returns the reason for the error.
//
// Returns:
//   - error: The reason for the error.
func (e ErrIn[T]) Unwrap() error {
	return e.Reason
}

// NewErrIn creates a new error that occurs when an error is encountered in a
// rule.
//
// Parameters:
//   - type_: The type in which the error occurred.
//   - reason: The reason for the error.
//
// Returns:
//   - *ErrIn[T]: An error if the type is invalid. Never returns nil.
func NewErrIn[T internal.TokenTyper](type_ T, reason error) *ErrIn[T] {
	return &ErrIn[T]{
		Type:   type_,
		Reason: reason,
	}
}

// ChangeReason changes the reason of the error.
//
// Parameters:
//   - reason: The new reason of the error.
func (e *ErrIn[T]) ChangeReason(reason error) {
	e.Reason = reason
}
