package grep

import (
	"regexp"
	"strings"
)

// matcher абстрагирует проверку совпадений строки с шаблоном.
type matcher interface {
	Match(line string) bool
}

type fixedMatcher struct{ needle string }

type regexMatcher struct{ re *regexp.Regexp }

func (m fixedMatcher) Match(line string) bool {
	return strings.Contains(line, m.needle)
}

func (m regexMatcher) Match(line string) bool {
	return m.re.MatchString(line)
}

func newMatcher(pat string, fixed, ignoreCase bool) (matcher, error) {
	if fixed {
		if ignoreCase {
			return fixedMatcher{needle: strings.ToLower(pat)}, nil
		}
		return fixedMatcher{needle: pat}, nil
	}
	flags := ""
	if ignoreCase {
		flags = "(?i)"
	}
	re, err := regexp.Compile(flags + pat)
	if err != nil {
		return nil, err
	}
	return regexMatcher{re: re}, nil
}

func normalizeLine(s string, fixed, ignoreCase bool) string {
	if fixed && ignoreCase {
		return strings.ToLower(s)
	}
	return s
}
