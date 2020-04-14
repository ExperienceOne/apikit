package stringutil

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func InStringSlice(strings []string, s string) bool {
	for _, i := range strings {
		if i == s {
			return true
		}
	}
	return false
}

func UnTitle(s string) string {

	if len(s) == 0 {
		return s
	}

	r, width := utf8.DecodeRuneInString(s)
	return fmt.Sprintf("%c%s", unicode.ToLower(r), s[width:])
}
