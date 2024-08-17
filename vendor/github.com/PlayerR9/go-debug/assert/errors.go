package assert

// AssertionFailed is the message that is shown when an assertion fails.
const AssertionFailed string = "assertion failed: "

// ErrAssertionFailed is the error returned when an assertion fails.
type ErrAssertionFailed struct {
	// Msg is the message that is shown when the assertion fails.
	Msg string
}

// Error implements the error interface.
//
// The message is prefixed with the AssertionFailed constant.
func (e *ErrAssertionFailed) Error() string {
	return AssertionFailed + e.Msg
}

// NewErrAssertionFailed is a constructor for ErrAssertionFailed.
//
// Parameters:
//   - msg: the message that is shown when the assertion fails.
//
// Returns:
//   - *ErrAssertionFailed: the error. Never returns nil.
func NewErrAssertionFailed(msg string) *ErrAssertionFailed {
	return &ErrAssertionFailed{Msg: msg}
}
