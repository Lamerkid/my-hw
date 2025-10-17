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
		if runesArray[i] == '\\' {
			// Если экранирование последнее в строке или экранируется не цифра или слэш, то возвращаем ошибку
			if i+1 >= len(runesArray) || !(runesArray[i+1] >= '0' && runesArray[i+1] <= '9') && runesArray[i+1] != '\\' {
				return "", ErrInvalidString
			}

			outputString.WriteRune(runesArray[i+1])
			i++
			continue
		}
		// Если успешное преобразование в int, то записываем предыдущую руну n раз
		if n, err := strconv.Atoi(string(runesArray[i])); err == nil {
			// Если цифра первая в строке или двузначное число, то возвращаем ошибку
			if i == 0 || i+1 < len(runesArray) && runesArray[i+1] >= '0' && runesArray[i+1] <= '9' {
				return "", ErrInvalidString
			}
			// Eсли в строке 0, то удалить последний записанный элемент
			if n == 0 {
				savedOutput := []rune(outputString.String())
				outputString.Reset()
				outputString.WriteString(string(savedOutput[:len(savedOutput)-1]))
			} else {
				outputString.WriteString(strings.Repeat(string(runesArray[i-1]), n-1))
			}
		} else {
			outputString.WriteRune(runesArray[i])
		}
	}
	return outputString.String(), nil
}
