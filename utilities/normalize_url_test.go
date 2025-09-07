package utilities

import (
	"testing"
)

func TestNormalizeUrl(t *testing.T) {
	testUrls := map[string]string{
		"http://example.com/foo%2a":                       "http://example.com/foo%2A",
		"HTTP://User@Example.COM/Foo":                     "http://User@example.com/Foo",
		"http://example.com/%7Efoo":                       "http://example.com/~foo",
		"http://example.com/foo/./bar/baz/../qux":         "http://example.com/foo/bar/qux",
		"http://example.com":                              "http://example.com/",
		"http://example.com:80/":                          "http://example.com/",
		"http://www.example.com/../a/b/../c/./d.html":     "http://example.com/a/c/d.html",
		"https://www.youtube.com/../watch?v=Rfpt__dWr2o":  "https://youtube.com/watch?v=Rfpt__dWr2o",
		"https://www.youtube.com/a/../watch?v=Rfpt__dWr2": "https://youtube.com/watch?v=Rfpt__dWr2",
	}

	for testUrl, wantedUrl := range testUrls {
		result, err := NormalizeURL(testUrl)
		if wantedUrl != result || err != nil {
			t.Errorf(`NormalizeUrl("%s") = %s, %v, want match for %s, nil`, testUrl, result, err, wantedUrl)
		}
	}
}

func TestRemoveDotSegments(t *testing.T) {
	testPath := "/a/b/c/./../../g"
	want := "/a/g"

	res := removeDotSegments(testPath)

	if want != res {
		t.Errorf(`removeDotSegments("%s") = %s want match for %s, `, testPath, res, want)
	}
}
