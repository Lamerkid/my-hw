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
		// Если успешное преобразование в int, то записываем предыдущую руну n раз
		if n, err := strconv.Atoi(string(runesArray[i])); err == nil {
			// Если n первая в строке или число двузначное, то возвращаем пустую строку и ошибку
			if i == 0 || i+1 < len(runesArray) && string(runesArray[i+1]) >= "0" && string(runesArray[i+1]) <= "9" {
				return "", ErrInvalidString
			}

			outputString.WriteString(strings.Repeat(string(runesArray[i-1]), n))

			// Если после руны нет числа, то записываем эту руну
		} else if !(i+1 < len(runesArray) && string(runesArray[i+1]) >= "0" && string(runesArray[i+1]) <= "9") {
			outputString.WriteString(string(runesArray[i]))
		}
	}
	return outputString.String(), nil
}
