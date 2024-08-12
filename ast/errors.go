package ast

import "fmt"

type ErrInvalidType[gr.] struct {
	Expected T
}

func (e *ErrInvalidType) Error() string {
	if tk == nil {
		return fmt.Errorf("expected %q, got nil instead", tk_type.String())
	}

	if tk.Type != tk_type {
		return fmt.Errorf("expected %q, got %q instead", tk_type.String(), tk.Type.String())
	}
}
