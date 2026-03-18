package fetcher

import (
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Result struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StatusCode  int    `json:"status_code"`
	LinkCount   int    `json:"link_count"`
	ImageCount  int    `json:"image_count"`
	Attempts    int    `json:"attempts"`
	Error       string `json:"error,omitempty"`
}

func Fetch(url string, timeout time.Duration) Result {
	return fetchWithRetry(url, timeout, 3)
}

func fetchWithRetry(url string, timeout time.Duration, maxAttempts int) Result {
	client := &http.Client{Timeout: timeout}

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		result, done := tryFetch(client, url, attempt)
		if done {
			return result
		}
		wait := time.Duration(1<<uint(attempt-1)) * time.Second
		time.Sleep(wait)
	}

	return Result{
		URL:      url,
		Attempts: maxAttempts,
		Error:    "failed after all attempts",
	}
}

func tryFetch(client *http.Client, url string, attempt int) (Result, bool) {
	resp, err := client.Get(url)
	if err != nil {
		return Result{URL: url, Attempts: attempt, Error: err.Error()}, false
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return Result{URL: url, Attempts: attempt, StatusCode: resp.StatusCode, Error: "server error"}, false
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return Result{URL: url, Attempts: attempt, Error: err.Error()}, false
	}

	return Result{
		URL:         url,
		Title:       doc.Find("title").Text(),
		Description: doc.Find(`meta[name="description"]`).AttrOr("content", ""),
		StatusCode:  resp.StatusCode,
		LinkCount:   doc.Find("a").Length(),
		ImageCount:  doc.Find("img").Length(),
		Attempts:    attempt,
	}, true
}