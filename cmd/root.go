package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/siderolabs/talos/pkg/machinery/client/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var configDir string // Global flag for the talosconfig directory
var version string   // Version of the tool

var rootCmd = &cobra.Command{
	Use:     "talswitcher",
	Short:   "Select or specify a Talos context to switch to.",
	Long:    `talswitcher is a CLI tool to switch between Talos contexts from multiple talosconfig files.`,
	Version: version,
	Args:    cobra.MaximumNArgs(1), // Accept at most one argument (the context name)
	Run: func(cmd *cobra.Command, args []string) {
		// Determine the taloscofnig directory
		if configDir == "" {
			configDir = os.Getenv("TALOSCONFIG_DIR")
			if configDir == "" {
				log.Fatal("talosconfig directory not provided.")
				log.Fatal("Please provide the directory containing talosconfig files via the --config-dir flag or TALOSCONFIG_DIR environment variable")
			}
		}
		if strings.HasPrefix(configDir, "~/") {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				log.Fatalf("Failed to determine home directory: %v", err)
			}
			configDir = filepath.Join(homeDir, configDir[2:])
		}

		// Get all files in the config directory
		files, err := os.ReadDir(configDir)
		if err != nil {
			log.Fatalf("Failed to read directory: %v", err)
		}

		// Parse all yaml files in the directory
		contextMap := make(map[string]string) // Map storing the context name (key) and the corresponding talosconfig file (value)
		var contextNames []string             // List of context names for the interactive prompt

		for _, file := range files {
			// Skip directories
			if file.IsDir() {
				continue
			}

			// Skip files that are not YAML
			if filepath.Ext(file.Name()) != ".yaml" && filepath.Ext(file.Name()) != ".yml" {
				continue
			}

			// Parse the talosconfig file
			path := filepath.Join(configDir, file.Name())
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
				contextMap[contextName] = path
				contextNames = append(contextNames, contextName)
			}
		}

		// Check if any contexts were found
		if len(contextMap) == 0 {
			log.Fatal("No Talos contexts found in the provided directory: ", configDir)
		}

		// Determine the target context
		var selectedContext string
		if len(args) == 1 {
			// Non-interactive mode: use the provided cluster name
			selectedContext = args[0]
			if _, exists := contextMap[selectedContext]; !exists {
				log.Fatalf("Context '%s' not found", selectedContext)
			}
		} else {
			// Interactive mode: show list of clusters
			prompt := &survey.Select{
				Message: "Choose a context:",
				Options: contextNames,
			}
			err = survey.AskOne(prompt, &selectedContext)
			if err != nil {
				log.Fatalf("Failed to get user input: %v", err)
			}
		}

		// Determine the target location for copying the file
		destPath := os.Getenv("TALOSCONFIG")
		if destPath == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				log.Fatalf("Failed to determine home directory: %v", err)
			}
			destPath = filepath.Join(homeDir, ".talos", "config")
		}

		// Ensure the destination directory exists
		destDir := filepath.Dir(destPath)
		err = os.MkdirAll(destDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", destDir, err)
		}

		// Load the talosconfig file for the selected context
		talosconfig, err := config.Open(contextMap[selectedContext])
		if err != nil {
			log.WithFields(log.Fields{"source": contextMap[selectedContext]}).Fatalf("Failed to parse talosconfig file: %v", err)
		}

		// Update the current context
		talosconfig.Context = selectedContext

		// Write the updated talosconfig back to the file
		err = talosconfig.Save(destPath)
		if err != nil {
			log.WithFields(log.Fields{"destination": destPath}).Fatalf("Failed to write to destination file: %v", err)
		}

		log.Infof("Switched to context '%s'", selectedContext)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Add any global flags here
	rootCmd.PersistentFlags().StringVarP(&configDir, "talosconfig-dir", "", "", "Directory containing talosconfig files")
}
