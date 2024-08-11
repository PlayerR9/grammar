package lexing

import (
	"fmt"
	"strconv"
	"strings"

	gcstr "github.com/PlayerR9/go-commons/strings"
)

// ErrLexing is an error that occurs while lexing.
type ErrLexing struct {
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
//	"an error occurred while lexing"
func (e *ErrLexing) Error() string {
	if e.Reason == nil {
		return "an error occurred while lexing"
	}

	return fmt.Sprintf("error while lexing: %s", e.Reason.Error())
}

// NewErrLexing creates a new error.
//
// Parameters:
//   - startPos: The start position of the error.
//   - delta: The delta of the error.
//   - reason: The reason of the error.
//
// Returns:
//   - *ErrLexing: The new error. Never returns nil.
func NewErrLexing(startPos int, delta int, reason error) *ErrLexing {
	return &ErrLexing{
		StartPos: startPos,
		Delta:    delta,
		Reason:   reason,
	}
}

// SetSuggestion sets the suggestion for solving the error.
//
// Parameters:
//   - suggestions: The suggestions for solving the error.
func (e *ErrLexing) SetSuggestion(suggestions ...string) {
	e.Suggestion = strings.Join(suggestions, " ")
}

// Unwrap returns the reason of the error.
//
// Returns:
//   - error: The reason of the error.
func (e *ErrLexing) Unwrap() error {
	return e.Reason
}

// ChangeReason changes the reason of the error.
//
// Parameters:
//   - reason: The new reason of the error.
func (e *ErrLexing) ChangeReason(reason error) {
	e.Reason = reason
}

// ErrUnexpectedRune is an error that occurs when an unexpected rune is
// encountered.
type ErrUnexpectedRune struct {
	// Expecteds is a list of expected runes.
	Expecteds []rune

	// Prev is the rune that was encountered before the expected rune.
	Prev *rune

	// Got is the rune that was encountered.
	Got *rune
}

// Error implements the error interface.
//
// Format:
//
//	"expected <expected> <prev>, got <actual> instead"
func (e *ErrUnexpectedRune) Error() string {
	var builder strings.Builder

	builder.WriteString("expected ")

	if len(e.Expecteds) == 0 {
		builder.WriteString("nothing")
	} else {
		elems := gcstr.SliceOfRunes(e.Expecteds)
		gcstr.QuoteStrings(elems)

		builder.WriteString(gcstr.EitherOrString(elems))
	}

	builder.WriteString(" ")

	if e.Prev == nil {
		builder.WriteString("at the beginning")
	} else {
		builder.WriteString("after ")
		builder.WriteString(strconv.QuoteRune(*e.Prev))
	}

	builder.WriteString(", got ")

	if e.Got == nil {
		builder.WriteString("nothing")
	} else {
		builder.WriteString(strconv.QuoteRune(*e.Got))
	}

	builder.WriteString(" instead")

	return builder.String()
}

// NewErrUnexpectedRune creates a new error.
//
// Parameters:
//   - expecteds: The expected runes.
//   - prev: The rune that was encountered before the expected rune.
//   - got: The rune that was encountered.
//
// Returns:
//   - *ErrUnexpectedRune: The new error. Never returns nil.
func NewErrUnexpectedRune(prev *rune, got *rune, expecteds ...rune) *ErrUnexpectedRune {
	return &ErrUnexpectedRune{
		Expecteds: expecteds,
		Prev:      prev,
		Got:       got,
	}
}

// ErrNoClosestWordFound is an error when no closest word is found.
type ErrNoClosestWordFound struct{}

// Error implements the error interface.
//
// Message: "no closest word was found"
func (e *ErrNoClosestWordFound) Error() string {
	return "no closest word was found"
}

// NewErrNoClosestWordFound creates a new ErrNoClosestWordFound.
//
// Returns:
//   - *ErrNoClosestWordFound: The new ErrNoClosestWordFound.
func NewErrNoClosestWordFound() *ErrNoClosestWordFound {
	return &ErrNoClosestWordFound{}
}
