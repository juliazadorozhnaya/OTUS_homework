package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/mailru/easyjson"
)

const bufferSize = 512

//easyjson:json
type UserEmail struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := DomainStat{}
	domainSuffix := "." + domain

	scanner := bufio.NewScanner(r)
	buffer := make([]byte, bufferSize)
	scanner.Buffer(buffer, 4*bufferSize)
	userEmail := &UserEmail{}
	for scanner.Scan() {
		*userEmail = UserEmail{} // Сброс структуры перед каждым использованием
		// Используйте easyjson для десериализации из сканера
		if err := easyjson.Unmarshal(scanner.Bytes(), userEmail); err != nil {
			return nil, err
		}

		email := userEmail.Email
		if strings.HasSuffix(email, domainSuffix) {
			if idx := strings.IndexRune(email, '@'); idx >= 0 {
				result[strings.ToLower(email[idx+1:])]++
			} else {
				return nil, fmt.Errorf("user email: %s format is not valid (doesn't contain \"@\" symbol)", userEmail.Email)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error while reading input: %w", err)
	}
	return result, nil
}
