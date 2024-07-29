package ast

// AstDoFunc is a function that does something with the AST.
//
// Parameters:
//   - a: The result of the AST.
//   - prev: The previous result of the function.
//
// Returns:
//   - any: The result of the function.
type AstDoFunc[N NodeTyper] func(a *AstResult[N], prev any) any

// AstPartsBuilder is a builder for AST parts.
type AstPartsBuilder[N NodeTyper] struct {
	// parts is the parts of the builder.
	parts []AstDoFunc[N]
}

// Add adds a part to the builder. Does nothing if the part is nil.
//
// Parameters:
//   - f: The part to add.
func (a *AstPartsBuilder[N]) Add(f AstDoFunc[N]) {
	if f != nil {
		a.parts = append(a.parts, f)
	}
}

// Build builds the builder.
//
// Returns:
//   - []AstDoFunc[N]: The parts of the builder.
func (a *AstPartsBuilder[N]) Build() []AstDoFunc[N] {
	return a.parts
}

// Reset resets the builder.
func (a *AstPartsBuilder[N]) Reset() {
	a.parts = a.parts[:0]
}
