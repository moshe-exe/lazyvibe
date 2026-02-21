// Package config provides configuration management for lazyvibe.
package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config represents the application configuration.
type Config struct {
	Theme            string `toml:"theme"`
	RefreshInterval  int    `toml:"refresh_interval"` // seconds
	DefaultTimeRange string `toml:"default_time_range"`
	ShowScrollbar    bool   `toml:"show_scrollbar"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		Theme:            "default",
		RefreshInterval:  10,
		DefaultTimeRange: "all",
		ShowScrollbar:    true,
	}
}

// configPath returns the path to the config file.
func configPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "lazyvibe", "config.toml")
}

// Load loads the configuration from file, creating defaults if needed.
func Load() (*Config, error) {
	cfg := DefaultConfig()
	path := configPath()
	if path == "" {
		return cfg, nil
	}

	// Check if config file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create default config file
		if err := cfg.Save(); err != nil {
			// Ignore save error, just use defaults
			return cfg, nil
		}
		return cfg, nil
	}

	// Load existing config
	if _, err := toml.DecodeFile(path, cfg); err != nil {
		return DefaultConfig(), err
	}

	return cfg, nil
}

// Save saves the configuration to file.
func (c *Config) Save() error {
	path := configPath()
	if path == "" {
		return nil
	}

	// Create directory if needed
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write config file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	return encoder.Encode(c)
}
