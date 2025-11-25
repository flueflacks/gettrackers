package cmd

import (
	"fmt"
	"os"

	"gettrackers/internal/filter"
	"gettrackers/internal/paths"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show cached sources or blocklist",
	Long:  `Display the cached source URLs or blocklist entries.`,
}

var showSourcesCmd = &cobra.Command{
	Use:   "sources",
	Short: "Display the cached source URLs",
	RunE:  runShowSources,
}

var showBlocklistCmd = &cobra.Command{
	Use:   "blocklist",
	Short: "Display blocklist entries",
	RunE:  runShowBlocklist,
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.AddCommand(showSourcesCmd)
	showCmd.AddCommand(showBlocklistCmd)
}

func runShowSources(cmd *cobra.Command, args []string) error {
	// Read and output the raw cache file content (unfiltered)
	cacheFile, err := paths.GetCacheFile()
	if err != nil {
		return fmt.Errorf("failed to get cache file path: %w", err)
	}

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("cache file does not exist")
		}
		return fmt.Errorf("failed to read cache file: %w", err)
	}

	// Output the raw cache file content
	fmt.Print(string(data))
	return nil
}

func runShowBlocklist(cmd *cobra.Command, args []string) error {
	blocklist, err := filter.LoadBlocklist()
	if err != nil {
		return fmt.Errorf("failed to load blocklist: %w", err)
	}

	if len(blocklist) == 0 {
		fmt.Fprintf(os.Stderr, "Blocklist is empty\n")
		return nil
	}

	for _, pattern := range blocklist {
		fmt.Println(pattern)
	}

	return nil
}

