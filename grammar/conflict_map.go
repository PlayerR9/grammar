package grammar

import (
	"iter"

	internal "github.com/PlayerR9/grammar/grammar/internal"

	"github.com/PlayerR9/go-commons/set"
)

type ConflictMap[T internal.TokenTyper] struct {
	table map[T]*set.Set[*Item[T]]
}

func NewConflictMap[T internal.TokenTyper]() *ConflictMap[T] {
	return &ConflictMap[T]{
		table: make(map[T]*set.Set[*Item[T]]),
	}
}

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

func (cm ConflictMap[T]) Len() int {
	return len(cm.table)
}
