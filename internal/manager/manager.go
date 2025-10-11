package manager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/siderolabs/talos/pkg/machinery/client/config"
	log "github.com/sirupsen/logrus"
)

// Manager handles talosconfig file operations and Talos context switching.
type Manager struct {
	talosconfigPath string
	talosconfigDir  string
	contextMap      map[string]string
	contextNames    []string
}

// Custom errors
var (
	ErrNoConfigDir      = errors.New("talosconfig directory not provided, please provide the directory containing talosconfig files via the --talosconfig-dir flag or TALOSCONFIG_DIR environment variable")
	ErrNotADirectory    = errors.New("the provided path is not a directory")
	ErrNoPreviousConfig = errors.New("no previous configuration found")
	ErrContextNotFound  = errors.New("context not found")
)

// NewManager creates a new talosconfig Manager instance.
// It validates and loads the config directory if provided.
func NewManager(configDir string) (*Manager, error) {
	m := &Manager{
		talosconfigPath: getTalosconfigPath(),
	}

	validatedDir, err := m.validateConfigDir(configDir)
	if err != nil {
		return nil, err
	}
	m.talosconfigDir = validatedDir

	// Load contexts from directory
	if err := m.loadContexts(); err != nil {
		return nil, fmt.Errorf("failed to load contexts: %w", err)
	}

	return m, nil
}

// GetAllContexts returns the available context names.
func (m *Manager) GetAllContexts() []string {
	return m.contextNames
}

// SwitchToContext switches to the specified Talos context.
func (m *Manager) SwitchToContext(contextName string) error {
	contextFilePath, exists := m.contextMap[contextName]
	if !exists {
		return fmt.Errorf("%w: %s", ErrContextNotFound, contextName)
	}

	if err := m.backup(); err != nil {
		log.Warnf("Failed to save current configuration as previous: %v", err)
	}

	talosconfig, err := config.Open(contextFilePath)
	if err != nil {
		return fmt.Errorf("failed to load talosconfig from %s: %w", contextFilePath, err)
	}

	talosconfig.Context = contextName

	if err := talosconfig.Save(m.talosconfigPath); err != nil {
		return fmt.Errorf("failed to write talosconfig: %w", err)
	}

	return nil
}

// Restore swaps the current talosconfig with the previous backup.
func (m *Manager) Restore() error {
	prevPath := m.getPreviousPath()

	if _, err := os.Stat(prevPath); os.IsNotExist(err) {
		return ErrNoPreviousConfig
	}

	currentConfig, err := os.ReadFile(m.talosconfigPath)
	if err != nil {
		return fmt.Errorf("failed to read current config: %w", err)
	}

	prevConfig, err := os.ReadFile(prevPath)
	if err != nil {
		return fmt.Errorf("failed to read previous config: %w", err)
	}

	if err := os.WriteFile(m.talosconfigPath, prevConfig, 0o600); err != nil {
		return fmt.Errorf("failed to write current config: %w", err)
	}

	if err := os.WriteFile(prevPath, currentConfig, 0o600); err != nil {
		return fmt.Errorf("failed to write previous config: %w", err)
	}

	return nil
}

// ================================================================================================
// Helper functions
// ================================================================================================
// getTalosconfigPath returns the path to the current talosconfig file.
func getTalosconfigPath() string {
	path := os.Getenv("TALOSCONFIG")
	if path == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Failed to determine home directory: %v", err)
		}
		path = filepath.Join(homeDir, ".talos", "config")
	}

	expandedPath, err := expandPath(path)
	if err != nil {
		log.Fatalf("Failed to expand TALOSCONFIG path: %v", err)
	}

	// Ensure the destination directory exists
	destDir := filepath.Dir(expandedPath)
	err = os.MkdirAll(destDir, 0o755)
	if err != nil {
		log.Fatalf("Failed to create directory %s: %v", destDir, err)
	}

	return expandedPath
}

// getPreviousPath returns the path where the previous talosconfig backup is stored.
func (m *Manager) getPreviousPath() string {
	return m.talosconfigPath + ".old"
}

// backup backs up the current talosconfig to config.old.
func (m *Manager) backup() error {
	if _, err := os.Stat(m.talosconfigPath); os.IsNotExist(err) {
		return fmt.Errorf("talosconfig file does not exist: %w", err)
	}

	data, err := os.ReadFile(m.talosconfigPath)
	if err != nil {
		return fmt.Errorf("failed to read current talosconfig: %w", err)
	}

	if err := os.WriteFile(m.getPreviousPath(), data, 0o600); err != nil {
		return fmt.Errorf("failed to write previous talosconfig: %w", err)
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

// validateConfigDir validates and expands the provided config directory path.
func (m *Manager) validateConfigDir(configDir string) (string, error) {
	if configDir == "" {
		configDir = os.Getenv("TALOSCONFIG_DIR")
		if configDir == "" {
			return "", ErrNoConfigDir
		}
	}

	expandedPath, err := expandPath(configDir)
	if err != nil {
		return "", err
	}

	info, err := os.Stat(expandedPath)
	if err != nil {
		return "", fmt.Errorf("failed to stat config directory: %w", err)
	}

	if !info.IsDir() {
		return "", ErrNotADirectory
	}

	return expandedPath, nil
}

// loadContexts scans the config directory for talosconfig files and loads all available contexts.
func (m *Manager) loadContexts() error {
	m.contextMap = make(map[string]string)
	m.contextNames = nil

	files, err := os.ReadDir(m.talosconfigDir)
	if err != nil {
		return fmt.Errorf("failed to read config directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		path := filepath.Join(m.talosconfigDir, file.Name())
		talosconfig, err := config.Open(path)
		if err != nil {
			log.WithField("file", file.Name()).Warnf("Failed to parse talosconfig file: %v", err)
			continue
		}

		for contextName := range talosconfig.Contexts {
			if existingPath, exists := m.contextMap[contextName]; exists {
				log.Warnf("Duplicate context name '%s' found in files:\n  - %s\n  - %s",
					contextName, existingPath, path)
				continue
			}
			m.contextMap[contextName] = path
			m.contextNames = append(m.contextNames, contextName)
		}
	}

	return nil
}
