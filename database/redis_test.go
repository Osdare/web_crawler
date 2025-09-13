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