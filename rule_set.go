package grammar

import (
	"fmt"
	"slices"
	"strings"

	utst "github.com/PlayerR9/go-commons/cmp"
	dbg "github.com/PlayerR9/go-debug/assert"
	internal "github.com/PlayerR9/grammar/internal"
)

// RuleSet is the rule set data structure.
type RuleSet[T internal.TokenTyper] struct {
	// rules is the list of all rules in the grammar.
	rules []*Rule[T]

	// items is the list of all items in the grammar.
	items map[T][]*Item[T]

	// symbols is the list of all symbols in the grammar.
	symbols *utst.Set[T]
}

// String implements the fmt.Stringer interface.
//
// Format:
//
//	// {rule}
//	// {rule}
//	...
func (rs RuleSet[T]) String() string {
	var elems []string

	for symbol := range rs.symbols.All() {
		item_list, ok := rs.items[symbol]
		dbg.AssertOk(ok, "items[%q]", symbol.String())

		vals := make([]string, 0, len(item_list))

		for _, item := range item_list {
			vals = append(vals, "// "+item.String())
		}

		elems = append(elems, strings.Join(vals, "\n"))
	}

	return strings.Join(elems, "\n\n")
}

// NewRuleSet creates a new RuleSet.
//
// Returns:
//   - *RuleSet[T]: The created RuleSet. Never returns nil.
func NewRuleSet[T internal.TokenTyper]() *RuleSet[T] {
	return &RuleSet[T]{
		rules:   make([]*Rule[T], 0),
		items:   make(map[T][]*Item[T]),
		symbols: utst.NewSet[T](),
	}
}

// MustAddRule adds a new rule to the rule set.
//
// Parameters:
//   - rule: The rule to add.
func (rs *RuleSet[T]) MustAddRule(rule *Rule[T]) {
	if rule == nil {
		return
	}

	if slices.ContainsFunc(rs.rules, rule.Equals) {
		panic("rule already exists")
	}

	rs.rules = append(rs.rules, rule)
}

// MustMakeRule adds a new rule to the rule set.
//
// Panics if the rule already exists or if the rhss is empty.
//
// Parameters:
//   - lhs: The left hand side of the rule.
//   - rhss: The right hand side of the rule.
func (rs *RuleSet[T]) MustMakeRule(lhs T, rhss []T) {
	rule, err := NewRule(lhs, rhss)
	dbg.AssertErr(err, "NewRule(%q, rhss)", lhs.String())

	if slices.ContainsFunc(rs.rules, rule.Equals) {
		panic("rule already exists")
	}

	rs.rules = append(rs.rules, rule)
}

// DetermineSymbols determines the symbols in the rule set.
func (rs *RuleSet[T]) DetermineSymbols() {
	rs.symbols = utst.NewSet[T]()

	for _, rule := range rs.rules {
		rs.symbols.Union(rule.Symbols())
	}
}

// DetermineItems determines the items in the rule set.
func (rs *RuleSet[T]) DetermineItems() {
	rs.DetermineSymbols()

	item_table := make(map[T][]*Item[T])

	for symbol := range rs.symbols.All() {
		var item_list []*Item[T]

		for _, rule := range rs.rules {
			indices := rule.IndicesOf(symbol)
			if len(indices) == 0 {
				continue
			}

			for _, idx := range indices {
				item, err := NewItem(rule, idx)
				dbg.AssertErr(err, "NewItem(rule, %d)", idx)

				item_list = append(item_list, item)
			}
		}

		item_table[symbol] = item_list
	}

	rs.items = item_table
}

// solve_lookbehinds is a helper function that solves the lookbehinds.
func (rs *RuleSet[T]) solve_lookbehinds() {
	cm := NewConflictMap[T]()
	defer cm.Cleanup()

	for {
		cannot_continue := true

		cm.Init(rs.items)

		if cm.Len() == 0 {
			break
		}

		for _, item := range cm.Entry() {
			ok := item.IncreaseLookbehind()
			if ok {
				cannot_continue = false
			}
		}

		if cannot_continue {
			break
		}
	}
}

// DetermineLookaheads determines the lookaheads with a specific offset of the specified item.
//
// Parameters:
//   - item: The item to determine the lookaheads for.
//   - offset: The offset to determine the lookaheads for.
//
// Note: The offset must be greater than 0.
func (rs RuleSet[T]) DetermineLookaheads(item *Item[T], offset int) {
	dbg.AssertThat("offset", dbg.NewOrderedAssert(offset).GreaterOrEqualThan(1)).Panic()

	if item == nil {
		return
	}

	next_rhs, ok := item.RhsAt(item.pos + offset)
	if !ok {
		return
	}

	solution := utst.NewSet[T]()

	todo := []T{next_rhs}
	seen := make(map[T]bool)

	for len(todo) > 0 {
		first := todo[0]
		todo = todo[1:]

		prev, ok := seen[first]
		if ok && prev {
			continue
		}

		seen[first] = true

		is_terminal := first.IsTerminal()

		if is_terminal {
			solution.Add(first)
			continue
		}

		nexts := rs.RulesWithLhs(first)

		lookaheads := ExtractRhsAt(nexts, 0)

		for _, lookahead := range lookaheads {
			is_teminal := lookahead.IsTerminal()

			if is_teminal {
				solution.Add(lookahead)
			} else {
				todo = append(todo, lookahead)
			}
		}
	}

	item.AppendLookahead(solution)
}

// solve_lookaheads is a helper function that solves the lookaheads.
func (rs *RuleSet[T]) solve_lookaheads() {
	cm := NewConflictMap[T]()
	defer cm.Cleanup()

	offset := 1

	for {
		cm.Init(rs.items)

		if cm.Len() == 0 {
			break
		}

		for _, item := range cm.Entry() {
			rs.DetermineLookaheads(item, offset)
		}

		offset++
	}
}

// SolveConflicts solves the conflicts in the rule set.
//
// Returns:
//   - bool: True if all conflicts were solved. False otherwise.
//
// If conflicts are not solved, this function will print out the conflicts.
func (rs *RuleSet[T]) SolveConflicts() bool {
	rs.solve_lookbehinds()
	rs.solve_lookaheads()

	cm := NewConflictMap[T]()
	defer cm.Cleanup()

	cm.Init(rs.items)

	if cm.Len() == 0 {
		return true
	}

	fmt.Println("Conflicts detected:")

	for symbol, items := range cm.All() {
		fmt.Println("\t" + symbol.String() + ":")

		for item := range items {
			fmt.Println("\t\t" + item.String())
		}

		fmt.Println()
	}

	fmt.Println()

	return false
}

// RulesWithLhs returns the rules with the specified left hand side.
//
// Parameters:
//   - lhs: The left hand side of the rules to return.
//
// Returns:
//   - []*internal.Rule[T]: The rules with the specified left hand side.
func (rs RuleSet[T]) RulesWithLhs(lhs T) []*Rule[T] {
	var rules []*Rule[T]

	for _, rule := range rs.rules {
		if rule.Lhs() == lhs {
			rules = append(rules, rule)
		}
	}

	return rules
}

func (rs RuleSet[T]) Decision(p *ActiveParser[T]) ([]*Item[T], error) {
	dbg.AssertNotNil(p, "p")

	top1, ok := p.Pop()
	dbg.AssertOk(ok, "p.Pop()")

	item_list, ok := rs.items[top1.Type]
	if !ok {
		return nil, fmt.Errorf("unexpected token: %s", top1.Type.String())
	}

	indices := make([]int, 0, len(item_list))

	for i := range item_list {
		indices = append(indices, i)
	}

	curr := top1.Type

	d, err := NewDecider(p, item_list)
	dbg.AssertErr(err, "NewDecider(p, item_list)")

	indices, err = d.Decision(indices, curr)
	if err != nil {
		return nil, err
	}

	offset := 1

	var solutions []int

	for {
		if len(indices) == 1 {
			solutions = indices
		} else if d.OnlyLookaheads(indices, offset) {
			indices, solutions = d.FilterLookaheads(indices, curr, top1)
		}

		if len(solutions) > 0 {
			break
		}

		indices, curr = d.ApplyPopRule(indices, curr, offset)

		indices, err = d.Decision(indices, curr)
		if err != nil {
			return nil, err
		}

		offset++
	}

	items := make([]*Item[T], 0, len(solutions))

	for _, sol := range solutions {
		item := item_list[sol]
		items = append(items, item)
	}

	return items, nil
}
