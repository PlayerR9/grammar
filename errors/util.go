package errors

import "bytes"

var (
	// Newline is the newline character in bytes.
	Newline []byte
)

func init() {
	Newline = []byte("\n")
}

// LimitReverseLines is a function that limits the lines of the data in reverse order.
//
// Parameters:
//   - data: The data to limit.
//   - limit: The limit of the lines.
//
// Returns:
//   - []byte: The limited data.
func LimitReverseLines(data []byte, limit int) []byte {
	if len(data) == 0 {
		return nil
	}

	lines := bytes.Split(data, Newline)

	if limit == -1 || limit > len(lines) {
		limit = len(lines)
	}

	start_idx := len(lines) - limit

	lines = lines[start_idx:]

	return bytes.Join(lines, Newline)
}

// LimitLines is a function that limits the lines of the data.
//
// Parameters:
//   - data: The data to limit.
//   - limit: The limit of the lines.
//
// Returns:
//   - []byte: The limited data.
func LimitLines(data []byte, limit int) []byte {
	if len(data) == 0 {
		return nil
	}

	lines := bytes.Split(data, Newline)

	if limit == -1 || limit > len(lines) {
		limit = len(lines)
	}

	lines = lines[:limit]

	return bytes.Join(lines, Newline)
}
