package slices

// PredicateFilter is a type that defines a slice filter function.
//
// Parameters:
//   - T: The type of the elements in the slice.
//
// Returns:
//   - bool: True if the element satisfies the filter function, otherwise false.
type PredicateFilter[T any] func(T) bool

// Intersect returns a PredicateFilter function that checks if an element
// satisfies all the PredicateFilter functions in funcs.
//
// Parameters:
//   - funcs: A slice of PredicateFilter functions.
//
// Returns:
//   - PredicateFilter: A PredicateFilter function that checks if a element satisfies
//     all the PredicateFilter functions in funcs.
//
// Behavior:
//   - If no filter functions are provided, then all elements are considered to satisfy
//     the filter function.
//   - It returns false as soon as it finds a function in funcs that the element
//     does not satisfy.
func Intersect[T any](funcs ...PredicateFilter[T]) PredicateFilter[T] {
	if len(funcs) == 0 {
		return func(elem T) bool { return true }
	}

	return func(elem T) bool {
		for _, f := range funcs {
			ok := f(elem)
			if !ok {
				return false
			}
		}

		return true
	}
}

// ParallelIntersect returns a PredicateFilter function that checks if an element
// satisfies all the PredicateFilter functions in funcs concurrently.
//
// Parameters:
//   - funcs: A slice of PredicateFilter functions.
//
// Returns:
//   - PredicateFilter: A PredicateFilter function that checks if a element satisfies
//     all the PredicateFilter functions in funcs.
//
// Behavior:
//   - If no filter functions are provided, then all elements are considered to satisfy
//     the filter function.
//   - It returns false as soon as it finds a function in funcs that the element
//     does not satisfy.
func ParallelIntersect[T any](funcs ...PredicateFilter[T]) PredicateFilter[T] {
	if len(funcs) == 0 {
		return func(elem T) bool { return true }
	}

	return func(elem T) bool {
		resultChan := make(chan bool, len(funcs))

		for _, f := range funcs {
			go func(f PredicateFilter[T]) {
				resultChan <- f(elem)
			}(f)
		}

		for range funcs {
			if !<-resultChan {
				return false
			}
		}

		return true
	}
}

// Union returns a PredicateFilter function that checks if an element
// satisfies at least one of the PredicateFilter functions in funcs.
//
// Parameters:
//   - funcs: A slice of PredicateFilter functions.
//
// Returns:
//   - PredicateFilter: A PredicateFilter function that checks if a element satisfies
//     at least one of the PredicateFilter functions in funcs.
//
// Behavior:
//   - If no filter functions are provided, then no elements are considered to satisfy
//     the filter function.
//   - It returns true as soon as it finds a function in funcs that the element
//     satisfies.
func Union[T any](funcs ...PredicateFilter[T]) PredicateFilter[T] {
	if len(funcs) == 0 {
		return func(elem T) bool { return false }
	}

	return func(elem T) bool {
		for _, f := range funcs {
			ok := f(elem)
			if ok {
				return true
			}
		}

		return false
	}
}

// ParallelUnion returns a PredicateFilter function that checks if an element
// satisfies at least one of the PredicateFilter functions in funcs concurrently.
//
// Parameters:
//   - funcs: A slice of PredicateFilter functions.
//
// Returns:
//   - PredicateFilter: A PredicateFilter function that checks if a element satisfies
//     at least one of the PredicateFilter functions in funcs.
//
// Behavior:
//   - If no filter functions are provided, then no elements are considered to satisfy
//     the filter function.
//   - It returns true as soon as it finds a function in funcs that the element
//     satisfies.
func ParallelUnion[T any](funcs ...PredicateFilter[T]) PredicateFilter[T] {
	if len(funcs) == 0 {
		return func(elem T) bool { return false }
	}

	return func(elem T) bool {
		resultChan := make(chan bool, len(funcs))

		for _, f := range funcs {
			go func(f PredicateFilter[T]) {
				resultChan <- f(elem)
			}(f)
		}

		for range funcs {
			if <-resultChan {
				return true
			}
		}

		return false
	}
}

// SliceFilter is a function that iterates over the slice and applies the filter
// function to each element.
//
// Parameters:
//   - S: slice of elements.
//   - filter: function that takes an element and returns a bool.
//
// Returns:
//   - []T: slice of elements that satisfy the filter function.
//
// Behavior:
//   - If S is empty, the function returns a nil slice.
//   - If S has only one element and it satisfies the filter function, the function
//     returns a slice with that element. Otherwise, it returns a nil slice.
//   - An element is said to satisfy the filter function if the function returns true
//     when applied to the element.
//   - If the filter function is nil, the function returns the original slice.
func SliceFilter[T any](S []T, filter PredicateFilter[T]) []T {
	if len(S) == 0 {
		return nil
	} else if filter == nil {
		return S
	}

	var top int

	for i := 0; i < len(S); i++ {
		ok := filter(S[i])
		if ok {
			S[top] = S[i]
			top++
		}
	}

	return S[:top]
}

// FilterNilValues is a function that iterates over the slice and removes the
// nil elements.
//
// Parameters:
//   - S: slice of elements.
//
// Returns:
//   - []*T: slice of elements that satisfy the filter function.
//
// Behavior:
//   - If S is empty, the function returns a nil slice.
func FilterNilValues[T any](S []*T) []*T {
	if len(S) == 0 {
		return nil
	}

	var top int

	for i := 0; i < len(S); i++ {
		if S[i] != nil {
			S[top] = S[i]
			top++
		}
	}

	return S[:top]
}

// FilterNilPredicates is a function that iterates over the slice and removes the
// nil predicate functions.
//
// Parameters:
//   - S: slice of predicate functions.
//
// Returns:
//   - []PredicateFilter: slice of predicate functions that are not nil.
//
// Behavior:
//   - If S is empty, the function returns a nil slice.
func FilterNilPredicates[T any](S []PredicateFilter[T]) []PredicateFilter[T] {
	if len(S) == 0 {
		return nil
	}

	var top int

	for i := 0; i < len(S); i++ {
		if S[i] != nil {
			S[top] = S[i]
			top++
		}
	}

	return S[:top]
}

// SFSeparate is a function that iterates over the slice and applies the filter
// function to each element. The returned slices contain the elements that
// satisfy and do not satisfy the filter function.
//
// Parameters:
//   - S: slice of elements.
//   - filter: function that takes an element and returns a bool.
//
// Returns:
//   - []T: slice of elements that satisfy the filter function.
//   - []T: slice of elements that do not satisfy the filter function.
//
// Behavior:
//   - If S is empty, the function returns two empty slices.
func SFSeparate[T any](S []T, filter PredicateFilter[T]) ([]T, []T) {
	if len(S) == 0 {
		return []T{}, []T{}
	}

	var failed []T
	var top int

	for i := 0; i < len(S); i++ {
		ok := filter(S[i])
		if ok {
			S[top] = S[i]
			top++
		} else {
			failed = append(failed, S[i])
		}
	}

	return S[:top], failed
}

// SFSeparateEarly is a variant of SFSeparate that returns all successful elements.
// If there are none, it returns the original slice and false.
//
// Parameters:
//   - S: slice of elements.
//   - filter: function that takes an element and returns a bool.
//
// Returns:
//   - []T: slice of elements that satisfy the filter function or the original slice.
//   - bool: true if there are successful elements, otherwise false.
//
// Behavior:
//   - If S is empty, the function returns an empty slice and true.
func SFSeparateEarly[T any](S []T, filter PredicateFilter[T]) ([]T, bool) {
	if len(S) == 0 {
		return []T{}, true
	}

	var top int

	for i := 0; i < len(S); i++ {
		ok := filter(S[i])
		if ok {
			S[top] = S[i]
			top++
		}
	}

	if top == 0 {
		return S, false
	} else {
		return S[:top], true
	}
}

// RemoveEmpty is a function that removes the empty elements from a slice.
//
// Parameters:
//   - elems: The slice of elements.
//
// Returns:
//   - []T: The slice of elements without the empty elements.
func RemoveEmpty[T comparable](elems []T) []T {
	var top int

	for i := 0; i < len(elems); i++ {
		empty := *new(T)
		if elems[i] != empty {
			elems[top] = elems[i]
			top++
		}
	}

	return elems[:top]
}
