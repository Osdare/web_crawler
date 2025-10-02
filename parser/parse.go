package parser

import (
	"slices"
	"strings"
	"web_crawler/types"

	"github.com/reiver/go-porterstemmer"
	"golang.org/x/net/html"
)

var stopWords = []string{"i", "me", "my", "myself", "we", "our", "ours", "ourselves", "you", "your", "yours", "yourself", "yourselves", "he", "him", "his", "himself", "she", "her", "hers", "herself", "it", "its", "itself", "they", "them", "their", "theirs", "themselves", "what", "which", "who", "whom", "this", "that", "these", "those", "am", "is", "are", "was", "were", "be", "been", "being", "have", "has", "had", "having", "do", "does", "did", "doing", "a", "an", "the", "and", "but", "if", "or", "because", "as", "until", "while", "of", "at", "by", "for", "with", "about", "against", "between", "into", "through", "during", "before", "after", "above", "below", "to", "from", "up", "down", "in", "out", "on", "off", "over", "under", "again", "further", "then", "once", "here", "there", "when", "where", "why", "how", "all", "any", "both", "each", "few", "more", "most", "other", "some", "such", "no", "nor", "not", "only", "own", "same", "so", "than", "too", "very", "s", "t", "can", "will", "just", "don", "should", "now"}

// Single pass over the html
func ParseBody(normUrl string, body *html.Node) (title string, rawUrls []string, images []types.Image, wordMap map[string]int) {
	//word and count
	wordMap = make(map[string]int)
	images = make([]types.Image, 0)
	rawUrls = make([]string, 0)

	var f func(*html.Node)
	f = func(n *html.Node) {
		switch n.Type {
		case html.TextNode:
			w := strings.SplitSeq(n.Data, " ")
			for word := range w {
				//normalizing and stemming
				word = strings.ToLower(word)
				word = strings.TrimSpace(word)
				word = removePunctuation(word)
				stem := porterstemmer.StemWithoutLowerCasing([]rune(word))
				if !slices.Contains(stopWords, string(stem)) {
					wordMap[string(stem)]++
				}
			}
		case html.ElementNode:
			if n.Data == "title" && n.FirstChild != nil {
				title = n.FirstChild.Data
			} else if n.Data == "a" {
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						rawUrls = append(rawUrls, attr.Val)
					}
				}
			} else if n.Data == "img" {
				image := types.Image{
					PageUrl: normUrl,
				}
				for _, attr := range n.Attr {
					if attr.Key == "src" {
						image.ImageUrl = attr.Val
					}
					if attr.Key == "alt" {
						image.AltText = attr.Val
					}
				}
				images = append(images, image)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(body)
	return title, rawUrls, images, wordMap
}
func removePunctuation(s string) string {
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
