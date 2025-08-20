package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type StudySessionConfig struct {
	ShowProgress    bool `json:"show_progress"`
	CardsPerSession int  `json:"cards_per_session"`
	NewCardsPerDay  int  `json:"new_cards_per_day"`
}

type Config struct {
	DataDirectory     string             `json:"data_directory"`
	AutoCreateDataDir bool               `json:"auto_create_data_dir"`
	DefaultEaseFactor float64            `json:"default_ease_factor"`
	Theme             string             `json:"theme"`
	BackupEnabled     bool               `json:"backup_enabled"`
	BackupDirectory   string             `json:"backup_directory"`
	StudySession      StudySessionConfig `json:"study_session"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	dataDir := getDefaultDataDir()

	// If we're in development (data directory exists locally), use it
	if _, err := os.Stat("./data"); err == nil {
		dataDir = "./data"
	}

	return &Config{
		DataDirectory:     dataDir,
		AutoCreateDataDir: true,
		DefaultEaseFactor: 2.5,
		Theme:             "default",
		BackupEnabled:     false,
		BackupDirectory:   "",
		StudySession: StudySessionConfig{
			ShowProgress:    true,
			CardsPerSession: 20,
			NewCardsPerDay:  10,
		},
	}
}

// getDefaultDataDir returns the default data directory using XDG specification
func getDefaultDataDir() string {
	// Try XDG_DATA_HOME first
	if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {
		return filepath.Join(xdgDataHome, "anktui")
	}

	// Fallback to ~/.local/share/anktui
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./data" // Last resort fallback
	}

	return filepath.Join(homeDir, ".local", "share", "anktui")
}

// getConfigDir returns the configuration directory using XDG specification
func getConfigDir() string {
	// Try XDG_CONFIG_HOME first
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, "anktui")
	}

	// Fallback to ~/.config/anktui
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./config" // Last resort fallback
	}

	return filepath.Join(homeDir, ".config", "anktui")
}

// GetConfigPath returns the full path to the config file
func GetConfigPath() string {
	return filepath.Join(getConfigDir(), "config.json")
}

// LoadConfig loads the configuration from the config file
func LoadConfig() (*Config, error) {
	configPath := GetConfigPath()

	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	// Read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Parse the JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig saves the configuration to the config file
func (c *Config) SaveConfig() error {
	configDir := getConfigDir()

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Marshal the config to JSON
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	// Write to the config file
	configPath := GetConfigPath()
	return os.WriteFile(configPath, data, 0644)
}

// EnsureDataDir creates the data directory if it doesn't exist and auto-create is enabled
func (c *Config) EnsureDataDir() error {
	if !c.AutoCreateDataDir {
		return nil
	}

	// Expand tilde in path if present
	dataDir := c.DataDirectory
	if len(dataDir) > 0 && dataDir[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		dataDir = filepath.Join(homeDir, dataDir[1:])
	}

	return os.MkdirAll(dataDir, 0755)
}

// GetExpandedDataDir returns the data directory with tilde expansion
func (c *Config) GetExpandedDataDir() (string, error) {
	dataDir := c.DataDirectory
	if len(dataDir) > 0 && dataDir[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dataDir = filepath.Join(homeDir, dataDir[1:])
	}

	return dataDir, nil
}
