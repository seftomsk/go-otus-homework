package hw10programoptimization

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	Email    string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain must not be empty")
	}

	var user User
	result := make(DomainStat)
	reader := bufio.NewReader(r)
	json := jsoniter.ConfigFastest
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return result, err
		}

		err = json.Unmarshal(line, &user)
		if err != nil {
			return result, err
		}

		if strings.HasSuffix(user.Email, domain) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
	return result, nil
}
