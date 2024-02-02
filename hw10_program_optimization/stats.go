package hw10programoptimization

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	decoder := json.NewDecoder(r)
	result := make(DomainStat)
	domainSuffix := "." + domain

	for {
		var user User
		if err := decoder.Decode(&user); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error decoding user: %w", err)
		}

		if emailDomain := getEmailDomain(user.Email); strings.HasSuffix(emailDomain, domainSuffix) {
			result[emailDomain]++
		}
	}

	return result, nil
}

// getEmailDomain извлекает доменную часть из email и приводит её к нижнему регистру
func getEmailDomain(email string) string {
	atIndex := strings.LastIndex(email, "@")
	if atIndex == -1 {
		return ""
	}
	return strings.ToLower(email[atIndex+1:])
}
