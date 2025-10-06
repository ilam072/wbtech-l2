package cut

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseFields парсит строку вида "1,3-5" и возвращает уникальные индексы (0-based)
func ParseFields(input string) ([]int, error) {
	fields := strings.Split(input, ",")
	resultMap := make(map[int]struct{})

	for _, f := range fields {
		if strings.Contains(f, "-") {
			parts := strings.Split(f, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid range: %s", f)
			}
			start, err := strconv.Atoi(parts[0])
			if err != nil || start < 1 {
				return nil, fmt.Errorf("invalid start of range: %s", parts[0])
			}
			end, err := strconv.Atoi(parts[1])
			if err != nil || end < start {
				return nil, fmt.Errorf("invalid end of range: %s", parts[1])
			}
			for i := start - 1; i <= end-1; i++ {
				resultMap[i] = struct{}{}
			}
		} else {
			idx, err := strconv.Atoi(f)
			if err != nil || idx < 1 {
				return nil, fmt.Errorf("invalid field: %s", f)
			}
			resultMap[idx-1] = struct{}{}
		}
	}

	result := make([]int, 0, len(resultMap))
	for k := range resultMap {
		result = append(result, k)
	}

	// Сортировка по возрастанию
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i] > result[j] {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result, nil
}
