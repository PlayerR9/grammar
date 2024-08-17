package bytes

import (
	"errors"
	"io"
)

// Unread is the error returned when an attempt is made to unread a byte that was not read.
//
// Readers must return Unread itself and not wrap it as callers will test this error using ==.
var Unread error

func init() {
	Unread = errors.New("unread byte")
}

// ByteStream is a linear stream of bytes.
type ByteStream struct {
	// data is the content of the stream.
	data []byte

	// pos is the current position in the stream.
	pos int

	// get_prev indicates whether the previous byte was unread.
	get_prev bool

	// prev is the previous byte.
	prev *byte
}

// ReadByte implements the io.ByteReader interface.
//
// Errors:
//   - io.EOF: When the stream is exhausted.
//   - *ErrInvalidUTF8Encoding: When the stream has an invalid UTF-8 encoding.
//
// Do err == io.EOF to check if the stream is exhausted. As in Go specification, do not wrap this io.EOF error
// if you want to propagate it as callers should also be able to do err == io.EOF to check the error.
func (s *ByteStream) ReadByte() (byte, error) {
	if s.get_prev {
		if s.prev == nil {
			panic(Unread.Error())
		}

		s.get_prev = false

		return *s.prev, nil
	}

	if len(s.data) == 0 {
		return '\000', io.EOF
	}

	b := s.data[0]

	s.data = s.data[1:]
	s.prev = &b
	s.pos++

	return b, nil
}

// UnreadByte implements the io.ByteUnreader interface.
//
// Errors:
//   - Unread: When no previous byte was read.
func (s *ByteStream) UnreadByte() error {
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
func (s *ByteStream) Init(data []byte) {
	s.data = data
	s.get_prev = false
	s.prev = nil
	s.pos = 0
}

// Pos returns the current position in the stream.
//
// Returns:
//   - int: The current position in the stream.
func (s ByteStream) Pos() int {
	if s.get_prev {
		return s.pos - 1
	}

	return s.pos
}

// IsExhausted checks if the stream is exhausted.
//
// Returns:
//   - bool: True if the stream is exhausted, false otherwise.
func (s ByteStream) IsExhausted() bool {
	if s.get_prev {
		return false
	}

	return len(s.data) == 0
}
