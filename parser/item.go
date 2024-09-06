package parser

import (
	"iter"
	"strings"

	gccmp "github.com/PlayerR9/go-commons/cmp"
	gcers "github.com/PlayerR9/go-commons/errors"
	gcint "github.com/PlayerR9/go-commons/ints"
	dbg "github.com/PlayerR9/go-debug/assert"
	"github.com/PlayerR9/grammar/internal"
)

// Item is an item in the parsing table.
type Item[T internal.TokenTyper] struct {
	// rule is the rule.
	rule *Rule[T]

	// pos is the position.
	pos int

	// act is the action.
	act internal.ActionType

	// prevs is the set of previous items.
	prevs *gccmp.Set[T]

	// lookaheads is the set of lookaheads.
	lookaheads []*gccmp.Set[T]
}

// Equals implements the pkg.Type interface.
//
// Two items are equal if their rules are equal, they are not nil, and their positions are equal.
func (item *Item[T]) Equals(other *Item[T]) bool {
	if other == nil {
		return false
	}

	return item.rule.Equals(other.rule) && item.pos == other.pos
}

// String implements the fmt.Stringer interface.
func (item Item[T]) String() string {
	var elems []string

	i := item.rule.Size()

	for rhs := range item.rule.Backwards() {
		i--

		if i == item.pos {
			elems = append(elems, "[", rhs.String(), "]")
		} else {
			elems = append(elems, rhs.String())
		}
	}

	elems = append(elems, "->", item.rule.Lhs().String(), ";", "("+item.act.String()+")")

	if len(item.lookaheads) > 0 {
		elems = append(elems, "---")

		for _, la := range item.lookaheads {
			elems = append(elems, "[", la.String(), "]")
		}
	}

	if item.prevs.Len() > 0 {
		elems = append(elems, "---")

		for prev := range item.prevs.All() {
			elems = append(elems, prev.String())
		}
	}

	return strings.Join(elems, " ")
}

// NewItem creates a new item.
//
// Parameters:
//   - rule: The rule.
//   - pos: The position.
//
// Returns:
//   - *Item[T]: The created item.
//   - error: An error if the position is invalid or the rule is nil.
func NewItem[T internal.TokenTyper](rule *Rule[T], pos int) (*Item[T], error) {
	if rule == nil {
		return nil, gcers.NewErrNilParameter("rule")
	}

	size := rule.Size()

	if pos < 0 || pos > size {
		return nil, gcers.NewErrInvalidParameter("pos", gcint.NewErrOutOfBounds(pos, 0, size))
	}

	var act internal.ActionType

	if pos == size {
		rhs, ok := rule.RhsAt(pos - 1)
		dbg.AssertOk(ok, "rule.RhsAt(%d)", pos)

		if rhs == T(0) {
			act = internal.ActAcceptType
		} else {
			act = internal.ActReduceType
		}
	} else {
		act = internal.ActShiftType
	}

	return &Item[T]{
		rule: rule,
		pos:  pos,
		act:  act,
	}, nil
}

// IsShift checks if the item is a shift.
//
// Returns:
//   - bool: True if the item is a shift, otherwise false.
func (item Item[T]) IsShift() bool {
	return item.pos < item.rule.Size()
}

// IsReduce checks if the item is a reduce.
//
// Returns:
//   - bool: True if the item is a reduce, otherwise false.
func (item *Item[T]) IncreaseLookbehind() bool {
	if item.pos == 0 {
		return false
	}

	pos := item.pos - item.prevs.Len() - 1

	prev, ok := item.rule.RhsAt(pos)
	if !ok {
		return false
	}

	item.prevs.Add(prev)

	return true
}

// IsInConflictWith checks if the item is in conflict with another item.
//
// Returns:
//   - bool: True if the item is in conflict with another item, otherwise false.
func (item Item[T]) IsInConflictWith(other *Item[T]) bool {
	if other == nil {
		return false
	}

	if item.IsShift() && other.IsShift() {
		return false
	}

	if !item.prevs.Equals(other.prevs) {
		return false
	}

	if len(item.lookaheads) != len(other.lookaheads) {
		return false
	}

	for i, la := range item.lookaheads {
		if !la.Equals(other.lookaheads[i]) {
			return false
		}
	}

	return true
}

// RhsAt returns the right hand side at the given position.
//
// Parameters:
//   - pos: The position.
//
// Returns:
//   - T: The right hand side.
//   - bool: True if the right hand side exists, otherwise false.
func (item Item[T]) RhsAt(pos int) (T, bool) {
	return item.rule.RhsAt(pos)
}

// AppendLookahead appends the given lookahead set to the item.
//
// Parameters:
//   - ls: The lookahead set.
//
// Returns:
//   - error: An error if the lookahead set is nil.
func (item *Item[T]) AppendLookahead(ls *gccmp.Set[T]) error {
	if ls == nil {
		return gcers.NewErrNilParameter("ls")
	}

	dbg.AssertThat("ls", dbg.NewOrderedAssert(ls.Len()).Equal(0)).Not().Panic()

	item.lookaheads = append(item.lookaheads, ls)

	return nil
}

// LookaheadAt returns the lookahead set at the given position.
//
// Parameters:
//   - pos: The position.
//
// Returns:
//   - *gccmp.Set[T]: The lookahead set.
//   - bool: True if the lookahead set exists, otherwise false.
func (item Item[T]) LookaheadAt(pos int) (*gccmp.Set[T], bool) {
	if pos < 0 || pos >= len(item.lookaheads) {
		return nil, false
	}

	elem := item.lookaheads[pos]

	return elem, true
}

// Rhs returns an iterator over the right hand side of the item.
//
// Returns:
//   - iter.Seq[T]: The iterator. Never returns nil.
func (item Item[T]) Rhs() iter.Seq[T] {
	return item.rule.Rhs()
}

// Lhs returns the left hand side of the item.
//
// Returns:
//   - T: The left hand side.
func (item Item[T]) Lhs() T {
	return item.rule.Lhs()
}

// Pos returns the position of the item in the rule.
//
// Returns:
//   - int: The position.
func (item Item[T]) Pos() int {
	return item.pos
}

// Advance advances the position of the item in the rule by one.
//
// Returns:
//   - bool: True if the position was advanced, otherwise false.
func (item *Item[T]) Advance() (*Item[T], bool) {
	if item.pos == item.rule.Size() {
		return item, false
	}

	return &Item[T]{
		rule:       item.rule,
		pos:        item.pos + 1,
		act:        item.act,
		lookaheads: item.lookaheads,
		prevs:      item.prevs,
	}, true
}
