package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	sr := []rune(s)
	var newStr strings.Builder

	for i, item := range sr {

		if unicode.IsDigit(item) && i == 0 {
			return "", ErrInvalidString
		}

		if i+1 < len(sr) && unicode.IsDigit(sr[i+1]) {
			count, _ := strconv.Atoi(string(sr[i+1]))
			newStr.WriteString(strings.Repeat(string(item), count))
			i++
		}
	}

	return newStr.String(), nil
}
