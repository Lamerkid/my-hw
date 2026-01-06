//go:generate easyjson -no_std_marshalers stats.go
package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	easyjson "github.com/mailru/easyjson"
)

//easyjson:json
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
	scanner := bufio.NewScanner(r)
	result := make(DomainStat)

	user := User{}

	for scanner.Scan() {
		if err := easyjson.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, fmt.Errorf("unmarshaling error: %w", err)
		}
		if strings.Contains(user.Email, "."+domain) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
	return result, nil
}
