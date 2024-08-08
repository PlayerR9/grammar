package parser

import (
	"fmt"
	"slices"

	gcers "github.com/PlayerR9/go-commons/errors"
	gr "github.com/PlayerR9/grammar/grammar"
	llq "github.com/PlayerR9/listlike/queue"
	lls "github.com/PlayerR9/listlike/stack"
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
	tokens *llq.ArrayQueue[*gr.Token[S]]

	// stack is the stack of the parser.
	stack *lls.ArrayStack[*gr.Token[S]]

	// popped is the stack of the parser.
	popped *lls.ArrayStack[*gr.Token[S]]

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
		stack:    lls.NewArrayStack[*gr.Token[S]](),
		popped:   lls.NewArrayStack[*gr.Token[S]](),
		decision: decision_func,
	}
}

// SetInputStream sets the input stream of the parser.
//
// Parameters:
//   - tokens: The input stream of the parser.
func (p *Parser[S]) SetInputStream(tokens []*gr.Token[S]) {
	p.tokens = llq.NewArrayQueue[*gr.Token[S]]()
	p.tokens.EnqueueMany(tokens)

	p.stack.Clear()
	p.popped.Clear()
}

// Pop pops a token from the stack.
//
// Returns:
//   - *Token[T]: The token if the stack is not empty, nil otherwise.
//   - bool: True if the stack is not empty, false otherwise.
func (p *Parser[S]) Pop() (*gr.Token[S], bool) {
	top, ok := p.stack.Pop()
	if !ok {
		return nil, false
	}

	p.popped.Push(top)

	return top, true
}

// Peek pops a token from the stack without removing it.
//
// Returns:
//   - *Token[T]: The token if the stack is not empty, nil otherwise.
//   - bool: True if the stack is not empty, false otherwise.
func (p *Parser[S]) Peek() (*gr.Token[S], bool) {
	top, ok := p.stack.Peek()
	if !ok {
		return nil, false
	}

	return top, true
}

// Shift shifts a token from the input stream to the stack.
//
// Returns:
//   - bool: True if the input stream is not empty, false otherwise.
func (p *Parser[S]) Shift() bool {
	first, ok := p.tokens.Dequeue()
	if !ok {
		return false
	}

	p.stack.Push(first)

	return true
}

// GetPopped returns the popped tokens.
//
// Returns:
//   - []*Token[S]: The popped tokens.
func (p *Parser[S]) GetPopped() []*gr.Token[S] {
	popped := p.popped.Slice()
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

	p.stack.Push(token)
}

// Refuse refuses all the tokens that were popped since the last
// call to Accept().
func (p *Parser[S]) Refuse() {
	for {
		top, ok := p.popped.Pop()
		if !ok {
			break
		}

		p.stack.Push(top)
	}
}

// Accept accepts all the tokens that were popped since the last
// call to Accept().
func (p *Parser[S]) Accept() {
	p.popped.Clear()
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
