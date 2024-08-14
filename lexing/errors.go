package lexing

import (
	"fmt"
	"strings"
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
