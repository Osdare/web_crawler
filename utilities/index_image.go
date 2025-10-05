package utilities

import (
	"github.com/reiver/go-porterstemmer"
	"slices"
	"strings"
	"web_crawler/consts"
	"web_crawler/types"
)

func IndexImage(image types.Image) map[string]int {
	m := make(map[string]int)

	w := strings.SplitSeq(image.AltText, " ")
	for word := range w {

		//normalizing and stemming
		word = strings.ToLower(word)
		word = strings.TrimSpace(word)
		word = RemovePunctuation(word)
		stem := porterstemmer.StemWithoutLowerCasing([]rune(word))

		if len(stem) >= 2 &&
			len(stem) <= 32 &&
			!slices.Contains(consts.StopWords, word) &&
			IsAlphanumeric(string(stem)) {
			m[string(stem)]++
		}
	}
	

	return m
}
