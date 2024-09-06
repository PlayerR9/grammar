package parser

import (
	"fmt"
	"slices"

	"github.com/PlayerR9/go-commons/cmp"
	"github.com/PlayerR9/go-commons/queue"
	"github.com/PlayerR9/go-commons/set"
	dbg "github.com/PlayerR9/go-debug/assert"
	"github.com/PlayerR9/grammar/internal"
)

// parse_table is the parsing table.
type parse_table[T internal.TokenTyper] struct {
	// symbols is the set of all symbols in the grammar.
	symbols *cmp.Set[T]

	// rule_set is the set of all rules in the grammar.
	rule_set *set.Set[*Rule[T]]

	// item_set is the set of all items in the grammar.
	item_set *set.Set[*Item[T]]

	// states is the set of all states in the grammar.
	states []*State[T]

	// action_table is the action table.
	action_table map[*State[T]]map[T]internal.ActionType

	// goto_table is the goto table.
	goto_table map[*State[T]]map[T]*State[T]
}

// make_symbols is a helper function that makes the symbols set.
func (pt *parse_table[T]) make_symbols() {
	dbg.AssertNotNil(pt, "pt")
	dbg.AssertNotNil(pt.rule_set, "pt.rule_set")
	dbg.Assert(pt.symbols.IsEmpty(), "symbols is not empty")

	for rule := range pt.rule_set.All() {
		_ = pt.symbols.Add(rule.Lhs())
		_ = pt.symbols.Union(rule.Symbols())
	}
}

// make_items is a helper function that makes the items set.
func (pt *parse_table[T]) make_items() {
	dbg.AssertNotNil(pt, "pt")
	dbg.AssertNotNil(pt.rule_set, "pt.rule_set")
	dbg.Assert(pt.item_set.IsEmpty(), "item_set is not empty")

	for rule := range pt.rule_set.All() {
		for i := 0; i <= rule.Size(); i++ {
			item, err := NewItem(rule, i)
			dbg.AssertErr(err, "NewItem(rule, %d)", i)

			pt.item_set.Add(item)
		}
	}
}

// new_parse_table creates a new parse table.
//
// Parameters:
//   - rules: The rules of the grammar.
//
// Returns:
//   - *parse_table[T]: The new parse table. Never returns nil.
func new_parse_table[T internal.TokenTyper](rules []*Rule[T]) *parse_table[T] {
	pt := &parse_table[T]{
		symbols:  cmp.NewSet[T](),
		rule_set: set.NewSetWithItems(rules),
		item_set: set.NewSet[*Item[T]](),
	}

	pt.make_symbols()
	pt.make_items()

	return pt
}

// get_items_with_lhs returns all items with the given lhs.
//
// Parameters:
//   - lhs: The left-hand side of the items.
//
// Returns:
//   - []*Item[T]: The items with the given lhs.
func (pt parse_table[T]) get_items_with_lhs(lhs T) []*Item[T] {
	var items []*Item[T]

	for item := range pt.item_set.All() {
		if item.Lhs() == lhs {
			items = append(items, item)
		}
	}

	return items
}

// closure returns the closure of the given item set.
//
// Parameters:
//   - seed: The seed items.
//
// Returns:
//   - []*Item[T]: The closure of the item set.
func (pt parse_table[T]) closure(seed []*Item[T]) []*Item[T] {
	if len(seed) == 0 {
		return nil
	}

	var result []*Item[T]

	q := queue.NewQueueWithElems(seed)

	for {
		first, ok := q.Dequeue()
		if !ok {
			break
		}

		if slices.Contains(result, first) {
			continue // already evaluated
		}

		for rhs := range first.Rhs() {
			if rhs.IsTerminal() {
				continue
			}

			tmp := pt.get_items_with_lhs(rhs)
			seed = append(seed, tmp...)
		}
	}

	return result
}

/* // Goto returns the closure of the item set after the given item is advanced.
//
// Parameters:
//   - item: The item.
//   - rhs: The right-hand side of the item.
//
// Returns:
//   - *set.Set[*Item[T]]: The closure of the item set.
//   - error: An error if goto failed.
func (pt ParseTable[T]) Goto(item *Item[T], rhs T) ([]*Item[T], error) {
	if item == nil {
		return nil, gcers.NewErrNilParameter("item")
	}

	after, ok := item.RhsAt(item.Pos() + 1)
	if !ok {
		return nil, fmt.Errorf("expected %q, got nothing instead", rhs.String())
	} else if after != rhs {
		return nil, fmt.Errorf("expected %q, got %q instead", rhs.String(), after.String())
	}

	ok = item.Advance()
	dbg.AssertOk(ok, "item.Advance()")

	return pt.Closure([]*Item[T]{item}), nil
} */

// make_all_states is a helper function that makes all states.
//
// Returns:
//   - error: An error if the closure failed.
func (pt *parse_table[T]) make_all_states() error {
	start_symbol := T(0)

	initial_items := pt.get_items_with_lhs(start_symbol)
	if len(initial_items) == 0 {
		return fmt.Errorf("there are no rules for the start symbol (%q)", start_symbol.String())
	} else if len(initial_items) > 1 {
		return fmt.Errorf("there are multiple rules for the start symbol (%q)", start_symbol.String())
	}

	state0 := NewState(initial_items[0], pt.closure(initial_items))

	pt.states = []*State[T]{state0}
	state_queue := queue.NewQueueWithElems([]*State[T]{state0})

	for {
		first, ok := state_queue.Dequeue()
		if !ok {
			break
		}

		for _, rule := range first.Rule() {
			pos := rule.Pos() + 1

			next, ok := rule.RhsAt(pos)
			if !ok || next.IsTerminal() {
				continue
			}

			rule, ok = rule.Advance()
			dbg.AssertOk(ok, "rule.Advance()")

			idx := -1

			for i := 0; i < len(pt.states) && idx == -1; i++ {
				if pt.states[i].IsOfSeed(rule) {
					idx = i
				}
			}

			if idx == -1 {
				new_state := NewState(rule, pt.closure([]*Item[T]{rule}))

				state_queue.Enqueue(new_state)
				pt.states = append(pt.states, new_state)

				idx = len(pt.states) - 1
			}

			first.AddNext(pt.states[idx])
		}
	}

	return nil
}

// init is a helper function that initializes the parsing table.
//
// Returns:
//   - error: An error if the initialization failed.
func (pt *parse_table[T]) init() error {
	err := pt.make_all_states()
	if err != nil {
		return err
	}

	pt.action_table = make(map[*State[T]]map[T]internal.ActionType)
	pt.goto_table = make(map[*State[T]]map[T]*State[T])

	for _, state := range pt.states {
		actions := make(map[T]internal.ActionType)
		gotos := make(map[T]*State[T])

		for symbol := range pt.symbols.All() {
			if symbol.IsTerminal() {
				seed := state.Seed()

				rhs, ok := seed.RhsAt(seed.Pos())
				if !ok {
					if symbol == T(0) {
						actions[symbol] = internal.ActAcceptType
					} else {
						actions[symbol] = internal.ActReduceType
					}
				} else if rhs != symbol {
					// Do a better handling here
					continue
				} else {
					actions[symbol] = internal.ActShiftType
				}

				gotos[symbol] = nil
			} else {
				var ns []*State[T]

				for next_state := range state.NextState() {
					seed := next_state.Seed()

					rhs, ok := seed.RhsAt(seed.Pos())
					if !ok || rhs != symbol {
						continue
					}

					ns = append(ns, next_state)
				}

				if len(ns) == 0 {
					gotos[symbol] = nil
				} else if len(ns) > 1 {
					return fmt.Errorf("ambiguous goto from %q", symbol.String())
				}

				actions[symbol] = internal.ActShiftType // FIXME: Make a new action type.
				gotos[symbol] = ns[0]
			}
		}

		pt.action_table[state] = actions
		pt.goto_table[state] = gotos
	}

	return nil
}
