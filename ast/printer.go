package ast

import (
	"slices"
	"strings"

	itrs "github.com/PlayerR9/iterators/simple"
)

// stack_element is a stack element.
type stack_element[N Noder] struct {
	// indent is the indentation.
	indent string

	// node is the node.
	node N

	// same_level is true if the node is on the same level.
	same_level bool

	// is_last is true if the node is the last node on the same level.
	is_last bool
}

// Printer is a tree printer.
type Printer[N Noder] struct {
	// lines is the list of lines.
	lines []string

	// seen is the list of seen nodes.
	seen map[Noder]bool
}

// PrintTree prints the tree.
//
// Parameters:
//   - root: The root node.
//
// Returns:
//   - string: The tree as a string.
//   - error: An error if printing fails.
func PrintTree[N Noder](root N) (string, error) {
	p := &Printer[N]{
		lines: make([]string, 0),
		seen:  make(map[Noder]bool),
	}

	se := &stack_element[N]{
		indent:     "",
		node:       root,
		same_level: false,
		is_last:    true,
	}

	stack := []*stack_element[N]{se}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		sub, err := p.trav(top)
		if err != nil {
			return "", err
		} else if len(sub) == 0 {
			continue
		}

		slices.Reverse(sub)

		stack = append(stack, sub...)
	}

	return strings.Join(p.lines, "\n"), nil
}

// trav traverses the tree.
//
// Parameters:
//   - elem: The stack element.
//
// Returns:
//   - []*StackElement: The list of stack elements.
//   - error: An error if traversing fails.
func (p *Printer[N]) trav(elem *stack_element[N]) ([]*stack_element[N], error) {
	// luc.AssertNil(elem, "elem")

	var builder strings.Builder

	if elem.indent != "" {
		builder.WriteString(elem.indent)

		ok := elem.node.IsLeaf()
		if !ok || elem.is_last {
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

		return nil, nil
	}

	builder.WriteString(elem.node.String())

	p.lines = append(p.lines, builder.String())

	p.seen[elem.node] = true

	iter := elem.node.Iterator()
	if iter == nil {
		return nil, nil
	}

	var elems []*stack_element[N]

	var indent strings.Builder

	indent.WriteString(elem.indent)

	if elem.same_level && !elem.is_last {
		indent.WriteString("│   ")
	} else {
		indent.WriteString("    ")
	}

	for {
		value, err := iter.Consume()
		if err == itrs.Exhausted {
			break
		} else if err != nil {
			return nil, err
		}

		node := value.(N)
		// node := luc.AssertConv[N](value, "value")

		se := &stack_element[N]{
			indent:     indent.String(),
			node:       node,
			same_level: false,
			is_last:    false,
		}

		elems = append(elems, se)
	}

	if len(elems) == 0 {
		return nil, nil
	}

	if len(elems) >= 2 {
		for _, e := range elems {
			e.same_level = true
		}
	}

	elems[len(elems)-1].is_last = true

	return elems, nil
}
