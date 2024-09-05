package grammar

import (
	"fmt"
	"iter"

	gcers "github.com/PlayerR9/go-commons/errors"
	uttr "github.com/PlayerR9/go-commons/tree"
	dbp "github.com/PlayerR9/go-debug/debug"
	ast "github.com/PlayerR9/grammar/ast"
	"github.com/PlayerR9/grammar/internal"
	lxr "github.com/PlayerR9/grammar/lexer"
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
}](lexer *lxr.Lexer[T], parser *Parser[T], builder *ast.AstBuilder[T, N]) (ParsingFunc[N], error) {
	if lexer == nil {
		return nil, gcers.NewErrNilParameter("lexer")
	} else if parser == nil {
		return nil, gcers.NewErrNilParameter("parser")
	} else if builder == nil {
		return nil, gcers.NewErrNilParameter("builder")
	}

	fn := func(data []byte, opt DebugSetting) (*uttr.Tree[N], error) {
		lexer.SetInputStream(data)

		err := lexer.Lex()

		tokens := lexer.Tokens()

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

		if err != nil {
			return nil, fmt.Errorf("error lexing: %w", err)
		}

		var eos gcers.ErrOrSol[error]

		next, stop := iter.Pull(parser.Parse(tokens))
		defer stop()

		for {
			parsed, ok := next()
			if !ok {
				break
			}

			err := parsed.Error()
			if err != nil {
				eos.AddSol(fmt.Errorf("error parsing: %w", err), 0)

				continue
			}

			forest := parsed.Forest()

			if opt&ShowTree != 0 {
				// DEBUG: Print the parse tree
				dbp.Print("Parse Tree:", func() []string {
					strs := make([]string, 0, len(forest))

					for _, tree := range forest {
						strs = append(strs, tree.GoString())
					}

					return strs
				})
			}

			if len(forest) == 0 {
				eos.AddSol(fmt.Errorf("empty parse tree"), 1)

				continue
			} else if len(forest) > 1 {
				eos.AddSol(fmt.Errorf("too many parse trees"), 1)

				continue
			}

			tree, err := builder.Make(forest[0])

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
				eos.AddSol(fmt.Errorf("error building AST: %w", err), 2)

				continue
			}

			return tree, nil
		}

		errs := eos.Solutions()
		if len(errs) == 0 {
			return nil, fmt.Errorf("no parse tree found")
		} else {
			return nil, errs[0]
		}
	}

	return fn, nil
}
