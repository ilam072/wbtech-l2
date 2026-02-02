package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	s, err := unpack(`a4bc2d5e`)
	fmt.Println(s, err)

	s, err = unpack(`abcd`)
	fmt.Println(s, err)

	s, err = unpack(`45`)
	fmt.Println(s, err)

	s, err = unpack(``)
	fmt.Println(s, err)

	s, err = unpack(`qwe\4\5`)
	fmt.Println(s, err)

	s, err = unpack(`qwe\45`)
	fmt.Println(s, err)

	s, err = unpack(`qwe\\5`)
	fmt.Println(s, err)
}

func unpack(s string) (string, error) {
	var (
		builder  strings.Builder
		num      int
		err      error
		escape   bool
		prevRune rune
		hasRune  bool
	)

	if len(s) == 0 {
		return "", nil
	}

	for _, r := range s {
		if r == '\\' && !escape {
			escape = true
			continue
		}

		if unicode.IsDigit(r) {
			if !hasRune {
				return "", errors.New("invalid string")
			}

			if escape {
				builder.WriteRune(r)
				prevRune = r
				hasRune = true
				escape = false
				continue
			}

			num, err = strconv.Atoi(string(r))
			if err != nil {
				return "", err
			}

			builder.WriteString(strings.Repeat(string(prevRune), num-1))
			continue
		}

		escape = false
		builder.WriteRune(r)
		prevRune = r
		hasRune = true
	}

	return builder.String(), nil
}
