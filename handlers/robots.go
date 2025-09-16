package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"web_crawler/types"
)

func DownloadRobots(domainName string) ([]string, error) {
	resp, err := http.Get(domainName + "/robots.txt")
	if err != nil {
		return make([]string, 0), fmt.Errorf("could not get robots from %v %v", domainName, err)
	}
	defer resp.Body.Close()

	bb, err := io.ReadAll(resp.Body)
	if err != nil {
		return make([]string, 0), fmt.Errorf("could not read bytes from body %v", err)
	}

	strings := make([]string, 0)
	s := ""
	for _, b := range bb {
		if b == '\n' {
			strings = append(strings, s)
			s = ""
		} else {
			s += string(b)
		}
	}

	return strings, nil
}

func RobotsToDomain(domainName string, robotsLines []string) (types.Domain, error) {
	domain := types.Domain{}
	domain.Name = domainName

	allow := make([]string, 0)
	disallow := make([]string, 0)
	isInUser := false
	for _, line := range robotsLines {

		if strings.HasPrefix(line, "User-agent: *") {
			isInUser = true
		} else if isInUser && strings.HasPrefix(line, "User-agent") {
			break
		}

		if isInUser {
			if strings.HasPrefix(line, "Allow") {
				s := strings.TrimSpace(strings.Replace(line, "Allow: ", "", 1))
				allow = append(allow, s)

			} else if strings.HasPrefix(line, "Disallow") {
				s := strings.TrimSpace(strings.Replace(line, "Disallow: ", "", 1))
				disallow = append(disallow, s)

			} else if strings.HasPrefix(strings.ToLower(line), "crawl-delay") {

				cd, err := strconv.Atoi(strings.Split(line, " ")[1])
				if err != nil {
					return types.Domain{}, fmt.Errorf("could not convert line: %v to int %v", line, err)
				}

				domain.CrawlDelay = cd
			}
		}
	}

	domain.Allowed = allow
	domain.Disallowed = disallow
	domain.LastCrawled = int(time.Now().Unix())

	return domain, nil
}
