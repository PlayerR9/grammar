package parser

import (
	"strconv"
	"strings"

	gcers "github.com/PlayerR9/go-commons/errors"
	gcstr "github.com/PlayerR9/go-commons/strings"
	"github.com/PlayerR9/grammar/PREV/internal"
)

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

// ErrParsing is the error for parsing errors.
type ErrParsing struct {
	// Err is the error.
	Err error

	// PossibleCause is the possible cause of the error.
	PossibleCause error
}

// Error implements the error interface.
//
// Message: "<err>, possible cause: <possible cause>".
func (e ErrParsing) Error() string {
	var builder strings.Builder

	builder.WriteString(gcers.Error(e.Err))

	if e.PossibleCause == nil {
		return builder.String()
	}

	builder.WriteString(", possible cause: ")
	builder.WriteString(e.PossibleCause.Error())

	return builder.String()
}

// Unwrap returns the underlying error.
//
// Returns:
//   - error: The underlying error.
func (e ErrParsing) Unwrap() error {
	return e.Err
}

// NewErrParsing creates a new ErrParsing.
//
// Parameters:
//   - err: The error.
//   - possibleCause: The possible cause of the error.
//
// Returns:
//   - *ErrParsing: A pointer to the new ErrParsing. Never returns nil.
func NewErrParsing(err error, possible_cause error) *ErrParsing {
	return &ErrParsing{
		Err:           err,
		PossibleCause: possible_cause,
	}
}
