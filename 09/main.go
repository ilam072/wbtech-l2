package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
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
		builder strings.Builder
		num     int
		err     error
		escape  bool
	)

	if utf8.RuneCountInString(s) == 0 {
		return "", errors.New("invalid string")
	}

	for i, r := range s {
		if r == '\\' && !escape {
			escape = true
			continue
		}

		if unicode.IsDigit(r) {
			if i == 0 {
				return "", errors.New("invalid string")
			}

			if escape {
				builder.WriteString(string(r))
				escape = false
				continue
			}

			num, err = strconv.Atoi(string(r))
			if err != nil {
				return "", err
			}
			builder.WriteString(strings.Repeat(string(s[i-1]), num-1))
			continue
		}

		escape = false
		builder.WriteString(string(r))
	}

	return builder.String(), nil
}
