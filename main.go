package main

import (
	"fmt"
	"web_crawler/database"
	"web_crawler/crawler"
)

func main() {
	fmt.Println(":)")

	db := database.DataBase{}
	db.Connect("localhost:6379", "0", "")

	err := db.PushUrl("https://osu.ppy.sh/")
	if err != nil {
		panic(err)
	}

	crawler.CrawlJob(&db)
}
