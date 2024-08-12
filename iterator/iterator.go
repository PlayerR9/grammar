package iterator

import (
	gcers "github.com/PlayerR9/go-commons/errors"
)

// IteratorFunc is the iterator function. The error ErrExausted signals the end of the iteration.
//
// Parameters:
//   - elem: The element to apply the function to.
//
// Returns:
//   - error: An error if the function failed.
type IteratorFunc func(elem any) error

// Iterater is an interface that defines the behavior of an iterator.
type Iterater interface {
	// Apply applies the iterator function on the current element.
	// The error io.EOF signals the successful end of the iteration.
	//
	// Parameters:
	//   - fn: The function to apply. Assumed to be non-nil.
	//
	// Returns:
	//   - error: An error if the function failed.
	//
	// Successful calls to Apply will also scan the next element.
	Apply(fn IteratorFunc) error

	// Reset resets the iterator. Used for initialization.
	Reset()
}

// Iterate applies the iterator function on the iterator.
// The error ErrExausted signals the end of the iteration.
//
// Parameters:
//   - it: The iterator. Assumed to be non-nil.
//   - fn: The function to apply. Assumed to be non-nil.
//
// Returns:
//   - error: An error if the function failed.
//
// Errors:
//   - *gcers.ErrNilParameter: If the iterator or function is nil.
//   - *gcint.ErrWhileAt: If the function failed at some point during the iteration.
func Iterate(it Iterater, fn IteratorFunc) error {
	if fn == nil {
		return gcers.NewErrNilParameter("fn")
	} else if it == nil {
		return gcers.NewErrNilParameter("it")
	}

	it.Reset()

	var err error
	idx := -1

	for err == nil {
		err = it.Apply(fn)
		idx++
	}

	if IsExausted(err) {
		return nil
	}

	return NewErrIteration(idx, err)
}
