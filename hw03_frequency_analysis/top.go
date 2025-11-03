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
	wordsSlice := strings.Fields(str)

	// От начала или конца строки есть пунктуация
	trimPunctuation := regexp.MustCompile(`^\p{P}*|\p{P}*$`)

	var cleanedWords []string
	matchWord := regexp.MustCompile(`\p{P}`)
	for _, word := range wordsSlice {
		// Не записываем пунктуации длиной меньше чем 2
		match := matchWord.MatchString(word)
		if len(word) < 2 && match {
			continue
		}
		cleaned := strings.ToLower(trimPunctuation.ReplaceAllString(word, ""))
		// Если слово полностью стерлось, значит оно состоит из пунктуаций, записываем его в неизменном виде
		if len(cleaned) == 0 {
			cleanedWords = append(cleanedWords, word)
		} else {
			cleanedWords = append(cleanedWords, cleaned)
		}
	}

	wordsMap := map[string]int{}
	for _, word := range cleanedWords {
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

	records := make([]record, 0, len(cleanedWords))
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
