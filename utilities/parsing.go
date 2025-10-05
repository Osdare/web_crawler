package utilities

import (
	"strings"
	"unicode"
)

func RemovePunctuation(s string) string {
	replacer := strings.NewReplacer(
		",", "",
		".", "",
		";", "",
		":", "",
		"!", "",
		"?", "",
	)

	return replacer.Replace(s)
}

func IsAlphanumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
