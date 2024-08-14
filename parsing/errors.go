package parsing

import (
	"strconv"
	"strings"

	gcstr "github.com/PlayerR9/go-commons/strings"
	gr "github.com/PlayerR9/grammar/grammar"
)

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
