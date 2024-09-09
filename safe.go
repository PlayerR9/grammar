package grammar

import (
	gcers "github.com/PlayerR9/go-commons/errors"
)

type Validater interface {
	Validate() error
}

func Validate(obj Validater, allow_nil bool) error {
	if obj == nil {
		if !allow_nil {
			return gcers.NewErrNilParameter("obj")
		} else {
			return nil
		}
	}

	err := obj.Validate()
	if err != nil {
		return err
	}

	return nil
}
