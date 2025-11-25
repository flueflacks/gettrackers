package cmd

import (
	"fmt"
	"os"

	"gettrackers/internal/filter"

	"github.com/spf13/cobra"
)

var blockCmd = &cobra.Command{
	Use:   "block <pattern>",
	Short: "Add a pattern to the blocklist",
	Long:  `Adds a pattern to the blocklist. If the pattern already exists, it is skipped.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runBlock,
}

func init() {
	rootCmd.AddCommand(blockCmd)
}

func runBlock(cmd *cobra.Command, args []string) error {
	pattern := args[0]

	if err := filter.AddToBlocklist(pattern); err != nil {
		return fmt.Errorf("failed to add to blocklist: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Added pattern to blocklist: %s\n", pattern)
	return nil
}

