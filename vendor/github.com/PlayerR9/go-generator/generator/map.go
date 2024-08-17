package generator

import (
	"cmp"
	"iter"
	"slices"
)

// ordered_map is a map that is ordered by the keys.
type ordered_map[K cmp.Ordered, V any] struct {
	// values is a map of the values in the map.
	values map[K]V

	// keys is a slice of the keys in the map.
	keys []K
}

// new_ordered_map creates a new OrderedMap.
//
// Returns:
//   - *OrderedMap: A pointer to the newly created OrderedMap.
//     Never returns nil.
func new_ordered_map[K cmp.Ordered, V any]() *ordered_map[K, V] {
	return &ordered_map[K, V]{
		values: make(map[K]V),
		keys:   make([]K, 0),
	}
}

// add adds a key-value pair to the map.
//
// Parameters:
//   - key: The key to add.
//   - value: The value to add.
//   - force: If true, the value will be added even if the key already exists. If
//     false, the value will not be added if the key already exists.
//
// Returns:
//   - bool: True if the value was added to the map, false otherwise.
func (m *ordered_map[K, V]) add(key K, value V, force bool) bool {
	pos, ok := slices.BinarySearch(m.keys, key)

	if !ok {
		m.keys = slices.Insert(m.keys, pos, key)
	}

	if ok && !force {
		return false
	}

	m.values[key] = value

	return true
}

// size is a method that returns the number of keys in the map.
//
// Returns:
//   - int: The number of keys in the map.
func (m ordered_map[K, V]) size() int {
	return len(m.keys)
}

// Map is a method that returns the map of the values in the map.
//
// Returns:
//   - map[K]V: The map of the values in the map. Never returns nil.
func (m ordered_map[K, V]) Map() map[K]V {
	return m.values
}

// Keys is a method that returns the keys in the map.
//
// Returns:
//   - []K: The keys in the map.
func (m ordered_map[K, V]) Keys() []K {
	return m.keys
}

// Entry returns an iterator that iterates over the entries in the map according
// to the order of the keys.
//
// Returns:
//   - iter.Seq2[K, V]: The iterator. Never returns nil.
func (m ordered_map[K, V]) Entry() iter.Seq2[K, V] {
	return func(yield func(key K, value V) bool) {
		for _, key := range m.keys {
			if !yield(key, m.values[key]) {
				break
			}
		}
	}
}
