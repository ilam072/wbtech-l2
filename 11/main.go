package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	fmt.Println(FindAnagrams([]string{"пятка", "пятак", "тяпка", "листок", "слиток", "столик", "стол"}))
}

func FindAnagrams(words []string) map[string][]string {
	result := make(map[string][]string)
	groups := make(map[string][]string)
	firstSeen := make(map[string]string)

	for _, w := range words {
		w = strings.ToLower(w)

		runes := []rune(w)
		sort.Slice(runes, func(i, j int) bool { return runes[i] < runes[j] })
		key := string(runes)

		if _, ok := firstSeen[key]; !ok {
			firstSeen[key] = w
		}

		groups[key] = append(groups[key], w)
	}

	for key, group := range groups {
		if len(group) > 1 {
			sort.Strings(group)
			result[firstSeen[key]] = group
		}
	}

	return result
}
