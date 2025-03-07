package config

import (
	"os"

	"github.com/siderolabs/talos/pkg/machinery/client/config"
	log "github.com/sirupsen/logrus"
)

// BackupCurrentConfig creates a backup of the current talosconfig file
func BackupCurrentConfig(configPath string) {
	backupPath := configPath + ".old"

	// Check if the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Debugf("No existing config file to backup at: %s", configPath)
		return
	}

	// Create a backup by copying the file
	configData, err := os.ReadFile(configPath)
	if err != nil {
		log.Warnf("Failed to read current config for backup: %v", err)
		return
	}
	err = os.WriteFile(backupPath, configData, 0o600)
	if err != nil {
		log.Warnf("Failed to create backup of current config: %v", err)
		return
	}

	log.Debugf("Created backup of current config at: %s", backupPath)
}

// RestorePreviousContext restores the previous talosconfig from backup
func RestorePreviousContext(configPath string) {
	backupPath := configPath + ".old"

	// Check if the backup file exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		log.Fatal("No previous context found. Backup file does not exist.")
		return
	}

	// Create a backup of the current config before restoring the previous one
	// This allows toggling back and forth between contexts
	currentData, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read current config: %v", err)
		return
	}

	// Read the backup file
	backupData, err := os.ReadFile(backupPath)
	if err != nil {
		log.Fatalf("Failed to read backup config: %v", err)
		return
	}

	// Write the backup data to the current config
	err = os.WriteFile(configPath, backupData, 0o600)
	if err != nil {
		log.Fatalf("Failed to restore previous config: %v", err)
		return
	}

	// Update the backup with the current config (to enable toggling)
	err = os.WriteFile(backupPath, currentData, 0o600)
	if err != nil {
		log.Warnf("Failed to create new backup of previous config: %v", err)
	}

	// Parse the restored config to get the context name for logging
	talosconfig, err := config.Open(configPath)
	if err != nil {
		log.Infof("Switched to previous context")
		return
	}

	log.Infof("Switched to previous context '%s'", talosconfig.Context)
}
