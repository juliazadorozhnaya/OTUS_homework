package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	if len(s) == 0 {
		return "", nil
	}

	var result strings.Builder
	var prevRune rune
	var isPrevRuneDigit bool

	for i, r := range s {
		isDigit := unicode.IsDigit(r)

		if i == 0 && isDigit {
			return "", ErrInvalidString
		}

		if isDigit && isPrevRuneDigit {
			return "", ErrInvalidString
		}

		if isDigit {
			count, _ := strconv.Atoi(string(r))
			if count > 0 && prevRune != 0 {
				result.WriteString(strings.Repeat(string(prevRune), count))
			}
		} else {
			if !isPrevRuneDigit && prevRune != 0 {
				result.WriteRune(prevRune)
			}
			prevRune = r
		}

		isPrevRuneDigit = isDigit
	}

	if !isPrevRuneDigit && prevRune != 0 {
		result.WriteRune(prevRune)
	}

	return result.String(), nil
}
