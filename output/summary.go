package output

import (
	"fmt"
	"strings"

	"scraper/fetcher"
)

func PrintSummary(results []fetcher.Result) {
	success := 0
	failed  := 0
	retried := 0

	for _, r := range results {
		if r.Error != "" {
			failed++
		} else {
			success++
		}
		if r.Attempts > 1 {
			retried++
		}
	}

	fmt.Println()
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println("  SCRAPE SUMMARY")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Printf("  ✓  Success   : %d\n", success)
	fmt.Printf("  ✗  Failed    : %d\n", failed)
	fmt.Printf("  ↻  Retried   : %d\n", retried)
	fmt.Printf("  ∑  Total     : %d\n", len(results))
	fmt.Println(strings.Repeat("─", 50))

	if failed > 0 {
		fmt.Println("\n  Failed URLs:")
		for _, r := range results {
			if r.Error != "" {
				fmt.Printf("  · %s\n    → %s\n", r.URL, r.Error)
			}
		}
	}
	fmt.Println()
}
