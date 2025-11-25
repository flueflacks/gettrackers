package paths

import (
	"os"
	"path/filepath"
)

// GetConfigDir returns the XDG config directory for gettrackers
func GetConfigDir() (string, error) {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configHome = filepath.Join(home, ".config")
	}
	configDir := filepath.Join(configHome, "gettrackers")
	return configDir, os.MkdirAll(configDir, 0755)
}

// GetCacheDir returns the XDG cache directory for gettrackers
func GetCacheDir() (string, error) {
	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		cacheHome = filepath.Join(home, ".cache")
	}
	cacheDir := filepath.Join(cacheHome, "gettrackers")
	return cacheDir, os.MkdirAll(cacheDir, 0755)
}

// GetConfigFile returns the path to the config file
func GetConfigFile() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.yaml"), nil
}

// GetBlocklistFile returns the path to the blocklist file
func GetBlocklistFile() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "blocklist.txt"), nil
}

// GetCacheFile returns the path to the cache file
func GetCacheFile() (string, error) {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cacheDir, "sources.txt"), nil
}

