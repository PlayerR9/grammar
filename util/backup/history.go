package backup

import (
	"iter"

	gcstck "github.com/PlayerR9/go-commons/stack"
)

// Execute executes a walker.
//
// The walker is expected to have a state machine behavior, where it can:
//   - Walk one event at a time using WalkOne.
//   - Return the next events using NextEvents.
//   - Return whether it has an error using HasError.
//
// This function will return an iterator of all possible solutions (i.e., all
// possible combinations of events that result in a solution).
//
// If the walker has an error, it is not considered a solution.
//
// The iterator will yield each solution in the order of the events (i.e., the
// first event is the one that is the first in the timeline).
//
// Parameters:
//   - init_fn: The function that initializes the walker.
//
// Returns:
//   - iter.Seq[W]: The iterator. Never returns nil.
func Execute[W interface {
	HasError() bool
	WalkOne(event T) bool
	NextEvents() []T
}, T any](init_fn func() W) iter.Seq[W] {
	if init_fn == nil {
		return func(yield func(W) bool) {}
	}

	fn := func(yield func(w W) bool) {
		first_walker := init_fn()

		if first_walker.HasError() {
			_ = yield(first_walker)

			return
		}

		first_events := first_walker.NextEvents()
		if len(first_events) == 0 {
			_ = yield(first_walker)

			return
		}

		is_solution := first_walker.WalkOne(first_events[0])

		if is_solution && !yield(first_walker) {
			return
		}

		possible_paths := gcstck.NewStack[[]T]()

		for i := len(first_events) - 1; i > 0; i-- {
			possible_paths.Push([]T{first_events[i]})
		}

		var invalid_walkers []W

		for {
			path, ok := possible_paths.Pop()
			if !ok {
				break
			}

			active_walker := init_fn()

			for _, event := range path {
				is_solution := active_walker.WalkOne(event)
				if is_solution {
					panic("should not be solution")
				} else if active_walker.HasError() {
					panic("should not have error")
				}
			}

			for {
				nexts := active_walker.NextEvents()
				if len(nexts) == 0 {
					if active_walker.HasError() {
						invalid_walkers = append(invalid_walkers, active_walker)
					} else if !yield(active_walker) {
						return
					} else {
						break
					}
				}

				is_solution := active_walker.WalkOne(nexts[0])

				if is_solution {
					if !yield(active_walker) {
						return
					}

					break
				}

				for i := len(nexts) - 1; i > 0; i-- {
					possible_paths.Push(append(path, nexts[i]))
				}

				if active_walker.HasError() {
					invalid_walkers = append(invalid_walkers, active_walker)

					break
				}
			}
		}

		for i := len(invalid_walkers) - 1; i >= 0; i-- {
			if !yield(invalid_walkers[i]) {
				return
			}
		}
	}

	return fn
}
