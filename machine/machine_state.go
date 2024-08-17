package machine

import "fmt"

// MachineStater is an interface that defines the behavior of a machine state (such as a FSM).
type MachineStater interface {
	// New initializes the state of the machine; making it ready to execute.
	New()

	// Init initializes/resets the state of the machine. This is used to initialize arguments
	// and/or information that allow the machine to be reused.
	//
	// Parameters:
	//   - args: The arguments to pass to the state machine.
	//
	// Returns:
	//   - Stater: The next state of the machine.
	//
	// If a nil state is returned, it is treated as an end state with nil reason.
	// (return nil equals NewEndState(nil))
	Init(args ...any) Stater

	// Cleanup cleans up the state of the machine. This is used to apply the necessary cleanup
	// operations that allows to clear memory and other resources.
	Cleanup()

	// Execute executes the current state of the machine.
	//
	// Parameters:
	//   - state: The current state of the machine.
	//
	// Returns:
	//   - S: The transitioned state of the machine.
	Execute(state Stater) Stater

	// Error returns the error of the machine. Only called when a critical error occurs.
	//
	// Returns:
	//   - error: The error of the machine.
	//
	// For non-critical errors, each implementation should offer methods to handle them
	// and never override this method.
	Error() error
}

// Run executes the state machine. Does nothing if ms is nil.
//
// Parameters:
//   - ms: The state machine to execute.
//   - args: The arguments to pass to the state machine.
func Run(ms MachineStater, args ...any) error {
	if ms == nil {
		return nil
	}

	ms.New()

	var state Stater = NewInitState()

	for {
		var next_state Stater

		switch state := state.(type) {
		case *SystemState:
			switch state.State {
			case InitSS:
				next_state = ms.Init(args)
				if next_state == nil {
					next_state = NewEndState(nil)
				}
			case EndSS:
				return nil
			case CriticalSS:
				err := ms.Error()

				ms.Cleanup()

				return err
			default:
				return fmt.Errorf("invalid state: %d", state.State)
			}
		default:
			new_state := ms.Execute(state)
			if new_state == nil {
				new_state = NewEndState(nil)
			}
		}

		state = next_state
	}

	return nil
}
