package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(text string) []string {
	words := strings.Fields(text)
	wordCount := make(map[string]int)

	// запишем все слова в мапу
	for _, word := range words {
		wordCount[word]++
	}

	// преобразуем в слайс для сортировки
	wordFreq := make([]struct {
		word  string
		count int
	}, 0, len(wordCount))

	// перезапишем в слайс
	for word, count := range wordCount {
		wordFreq = append(wordFreq, struct {
			word  string
			count int
		}{word, count})
	}

	// сортируем слайс, чтобы потом с него вытащить ТОП по count'ам
	sort.Slice(wordFreq, func(i, j int) bool {
		if wordFreq[i].count == wordFreq[j].count {
			return wordFreq[i].word < wordFreq[j].word
		} // если количество равно, то смотрим на лексикографический порядок самого слова
		return wordFreq[i].count > wordFreq[j].count
	})

	topWords := make([]string, 0, 10)
	for i := 0; i < 10 && i < len(wordFreq); i++ {
		topWords = append(topWords, wordFreq[i].word)
	}

	return topWords
}
