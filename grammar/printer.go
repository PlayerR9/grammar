package grammar

import (
	"fmt"
	"slices"
	"strings"

	luc "github.com/PlayerR9/lib_units/common"
)

// Noder is an interface that defines the behavior of a node.
type Noder interface {
	// IsLeaf returns true if the node is a leaf node.
	//
	// Returns:
	//   - bool: True if the node is a leaf node. False otherwise.
	IsLeaf() bool

	luc.Iterable[Noder]
	fmt.Stringer
}

// StackElement is a stack element.
type StackElement[T Noder] struct {
	// indent is the indentation.
	indent string

	// node is the node.
	node T

	// same_level is true if the node is on the same level.
	same_level bool

	// is_last is true if the node is the last node on the same level.
	is_last bool
}

// Printer is a tree printer.
type Printer[T Noder] struct {
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
func PrintTree[T Noder](root T) (string, error) {
	p := &Printer[T]{
		lines: make([]string, 0),
		seen:  make(map[Noder]bool),
	}

	se := &StackElement[T]{
		indent:     "",
		node:       root,
		same_level: false,
		is_last:    true,
	}

	stack := []*StackElement[T]{se}

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
func (p *Printer[T]) trav(elem *StackElement[T]) ([]*StackElement[T], error) {
	luc.AssertNil(elem, "elem")

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

	var elems []*StackElement[T]

	var indent strings.Builder

	indent.WriteString(elem.indent)

	if elem.same_level && !elem.is_last {
		indent.WriteString("│   ")
	} else {
		indent.WriteString("    ")
	}

	for {
		value, err := iter.Consume()
		ok := luc.IsDone(err)
		if ok {
			break
		} else if err != nil {
			return nil, err
		}

		node, ok := value.(T)
		luc.Assert(ok, fmt.Sprintf("expected %T, got %T", *new(T), value))

		se := &StackElement[T]{
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
