package sorter

import (
	"bufio"
	"cmp"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Sorter struct {
	data []string
}

var titleCaser = cases.Title(language.Und)

func New() Sorter {
	return Sorter{[]string{}}
}

func (s *Sorter) Sort() {
	slices.Sort(s.data)
}

func (s *Sorter) SortByColumn(column int, numeric bool) {
	pos := column
	if column != 0 {
		pos = column - 1
	}

	sort.Slice(s.data, func(i, j int) bool {
		keyI := keyFromColumn(s.data[i], pos)
		keyJ := keyFromColumn(s.data[j], pos)

		if numeric {
			numI, _ := leadingNumber(keyI)
			numJ, _ := leadingNumber(keyJ)
			if numI != numJ {
				return numI < numJ
			}
			// tie-breaker: вся строка
			return s.data[i] < s.data[j]
		}

		if keyI != keyJ {
			return keyI < keyJ
		}
		// tie-breaker: вся строка
		return s.data[i] < s.data[j]
	})
}

func (s *Sorter) SortNumeric() {
	slices.SortFunc(s.data, func(a, b string) int {
		na, okA := leadingNumber(a)
		nb, okB := leadingNumber(b)

		if !okA {
			na = 0
		}
		if !okB {
			nb = 0
		}

		if na < nb {
			return -1
		}
		if na > nb {
			return 1
		}
		return cmp.Compare(a, b)
	})
}

func (s *Sorter) Reverse() {
	slices.Reverse(s.data)
}

func (s *Sorter) Unique() {
	unique := make(map[string]struct{})
	for _, line := range s.data {
		unique[line] = struct{}{}
	}

	storage := make([]string, 0, len(s.data))
	for k := range unique {
		storage = append(storage, k)
	}

	s.data = storage
}

func (s *Sorter) SortByMonth(column int) {
	monthMap := map[string]time.Month{
		"Jan": time.January, "Feb": time.February, "Mar": time.March,
		"Apr": time.April, "May": time.May, "Jun": time.June,
		"Jul": time.July, "Aug": time.August, "Sep": time.September,
		"Oct": time.October, "Nov": time.November, "Dec": time.December,
	}

	pos := column - 1
	if pos < 0 {
		pos = 0
	}

	slices.SortFunc(s.data, func(a, b string) int {
		fieldsA := strings.Fields(a)
		fieldsB := strings.Fields(b)

		var ma, mb time.Month

		getMonth := func(fields []string) time.Month {
			if len(fields) <= pos {
				return 0
			}
			field := strings.TrimSpace(fields[pos])
			monthStr := titleCaser.String(strings.ToLower(field))
			return monthMap[monthStr]
		}

		ma = getMonth(fieldsA)
		mb = getMonth(fieldsB)

		aVal := int(ma)
		bVal := int(mb)

		if aVal == 0 {
			aVal = -1
		}
		if bVal == 0 {
			bVal = -1
		}

		if aVal < bVal {
			return -1
		} else if aVal > bVal {
			return 1
		}
		return 0
	})
}

func (s *Sorter) TrimTrailingSpaces() {
	for i, v := range s.data {
		s.data[i] = strings.TrimRight(v, " \t\r")
	}
}

func (s *Sorter) IsSorted(fileName string) {
	if fileName == "" {
		fileName = "-"
	}

	messageFormat := "sort: %s:%d: disorder: %s"
	for i := 1; i < len(s.data); i++ {
		if s.data[i] < s.data[i-1] {
			fmt.Printf(messageFormat, fileName, i+1, s.data[i])
			return
		}
	}
}

func (s *Sorter) SortHumanNumeric() {
	multiplier := map[byte]float64{
		'K': 1 << 10, 'k': 1 << 10,
		'M': 1 << 20, 'm': 1 << 20,
		'G': 1 << 30, 'g': 1 << 30,
	}

	parse := func(str string) (float64, float64) {
		str = strings.TrimSpace(str)
		if str == "" {
			return 0, 0
		}
		last := str[len(str)-1]
		if mul, ok := multiplier[last]; ok {
			num, err := strconv.ParseFloat(str[:len(str)-1], 64)
			if err != nil {
				return 0, 0
			}
			return num, mul
		}
		num, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return 0, 0
		}
		return num, 1
	}

	slices.SortFunc(s.data, func(a, b string) int {
		na, ma := parse(a)
		nb, mb := parse(b)

		if ma == 1 && mb == 1 {
			if na < nb {
				return -1
			} else if na > nb {
				return 1
			}
			return 0
		}

		if ma == 1 && mb != 1 {
			// a без суффикса, b с суффиксом -> a идёт раньше
			return -1
		}
		if ma != 1 && mb == 1 {
			// a с суффиксом, b без суффикса -> b идёт раньше
			return 1
		}

		valA := na * ma
		valB := nb * mb
		if valA < valB {
			return -1
		} else if valA > valB {
			return 1
		}
		return 0
	})
}

func (s *Sorter) Data() []string {
	return s.data
}

func (s *Sorter) ReadData(input io.Reader) {
	reader := bufio.NewReader(input)
	for {
		nextString, err := reader.ReadString('\n')
		if err == io.EOF {
			nextString = nextString + "\n"
		}

		s.data = append(s.data, nextString)

		if err == io.EOF {
			break
		}
	}
}

func leadingNumber(s string) (float64, bool) {
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	if i >= len(s) {
		return 0, false
	}
	start := i

	if s[i] == '+' || s[i] == '-' {
		i++
	}

	digitsBefore := 0
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
		digitsBefore++
	}

	digitsAfter := 0
	if i < len(s) && s[i] == '.' {
		i++
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			i++
			digitsAfter++
		}
	}

	if digitsBefore == 0 && digitsAfter == 0 {
		return 0, false
	}

	if i < len(s) && (s[i] == 'e' || s[i] == 'E') {
		j := i + 1
		if j < len(s) && (s[j] == '+' || s[j] == '-') {
			j++
		}
		expDigits := 0
		for j < len(s) && s[j] >= '0' && s[j] <= '9' {
			j++
			expDigits++
		}
		if expDigits > 0 {
			i = j
		}
	}

	numStr := s[start:i]
	if numStr == "" {
		return 0, false
	}
	f, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, false
	}
	return f, true
}

func keyFromColumn(line string, pos int) string {
	fields := strings.Fields(line)
	if len(fields) <= pos {
		return "" // нет такого столбца
	}
	return strings.Join(fields[pos:], " ")
}
