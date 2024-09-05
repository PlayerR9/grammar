package grammar

import (
	"iter"
	"slices"

	dbg "github.com/PlayerR9/go-debug/assert"
	internal "github.com/PlayerR9/grammar/internal"
)

// State is the state of the goto graph.
type State[T internal.TokenTyper] struct {
	// items is the list of all items in the state. The first item is the seed item.
	items []*Item[T]

	// nexts is the list of next states.
	nexts []*State[T]
}

// NewState creates a new state.
//
// Parameters:
//   - seed: The seed item.
//   - closure: The closure of the state.
//
// Returns:
//   - *State[T]: The created state. Never returns nil.
func NewState[T internal.TokenTyper](seed *Item[T], closure []*Item[T]) *State[T] {
	idx := slices.Index(closure, seed)
	if idx != -1 {
		closure = slices.Delete(closure, idx, idx+1)
	}

	closure = append([]*Item[T]{seed}, closure...)

	return &State[T]{
		items: closure,
	}
}

// AddNext adds the next state to the state. Nil or already added states are ignored.
//
// Parameters:
//   - next: The next state.
func (s *State[T]) AddNext(next *State[T]) {
	if next != nil && !slices.Contains(s.nexts, next) {
		s.nexts = append(s.nexts, next)
	}
}

// Rule returns an iterator over the rules in the state. This excludes the seed item.
//
// Returns:
//   - iter.Seq2[int, *Item[T]]: The iterator. Never returns nil.
func (s *State[T]) Rule() iter.Seq2[int, *Item[T]] {
	return func(yield func(int, *Item[T]) bool) {
		for i := 0; i < len(s.items); i++ {
			if !yield(i, s.items[i]) {
				break
			}
		}
	}
}

// IsOfSeed checks if the state is of the seed item.
//
// Parameters:
//   - item: The item to check.
//
// Returns:
//   - bool: True if the state is of the seed item, false otherwise.
func (s State[T]) IsOfSeed(item *Item[T]) bool {
	if item == nil {
		return false
	}

	dbg.AssertThat("s.items", dbg.NewOrderedAssert(len(s.items)).GreaterThan(0)).Panic()

	return s.items[0].Equals(item)
}

// NextState returns an iterator over the next states.
//
// Returns:
//   - iter.Seq[*State[T]]: The iterator. Never returns nil.
func (s State[T]) NextState() iter.Seq[*State[T]] {
	return func(yield func(*State[T]) bool) {
		for _, next := range s.nexts {
			if !yield(next) {
				break
			}
		}
	}
}

// Seed returns the seed item.
//
// Returns:
//   - *Item[T]: The seed item. Never returns nil.
func (s State[T]) Seed() *Item[T] {
	return s.items[0]
}
