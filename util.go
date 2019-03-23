package md

import (
	"strings"
	"unicode"
)

func countLeft(s string, r rune) int {
	var count int
	for _, sr := range s {
		if sr != r {
			return count
		}
		count++
	}
	return count
}

func isEmpty(s string) bool {
	return strings.TrimLeftFunc(s, unicode.IsSpace) == ""
}
