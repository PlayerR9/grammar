package grammar

import (
	"iter"

	gcers "github.com/PlayerR9/go-commons/errors"
	uttr "github.com/PlayerR9/go-commons/tree"
)

type ErrParsing[T interface {
	Child() iter.Seq[T]
	BackwardChild() iter.Seq[T]

	uttr.TreeNoder
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

	uttr.TreeNoder
}](err error, forest []*uttr.Tree[T]) *ErrParsing[T] {
	return &ErrParsing[T]{
		Err:    err,
		Forest: forest,
	}
}
