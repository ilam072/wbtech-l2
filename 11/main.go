package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	fmt.Println(FindAnagrams([]string{
		"пятка", "пятак", "тяпка",
		"листок", "слиток", "столик",
		"стол",
	}))
}

type anagramGroup struct {
	first string
	words []string
}

func FindAnagrams(words []string) map[string][]string {
	groups := make(map[string]*anagramGroup)

	for _, w := range words {
		w = strings.ToLower(w)

		runes := []rune(w)
		sort.Slice(runes, func(i, j int) bool { return runes[i] < runes[j] })
		key := string(runes)

		if _, ok := groups[key]; !ok {
			groups[key] = &anagramGroup{
				first: w,
				words: []string{w},
			}
			continue
		}

		groups[key].words = append(groups[key].words, w)
	}

	result := make(map[string][]string)

	for _, group := range groups {
		if len(group.words) > 1 {
			sort.Strings(group.words)
			result[group.first] = group.words
		}
	}

	return result
}
