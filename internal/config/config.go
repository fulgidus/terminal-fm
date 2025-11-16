// Package config provides application configuration management.
package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds all application configuration.
type Config struct {
	Server  ServerConfig
	Player  PlayerConfig
	Storage StorageConfig
	I18n    I18nConfig
	DevMode bool
}

// ServerConfig contains SSH server settings.
type ServerConfig struct {
	Host        string
	Port        int
	KeyPath     string
	MaxTimeout  int // seconds
	IdleTimeout int // seconds
}

// PlayerConfig contains audio player settings.
type PlayerConfig struct {
	DefaultPlayer string // "ffplay" or "mpv"
	FFplayPath    string
	MpvPath       string
	BufferSize    int // seconds
	MaxRetries    int
}

// StorageConfig contains database settings.
type StorageConfig struct {
	DBPath     string
	BackupPath string
}

// I18nConfig contains internationalization settings.
type I18nConfig struct {
	DefaultLocale string
	LocalesPath   string
}

// New creates a new Config with default values.
func New() *Config {
	homeDir, _ := os.UserHomeDir()
	dataDir := filepath.Join(homeDir, ".terminal-fm")

	return &Config{
		Server: ServerConfig{
			Host:        "0.0.0.0",
			Port:        22,
			KeyPath:     filepath.Join(dataDir, "ssh_host_key"),
			MaxTimeout:  3600, // 1 hour
			IdleTimeout: 1800, // 30 minutes
		},
		Player: PlayerConfig{
			DefaultPlayer: "ffplay",
			FFplayPath:    "ffplay",
			MpvPath:       "mpv",
			BufferSize:    5,
			MaxRetries:    3,
		},
		Storage: StorageConfig{
			DBPath:     filepath.Join(dataDir, "terminal-fm.db"),
			BackupPath: filepath.Join(dataDir, "backups"),
		},
		I18n: I18nConfig{
			DefaultLocale: "en",
			LocalesPath:   "pkg/i18n/locales",
		},
		DevMode: false,
	}
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Player.DefaultPlayer != "ffplay" && c.Player.DefaultPlayer != "mpv" {
		return fmt.Errorf("invalid player: %s (must be 'ffplay' or 'mpv')", c.Player.DefaultPlayer)
	}

	if c.I18n.DefaultLocale != "en" && c.I18n.DefaultLocale != "it" {
		return fmt.Errorf("unsupported locale: %s (must be 'en' or 'it')", c.I18n.DefaultLocale)
	}

	return nil
}

// EnsureDataDir creates the data directory if it doesn't exist.
func (c *Config) EnsureDataDir() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	dataDir := filepath.Join(homeDir, ".terminal-fm")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Create backups directory
	if err := os.MkdirAll(c.Storage.BackupPath, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	return nil
}
