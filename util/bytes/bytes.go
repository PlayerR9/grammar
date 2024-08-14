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
