package grammar

import (
	"iter"

	"github.com/PlayerR9/go-commons/set"
	"github.com/PlayerR9/grammar/internal"
)

// ConflictMap is the conflict map.
type ConflictMap[T internal.TokenTyper] struct {
	// table is the conflict map.
	table map[T]*set.Set[*Item[T]]
}

// NewConflictMap creates a new conflict map.
//
// Returns:
//   - *ConflictMap[T]: The new conflict map. Never returns nil.
func NewConflictMap[T internal.TokenTyper]() *ConflictMap[T] {
	return &ConflictMap[T]{
		table: make(map[T]*set.Set[*Item[T]]),
	}
}

// Cleanup cleans up the conflict map for the garbage collector.
func (cm *ConflictMap[T]) Cleanup() {
	if len(cm.table) == 0 {
		return
	}

	for k, v := range cm.table {
		v.Clear()
		cm.table[k] = nil
		delete(cm.table, k)
	}

	cm.table = nil
}

// Reset resets the conflict map.
func (cm *ConflictMap[T]) Reset() {
	if cm.table == nil {
		cm.table = make(map[T]*set.Set[*Item[T]])
	}

	for k, v := range cm.table {
		v.Clear()
		cm.table[k] = nil
		delete(cm.table, k)
	}
}

// conflicting_items is a helper function that returns the conflicting items.
//
// Parameters:
//   - items: The items.
//
// Returns:
//   - *set.Set[*Item[T]]: The conflicting items.
func conflicting_items[T internal.TokenTyper](items []*Item[T]) *set.Set[*Item[T]] {
	if len(items) < 2 {
		return nil
	}

	conflict_map := make(map[*Item[T]][]*Item[T])

	for i := 0; i < len(items)-1; i++ {
		var conflicts []*Item[T]

		for j := i + 1; j < len(items); j++ {
			if !items[i].IsInConflictWith(items[j]) {
				continue
			}

			conflicts = append(conflicts, items[j])
		}

		if len(conflicts) > 0 {
			conflict_map[items[i]] = conflicts
		}
	}

	if len(conflict_map) == 0 {
		return nil
	}

	conflicts := set.NewSet[*Item[T]]()

	for k, vals := range conflict_map {
		_ = conflicts.Add(k)

		for _, item := range vals {
			_ = conflicts.Add(item)
		}
	}

	return conflicts
}

// Init initializes the conflict map.
//
// Parameters:
//   - items: The items.
func (cm *ConflictMap[T]) Init(items map[T][]*Item[T]) {
	if len(items) == 0 {
		return
	}

	cm.Reset()

	for symbol, item_list := range items {
		items := conflicting_items(item_list)
		if items == nil || items.Size() == 0 {
			continue
		}

		cm.table[symbol] = items
	}
}

// Entry returns an iterator over the conflict map where the key is the token symbol
// and the value is the conflicting items.
//
// Returns:
//   - iter.Seq2[T, *Item[T]]: The iterator. Never returns nil.
func (cm ConflictMap[T]) Entry() iter.Seq2[T, *Item[T]] {
	fn := func(yield func(T, *Item[T]) bool) {
		for s, items := range cm.table {
			for item := range items.All() {
				if !yield(s, item) {
					return
				}
			}
		}
	}

	return fn
}

// All returns an iterator over the conflict map where the key is the token symbol
// and the value is an iterator over the conflicting items.
//
// Returns:
//   - iter.Seq2[T, iter.Seq[*Item[T]]]: The iterator. Never returns nil.
func (cm ConflictMap[T]) All() iter.Seq2[T, iter.Seq[*Item[T]]] {
	fn := func(yield func(T, iter.Seq[*Item[T]]) bool) {
		for s, items := range cm.table {
			sub_fn := func(yield func(*Item[T]) bool) {

				for item := range items.All() {
					if !yield(item) {
						return
					}
				}
			}

			if !yield(s, sub_fn) {
				return
			}
		}
	}

	return fn
}

// Len returns the number of items in the conflict map.
//
// Returns:
//   - int: The number of items in the conflict map.
func (cm ConflictMap[T]) Len() int {
	return len(cm.table)
}
