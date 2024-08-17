package matcher

import (
	"errors"
	"fmt"
	"go/scanner"
	"io"

	gcers "github.com/PlayerR9/go-commons/errors"
	gcslc "github.com/PlayerR9/go-commons/slices"
	dbg "github.com/PlayerR9/go-debug/assert"
	gr "github.com/PlayerR9/grammar/grammar"
)

// LexerMatcher is an interface that defines the behavior of a lexer matcher.
type LexerMatcher interface {
	// FirstMatch matches the given character and it is responsible for
	// initializing the state of the matcher in order to guarantee successful
	// calls to Match() function.
	//
	// Parameters:
	//   - char: The character to match.
	//
	// Returns:
	//   - error: An error that specifies why the char does not match and the
	//     possible suggestions for solving the error.
	FirstMatch(char rune) error

	// Match matches the given character. The result depends on previous calls
	// to this method.
	//
	// Parameters:
	//   - char: The character to match.
	//
	// Returns:
	//   - bool: True if the rule is completely matched, false otherwise.
	//   - error: An error if the char does not match and the rule is in an invalid
	//     state or if something terrible happens internally.
	Match(char rune) (bool, error)

	// IsValid checks if the current state of the matcher is valid.
	//
	// Returns:
	//   - bool: True if the matcher is valid, false otherwise.
	IsValid() bool
}

type Rule[S gr.TokenTyper] struct {
	symbol S
	rule   LexerMatcher
}

type RuleMatched struct {
	symbol int
	chars  []rune
}

type MatcherState int

const (
	CriticalState MatcherState = iota - 2
	EndState
	InitialState
	IndexLoadingState
	FilteringStepState
)

type RuleMatcher[S gr.TokenTyper] struct {
	rules         []Rule[S]
	indices       []int
	eval          gcers.ErrOrSol[*RuleMatched]
	chars         []rune
	err           error
	eof_completed bool
	scanner       io.RuneScanner
}

// Init implements the MachineStater interface.
func (rm *RuleMatcher[S]) Init() {
	rm.indices = rm.indices[:0]
	rm.eval = gcers.ErrOrSol[*RuleMatched]{}
	// When go-commons is updated, replace the following line with: rm.eval.Reset()
	rm.chars = rm.chars[:0]

	rm.eof_completed = false

	rm.err = nil
	rm.scanner = nil
}

func (rm *RuleMatcher[S]) Execute(state MatcherState, args ...any) MatcherState {
	switch state {
	case InitialState:
		if len(args) == 0 {
			return EndState, errors.New("required first argument to be of type io.RuneScanner, got nil instead")
		}

		scanner, ok := args[0].(io.RuneScanner)
		if !ok {
			return EndState, fmt.Errorf("required first argument to be of type io.RuneScanner, got %T instead", args[0])
		}

		rm.scanner = scanner

		return IndexLoadingState, nil
	case IndexLoadingState:
		/* if scanner == nil {
			rm.err = NewErrNoMatch(errors.New("empty scanner"))

			return EndState
		} */

		char, _, err := rm.scanner.ReadRune()
		if err == io.EOF {
			rm.err = NewErrNoMatch(err)

			return EndState, nil
		}

		if err != nil {
			rm.err = err

			return EndState, nil
		}

		for i, rule := range rm.rules {
			err := rule.rule.FirstMatch(char)
			if err == nil {
				rm.indices = append(rm.indices, i)
			} else {
				rm.eval.AddErr(err, 0)
			}
		}

		if len(rm.indices) > 0 {
			rm.chars = append(rm.chars, char)

			return FilteringStepState
		}

		err = scanner.UnreadRune()
		dbg.AssertErr(err, "scanner.UnreadRune()")

		errs := rm.eval.Errors()

		rm.err = NewErrNoMatch(errs[0])

		return EndState // TODO: handle multiple errors
	case FilteringStepState:
		if len(rm.indices) == 0 {
			if !rm.eof_completed {
				err := scanner.UnreadRune()
				dbg.AssertErr(err, "scanner.UnreadRune()")

				rm.chars = rm.chars[:len(rm.chars)-1]
			} else {
				fn := func(idx int) bool {
					rule := rm.rules[idx]

					return rule.rule.IsValid()
				}

				tmp, ok := gcslc.SFSeparateEarly(rm.indices, fn)
				if !ok {
					rm.err = NewErrNoMatch(rm.make_error())
					return EndState
				}

				rm.indices = tmp
			}

			return EndState
		}

		char, _, err := scanner.ReadRune()
		if err == io.EOF {
			rm.eof_completed = true

			break
		} else if err != nil {
			rm.err = err

			return EndState
		}

		rm.step(char)

		return FilteringStepState
	case EndState:
		// Do nothing
	default:
		rm.err = fmt.Errorf("invalid state: %d", state)

		return EndState
	}
}

func (rm *RuleMatcher[S]) GetSolution() ([]*RuleMatched, error) {
	ok := rm.eval.HasError()
	if !ok {
		sols := rm.eval.Solutions()

		return sols, nil
	}

	errs := rm.eval.Errors()

	return nil, NewErrNoMatch(errs[0]) // TODO: deal with multiple errors
}

func (rm *RuleMatcher[S]) AddRule(symbol S, rule LexerMatcher) {
	if rule == nil {
		return
	}

	rm.rules = append(rm.rules, Rule[S]{
		symbol: symbol,
		rule:   rule,
	})
}

func (rm *RuleMatcher[S]) step(char rune) {
	level := len(rm.chars)
	rm.chars = append(rm.chars, char)

	fn := func(idx int) bool {
		dbg.AssertThat("idx", idx).InRange(0, len(rm.indices)-1).Panic()
		rule := rm.rules[idx]

		ok, err := rule.rule.Match(char)
		if err == nil && !ok {
			return true
		}

		if err != nil {
			rm.eval.AddErr(err, level)

			return false
		}

		match := &RuleMatched{
			symbol: int(rule.symbol),
			chars:  make([]rune, len(rm.chars)),
		}
		copy(match.chars, rm.chars)

		rm.eval.AddSol(match, level)

		return false
	}

	rm.indices = gcslc.SliceFilter(rm.indices, fn)
}

func (rm *RuleMatcher[S]) make_error() error {
	return nil
}
