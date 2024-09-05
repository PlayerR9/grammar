package internal

//go:generate stringer -type=ActionType -linecomment

// ActionType is the action type.
type ActionType int8

const (
	// ActErrorType is the error action type.
	ActErrorType ActionType = iota // ERROR

	// ActShiftType is the shift action type.
	ActShiftType // SHIFT

	// ActReduceType is the reduce action type.
	ActReduceType // REDUCE

	// ActAcceptType is the accept action type.
	ActAcceptType // ACCEPT
)
