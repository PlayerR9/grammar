package lexing

import (
	"io"
	"testing"

	gcch "github.com/PlayerR9/go-commons/runes"
)

func TestRightLex(t *testing.T) {
	var scanner gcch.CharStream

	scanner.Init([]byte("\r\n\r\nt"))

	f := func(scanner io.RuneScanner) ([]rune, error) {
		c, _, err := scanner.ReadRune()
		if err != nil {
			return nil, err
		}

		if c != '\r' {
			_ = scanner.UnreadRune()

			return nil, NoMatch
		}

		c, _, err = scanner.ReadRune()
		if err != nil {
			return nil, err
		}

		if c != '\n' {
			_ = scanner.UnreadRune()

			return nil, nil
		}

		return []rune{'\r', '\n'}, nil
	}

	chars, err := RightLex(&scanner, f)
	if err != nil {
		t.Errorf("expected no error, got %s instead", err.Error())
	}

	if len(chars) != 4 {
		t.Errorf("expected 4 but got %d", len(chars))
	}

	if chars[0] != '\r' {
		t.Errorf("expected '\r' but got %c", chars[0])
	}

	if chars[1] != '\n' {
		t.Errorf("expected '\n' but got %c", chars[1])
	}

	if chars[2] != '\r' {
		t.Errorf("expected '\r' but got %c", chars[2])
	}

	if chars[3] != '\n' {
		t.Errorf("expected '\n' but got %c", chars[3])
	}

	c, _, err := scanner.ReadRune()
	if err != nil {
		t.Errorf("expected no error, got %s instead", err.Error())
	}

	if c != 't' {
		t.Errorf("expected 't' but got %c", c)
	}
}
