package output

import (
	"encoding/csv"
	"encoding/json"
	"os"

	"scraper/fetcher"
)

func SaveJSON(results []fetcher.Result, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}

func SaveCSV(results []fetcher.Result, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"URL", "Title", "Error"})

	for _, r := range results {
		writer.Write([]string{r.URL, r.Title, r.Error})
	}

	return nil
}