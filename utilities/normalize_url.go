package utilities

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// Normalize url according to RFC 3986
func NormalizeURL(rawURL string) (string, error) {

	octetsMap := map[string]string{
		//A-Z
		"%41": "A",
		"%42": "B",
		"%43": "C",
		"%44": "D",
		"%45": "E",
		"%46": "F",
		"%47": "G",
		"%48": "H",
		"%49": "I",
		"%4A": "J",
		"%4B": "K",
		"%4C": "L",
		"%4D": "M",
		"%4E": "N",
		"%4F": "O",
		"%50": "P",
		"%51": "Q",
		"%52": "R",
		"%53": "S",
		"%54": "T",
		"%55": "U",
		"%56": "V",
		"%57": "W",
		"%58": "X",
		"%59": "Y",
		"%5A": "Z",

		//a-z
		"%61": "a",
		"%62": "b",
		"%63": "c",
		"%64": "d",
		"%65": "e",
		"%66": "f",
		"%67": "g",
		"%68": "h",
		"%69": "i",
		"%6A": "j",
		"%6B": "k",
		"%6C": "l",
		"%6D": "m",
		"%6E": "n",
		"%6F": "o",
		"%70": "p",
		"%71": "q",
		"%72": "r",
		"%73": "s",
		"%74": "t",
		"%75": "u",
		"%76": "v",
		"%77": "w",
		"%78": "x",
		"%79": "y",
		"%7A": "z",

		//0-9
		"%30": "0",
		"%31": "1",
		"%32": "2",
		"%33": "3",
		"%34": "4",
		"%35": "5",
		"%36": "6",
		"%37": "7",
		"%38": "8",
		"%39": "9",

		//Special
		"%2D": "-",
		"%2E": ".",
		"%5F": "_",
		"%7E": "~",
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	if u.Scheme != "https" && u.Scheme != "http" {
		return "", fmt.Errorf("url has invalid field 'Scheme'")
	}

	if u.Host == "" {
		return "", fmt.Errorf("urls has no field 'Host'")
	}

	//Scheme
	if u.Scheme == "http" {
		u.Host = strings.Replace(u.Host, ":80", "", 1)
	}

	if u.Scheme == "https" {
		u.Host = strings.Replace(u.Host, ":443", "", 1)
	}

	//Host
	u.Host = strings.ToLower(u.Host)
	u.Host = strings.TrimPrefix(u.Host, "www.")

	//Path
	u.Path = removeDotSegments(u.Path)	

	if u.Path == "" {
		u.Path = "/"
	}

	//RawPath
	//Make sure all octets are upper case (%2a -> %2A)
	rp := u.RawPath
	start := 0
	for {
		idx := strings.Index(rp[start:], "%")
		if idx == -1 {
			break
		}

		absoluteIdx := start + idx
		end := absoluteIdx + 3
		rp = rp[:absoluteIdx] + strings.ToUpper(rp[absoluteIdx:end]) + rp[end:]

		start = end + 1
		if start > len(rp) {
			break
		}
	}
	u.RawPath = rp

	var oldNew []string
	for old, new := range octetsMap {
		oldNew = append(oldNew, old, new)
	}
	replacer := strings.NewReplacer(oldNew...)
	u.RawPath = replacer.Replace(u.RawPath)

	//Query
	query := u.Query()
	if len(query) > 0 {
		keys := make([]string, 0, len(query))
		for k := range query {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		sortedQuery := make(url.Values)
		for _, k := range keys {
			sortedQuery[k] = query[k]
		}
		u.RawQuery = sortedQuery.Encode()
	}

	//Fragment
	u.Fragment = ""

	return u.String(), nil
}

func removeDotSegments(path string) string {
	parts := strings.Split(path, "/")
	stack := []string{}

	for _, part := range parts {
		switch part {
		case "", ".":
			continue
		case "..":
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		default:
			stack = append(stack, part)
		}
	}

	return "/" + strings.Join(stack, "/")
}
