package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type wordContainer struct {
	word      string
	frequency int
}

func Top10(rawStr string) []string {
	rawStr = strings.TrimSpace(rawStr)

	if len(rawStr) == 0 {
		return []string{}
	}

	rawWordsSlice := strings.Fields(rawStr)
	if len(rawWordsSlice) == 1 {
		return rawWordsSlice
	}
	wordsMap := make(map[string]int)

	for _, word := range rawWordsSlice {
		wordsMap[word]++
	}

	wordsContainer := make([]wordContainer, 0, len(wordsMap))

	for word, frequency := range wordsMap {
		wc := wordContainer{word: word, frequency: frequency}
		wordsContainer = append(wordsContainer, wc)
	}

	sort.Slice(wordsContainer, func(i, j int) bool {
		if wordsContainer[i].frequency == wordsContainer[j].frequency {
			return wordsContainer[i].word < wordsContainer[j].word
		}

		return wordsContainer[i].frequency > wordsContainer[j].frequency
	})

	if len(wordsContainer) > 10 {
		wordsContainer = wordsContainer[:10]
	}

	resultWordsSlice := make([]string, 0, len(wordsContainer))

	for _, wc := range wordsContainer {
		resultWordsSlice = append(resultWordsSlice, wc.word)
	}

	return resultWordsSlice
}
