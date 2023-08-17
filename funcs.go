package xrss

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
)

var funcLibrary template.FuncMap = template.FuncMap{
	"fetchFeed": funcFetchFeed,
}

type FetchResult struct {
	StatusCode      int           `json:"status_code"`
	URL             string        `json:"url"`
	Etag            string        `json:"etag"`
	LastModified    time.Time     `json:"last_modified"`
	FetchedOn       time.Time     `json:"fetched_on"`
	FetchedDuration time.Duration `json:"fetched_duration"`
	FetchedBytes    int           `json:"fetched_bytes"`
	FetchedCount    int           `json:"fetched_count"`
	Feed            *gofeed.Feed  `json:"feed,omitempty"`
}

var gmtTimeZoneLocation *time.Location = must(time.LoadLocation("GMT"))

func funcFetchFeed(url string, etag string, lastModified string) (FetchResult, error) {
	start := time.Now().In(gmtTimeZoneLocation)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return FetchResult{}, err
	}
	req.Header.Set("User-Agent", "xrss/1.0")

	if etag != "" {
		req.Header.Set("If-None-Match", fmt.Sprintf(`"%s"`, etag))
	}

	if lastModified != "" {
		req.Header.Set("If-Modified-Since", lastModified)
	}

	resp, err := http.DefaultClient.Do(req)

	result := FetchResult{
		StatusCode:      resp.StatusCode,
		URL:             url,
		FetchedOn:       start,
		FetchedDuration: time.Since(start),
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

	counter := &CountingReader{reader: resp.Body}
	result.Feed, err = parser.Parse(counter)
	if err != nil {
		return result, err
	}
	result.FetchedBytes = counter.BytesRead
	result.FetchedCount = len(result.Feed.Items)

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

type CountingReader struct {
	reader    io.Reader
	BytesRead int
}

func (r *CountingReader) Read(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)
	r.BytesRead += n
	return n, err
}
