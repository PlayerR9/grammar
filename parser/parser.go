package parser

import (
	"fmt"
	"slices"

	gr "github.com/PlayerR9/grammar/grammar"
)

// ParseFunc is a function that parses a token.
//
// Parameters:
//   - parser: The parser. Assumed to be non-nil.
//   - top1: The first token. Assumed to be non-nil.
//   - la: The lookahead token.
//
// Returns:
//   - Actioner: The action to perform.
//   - error: An error if the decision is invalid.
type ParseFunc[T gr.Enumer] func(parser *Parser[T], top1 *gr.Token[T], la *gr.Token[T]) (Actioner, error)

// Parser is a parser.
type Parser[T gr.Enumer] struct {
	// table is the table of rules.
	table map[T]ParseFunc[T]

	// tokens is the list of tokens to parse.
	tokens []*gr.Token[T]

	// stack is the stack of tokens.
	stack []*gr.Token[T]

	// popped is the list of tokens that have been popped.
	popped []*gr.Token[T]
}

// Pop pops a token from the stack.
//
// Returns:
//   - *gr.Token[T]: The popped token.
//   - bool: True if the token was popped, false otherwise.
func (p *Parser[T]) Pop() (*gr.Token[T], bool) {
	if len(p.tokens) == 0 {
		return nil, false
	}

	tk := p.tokens[0]
	p.tokens = p.tokens[1:]

	p.popped = append(p.popped, tk)

	return tk, true
}

// decision is a helper function that decides what to do next.
//
// Returns:
//   - Actioner: The action to perform.
//   - error: An error if the decision is invalid.
func (p *Parser[T]) decision() (Actioner, error) {
	top1, ok := p.Pop()
	if !ok {
		return nil, fmt.Errorf("unexpected EOF")
	}

	fn, ok := p.table[top1.Type]
	if !ok {
		return nil, fmt.Errorf("unexpected token: %v", top1)
	}

	act, err := fn(p, top1, top1.Lookahead)
	if err != nil {
		return nil, err
	}

	return act, nil
}

// shift is a helper function that shifts a token.
//
// Returns:
//   - bool: True if the token was shifted, false otherwise.
func (p *Parser[T]) shift() bool {
	if len(p.tokens) == 0 {
		return false
	}

	if len(p.popped) > 0 {
		panic("popped should be empty when shifting")
	}

	top := p.tokens[0]
	p.tokens = p.tokens[1:]

	p.stack = append(p.stack, top)

	return true
}

// refuse is a helper function that refuses all tokens that were popped.
func (p *Parser[T]) refuse() {
	for len(p.popped) > 0 {
		top := p.popped[0]
		p.popped = p.popped[1:]

		p.stack = append(p.stack, top)
	}
}

// accept is a helper function that accepts all tokens that were popped.
func (p *Parser[T]) accept() {
	p.popped = p.popped[:0]
}

// get_popped returns the list of tokens that have been popped.
//
// Returns:
//   - []*gr.Token[T]: The list of tokens that have been popped.
func (p Parser[T]) get_popped() []*gr.Token[T] {
	popped := make([]*gr.Token[T], len(p.popped))
	copy(popped, p.popped)

	slices.Reverse(popped)

	return popped
}

// reduce is a helper function that reduces a rule.
//
// Parameters:
//   - rule: The rule to reduce.
//
// Returns:
//   - error: An error if the rule could not be reduced.
func (p *Parser[T]) reduce(rule *Rule[T]) error {
	if rule == nil {
		panic("rule should not be nil")
	}

	for rhs := range rule.BackwardRhs() {
		top, ok := p.Pop()
		if !ok {
			return NewErrUnexpectedToken(rhs, rhs, nil)
		} else if top.Type != rhs {
			return NewErrUnexpectedToken(rhs, rhs, &top.Type)
		}
	}

	popped := p.get_popped()
	if len(popped) == 0 {
		panic("popped should not be empty")
	}

	tk, err := gr.NewToken(rule.Lhs(), "", popped)
	if err != nil {
		panic(fmt.Sprintf("could not create token: %v", err))
	}

	p.stack = append(p.stack, tk)

	return nil
}

// Parse parses a list of tokens.
//
// Parameters:
//   - tokens: The list of tokens to parse.
//
// Returns:
//   - *gr.Token[T]: The root token of the parse tree.
//   - error: An error if the parse failed.
func (p *Parser[T]) Parse(tokens []*gr.Token[T]) (*gr.Token[T], error) {
	if !p.shift() {
		return nil, fmt.Errorf("nothing to parse")
	}

	for {
		act, err := p.decision()
		p.refuse()

		if err != nil {
			return nil, err
		} else if act == nil {
			return nil, fmt.Errorf("no decision was made")
		}

		switch act := act.(type) {
		case *ShiftAct:
			if !p.shift() {
				return nil, fmt.Errorf("could not shift")
			}
		case *ReduceAct[T]:
			err := p.reduce(act.Rule())
			if err != nil {
				return nil, err
			}

			p.accept()
		case *AcceptAct[T]:
			err := p.reduce(act.Rule())
			if err != nil {
				return nil, err
			}

			p.accept()

			forest := make([]*gr.Token[T], len(p.stack))
			copy(forest, p.stack)

			slices.Reverse(forest)

			if len(forest) != 1 {
				return nil, fmt.Errorf("expected exactly one root but got %d", len(forest))
			}

			root := forest[0]

			return root, nil
		default:
			return nil, fmt.Errorf("unexpected action: %T", act)
		}
	}
}
