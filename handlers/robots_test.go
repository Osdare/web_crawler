package handlers

import (
	"testing"
	"web_crawler/types"
)

func TestRobotsToDomain(t *testing.T) {

	domainName := "https://google.com"	

	domain, err := GetRobotsFromDomain(domainName)
	if err != nil {
		t.Error(err)
	}

	if domain.Disallowed[0] != "/search" {
		t.Error("expected /search got: ", domain.Disallowed[0])
	}

}

func TestCanCrawl(t *testing.T) {
	domainG := types.Domain{
		Name: "https://google.com",
		CrawlDelay: 1,
		LastCrawled: 1234,
		Allowed: []string{"/search/about", "/search/howsearchworks", "/example/page/", "/example/allowed.gif"},
		Disallowed: []string{"/search", "/example/page/disallowed.gif", "*.gif$"},
	}

	expected := map[string]bool{
		"https://google.com/search/about": true,
		"https://google.com/search/foo": false,
		"https://google.com/search/howsearchworks": true,
		"https://google.com/search/yeah": false,
		"https://google.com/example/page/": true,
		"https://google.com/example/page/disallowed.gif": false,
		"https://google.com/example/page/asdf.gif?size=large": true,
		"https://google.com/example/allowed.gif": true,
		"https://google.com/lolo.gif": false,
	}

	for key, val := range expected {
		res, err := CanCrawl(key, domainG)
		if err != nil {
			t.Error(err)
		}

		if res != val {
			t.Errorf("got %v expected %v for domain %v", res, val, key)
		}
	}
}
