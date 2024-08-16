package parsing

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	gcos "github.com/PlayerR9/go-commons/os"
	gcstr "github.com/PlayerR9/go-commons/strings"
	dbg "github.com/PlayerR9/go-debug/assert"
	"github.com/PlayerR9/grammar/ast"
	displ "github.com/PlayerR9/grammar/displayer"
	gr "github.com/PlayerR9/grammar/grammar"
)

// DecisionFunc is the function that returns the decision of the parser.
//
// Parameters:
//   - parser: The parser.
//   - lookahead: The lookahead token.
//
// Returns:
//   - Actioner: The action of the decision.
//   - error: An error if the decision is invalid.
type DecisionFunc[S gr.TokenTyper] func(parser *Parser[S], lookahead *gr.Token[S]) (Actioner, error)

// Parser is the parser of the grammar.
type Parser[S gr.TokenTyper] struct {
	// tokens is the tokens of the parser.
	tokens []*gr.Token[S]

	// stack is the stack of the parser.
	stack []*gr.Token[S]

	// popped is the stack of the parser.
	popped []*gr.Token[S]

	// decision is the function that returns the decision of the parser.
	decision DecisionFunc[S]

	// Err is the error reason of the parser.
	Err *displ.ErrParsing

	// last_action is the last action of the parser.
	last_action Actioner
}

// NewParser creates a new parser.
//
// Parameters:
//   - decision_func: The function that returns the decision of the parser.
//
// Returns:
//   - *Parser: The new parser.
//
// This function returns nil if the decision_func is nil.
func NewParser[S gr.TokenTyper](decision_func DecisionFunc[S]) *Parser[S] {
	if decision_func == nil {
		return nil
	}

	return &Parser[S]{
		decision: decision_func,
	}
}

// SetInputStream sets the input stream of the parser.
//
// Parameters:
//   - tokens: The input stream of the parser.
func (p *Parser[S]) SetInputStream(tokens []*gr.Token[S]) {
	// Clean the previous parser's state.
	gr.CleanTokens(p.tokens)
	p.tokens = p.tokens[:0]

	gr.CleanTokens(p.stack)
	p.stack = p.stack[:0]

	gr.CleanTokens(p.popped)
	p.popped = p.popped[:0]

	p.Err = nil

	// Set the new input stream.

	p.tokens = tokens
}

// Pop pops a token from the stack.
//
// Returns:
//   - *Token[T]: The token if the stack is not empty, nil otherwise.
//   - bool: True if the stack is not empty, false otherwise.
func (p *Parser[S]) Pop() (*gr.Token[S], bool) {
	if len(p.stack) == 0 {
		return nil, false
	}

	top := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]

	p.popped = append(p.popped, top)

	return top, true
}

// Peek pops a token from the stack without removing it.
//
// Returns:
//   - *Token[T]: The token if the stack is not empty, nil otherwise.
//   - bool: True if the stack is not empty, false otherwise.
func (p *Parser[S]) Peek() (*gr.Token[S], bool) {
	if len(p.stack) == 0 {
		return nil, false
	}

	return p.stack[len(p.stack)-1], true
}

// Shift shifts a token from the input stream to the stack.
//
// Returns:
//   - bool: True if the input stream is not empty, false otherwise.
func (p *Parser[S]) Shift() bool {
	if len(p.tokens) == 0 {
		return false
	}

	first := p.tokens[0]
	p.tokens = p.tokens[1:]

	p.stack = append(p.stack, first)

	return true
}

// GetPopped returns the popped tokens.
//
// Returns:
//   - []*Token[S]: The popped tokens.
func (p *Parser[S]) GetPopped() []*gr.Token[S] {
	popped := make([]*gr.Token[S], 0, len(p.popped))

	for i := 0; i < len(p.popped); i++ {
		popped = append(popped, p.popped[i])
	}

	slices.Reverse(popped)

	return popped
}

// Push pushes a token to the stack. Does nothing if the token is nil.
//
// Parameters:
//   - token: The token to push.
func (p *Parser[S]) Push(token *gr.Token[S]) {
	if token == nil {
		return
	}

	p.stack = append(p.stack, token)
}

// Refuse refuses all the tokens that were popped since the last
// call to Accept().
func (p *Parser[S]) Refuse() {
	for len(p.popped) > 0 {
		top := p.popped[len(p.popped)-1]
		p.popped = p.popped[:len(p.popped)-1]

		p.stack = append(p.stack, top)
	}
}

// Accept accepts all the tokens that were popped since the last
// call to Accept().
func (p *Parser[S]) Accept() {
	p.popped = p.popped[:0]
}

// get_forest returns the syntax forest of the parser.
//
// Parameters:
//   - parser: The parser.
//
// Returns:
//   - []*grammar.Token[S]: The syntax forest of the parser.
func get_forest[S gr.TokenTyper](parser *Parser[S]) []*gr.Token[S] {
	dbg.AssertNotNil(parser, "parser")

	var forest []*gr.Token[S]

	for {
		top, ok := parser.Pop()
		if !ok {
			break
		}

		dbg.AssertNotNil(top, "top")

		forest = append(forest, top)
	}

	return forest
}

// apply_reduce applies a reduce action to the parser.
//
// Parameters:
//   - parser: The parser.
//   - rule: The rule to reduce.
//
// Returns:
//   - error: An error if the parser encounters an error while applying the reduce action.
func apply_reduce[S gr.TokenTyper](parser *Parser[S], rule *Rule[S]) error {
	if parser == nil {
		panic("parser cannot be nil")
	} else if rule == nil {
		panic("rule cannot be nil")
	}

	var prev *S

	for _, rhs := range rule.GetRhss() {
		top, ok := parser.Pop()
		if !ok {
			return NewErrUnexpectedToken(prev, nil, rhs)
		}

		top_type := top.GetType()

		if top_type != rhs {
			return NewErrUnexpectedToken(prev, &top_type, rhs)
		}
	}

	popped := parser.GetPopped()
	last_token := popped[len(popped)-1]

	parser.Accept()

	tk := gr.NewToken(rule.lhs, "", popped[0].At, last_token.Lookahead)
	tk.AddChildren(popped)

	parser.Push(tk)

	return nil
}

// FullParse is just a wrapper around the Grammar.FullParse function.
//
// Parameters:
//   - tokens: The input stream of the parser.
//
// Returns:
//   - []*gr.Token[S]: The syntax forest of the input stream.
func (p *Parser[S]) FullParse(tokens []*gr.Token[S]) []*gr.Token[S] {
	p.SetInputStream(tokens)

	ok := p.Shift() // initial shift
	if !ok {
		forest := get_forest(p)

		p.Err = displ.NewErrParsing(0, -1, errors.New("no tokens were specified"))

		return forest
	}

	for p.Err == nil {
		top, _ := p.Peek()
		// luc.AssertOk(ok, "parser.Peek()")

		act, err := p.decision(p, top.Lookahead)
		if err != nil {
			p.Err = displ.NewErrParsing(top.At, -1, err)
			p.Refuse()
			break
		}

		p.Refuse()

		switch act := act.(type) {
		case *ShiftAction:
			_ = p.Shift()
			// luc.AssertOk(ok, "parser.Shift()")
		case *ReduceAction[S]:
			err := apply_reduce(p, act.rule)
			if err != nil {
				p.Err = displ.NewErrParsing(top.At, -1, err)
			}
		case *AcceptAction[S]:
			err := apply_reduce(p, act.rule)
			if err == nil {
				forest := get_forest(p)

				return forest
			}

			p.Err = displ.NewErrParsing(top.At, -1, err)
		default:
			p.Err = displ.NewErrParsing(top.At, -1, errors.New("invalid action type"))
		}
	}

	p.Refuse()
	forest := get_forest(p)

	return forest
}

// FullParseWithSteps is like FullParse but, for each step, it pauses and prints
// its debug state.
//
// Parameters:
//   - tokens: The input stream of the parser.
//
// Returns:
//   - []*gr.Token[S]: The syntax forest of the input stream.
func (p *Parser[S]) FullParseWithSteps(tokens []*gr.Token[S], data []byte, tab_size int) []*gr.Token[S] {
	p.SetInputStream(tokens)

	err := p.Step("\t\t**Initial State:**\n", data, tab_size)
	dbg.AssertErr(err, "parser.Step()")

	ok := p.Shift() // initial shift
	if !ok {
		forest := get_forest(p)

		p.Err = displ.NewErrParsing(0, -1, errors.New("no tokens were specified"))

		return forest
	}

	p.last_action = NewShiftAction()

	err = p.Step("\t\t**Initial Shift:**\n", data, tab_size)
	dbg.AssertErr(err, "parser.Step()")

	for p.Err == nil {
		top, _ := p.Peek()
		// luc.AssertOk(ok, "parser.Peek()")

		act, err := p.decision(p, top.Lookahead)
		if err != nil {
			p.Err = displ.NewErrParsing(top.At, -1, err)
			p.Refuse()
			break
		}

		p.Refuse()

		p.last_action = act

		err = p.Step("\t\t**Decision:**\n", data, tab_size)
		dbg.AssertErr(err, "parser.Step()")

		switch act := act.(type) {
		case *ShiftAction:
			_ = p.Shift()
			// luc.AssertOk(ok, "parser.Shift()")
		case *ReduceAction[S]:
			err := apply_reduce(p, act.rule)
			if err != nil {
				p.Err = displ.NewErrParsing(top.At, -1, err)
			}
		case *AcceptAction[S]:
			err := apply_reduce(p, act.rule)
			if err == nil {
				forest := get_forest(p)

				return forest
			}

			p.Err = displ.NewErrParsing(top.At, -1, err)
		default:
			p.Err = displ.NewErrParsing(top.At, -1, errors.New("invalid action type"))
		}

		p.last_action = nil

		err = p.Step("\t\t**Apply Action:**\n", data, tab_size)
		dbg.AssertErr(err, "parser.Step()")
	}

	p.Refuse()
	forest := get_forest(p)

	err = p.Step("\t\t**Final State:**\n", data, tab_size)
	dbg.AssertErr(err, "parser.Step()")

	return forest
}

// display_stack is a helper function that displays the stack.
func (p Parser[S]) display_stack() {
	var pr ast.AstPrinter[*gr.Token[S]]

	for _, elem := range p.stack {
		err := ast.Apply(&pr, elem)
		dbg.AssertErr(err, "traversing.Apply(&printer, %s)", elem.String())

		fmt.Println(pr.String())
		fmt.Println()
	}
}

// display_tokens is a helper function that displays the tokens.
func (p Parser[S]) display_tokens(width int) {
	elems := make([]string, 0, len(p.tokens)+1)
	elems = append(elems, "")

	for _, tok := range p.tokens {
		elems = append(elems, tok.String())
	}

	var str string

	if width <= 0 {
		str = strings.Join(elems, " <- ")
	} else {
		tmp, n := gcstr.AdaptToScreenWidth(elems, width, " <- ")
		str = tmp

		if n != 0 {
			str += fmt.Sprintf("\n+ %d more", n)
		}
	}

	fmt.Println(str)
}

// display_data is a helper function that displays the data.
//
// It is useful for debugging.
//
// Parameters:
//   - data: The data to display.
//   - tab_size: The size of the tab.
func (p Parser[S]) display_data(data []byte, tab_size int) {
	var at int

	if len(p.tokens) < 2 {
		at = 0
	} else {
		first_token := p.tokens[0]
		at = first_token.At
	}

	res := displ.PrintBoxedData(data, at,
		displ.WithDelta(1),
		displ.WithLimitNextLines(1),
		displ.WithLimitPrevLines(1),
		displ.WithFixedTabSize(tab_size),
	)

	fmt.Println(string(res))
}

// Step is a function that pauses and prints the current state of the parser.
//
// It is useful for debugging.
//
// Parameters:
//   - title: The title of the step. This is used for displaying the step title.
//
// Returns:
//   - error: Any error that might have occurred. This is used for fatal errors.
func (p *Parser[S]) Step(title string, data []byte, tab_size int) error {
	gcos.ClearScreen()

	p.display_data(data, tab_size)
	fmt.Println()

	if title != "" {
		fmt.Println()
		fmt.Println(title)
		fmt.Println()
	}

	if p.last_action == nil {
		p.display_tokens(3 * 80)
		p.display_stack()

		fmt.Println()

		fmt.Println("Press ENTER to continue...")
		fmt.Scanln()

		return nil
	}

	switch act := p.last_action.(type) {
	case *ShiftAction:
		if len(p.tokens) == 0 {
			return errors.New("no tokens left to shift")
		}

		first := p.tokens[0]

		fmt.Printf("Shifting: %s...\n", first.String())
	case *ReduceAction[S]:
		if act.rule == nil {
			return errors.New("no rule to reduce")
		}

		fmt.Printf("Reducing: %q...\n", act.rule.String())
	case *AcceptAction[S]:
		if act.rule == nil {
			return errors.New("no rule to accept")
		}

		fmt.Printf("Accepting: %q...\n", act.rule.String())
	default:
		return fmt.Errorf("invalid action type: %T", act)
	}

	fmt.Println()

	p.display_tokens(3 * 80)
	p.display_stack()
	fmt.Println()

	fmt.Println("Press ENTER to continue...")
	fmt.Scanln()

	return nil
}
