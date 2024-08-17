package slices

// Equaler is an interface that defines a method to compare two objects of the
// same type for equality.
type Equaler interface {
	// Equals returns true if the object is equal to the other object.
	//
	// Parameters:
	//   - other: The other object to compare to.
	//
	// Returns:
	//   - bool: True if the objects are equal, false otherwise.
	Equals(other Equaler) bool
}
