package pkg

import (
	"fmt"
	"log"

	ggen "github.com/PlayerR9/lib_units/generator"
)

var (
	// Logger is the logger.
	Logger *log.Logger
)

func init() {
	Logger = ggen.InitLogger("node")
}

type GenData struct {
	PackageName string

	TypeName string

	NodeName string
	NodeSig  string

	IteratorName string
	IteratorSig  string

	Generics string

	Noder string
}

// SetPackageName implements the generator.Generater interface.
func (gd *GenData) SetPackageName(pkg_name string) {
	gd.PackageName = pkg_name
}

var (
	Generator *ggen.CodeGenerator[*GenData]
)

func init() {
	tmp, err := ggen.NewCodeGeneratorFromTemplate[*GenData]("", templ)
	if err != nil {
		Logger.Fatalf("Failed to create code generator: %s", err.Error())
	}

	tmp.AddDoFunc(func(gd *GenData) error {
		sig, err := ggen.MakeTypeSig(gd.NodeName, "")
		if err != nil {
			return fmt.Errorf("failed to make type sig: %w", err)
		}

		gd.NodeSig = sig

		return nil
	})

	tmp.AddDoFunc(func(gd *GenData) error {
		gd.Generics = ggen.GenericsSigFlag.String()

		return nil
	})

	tmp.AddDoFunc(func(gd *GenData) error {
		sig, err := ggen.MakeTypeSig(gd.NodeName, "Iterator")
		if err != nil {
			return fmt.Errorf("failed to make iterator sig: %w", err)
		}

		gd.IteratorSig = sig

		return nil
	})

	tmp.AddDoFunc(func(gd *GenData) error {
		gd.IteratorName = gd.NodeName + "Iterator"

		return nil
	})

	tmp.AddDoFunc(func(gd *GenData) error {
		if gd.PackageName == "ast" {
			gd.Noder = "Noder"
		} else {
			gd.Noder = "ast.Noder"
		}

		return nil
	})

	Generator = tmp
}

// templ is the template for the ast node.
const templ = `// Code generated by go generate; EDIT THIS FILE DIRECTLY
package {{ .PackageName }}

import (
	"strconv"
	"strings"

	{{ if ne .PackageName "ast" }}"github.com/PlayerR9/grammar/ast"{{ end }}
	"github.com/PlayerR9/lib_units/common"
)

// {{ .IteratorName }} is a pull-based iterator that iterates
// over the children of a {{ .NodeName }}.
type {{ .IteratorName }}{{ .Generics }} struct {
	parent, current *{{ .NodeSig }}
}

// Consume implements the common.Iterater interface.
//
// The only error type that can be returned by this function is the *common.ErrExhaustedIter type.
//
// Moreover, the return value is always of type *{{ .NodeSig }} and never nil; unless the iterator
// has reached the end of the branch.
func (iter *{{ .IteratorSig }}) Consume() ({{ .Noder }}, error) {
	if iter.current == nil {
		return nil, common.NewErrExhaustedIter()
	}

	node := iter.current
	iter.current = iter.current.NextSibling

	return node, nil
}

// Restart implements the common.Iterater interface.
func (iter *{{ .IteratorSig }}) Restart() {
	iter.current = iter.parent.FirstChild
}

// {{ .NodeName }} is a node in a ast.
type {{ .NodeName }}{{ .Generics }} struct {
	Parent, FirstChild, NextSibling, LastChild, PrevSibling *{{ .NodeSig }}

	Type {{ .TypeName }}
	Data string
}

// IsLeaf implements the {{ .Noder }} interface.
func (tn *{{ .NodeSig }}) IsLeaf() bool {
	return tn.FirstChild == nil
}

// AddChild implements the {{ .Noder }} interface.
func (tn *{{ .NodeSig }}) AddChild(target {{ .Noder }}) {
	if target == nil {
		return
	}

	tmp, ok := target.(*{{ .NodeSig }})
	if !ok {
		return
	}
	
	tmp.NextSibling = nil
	tmp.PrevSibling = nil

	last_child := tn.LastChild

	if last_child == nil {
		tn.FirstChild = tmp
	} else {
		last_child.NextSibling = tmp
		tmp.PrevSibling = last_child
	}

	tmp.Parent = tn
	tn.LastChild = tmp
}

// AddChildren implements the {{ .Noder }} interface.
func (tn *{{ .NodeSig }}) AddChildren(children []{{ .Noder }}) {
	if len(children) == 0 {
		return
	}
	
	var valid_children []*{{ .NodeSig }}

	for _, child := range children {
		if child == nil {
			continue
		}

		c, ok := child.(*{{ .NodeSig }})
		if !ok {
			continue
		}

		valid_children = append(valid_children, c)
	}

	if len(valid_children) == 0 {
		return
	}

	// Deal with the first child
	first_child := valid_children[0]

	first_child.NextSibling = nil
	first_child.PrevSibling = nil

	last_child := tn.LastChild

	if last_child == nil {
		tn.FirstChild = first_child
	} else {
		last_child.NextSibling = first_child
		first_child.PrevSibling = last_child
	}

	first_child.Parent = tn
	tn.LastChild = first_child

	// Deal with the rest of the children
	for i := 1; i < len(valid_children); i++ {
		child := valid_children[i]

		child.NextSibling = nil
		child.PrevSibling = nil

		last_child := tn.LastChild
		last_child.NextSibling = child
		child.PrevSibling = last_child

		child.Parent = tn
		tn.LastChild = child
	}
}

// Iterator implements the {{ .Noder }} interface.
//
// This function returns an iterator that iterates over the direct children of the node.
// Implemented as a pull-based iterator, this function never returns nil and any of the
// values is guaranteed to be a non-nil node of type {{ .NodeSig }}.
func (tn *{{ .NodeSig }}) Iterator() common.Iterater[{{ .Noder }}] {
	return &{{ .IteratorSig }}{
		parent: tn,
		current: tn.FirstChild,
	}
}

// String implements the {{ .Noder }} interface.
func (tn *{{ .NodeSig }}) String() string {
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

// New{{ .NodeName }} creates a new node with the given data.
//
// Parameters:
//   - n_type: The type of the node.
//   - data: The data of the node.
//
// Returns:
//   - *{{ .NodeSig }}: A pointer to the newly created node. It is
//   never nil.
func New{{ .NodeName }}{{ .Generics }}(n_type {{ .TypeName }}, data string) *{{ .NodeSig }} {
	return &{{ .NodeSig }}{
		Type: n_type,
		Data: data,
	}
}
`
