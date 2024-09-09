package parser

import (
	"errors"
	"fmt"
	"slices"

	"github.com/PlayerR9/go-commons/stack"
	"github.com/PlayerR9/go-commons/tree"
	dbg "github.com/PlayerR9/go-debug/assert"
	gr "github.com/PlayerR9/grammar/PREV/grammar"
	internal "github.com/PlayerR9/grammar/PREV/internal"
)

// ActiveParser is the active parser (i.e., the one that is currently parsing).
type ActiveParser[T internal.TokenTyper] struct {
	// global contains the shared information between active parsers.
	global *Parser[T]

	// reader is the token reader.
	reader gr.TokenReader[T]

	// token_stack is the token token_stack.
	token_stack *stack.RefusableStack[*gr.Token[T]]

	// err is the reason to why the active parser has failed. Nil if it has succeded.
	err error

	// possible_cause is the possible cause of the error.
	possible_cause error

	// accept_found is true if an accept was found. False otherwise.
	accept_found bool
}

// HasError checks if the error is not nil.
//
// Returns:
//   - bool: True if the error is not nil.
func (ap ActiveParser[T]) HasError() bool {
	return ap.err != nil
}

// apply is a helper function that applies the action to the stack.
//
// Parameters:
//   - item: The item to apply.
//
// Returns:
//   - bool: True if the action is accepted. False otherwise.
func (ap *ActiveParser[T]) WalkOne(item *Item[T]) bool {
	dbg.AssertNotNil(item, "item")

	ap.accept_found = false

	act := item.act

	switch act {
	case internal.ActShiftType:
		err := ap.shift()

		if err != nil {
			ap.err = fmt.Errorf("error shifting: %w", err)
		}
	case internal.ActReduceType:
		err := ap.reduce(item.rule)
		if err != nil {
			ap.token_stack.Refuse()

			ap.err = fmt.Errorf("error reducing: %w", err)
		}
	case internal.ActAcceptType:
		err := ap.reduce(item.rule)
		if err == nil {
			if ap.token_stack.Size() == 1 {
				return true
			}

			ap.err = errors.New("not a valid parse")
			ap.possible_cause = nil
		} else {
			ap.token_stack.Refuse()

			ap.err = fmt.Errorf("error reducing: %w", err)
		}
	default:
		ap.err = fmt.Errorf("invalid action: %v", act)
	}

	return false
}

// exec executes the active parser.
//
// Parameters:
//   - history: The history of the parser.
//
// Returns:
//   - []*Item[T]: The possible paths.
func (ap *ActiveParser[T]) NextEvents() []*Item[T] {
	items, decision_err := ap.global.rule_set.Decision(ap)
	ap.token_stack.Refuse()

	if len(items) == 0 {
		if decision_err == nil {
			decision_err = errors.New("no action available")
		}

		ap.err = decision_err
		ap.possible_cause = nil

		return nil
	}

	if decision_err != nil {
		ap.possible_cause = decision_err
	}

	return items
}

// Pop pops a token from the stack.
//
// Returns:
//   - *Token[T]: The popped token.
//   - bool: True if the token was popped, false otherwise.
func (ap *ActiveParser[T]) Pop() (*gr.Token[T], bool) {
	return ap.token_stack.Pop()
}

/* // exec_witn_fn executes the active parser with a custom decision function.
//
// Parameters:
//   - history: The history of the parser.
//
// Returns:
//   - []*util.History[*Item[T]]: The possible paths.
func (ap *ActiveParser[T]) exec_witn_fn(history *util.History[*Item[T]]) []*util.History[*Item[T]] {
	dbg.AssertNotNil(history, "history")

	var possible_paths []*util.History[*Item[T]]

	for {
		items, decision_err := ap.global.decision_fn(ap)
		ap.token_stack.Refuse()

		dbg.AssertThat("items", dbg.NewOrderedAssert(len(items)).GreaterThan(0))

		if len(items) == 0 {
			if decision_err == nil {
				decision_err = errors.New("no action available")
			}

			ap.err = decision_err
			ap.possible_cause = nil

			return possible_paths
		}

		if decision_err != nil {
			ap.possible_cause = decision_err
		}

		if len(items) == 1 {
			history.AddEvent(items[0])
		} else {
			original_history := history.Copy()

			history.AddEvent(items[0])

			for _, item := range items[1:] {
				new_history := original_history.Copy()
				new_history.AddEvent(item)

				possible_paths = append(possible_paths, new_history)
			}
		}

		dbg.AssertOk(history.CanWalk(), "p.history.CanWalk()")

		err := history.WalkOnce()
		if err != nil {
			ap.err = err
			ap.possible_cause = nil

			return possible_paths
		}

		if ap.accept_found {
			if ap.token_stack.Size() != 1 {
				ap.err = errors.New("not a valid parse")
				ap.possible_cause = nil

				return possible_paths
			}

			break
		}
	}

	return possible_paths
} */

// reduce is a helper function that reduces the stack.
//
// Parameters:
//   - lhs: The left hand side token.
//   - rhss: The right hand side tokens.
//
// Returns:
//   - error: An error of type *ErrUnexpectedToken if any.
func (ap *ActiveParser[T]) reduce(rule *Rule[T]) error {
	dbg.AssertNotNil(rule, "rule")

	var prev *T

	for rhs := range rule.Backwards() {
		top, ok := ap.token_stack.Pop()
		if !ok {
			return gr.NewErrUnexpectedToken(prev, nil, rhs)
		} else if top.Type != rhs {
			return gr.NewErrUnexpectedToken(prev, &top.Type, rhs)
		}

		prev = &top.Type
	}

	popped := ap.token_stack.Popped()

	ap.token_stack.Accept()

	tk := gr.NewToken(rule.Lhs(), "", popped[len(popped)-1].Lookahead)
	tk.AddChildren(popped)

	ap.token_stack.Push(tk)

	return nil
}

// shift is a helper function that shifts the token.
//
// Returns:
//   - error: An error if any.
func (ap *ActiveParser[T]) shift() error {
	tk, err := ap.reader.ReadToken()
	if err != nil {
		return err
	}

	ap.token_stack.Push(tk)

	return nil
}

// Forest returns the tree that were parsed.
//
// Returns:
//   - []*uttr.Tree[*grammar.Token[T]]: The forest.
func (ap ActiveParser[T]) Forest() []*tree.Tree[*gr.Token[T]] {
	var forest []*tree.Tree[*gr.Token[T]]

	for {
		top, ok := ap.token_stack.Pop()
		if !ok {
			break
		}

		forest = append(forest, tree.NewTree(top))
	}

	slices.Reverse(forest)

	return forest
}

// Error returns the error if any.
//
// Returns:
//   - error: An error if any.
func (ap ActiveParser[T]) Error() error {
	if ap.err == nil {
		return nil
	}

	return NewErrParsing(ap.err, ap.possible_cause)
}
