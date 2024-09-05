package grammar

import (
	"fmt"
	"iter"

	gcers "github.com/PlayerR9/go-commons/errors"
	uttr "github.com/PlayerR9/go-commons/tree"
	dbp "github.com/PlayerR9/go-debug/debug"
	"github.com/PlayerR9/grammar/grammar/internal"
)

type Noder interface {
}

func ParseData[T internal.TokenTyper, N interface {
	Child() iter.Seq[N]
	BackwardChild() iter.Seq[N]

	uttr.TreeNoder
}](lexer Lexer[T], parser *Parser[T], builder *AstBuilder[T, N]) func(data []byte) (*uttr.Tree[N], error) {
	if parser == nil {
		return func(data []byte) (*uttr.Tree[N], error) {
			return nil, gcers.NewErrNilParameter("parser")
		}
	} else if builder == nil {
		return func(data []byte) (*uttr.Tree[N], error) {
			return nil, gcers.NewErrNilParameter("builder")
		}
	}

	return func(data []byte) (*uttr.Tree[N], error) {
		lexer.SetInputStream(data)

		err := lexer.Lex()

		tokens := lexer.Tokens()

		// DEBUG: Print the tokens
		dbp.Print("Tokens:", func() []string {
			strs := make([]string, 0, len(tokens)+1)

			for _, token := range tokens {
				strs = append(strs, token.GoString())
			}

			return strs
		})

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

			// DEBUG: Print the parse tree
			dbp.Print("Parse Tree:", func() []string {
				strs := make([]string, 0, len(forest))

				for _, tree := range forest {
					strs = append(strs, tree.GoString())
				}

				return strs
			})

			if len(forest) == 0 {
				eos.AddSol(fmt.Errorf("empty parse tree"), 1)

				continue
			} else if len(forest) > 1 {
				eos.AddSol(fmt.Errorf("too many parse trees"), 1)

				continue
			}

			tree, err := builder.Make(forest[0])

			// DEBUG: Print the AST
			dbp.Print("AST:", func() []string {
				strs := make([]string, 0, 1)

				if tree != nil {
					strs = append(strs, tree.GoString())
				}

				return strs
			})

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
}
