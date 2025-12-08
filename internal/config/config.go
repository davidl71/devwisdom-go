package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds wisdom configuration
type Config struct {
	Source        string `json:"source"`
	HebrewEnabled bool   `json:"hebrew_enabled"`
	HebrewOnly    bool   `json:"hebrew_only"`
	Disabled      bool   `json:"disabled"`
	configPath    string
}

// NewConfig creates a new config with defaults
func NewConfig() *Config {
	return &Config{
		Source:        "pistis_sophia", // Default source
		HebrewEnabled: false,
		HebrewOnly:    false,
		Disabled:      false,
		configPath:    GetConfigPath(),
	}
}

// Load loads configuration from file or environment
func (c *Config) Load() error {
	// First check environment variables
	if source := os.Getenv("EXARP_WISDOM_SOURCE"); source != "" {
		c.Source = source
	}

	if os.Getenv("EXARP_WISDOM_HEBREW") == "1" {
		c.HebrewEnabled = true
	}

	if os.Getenv("EXARP_WISDOM_HEBREW_ONLY") == "1" {
		c.HebrewOnly = true
	}

	if os.Getenv("EXARP_DISABLE_WISDOM") == "1" {
		c.Disabled = true
	}

	// Check for .exarp_no_wisdom marker file
	if _, err := os.Stat(".exarp_no_wisdom"); err == nil {
		c.Disabled = true
	}

	// Try to load from config file
	if _, err := os.Stat(c.configPath); err == nil {
		data, err := os.ReadFile(c.configPath)
		if err == nil {
			if err := json.Unmarshal(data, c); err == nil {
				// Config loaded successfully
				return nil
			}
		}
	}

	return nil // Config file is optional
}

// Save saves configuration to file
func (c *Config) Save() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(c.configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(c.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	// Look for .exarp_wisdom_config in current directory
	// TODO: Also check home directory
	home, err := os.UserHomeDir()
	if err == nil {
		// Check home first
		homeConfig := filepath.Join(home, ".exarp_wisdom_config")
		if _, err := os.Stat(homeConfig); err == nil {
			return homeConfig
		}
	}

	// Default to current directory
	return filepath.Join(".", ".exarp_wisdom_config")
}
