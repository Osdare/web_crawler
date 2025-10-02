package parser

import (
	"fmt"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestWordsFromBody(t *testing.T) {
	r := strings.NewReader(Body)
	h, err := html.Parse(r)
	if err != nil {
		t.Error(err)
	}

	_, _, _, words := ParseBody("osu.ppy.sh", h)
	fmt.Println(words)
}
