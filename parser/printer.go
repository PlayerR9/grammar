package parser

import (
	gr "github.com/PlayerR9/grammar/grammar"
	luc "github.com/PlayerR9/lib_units/common"
)

// PrintParseTree is a helper function that prints the parse tree.
//
// Parameters:
//   - root: The root of the parse tree.
//
// Returns:
//   - string: The parse tree data.
func PrintParseTree[T gr.TokenTyper](root *gr.Token[T]) string {
	if root == nil {
		return ""
	}

	str, err := gr.PrintTree(root)
	luc.AssertErr(err, "grammar.PrintTree(root)")

	return str
}
