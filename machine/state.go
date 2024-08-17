package machine

// Stater is an interface that defines the behavior of a state. 0 is reserved for the initial state
// while -1 is reserved for the terminal state. -2 is reserved for critical errors.
type Stater interface {
}

// SystemStater is the state of the system.
type SystemStater int

const (
	// InitSS is the initial state of the system.
	InitSS SystemStater = iota

	// EndSS is the end state of the system.
	EndSS

	// CriticalSS is the critical state of the system.
	CriticalSS
)

// SystemState is the state of the system.
type SystemState struct {
	// State is the state of the system.
	State SystemStater

	// Reason is the error of the system. (if any)
	Reason error
}

// NewInitState creates a new initial state.
//
// Returns:
//   - *SystemState: The new initial state. Never returns nil.
func NewInitState() *SystemState {
	return &SystemState{
		State: InitSS,
	}
}

// NewEndState creates a new end state.
//
// Parameters:
//   - err: The error of the end state. (If any)
//
// Returns:
//   - *SystemState: The new end state. Never returns nil.
func NewEndState(err error) *SystemState {
	return &SystemState{
		State:  EndSS,
		Reason: err,
	}
}

// NewCriticalState creates a new critical state.
//
// Parameters:
//   - err: The error of the critical state.
//
// Returns:
//   - *SystemState: The new critical state. Never returns nil.
func NewCriticalState(err error) *SystemState {
	return &SystemState{
		State:  CriticalSS,
		Reason: err,
	}
}
