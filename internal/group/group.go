package group

import (
	"fmt"
	"math/rand"
	"net/url"
	"time"
)

// Group represents a group of URLs with the same domain
type Group struct {
	Domain string
	URLs   []string
}

// GroupByDomain groups URLs by their domain
func GroupByDomain(urls []string) []Group {
	groups := make(map[string][]string)
	noDomain := []string{}

	for _, u := range urls {
		domain, err := extractDomain(u)
		if err != nil {
			noDomain = append(noDomain, u)
			continue
		}

		if groups[domain] == nil {
			groups[domain] = []string{}
		}
		groups[domain] = append(groups[domain], u)
	}

	// Convert map to slice
	result := []Group{}
	for domain, urls := range groups {
		result = append(result, Group{
			Domain: domain,
			URLs:   urls,
		})
	}

	// Add "no domain" group if there are any
	if len(noDomain) > 0 {
		result = append(result, Group{
			Domain: "",
			URLs:   noDomain,
		})
	}

	// Shuffle groups in random order
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})

	return result
}

// extractDomain extracts the domain from a URL
func extractDomain(urlStr string) (string, error) {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	host := parsed.Hostname()
	if host == "" {
		return "", fmt.Errorf("no hostname in URL")
	}

	return host, nil
}

