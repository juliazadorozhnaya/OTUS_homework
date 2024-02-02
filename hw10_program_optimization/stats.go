package hw10programoptimization

import (
	"fmt"
	"io"
	"strings"

	"github.com/mailru/easyjson"
)

//easyjson:json
type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	domainSuffix := "." + domain

	for {
		var user User
		if err := easyjson.UnmarshalFromReader(r, &user); err != nil {
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

// getEmailDomain extracts the domain part from the email and converts it to lowercase.
func getEmailDomain(email string) string {
	atIndex := strings.LastIndex(email, "@")
	if atIndex == -1 {
		return ""
	}
	return strings.ToLower(email[atIndex+1:])
}
