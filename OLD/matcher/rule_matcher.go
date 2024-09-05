package matcher

import (
	"errors"
	"io"
	"iter"

	gcers "github.com/PlayerR9/go-commons/errors"
	gcslc "github.com/PlayerR9/go-commons/slices"
	grmch "github.com/PlayerR9/grammar/OLD/machine"
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

type Rule struct {
	symbol int
	rule   LexerMatcher
}

type RuleMatched struct {
	symbol int
	chars  []rune
}

const (
	IndexLoadingState grmch.SystemState = iota + 1
	FilteringStepState
)

var (
	rm_fsm *grmch.MachineState[*RuleMatcherInfo]
)

func init() {
	init_fn := func(info *RuleMatcherInfo, _ *rune) (grmch.SystemState, error) {
		info.indices = info.indices[:0]

		info.eval = gcers.ErrOrSol[*RuleMatched]{}
		// When go-commons is updated, replace the following line with: rm.eval.Reset()

		info.chars = info.chars[:0]

		info.eof_completed = false
		info.err = nil

		return IndexLoadingState, nil
	}

	cleanup_fn := func(info *RuleMatcherInfo) {
		if len(info.indices) > 0 {
			info.indices = info.indices[:0]
		}
		info.indices = nil

		info.eval = gcers.ErrOrSol[*RuleMatched]{}
		// When go-commons is updated, replace the following line with: rm.eval.Reset()

		if len(info.chars) > 0 {
			info.chars = info.chars[:0]
		}
		info.chars = nil

		info.eof_completed = false
		info.err = nil
	}

	rm_fsm = grmch.NewMachineState(init_fn)
	rm_fsm.WithCleanup(cleanup_fn)

	index_loading_fn := func(info *RuleMatcherInfo, char *rune) (grmch.SystemState, error) {
		if char == nil {
			info.err = NewErrNoMatch(io.EOF)

			return grmch.EndSS, nil
		}

		c := *char

		for i, rule := range info.global.All() {
			err := rule.FirstMatch(c)
			if err == nil {
				info.indices = append(info.indices, i)
			} else {
				info.eval.AddErr(err, 0)
			}
		}

		if len(info.indices) > 0 {
			info.chars = append(info.chars, c)
			return FilteringStepState, nil
		}

		/* 	err = scanner.UnreadRune()
		dbg.AssertErr(err, "scanner.UnreadRune()") */

		errs := info.eval.Errors()

		info.err = NewErrNoMatch(errs[0]) // TODO: handle multiple errors

		return grmch.EndSS, nil
	}

	rm_fsm.AddState(IndexLoadingState, index_loading_fn)

	filtering_step_fn := func(info *RuleMatcherInfo, char *rune) (grmch.SystemState, error) {
		if len(info.indices) == 0 {
			if !info.eof_completed {
				/* err := scanner.UnreadRune()
				dbg.AssertErr(err, "scanner.UnreadRune()") */

				info.chars = info.chars[:len(info.chars)-1]
			} else {
				fn := func(idx int) bool {
					rule := info.global.get_rule_at(idx)
					return rule.IsValid()
				}

				tmp, ok := gcslc.SFSeparateEarly(info.indices, fn)
				if !ok {
					info.err = NewErrNoMatch(errors.New("[MAKE SUGGESTIONS HERE]"))
					return grmch.EndSS, nil
				}

				info.indices = tmp
			}

			return grmch.EndSS, nil
		}

		if char == nil {
			info.eof_completed = true

			return grmch.EndSS, nil
		}

		level := len(info.indices)
		info.chars = append(info.chars, *char)

		fn := func(idx int) bool {
			// dbg.AssertThat("idx", idx).InRange(0, len(info.indices)-1).Panic()
			rule := info.global.get_rule_at(idx)

			ok, err := rule.Match(*char)
			if err == nil && !ok {
				return true
			}

			if err != nil {
				info.eval.AddErr(err, level)

				return false
			}

			match := &RuleMatched{
				symbol: info.global.get_symbol_at(idx),
				chars:  make([]rune, len(info.chars)),
			}
			copy(match.chars, info.chars)

			info.eval.AddSol(match, level)

			return false
		}

		info.indices = gcslc.SliceFilter(info.indices, fn)

		return FilteringStepState, nil
	}

	rm_fsm.AddState(FilteringStepState, filtering_step_fn)
}

type RuleMatcherInfo struct {
	global        *RuleMatcher
	indices       []int
	eval          gcers.ErrOrSol[*RuleMatched]
	chars         []rune
	err           error
	eof_completed bool
}

type RuleMatcher struct {
	rules []Rule
}

func (rm *RuleMatcher) AddRule(symbol int, rule LexerMatcher) {
	if rule == nil {
		return
	}

	rm.rules = append(rm.rules, Rule{
		symbol: symbol,
		rule:   rule,
	})
}

func (rm *RuleMatcher) Match(scanner io.RuneScanner) ([]*RuleMatched, error) {
	run, clean := rm_fsm.Make(&RuleMatcherInfo{
		global: rm,
	})
	defer clean()

	res, err := run(scanner)
	if err != nil {
		return nil, err
	}

	ok := res.eval.HasError()
	if !ok {
		sols := res.eval.Solutions()

		return sols, nil
	}

	errs := res.eval.Errors()

	return nil, NewErrNoMatch(errs[0]) // TODO: deal with multiple errors
}

func (rm *RuleMatcher) All() iter.Seq2[int, LexerMatcher] {
	return func(yield func(int, LexerMatcher) bool) {
		for i, rule := range rm.rules {
			yield(i, rule.rule)
		}
	}
}

func (rm *RuleMatcher) get_rule_at(idx int) LexerMatcher {
	rule := rm.rules[idx]
	return rule.rule
}

func (rm *RuleMatcher) get_symbol_at(idx int) int {
	rule := rm.rules[idx]
	return rule.symbol
}
