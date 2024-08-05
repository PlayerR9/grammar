package parser

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
	var expected string

	if len(e.Expecteds) == 0 {
		expected = "nothing"
	} else {
		values := make([]string, 0, len(e.Expecteds))

		for _, expected := range e.Expecteds {
			values = append(values, expected.String())
		}

		expected = gcstr.EitherOrString(values, true)
	}

	var got string

	if e.Got == nil {
		got = "nothing"
	} else {
		got = strconv.Quote((*e.Got).String())
	}

	var builder strings.Builder

	builder.WriteString("expected ")
	builder.WriteString(expected)

	if e.After != nil {
		builder.WriteString(" after ")
		builder.WriteString((*e.After).String())
	}

	builder.WriteString(", got ")
	builder.WriteString(got)
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
