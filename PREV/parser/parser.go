package parser

import (
	"iter"

	util "github.com/PlayerR9/go-commons/backup"
	gcers "github.com/PlayerR9/go-commons/errors"
	"github.com/PlayerR9/go-commons/stack"
	dbg "github.com/PlayerR9/go-debug/assert"
	gr "github.com/PlayerR9/grammar/PREV/grammar"
	"github.com/PlayerR9/grammar/PREV/internal"
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
	tokens []*gr.Token[T]

	// rule_set is the rule set.
	rule_set *RuleSet[T]

	// table is the parsing table.
	table *parse_table[T]

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

	pt := new_parse_table(rule_set.rules)
	err := pt.init()
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

// new_active_parser creates a new active parser.
//
// Parameters:
//   - global: The shared information between active parsers.
//
// Returns:
//   - *ActiveParser: The new active parser.
//   - error: An error if shifting the first token failed.
func (p *Parser[T]) active_parser_of() *ActiveParser[T] {
	dbg.AssertThat("len(p.tokens)", dbg.NewOrderedAssert(len(p.tokens)).GreaterThan(0)).Panic()

	tokens := make([]*gr.Token[T], 0, len(p.tokens))
	for i := 0; i < len(p.tokens); i++ {
		tokens = append(tokens, p.tokens[i].Copy())
	}

	for i := 0; i < len(tokens)-1; i++ {
		tokens[i].Lookahead = tokens[i+1]
	}

	new_ap := &ActiveParser[T]{
		global:         p,
		reader:         gr.NewTokenStream(tokens),
		token_stack:    stack.NewRefusableStack[*gr.Token[T]](),
		err:            nil,
		possible_cause: nil,
	}

	err := new_ap.shift() // initial shift
	if err != nil {
		new_ap.err = err

		return nil
	}

	return new_ap
}

// Parse is the main function of the parser.
//
// Parameters:
//   - tokens: The tokens to be parsed.
//
// Returns:
//   - *ActiveParser[T]: The parser.
//   - error: An error if any.
func (p *Parser[T]) Parse(tokens []*gr.Token[T]) iter.Seq[*ActiveParser[T]] {
	p.tokens = tokens

	return util.Execute(p.active_parser_of)
}
