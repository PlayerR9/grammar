package grammar

import (
	"iter"

	gcers "github.com/PlayerR9/go-commons/errors"
	utstk "github.com/PlayerR9/go-commons/stack"
	dbg "github.com/PlayerR9/go-debug/assert"
	internal "github.com/PlayerR9/grammar/grammar/internal"
)

// DecisionFn is the decision function.
//
// Parameters:
//   - ap: The active parser.
//
// Returns:
//   - []*Item[T]: The list of items.
//   - error: An error.
type DecisionFn[T internal.TokenTyper] func(ap *ActiveParser[T]) ([]*Item[T], error)

// Parser is the grammar parser.
type Parser[T internal.TokenTyper] struct {
	// tokens is the token stream.
	tokens []*Token[T]

	// rule_set is the rule set.
	rule_set *RuleSet[T]

	// table is the parsing table.
	table *ParseTable[T]

	// decision_fn is the decision function.
	decision_fn DecisionFn[T]
}

// NewParser creates a new parser with the given rule set.
//
// Parameters:
//   - rule_set: The rule set.
//
// Returns:
//   - *Parser[T]: The new parser.
//   - error: An error of type *errors.ErrInvalidParameter if rule_set is nil.
func NewParser[T internal.TokenTyper](rule_set *RuleSet[T]) (*Parser[T], error) {
	if rule_set == nil {
		return nil, gcers.NewErrNilParameter("rule_set")
	}

	pt := NewParseTable(rule_set.rules)
	err := pt.Init()
	if err != nil {
		return nil, err
	}

	return &Parser[T]{
		rule_set: rule_set,
		table:    pt,
	}, nil
}

// NewParserWithFunc creates a new parser with the given rule set.
//
// Parameters:
//   - decision_fn: The decision function.
//
// Returns:
//   - *Parser[T]: The new parser.
//   - error: An error of type *errors.ErrInvalidParameter if rule_set is nil.
func NewParserWithFunc[T internal.TokenTyper](decision_fn DecisionFn[T]) (*Parser[T], error) {
	if decision_fn == nil {
		return nil, gcers.NewErrNilParameter("decision_fn")
	}

	return &Parser[T]{
		decision_fn: decision_fn,
	}, nil
}

// Parse is the main function of the parser.
//
// Parameters:
//   - tokens: The tokens to be parsed.
//
// Returns:
//   - *ActiveParser[T]: The parser.
//   - error: An error if any.
func (p *Parser[T]) Parse(tokens []*Token[T]) iter.Seq[*ActiveParser[T]] {
	var fn func(yield func(*ActiveParser[T]) bool)

	p.tokens = tokens

	if p.decision_fn == nil {
		fn = func(yield func(*ActiveParser[T]) bool) {
			active, err := NewActiveParser(p, nil)
			dbg.AssertErr(err, "NewActiveParser(p, nil)")

			var invalid_parsers []*ActiveParser[T]

			stack := utstk.NewStack(active)

			for {
				top, ok := stack.Pop()
				if !ok {
					break
				}

				ok = top.WalkAll()
				if !ok {
					invalid_parsers = append(invalid_parsers, top)

					continue
				}

				possible_paths := top.Exec()
				dbg.AssertNotNil(possible_paths, "possible_paths")

				for _, path := range possible_paths {
					stack.Push(path)
				}

				if top.HasError() {
					invalid_parsers = append(invalid_parsers, top)
				} else if !yield(top) {
					return
				}
			}

			// For now we will assume that the last invalid parser is the most likely error.
			last_invalid := invalid_parsers[len(invalid_parsers)-1]

			_ = yield(last_invalid)
		}
	} else {
		fn = func(yield func(*ActiveParser[T]) bool) {
			active, err := NewActiveParser(p, nil)
			dbg.AssertErr(err, "NewActiveParser(p, nil)")

			var invalid_parsers []*ActiveParser[T]

			stack := utstk.NewStack(active)

			for {
				top, ok := stack.Pop()
				if !ok {
					break
				}

				ok = top.WalkAll()
				if !ok {
					invalid_parsers = append(invalid_parsers, top)

					continue
				}

				possible_paths := top.ExecWithFn()
				dbg.AssertNotNil(possible_paths, "possible_paths")

				for _, path := range possible_paths {
					stack.Push(path)
				}

				if top.HasError() {
					invalid_parsers = append(invalid_parsers, top)
				} else if !yield(top) {
					return
				}
			}

			// For now we will assume that the last invalid parser is the most likely error.
			last_invalid := invalid_parsers[len(invalid_parsers)-1]

			_ = yield(last_invalid)
		}
	}

	return fn
}
