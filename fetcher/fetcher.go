package fetcher

import (
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Result struct {
	URL      string
	Title    string
	Status   int
	Attempts int
	Error    string
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

		// Exponential backoff: wait 1s, 2s, 4s between retries
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

	// Retry on server errors (5xx)
	if resp.StatusCode >= 500 {
		return Result{URL: url, Attempts: attempt, Status: resp.StatusCode, Error: "server error"}, false
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return Result{URL: url, Attempts: attempt, Error: err.Error()}, false
	}

	title := doc.Find("title").Text()

	return Result{
		URL:      url,
		Title:    title,
		Status:   resp.StatusCode,
		Attempts: attempt,
	}, true
}