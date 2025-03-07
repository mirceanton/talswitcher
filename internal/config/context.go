package config

import (
	"github.com/siderolabs/talos/pkg/machinery/client/config"
	log "github.com/sirupsen/logrus"
)

// SwitchContext switches to the specified context
func SwitchContext(contextName string, contextPath string, destPath string) {
	// Load the talosconfig file for the selected context
	log.WithFields(log.Fields{"source": contextPath}).Debugf("Loading talosconfig file for context: %s", contextName)
	talosconfig, err := config.Open(contextPath)
	if err != nil {
		log.WithFields(log.Fields{"source": contextPath}).Fatalf("Failed to parse talosconfig file: %v", err)
	}

	// Update the current context
	log.Debugf("Updating talosconfig with new context: %s", contextName)
	talosconfig.Context = contextName

	// Write the updated talosconfig back to the file
	log.WithFields(log.Fields{"destination": destPath}).Debugf("Writing updated talosconfig file")
	err = talosconfig.Save(destPath)
	if err != nil {
		log.WithFields(log.Fields{"destination": destPath}).Fatalf("Failed to write to destination file: %v", err)
	}
	log.Debugf("Successfully wrote talosconfig file to: %s", destPath)

	log.Infof("Switched to context '%s'", contextName)
}
