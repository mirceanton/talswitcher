package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	TalosconfigDir string
	Talosconfig    string
	LogLevel       log.Level
	LogFormat      log.Formatter
}

const (
	// Configuration keys
	keyTalosconfigDir = "talosconfig-dir"
	keyTalosconfig    = "talosconfig"
	keyLogLevel       = "log-level"
	keyLogFormat      = "log-format"

	// Default values
	defaultLogLevel  = "info"
	defaultLogFormat = "text"
)

var (
	defaultTalosconfigDir = filepath.Join(os.Getenv("HOME"), ".talos", "configs/")
	defaultTalosconfig    = filepath.Join(os.Getenv("HOME"), ".talos", "config")
)

// Init initializes Viper configuration
func Init() {
	// Enable reading from environment variables
	viper.AutomaticEnv()

	// Replace hyphens with underscores in env vars
	// This allows --talosconfig-dir to map to TALOSCONFIG_DIR
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Set default values
	viper.SetDefault(keyTalosconfigDir, defaultTalosconfigDir)
	viper.SetDefault(keyTalosconfig, defaultTalosconfig)
	viper.SetDefault(keyLogLevel, defaultLogLevel)
	viper.SetDefault(keyLogFormat, defaultLogFormat)
}

// Load returns the current configuration
func Load() (*Config, error) {
	cfg := &Config{}

	// Parse log level
	levelStr := viper.GetString(keyLogLevel)
	level, err := log.ParseLevel(levelStr)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %s", levelStr)
	}
	cfg.LogLevel = level

	// Parse log format
	formatStr := viper.GetString(keyLogFormat)
	switch strings.ToLower(formatStr) {
	case "json":
		cfg.LogFormat = &log.JSONFormatter{}
	case "text":
		cfg.LogFormat = &log.TextFormatter{FullTimestamp: true}
	default:
		return nil, fmt.Errorf("invalid log format: %s", formatStr)
	}

	// Expand and validate talosconfig dir path
	cfg.TalosconfigDir, err = expandPath(viper.GetString(keyTalosconfigDir))
	if err != nil {
		return nil, fmt.Errorf("failed to expand talosconfig directory path: %w", err)
	}
	if err := cfg.validateTalosconfigDir(); err != nil {
		return nil, err
	}

	// Expand and validate talosconfig dir path
	cfg.Talosconfig, err = expandPath(viper.GetString(keyTalosconfig))
	if err != nil {
		return nil, fmt.Errorf("failed to expand talosconfig path: %w", err)
	}
	if err := cfg.validateTalosconfig(); err != nil {
		log.Warnf("Talosconfig file validation failed: %v", err)
		log.Warn("You can still switch contexts, but namespace operations will not be available")
	}

	return cfg, nil
}

// validateTalosconfigDir validates that the talosconfig directory exists and is a directory
func (c *Config) validateTalosconfigDir() error {
	info, err := os.Stat(c.TalosconfigDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("talosconfig directory does not exist: %s", c.TalosconfigDir)
		}
		return fmt.Errorf("failed to stat talosconfig directory: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("talosconfig directory path is not a directory: %s", c.TalosconfigDir)
	}

	return nil
}

// validateTalosconfig validates that the talosconfig file exists and is a file
func (c *Config) validateTalosconfig() error {
	info, err := os.Stat(c.Talosconfig)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("talosconfig file does not exist: %s", c.Talosconfig)
		}
		return fmt.Errorf("failed to stat talosconfig file: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("talosconfig path is a directory, expected a file: %s", c.Talosconfig)
	}

	return nil
}

// expandPath expands a path starting with ~/ to the full home directory path.
func expandPath(path string) (string, error) {
	if !strings.HasPrefix(path, "~/") {
		return path, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(homeDir, path[2:]), nil
}
