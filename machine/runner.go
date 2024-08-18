package machine

import (
	"fmt"
	"io"
	"iter"
)

type MachineStepRunner[T any] struct {
	table      map[SystemState]StepFunc[T]
	transition map[SystemState]SystemState
}

func (msr MachineStepRunner[T]) All(scanner io.RuneScanner) iter.Seq2[StepFunc[T], *rune] {
	state := InitSS

	var fn func(yield func(StepFunc[T], *rune) bool)

	if scanner == nil {
		fn = func(yield func(StepFunc[T], *rune) bool) {
			for state != EndSS {
				step_fn, ok := msr.table[state]
				if !ok {
					panic(fmt.Sprintf("Invalid state: %d", state))
				}

				if !yield(step_fn, nil) {
					break
				}

				tmp, ok := msr.transition[state]
				if !ok {
					panic(fmt.Sprintf("Missing transition from state: %d", state))
				}

				state = tmp
			}
		}
	} else {
		fn = func(yield func(StepFunc[T], *rune) bool) {
			for state != EndSS {
				step_fn, ok := msr.table[state]
				if !ok {
					panic(fmt.Sprintf("Invalid state: %d", state))
				}

				var char *rune

				if state != InitSS {
					c, _, err := scanner.ReadRune()
					if err == nil {
						char = &c
					} else if err != io.EOF {
						panic(err)
					}
				}

				if !yield(step_fn, char) {
					break
				}

				tmp, ok := msr.transition[state]
				if !ok {
					panic(fmt.Sprintf("Missing transition from state: %d", state))
				}

				state = tmp
			}
		}
	}

	return fn
}

/*
func (msr MachineStepRunner[T]) Run(info T, scanner io.RuneScanner) error {
	state, err := msr.init_fn(info)
	if err != nil {
		return err
	}

	states := make([]SystemState, 0, 1)
	states = append(states, state)

	for {
		state := states[0]
		states = states[1:]

		if state == EndSS {
			break
		}

		step_fn, ok := msr.table[state]
		if !ok {
			return fmt.Errorf("invalid state: %d", state)
		}

		var char *rune

		c, _, err := scanner.ReadRune()
		if err == nil {
			char = &c
		} else if err != io.EOF {
			return err
		}

		state, err = step_fn(info, char)
		if err != nil {
			return err
		}

		states = append(states, state)
	}

	return nil
} */
