package errors

import (
	"cmp"
	"errors"
	"fmt"
	"reflect"
)

var (
	// NilValue is the error returned when a pointer is nil. While readers are not expected to return this
	// error by itself, if it does, readers must not wrap it as callers will test this error using ==.
	NilValue error
)

func init() {
	NilValue = errors.New("pointer must not be nil")
}

// ErrEmpty represents an error when a value is empty.
type ErrEmpty struct {
	// Type is the type of the empty value.
	Type any
}

// Error implements the error interface.
//
// Message: "{{ .Type }} must not be empty"
func (e *ErrEmpty) Error() string {
	var t_string string

	if e.Type == nil {
		t_string = "nil"
	} else {
		to := reflect.TypeOf(e.Type)
		t_string = to.String()
	}

	return t_string + " must not be empty"
}

// NewErrEmpty creates a new ErrEmpty error.
//
// Parameters:
//   - var_type: The type of the empty value.
//
// Returns:
//   - *ErrEmpty: A pointer to the newly created ErrEmpty. Never returns nil.
func NewErrEmpty(var_type any) *ErrEmpty {
	return &ErrEmpty{
		Type: var_type,
	}
}

// ErrGT represents an error when a value is less than or equal to a specified value.
type ErrGT[T cmp.Ordered] struct {
	// Value is the value that caused the error.
	Value T
}

// Error implements the error interface.
//
// Message: "value must be greater than <value>"
func (e *ErrGT[T]) Error() string {
	return fmt.Sprintf("value must ge greater than %v", e.Value)
}

// NewErrGT creates a new ErrGT error with the specified value.
//
// Parameters:
//   - value: The minimum value that is not allowed.
//
// Returns:
//   - *ErrGT: A pointer to the newly created ErrGT.
func NewErrGT[T cmp.Ordered](value T) *ErrGT[T] {
	e := &ErrGT[T]{
		Value: value,
	}
	return e
}

// ErrLT represents an error when a value is greater than or equal to a specified value.
type ErrLT[T cmp.Ordered] struct {
	// Value is the value that caused the error.
	Value T
}

// Error implements the error interface.
//
// Message: "value must be less than <value>"
func (e *ErrLT[T]) Error() string {
	return fmt.Sprintf("value must be less than %v", e.Value)
}

// NewErrLT creates a new ErrLT error with the specified value.
//
// Parameters:
//   - value: The maximum value that is not allowed.
//
// Returns:
//   - *ErrLT: A pointer to the newly created ErrLT.
func NewErrLT[T cmp.Ordered](value T) *ErrLT[T] {
	e := &ErrLT[T]{
		Value: value,
	}
	return e
}

// ErrGTE represents an error when a value is less than a specified value.
type ErrGTE[T cmp.Ordered] struct {
	// Value is the value that caused the error.
	Value T
}

// Error implements the error interface.
//
// Message: "value must be greater than or equal to <value>"
func (e *ErrGTE[T]) Error() string {
	return fmt.Sprintf("value must be greater than or equal to %v", e.Value)
}

// NewErrGTE creates a new ErrGTE error with the specified value.
//
// Parameters:
//   - value: The minimum value that is allowed.
//
// Returns:
//   - *ErrGTE: A pointer to the newly created ErrGTE.
func NewErrGTE[T cmp.Ordered](value T) *ErrGTE[T] {
	e := &ErrGTE[T]{
		Value: value,
	}
	return e
}

// ErrLTE represents an error when a value is greater than a specified value.
type ErrLTE[T cmp.Ordered] struct {
	// Value is the value that caused the error.
	Value T
}

// Error implements the error interface.
//
// Message: "value must be less than or equal to <value>"
func (e *ErrLTE[T]) Error() string {
	return fmt.Sprintf("value must be less than or equal to %v", e.Value)
}

// NewErrLTE creates a new ErrLTE error with the specified value.
//
// Parameters:
//   - value: The maximum value that is allowed.
//
// Returns:
//   - *ErrLTE: A pointer to the newly created ErrLTE.
func NewErrLTE[T cmp.Ordered](value T) *ErrLTE[T] {
	e := &ErrLTE[T]{
		Value: value,
	}
	return e
}
