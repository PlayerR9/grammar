package grammar

/*

// GotoGraph is the goto graph data structure.
type GotoGraph[T internal.TokenTyper] struct {
	all_symbols *cmp.Set[T]

	states []*State[T]

	// goto_table is a map[state_id]map[non-terminal symbol]next_state_id
	goto_table map[int]map[T]int

	// action_table is a map[state_id]map[terminal symbol]action
	action_table map[int]map[T]ActionType
}

func NewGotoGraph[T internal.TokenTyper](is *set.Set[*Item[T]], initial_symbol T) (*GotoGraph[T], error) {
	if initial_symbol.IsTerminal() {
		return nil, fmt.Errorf("initial symbol (%q) must not be terminal", initial_symbol.String())
	}

	initial_states := GetItemsWithLhs(is, initial_symbol)
	if len(initial_states) == 0 {
		return nil, fmt.Errorf("no items with lhs %q", initial_symbol.String())
	} else if len(initial_states) > 1 {
		return nil, fmt.Errorf("multiple items with lhs %q", initial_symbol.String())
	}

	state_0 := NewState[T](is, initial_states[0])

	for symbol := range






	state_queue := queue.NewQueueWithElems([]*State[T]{state_0})
	state_done := []*State[T]{state_0}

	table := make(map[int]map[int]int)

	for {
		first, ok := state_queue.Dequeue()
		if !ok {
			break
		}

		goto_map := make(map[int]int)

		for rule_idx, seed := range first.Rule() {
			state_idx := -1

			for i := 0; i < len(state_done) && state_idx == -1; i++ {
				if state_done[i].IsOfSeed(seed) {
					state_idx = i
				}
			}

			if state_idx == -1 {
				new_state := NewState(is, seed)

				state_queue.Enqueue(new_state)
				state_done = append(state_done, new_state)

				state_idx = len(state_done) - 1
			}

			goto_map[rule_idx] = state_idx
		}
	}

	return &GotoGraph[T]{
		states: state_done,
		table:  table,
	}, nil
}
*/
