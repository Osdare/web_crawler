package crawler

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/url"
	"web_crawler/database"
	"web_crawler/handlers"
	"web_crawler/parser"
	"web_crawler/types"
	"web_crawler/utilities"
)

// returns html as node
func Crawl(normUrl string) (*html.Node, string, error) {
	resp, err := http.Get(normUrl)
	if err != nil {
		return nil, "", fmt.Errorf("could not get url: %v %v", normUrl, err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("could not read body %v", err)
	}

	html, err := html.Parse(bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, "", fmt.Errorf("could not parse response body %v", err)
	}

	content := string(bodyBytes)

	return html, content, nil
}

func CrawlJob(db *database.DataBase) {
	//infinite loop
	ql, err := db.UrlQueueLength()
	if err != nil {
		log.Printf("aaaaaaaaa")
	}
	for ql > 0 {
		fmt.Println("hello :D")
		//get url from database
		link, err := db.PopUrl()
		if err != nil {
			log.Printf("could not pop url from db %v\n", err)
			continue
		}
		log.Printf("Crawling: %v", link)

		u, err := url.Parse(link)
		if err != nil {
			log.Printf("could not parse url: %v %v\n", link, err)
			continue
		}

		//check if domain exists and create a new one if not

		domainExists, err := db.DomainExists(u.Host)
		if err != nil {
			log.Printf("couldn not check if domain: %v exists %v\n", u.Host, err)
			continue
		}

		var domain types.Domain
		if domainExists {
			domain, err = db.GetDomain(u.Host)
			if err != nil {
				log.Printf("could not get domain: %v %v\n", u.Host, err)
				continue
			}
		} else {
			domain, err = handlers.GetRobotsFromDomain(u.Scheme + "://" + u.Host)
			if err != nil {
				log.Printf("could not get new domain: %v reason: %v", u.Scheme+"://"+u.Host, err)
				continue
			}

			err = db.AddDomain(domain)
			if err != nil {
				log.Printf("could not add domain: %v to db. Reason: %v", u.Host, err)
				continue
			}
		}

		//check domain and if we are allowed to crawl else continue
		//and put url back in the queue
		canCrawl, reason, err := handlers.CanCrawl(link, domain)
		if err != nil {
			log.Printf("error from CanCrawl function %v\n", err)
			continue
		}
		if !canCrawl {
			log.Printf("not allowed to crawl url: %v in domain: %v\n", link, domain.Name)
			if reason == handlers.ReasonCrawlDelay {
				err = db.PushUrl(link)
			}
			if err != nil {
				log.Printf("could not push url to queue. Reason: %v\n", err)
			}
			continue
		}

		//crawl it
		html, content, err := Crawl(link)
		if err != nil {
			log.Printf("Could not crawl url: %v %v\n", link, err)
			continue
		}

		//create new page
		//nu := strings.Replace(link, u.Scheme+"://", "", 1)
		page := types.Page{
			NormUrl: link,
			Content: content,
		}
		err = db.AddPage(page)
		if err != nil {
			log.Printf("page: %v could not be added to database %v\n", link, err)
		}

		//normalize urls and put new urls in database
		newUrls, err := parser.UrlsFromBody(html)
		if err != nil {
			log.Printf("Could not parse urls from html %v\n", err)
			continue
		}

		urlsToAdd := utilities.NormalizeUrlSlice(newUrls)
		for _, newUrl := range urlsToAdd {
			//log.Printf("attempting to add url: %v to queue\n", newUrl)
			err = db.PushUrl(newUrl)
			if err != nil {
				log.Printf("Could not add url: %v to queue %v\n", newUrl, err)
			}
		}

		ql, err = db.UrlQueueLength()
		log.Printf("queue length after adding urls: %v\n",ql)
		if err != nil {
			log.Printf("ASDFASDF")
		}
	}
	fmt.Println("empty :D")
}
