package parser

import (
	"golang.org/x/net/html"
)

func UrlsFromBody(body *html.Node) []string {
	urls := make([]string, 0)
		
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					urls = append(urls, attr.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(body)

	return urls
}