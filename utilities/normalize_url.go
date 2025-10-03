package utilities

import (
	"log"
	"net/url"
)

func NormalizeLink(base string, link string) (string, error) {
	baseURL, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	ref, err := url.Parse(link)
	if err != nil {
		return "", err
	}

	absURL := baseURL.ResolveReference(ref)
	return absURL.String(), nil
}

func NormalizeUrlSlice(base string, rawUrls []string) []string {
	normalizedUrls := make([]string, 0)

	for _, rawUrl := range rawUrls {
		normUrl, err := NormalizeLink(base, rawUrl)
		if err != nil {
			log.Println(err)
		}

		if normUrl != "" {
			normalizedUrls = append(normalizedUrls, normUrl)
		}
	}

	return normalizedUrls
}
