package fetcher

import (
	"encoding/json"
	"os"
)

const checkpointFile = ".scraper_checkpoint.json"

func LoadCheckpoint() map[string]bool {
	visited := make(map[string]bool)

	data, err := os.ReadFile(checkpointFile)
	if err != nil {
		return visited
	}

	json.Unmarshal(data, &visited)
	return visited
}

func SaveCheckpoint(results []Result) {
	visited := make(map[string]bool)
	for _, r := range results {
		if r.Error == "" {
			visited[r.URL] = true
		}
	}
	data, _ := json.Marshal(visited)
	os.WriteFile(checkpointFile, data, 0644)
}

func ClearCheckpoint() {
	os.Remove(checkpointFile)
}