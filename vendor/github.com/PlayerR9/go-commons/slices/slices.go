package slices

import (
	"cmp"

	"golang.org/x/exp/slices"
)

// Find returns the index of the first occurrence of an element in the slice.
//
// Parameters:
//   - S: slice of elements.
//   - elem: element to find.
//
// Returns:
//   - int: index of the first occurrence of the element or -1 if not found.
func Find[T comparable](S []T, elem T) int {
	if len(S) == 0 {
		return -1
	}

	for i, e := range S {
		if e == elem {
			return i
		}
	}

	return -1
}

// FindEquals is the same as Find but uses the Equals method of the elements.
//
// Parameters:
//   - S: slice of elements.
//   - elem: element to find.
//
// Returns:
//   - int: index of the first occurrence of the element or -1 if not found.
func FindEquals[T Equaler](S []T, elem T) int {
	if len(S) == 0 {
		return -1
	}

	for i, e := range S {
		ok := e.Equals(elem)
		if ok {
			return i
		}
	}

	return -1
}

// Uniquefy removes duplicate elements from the slice.
//
// Parameters:
//   - S: slice of elements.
//   - prioritizeFirst: If true, the first occurrence of an element is kept.
//     If false, the last occurrence of an element is kept.
//
// Returns:
//   - []T: slice of elements with duplicates removed.
//
// Behavior:
//   - The function preserves the order of the elements in the slice.
func Uniquefy[T comparable](S []T, prioritizeFirst bool) []T {
	if len(S) < 2 {
		return S
	}

	var unique []T

	if prioritizeFirst {
		seen := make(map[T]bool)

		for _, e := range S {
			_, ok := seen[e]
			if !ok {
				unique = append(unique, e)
				seen[e] = true
			}
		}
	} else {
		seen := make(map[T]int)

		for _, e := range S {
			pos, ok := seen[e]
			if !ok {
				seen[e] = len(unique)
				unique = append(unique, e)
			} else {
				unique[pos] = e
			}
		}
	}

	return unique
}

// UniquefyLeft is a helper function that removes duplicate elements from the slice.
//
// Parameters:
//   - S: slice of elements.
//
// Returns:
//   - []T: slice of elements with duplicates removed.
func uniquefyLeft[T Equaler](S []T) []T {
	if len(S) < 2 {
		return S
	}

	for i := 0; i < len(S)-1; i++ {
		elem := S[i]
		top := i + 1

		for j := i + 1; j < len(S); j++ {
			ok := elem.Equals(S[j])
			if !ok {
				S[top] = S[j]
				top++
			}
		}

		S = S[:top]
	}

	return S
}

// UniquefyEquals is the same as Uniquefy but uses the Equals method of the elements.
//
// Parameters:
//   - S: slice of elements.
//   - prioritizeFirst: If true, the first occurrence of an element is kept.
//     If false, the last occurrence of an element is kept.
//
// Returns:
//   - []T: slice of elements with duplicates removed.
//
// Behavior:
//   - The function preserves the order of the elements in the slice.
//   - This can modify the original slice.
func UniquefyEquals[T Equaler](S []T, prioritizeFirst bool) []T {
	if len(S) < 2 {
		return S
	}

	if !prioritizeFirst {
		slices.Reverse(S)

		S = uniquefyLeft(S)

		slices.Reverse(S)
	} else {
		S = uniquefyLeft(S)
	}

	return S
}

// MergeUnique merges two slices and removes duplicate elements.
//
// Parameters:
//   - S1: first slice of elements.
//   - S2: second slice of elements.
//
// Returns:
//   - []T: slice of elements with duplicates removed.
//
// Behaviors:
//   - The function does not preserve the order of the elements in the slices.
func MergeUnique[T comparable](S1, S2 []T) []T {
	seen := make(map[T]bool)

	for _, e := range S1 {
		_, ok := seen[e]
		if !ok {
			seen[e] = true
		}
	}

	for _, e := range S2 {
		_, ok := seen[e]
		if !ok {
			seen[e] = true
		}
	}

	merged := make([]T, 0, len(seen))

	for e := range seen {
		merged = append(merged, e)
	}

	return merged
}

// MergeUniqueEquals is the same as MergeUnique but uses the Equals method of the elements.
//
// Parameters:
//   - S1: first slice of elements.
//   - S2: second slice of elements.
//
// Returns:
//   - []T: slice of elements with duplicates removed.
//
// Behaviors:
//   - The function does preserve the order of the elements in the slices.
func MergeUniqueEquals[T Equaler](S1, S2 []T) []T {
	S1 = UniquefyEquals(S1, true)
	S2 = UniquefyEquals(S2, true)

	if len(S1) == 0 {
		return S2
	} else if len(S2) == 0 {
		return S1
	}

	elems := make([]T, len(S1))
	copy(elems, S1)
	limit := len(elems)

	for _, e := range S2 {
		found := false

		for i := 0; i < limit && !found; i++ {
			found = elems[i].Equals(e)
		}

		if !found {
			elems = append(elems, e)
		}
	}

	return elems
}

// IndexOfDuplicate returns the index of the first duplicate element in the slice.
//
// Parameters:
//   - S: slice of elements.
//
// Returns:
//   - int: index of the first duplicate element or -1 if there are no duplicates.
func IndexOfDuplicate[T comparable](S []T) int {
	if len(S) < 2 {
		return -1
	}

	seen := make(map[T]bool)

	for i, e := range S {
		if _, ok := seen[e]; ok {
			return i
		}

		seen[e] = true
	}

	return -1
}

// IndexOfDuplicateEquals is the same as IndexOfDuplicate but uses the Equals method of the elements.
//
// Parameters:
//   - S: slice of elements.
//
// Returns:
//   - int: index of the first duplicate element or -1 if there are no duplicates.
func IndexOfDuplicateEquals[T Equaler](S []T) int {
	if len(S) < 2 {
		return -1
	}

	for i := 0; i < len(S)-1; i++ {
		elem := S[i]

		for j := i + 1; j < len(S); j++ {
			ok := elem.Equals(S[j])
			if ok {
				return j
			}
		}
	}

	return -1
}

// computeLPSArray is a helper function that computes the Longest Prefix
// Suffix (LPS) array for the Knuth-Morris-Pratt algorithm.
//
// Parameters:
//   - subS: The subslice to compute the LPS array for.
//   - lps: The LPS array to store the results in.
//
// Behavior:
//   - The function modifies the lps array in place.
//   - The lps array is initialized with zeros.
//   - The lps array is used to store the length of the longest prefix
//     that is also a suffix for each index in the subslice.
//   - The first element of the lps array is always 0.
func computeLPSArray[T comparable](subS []T, lps []int) {
	length := 0
	i := 1
	lps[0] = 0 // lps[0] is always 0

	// the loop calculates lps[i] for i = 1 to len(subS)-1
	for i < len(subS) {
		if subS[i] == subS[length] {
			length++
			lps[i] = length
			i++
		} else {
			if length != 0 {
				length = lps[length-1]
			} else {
				lps[i] = 0
				i++
			}
		}
	}
}

// FindSubBytesFrom finds the first occurrence of a subslice in a byte
// slice starting from a given index.
//
// Parameters:
//   - S: The byte slice to search in.
//   - subS: The byte slice to search for.
//   - at: The index to start searching from.
//
// Returns:
//   - int: The index of the first occurrence of the subslice.
//
// Behavior:
//   - The function uses the Knuth-Morris-Pratt algorithm to find the subslice.
//   - If S or subS is empty, the function returns -1.
//   - If the subslice is not found, the function returns -1.
//   - If at is negative, it is set to 0.
func FindSubsliceFrom[T comparable](S []T, subS []T, at int) int {
	if len(subS) == 0 || len(S) == 0 || at+len(subS) > len(S) {
		return -1
	}

	if at < 0 {
		at = 0
	}

	lps := make([]int, len(subS))
	computeLPSArray(subS, lps)

	i := at
	j := 0
	for i < len(S) {
		if S[i] == subS[j] {
			i++
			j++
		}

		if j == len(subS) {
			return i - j
		} else if i < len(S) && S[i] != subS[j] {
			if j != 0 {
				j = lps[j-1]
			} else {
				i = i + 1
			}
		}
	}

	return -1
}

// computeLPSArrayEquals is a helper function that computes the Longest Prefix
// Suffix (LPS) array for the Knuth-Morris-Pratt algorithm using a custom
// comparison function.
//
// Parameters:
//   - subS: The subslice to compute the LPS array for.
//   - lps: The LPS array to store the results in.
//
// Behavior:
//   - The function modifies the lps array in place.
//   - The lps array is initialized with zeros.
//   - The lps array is used to store the length of the longest prefix
//     that is also a suffix for each index in the subslice.
//   - The first element of the lps array is always 0.
func computeLPSArrayEquals[T Equaler](subS []T, lps []int) {
	length := 0
	i := 1
	lps[0] = 0 // lps[0] is always 0

	// the loop calculates lps[i] for i = 1 to len(subS)-1
	for i < len(subS) {
		ok := subS[i].Equals(subS[length])
		if ok {
			length++
			lps[i] = length
			i++
		} else {
			if length != 0 {
				length = lps[length-1]
			} else {
				lps[i] = 0
				i++
			}
		}
	}
}

// FindSubsliceFromEquals finds the first occurrence of a subslice in a byte
// slice starting from a given index using a custom comparison function.
//
// Parameters:
//   - S: The byte slice to search in.
//   - subS: The byte slice to search for.
//   - at: The index to start searching from.
//
// Returns:
//   - int: The index of the first occurrence of the subslice.
//
// Behavior:
//   - The function uses the Knuth-Morris-Pratt algorithm to find the subslice.
//   - If S or subS is empty, the function returns -1.
//   - If the subslice is not found, the function returns -1.
//   - If at is negative, it is set to 0.
func FindSubsliceFromEquals[T Equaler](S []T, subS []T, at int) int {
	if len(subS) == 0 || len(S) == 0 || at+len(subS) > len(S) {
		return -1
	}

	if at < 0 {
		at = 0
	}

	lps := make([]int, len(subS))
	computeLPSArrayEquals(subS, lps)

	i := at
	j := 0
	for i < len(S) {
		ok := S[i].Equals(subS[j])
		if ok {
			i++
			j++
		}

		if j == len(subS) {
			return i - j
		} else if i < len(S) {
			ok := S[i].Equals(subS[j])
			if !ok {
				if j != 0 {
					j = lps[j-1]
				} else {
					i = i + 1
				}
			}
		}
	}

	return -1
}

/*
// Difference returns the elements that are in S1 but not in S2.
//
// Parameters:
//   - S1: The first slice of elements.
//   - S2: The second slice of elements.
func Difference[T comparable](S1, S2 []T) []T {
	if len(S1) == 0 {
		return S2
	} else if len(S2) == 0 {
		return S1
	}

	seen := make(map[T]bool)

	for _, e := range S2 {
		seen[e] = true
	}

	diff := make([]T, 0)

	for _, e := range S1 {
		if _, ok := seen[e]; !ok {
			diff = append(diff, e)
		}
	}

	return diff
}
*/

// TryInsert is a helper function that inserts an element into a slice only
// if the element is not already in the slice.
//
// Parameters:
//   - slc: The slice to insert into.
//   - e: The element to insert.
//
// Returns:
//   - []T: The slice with the inserted element.
//
// This function only works if the slice is sorted.
func TryInsert[T cmp.Ordered](slc []T, e T) []T {
	pos, ok := slices.BinarySearch(slc, e)
	if ok {
		return slc
	}

	slc = slices.Insert(slc, pos, e)

	return slc
}

// OrderedUniquefy returns a copy of elems without duplicates.
//
// Parameters:
//   - elems: The elements to uniquefy.
//
// Returns:
//   - []T: The unique elements.
//
// This function also sorts the elements.
func OrderedUniquefy[T cmp.Ordered](elems []T) []T {
	if len(elems) == 0 {
		return nil
	}

	var unique []T

	for _, elem := range elems {
		pos, ok := slices.BinarySearch(unique, elem)
		if !ok {
			unique = slices.Insert(unique, pos, elem)
		}
	}

	return unique
}

// FilterNonEmpty removes zero values from a slice of elements.
//
// Parameters:
//   - values: The slice of elements to trim.
//
// Returns:
//   - []T: The slice of elements with zero values removed.
//
// For example, in a slice of ints, all the values that are 0 will be removed.
func FilterNonEmpty[T cmp.Ordered](values []T) []T {
	if len(values) == 0 {
		return nil
	}

	empty := *new(T)

	var top int

	for i := 0; i < len(values); i++ {
		if values[i] != empty {
			values[top] = values[i]
			top++
		}
	}

	return values[:top:top]
}
