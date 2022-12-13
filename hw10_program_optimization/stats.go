package hw10programoptimization

import (
	"encoding/json"
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
	user := &User{}
	result := make(DomainStat)
	decoder := json.NewDecoder(r)

	for decoder.More() {
		err := decoder.Decode(user)
		if err != nil {
			return nil, err
		}
		if strings.HasSuffix(user.Email, domain) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
	return result, nil
}
