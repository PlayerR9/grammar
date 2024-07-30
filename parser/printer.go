package parser

import (
	gr "github.com/PlayerR9/grammar/grammar"
	luc "github.com/PlayerR9/lib_units/common"
	tr "github.com/PlayerR9/tree/tree"
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

	str, err := tr.PrintTree(root)
	luc.AssertErr(err, "grammar.PrintTree(root)")

	return str
}
