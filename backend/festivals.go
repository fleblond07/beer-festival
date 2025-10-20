package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func loadFestivals(filePath string) ([]Festival, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read festivals file: %w", err)
	}

	var festivals []Festival
	if err := json.Unmarshal(file, &festivals); err != nil {
		return nil, fmt.Errorf("failed to parse festivals JSON: %w", err)
	}

	return festivals, nil
}
