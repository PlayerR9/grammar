package displayer

// PrintOptions are options that can be passed to the Print function.
type PrintOption func(s *PrintSettings)

// WithLimitPrevLines sets the limit of the previous lines to print.
// If the limit is negative, it is not set.
//
// Parameters:
//   - prev_lines: The limit of the previous lines to print.
//
// Returns:
//   - PrintOption: The function that sets the limit of the previous lines to print.
func WithLimitPrevLines(prev_lines int) PrintOption {
	if prev_lines < 0 {
		prev_lines = -1
	}

	return func(s *PrintSettings) {
		s.prev_lines = prev_lines
	}
}

// WithLimitNextLines sets the limit of the next lines to print.
// If the limit is negative, it is not set.
//
// Parameters:
//   - next_lines: The limit of the next lines to print.
//
// Returns:
//   - PrintOption: The function that sets the limit of the next lines to print.
func WithLimitNextLines(next_lines int) PrintOption {
	if next_lines < 0 {
		next_lines = -1
	}

	return func(s *PrintSettings) {
		s.next_lines = next_lines
	}
}

// WithDelta sets the delta to print.
// If the delta is negative, it is not set.
// If the delta is 0, it is set to 1.
//
// Parameters:
//   - delta: The delta to print.
//
// Returns:
//   - PrintOption: The function that sets the delta to print.
func WithDelta(delta int) PrintOption {
	if delta < 0 {
		delta = -1
	} else if delta == 0 {
		delta = 1
	}

	return func(s *PrintSettings) {
		s.delta = delta
	}
}

// WithFixedTabSize sets the fixed tab size to print.
// If the tab size is negative, it is not set.
// If the tab size is 0, it is set to 3.
//
// Parameters:
//   - tab_size: The fixed tab size to print.
//
// Returns:
//   - PrintOption: The function that sets the fixed tab size to print.
func WithFixedTabSize(tab_size int) PrintOption {
	if tab_size < 0 {
		tab_size = -1
	} else if tab_size == 0 {
		tab_size = 3
	}

	return func(s *PrintSettings) {
		s.tab_size = tab_size
	}
}
