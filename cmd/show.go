package cmd

import (
	"fmt"
	"os"

	"gettrackers/internal/fetch"
	"gettrackers/internal/filter"

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
	urls, err := fetch.LoadCache()
	if err != nil {
		return fmt.Errorf("failed to load cache: %w", err)
	}

	for _, url := range urls {
		fmt.Println(url)
	}

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

