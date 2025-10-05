package parser

import (
	"slices"
	"strings"
	"web_crawler/consts"
	"web_crawler/types"
	"web_crawler/utilities"

	"github.com/reiver/go-porterstemmer"
	"golang.org/x/net/html"
)

// Single pass over the html
func ParseBody(normUrl string, body *html.Node) (title string, rawUrls []string, images []types.Image, wordMap map[string]int) {
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
				word = utilities.RemovePunctuation(word)
				stem := porterstemmer.StemWithoutLowerCasing([]rune(word))

				if len(stem) >= 2 &&
					len(stem) <= 32 &&
					!slices.Contains(consts.StopWords, word) &&
					utilities.IsAlphanumeric(string(stem)) {
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
