package crawler

import (
	"fmt"
	"testing"
	"web_crawler/database"
	"web_crawler/parser"
	"web_crawler/utilities"
)

func TestCrawlAndParse(t *testing.T) {
	html, _, err := Crawl("https://osu.ppy.sh/users/5070783")
	if err != nil {
		t.Error(err)
	}

	urls, err := parser.UrlsFromBody(html)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(len(urls))
	normUrls := utilities.NormalizeUrlSlice(urls)
	fmt.Println(len(normUrls))
}

func TestCrawlJob(t *testing.T) {
	db := database.DataBase{}
	db.Connect("localhost:6379", "0", "")

	err := db.PushUrl("https://osu.ppy.sh/")
	if err != nil {
		t.Error(err)
	}

	CrawlJob(&db)
}
