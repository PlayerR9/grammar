package grammar

import (
	"errors"

	gcers "github.com/PlayerR9/go-commons/errors"
	internal "github.com/PlayerR9/grammar/internal"
)

// History is the history of the parser.
type History[T internal.TokenTyper] struct {
	// timeline is the timeline of the history.
	timeline []*Item[T]

	// arrow is the current position in the timeline.
	arrow int
}

// NewHistory creates a new history with the given timeline.
//
// Parameters:
//   - timeline: The timeline of the history.
//
// Returns:
//   - *History: The new history. Never returns nil.
func NewHistory[T internal.TokenTyper](timeline []*Item[T]) *History[T] {
	return &History[T]{
		timeline: timeline,
		arrow:    0,
	}
}

// CanWalk checks if the history can walk.
//
// Returns:
//   - bool: True if the history can walk, false otherwise.
func (h History[T]) CanWalk() bool {
	return h.arrow < len(h.timeline)
}

// Walk walks the history.
//
// Parameters:
//   - do: The function to be called for each item in the history.
//
// Returns:
//   - error: An error if walk errors.
func (h *History[T]) Walk(do func(item *Item[T]) error) error {
	if h.arrow >= len(h.timeline) {
		return errors.New("already at the end of the history")
	} else if do == nil {
		return gcers.NewErrNilParameter("do")
	}

	item := h.timeline[h.arrow]
	h.arrow++

	err := do(item)
	if err != nil {
		return err
	}

	return nil
}

// AddEvent adds an event to the history. It ignores nil events.
//
// Parameters:
//   - item: The item to be added to the history.
func (h *History[T]) AddEvent(item *Item[T]) {
	if item == nil {
		return
	}

	h.timeline = append(h.timeline, item)
}

// Copy creates a copy of the history.
//
// Returns:
//   - *History: The copy. Never returns nil.
func (h History[T]) Copy() *History[T] {
	h_copy := make([]*Item[T], len(h.timeline))
	copy(h_copy, h.timeline)

	return &History[T]{
		timeline: h_copy,
		arrow:    h.arrow,
	}
}
