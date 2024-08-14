package displayer

import (
	"fmt"
	"strings"
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
