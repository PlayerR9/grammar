package traversing

import (
	"strings"
)

// stack_element is a stack element.
type stack_element[S TokenTyper] struct {
	// indent is the indentation.
	indent string

	// node is the node.
	node *Token[S]

	// same_level is true if the node is on the same level.
	same_level bool

	// is_last is true if the node is the last node on the same level.
	is_last bool
}

// token_printer is a tree printer.
type token_printer[S TokenTyper] struct {
	// lines is the list of lines.
	lines []string

	// seen is the list of seen nodes.
	seen map[*Token[S]]bool
}

// trav traverses the tree.
//
// Parameters:
//   - elem: The stack element.
//
// Returns:
//   - []*StackElement: The list of stack elements.
func (p *token_printer[S]) trav(elem *stack_element[S]) []*stack_element[S] {
	// luc.AssertNil(elem, "elem")

	var builder strings.Builder

	if elem.indent != "" {
		builder.WriteString(elem.indent)

		if elem.node.FirstChild != nil || elem.is_last {
			builder.WriteString("└── ")
		} else {
			builder.WriteString("├── ")
		}
	}

	// Prevent cycles.
	_, ok := p.seen[elem.node]
	if ok {
		builder.WriteString("... WARNING: Cycle detected!")

		p.lines = append(p.lines, builder.String())

		return nil
	}

	builder.WriteString(elem.node.String())

	p.lines = append(p.lines, builder.String())

	p.seen[elem.node] = true

	if elem.node.FirstChild == nil {
		return nil
	}

	var elems []*stack_element[S]

	var indent strings.Builder

	indent.WriteString(elem.indent)

	if elem.same_level && !elem.is_last {
		indent.WriteString("│   ")
	} else {
		indent.WriteString("    ")
	}

	for c := elem.node.FirstChild; c != nil; c = c.NextSibling {
		se := &stack_element[S]{
			indent:     indent.String(),
			node:       c,
			same_level: false,
			is_last:    false,
		}

		elems = append(elems, se)
	}

	// luc.Assert(len(elems) > 0, "len(elems) > 0")

	if len(elems) >= 2 {
		for _, e := range elems {
			e.same_level = true
		}
	}

	elems[len(elems)-1].is_last = true

	return elems
}
