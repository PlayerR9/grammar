package machine

import "io"

// SystemState is a system state.
type SystemState int

const (
	// EndSS is the end state of the system.
	EndSS SystemState = iota - 1

	// InitSS is the initial state of the system.
	InitSS
)

// CleanupFunc is a function that cleans up the state of the machine.
//
// Parameters:
//   - info: The state machine information.
//
// This is used to apply the necessary cleanup operations that allows to clear
// memory and other resources.
type CleanupFunc[T any] func(info T)

// StepFunc is a function that executes a state machine step.
//
// Parameters:
//   - info: The state machine information.
//   - char: The character to process. Nil if the stream is ended.
//
// Returns:
//   - SystemState: The next state of the machine.
//   - error: The error of the machine. Only used for when a panic-level of error
//     occurs.
type StepFunc[T any] func(info T, char *rune) (SystemState, error)

// RunFunc is a function that executes a state machine and returns it.
//
// Parameters:
//   - scanner: The scanner to process.
//
// Returns:
//   - T: The state machine information.
//   - error: The error of the machine.
type RunFunc[T any] func(scanner io.RuneScanner) (T, error)
