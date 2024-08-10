package errors

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
