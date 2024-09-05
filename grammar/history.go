package grammar

import (
	"errors"

	gcers "github.com/PlayerR9/go-commons/errors"
	internal "github.com/PlayerR9/grammar/grammar/internal"
)

type History[T internal.TokenTyper] struct {
	timeline []*Item[T]
	arrow    int
}

func NewHistory[T internal.TokenTyper](timeline []*Item[T]) *History[T] {
	return &History[T]{
		timeline: timeline,
		arrow:    0,
	}
}

func (h History[T]) CanWalk() bool {
	return h.arrow < len(h.timeline)
}

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

func (h *History[T]) AddEvent(item *Item[T]) {
	if item == nil {
		return
	}

	h.timeline = append(h.timeline, item)
}

func (h History[T]) Copy() *History[T] {
	h_copy := make([]*Item[T], len(h.timeline))
	copy(h_copy, h.timeline)

	return &History[T]{
		timeline: h_copy,
		arrow:    h.arrow,
	}
}
