package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/flueflacks/gettrackers/internal/config"
	"github.com/flueflacks/gettrackers/internal/fetch"
	"github.com/flueflacks/gettrackers/internal/filter"
	"github.com/flueflacks/gettrackers/internal/group"

	"github.com/spf13/cobra"
)

var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "Output grouped tracker URLs (default command)",
	Long:  `Downloads, filters, and outputs tracker URLs grouped by domain in random order.`,
	RunE:  runGroups,
}

var startPriority int

func init() {
	rootCmd.AddCommand(groupsCmd)
	groupsCmd.Flags().IntVarP(&startPriority, "start-priority", "p", 0, "Output N blank lines before tracker groups (default: 0)")
}

func runGroups(cmd *cobra.Command, args []string) error {
	// Validate start-priority flag
	if startPriority < 0 {
		return fmt.Errorf("start-priority must be a non-negative integer, got: %d", startPriority)
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if we need to fetch
	shouldFetch, err := fetch.ShouldFetch(false)
	if err != nil {
		return fmt.Errorf("failed to check cache: %w", err)
	}

	if shouldFetch {
		if err := fetch.Fetch(cfg.SourceURLs); err != nil {
			// Try to use stale cache if available
			fmt.Fprintf(os.Stderr, "Warning: failed to fetch new data: %v\n", err)
			fmt.Fprintf(os.Stderr, "Attempting to use stale cache...\n")
		}
	}

	// Load cache
	urls, err := fetch.LoadCache()
	if err != nil {
		return fmt.Errorf("failed to load cache: %w", err)
	}

	// Load blocklist
	blocklist, err := filter.LoadBlocklist()
	if err != nil {
		return fmt.Errorf("failed to load blocklist: %w", err)
	}

	// Filter URLs
	filtered := filter.Filter(urls, blocklist)

	if len(filtered) == 0 {
		if len(urls) > 0 {
			fmt.Fprintf(os.Stderr, "all urls blocked\n")
		}
		return nil
	}

	// Group by domain
	groups := group.GroupByDomain(filtered)

	// Get output writer
	writer, err := getOutputWriter()
	if err != nil {
		return err
	}
	defer func() {
		if writer != os.Stdout {
			writer.Close()
		}
	}()

	// Output blank lines for start-priority
	if err := outputStartPriority(writer, startPriority); err != nil {
		return err
	}

	// Output groups
	return outputGroups(writer, groups)
}

func outputStartPriority(writer io.Writer, count int) error {
	for i := 0; i < count; i++ {
		if _, err := fmt.Fprintln(writer); err != nil {
			return err
		}
	}
	return nil
}

func outputGroups(writer io.Writer, groups []group.Group) error {
	for i, g := range groups {
		if i > 0 {
			if _, err := fmt.Fprintln(writer); err != nil {
				return err
			}
		}
		for _, url := range g.URLs {
			if _, err := fmt.Fprintln(writer, url); err != nil {
				return err
			}
		}
	}
	return nil
}

