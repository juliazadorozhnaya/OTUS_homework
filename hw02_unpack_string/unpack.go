package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var builder strings.Builder
	runeSlice := []rune(s)

	for i := 0; i < len(runeSlice); i++ {
		current := runeSlice[i]

		if unicode.IsLetter(current) || (current == '\\' && i+1 < len(runeSlice) && runeSlice[i+1] == '\\') {
			builder.WriteRune(current)
		} else if current == '\\' {
			i++
			if i >= len(runeSlice) {
				return "", ErrInvalidString
			}

			next := runeSlice[i]
			if unicode.IsDigit(next) || next == '\\' {
				builder.WriteRune(next)
			} else {
				return "", ErrInvalidString
			}
		} else if unicode.IsDigit(current) {
			if i == 0 || !(unicode.IsLetter(runeSlice[i-1]) || runeSlice[i-1] == '\\') {
				return "", ErrInvalidString
			}

			repeatCount, err := strconv.Atoi(string(current))
			if err != nil {
				return "", ErrInvalidString
			}

			prev := runeSlice[i-1]
			if repeatCount > 1 {
				builder.WriteString(strings.Repeat(string(prev), repeatCount-1))
			} else if repeatCount == 0 {
				result := builder.String()
				if len(result) > 0 {
					builder.Reset()
					builder.WriteString(result[:len(result)-1])
				}
			}
		} else {
			return "", ErrInvalidString
		}
	}

	return builder.String(), nil
}
