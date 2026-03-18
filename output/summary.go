package output

import (
	"fmt"

	"github.com/Prateek-1110/web_scrapper/fetcher"

	"github.com/fatih/color"
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

	green  := color.New(color.FgGreen, color.Bold).SprintfFunc()
	red    := color.New(color.FgRed, color.Bold).SprintfFunc()
	yellow := color.New(color.FgYellow, color.Bold).SprintfFunc()
	cyan   := color.New(color.FgCyan, color.Bold).SprintfFunc()

	fmt.Println()
	fmt.Println(cyan("─────────────────────────────────────────────────"))
	fmt.Println(cyan("  SCRAPE SUMMARY"))
	fmt.Println(cyan("─────────────────────────────────────────────────"))
	fmt.Printf("  %s  Success   : %d\n", green("✓"), success)
	fmt.Printf("  %s  Failed    : %d\n", red("✗"), failed)
	fmt.Printf("  %s  Retried   : %d\n", yellow("↻"), retried)
	fmt.Printf("  %s  Total     : %d\n", cyan("∑"), len(results))
	fmt.Println(cyan("─────────────────────────────────────────────────"))

	if failed > 0 {
		fmt.Println(red("\n  Failed URLs:"))
		for _, r := range results {
			if r.Error != "" {
				fmt.Printf("  · %s\n    → %s\n",
					red(r.URL),
					yellow(r.Error),
				)
			}
		}
	}
	fmt.Println()
}