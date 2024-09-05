package internal

// TokenTyper is an interface for token types.
type TokenTyper interface {
	~int

	// String returns the literal name of the token type.
	//
	// Returns:
	//   - string: The literal name of the token type.
	String() string

	// IsTerminal checks whether the token type is a terminal.
	//
	// Returns:
	//   - bool: True if the token type is a terminal, false otherwise.
	IsTerminal() bool
}
