package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// ErrInvalidString is the error returned when the input string is invalid.
var ErrInvalidString = errors.New("invalid string")

// Unpack expands a string according to specific rules.
func Unpack(s string) (string, error) {
	var builder strings.Builder
	runeSlice := []rune(s)

	for i := 0; i < len(runeSlice); i++ {
		current := runeSlice[i]

		switch {
		case isLetterOrEscapedBackslash(current, i, runeSlice):
			builder.WriteRune(current)
		case current == '\\':
			next, err := handleBackslash(&i, runeSlice)
			if err != nil {
				return "", err
			}
			builder.WriteRune(next)
		case unicode.IsDigit(current):
			err := handleDigit(current, i, runeSlice, &builder)
			if err != nil {
				return "", err
			}
		default:
			return "", ErrInvalidString
		}
	}

	return builder.String(), nil
}

// isLetterOrEscapedBackslash checks if the current rune is a letter or an escaped backslash.
func isLetterOrEscapedBackslash(current rune, i int, runeSlice []rune) bool {
	return unicode.IsLetter(current) || (current == '\\' && i+1 < len(runeSlice) && runeSlice[i+1] == '\\')
}

// handleBackslash processes a backslash in the string.
func handleBackslash(i *int, runeSlice []rune) (rune, error) {
	*i++
	if *i >= len(runeSlice) {
		return 0, ErrInvalidString
	}
	next := runeSlice[*i]
	if unicode.IsDigit(next) || next == '\\' {
		return next, nil
	}
	return 0, ErrInvalidString
}

// handleDigit processes a digit in the string.
func handleDigit(current rune, i int, runeSlice []rune, builder *strings.Builder) error {
	if i == 0 || !(unicode.IsLetter(runeSlice[i-1]) || runeSlice[i-1] == '\\') {
		return ErrInvalidString
	}
	repeatCount, err := strconv.Atoi(string(current))
	if err != nil {
		return ErrInvalidString
	}
	prev := runeSlice[i-1]
	processRepeat(repeatCount, prev, builder)
	return nil
}

// processRepeat repeats the previous character based on the repeat count.
func processRepeat(repeatCount int, prev rune, builder *strings.Builder) {
	if repeatCount > 1 {
		builder.WriteString(strings.Repeat(string(prev), repeatCount-1))
	} else if repeatCount == 0 {
		removeLastRuneIfAny(builder)
	}
}

// removeLastRuneIfAny removes the last rune from the builder if any.
func removeLastRuneIfAny(builder *strings.Builder) {
	result := builder.String()
	if len(result) > 0 {
		builder.Reset()
		builder.WriteString(result[:len(result)-1])
	}
}
