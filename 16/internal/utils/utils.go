package utils

import "net/url"

func Normalize(u *url.URL) string {
	uCopy := *u
	uCopy.Fragment = "" // выкидываем #
	q := uCopy.Query()
	for k := range q {
		if len(k) >= 4 && k[:4] == "utm_" {
			q.Del(k)
		}
	}
	uCopy.RawQuery = q.Encode()
	return uCopy.String()
}
