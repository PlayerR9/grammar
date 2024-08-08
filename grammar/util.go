package grammar

// CleanSlice cleans all the elements in a slice of type T that implements the Cleanup method
// and that returns a slice of type *T. This is because the returned elements are also cleaned.
//
// Parameters:
//   - s: The slice to clean.
//
// Notes: Remember to do s = s[:0] after having called this function.
func CleanTokens[S TokenTyper](s []*Token[S]) {
	if len(s) == 0 {
		return
	}

	var stack []*Token[S]

	for i := len(s) - 1; i >= 0; i-- {
		elem := s[i]
		if elem == nil {
			continue
		}

		tmp := (*elem).Cleanup()

		if len(tmp) > 0 {
			stack = append(stack, tmp...)
		}

		s[i] = nil
	}

	for len(stack) > 0 {
		top := stack[0]
		stack = stack[1:]

		if top == nil {
			continue
		}

		tmp := (*top).Cleanup()

		if len(tmp) > 0 {
			stack = append(stack, tmp...)
		}
	}
}
