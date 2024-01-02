package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var result string
	var prevRune rune
	var skipNext bool
	var prevDigit bool

	if len(s) == 0 {
		return "", nil
	}

	for i, r := range s {
		if i == 0 && unicode.IsDigit(r) {
			return "", ErrInvalidString
		}

		if unicode.IsDigit(r) {
			if prevDigit {
				return "", ErrInvalidString
			}
			count, _ := strconv.Atoi(string(r))
			if count == 0 {
				skipNext = true
			} else {
				result += strings.Repeat(string(prevRune), count-1)
			}
			prevDigit = true
		} else {
			if !unicode.IsDigit(prevRune) && !skipNext && i != 0 {
				result += string(prevRune)
			}
			prevRune = r
			skipNext = false
			prevDigit = false
		}
	}

	if !unicode.IsDigit(prevRune) && !skipNext {
		result += string(prevRune)
	}

	return result, nil
}
