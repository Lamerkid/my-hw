package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var outputString strings.Builder
	for i := 0; i < len(str); i++ {
		n, err := strconv.Atoi(string(str[i]))
		if i+1 < len(str) && string(str[i+1]) == "0" {
			continue
		}
		if err != nil {
			outputString.WriteString(string(str[i]))
		} else {
			if n == 0 {
				outputString.WriteString(strings.Repeat(string(str[i-1]), n))
			} else {
				outputString.WriteString(strings.Repeat(string(str[i-1]), n-1))
			}
		}
	}
	return outputString.String(), nil
}
