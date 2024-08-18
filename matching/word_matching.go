package matching

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	gcers "github.com/PlayerR9/go-commons/errors"
	gcch "github.com/PlayerR9/go-commons/runes"
	gcslc "github.com/PlayerR9/go-commons/slices"
	dbg "github.com/PlayerR9/go-debug/assert"
)

type WordMatcher struct {
	chars []rune
	size  int
}

func NewWordMatcher(word string) (*WordMatcher, error) {
	if word == "" {
		return nil, nil
	}

	chars, err := gcch.StringToUtf8(word)
	if err != nil {
		return nil, err
	}

	return &WordMatcher{
		chars: chars,
		size:  len(chars),
	}, nil
}

func (wm WordMatcher) Match() WordMatcherInfo {
	return WordMatcherInfo{
		global: &wm,
		pos:    0,
	}
}

func (wm WordMatcher) Size() int {
	return wm.size
}

func (wm WordMatcher) CharAt(idx int) (rune, bool) {
	if idx < 0 || idx >= wm.size {
		return 0, false
	}

	return wm.chars[idx], true
}

var (
	ErrDone error
)

func init() {
	ErrDone = errors.New("done")
}

type WordMatcherInfo struct {
	global *WordMatcher
	pos    int
	err    error
}

func (wmi *WordMatcherInfo) Step(char *rune) {
	if wmi.err != nil {
		return
	}

	if wmi.pos >= wmi.global.Size() {
		wmi.err = ErrDone
		return
	}

	c, ok := wmi.global.CharAt(wmi.pos)
	dbg.AssertOk(ok, "wmi.global.CharAt(%d)", wmi.pos)

	if char == nil {
		wmi.err = fmt.Errorf("expected %q, got nothing instead", strconv.QuoteRune(c))
		return
	} else if c != *char {
		wmi.err = fmt.Errorf("expected %q, got %q instead", strconv.QuoteRune(c), strconv.QuoteRune(*char))
		return
	}

	wmi.pos++
}

func (wmi *WordMatcherInfo) IsDone() bool {
	return wmi.err == ErrDone
}

func (wmi *WordMatcherInfo) Err() error {
	return wmi.err
}

type StateMachiner interface {
	Step(char *rune)
	IsDone() bool
	Err() error
}

type Parallel struct {
	matchers []StateMachiner
}

func (p Parallel) Match() ParallelInfo {
	active := make([]int, 0, len(p.matchers))

	for i := 0; i < len(p.matchers); i++ {
		active = append(active, i)
	}

	return ParallelInfo{
		global: &p,
		active: active,
	}
}

func (p Parallel) get_matcher_at(idx int) (StateMachiner, bool) {
	if idx < 0 || idx >= len(p.matchers) {
		return nil, false
	}

	return p.matchers[idx], true
}

type ParallelInfo struct {
	global *Parallel
	active []int
	eval   gcers.ErrOrSol[int]
	err    error
	level  int
}

func (pi *ParallelInfo) Step(char *rune) {
	if pi.err != nil {
		return
	}

	if len(pi.active) == 0 {
		pi.err = ErrDone
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(pi.active))

	for _, idx := range pi.active {
		fn := func(idx int) {
			defer wg.Done()

			matcher, ok := pi.global.get_matcher_at(idx)
			if ok {
				matcher.Step(char)
			}
		}

		go fn(idx)
	}

	wg.Wait()

	level := pi.level

	fn := func(idx int) bool {
		matcher := pi.global.matchers[idx]

		if !matcher.IsDone() {
			return true
		}

		err := matcher.Err()
		if err != nil {
			pi.eval.AddErr(err, level)
		} else {
			pi.eval.AddSol(idx, level)
		}

		return false
	}

	pi.active = gcslc.SliceFilter(pi.active, fn)

	pi.level++

	if len(pi.active) == 0 {
		pi.err = ErrDone
	}
}

func MaxParallel(elems []func(), max int) {
	if len(elems) == 0 {
		return
	} else if max <= 0 {
		panic("invalid max")
	}

	elem_size := len(elems)

	if elem_size <= max {
		var wg sync.WaitGroup

		wg.Add(elem_size)

		for _, e := range elems {
			go func() {
				defer wg.Done()

				e()
			}()
		}

		wg.Wait()
	} else {
		groups := make(map[int][]func())

		for i := 0; i < max; i++ {
			groups[i] = make([]func(), 0, max)
		}

		for i := 0; i < len(elems); i++ {
			group_id := i % max

			groups[group_id] = append(groups[group_id], elems[i])
		}

		var wg sync.WaitGroup

		wg.Add(max)

		for i := 0; i < max; i++ {
			go func() {
				defer wg.Done()

				elems[i]()
			}()
		}

		wg.Wait()
	}
}

func SplitIntoGroups[T any](elems []T, count int) (map[int][]T, error) {
	if len(elems) == 0 {
		return nil, nil
	} else if count <= 0 {
		return nil, gcers.NewErrInvalidParameter("count", gcers.NewErrGT(0))
	}

	if count == 1 {
		groups := make(map[int][]T, 1)
		groups[0] = elems

		return groups, nil
	}

	if count < len(elems) {
		multiple := -1

		for i := count; i > 0 && multiple == -1; i-- {
			if len(elems)%i == 0 {
				multiple = i
			}
		}

		if multiple == -1 {

		} else {

		}
	}

	max := len(elems) / count

	groups := make(map[int][]T, count)

	for i := 0; i < count; i++ {
		groups[i] = make([]T, 0, max)
	}

	for i, e := range elems {
		group_id := i % count
		groups[group_id] = append(groups[group_id], e)
	}

	return groups
}
