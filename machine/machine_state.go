package machine

import (
	"fmt"
	"io"
)

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

// MachineState is a state machine.
type MachineState[T any] struct {
	// table is the state table of the state machine.
	table map[SystemState]StepFunc[T]

	// cleanup is the cleanup function of the state machine.
	//
	// This is used to apply the necessary cleanup operations that allows to clear
	// memory and other resources.
	//
	// Parameters:
	//   - info: The state machine information.
	cleanup func(info T)
}

// NewMachineState creates a new state machine. If the init function is nil,
// it will be set to func(_ T, _ *rune) SystemState { return EndSS }.
//
// Parameters:
//   - init: The initialization function of the state machine.
//
// Returns:
//   - *MachineState: The new state machine. Never return nil.
func NewMachineState[T any](init func(info T, _ *rune) (SystemState, error)) *MachineState[T] {
	if init == nil {
		init = func(_ T, _ *rune) (SystemState, error) { return EndSS, nil }
	}

	table := make(map[SystemState]StepFunc[T])
	table[InitSS] = init

	return &MachineState[T]{
		table: table,
	}
}

// AddState adds a new state to the state machine. Does nothing if f is nil.
//
// Parameters:
//   - state: The state to add.
//   - f: The function to execute when the state is reached.
//
// If the state already exists, it is overwritten.
// InitSS, CriticalSS and EndSS are never added.
func (ms *MachineState[T]) AddState(state SystemState, f StepFunc[T]) {
	if f == nil {
		return
	} else if state == InitSS || state == EndSS {
		return
	}

	ms.table[state] = f
}

// WithCleanup sets the cleanup function of the state machine.
//
// Parameters:
//   - cleanup: The cleanup function of the state machine.
func (ms *MachineState[T]) WithCleanup(cleanup CleanupFunc[T]) {
	ms.cleanup = cleanup
}

// RunFunc is a function that executes a state machine and returns it.
//
// Parameters:
//   - scanner: The scanner to process.
//
// Returns:
//   - T: The state machine information.
//   - error: The error of the machine.
type RunFunc[T any] func(io.RuneScanner) (T, error)

// Make makes a new state machine.
//
// Parameters:
//   - info: The state machine information.
//
// Returns:
//   - RunFunc: The run function of the state machine.
//   - func(): The cleanup function of the state machine.
func (ms MachineState[T]) Make(info T) (RunFunc[T], func()) {
	fn := func(scanner io.RuneScanner) (T, error) {
		f, ok := ms.table[InitSS]
		if !ok {
			return info, fmt.Errorf("invalid state: %d", InitSS)
		}

		state, err := f(info, nil)
		if err != nil {
			return info, err
		}

		for state != EndSS {
			f, ok := ms.table[state]
			if !ok {
				return info, fmt.Errorf("invalid state: %d", state)
			}

			var char *rune

			c, _, err := scanner.ReadRune()
			if err == nil {
				char = &c
			} else if err != io.EOF {
				return info, err
			}

			state, err = f(info, char)
			if err != nil {
				return info, err
			}
		}

		return info, nil
	}

	var cleanup func()

	if ms.cleanup == nil {
		cleanup = func() {}
	} else {
		cleanup = func() {
			ms.cleanup(info)
		}
	}

	return fn, cleanup
}
