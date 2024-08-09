package parsing

import (
	"fmt"
	"slices"

	gcers "github.com/PlayerR9/go-commons/errors"
	gr "github.com/PlayerR9/grammar/grammar"
)

// DecisionFunc is the function that returns the decision of the parser.
//
// Parameters:
//   - parser: The parser.
//   - lookahead: The lookahead token.
//
// Returns:
//   - Actioner: The action of the decision.
//   - error: An error if the decision is invalid.
type DecisionFunc[S gr.TokenTyper] func(parser *Parser[S], lookahead *gr.Token[S]) (Actioner, error)

// Parser is the parser of the grammar.
type Parser[S gr.TokenTyper] struct {
	// tokens is the tokens of the parser.
	tokens []*gr.Token[S]

	// stack is the stack of the parser.
	stack []*gr.Token[S]

	// popped is the stack of the parser.
	popped []*gr.Token[S]

	// decision is the function that returns the decision of the parser.
	decision DecisionFunc[S]
}

// NewParser creates a new parser.
//
// Parameters:
//   - decision_func: The function that returns the decision of the parser.
//
// Returns:
//   - *Parser: The new parser.
//
// This function returns nil if the decision_func is nil.
func NewParser[S gr.TokenTyper](decision_func DecisionFunc[S]) *Parser[S] {
	if decision_func == nil {
		return nil
	}

	return &Parser[S]{
		decision: decision_func,
	}
}

// SetInputStream sets the input stream of the parser.
//
// Parameters:
//   - tokens: The input stream of the parser.
func (p *Parser[S]) SetInputStream(tokens []*gr.Token[S]) {
	// Clean the previous parser's state.
	gr.CleanTokens(p.tokens)
	p.tokens = p.tokens[:0]

	gr.CleanTokens(p.stack)
	p.stack = p.stack[:0]

	gr.CleanTokens(p.popped)
	p.popped = p.popped[:0]

	// Set the new input stream.

	p.tokens = tokens
}

// Pop pops a token from the stack.
//
// Returns:
//   - *Token[T]: The token if the stack is not empty, nil otherwise.
//   - bool: True if the stack is not empty, false otherwise.
func (p *Parser[S]) Pop() (*gr.Token[S], bool) {
	if len(p.stack) == 0 {
		return nil, false
	}

	top := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]

	p.popped = append(p.popped, top)

	return top, true
}

// Peek pops a token from the stack without removing it.
//
// Returns:
//   - *Token[T]: The token if the stack is not empty, nil otherwise.
//   - bool: True if the stack is not empty, false otherwise.
func (p *Parser[S]) Peek() (*gr.Token[S], bool) {
	if len(p.stack) == 0 {
		return nil, false
	}

	return p.stack[len(p.stack)-1], true
}

// Shift shifts a token from the input stream to the stack.
//
// Returns:
//   - bool: True if the input stream is not empty, false otherwise.
func (p *Parser[S]) Shift() bool {
	if len(p.tokens) == 0 {
		return false
	}

	first := p.tokens[0]
	p.tokens = p.tokens[1:]

	p.stack = append(p.stack, first)

	return true
}

// GetPopped returns the popped tokens.
//
// Returns:
//   - []*Token[S]: The popped tokens.
func (p *Parser[S]) GetPopped() []*gr.Token[S] {
	popped := make([]*gr.Token[S], 0, len(p.popped))

	for i := 0; i < len(p.popped); i++ {
		popped = append(popped, p.popped[i])
	}

	slices.Reverse(popped)

	return popped
}

// Push pushes a token to the stack. Does nothing if the token is nil.
//
// Parameters:
//   - token: The token to push.
func (p *Parser[S]) Push(token *gr.Token[S]) {
	if token == nil {
		return
	}

	p.stack = append(p.stack, token)
}

// Refuse refuses all the tokens that were popped since the last
// call to Accept().
func (p *Parser[S]) Refuse() {
	for len(p.popped) > 0 {
		top := p.popped[len(p.popped)-1]
		p.popped = p.popped[:len(p.popped)-1]

		p.stack = append(p.stack, top)
	}
}

// Accept accepts all the tokens that were popped since the last
// call to Accept().
func (p *Parser[S]) Accept() {
	p.popped = p.popped[:0]
}

// get_forest returns the syntax forest of the parser.
//
// Parameters:
//   - parser: The parser.
//
// Returns:
//   - []*grammar.TokenTree[S]: The syntax forest of the parser.
func get_forest[S gr.TokenTyper](parser *Parser[S]) []*gr.TokenTree[S] {
	if parser == nil {
		panic("parser cannot be nil")
	}

	// luc.Assert(parser != nil, "parser must not be nil")

	var forest []*gr.TokenTree[S]

	for {
		top, ok := parser.Pop()
		if !ok {
			break
		}

		tree, _ := gr.NewTokenTree(top)
		// luc.AssertErr(err, "gr.NewTokenTree(top)")

		forest = append(forest, tree)
	}

	return forest
}

// apply_reduce applies a reduce action to the parser.
//
// Parameters:
//   - parser: The parser.
//   - rule: The rule to reduce.
//
// Returns:
//   - error: An error if the parser encounters an error while applying the reduce action.
func apply_reduce[S gr.TokenTyper](parser *Parser[S], rule *Rule[S]) error {
	if parser == nil {
		panic("parser cannot be nil")
	} else if rule == nil {
		panic("rule cannot be nil")
	}

	var prev *S

	for _, rhs := range rule.GetRhss() {
		top, ok := parser.Pop()
		if !ok {
			return NewErrUnexpectedToken(prev, nil, rhs)
		}

		top_type := top.GetType()

		if top_type != rhs {
			return NewErrUnexpectedToken(prev, &top_type, rhs)
		}
	}

	popped := parser.GetPopped()
	last_token := popped[len(popped)-1]

	parser.Accept()

	tk := gr.NewToken(rule.lhs, "", last_token.At, last_token.Lookahead)
	tk.AddChildren(popped)

	parser.Push(tk)

	return nil
}

// FullParse is just a wrapper around the Grammar.FullParse function.
//
// Parameters:
//   - tokens: The input stream of the parser.
//
// Returns:
//   - []*gr.TokenTree[S]: The syntax forest of the input stream.
//   - error: An error if the parser encounters an error while parsing the input stream.
func FullParse[S gr.TokenTyper](parser *Parser[S], tokens []*gr.Token[S]) ([]*gr.TokenTree[S], error) {
	if parser == nil {
		forest := get_forest(parser)

		return forest, gcers.NewErrNilParameter("parser")
	}

	parser.SetInputStream(tokens)

	ok := parser.Shift() // initial shift
	if !ok {
		forest := get_forest(parser)

		return forest, fmt.Errorf("no tokens in input stream")
	}

	for {
		top, _ := parser.Peek()
		// luc.AssertOk(ok, "parser.Peek()")

		act, err := parser.decision(parser, top.Lookahead)
		parser.Refuse()

		if err != nil {
			forest := get_forest(parser)

			return forest, fmt.Errorf("error getting decision: %w", err)
		}

		switch act := act.(type) {
		case *ShiftAction:
			_ = parser.Shift()
			// luc.AssertOk(ok, "parser.Shift()")
		case *ReduceAction[S]:
			err := apply_reduce(parser, act.rule)
			if err != nil {
				parser.Refuse()

				forest := get_forest(parser)

				return forest, fmt.Errorf("error applying reduce: %w", err)
			}
		case *AcceptAction[S]:
			err := apply_reduce(parser, act.rule)
			if err != nil {
				parser.Refuse()

				forest := get_forest(parser)

				return forest, fmt.Errorf("error applying accept: %w", err)
			}

			forest := get_forest(parser)

			return forest, nil
		default:
			forest := get_forest(parser)

			return forest, fmt.Errorf("unexpected action: %v", act)
		}
	}
}
