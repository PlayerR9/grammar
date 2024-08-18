package matcher

import (
	"io"

	gcch "github.com/PlayerR9/go-commons/runes"
	grmch "github.com/PlayerR9/grammar/machine"
)

var (
	WordMatcher *grmch.MachineState[*WordMatcherInfo]
)

func init() {
	fn := func(info *WordMatcherInfo, _ *rune) (grmch.SystemState, error) {

	}

	WordMatcher = grmch.NewMachineState(fn)

	WordMatcher.AddState()
}

type WordMatcherInfo struct {
	chars []rune
	pos   int
}

type WordMatcherMS struct {
	chars []rune
}

func NewWordMatcherMS(word string) (grmch.RunFunc[*WordMatcherInfo], func()) {
	chars, err := gcch.StringToUtf8(word)
	if err != nil {
		run_fn := func(scanner io.RuneScanner) (*WordMatcherInfo, error) {
			return nil, err
		}

		cleanup_fn := func() {}

		return run_fn, cleanup_fn
	}

	inf := &WordMatcherInfo{
		chars: chars,
		pos:   0,
	}

	run, cleanup := WordMatcher.Make(inf)
	return run, cleanup
}
