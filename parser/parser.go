package parser

import (
	"errors"
	"fmt"

	gr "github.com/PlayerR9/grammar/grammar"
	luc "github.com/PlayerR9/lib_units/common"
)

// Parser is an interface that defines the behavior of a parser.
type Parser[T gr.TokenTyper] interface {
	// SetInputStream sets the input stream of the parser.
	//
	// Parameters:
	//   - tokens: The input stream of the parser.
	SetInputStream(tokens []*gr.Token[T])

	// GetDecision is a function that gets the decision of the parser.
	//
	// Parameters:
	//   - lookahead: The lookahead token.
	//
	// Returns:
	//   - Actioner: The decision of the parser.
	//   - error: An error if the parser encounters an error while getting the decision.
	GetDecision(lookahead *gr.Token[T]) (Actioner, error)

	// Shift is a function that shifts the input stream of the parser.
	//
	// Returns:
	//   - bool: True if the parser could shift the input stream, false otherwise.
	Shift() bool

	// Pop pops the top token of the stack.
	//
	// Returns:
	//   - *Token[T]: The top token of the stack.
	//   - bool: True if the stack is not empty, false otherwise.
	Pop() (*gr.Token[T], bool)

	// Peek peeks the top token of the stack.
	//
	// Returns:
	//   - *Token[T]: The top token of the stack.
	//   - bool: True if the stack is not empty, false otherwise.
	Peek() (*gr.Token[T], bool)

	// GetPopped returns the popped tokens of the parser.
	//
	// Returns:
	//   - []*Token[T]: The popped tokens of the parser.
	//
	// The last token returned is the furthest token in the rule.
	GetPopped() []*gr.Token[T]

	// Push pushes a token onto the stack. Does nothing if the token is nil.
	//
	// Parameters:
	//   - token: The token to push onto the stack.
	Push(token *gr.Token[T])

	// Refuse is a function that refuses any token that was popped from the stack.
	Refuse()

	// Accept is a function that accepts any token that was popped from the stack.
	Accept()
}

// apply_reduce applies a reduce action to the parser.
//
// Parameters:
//   - parser: The parser.
//   - rule: The rule to reduce.
//
// Returns:
//   - error: An error if the parser encounters an error while applying the reduce action.
func apply_reduce[T gr.TokenTyper](parser Parser[T], rule *Rule[T]) error {
	luc.AssertParam("parser", parser != nil, errors.New("value cannot be nil"))
	luc.AssertNil(rule, "rule")

	iter := rule.Iterator()
	luc.Assert(iter != nil, "iter must not be nil")

	var prev *T

	for {
		value, err := iter.Consume()
		if err != nil {
			break
		}

		top, ok := parser.Pop()
		if !ok {
			return NewErrUnexpectedToken(prev, nil, value)
		} else if top.Type != value {
			return NewErrUnexpectedToken(prev, &top.Type, value)
		}
	}

	popped := parser.GetPopped()
	last_token := popped[len(popped)-1]

	parser.Accept()

	tk, err := gr.NewToken(rule.lhs, popped, last_token.At, last_token.Lookahead)
	luc.AssertErr(err, "NewToken(%s, popped, %d, %d)", rule.lhs.String(), last_token.At, last_token.Lookahead)

	parser.Push(tk)

	return nil
}

// get_forest returns the syntax forest of the parser.
//
// Parameters:
//   - parser: The parser.
//
// Returns:
//   - []*Token[T]: The syntax forest of the parser.
func get_forest[T gr.TokenTyper](parser Parser[T]) []*gr.Token[T] {
	luc.Assert(parser != nil, "parser must not be nil")

	var forest []*gr.Token[T]

	for {
		top, ok := parser.Pop()
		if !ok {
			break
		}

		forest = append(forest, top)
	}

	return forest
}

// FullParse parses the input stream of the parser.
//
// Parameters:
//   - parser: The parser.
//   - tokens: The input stream of the parser.
//
// Returns:
//   - []*Token[T]: The syntax forest of the input stream.
//   - error: An error if the parser encounters an error while parsing the input stream.
func FullParse[T gr.TokenTyper](parser Parser[T], tokens []*gr.Token[T]) ([]*gr.Token[T], error) {
	if parser == nil {
		forest := get_forest(parser)

		return forest, luc.NewErrNilParameter("parser")
	}

	parser.SetInputStream(tokens)

	ok := parser.Shift() // initial shift
	if !ok {
		forest := get_forest(parser)

		return forest, fmt.Errorf("no tokens in input stream")
	}

	for {
		top, ok := parser.Peek()
		luc.AssertOk(ok, "parser.Peek()")

		act, err := parser.GetDecision(top.Lookahead)
		if err != nil {
			forest := get_forest(parser)

			return forest, fmt.Errorf("error getting decision: %w", err)
		}

		switch act := act.(type) {
		case *ShiftAction:
			ok := parser.Shift()
			luc.AssertOk(ok, "parser.Shift()")
		case *ReduceAction[T]:
			err := apply_reduce(parser, act.rule)
			if err != nil {
				parser.Refuse()

				forest := get_forest(parser)

				return forest, fmt.Errorf("error applying reduce: %w", err)
			}
		case *AcceptAction[T]:
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
