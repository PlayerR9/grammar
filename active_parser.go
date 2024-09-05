package grammar

import (
	"errors"
	"fmt"
	"slices"

	gcers "github.com/PlayerR9/go-commons/errors"
	"github.com/PlayerR9/go-commons/stack"
	"github.com/PlayerR9/go-commons/tree"
	dbg "github.com/PlayerR9/go-debug/assert"
	gr "github.com/PlayerR9/grammar/grammar"
	internal "github.com/PlayerR9/grammar/internal"
)

// ActiveParser is the active parser (i.e., the one that is currently parsing).
type ActiveParser[T internal.TokenTyper] struct {
	// global contains the shared information between active parsers.
	global *Parser[T]

	// reader is the token reader.
	reader TokenReader[T]

	// token_stack is the token token_stack.
	token_stack *stack.RefusableStack[*gr.Token[T]]

	// history is the history of the parser.
	history *History[T]

	// err is the reason to why the active parser has failed. Nil if it has succeded.
	err error
}

// new_active_parser creates a new active parser.
//
// Parameters:
//   - global: The shared information between active parsers.
//   - history: The history of the parser.
//
// Returns:
//   - *ActiveParser: The new active parser.
//   - error: An error if 'global' is nil or no tokens are provided.
func new_active_parser[T internal.TokenTyper](global *Parser[T], history *History[T]) (*ActiveParser[T], error) {
	if global == nil {
		return nil, gcers.NewErrNilParameter("global")
	} else if len(global.tokens) == 0 {
		return nil, errors.New("no tokens provided")
	}

	tokens := make([]*gr.Token[T], 0, len(global.tokens))
	for i := 0; i < len(global.tokens); i++ {
		tokens = append(tokens, global.tokens[i].Copy())
	}

	for i := 0; i < len(tokens)-1; i++ {
		tokens[i].Lookahead = tokens[i+1]
	}

	if history == nil {
		history = NewHistory[T](nil)
	}

	return &ActiveParser[T]{
		global:      global,
		reader:      NewTokenStream(tokens),
		token_stack: stack.NewRefusableStack[*gr.Token[T]](),
		history:     history,
		err:         nil,
	}, nil
}

// Pop pops a token from the stack.
//
// Returns:
//   - *Token[T]: The popped token.
//   - bool: True if the token was popped, false otherwise.
func (ap *ActiveParser[T]) Pop() (*gr.Token[T], bool) {
	return ap.token_stack.Pop()
}

// can_walk checks if the active parser can walk.
//
// Returns:
//   - bool: True if the active parser can walk, false otherwise.
func (ap *ActiveParser[T]) can_walk() bool {
	return ap.history.CanWalk()
}

// walk walks the active parser.
//
// Parameters:
//   - decision_err: The decision error.
//
// Returns:
//   - bool: True if the walk is an accept action, false otherwise.
func (ap *ActiveParser[T]) walk(decision_err error) bool {
	var ok bool

	fn := func(item *Item[T]) error {
		tmp, err := ap.apply(decision_err, item)
		if err != nil {
			return err
		}

		ok = tmp

		return nil
	}

	err := ap.history.Walk(fn)
	if err != nil {
		ap.err = NewErrParsing(err, nil)

		return false
	}

	return ok
}

// exec executes the active parser.
//
// Returns:
//   - []*ActiveParser[T]: The possible paths.
func (ap *ActiveParser[T]) exec() []*ActiveParser[T] {
	var possible_paths []*ActiveParser[T]

	for {
		items, decision_err := ap.global.rule_set.Decision(ap)
		ap.token_stack.Refuse()

		if len(items) == 0 {
			if decision_err == nil {
				decision_err = errors.New("no action available")
			}

			ap.err = NewErrParsing(decision_err, nil)

			return possible_paths
		}

		if len(items) == 1 {
			ap.history.AddEvent(items[0])
		} else {
			original_history := ap.history.Copy()

			ap.history.AddEvent(items[0])

			for _, item := range items[1:] {
				new_history := original_history.Copy()
				new_history.AddEvent(item)

				new_active, err := new_active_parser(ap.global, new_history)
				dbg.AssertErr(err, "NewActiveParser(ap.global, new_history)")

				possible_paths = append(possible_paths, new_active)
			}
		}

		dbg.AssertOk(ap.can_walk(), "p.CanWalk()")

		is_accept := ap.walk(decision_err)
		if ap.HasError() {
			return possible_paths
		}

		if is_accept {
			if ap.token_stack.Size() != 1 {
				ap.err = NewErrParsing(errors.New("not a valid parse"), nil)

				return possible_paths
			}

			break
		}
	}

	return possible_paths
}

// exec_witn_fn executes the active parser with a custom decision function.
//
// Returns:
//   - []*ActiveParser[T]: The possible paths.
func (ap *ActiveParser[T]) exec_witn_fn() []*ActiveParser[T] {
	var possible_paths []*ActiveParser[T]

	for {
		items, decision_err := ap.global.decision_fn(ap)
		ap.token_stack.Refuse()

		dbg.AssertThat("items", dbg.NewOrderedAssert(len(items)).GreaterThan(0))

		if len(items) == 0 {
			if decision_err == nil {
				decision_err = errors.New("no action available")
			}

			ap.err = NewErrParsing(decision_err, nil)

			return possible_paths
		}

		if len(items) == 1 {
			ap.history.AddEvent(items[0])
		} else {
			original_history := ap.history.Copy()

			ap.history.AddEvent(items[0])

			for _, item := range items[1:] {
				new_history := original_history.Copy()
				new_history.AddEvent(item)

				new_active, err := new_active_parser(ap.global, new_history)
				dbg.AssertErr(err, "NewActiveParser(ap.global, new_history)")

				possible_paths = append(possible_paths, new_active)
			}
		}

		dbg.AssertOk(ap.can_walk(), "p.CanWalk()")

		is_accept := ap.walk(decision_err)
		if ap.HasError() {
			return possible_paths
		}

		if is_accept {
			if ap.token_stack.Size() != 1 {
				ap.err = NewErrParsing(errors.New("not a valid parse"), nil)

				return possible_paths
			}

			break
		}
	}

	return possible_paths
}

// walk_all walks the active parser.
func (ap *ActiveParser[T]) walk_all() {
	err := ap.shift() // initial shift
	if err != nil {
		ap.err = NewErrParsing(err, nil)

		return
	}

	for ap.can_walk() {
		is_accept := ap.walk(nil)
		if ap.HasError() || is_accept {
			break
		}
	}
}

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

// apply is a helper function that applies the action to the stack.
//
// Returns:
//   - bool: True if the parse was accepted, false otherwise.
//   - error: An error if any.
func (ap *ActiveParser[T]) apply(decision_err error, item *Item[T]) (bool, error) {
	dbg.AssertNotNil(item, "item")

	act := item.act

	switch act {
	case internal.ActShiftType:
		err := ap.shift()

		if err != nil {
			return false, NewErrParsing(fmt.Errorf("error shifting: %w", err), decision_err)
		}
	case internal.ActReduceType:
		err := ap.reduce(item.rule)
		if err != nil {
			ap.token_stack.Refuse()

			return false, NewErrParsing(fmt.Errorf("error reducing: %w", err), decision_err)
		}
	case internal.ActAcceptType:
		err := ap.reduce(item.rule)
		if err != nil {
			ap.token_stack.Refuse()

			return false, NewErrParsing(fmt.Errorf("error accepting: %w", err), decision_err)
		}

		return true, nil
	default:
		return false, NewErrParsing(fmt.Errorf("invalid action: %v", act), decision_err)
	}

	return false, nil
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
	return ap.err
}

// HasError checks if the error is not nil.
//
// Returns:
//   - bool: True if the error is not nil.
func (ap ActiveParser[T]) HasError() bool {
	return ap.err != nil
}
