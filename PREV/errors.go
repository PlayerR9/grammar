package grammar

import (
	"iter"

	gcers "github.com/PlayerR9/go-commons/errors"
	uttr "github.com/PlayerR9/tree/tree"
)

type ErrParsing[T interface {
	Child() iter.Seq[T]
	BackwardChild() iter.Seq[T]
	Cleanup() []T
	LinkChildren(children []T)
	Copy() T

	uttr.Noder
}] struct {
	Err    error
	Forest []*uttr.Tree[T]
}

func (e ErrParsing[T]) Error() string {
	return gcers.Error(e.Err)
}

func (e ErrParsing[T]) Unwrap() error {
	return e.Err
}

func NewErrParsing[T interface {
	Child() iter.Seq[T]
	BackwardChild() iter.Seq[T]
	Cleanup() []T
	LinkChildren(children []T)
	Copy() T

	uttr.Noder
}](err error, forest []*uttr.Tree[T]) *ErrParsing[T] {
	return &ErrParsing[T]{
		Err:    err,
		Forest: forest,
	}
}
