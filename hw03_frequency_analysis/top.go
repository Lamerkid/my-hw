package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(str string) []string {
	if len(str) == 0 {
		return nil
	}

	wordsSlice := strings.Fields(str)
	wordsMap := map[string]int{}

	for _, word := range wordsSlice {
		if val, ok := wordsMap[word]; ok {
			wordsMap[word] = val + 1
			continue
		}
		wordsMap[word] = 1
	}

	type record struct {
		Word  string
		Count int
	}

	records := make([]record, 0, len(wordsSlice))
	for k, v := range wordsMap {
		records = append(records, record{k, v})
	}

	sort.Slice(records, func(i, j int) bool {
		if records[i].Count == records[j].Count {
			return records[i].Word < records[j].Word
		}
		return records[i].Count > records[j].Count
	})

	var output []string
	for i := 0; i < 10; i++ {
		output = append(output, records[i].Word)
	}

	return output
}
