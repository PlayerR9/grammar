package parser

import (
	"fmt"
	"strings"

	gr "github.com/PlayerR9/grammar/grammar"
)

// ErrAfter is an error that occurs after a certain type.
type ErrAfter[T gr.Enumer] struct {
	// Type is the type before the one that caused the error.
	Type T

	// Err is the underlying error.
	Err error
}

// Error implements the error interface.
//
// Message: "after <type>: <error>" or "something went wrong after <type>"
func (e ErrAfter[T]) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("something went wrong after %q", e.Type.String())
	} else {
		return fmt.Sprintf("after %q: %s", e.Type.String(), e.Err.Error())
	}
}

// Unwrap implements the error interface.
func (e ErrAfter[T]) Unwrap() error {
	return e.Err
}

// NewErrAfter creates a new ErrAfter error.
//
// Parameters:
//   - type_: The type before the one that caused the error.
//   - err: The underlying error.
//
// Returns:
//   - *ErrAfter: The new error. Never returns nil.
func NewErrAfter[T gr.Enumer](type_ T, err error) *ErrAfter[T] {
	return &ErrAfter[T]{
		Type: type_,
		Err:  err,
	}
}

// ErrBefore is an error that occurs before a certain type.
type ErrBefore[T gr.Enumer] struct {
	// Type is the type after the one that caused the error.
	Type T

	// Err is the underlying error.
	Err error
}

// Error implements the error interface.
//
// Message: "before <type>: <error>" or "something went wrong before <type>"
func (e ErrBefore[T]) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("something went wrong before %q", e.Type.String())
	} else {
		return fmt.Sprintf("before %q: %s", e.Type.String(), e.Err.Error())
	}
}

// Unwrap implements the error interface.
func (e ErrBefore[T]) Unwrap() error {
	return e.Err
}

// NewErrBefore creates a new ErrBefore error.
//
// Parameters:
//   - type_: The type after the one that caused the error.
//   - err: The underlying error.
//
// Returns:
//   - *ErrBefore: The new error. Never returns nil.
func NewErrBefore[T gr.Enumer](type_ T, err error) *ErrBefore[T] {
	return &ErrBefore[T]{
		Type: type_,
		Err:  err,
	}
}

// ErrUnexpectedToken is an error that occurs when an unexpected token is found.
type ErrUnexpectedToken[T gr.Enumer] struct {
	// Left is the expected type.
	Left T

	// Right is the unexpected type.
	Right T

	// Got is the token that was found.
	Got *T
}

// Error implements the error interface.
//
// Message: "expected either <left> or <right> but got <got> instead"
func (e ErrUnexpectedToken[T]) Error() string {
	var builder strings.Builder

	builder.WriteString("expected either")
	builder.WriteString(e.Left.String())
	builder.WriteString(" or ")
	builder.WriteString(e.Right.String())
	builder.WriteString(" but got ")

	if e.Got == nil {
		builder.WriteString("nothing")
	} else {
		builder.WriteString((*e.Got).String())
	}

	builder.WriteString(" instead")

	return builder.String()
}

// NewErrUnexpectedToken creates a new ErrUnexpectedToken error.
//
// Parameters:
//   - left: The expected type.
//   - right: The unexpected type.
//   - got: The token that was found.
//
// Returns:
//   - *ErrUnexpectedToken: The new error. Never returns nil.
func NewErrUnexpectedToken[T gr.Enumer](left, right T, got *T) *ErrUnexpectedToken[T] {
	return &ErrUnexpectedToken[T]{
		Left:  left,
		Right: right,
		Got:   got,
	}
}
