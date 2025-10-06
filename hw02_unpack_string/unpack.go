package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var outputString strings.Builder
	runesArray := []rune(str)
	for i := 0; i < len(runesArray); i++ {

		if i+1 < len(runesArray) && string(runesArray[i+1]) == "0" {
			continue
		}

		if i+1 < len(runesArray) && string(runesArray[i]) >= "0" && string(runesArray[i]) <= "9" && string(runesArray[i+1]) >= "0" && string(runesArray[i+1]) <= "9" {
			return "", ErrInvalidString
		}

		n, err := strconv.Atoi(string(runesArray[i]))

		if err != nil {
			outputString.WriteString(string(runesArray[i]))
		} else {
			if i == 0 {
				return "", ErrInvalidString
			}

			if n == 0 {
				outputString.WriteString(strings.Repeat(string(runesArray[i-1]), n))
			} else {
				outputString.WriteString(strings.Repeat(string(runesArray[i-1]), n-1))
			}
		}
	}
	return outputString.String(), nil
}
