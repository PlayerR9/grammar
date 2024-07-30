package main

func main() {

}

const templ string = `
// Node is a node in the AST.
type Node[N NodeTyper] struct {
	// Parent is the parent of the node.
	Parent *Node[N]

	// Children is the children of the node.
	Children []*Node[N]

	// Type is the type of the node.
	Type N

	// Data is the data of the node.
	Data string
}

// IsLeaf implements the grammar.Noder interface.
func (n *Node[N]) IsLeaf() bool {
	return len(n.Children) == 0
}

// Iterator implements the grammar.Noder interface.
func (n *Node[N]) Iterator() luc.Iterater[gr.Noder] {
	if len(n.Children) == 0 {
		return nil
	}

	nodes := make([]gr.Noder, 0, len(n.Children))
	for _, child := range n.Children {

		nodes = append(nodes, child)
	}

	return luc.NewSimpleIterator(nodes)
}

// String implements the grammar.Noder interface.
func (n *Node[N]) String() string {
	var builder strings.Builder

	builder.WriteString("Node[")
	builder.WriteString(n.Type.String())

	if n.Data != "" {
		builder.WriteString(" (")
		builder.WriteString(strconv.Quote(n.Data))
		builder.WriteRune(')')
	}

	builder.WriteRune(']')

	return builder.String()
}

// NewNode creates a new node.
//
// Parameters:
//   - t: The type of the node.
//   - data: The data of the node.
//
// Returns:
//   - *Node[N]: The new node. Never returns nil.
func NewNode[N NodeTyper](t N, data string) *Node[N] {
	return &Node[N]{
		Type: t,
		Data: data,
	}
}

// AppendChildren implements the *Node[N] interface.
func (n *Node[N]) AppendChildren(children []*Node[N]) {
	children = lus.FilterNilValues(children)
	if len(children) == 0 {
		return
	}

	for _, child := range children {
		child.Parent = n
	}

	n.Children = append(n.Children, children...)
}

// SetChildren sets the children of the node. Nil children are ignored.
func (n *Node[N]) SetChildren(children []*Node[N]) {
	children = lus.FilterNilValues(children)
	if len(children) == 0 {
		return
	}

	for _, child := range children {
		child.Parent = n
	}

	n.Children = children
}


// Result is the result of the AST.
type Result[N NodeTyper] struct {
	// nodes is the nodes of the result.
	nodes []*Node[N]
}

// NewResult creates a new AstResult.
//
// Returns:
//   - *AstResult[N]: The new AstResult. Never returns nil.
func NewResult[N NodeTyper]() *Result[N] {
	return &Result[N]{}
}

// MakeNode creates a new node and adds it to the result; replacing any existing nodes.
//
// Parameters:
//   - t: The type of the node.
//   - data: The data of the node.
func (a *Result[N]) MakeNode(t N, data string) {
	n := NewNode(t, data)

	a.nodes = []*Node[N]{n}
}

// SetNodes sets the nodes of the result. It ignores the nodes that are nil.
//
// Parameters:
//   - nodes: The nodes to set.
func (a *Result[N]) SetNodes(nodes []*Node[N]) {
	nodes = lus.FilterNilValues(nodes)
	if len(nodes) > 0 {
		a.nodes = nodes
	}
}

// AppendNodes appends the nodes of the result. It ignores the nodes that are nil.
//
// Parameters:
//   - nodes: The nodes to append.
func (a *Result[N]) AppendNodes(nodes []*Node[N]) {
	nodes = lus.FilterNilValues(nodes)

	if len(nodes) > 0 {
		a.nodes = append(a.nodes, nodes...)
	}
}

// AppendChildren appends the children to the node. It ignores the children that are nil.
//
// Parameters:
//   - children: The children to append.
//
// Returns:
//   - error: An error if the result is an error.
func (a *Result[N]) AppendChildren(children []*Node[N]) error {
	children = lus.FilterNilValues(children)

	if len(children) == 0 {
		return nil
	}

	if len(a.nodes) == 0 {
		return errors.New("no node to append children to")
	} else if len(a.nodes) > 1 {
		return errors.New("cannot append children to multiple nodes")
	}

	a.nodes[0].AppendChildren(children)

	return nil
}

// Apply applies the result.
//
// Returns:
//   - []*Node[N]: The nodes of the result.
func (a *Result[N]) Apply() []*Node[N] {
	return a.nodes
}

// DoFunc does something with the result.
//
// Parameters:
//   - f: The function to do something with the result.
//   - prev: The previous result of the function.
//
// Returns:
//   - any: The result of the function.
//   - error: An error if the function failed.
//
// Errors:
//   - *common.ErrInvalidParameter: If the f is nil.
//   - error: Any error returned by the f function.
func (a *Result[N]) DoFunc(f DoFunc[N], prev any) (any, error) {
	if f == nil {
		return nil, luc.NewErrNilParameter("f")
	}

	res, err := f(a, prev)
	if err != nil {
		return res, err
	}

	return res, nil
}

// TransformNodes transforms the nodes of the result.
//
// Parameters:
//   - new_type: The new type of the nodes.
//   - new_data: The new data of the nodes.
func (a *Result[N]) TransformNodes(new_type N, new_data string) {
	if len(a.nodes) == 0 {
		return
	}

	for _, node := range a.nodes {
		node.Type = new_type
		node.Data = new_data
	}
}

// Make is the constructor for the AST.
type Make[N NodeTyper, T gr.TokenTyper] struct {
	// ast_map is the map of the AST.
	ast_map map[T][]DoFunc[N]
}

// NewMake creates a new Make.
//
// Returns:
//   - *Make[N, T]: The new Make.
func NewMake[N NodeTyper, T gr.TokenTyper]() *Make[N, T] {
	return &Make[N, T]{
		ast_map: make(map[T][]DoFunc[N]),
	}
}

// AddEntry adds an entry to the AST. Nil steps are ignored.
//
// Parameters:
//   - t: The type of the entry.
//   - steps: The steps of the entry.
//
// Returns:
//   - error: An error if no steps were provided or if the entry already exists.
func (m *Make[N, T]) AddEntry(t T, steps []DoFunc[N]) error {
	if len(steps) == 0 {
		return errors.New("no steps provided")
	}

	var top int

	for i := 0; i < len(steps); i++ {
		if steps[i] != nil {
			steps[top] = steps[i]
			top++
		}
	}

	steps = steps[:top]

	if len(steps) == 0 {
		return errors.New("no steps provided")
	}

	if m.ast_map == nil {
		m.ast_map = make(map[T][]DoFunc[N])
	}

	_, ok := m.ast_map[t]
	if ok {
		return fmt.Errorf("entry with type %q already exists", t.String())
	}

	m.ast_map[t] = steps

	return nil
}

// Apply creates the AST given the root.
//
// Parameters:
//   - root: The root of the AST.
//
// Returns:
//   - []*Node[N]: The AST.
//   - error: An error if the AST could not be created.
func (m *Make[N, T]) Apply(root *gr.Token[T]) ([]*Node[N], error) {
	if root == nil {
		return nil, luc.NewErrNilParameter("root")
	}

	steps, ok := m.ast_map[root.Type]
	if !ok {
		return nil, fmt.Errorf("unexpected token type: %q", root.Type.String())
	}

	res := NewResult[N]()

	var prev any = root
	var err error

	for _, step := range steps {
		prev, err = step(res, prev)
		if err != nil {
			nodes := res.Apply()

			return nodes, fmt.Errorf("in %q: %w", root.Type.String(), err)
		}
	}

	if prev != nil {
		panic(luc.NewErrPossibleError(
			fmt.Errorf("last function returned (%v) instead of nil", prev),
			errors.New("you may have forgotten to specify a function"),
		))
	}

	nodes := res.Apply()

	return nodes, nil
}

// PrintAst stringifies the AST.
//
// Parameters:
//   - root: The root of the AST.
//
// Returns:
//   - string: The AST as a string.
func PrintAst[N NodeTyper](root *Node[N]) string {
	if root == nil {
		return ""
	}

	str, err := gr.PrintTree(root)
	luc.AssertErr(err, "Strings.PrintTree(root)")

	return str
}

// LeftAstFunc is a function that parses the left-recursive AST.
//
// Parameters:
//   - children: The children of the current node.
//
// Returns:
//   - []*Node[N]: The left-recursive AST.
//   - error: An error if the left-recursive AST could not be parsed.
type LeftAstFunc[N NodeTyper, T gr.TokenTyper] func(children []*gr.Token[T]) ([]*Node[N], error)

// LeftRecursive parses the left-recursive AST.
//
// Parameters:
//   - root: The root of the left-recursive AST.
//   - lhs_type: The type of the left-hand side.
//   - f: The function that parses the left-recursive AST.
//
// Returns:
//   - []*Node[N]: The left-recursive AST.
//   - error: An error if the left-recursive AST could not be parsed.
func LeftRecursive[N NodeTyper, T gr.TokenTyper](root *gr.Token[T], lhs_type T, f LeftAstFunc[N, T]) ([]*Node[N], error) {
	luc.AssertNil(root, "root")

	var nodes []*Node[N]

	for root != nil {
		if root.Type != lhs_type {
			return nodes, fmt.Errorf("expected %q, got %q instead", lhs_type.String(), root.Type.String())
		}

		children, err := ExtractChildren(root)
		if err != nil {
			return nodes, err
		} else if len(children) == 0 {
			return nodes, fmt.Errorf("expected at least 1 child, got 0 children instead")
		}

		last_child := children[len(children)-1]

		if last_child.Type == lhs_type {
			children = children[:len(children)-1]
			root = last_child
		} else {
			root = nil
		}

		sub_nodes, err := f(children)
		if len(sub_nodes) > 0 {
			nodes = append(nodes, sub_nodes...)
		}

		if err != nil {
			return nodes, fmt.Errorf("in %q: %w", root.Type.String(), err)
		}
	}

	return nodes, nil
}

// ToAstFunc is a function that parses the AST.
//
// Parameters:
//   - root: The root of the AST.
//
// Returns:
//   - []*Node[N]: The AST.
//   - error: An error if the AST could not be parsed.
type ToAstFunc[N NodeTyper, T gr.TokenTyper] func(root *gr.Token[T]) ([]*Node[N], error)

// ToAst parses the AST.
//
// Parameters:
//   - root: The root of the AST.
//   - to_ast: The function that parses the AST.
//
// Returns:
//   - []N: The AST.
//   - error: An error if the AST could not be parsed.
//
// Errors:
//   - *common.ErrInvalidParameter: If the root is nil or the to_ast is nil.
//   - error: Any error returned by the to_ast function.
func ToAst[N NodeTyper, T gr.TokenTyper](root *gr.Token[T], to_ast ToAstFunc[N, T]) ([]*Node[N], error) {
	if root == nil {
		return nil, luc.NewErrNilParameter("root")
	} else if to_ast == nil {
		return nil, luc.NewErrNilParameter("to_ast")
	}

	nodes, err := to_ast(root)
	if err != nil {
		return nodes, err
	}

	return nodes, nil
}

// ExtractData extracts the data from a token.
//
// Parameters:
//   - node: The token to extract the data from.
//
// Returns:
//   - string: The data of the token.
//   - error: An error if the data is not of type string or if the token is nil.
func ExtractData[T gr.TokenTyper](node *gr.Token[T]) (string, error) {
	if node == nil {
		return "", luc.NewErrNilParameter("node")
	}

	data, ok := node.Data.(string)
	if !ok {
		return "", fmt.Errorf("expected string, got %T instead", node.Data)
	}

	return data, nil
}

// ExtractChildren extracts the children from a token.
//
// Parameters:
//   - node: The token to extract the children from.
//
// Returns:
//   - []*gr.Token[T]: The children of the token.
//   - error: An error if the children is not of type []*gr.Token[T] or if the token is nil.
func ExtractChildren[T gr.TokenTyper](node *gr.Token[T]) ([]*gr.Token[T], error) {
	if node == nil {
		return nil, luc.NewErrNilParameter("node")
	}

	children, ok := node.Data.([]*gr.Token[T])
	if !ok {
		return nil, fmt.Errorf("expected []*Token, got %T instead", node.Data)
	}

	return children, nil
}

// DoFunc is a function that does something with the AST.
//
// Parameters:
//   - a: The result of the AST.
//   - prev: The previous result of the function.
//
// Returns:
//   - any: The result of the function.
//   - error: An error if the function failed.
type DoFunc[N NodeTyper] func(a *Result[N], prev any) (any, error)

// PartsBuilder is a builder for AST parts.
type PartsBuilder[N NodeTyper] struct {
	// parts is the parts of the builder.
	parts []DoFunc[N]
}

// NewPartsBuilder creates a new parts builder.
//
// Returns:
//   - *PartsBuilder[N]: The parts builder.
func NewPartsBuilder[N NodeTyper]() *PartsBuilder[N] {
	return &PartsBuilder[N]{
		parts: make([]DoFunc[N], 0),
	}
}

// Add adds a part to the builder. Does nothing if the part is nil.
//
// Parameters:
//   - f: The part to add.
func (a *PartsBuilder[N]) Add(f DoFunc[N]) {
	if f != nil {
		a.parts = append(a.parts, f)
	}
}

// Build builds the builder.
//
// Returns:
//   - []AstDoFunc[N]: The parts of the builder.
func (a *PartsBuilder[N]) Build() []DoFunc[N] {
	if len(a.parts) == 0 {
		return nil
	}

	steps := make([]DoFunc[N], len(a.parts))
	copy(steps, a.parts)

	return steps
}

// Reset resets the builder.
func (a *PartsBuilder[N]) Reset() {
	a.parts = a.parts[:0]
}

`
