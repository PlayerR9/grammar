package parsing

import (
	"fmt"
	"strconv"
	"strings"

	gcstr "github.com/PlayerR9/go-commons/strings"
	gr "github.com/PlayerR9/grammar/grammar"
)

// ErrParsing is an error that occurs while lexing.
type ErrParsing struct {
	// StartPos is the start position of the error.
	StartPos int

	// Delta is the delta of the error.
	Delta int

	// Reason is the reason of the error.
	Reason error

	// Suggestion is the suggestion for solving the error.
	Suggestion string
}

// Error implements the error interface.
//
// Format:
//
//	"an error occurred while parsing"
func (e *ErrParsing) Error() string {
	if e.Reason == nil {
		return "an error occurred while parsing"
	}

	return fmt.Sprintf("error while parsing: %s", e.Reason.Error())
}

// NewErrParsing creates a new error.
//
// Parameters:
//   - startPos: The start position of the error.
//   - delta: The delta of the error.
//   - reason: The reason of the error.
//
// Returns:
//   - *ErrParsing: The new error. Never returns nil.
func NewErrParsing(startPos int, delta int, reason error) *ErrParsing {
	return &ErrParsing{
		StartPos: startPos,
		Delta:    delta,
		Reason:   reason,
	}
}

// SetSuggestion sets the suggestion for solving the error.
//
// Parameters:
//   - suggestions: The suggestions for solving the error.
func (e *ErrParsing) SetSuggestion(suggestions ...string) {
	e.Suggestion = strings.Join(suggestions, " ")
}

// Unwrap returns the reason of the error.
//
// Returns:
//   - error: The reason of the error.
func (e *ErrParsing) Unwrap() error {
	return e.Reason
}

// ChangeReason changes the reason of the error.
//
// Parameters:
//   - reason: The new reason of the error.
func (e *ErrParsing) ChangeReason(reason error) {
	e.Reason = reason
}

// ErrUnexpectedToken is an error that occurs when an unexpected token is
// encountered.
type ErrUnexpectedToken[T gr.TokenTyper] struct {
	// Expecteds is the list of expected tokens.
	Expecteds []T

	// Got is the token that was encountered.
	Got *T

	// After is the token that was encountered after the expected token.
	After *T
}

// Error implements the error interface.
//
// Format:
//
//	"expected either <value 0>, <value 1>, <value 2>, ..., or <value n> after <after> instead, got <actual> instead"
func (e *ErrUnexpectedToken[T]) Error() string {
	var builder strings.Builder

	builder.WriteString("expected ")

	if len(e.Expecteds) == 0 {
		builder.WriteString("nothing")
	} else {
		elems := gcstr.SliceOfStringer(e.Expecteds)
		gcstr.QuoteStrings(elems)

		builder.WriteString(gcstr.EitherOrString(elems))
	}

	if e.After != nil {
		builder.WriteString(" after ")
		builder.WriteString(strconv.Quote((*e.After).String()))
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

// NewErrUnexpectedToken creates a new unexpected token error.
//
// Parameters:
//   - after: The token that was encountered after the expected token.
//   - got: The token that was encountered.
//   - expecteds: The expected tokens.
//
// Returns:
//   - *ErrUnexpectedToken[T]: The new error. Never returns nil.
func NewErrUnexpectedToken[T gr.TokenTyper](after *T, got *T, expecteds ...T) *ErrUnexpectedToken[T] {
	return &ErrUnexpectedToken[T]{
		Expecteds: expecteds,
		Got:       got,
		After:     after,
	}
}
