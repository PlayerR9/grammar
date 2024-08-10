package grammar

import (
	"errors"
	"fmt"

	"github.com/PlayerR9/grammar/ast"
	"github.com/PlayerR9/grammar/grammar"
	"github.com/PlayerR9/grammar/lexing"
	"github.com/PlayerR9/grammar/parsing"
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
)

// Parser is the parser of the grammar.
type Parser[T ast.Noder, S grammar.TokenTyper] struct {
	// l is the lexer.
	l *lexing.Lexer[S]

	// p is the parser.
	p *parsing.Parser[S]

	// b is the ast builder.
	b *ast.Make[T, S]

	// debug is the debug setting.
	debug DebugSetting

	// data is the data of the parser.
	data []byte
}

// NewParser creates a new parser.
//
// Parameters:
//   - l: The lexer.
//   - p: The parser.
//   - b: The ast builder.
//
// Returns:
//   - *Parser: The new parser.
//
// This function returns nil iff any of the parameters is nil.
func NewParser[T ast.Noder, S grammar.TokenTyper](l *lexing.Lexer[S], p *parsing.Parser[S], b *ast.Make[T, S]) *Parser[T, S] {
	if l == nil || p == nil || b == nil {
		return nil
	}

	return &Parser[T, S]{
		l:     l,
		p:     p,
		b:     b,
		debug: ShowNone,
	}
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
//   - debug: The debug setting.
//
// Returns:
//   - *ast.Node[NodeType]: The AST tree.
//   - error: An error if the parsing failed.
func (p *Parser[T, S]) Parse(data []byte) (T, error) {
	if len(data) == 0 {
		return *new(T), errors.New("parameter (\"data\") is invalid: value must not be empty")
	}

	p.data = data

	if p.debug&ShowData != 0 {
		fmt.Println("Debug option show_data is enabled, printing data:")
		fmt.Println(string(p.data))
		fmt.Println()
	}

	tokens := lexing.FullLex(p.l, p.data)

	if p.debug&ShowLex != 0 {
		fmt.Println("Debug option show_lex is enabled, printing tokens:")
		for _, token := range tokens {
			fmt.Println("\t-", token.String())
		}
		fmt.Println()
	}

	if err := p.l.Error(); err != nil {
		return *new(T), err
	}

	forest, err := parsing.FullParse(p.p, tokens)

	if p.debug&ShowTree != 0 {
		fmt.Println("Debug option show_tree is enabled, printing forest:")

		for _, tree := range forest {
			fmt.Println(tree.String())
			fmt.Println()
		}

		fmt.Println()
	}

	if err != nil {
		return *new(T), fmt.Errorf("error while parsing: %w", err)
	} else if len(forest) != 1 {
		return *new(T), fmt.Errorf("expected 1 tree, got %d trees instead", len(forest))
	}

	nodes, err := p.b.Apply(forest[0])

	if p.debug&ShowAst != 0 {
		fmt.Println("Debug option show_ast is enabled, printing nodes:")

		for _, node := range nodes {
			fmt.Println(ast.PrintAst(node))
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
