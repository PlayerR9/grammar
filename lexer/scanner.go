package lexer

import (
	"errors"
	"unicode/utf8"
)

type CharStream struct {
	data []byte
	pos  int

	get_prev bool

	prev      *rune
	prev_size int
}

func (s *CharStream) ReadRune() (rune, int, error) {
	if s.get_prev {
		if s.prev == nil {
			panic("no previous rune was read")
		}

		s.get_prev = false

		return *s.prev, s.prev_size, nil
	}

	if len(s.data) == 0 {
		return '\000', 0, NewErrInputStreamExhausted()
	}

	c, size := utf8.DecodeRune(s.data)
	if c == utf8.RuneError {
		return '\000', 0, NewErrInvalidUTF8Encoding(s.pos)
	}

	s.data = s.data[size:]
	s.prev = &c
	s.prev_size = size
	s.pos += size

	return c, size, nil
}

func (s *CharStream) UnreadRune() error {
	if s.prev == nil {
		return errors.New("no previous rune was read")
	}

	s.get_prev = true

	return nil
}

func NewCharStream() *CharStream {
	return &CharStream{}
}

func (cs *CharStream) Init(data []byte) {
	cs.data = data
	cs.get_prev = false
	cs.prev = nil
	cs.pos = 0
}

func (cs *CharStream) Pos() int {
	if cs.get_prev {
		return cs.pos - cs.prev_size
	}

	return cs.pos
}

func (cs *CharStream) IsExhausted() bool {
	if cs.get_prev {
		return false
	}

	return len(cs.data) == 0
}
