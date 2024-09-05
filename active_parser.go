package grammar

import (
	"errors"
	"fmt"
	"slices"

	gcers "github.com/PlayerR9/go-commons/errors"
	"github.com/PlayerR9/go-commons/stack"
	"github.com/PlayerR9/go-commons/tree"
	dbg "github.com/PlayerR9/go-debug/assert"
	internal "github.com/PlayerR9/grammar/internal"
)

// ActiveParser is the active parser (i.e., the one that is currently parsing).
type ActiveParser[T internal.TokenTyper] struct {
	// global contains the shared information between active parsers.
	global *Parser[T]

	// reader is the token reader.
	reader TokenReader[T]

	// token_stack is the token token_stack.
	token_stack *stack.RefusableStack[*Token[T]]

	// history is the history of the parser.
	history *History[T]

	// err is the reason to why the active parser has failed. Nil if it has succeded.
	err error
}

// NewActiveParser creates a new active parser.
//
// Parameters:
//   - global: The shared information between active parsers.
//   - history: The history of the parser.
//
// Returns:
//   - *ActiveParser: The new active parser.
//
// Panic with an error of type *pkg.Err if 'global' is nil.
func NewActiveParser[T internal.TokenTyper](global *Parser[T], history *History[T]) (*ActiveParser[T], error) {
	if global == nil {
		return nil, gcers.NewErrNilParameter("global")
	} else if len(global.tokens) == 0 {
		return nil, errors.New("no tokens provided")
	}

	tokens := make([]*Token[T], 0, len(global.tokens))
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
		token_stack: stack.NewRefusableStack[*Token[T]](),
		history:     history,
		err:         nil,
	}, nil
}

// Pop pops a token from the stack.
//
// Returns:
//   - *Token[T]: The popped token.
//   - bool: True if the token was popped, false otherwise.
func (ap *ActiveParser[T]) Pop() (*Token[T], bool) {
	return ap.token_stack.Pop()
}

func (ap ActiveParser[T]) IsDone() bool {
	return ap.token_stack.IsEmpty()
}

func (ap *ActiveParser[T]) CanWalk() bool {
	return ap.history.CanWalk()
}

func (ap *ActiveParser[T]) Walk(decision_err error) bool {
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

// Exec executes the active parser.
//
// Returns:
//   - *sdslices.Slice[*sdt.Wrap[*ActiveParser[T]]]: The possible paths. Never returns nil.
//
// Panic if the active parser fails.
func (ap *ActiveParser[T]) Exec() []*ActiveParser[T] {
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

				new_active, err := NewActiveParser(ap.global, new_history)
				dbg.AssertErr(err, "NewActiveParser(ap.global, new_history)")

				possible_paths = append(possible_paths, new_active)
			}
		}

		dbg.AssertOk(ap.CanWalk(), "p.CanWalk()")

		is_accept := ap.Walk(decision_err)
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

// ExecWithFn executes the active parser with a custom decision function.
//
// Returns:
//   - *sdslices.Slice[*sdt.Wrap[*ActiveParser[T]]]: The possible paths. Never returns nil.
//
// Panic if the active parser fails.
func (ap *ActiveParser[T]) ExecWithFn() []*ActiveParser[T] {
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

				new_active, err := NewActiveParser(ap.global, new_history)
				dbg.AssertErr(err, "NewActiveParser(ap.global, new_history)")

				possible_paths = append(possible_paths, new_active)
			}
		}

		dbg.AssertOk(ap.CanWalk(), "p.CanWalk()")

		is_accept := ap.Walk(decision_err)
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

func (ap *ActiveParser[T]) WalkAll() bool {
	err := ap.shift() // initial shift
	if err != nil {
		ap.err = NewErrParsing(err, nil)

		return false
	}

	for ap.CanWalk() {
		is_accept := ap.Walk(nil)
		if ap.HasError() {
			return false
		} else if is_accept {
			return true
		}
	}

	return true
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
	if rule == nil {
		return fmt.Errorf("no rule provided")
	}

	var prev *T

	for rhs := range rule.Backwards() {
		top, ok := ap.token_stack.Pop()
		if !ok {
			return NewErrUnexpectedToken(prev, nil, rhs)
		} else if top.Type != rhs {
			return NewErrUnexpectedToken(prev, &top.Type, rhs)
		}

		prev = &top.Type
	}

	popped := ap.token_stack.Popped()

	ap.token_stack.Accept()

	tk := NewToken(rule.Lhs(), "", popped[len(popped)-1].Lookahead)
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
//   - []*uttr.Tree[*Token[T]]: The forest.
func (ap ActiveParser[T]) Forest() []*tree.Tree[*Token[T]] {
	var forest []*tree.Tree[*Token[T]]

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
