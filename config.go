package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the monitor configuration
type Config struct {
	Files     []string `json:"files"`
	IntervalS int      `json:"interval_s"`
}

// LoadConfig reads and validates the configuration file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config.json: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config.json: %w", err)
	}

	// Validate configuration
	if len(config.Files) == 0 {
		return nil, fmt.Errorf("empty file list")
	}

	if config.IntervalS <= 0 {
		config.IntervalS = 5 // default 5 seconds
	}

	// Convert relative paths to absolute
	for i, file := range config.Files {
		absPath, err := filepath.Abs(file)
		if err != nil {
			return nil, fmt.Errorf("invalid path %q: %w", file, err)
		}
		config.Files[i] = absPath
	}

	return &config, nil
}
