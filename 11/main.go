package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	fmt.Println(FindAnagrams([]string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}))
}

func FindAnagrams(words []string) map[string][]string {
	result := make(map[string][]string)
	temp := make(map[string][]string)

	for _, w := range words {
		w = strings.ToLower(w)
		runes := []rune(w)
		sort.Slice(runes, func(i, j int) bool { return runes[i] < runes[j] })
		key := string(runes)
		temp[key] = append(temp[key], w)
	}

	for _, group := range temp {
		if len(group) > 1 {
			sort.Strings(group)
			result[group[0]] = group
		}
	}

	return result
}
