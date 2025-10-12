package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/siderolabs/talos/pkg/machinery/client/config"
	log "github.com/sirupsen/logrus"
)

// Manager handles talosconfig file operations and Talos context switching.
type Manager struct {
	talosconfigPath string
	backupPath      string
	talosconfigDir  string
	contextMap      map[string]string
	contextNames    []string
}

// NewManager creates a new talosconfig Manager instance.
// It validates and loads the config directory if provided.
func NewManager(talosconfigPath, talosconfigDir string) (*Manager, error) {
	m := &Manager{
		talosconfigDir:  talosconfigDir,
		talosconfigPath: talosconfigPath,
		backupPath:      talosconfigPath + ".old",
	}

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
	// Find the talosconfig file containing the specified context
	contextFilePath, exists := m.contextMap[contextName]
	if !exists {
		return fmt.Errorf("context %s not found", contextName)
	}

	// Load the talosconfig file containing the desired context
	talosconfig, err := config.Open(contextFilePath)
	if err != nil {
		return fmt.Errorf("failed to load talosconfig from %s: %w", contextFilePath, err)
	}

	// Update the context in the loaded talosconfig
	talosconfig.Context = contextName

	// Backup the current talosconfig before switching
	if err := m.backup(); err != nil {
		log.Warnf("Failed to save current configuration as previous: %v", err)
	}

	// Save the updated talosconfig to the active path
	if err := talosconfig.Save(m.talosconfigPath); err != nil {
		return fmt.Errorf("failed to write talosconfig: %w", err)
	}

	return nil
}

// Restore swaps the current talosconfig with the previous backup.
func (m *Manager) Restore() error {
	// Read the current talosconfig
	currentConfig, err := os.ReadFile(m.talosconfigPath)
	if err != nil {
		return fmt.Errorf("failed to read current config: %w", err)
	}

	// Read the previous talosconfig backup
	prevConfig, err := os.ReadFile(m.backupPath)
	if err != nil {
		return fmt.Errorf("failed to read previous config: %w", err)
	}

	// Swap the files
	if err := os.WriteFile(m.talosconfigPath, prevConfig, 0o600); err != nil {
		return fmt.Errorf("failed to write current config: %w", err)
	}
	if err := os.WriteFile(m.backupPath, currentConfig, 0o600); err != nil {
		return fmt.Errorf("failed to write previous config: %w", err)
	}

	return nil
}

// ================================================================================================
// Helper functions
// ================================================================================================

// backup backs up the current talosconfig to config.old.
func (m *Manager) backup() error {
	data, err := os.ReadFile(m.talosconfigPath)
	if err != nil {
		return fmt.Errorf("failed to read current talosconfig: %w", err)
	}

	if err := os.WriteFile(m.backupPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write previous talosconfig: %w", err)
	}

	return nil
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
