package grammar

import (
	"strconv"
	"strings"

	gcstr "github.com/PlayerR9/go-commons/strings"
	internal "github.com/PlayerR9/grammar/grammar/internal"
)

// ErrUnexpectedToken is the error for unexpected tokens.
type ErrUnexpectedToken[T internal.TokenTyper] struct {
	// Expected is the expected token.
	Expecteds []T

	// Prev is the previous token.
	Prev *T

	// Got is the actual token.
	Got *T
}

// Error implements the error interface.
//
// Message: "expected {expected} after {prev}, got {got} instead".
func (e ErrUnexpectedToken[T]) Error() string {
	var builder strings.Builder

	builder.WriteString("expected ")

	if len(e.Expecteds) == 0 {
		builder.WriteString("nothing")
	} else {
		values := make([]string, 0, len(e.Expecteds))
		for _, expected := range e.Expecteds {
			values = append(values, expected.String())
		}
		gcstr.QuoteStrings(values)
		builder.WriteString(gcstr.EitherOrString(values))
	}

	if e.Prev == nil {
		builder.WriteString(" at the start of the input")
	} else {
		builder.WriteString(" after ")
		builder.WriteString(strconv.Quote((*e.Prev).String()))
	}

	builder.WriteString(", got ")

	if e.Got == nil {
		builder.WriteString("nothing")
	} else {
		builder.WriteString(strconv.Quote((*e.Got).String()))
	}

	builder.WriteString(" instead")

	return builder.String()
}

// NewErrUnexpectedToken creates a new ErrUnexpectedToken.
//
// Parameters:
//   - prev: The previous token.
//   - got: The actual token.
//   - expecteds: The expected tokens.
//
// Returns:
//   - *ErrUnexpectedToken[T]: A pointer to the new ErrUnexpectedToken. Never returns nil.
func NewErrUnexpectedToken[T internal.TokenTyper](prev, got *T, expecteds ...T) *ErrUnexpectedToken[T] {
	return &ErrUnexpectedToken[T]{
		Expecteds: expecteds,
		Prev:      prev,
		Got:       got,
	}
}

// ErrUnexpectedLookahead is the error for unexpected tokens.
type ErrUnexpectedLookahead[T internal.TokenTyper] struct {
	// Expected is the expected token.
	Expecteds []T

	// Prev is the previous token.
	Prev T

	// Got is the actual token.
	Got *T
}

// Error implements the error interface.
//
// Message: "expected {expected} after {prev}, got {got} instead".
func (e ErrUnexpectedLookahead[T]) Error() string {
	var builder strings.Builder

	builder.WriteString("expected ")

	if len(e.Expecteds) == 0 {
		builder.WriteString("nothing")
	} else {
		values := make([]string, 0, len(e.Expecteds))
		for _, expected := range e.Expecteds {
			values = append(values, expected.String())
		}
		gcstr.QuoteStrings(values)
		builder.WriteString(gcstr.EitherOrString(values))
	}

	builder.WriteString(" after ")
	builder.WriteString(strconv.Quote(e.Prev.String()))
	builder.WriteString(", got ")

	if e.Got == nil {
		builder.WriteString("nothing")
	} else {
		builder.WriteString(strconv.Quote((*e.Got).String()))
	}

	builder.WriteString(" instead")

	return builder.String()
}

// NewErrUnexpectedLookahead creates a new ErrUnexpectedLookahead.
//
// Parameters:
//   - prev: The previous token.
//   - got: The actual token.
//   - expecteds: The expected tokens.
//
// Returns:
//   - *ErrUnexpectedLookahead[T]: A pointer to the new ErrUnexpectedLookahead. Never returns nil.
func NewErrUnexpectedLookahead[T internal.TokenTyper](prev T, got *T, expecteds ...T) *ErrUnexpectedLookahead[T] {
	return &ErrUnexpectedLookahead[T]{
		Expecteds: expecteds,
		Prev:      prev,
		Got:       got,
	}
}

type ErrIn[T internal.TokenTyper] struct {
	Type   T
	Reason error
}

func (e ErrIn[T]) Error() string {
	var builder strings.Builder

	builder.WriteString("in ")
	builder.WriteString(e.Type.String())

	if e.Reason == nil {
		builder.WriteString(": something went wrong")
	} else {
		builder.WriteString(": ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

func (e ErrIn[T]) Unwrap() error {
	return e.Reason
}

func NewErrIn[T internal.TokenTyper](type_ T, reason error) *ErrIn[T] {
	return &ErrIn[T]{
		Type:   type_,
		Reason: reason,
	}
}

func (e *ErrIn[T]) ChangeReason(reason error) {
	e.Reason = reason
}

/* // ErrUnexpectedNode is the error for unexpected tokens.
type ErrUnexpectedNode[T internal.TokenTyper] struct {
	// Expected is the expected token.
	Expecteds []T

	// Prev is the previous token.
	Prev *T

	// Got is the actual token.
	Got *T
}

// Error implements the error interface.
//
// Message: "expected {expected} after {prev}, got {got} instead".
func (e ErrUnexpectedNode[N]) Error() string {
	var builder strings.Builder

	builder.WriteString("expected ")

	if len(e.Expecteds) == 0 {
		builder.WriteString("nothing")
	} else {
		values := make([]string, 0, len(e.Expecteds))
		for _, expected := range e.Expecteds {
			values = append(values, expected.String())
		}
		gcstr.QuoteStrings(values)
		builder.WriteString(gcstr.EitherOrString(values))
	}

	if e.Prev == nil {
		builder.WriteString(" at the start of the input")
	} else {
		builder.WriteString(" after ")
		builder.WriteString(strconv.Quote((*e.Prev).String()))
	}

	builder.WriteString(", got ")

	if e.Got == nil {
		builder.WriteString("nothing")
	} else {
		builder.WriteString(strconv.Quote((*e.Got).String()))
	}

	builder.WriteString(" instead")

	return builder.String()
}

// NewErrUnexpectedNode creates a new ErrUnexpectedNode.
//
// Parameters:
//   - prev: The previous token.
//   - got: The actual token.
//   - expecteds: The expected tokens.
//
// Returns:
//   - *ErrUnexpectedNode: A pointer to the new ErrUnexpectedNode. Never returns nil.
func NewErrUnexpectedNode[N Noder](prev, got *N, expecteds ...N) *ErrUnexpectedNode[N] {
	return &ErrUnexpectedNode[N]{
		Expecteds: expecteds,
		Prev:      prev,
		Got:       got,
	}
}
*/

type ErrParsing struct {
	Err           error
	PossibleCause error
}

func (e ErrParsing) Error() string {
	var builder strings.Builder

	if e.Err == nil {
		builder.WriteString("something went wrong")
	} else {
		builder.WriteString(e.Err.Error())
	}

	if e.PossibleCause == nil {
		return builder.String()
	}

	builder.WriteString(", possible cause: ")
	builder.WriteString(e.PossibleCause.Error())

	return builder.String()
}

func (e ErrParsing) Unwrap() error {
	return e.Err
}

func NewErrParsing(err error, possibleCause error) *ErrParsing {
	return &ErrParsing{
		Err:           err,
		PossibleCause: possibleCause,
	}
}
