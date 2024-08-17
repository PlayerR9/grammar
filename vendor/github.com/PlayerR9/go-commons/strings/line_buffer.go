package strings

import "strings"

// LineBuffer is a struct that represents a line buffer.
type LineBuffer struct {
	// builder is the line builder.
	builder strings.Builder

	// lines is the line buffer.
	lines []string
}

// NewLineBuffer creates a new line buffer. This is unnecessary as it is equivalent
// to var lb LineBuffer
//
// Returns:
//   - LineBuffer: The new line buffer.
func NewLineBuffer() LineBuffer {
	return LineBuffer{
		lines: make([]string, 0),
	}
}

// String returns the lines in the line buffer as a string joined by newlines.
func (lb LineBuffer) String() string {
	if lb.builder.Len() > 0 {
		lb.lines = append(lb.lines, lb.builder.String())
		lb.builder.Reset()
	}

	return strings.Join(lb.lines, "\n")
}

// AddLine adds a line to the line buffer.
//
// Parameters:
//   - line: The line to add.
func (lb *LineBuffer) AddLine(line string) {
	if lb.builder.Len() > 0 {
		lb.lines = append(lb.lines, lb.builder.String())
		lb.builder.Reset()
	}

	lb.lines = append(lb.lines, line)
}

// AddString adds a string to the line buffer.
//
// Parameters:
//   - line: The string to add.
func (lb *LineBuffer) AddString(line string) {
	lb.builder.WriteString(line)
}

// Accept accepts the current line buffer.
func (lb *LineBuffer) Accept() {
	if lb.builder.Len() == 0 {
		return
	}

	lb.lines = append(lb.lines, lb.builder.String())
	lb.builder.Reset()
}

// Reset resets the line buffer.
func (lb *LineBuffer) Reset() {
	lb.lines = lb.lines[:0]
	lb.builder.Reset()
}
