package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Prateek-1110/web_scrapper/fetcher"
	"github.com/Prateek-1110/web_scrapper/output"

	"github.com/schollz/progressbar/v3"
)

func main() {
	// CLI flags — all flags must be inside main()
	workers := flag.Int("workers", 10, "number of concurrent workers")
	timeout := flag.Duration("timeout", 10*time.Second, "per-request timeout")
	outJSON := flag.String("json", "results.json", "JSON output file")
	outCSV  := flag.String("csv", "results.csv", "CSV output file")
	resume  := flag.Bool("resume", false, "skip already scraped URLs")
	rate    := flag.Int("rate", 0, "max req per sec (0 = unlimited)")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Usage: go run main.go [flags] urls.txt")
		os.Exit(1)
	}

	urls, err := readURLs(args[0])
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	// Resume: filter out already scraped URLs
	if *resume {
		visited := fetcher.LoadCheckpoint()
		var remaining []string
		for _, u := range urls {
			if !visited[u] {
				remaining = append(remaining, u)
			}
		}
		skipped := len(urls) - len(remaining)
		if skipped > 0 {
			fmt.Printf("Resuming — skipping %d already scraped URLs\n", skipped)
		}
		urls = remaining
	}

	fmt.Printf("Loaded %d URLs · %d workers · timeout %s\n\n", len(urls), *workers, *timeout)

	// Graceful shutdown: listen for Ctrl+C
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	jobs    := make(chan string, len(urls))
	results := make(chan fetcher.Result, len(urls))
	bar     := progressbar.Default(int64(len(urls)))

	// Rate limiter
	var rateLimiter <-chan time.Time
	if *rate > 0 {
		ticker := time.NewTicker(time.Second / time.Duration(*rate))
		defer ticker.Stop()
		rateLimiter = ticker.C
	}

	// Launch workers
	var wg sync.WaitGroup
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range jobs {
				if rateLimiter != nil {
					<-rateLimiter
				}
				results <- fetcher.Fetch(url, *timeout)
				bar.Add(1)
			}
		}()
	}

	// Feed jobs in a goroutine so Ctrl+C can interrupt
	go func() {
		for _, u := range urls {
			select {
			case jobs <- u:
			case <-stop:
				fmt.Println("\n\nInterrupt received — saving progress...")
				close(jobs)
				return
			}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var all []fetcher.Result
	for r := range results {
		all = append(all, r)
	}

	save(all, *outJSON, *outCSV)
	fetcher.SaveCheckpoint(all)
}

func save(all []fetcher.Result, jsonFile, csvFile string) {
	output.PrintSummary(all)
	output.SaveJSON(all, jsonFile)
	output.SaveCSV(all, csvFile)
	fmt.Printf("Saved → %s\n", jsonFile)
	fmt.Printf("Saved → %s\n", csvFile)
}

func readURLs(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			urls = append(urls, line)
		}
	}
	return urls, scanner.Err()
}