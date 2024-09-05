package machine

import (
	"fmt"
	"io"
)

// MachineState is a state machine.
type MachineState[T any] struct {
	// table is the state table of the state machine.
	table map[SystemState]StepFunc[T]

	// transition is the transition table of the state machine.
	transition map[SystemState]SystemState

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
func NewMachineState[T any](to SystemState, init func(info T, _ *rune) (SystemState, error)) *MachineState[T] {
	if init == nil {
		init = func(_ T, _ *rune) (SystemState, error) { return EndSS, nil }
	}

	table := make(map[SystemState]StepFunc[T])
	table[InitSS] = init

	return &MachineState[T]{
		table:      table,
		transition: make(map[SystemState]SystemState),
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
func (ms *MachineState[T]) AddState(from, to SystemState, f StepFunc[T]) {
	if f == nil || from == InitSS || from == EndSS {
		return
	}

	ms.table[from] = f
	ms.transition[from] = to
}

// WithCleanup sets the cleanup function of the state machine.
//
// Parameters:
//   - cleanup: The cleanup function of the state machine.
func (ms *MachineState[T]) WithCleanup(cleanup CleanupFunc[T]) {
	ms.cleanup = cleanup
}

// Make makes a new state machine.
//
// Parameters:
//   - info: The state machine information.
//
// Returns:
//   - RunFunc: The run function of the state machine.
//   - func(): The cleanup function of the state machine.
//
// Be sure to check the cleanup function if it is not nil. If it is only when
// WithCleanup was called with a nil value or never called.
func (ms MachineState[T]) Make(info T) (RunFunc[T], func()) {
	msr := MachineStepRunner[T]{
		table: ms.table,
	}

	f, ok := ms.table[InitSS]
	if !ok {
		run_fn := func(scanner io.RuneScanner) (T, error) {
			return info, fmt.Errorf("invalid state: %d", InitSS)
		}

		return run_fn, nil
	}

	fn := func(scanner io.RuneScanner) (T, error) {
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
			if err != nil {
				if err != io.EOF {
					return info, err
				}
			} else {
				char = &c
			}

			state, err = f(info, char)
			if err != nil {
				return info, err
			}
		}

		return info, nil
	}

	if ms.cleanup == nil {
		return fn, nil
	}

	cleanup_fn := func() {
		ms.cleanup(info)
	}

	return fn, cleanup_fn
}

/*
// Make makes a new state machine.
//
// Parameters:
//   - info: The state machine information.
//
// Returns:
//   - RunFunc: The run function of the state machine.
//   - func(): The cleanup function of the state machine.
//
// Be sure to check the cleanup function if it is not nil. If it is only when
// WithCleanup was called with a nil value or never called.
func (ms MachineState[T]) MakeSteps(info T) (RunFunc[T], func()) {
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

	if ms.cleanup == nil {
		return fn, nil
	}

	cleanup_fn := func() {
		ms.cleanup(info)
	}

	return fn, cleanup_fn
}
*/
