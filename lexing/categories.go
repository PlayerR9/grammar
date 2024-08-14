package lexing

import (
	"io"
	"unicode"
)

var (
	// CatDecimal is the category of decimal digits (i.e., from 0 to 9).
	CatDecimal LexFunc

	// CatUppercase is the category of uppercase letters (i.e., from A to Z).
	CatUppercase LexFunc

	// CatLowercase is the category of lowercase letters (i.e., from a to z).
	CatLowercase LexFunc
)

func init() {
	CatDecimal = func(scanner io.RuneScanner) ([]rune, error) {
		// [0-9]

		c, _, err := scanner.ReadRune()
		if err != nil {
			return nil, err
		}

		if !unicode.IsDigit(c) {
			_ = scanner.UnreadRune()

			return nil, nil
		}

		return []rune{c}, nil
	}

	CatUppercase = func(scanner io.RuneScanner) ([]rune, error) {
		// [A-Z]

		c, _, err := scanner.ReadRune()
		if err != nil {
			return nil, err
		}

		if !unicode.IsUpper(c) {
			_ = scanner.UnreadRune()

			return nil, nil
		}

		return []rune{c}, nil
	}

	CatLowercase = func(scanner io.RuneScanner) ([]rune, error) {
		// [a-z]

		c, _, err := scanner.ReadRune()
		if err != nil {
			return nil, err
		}

		if !unicode.IsLower(c) {
			_ = scanner.UnreadRune()

			return nil, nil
		}

		return []rune{c}, nil
	}
}
