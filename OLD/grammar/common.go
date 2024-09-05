package grammar

import (
	"fmt"
)

// TokenTyper is an interface that defines the behavior of a token type.
//
// Value of 0 is reserved for the EOF token.
type TokenTyper interface {
	~int

	fmt.Stringer
	fmt.GoStringer
}
