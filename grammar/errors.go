package grammar

import (
	"strconv"
	"strings"

	gcstr "github.com/PlayerR9/go-commons/strings"
	internal "github.com/PlayerR9/grammar/internal"
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
