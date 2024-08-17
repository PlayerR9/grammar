package matcher

// ErrNoMatch is the error that occurs when the matcher does not match any rule.
type ErrNoMatch struct {
	// Reason is the reason of the error.
	Reason error
}

// Error implements the error interface.
//
// Message: "no match: <reason>"
func (e ErrNoMatch) Error() string {
	if e.Reason == nil {
		return "no match"
	}

	return "no match: " + e.Reason.Error()
}

// Unwrap returns the reason of the error.
//
// Returns:
//   - error: The reason of the error.
func (e ErrNoMatch) Unwrap() error {
	return e.Reason
}

// NewErrNoMatch creates a new error that occurs when the matcher does not match any rule.
//
// Parameters:
//   - reason: The reason of the error.
//
// Returns:
//   - *ErrNoMatch: The new error. Never returns nil.
func NewErrNoMatch(reason error) *ErrNoMatch {
	return &ErrNoMatch{
		Reason: reason,
	}
}
