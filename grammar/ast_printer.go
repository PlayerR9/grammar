package grammar

import (
	"iter"
	"strings"

	gcstr "github.com/PlayerR9/go-commons/strings"
)

type TokenTreePrinter[N interface {
	comparable
	DirectChild() iter.Seq[N]
	IsLeaf() bool
	String() string
}] struct {
	lines gcstr.LineBuffer
}

func NewTokenTreePrinter[N interface {
	comparable
	DirectChild() iter.Seq[N]
	IsLeaf() bool
	String() string
}]() *TokenTreePrinter[N] {
	return &TokenTreePrinter[N]{}
}

func (p *TokenTreePrinter[N]) String() string {
	type Printer struct {
		// lines is the list of lines.
		lines *gcstr.LineBuffer

		// seen is the list of seen nodes.
		seen map[N]bool

		// same_level is true if the node is on the same level.
		same_level bool

		// is_last is true if the node is the last node on the same level.
		is_last bool

		// indent is the indentation.
		indent string
	}

	type TravData struct {
		// Node is the node.
		Node N

		// Data is the data associated with the node before the node is visited.
		Data *Printer
	}

	sp := &Printer{
		lines:      &p.lines,
		seen:       make(map[N]bool),
		is_last:    true,
		indent:     "",
		same_level: false,
	}

	stack := []TravData{{root, sp}}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		children, err := func(node N) ([]TravData, error) {
			if sp.indent != "" {
				p.lines.AddString(sp.indent)

				ok := node.IsLeaf()
				if !ok || sp.is_last {
					p.lines.AddString("└── ")
				} else {
					p.lines.AddString("├── ")
				}
			}

			// Prevent cycles.
			_, ok := sp.seen[node]
			if ok {
				p.lines.AddString("... WARNING: Cycle detected!")
				p.lines.Accept()

				return nil, nil
			}

			p.lines.AddString(node.String())
			p.lines.Accept()

			sp.seen[node] = true

			var indent strings.Builder

			indent.WriteString(sp.indent)

			if sp.same_level && !sp.is_last {
				indent.WriteString("│   ")
			} else {
				indent.WriteString("    ")
			}

			sp.indent = indent.String()
			sp.same_level = false
			sp.is_last = false

			var children []TravData

			for c := range node.DirectChild() {
				td := TravData{
					Node: c,
					Data: &Printer{
						lines:      sp.lines,
						seen:       sp.seen,
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
					c.Data.same_level = true
				}
			}

			last_child := children[len(children)-1].Data

			last_child.is_last = true

			return children, nil
		}(top.Node)
		if err != nil {
			return err
		}

		for i := len(children) - 1; i >= 0; i-- {
			stack = append(stack, children[i])
		}
	}

	return nil
}
