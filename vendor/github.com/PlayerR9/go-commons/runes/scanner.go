package runes

import (
	"errors"
	"io"
	"unicode/utf8"
)

// Unread is the error returned when an attempt is made to unread a rune that was not read.
//
// Readers must return Unread itself and not wrap it as callers will test this error using ==.
var Unread error

func init() {
	Unread = errors.New("unread rune")
}

// CharStream is a linear stream of runes.
type CharStream struct {
	// data is the content of the stream.
	data []byte

	// pos is the current position in the stream.
	pos int

	// get_prev indicates whether the previous rune was unread.
	get_prev bool

	// prev is the previous rune.
	prev *rune

	// prev_size is the size of the previous rune.
	prev_size int
}

// ReadRune implements the io.RuneReader interface.
//
// Errors:
//   - io.EOF: When the stream is exhausted.
//   - *ErrInvalidUTF8Encoding: When the stream has an invalid UTF-8 encoding.
//
// Do err == io.EOF to check if the stream is exhausted. As in Go specification, do not wrap this io.EOF error
// if you want to propagate it as callers should also be able to do err == io.EOF to check the error.
func (s *CharStream) ReadRune() (rune, int, error) {
	if s.get_prev {
		if s.prev == nil {
			panic(Unread.Error())
		}

		s.get_prev = false

		return *s.prev, s.prev_size, nil
	}

	if len(s.data) == 0 {
		return '\000', 0, io.EOF
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

// UnreadRune implements the io.RuneUnreader interface.
//
// Errors:
//   - Unread: When no previous rune was read.
func (s *CharStream) UnreadRune() error {
	if s.prev == nil {
		return Unread
	}

	s.get_prev = true

	return nil
}

// Init initializes the CharStream.
//
// Parameters:
//   - data: The content of the stream.
func (s *CharStream) Init(data []byte) {
	s.data = data
	s.get_prev = false
	s.prev = nil
	s.pos = 0
}

// Pos returns the current position in the stream.
//
// Returns:
//   - int: The current position in the stream.
func (s CharStream) Pos() int {
	if s.get_prev {
		return s.pos - s.prev_size
	}

	return s.pos
}

// IsExhausted checks if the stream is exhausted.
//
// Returns:
//   - bool: True if the stream is exhausted, false otherwise.
func (s CharStream) IsExhausted() bool {
	if s.get_prev {
		return false
	}

	return len(s.data) == 0
}

// Copy creates a copy of the stream.
//
// Returns:
//   - CharStream: A copy of the stream.
func (s CharStream) Copy() CharStream {
	return CharStream{
		data:      s.data,
		pos:       s.pos,
		get_prev:  s.get_prev,
		prev:      s.prev,
		prev_size: s.prev_size,
	}
}
