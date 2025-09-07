package utilities

import (
	"fmt"
	"net/url"
)

func GetDomainFromUrl(rawUrl string) (string, error) {
	parsedUrl, err := url.Parse(rawUrl)

	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	return parsedUrl.Host, nil
}
