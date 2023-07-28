package xrss

import (
	"encoding/json"

	"github.com/mmcdole/gofeed"
)

func funcParseRSS(url string) (string, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return "", err
	}
	b, err := json.Marshal(feed)
	return string(b), err
}
