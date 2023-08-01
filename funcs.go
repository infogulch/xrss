package xrss

import (
	"encoding/json"
	"html/template"

	"github.com/mmcdole/gofeed"
)

var funcLibrary template.FuncMap = template.FuncMap{
	"fetchRSS": funcFetchRSS,
}

func funcFetchRSS(url string) (string, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return "", err
	}
	b, err := json.Marshal(feed)
	return string(b), err
}
