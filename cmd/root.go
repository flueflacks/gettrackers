package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gettrackers",
	Short: "Download, filter, and output tracker URLs grouped by domain",
	Long:  `gettrackers is a CLI tool that downloads tracker URLs from configurable sources, filters them using a blocklist, and outputs them grouped by domain.`,
	RunE:  runGroups,
}

var outputFile string

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "Write output to file instead of stdout")
	rootCmd.Flags().IntVarP(&startPriority, "start-priority", "p", 0, "Output N blank lines before tracker groups (default: 0)")
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

// getOutputWriter returns the appropriate writer for output
func getOutputWriter() (*os.File, error) {
	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create output file: %w", err)
		}
		return file, nil
	}
	return os.Stdout, nil
}

