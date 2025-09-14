package database

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
	"web_crawler/types"
)

type DataBase struct {
	client *redis.Client
	ctx    context.Context
}

const pageTag = "page"
const domainTag = "domain"

func (db *DataBase) Connect() error {
	addr := os.Getenv("REDIS_ADDR")
	database := os.Getenv("REDIS_DB")
	password := os.Getenv("REDIS_PASSWORD")

	dbId, err := strconv.Atoi(database)
	if err != nil {
		return err
	}

	db.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       dbId,
		Password: password,
	})

	db.ctx = context.Background()

	_, err = db.client.Ping(db.ctx).Result()
	if err != nil {
		return fmt.Errorf("couldn't connect do db %v %v", addr, err)
	}

	return nil
}

func (db *DataBase) PageExists(normPageUrl string) (bool, error) {
	res, err := db.client.Exists(db.ctx, pageTag+":"+normPageUrl).Result()
	if err != nil {
		return false, err
	}

	if res < 1 {
		return false, nil
	}

	return true, nil
}

func (db *DataBase) AddPage(page types.Page) error {
	h := sha256.New()
	h.Write([]byte(page.Content))

	checksum := h.Sum(nil)

	res, err := db.client.SIsMember(db.ctx, "contenthashes", checksum).Result()
	if err != nil {
		return err
	}
	if res {
		//Silent fail
		return nil
	}
	db.client.SAdd(db.ctx, "contenthashes", checksum)

	hashFields := []string{
		"content", page.Content,
	}

	err = db.client.HSet(db.ctx, pageTag+":"+page.NormUrl, hashFields).Err()
	if err != nil {
		return fmt.Errorf("could not add page to db %v", err)
	}

	return nil
}

func (db *DataBase) AddDomain(domain types.Domain) error {
	crawlDelay := strconv.Itoa(domain.CrawlDelay)
	lastCrawled := strconv.Itoa(domain.LastCrawled)
	disallowed := strconv.FormatBool(domain.Disallowed)

	hashFields := []string{
		"crawldelay", crawlDelay,
		"lastcrawled", lastCrawled,
		"disallowed", disallowed,
	}

	err := db.client.HSet(db.ctx, domainTag+":"+domain.Name, hashFields).Err()
	if err != nil {
		return fmt.Errorf("could not add domain %v to database %v", domain.Name, err)
	}

	return nil
}

func (db *DataBase) GetDomain(domainName string) (types.Domain, error) {

	res, err := db.client.HGetAll(db.ctx, domainTag+":"+domainName).Result()
	if err != nil {
		return types.Domain{}, fmt.Errorf("could not get domain %v %v", domainName, err)
	}

	crawlDelay, err := strconv.Atoi(res["crawldelay"])
	if err != nil {
		return types.Domain{}, fmt.Errorf("crawldelay could not be converted to int %v %v", res["crawldelay"], err)
	}

	lastCrawled, err := strconv.Atoi(res["lastcrawled"])
	if err != nil {
		return types.Domain{}, fmt.Errorf("lastcrawled could not be converted to int %v %v", res["lastcrawled"], err)
	}

	disallowed := res["disallowed"] == "true"

	return types.Domain{
		Name:        domainName,
		CrawlDelay:  crawlDelay,
		LastCrawled: lastCrawled,
		Disallowed:  disallowed,
	}, nil
}

func (db *DataBase) UpdateDomainLastCrawled(domain string, lastCrawled int) error {

	err := db.client.HSet(db.ctx, domainTag+":"+domain, "lastcrawled", lastCrawled).Err()
	if err != nil {
		return fmt.Errorf("could not update lastcrawled for domain: %v %v", domain, err)
	}

	return nil
}

// Queue stuff
func (db *DataBase) PushUrl(normUrl string) error {
	exists, err := db.PageExists(normUrl)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	if exists {
		//Silent fail
		return nil
	}

	res, err := db.client.SAdd(db.ctx, "urlset", normUrl).Result()
	if err != nil {
		return fmt.Errorf("could not add url %v to set %v", normUrl, err)
	}
	if res == 0 {
		//Silent fail
		return nil
	}

	err = db.client.LPush(db.ctx, "urlqueue", normUrl).Err()
	if err != nil {
		return fmt.Errorf("could not push %v to urlqueue %v", normUrl, err)
	}

	return nil
}

func (db *DataBase) PopUrl() (string, error) {
	res, err := db.client.BRPop(db.ctx, 0, "urlqueue").Result()
	if err != nil {
		return "", fmt.Errorf("urlqueue could not be popped %v", err)
	}

	err = db.client.SRem(db.ctx, "urlset", res[1]).Err()
	if err != nil {
		return "", fmt.Errorf("could not remove %v from urlset %v", res[1], err)
	}

	return res[1], nil
}

func (db *DataBase) UrlQueueLength() (int64, error) {
	res, err := db.client.LLen(db.ctx, "urlqueue").Result()
	if err != nil {
		return 0, fmt.Errorf("could not retrieve length of urlqueue %v", err)
	}

	return res, nil
}
