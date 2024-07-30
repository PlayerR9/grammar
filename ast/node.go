package ast

import (
	"fmt"
	"strconv"
	"strings"
)

// String implements the tree.Noder interface.
func (tn *Node[N]) String() string {
	var builder strings.Builder

	builder.WriteString("Node[")
	builder.WriteString(tn.Type.String())

	if tn.Data != "" {
		builder.WriteString(" (")
		builder.WriteString(strconv.Quote(tn.Data))
		builder.WriteRune(')')
	}

	builder.WriteRune(']')

	return builder.String()
}

// NodeTyper is an interface that defines the behavior of a node type.
type NodeTyper interface {
	~int

	fmt.Stringer
}

// AstNoder is an interface that defines the behavior of an AST node.
type AstNoder interface{}
