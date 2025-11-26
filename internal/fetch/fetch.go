package fetch

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/flueflacks/gettrackers/internal/paths"
)

const cacheMaxAge = 24 * time.Hour

// ShouldFetch checks if the cache should be refreshed
func ShouldFetch(force bool) (bool, error) {
	if force {
		return true, nil
	}

	cacheFile, err := paths.GetCacheFile()
	if err != nil {
		return false, err
	}

	info, err := os.Stat(cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, fmt.Errorf("failed to stat cache file: %w", err)
	}

	age := time.Since(info.ModTime())
	return age > cacheMaxAge, nil
}

// Fetch downloads tracker URLs from the given source URLs and saves to cache
func Fetch(sourceURLs []string) error {
	if len(sourceURLs) == 0 {
		return fmt.Errorf("no source URLs provided")
	}

	// Shuffle URLs for random order
	urls := make([]string, len(sourceURLs))
	copy(urls, sourceURLs)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(urls), func(i, j int) {
		urls[i], urls[j] = urls[j], urls[i]
	})

	// Try each URL until one succeeds
	var lastErr error
	for _, url := range urls {
		content, err := downloadURL(url)
		if err != nil {
			lastErr = err
			continue
		}

		// Save to cache
		cacheFile, err := paths.GetCacheFile()
		if err != nil {
			return err
		}

		if err := os.WriteFile(cacheFile, content, 0644); err != nil {
			return fmt.Errorf("failed to write cache file: %w", err)
		}

		// Update modification time to now
		now := time.Now()
		if err := os.Chtimes(cacheFile, now, now); err != nil {
			return fmt.Errorf("failed to update cache file time: %w", err)
		}

		return nil
	}

	return fmt.Errorf("failed to download from all mirrors: %w", lastErr)
}

// downloadURL downloads content from a URL
func downloadURL(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d from %s", resp.StatusCode, url)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response from %s: %w", url, err)
	}

	return content, nil
}

// LoadCache loads the cached tracker URLs
func LoadCache() ([]string, error) {
	cacheFile, err := paths.GetCacheFile()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("cache file does not exist")
		}
		return nil, fmt.Errorf("failed to read cache file: %w", err)
	}

	// Split by lines and filter empty lines
	lines := []string{}
	for _, line := range splitLines(string(data)) {
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines, nil
}

// splitLines splits a string by newlines, handling both \n and \r\n
func splitLines(s string) []string {
	var lines []string
	var current []rune
	for _, r := range s {
		if r == '\n' {
			lines = append(lines, string(current))
			current = []rune{}
		} else if r != '\r' {
			current = append(current, r)
		}
	}
	if len(current) > 0 {
		lines = append(lines, string(current))
	}
	return lines
}

