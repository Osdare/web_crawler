package utilities

import (
	"testing"
)

func TestNormalizeUrl(t *testing.T) {
	testUrls := map[string]string{
		"http://example.com/foo%2a":               "http://example.com/foo%2A",
		"http://user@example.com/foo":             "http://user@example.com/foo",
		"http://example.com/%7Efoo":               "http://example.com/~foo",
		"http://example.com/foo/./bar/baz/../qux": "http://example.com/foo/bar/qux",
		"http://example.com":                      "http://example.com/",
		"http://example.com:80/":                  "http://example.com/",
	}

	for testUrl, wantedUrl := range testUrls {
		result, err := NormalizeURL(testUrl)
		if wantedUrl != result || err != nil {
			t.Errorf(`NormalizeUrl("%s") = %s, %v, want match for %s, nil`, testUrl, result, err, wantedUrl)
		}
	}
}
