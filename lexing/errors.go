package lexing

import (
	"fmt"
	"strconv"
	"strings"

	gcstr "github.com/PlayerR9/go-commons/strings"
)

type ErrLexing struct {
	StartPos   int
	Delta      int
	Reason     error
	Suggestion string
}

func (e *ErrLexing) Error() string {
	if e.Reason == nil {
		return "an error occurred while lexing"
	}

	return fmt.Sprintf("error while lexing: %s", e.Reason.Error())
}

func NewErrLexing(startPos int, delta int, reason error) *ErrLexing {
	return &ErrLexing{
		StartPos: startPos,
		Delta:    delta,
		Reason:   reason,
	}
}

func (e *ErrLexing) SetSuggestion(suggestions ...string) {
	e.Suggestion = strings.Join(suggestions, " ")
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
