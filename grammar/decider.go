package grammar

import (
	"fmt"

	utst "github.com/PlayerR9/go-commons/cmp"
	gcers "github.com/PlayerR9/go-commons/errors"
	gcslc "github.com/PlayerR9/go-commons/slices"
	dbg "github.com/PlayerR9/go-debug/assert"
	internal "github.com/PlayerR9/grammar/grammar/internal"
)

type Decider[T internal.TokenTyper] struct {
	p         *ActiveParser[T]
	item_list []*Item[T]
	err       error
}

func NewDecider[T internal.TokenTyper](p *ActiveParser[T], item_list []*Item[T]) (*Decider[T], error) {
	if p == nil {
		return nil, gcers.NewErrNilParameter("p")
	}

	return &Decider[T]{
		p:         p,
		item_list: item_list,
		err:       nil,
	}, nil
}

func (d *Decider[T]) FilterLookaheads(indices []int, prev T, top1 *Token[T]) ([]int, []int) {
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

func (d *Decider[T]) ApplyPopRule(indices []int, prev T, offset int) ([]int, T) {
	dbg.AssertThat("offset", dbg.NewOrderedAssert(offset).GreaterThan(0)).Panic()

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
		d.err = NewErrUnexpectedToken(&prev, nil, expected.Slice()...)
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

func (d *Decider[T]) Decision(indices []int, prev T) ([]int, error) {
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

func (d Decider[T]) filter_no_prev_items(indices []int, offset int) []int {
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

func (d Decider[T]) OnlyLookaheads(indices []int, offset int) bool {
	if offset < 1 {
		return false
	}

	indices_copy := make([]int, len(indices))
	copy(indices_copy, indices)

	indices_copy = d.filter_no_prev_items(indices_copy, offset)
	return len(indices) == len(indices_copy)
}

func (d Decider[T]) EvaluateSolution(curr T, solutions []int) []internal.ActionType {
	var final_sol []internal.ActionType

	for _, sol := range solutions {
		item := d.item_list[sol]

		final_sol = append(final_sol, item.act)
	}

	return final_sol
}
