package ast

// DoFunc is a function that does something with the AST.
//
// Parameters:
//   - a: The result of the AST.
//   - prev: The previous result of the function.
//
// Returns:
//   - any: The result of the function.
//   - error: An error if the function failed.
type DoFunc[N Noder] func(a *Result[N], prev any) (any, error)

// PartsBuilder is a builder for AST parts.
type PartsBuilder[N Noder] struct {
	// parts is the parts of the builder.
	parts []DoFunc[N]
}

// NewPartsBuilder creates a new parts builder.
//
// Returns:
//   - PartsBuilder[N]: The parts builder.
func NewPartsBuilder[N Noder]() PartsBuilder[N] {
	return PartsBuilder[N]{
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
func (a PartsBuilder[N]) Build() []DoFunc[N] {
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
