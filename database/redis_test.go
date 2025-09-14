package database

import (
	"reflect"
	"testing"
	"time"
	"web_crawler/types"
	"os"
)

func TestAddDomainGetDomain(t *testing.T) {
	os.Setenv("REDIS_ADDR", "localhost:6379")
	os.Setenv("REDIS_DB","0")
	os.Setenv("REDIS_PROTOCOL", "2")
	
	db := DataBase{}
	err := db.Connect()
	if err != nil {
		t.Errorf("could not connect to database %v", err)
	}
	
	domain := types.Domain{
		Name: "google.com",
		CrawlDelay: 10,
		LastCrawled: int(time.Now().Unix()),
		Disallowed: true,
	}

	err = db.AddDomain(domain)
	if err != nil {
		t.Errorf("domain: %v could not be added to database %v", domain.Name, err)
	}

	dbDomain, err := db.GetDomain(domain.Name)
	if err != nil {
		t.Errorf("got en error when trying to fetch domain %v, %v", domain.Name, err)	
	}

	if !reflect.DeepEqual(dbDomain, domain) {
		t.Error("domains are not equal ", domain, dbDomain )
	}
}

func TestUrlQueue(t *testing.T) {
	tesurls := []string{
		"google.com/oogabooga",
		"x.com/nooo",
		"instagram.com/shit",
		"instagram.com/shit",
		"instagram.com/shit",
		"example.ax/fisk",
	}

	db := DataBase{}
	err := db.Connect()
	if err != nil {
		t.Errorf("could not connect to db %v", err)
	}

	for _, url := range tesurls {
		err = db.PushUrl(url)
		if err != nil {
			t.Errorf("could not push url %v %v", url, err)
		}

	}

	length, err := db.UrlQueueLength()
	if err != nil {
		t.Error(err)
	}

	if length != 4 {
		t.Errorf("expected 4 got %d", length)
	}

	//ADD a url to the db and then try to re-add it to the queue

	url, err := db.PopUrl()
	if err != nil {
		t.Error(err)
	}

	urlPage := types.Page{
		NormUrl: url,
		Content: "i am a random piece of html",
	}

	err = db.AddPage(urlPage)
	if err != nil {
		t.Error(err)
	}

	err = db.PushUrl(url)
	if err != nil {
		t.Error(err)
	}

	length, err = db.UrlQueueLength()
	if err != nil {
		t.Error(err)
	}

	if length != 3 {
		t.Errorf("length was %d expected 3", length)
	}


}