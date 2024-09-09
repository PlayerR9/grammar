package parser

import (
	"fmt"

	utst "github.com/PlayerR9/go-commons/cmp"
	gcslc "github.com/PlayerR9/go-commons/slices"
	gr "github.com/PlayerR9/grammar/PREV/grammar"
	"github.com/PlayerR9/grammar/PREV/internal"
)

// decider is the decider of the active parser.
type decider[T internal.TokenTyper] struct {
	// p is the active parser.
	p *ActiveParser[T]

	// item_list is the list of items.
	item_list []*Item[T]

	// err is the reason to why the active parser has failed. Nil if it has succeded.
	err error
}

// new_decider is a helper function that creates a new decider.
//
// Parameters:
//   - p: The active parser.
//   - item_list: The list of items.
//
// Returns:
//   - *Decider[T]: The new decider. Never returns nil.
func new_decider[T internal.TokenTyper](p *ActiveParser[T], item_list []*Item[T]) *decider[T] {
	// dbg.AssertNotNil(p, "p")

	return &decider[T]{
		p:         p,
		item_list: item_list,
		err:       nil,
	}
}

// filter_lookaheads is a helper function that filters the lookahead sets.
//
// Parameters:
//   - indices: The indices.
//   - top1: The top1 token.
//
// Returns:
//   - []int: The filtered indices.
//   - []int: The solutions.
func (d *decider[T]) filter_lookaheads(indices []int, top1 *gr.Token[T]) ([]int, []int) {
	var solutions []int

	la := top1
	ok := true

	offset := 0

	for la != nil && ok {
		la = la.Lookahead

		var partial []int

		fn := func(idx int) bool {
			item := d.item_list[idx]

			ls, ok := item.LookaheadAt(offset)
			if !ok {
				partial = append(partial, idx)
			}

			return ok && la != nil && ls.Contains(la.Type)
		}

		indices, ok = gcslc.SFSeparateEarly(indices, fn)
		if len(partial) > 0 {
			solutions = append(solutions, partial...)
		}

		offset++
	}

	/* 	// Sort the solutions by their length. (longest first)
	   	fn := func(a, b []int) int {
	   		return len(b) - len(a)
	   	}

	   	slices.SortFunc(solutions, fn) */

	return indices, solutions
}

// apply_pop_rule is a helper function that applies the pop rule.
//
// Parameters:
//   - indices: The indices.
//   - prev: The previous token.
//   - offset: The offset.
//
// Returns:
//   - []int: The new indices.
//   - T: The new previous token.
func (d *decider[T]) apply_pop_rule(indices []int, prev T, offset int) ([]int, T) {
	// dbg.AssertThat("offset", dbg.NewOrderedAssert(offset).GreaterThan(0)).Panic()

	top, pop_ok := d.p.Pop()

	expected := utst.NewSet[T]()
	all_done := true

	fn := func(idx int) bool {
		item := d.item_list[idx]

		pos := item.pos - offset
		if pos > 0 {
			all_done = false
		}

		rhs, ok := item.RhsAt(pos)
		if ok {
			expected.Add(rhs)
		}

		return !ok || (pop_ok && rhs == top.Type)
	}

	tmp, ok := gcslc.SFSeparateEarly(indices, fn)
	if !ok {
		d.err = gr.NewErrUnexpectedToken(&prev, nil, expected.Slice()...)
	} else {
		if top != nil {
			prev = top.Type
		}

		indices = tmp
	}

	if all_done {
		// If all are dones, prioritize the items with the longest lookbehinds.

		indices = gcslc.MaxsFunc(indices, func(idx int) int {
			item := d.item_list[idx]

			return item.prevs.Len()
		})
	}

	return indices, prev
}

// decision is a helper function that decides which rule is the next one.
//
// Parameters:
//   - indices: The indices.
//   - prev: The previous token.
//
// Returns:
//   - []int: The new indices.
//   - error: The error. Never returns nil.
func (d *decider[T]) decision(indices []int, prev T) ([]int, error) {
	if d.err != nil {
		return nil, d.err
	}

	if len(indices) == 0 {
		return nil, fmt.Errorf("no rules available for %s", prev.String())
	}

	if len(indices) == 1 {
		// If there's only one rule, then we already know which one it is.

		return indices, nil
	}

	// If all rules are shift rules, then we don't care about which rule is chosen.

	all_shifts_rule := true

	for _, idx := range indices {
		item := d.item_list[idx]

		if !item.IsShift() {
			all_shifts_rule = false
			break
		}
	}

	if all_shifts_rule {
		return []int{indices[0]}, nil
	}

	return indices, nil
}

// filter_no_prev_items is a helper function that filters out the items that
// don't have the previous token.
//
// Parameters:
//   - indices: The indices.
//   - offset: The offset.
//
// Returns:
//   - []int: The filtered indices.
func (d decider[T]) filter_no_prev_items(indices []int, offset int) []int {
	if offset < 1 {
		return nil
	}

	fn := func(idx int) bool {
		item := d.item_list[idx]

		_, ok := item.RhsAt(item.pos - offset)
		return !ok
	}

	return gcslc.SliceFilter(indices, fn)
}

// only_lookaheads is a helper function that checks whether the indices only
// contain lookaheads.
//
// Parameters:
//   - indices: The indices.
//   - offset: The offset.
//
// Returns:
//   - bool: Whether the indices only contain lookaheads.
func (d decider[T]) only_lookaheads(indices []int, offset int) bool {
	if offset < 1 {
		return false
	}

	indices_copy := make([]int, len(indices))
	copy(indices_copy, indices)

	indices_copy = d.filter_no_prev_items(indices_copy, offset)
	return len(indices) == len(indices_copy)
}

/* // evaluate_solution is a helper function that evaluates the solution.
//
// Parameters:
//   - curr: The current token.
//   - solutions: The solutions.
//
// Returns:
//   - []internal.ActionType: The evaluated solution.
func (d decider[T]) evaluate_solution(curr T, solutions []int) []internal.ActionType {
	var final_sol []internal.ActionType

	for _, sol := range solutions {
		item := d.item_list[sol]

		final_sol = append(final_sol, item.act)
	}

	return final_sol
}
*/
