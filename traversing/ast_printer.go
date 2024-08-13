package traversing

import (
	"strings"

	itr "github.com/PlayerR9/go-commons/iterator"
	dbg "github.com/PlayerR9/go-debug/assert"

	ustr "github.com/PlayerR9/grammar/util/strings"
)

// AstPrinter is a tree printer.
type AstPrinter struct {
	// lines is the list of lines.
	lines *ustr.LineBuffer

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
	p.lines.Reset()

	if len(p.seen) > 0 {
		for k := range p.seen {
			delete(p.seen, k)
		}
	}

	p.seen = make(map[TreeNoder]bool)
	p.indent = ""
	p.same_level = false
	p.is_last = true
}

// Apply implements the Traverser interface.
func (p *AstPrinter) Apply(node TreeNoder) ([]TravData, error) {
	dbg.AssertNotNil(p, "info")

	if p.indent != "" {
		p.lines.AddString(p.indent)

		ok := node.IsLeaf()
		if !ok || p.is_last {
			p.lines.AddString("└── ")
		} else {
			p.lines.AddString("├── ")
		}
	}

	// Prevent cycles.
	_, ok := p.seen[node]
	if ok {
		p.lines.AddString("... WARNING: Cycle detected!")
		p.lines.Accept()

		return nil, nil
	}

	p.lines.Accept()

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

	var children []TravData

	fn := func(elem any) error {
		value := dbg.AssertConv[TreeNoder](elem, "elem")

		td := TravData{
			Node: value,
			Data: &AstPrinter{
				lines:      p.lines,
				seen:       p.seen,
				indent:     indent.String(),
				same_level: false,
				is_last:    false,
			},
		}

		children = append(children, td)

		return nil
	}

	err := itr.Iterate(node.Iterator(), fn)
	if err != nil {
		return nil, err
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

// String implements the fmt.Stringer interface.
//
// Returns the printed tree as a string with newlines.
func (p *AstPrinter) String() string {
	return p.lines.String()
}
