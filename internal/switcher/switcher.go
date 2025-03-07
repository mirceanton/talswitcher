package switcher

import (
	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"

	"github.com/mirceanton/talswitcher/internal/config"
	"github.com/mirceanton/talswitcher/pkg/types"
)

// Switcher handles the context switching logic
type Switcher struct {
	ConfigDir string
}

// NewSwitcher creates a new Switcher instance
func NewSwitcher(configDir string) *Switcher {
	return &Switcher{
		ConfigDir: configDir,
	}
}

// Switch handles the context switching process
func (s *Switcher) Switch(contextArg string) {
	// Get the destination path for talosconfig
	destPath := config.GetDestinationPath()

	// Check if switching to previous context
	if contextArg == types.PreviousContextSymbol {
		config.RestorePreviousContext(destPath)
		return
	}

	// Load all available contexts
	contextMap, contextNames := config.LoadContexts(s.ConfigDir)

	// Determine the target context
	var selectedContext string

	if contextArg != "" {
		// Non-interactive mode: use the provided context name
		selectedContext = contextArg
		log.Debugf("Using context from command line: %s", selectedContext)
		if _, exists := contextMap[selectedContext]; !exists {
			log.Fatalf("Context '%s' not found", selectedContext)
		}
	} else {
		// Interactive mode: show list of contexts
		log.Debugf("Using interactive mode")
		prompt := &survey.Select{
			Message: "Choose a context:",
			Options: contextNames,
		}
		err := survey.AskOne(prompt, &selectedContext)
		if err != nil {
			log.Fatalf("Failed to get user input: %v", err)
		}
	}

	// Create a backup of the current config
	config.BackupCurrentConfig(destPath)

	// Switch to the selected context
	config.SwitchContext(selectedContext, contextMap[selectedContext], destPath)
}
