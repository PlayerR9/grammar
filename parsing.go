package grammar

import (
	"fmt"
	"iter"

	gcers "github.com/PlayerR9/go-commons/errors"
	gchlp "github.com/PlayerR9/go-commons/helpers"
	uttr "github.com/PlayerR9/go-commons/tree"
	dbp "github.com/PlayerR9/go-debug/debug"
	ast "github.com/PlayerR9/grammar/ast"
	gr "github.com/PlayerR9/grammar/grammar"
	"github.com/PlayerR9/grammar/internal"
	lxr "github.com/PlayerR9/grammar/lexer"
	prx "github.com/PlayerR9/grammar/parser"
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

// ParsingFunc is the parsing function.
//
// Parameters:
//   - data: The data.
//   - opt: The debug option.
//
// Returns:
//   - *uttr.Tree[N]: The tree.
//   - error: An error.
type ParsingFunc[N interface {
	Child() iter.Seq[N]
	BackwardChild() iter.Seq[N]

	uttr.TreeNoder
}] func(data []byte, opt DebugSetting) (*uttr.Tree[N], error)

// Noder is the interface for nodes.
type Noder interface {
}

// ParseData returns a function that parses data.
//
// Parameters:
//   - lexer: The lexer.
//   - parser: The parser.
//   - builder: The builder.
//
// Returns:
//   - ParsingFunc[N]: The function that parses data. Never returns nil.
func ParseData[T internal.TokenTyper, N interface {
	Child() iter.Seq[N]
	BackwardChild() iter.Seq[N]

	uttr.TreeNoder
}](lexer *lxr.Lexer[T], parser *prx.Parser[T], builder *ast.AstBuilder[T, N]) (ParsingFunc[N], error) {
	if lexer == nil {
		return nil, gcers.NewErrNilParameter("lexer")
	} else if parser == nil {
		return nil, gcers.NewErrNilParameter("parser")
	} else if builder == nil {
		return nil, gcers.NewErrNilParameter("builder")
	}

	fn := func(data []byte, opt DebugSetting) (*uttr.Tree[N], error) {
		err := lexer.SetInputStream(data)
		if err != nil {
			return nil, fmt.Errorf("error setting input stream: %w", err)
		}

		// var eor gcers.ErrOrSol[]

		for lexed := range lexer.Lex() {
			err := lexed.Error()
			if err != nil {
				return nil, fmt.Errorf("error lexing: %w", err)
			}

			tokens := lexed.Tokens()

			if opt&ShowLex != 0 {
				// DEBUG: Print the tokens
				dbp.Print("Tokens:", func() []string {
					strs := make([]string, 0, len(tokens)+1)

					for _, token := range tokens {
						strs = append(strs, token.GoString())
					}

					return strs
				})
			}

			// if opt&ShowTree != 0 {
			// 	// DEBUG: Print the parse tree
			// 	dbp.Print("Parse Tree:", func() []string {
			// 		strs := make([]string, 0, len(forest))

			// 		for _, tree := range forest {
			// 			strs = append(strs, tree.GoString())
			// 		}

			// 		return strs
			// 	})
			// }

			tree, err := FullParse(parser, builder, tokens)

			if opt&ShowAst != 0 {
				// DEBUG: Print the AST
				dbp.Print("AST:", func() []string {
					strs := make([]string, 0, 1)

					if tree != nil {
						strs = append(strs, tree.GoString())
					}

					return strs
				})
			}

			if err != nil {
				return nil, err
			}

			return tree, nil
		}
	}

	return fn, nil
}

type Trees[T internal.TokenTyper, N interface {
	Child() iter.Seq[N]
	BackwardChild() iter.Seq[N]

	uttr.TreeNoder
}] struct {
	ParsingTree *uttr.Tree[*gr.Token[T]]
	AST         *uttr.Tree[N]
}

func FullParse[T internal.TokenTyper, N interface {
	Child() iter.Seq[N]
	BackwardChild() iter.Seq[N]

	uttr.TreeNoder
}](parser *prx.Parser[T], builder *ast.AstBuilder[T, N], tokens []*gr.Token[T]) (*uttr.Tree[*gr.Token[T]], *uttr.Tree[N], error) {
	var eos_parsed []*gchlp.SimpleHelper[[]*uttr.Tree[*gr.Token[T]]]
	var eos_ast []*gchlp.SimpleHelper[*uttr.Tree[N]]
	var found_ast bool

	for parsed := range parser.Parse(tokens) {
		err := parsed.Error()
		forest := parsed.Forest()

		if err != nil {
			if !found_ast {
				h := gchlp.NewSimpleHelper(forest, fmt.Errorf("error parsing: %w", err))
				eos_parsed = append(eos_parsed, h)
			}

			break
		}

		if len(forest) != 1 {
			if !found_ast {
				if len(forest) == 0 {
					h := gchlp.NewSimpleHelper[[]*uttr.Tree[*gr.Token[T]]](nil, fmt.Errorf("empty parse tree"))
					eos_parsed = append(eos_parsed, h)
				} else {
					h := gchlp.NewSimpleHelper(forest, fmt.Errorf("too many parse trees"))
					eos_parsed = append(eos_parsed, h)
				}
			}

			continue
		}

		tree, err := builder.Make(forest[0])
		if err != nil {
			if !found_ast {
				h_ast := gchlp.NewSimpleHelper(tree, fmt.Errorf("error building AST: %w", err))
				eos_ast = []*gchlp.SimpleHelper[*uttr.Tree[N]]{h_ast}

				h_parsed := gchlp.NewSimpleHelper(forest, nil)
				eos_parsed = []*gchlp.SimpleHelper[[]*uttr.Tree[*gr.Token[T]]]{h_parsed}

				found_ast = true
			} else {
				h_ast := gchlp.NewSimpleHelper(tree, fmt.Errorf("error building AST: %w", err))
				eos_ast = append(eos_ast, h_ast)

				h_parsed := gchlp.NewSimpleHelper(forest, nil)
				eos_parsed = append(eos_parsed, h_parsed)
			}

			continue
		}

		return forest[0], tree, nil
	}

	errs := eos.Solutions()
	if len(errs) == 0 {
		return nil, fmt.Errorf("no parse tree found")
	} else {
		return nil, errs[0]
	}
}
