package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/siderolabs/talos/pkg/machinery/client/config"
	log "github.com/sirupsen/logrus"
)

// GetConfigDirectory resolves and validates the talosconfig directory
func GetConfigDirectory(configDir string) string {
	if configDir == "" {
		configDir = os.Getenv("TALOSCONFIG_DIR")
		if configDir == "" {
			log.Fatal("talosconfig directory not provided.")
			log.Fatal("Please provide the directory containing talosconfig files via the --config-dir flag or TALOSCONFIG_DIR environment variable")
		}
	}

	// Expand tilde to user's home directory
	if strings.HasPrefix(configDir, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Failed to determine home directory: %v", err)
		}
		configDir = filepath.Join(homeDir, configDir[2:])
	}

	log.Debugf("Using talosconfig directory: %s", configDir)
	return configDir
}

// GetDestinationPath returns the path where the talosconfig should be written
func GetDestinationPath() string {
	destPath := os.Getenv("TALOSCONFIG")
	if destPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Failed to determine home directory: %v", err)
		}
		destPath = filepath.Join(homeDir, ".talos", "config")
		log.Debugf("Using default talosconfig path: %s", destPath)
	} else {
		log.Debugf("Using talosconfig path from TALOSCONFIG env var: %s", destPath)
	}

	// Ensure the destination directory exists
	destDir := filepath.Dir(destPath)
	err := os.MkdirAll(destDir, 0o755)
	if err != nil {
		log.Fatalf("Failed to create directory %s: %v", destDir, err)
	}

	return destPath
}

// LoadContexts loads all available contexts from the config directory
func LoadContexts(configDir string) (map[string]string, []string) {
	// Get all files in the config directory
	files, err := os.ReadDir(configDir)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}
	log.Debugf("Found %d files in directory: %s", len(files), configDir)

	// Parse all yaml files in the directory
	contextMap := make(map[string]string) // Map storing the context name (key) and the corresponding talosconfig file (value)
	var contextNames []string             // List of context names for the interactive prompt

	for _, file := range files {
		// Skip directories
		if file.IsDir() {
			log.Debugf("Skipping directory: %s", file.Name())
			continue
		}

		// Parse the talosconfig file
		path := filepath.Join(configDir, file.Name())
		log.Debugf("Parsing yaml file: %s", path)
		talosconfig, err := config.Open(path)
		if err != nil {
			log.WithFields(log.Fields{"file": file.Name()}).Warnf("Failed to parse talosconfig file: %v", err)
			continue
		}

		// Add the context to the map
		for contextName := range talosconfig.Contexts {
			if _, exists := contextMap[contextName]; exists {
				log.Fatalf("Duplicate context name '%s' found in files:\n- %s\n- %s", contextName, contextMap[contextName], path)
			}
			log.Debugf("Found context: %s", contextName)
			contextMap[contextName] = path
			contextNames = append(contextNames, contextName)
		}
	}

	// Check if any contexts were found
	if len(contextMap) == 0 {
		log.Fatal("No Talos contexts found in the provided directory: ", configDir)
	}
	log.Debugf("Found %d unique contexts", len(contextMap))

	return contextMap, contextNames
}
