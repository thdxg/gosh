package readline

import (
	"errors"
	"io"
	"os"
	"strings"
	"unicode"
)

const (
	del = '\n'
)

func ReadLine() ([]string, error) {
	args := make([]string, 0)

	var current rune

	for current != del {
		tok := new(strings.Builder)
		for {
			char := make([]byte, 1)
			_, err := os.Stdin.Read(char)
			if err != nil {
				if errors.Is(err, io.EOF) {
					err = nil
				}
				return nil, err
			}
			current = rune(char[0])

			if unicode.IsSpace(current) {
				break
			} else if unicode.IsPrint(current) {
				tok.WriteRune(current)
			}
		}
		args = append(args, tok.String())
	}

	return args, nil
}
