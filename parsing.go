package grammar

import (
	"errors"
	"fmt"
	"iter"

	"github.com/PlayerR9/grammar/ast"
	"github.com/PlayerR9/grammar/grammar"
	"github.com/PlayerR9/grammar/lexing"
	"github.com/PlayerR9/grammar/parsing"

	dbg "github.com/PlayerR9/go-debug/assert"
)

// DebugSetting is the debug setting.
type DebugSetting int

const (
	// ShowNone shows no debug information.
	ShowNone DebugSetting = 0

	// ShowLex shows the lexer debug information.
	ShowLex DebugSetting = 1

	// ShowTree shows the parser debug information.
	ShowTree DebugSetting = 2

	// ShowAst shows the ast debug information.
	ShowAst DebugSetting = 4

	// ShowData shows the file reader debug information.
	ShowData DebugSetting = 8

	// ShowAll shows all debug information.
	ShowAll DebugSetting = ShowLex | ShowTree | ShowAst | ShowData

	// ShowParsing runs the parser step by step.
	ShowParsing DebugSetting = 16
)

// Parser is the parser of the grammar.
type Parser[T interface {
	comparable
	DirectChild() iter.Seq[T]
	BackwardChild() iter.Seq[T]
	ast.Noder
}, S grammar.TokenTyper] struct {
	// lexer is the lexer.
	lexer lexing.Lexer[S]

	// parser is the parser.
	parser parsing.Parser[S]

	// builder is the ast builder.
	builder ast.Make[T, S]

	// debug is the debug setting.
	debug DebugSetting
}

// Init initializes the parser with the given lexer, parser and builder.
//
// Parameters:
//   - l: The lexer.
//   - p: The parser.
//   - b: The ast builder.
//
// Also, by default no debug information is shown. Use SetDebug to show debug information.
func (p *Parser[T, S]) Init(lex lexing.Lexer[S], prs parsing.Parser[S], b ast.Make[T, S]) {
	p.lexer = lex
	p.parser = prs
	p.builder = b
	p.debug = ShowNone
}

// SetDebug sets the debug setting.
//
// Parameters:
//   - debug: The debug setting.
func (p *Parser[T, S]) SetDebug(debug DebugSetting) {
	p.debug = debug
}

// Parse parses the given data and returns the AST tree.
//
// Parameters:
//   - data: The data to parse.
//
// Returns:
//   - T: The AST tree.
//   - error: An error if the parsing failed.
func (p Parser[T, S]) Parse(data []byte) (T, error) {
	if len(data) == 0 {
		return *new(T), errors.New("parameter (\"data\") is invalid: value must not be empty")
	}

	if p.debug&ShowData != 0 {
		fmt.Println("Debug option show_data is enabled, printing data:")
		fmt.Println(string(data))
		fmt.Println()
	}

	tokens, err := p.lexer.FullLex(data)

	if p.debug&ShowLex != 0 {
		fmt.Println("Debug option show_lex is enabled, printing tokens:")
		for _, token := range tokens {
			fmt.Println("\t-", token.String())
		}
		fmt.Println()
	}

	if err != nil {
		return *new(T), err
	}

	var forest []*grammar.Token[S]

	if p.debug&ShowParsing != 0 {
		forest = p.parser.FullParseWithSteps(tokens, data, 3)
	} else {
		forest = p.parser.FullParse(tokens)
	}

	if p.debug&ShowTree != 0 {
		fmt.Println("Debug option show_tree is enabled, printing forest:")

		var p ast.AstPrinter[*grammar.Token[S]]

		for _, tree := range forest {
			dbg.AssertNotNil(tree, "tree")

			err := ast.Apply(&p, tree)
			dbg.AssertErr(err, "traversing.Apply(p, %s)", tree.String())

			fmt.Println(p.String())
			fmt.Println()
		}

		fmt.Println()
	}

	if p.parser.Err != nil {
		return *new(T), p.parser.Err
	} else if len(forest) == 0 {
		return *new(T), fmt.Errorf("expected at least 1 tree, got 0 trees instead")
	}

	nodes, err := p.builder.Apply(forest[0])

	if p.debug&ShowAst != 0 {
		fmt.Println("Debug option show_ast is enabled, printing nodes:")

		var pr ast.AstPrinter[T]

		for _, node := range nodes {
			dbg.AssertNotNil(node, "node")

			err := ast.Apply(&pr, node)
			dbg.AssertErr(err, "traversing.Apply(pr, %s)", node.String())

			fmt.Println(pr.String())
			fmt.Println()
		}

		fmt.Println()
	}

	if err != nil {
		return *new(T), fmt.Errorf("error while converting to AST: %w", err)
	} else if len(nodes) != 1 {
		return *new(T), fmt.Errorf("expected 1 node, got %d nodes instead", len(nodes))
	}

	return nodes[0], nil
}
