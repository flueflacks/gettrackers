package config

import (
	"fmt"
	"os"

	"gettrackers/internal/paths"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	SourceURLs []string `yaml:"source_urls"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		SourceURLs: []string{
			"https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_all.txt",
			"https://ngosang.github.io/trackerslist/trackers_all.txt",
			"https://cdn.jsdelivr.net/gh/ngosang/trackerslist@master/trackers_all.txt",
		},
	}
}

// Load loads the configuration from file, or returns default if file doesn't exist
func Load() (*Config, error) {
	configFile, err := paths.GetConfigFile()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	configFile, err := paths.GetConfigFile()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Get returns a config value by key
func (c *Config) Get(key string) (interface{}, error) {
	switch key {
	case "source_urls":
		return c.SourceURLs, nil
	default:
		return nil, fmt.Errorf("unknown config key: %s", key)
	}
}

// Set sets a config value by key
func (c *Config) Set(key string, value interface{}) error {
	switch key {
	case "source_urls":
		urls, ok := value.([]string)
		if !ok {
			return fmt.Errorf("source_urls must be a list of strings")
		}
		c.SourceURLs = urls
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}
	return nil
}

