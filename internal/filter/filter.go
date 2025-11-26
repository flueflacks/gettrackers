package filter

import (
	"os"
	"strings"

	"github.com/flueflacks/gettrackers/internal/paths"
)

// LoadBlocklist loads blocklist patterns from file
func LoadBlocklist() ([]string, error) {
	blocklistFile, err := paths.GetBlocklistFile()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(blocklistFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	patterns := []string{}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			patterns = append(patterns, line)
		}
	}

	return patterns, nil
}

// Filter removes URLs that match any blocklist pattern
func Filter(urls []string, blocklist []string) []string {
	if len(blocklist) == 0 {
		return urls
	}

	filtered := []string{}
	for _, url := range urls {
		shouldBlock := false
		for _, pattern := range blocklist {
			if strings.Contains(url, pattern) {
				shouldBlock = true
				break
			}
		}
		if !shouldBlock {
			filtered = append(filtered, url)
		}
	}

	return filtered
}

// AddToBlocklist adds a pattern to the blocklist if it doesn't already exist
func AddToBlocklist(pattern string) error {
	blocklistFile, err := paths.GetBlocklistFile()
	if err != nil {
		return err
	}

	// Load existing blocklist
	existing, err := LoadBlocklist()
	if err != nil {
		return err
	}

	// Check for duplicates
	pattern = strings.TrimSpace(pattern)
	for _, existingPattern := range existing {
		if existingPattern == pattern {
			// Already exists, return success
			return nil
		}
	}

	// Append new pattern
	existing = append(existing, pattern)
	content := strings.Join(existing, "\n") + "\n"

	if err := os.WriteFile(blocklistFile, []byte(content), 0644); err != nil {
		return err
	}

	return nil
}

