package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

func Top10(str string) []string {
	if len(str) == 0 {
		return nil
	}
	// Найти все слова (включая знаки препинания внутри слов)
	re := regexp.MustCompile(`[\wа-яА-Я.,!?;:-]+`)
	wordsSlice := re.FindAllString(str, -1)

	// Начало или конец строки НЕ слова
	nonWords := regexp.MustCompile(`^[^\wа-яА-Я-]*|[^\wа-яА-Я-]*$`)

	for i := 0; i < len(wordsSlice); i++ {
		// Удаление всего подходящего под nonWords и приведение строки в нижний регистр
		wordsSlice[i] = nonWords.ReplaceAllString(wordsSlice[i], "")
		wordsSlice[i] = strings.ToLower(wordsSlice[i])
	}

	wordsMap := map[string]int{}

	for _, word := range wordsSlice {
		// "-" не является словом - не записываем
		if word == "-" {
			continue
		}
		if count, ok := wordsMap[word]; ok {
			wordsMap[word] = count + 1
			continue
		}
		wordsMap[word] = 1
	}

	type record struct {
		Word  string
		Count int
	}

	records := make([]record, 0, len(wordsSlice))
	for w, c := range wordsMap {
		records = append(records, record{w, c})
	}
	// Сортируем по кол-ву, далее по алфавиту
	sort.Slice(records, func(i, j int) bool {
		if records[i].Count == records[j].Count {
			return records[i].Word < records[j].Word
		}
		return records[i].Count > records[j].Count
	})

	var output []string
	for i := 0; i < 10; i++ {
		if i == len(records) {
			break
		}
		output = append(output, records[i].Word)
	}

	return output
}
