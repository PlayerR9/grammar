package bytes

import "bytes"

// FixTabSize fixes the tab size by replacing it with a specified bytes iff
// the tab size is greater than 0. The replacement bytes are repeated for the
// specified number of times.
//
// Parameters:
//   - size: The size of the tab.
//   - rep: The replacement bytes.
//
// Returns:
//   - []byte: The fixed tab size.
func FixTabSize(size int, rep []byte) []byte {
	if size <= 0 {
		return []byte{'\t'}
	}

	return bytes.Repeat(rep, size)
}

// DetermineCoords is a helper function that determines the coordinates of the given position.
//
// Parameters:
//   - data: The data read from the input stream.
//   - pos: The position of the faulty token.
//
// Returns:
//   - int: The x coordinate of the faulty token.
//   - int: The y coordinate of the faulty token.
func DetermineCoords(data []byte, pos int) (int, int) {
	if len(data) == 0 {
		return 0, 0
	}

	if pos < 0 {
		pos = len(data) + pos
	} else if pos >= len(data) {
		pos = len(data) - 1
	}

	var x int
	var y int

	for i := 0; i < pos; i++ {
		if data[i] == '\n' {
			x = 0
			y++
		} else {
			x++
		}
	}

	return x, y
}
