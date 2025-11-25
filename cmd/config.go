package cmd

import (
	"fmt"
	"os"
	"strings"

	"gettrackers/internal/config"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `Get or set configuration values.`,
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Show config values",
	Long:  `Show all config values if no key is specified, or a specific value if key is given.`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runConfigGet,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a config value",
	Long:  `Set a configuration value. For source_urls, provide comma-separated URLs.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runConfigSet,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if len(args) == 0 {
		// Show all config
		fmt.Printf("source_urls:\n")
		for _, url := range cfg.SourceURLs {
			fmt.Printf("  - %s\n", url)
		}
		return nil
	}

	// Show specific key
	key := args[0]
	value, err := cfg.Get(key)
	if err != nil {
		return err
	}

	switch v := value.(type) {
	case []string:
		for _, item := range v {
			fmt.Println(item)
		}
	default:
		fmt.Println(value)
	}

	return nil
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	key := args[0]
	valueStr := args[1]

	var value interface{}
	switch key {
	case "source_urls":
		// Split by comma and trim spaces
		urls := strings.Split(valueStr, ",")
		trimmed := make([]string, 0, len(urls))
		for _, url := range urls {
			trimmed = append(trimmed, strings.TrimSpace(url))
		}
		value = trimmed
	default:
		value = valueStr
	}

	if err := cfg.Set(key, value); err != nil {
		return err
	}

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Set %s = %v\n", key, value)
	return nil
}

