package traversing

import (
	"io"
	"strings"

	dbg "github.com/PlayerR9/go-debug/assert"
)

type TreeNoder interface {
	IsLeaf() bool
	String() string
	Iterator() Iterater
}

// AstPrinter is a tree printer.
type AstPrinter struct {
	// lines is the list of lines.
	lines []string

	// seen is the list of seen nodes.
	seen map[TreeNoder]bool

	// same_level is true if the node is on the same level.
	same_level bool

	// is_last is true if the node is the last node on the same level.
	is_last bool

	// indent is the indentation.
	indent string
}

// Reset implements the Traverser interface.
func (p *AstPrinter) Reset() {
	p.lines = p.lines[:0]

	for k := range p.seen {
		delete(p.seen, k)
	}

	p.seen = make(map[TreeNoder]bool)
	p.indent = ""
	p.same_level = false
	p.is_last = true
}

// Apply implements the Traverser interface.
func (p *AstPrinter) Apply(node TreeNoder) ([]TravData, error) {
	dbg.AssertNotNil(p, "info")

	var builder strings.Builder

	if p.indent != "" {
		builder.WriteString(p.indent)

		ok := node.IsLeaf()
		if !ok || p.is_last {
			builder.WriteString("└── ")
		} else {
			builder.WriteString("├── ")
		}
	}

	// Prevent cycles.
	_, ok := p.seen[node]
	if ok {
		builder.WriteString("... WARNING: Cycle detected!")

		p.lines = append(p.lines, builder.String())

		return nil, nil
	}

	builder.WriteString(node.String())
	p.lines = append(p.lines, builder.String())
	p.seen[node] = true

	var indent strings.Builder

	indent.WriteString(p.indent)

	if p.same_level && !p.is_last {
		indent.WriteString("│   ")
	} else {
		indent.WriteString("    ")
	}

	p.indent = indent.String()
	p.same_level = false
	p.is_last = false

	iter := node.Iterator()
	dbg.AssertNotNil(iter, "iter")

	for {
		err := iter.Consume()
		if err == io.EOF {
			break
		}

		dbg.AssertErr(err, "iter.Consume()")

		td := TravData{
			Node: value,
			Data: &AstPrinter{
				indent:     indent.String(),
				same_level: false,
				is_last:    false,
			},
		}

		children = append(children, td)
	}

	if len(children) == 0 {
		return nil, nil
	}

	if len(children) >= 2 {
		for _, c := range children {
			data := dbg.AssertConv[*AstPrinter](c.Data, "c.Data")

			data.same_level = true
		}
	}

	last_child := children[len(children)-1].Data

	tmp := dbg.AssertConv[*AstPrinter](last_child, "last_child")

	tmp.is_last = true

	return children, nil
}

// Copy creates a copy of the Printer. This method returns the Printer itself.
//
// Returns:
//   - *Printer: The copy.
func (p *AstPrinter) Copy() Traverser {
	return p
}

// String implements the fmt.Stringer interface.
//
// Returns the printed tree as a string with newlines.
func (p *AstPrinter) String() string {
	return strings.Join(p.lines, "\n")
}
