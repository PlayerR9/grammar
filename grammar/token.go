package grammar

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// String implements the tree.Noder interface.
//
// Format:
//
//	"Token[T][{{ .Type }} ({{ .Data }})] : {{ .At }}]"
func (tn *Token[T]) String() string {
	var builder strings.Builder

	builder.WriteString("Token[T][")
	builder.WriteString(tn.Type.String())

	if tn.Data != "" {
		builder.WriteString(" (")
		builder.WriteString(strconv.Quote(tn.Data))
		builder.WriteRune(')')
	}

	builder.WriteString(" : ")
	builder.WriteString(strconv.Itoa(tn.At))
	builder.WriteRune(']')

	return builder.String()
}

// TokenTyper is an interface that defines the behavior of a token type.
//
// Value of 0 is reserved for the EOF token.
type TokenTyper interface {
	~int

	fmt.Stringer
	fmt.GoStringer
}

// Size returns the number of runes in the token's data.
//
// Returns:
//   - int: The number of runes in the token's data.
func (t *Token[T]) Size() int {
	if t.Data != "" {
		return utf8.RuneCountInString(t.Data)
	}

	var size int

	for c := t.FirstChild; c != nil; c = c.NextSibling {
		size += c.Size()
	}

	return size
}
