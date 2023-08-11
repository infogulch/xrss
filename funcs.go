package xrss

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
)

var funcLibrary template.FuncMap = template.FuncMap{
	"fetchFeed": funcFetchFeed,
}

type FetchResult struct {
	ResponseStatus int
	URL            string
	Etag           string
	LastModified   time.Time
	FetchTime      time.Time
	FetchDuration  time.Duration
	*gofeed.Feed
}

var gmtTimeZoneLocation *time.Location = must(time.LoadLocation("GMT"))

func funcFetchFeed(url string, etag string, lastModified time.Time) (FetchResult, error) {
	start := time.Now().In(gmtTimeZoneLocation)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return FetchResult{}, err
	}
	req.Header.Set("User-Agent", "xrss/1.0")

	if etag != "" {
		req.Header.Set("If-None-Match", fmt.Sprintf(`"%s"`, etag))
	}

	if !lastModified.IsZero() {
		req.Header.Set("If-Modified-Since", lastModified.In(gmtTimeZoneLocation).Format(time.RFC1123))
	}

	resp, err := http.DefaultClient.Do(req)

	result := FetchResult{
		ResponseStatus: resp.StatusCode,
		URL:            url,
		FetchTime:      start,
		FetchDuration:  time.Since(start),
	}

	if err != nil {
		return result, err
	}

	if resp != nil {
		defer func() {
			ce := resp.Body.Close()
			if ce != nil {
				err = ce
			}
		}()
	}

	if resp.StatusCode == http.StatusNotModified {
		return result, nil
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return result, fmt.Errorf("feed fetch failed: status=%d", resp.StatusCode)
	}

	parser := gofeed.NewParser()
	// maybe expose .AuthConfig or .Client

	result.Feed, err = parser.Parse(resp.Body)
	if err != nil {
		return result, err
	}

	if eTag := resp.Header.Get("Etag"); eTag != "" {
		result.Etag = eTag
	}

	if lastModified := resp.Header.Get("Last-Modified"); lastModified != "" {
		parsed, err := time.ParseInLocation(time.RFC1123, lastModified, gmtTimeZoneLocation)
		if err == nil {
			result.LastModified = parsed
		}
	}

	return result, nil
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
