package cmd

import (
	"fmt"
	"os"

	"github.com/flueflacks/gettrackers/internal/config"
	"github.com/flueflacks/gettrackers/internal/fetch"

	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Force download/update the cached sources file",
	Long:  `Downloads tracker URLs from configured sources and updates the cache file.`,
	RunE:  runFetch,
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}

func runFetch(cmd *cobra.Command, args []string) error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Force fetch
	if err := fetch.Fetch(cfg.SourceURLs); err != nil {
		return fmt.Errorf("failed to fetch: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Successfully updated cache\n")
	return nil
}

