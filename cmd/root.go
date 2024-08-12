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
var logLevel string  // Global flag for the log level
var logFormat string // Global flag for the log format
var version string   // Version of the tool

func setConfigDirectory() {
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
}

func setLogLevel() {
	if logLevel == "" {
		logLevel = os.Getenv("TALSWITCHER_LOG_LEVEL")
		if logLevel == "" {
			logLevel = "info"
		}
	}
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatalf("Invalid log level: %v", err)
	}
	log.SetLevel(level)
}

func setLogFormat() {
	if logFormat == "" {
		logFormat = os.Getenv("TALSWITCHER_LOG_FORMAT")
		if logFormat == "" {
			logFormat = "text"
		}
	}
	switch logFormat {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	case "text":
		log.SetFormatter(&log.TextFormatter{})

	default:
		log.Fatalf("Invalid log format: %s", logFormat)
	}
}

var rootCmd = &cobra.Command{
	Use:     "talswitcher",
	Short:   "Select or specify a Talos context to switch to.",
	Long:    `talswitcher is a CLI tool to switch between Talos contexts from multiple talosconfig files.`,
	Version: version,
	Args:    cobra.MaximumNArgs(1), // Accept at most one argument (the context name)
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		setLogLevel()
		setLogFormat()
		log.Debugf("Using log level: %s", log.GetLevel())
		log.Debugf("Using log format: %s", logFormat)

		setConfigDirectory()
		log.Debugf("Using talosconfig directory: %s", configDir)
	},
	Run: func(cmd *cobra.Command, args []string) {
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

			// Skip files that are not YAML
			if filepath.Ext(file.Name()) != ".yaml" && filepath.Ext(file.Name()) != ".yml" {
				log.Debugf("Skipping non-YAML file: %s", file.Name())
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

		// Determine the target context
		var selectedContext string
		if len(args) == 1 {
			// Non-interactive mode: use the provided cluster name
			selectedContext = args[0]
			log.Debugf("Using context from command line: %s", selectedContext)
			if _, exists := contextMap[selectedContext]; !exists {
				log.Fatalf("Context '%s' not found", selectedContext)
			}
		} else {
			// Interactive mode: show list of clusters
			log.Debugf("Using interactive mode")
			prompt := &survey.Select{
				Message: "Choose a context:",
				Options: contextNames,
			}
			err = survey.AskOne(prompt, &selectedContext)
			if err != nil {
				log.Fatalf("Failed to get user input: %v", err)
			}
		}
		log.Debugf("Selected context: %s", selectedContext)

		// Determine the target location for copying the file
		destPath := os.Getenv("TALOSCONFIG")
		if destPath == "" {
			log.Debugf("TALOSCONFIG environment variable not set. Using default location")
			homeDir, err := os.UserHomeDir()
			if err != nil {
				log.Fatalf("Failed to determine home directory: %v", err)
			}
			log.Debugf("Using home directory: %s", homeDir)
			destPath = filepath.Join(homeDir, ".talos", "config")
			log.Debugf("Using default destination: %s", destPath)
		}

		// Ensure the destination directory exists
		destDir := filepath.Dir(destPath)
		log.Debugf("Ensuring destination directory exists: %s", destDir)
		err = os.MkdirAll(destDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", destDir, err)
		}

		// Load the talosconfig file for the selected context
		log.WithFields(log.Fields{"source": contextMap[selectedContext]}).Debugf("Loading talosconfig file for context: %s", selectedContext)
		talosconfig, err := config.Open(contextMap[selectedContext])
		if err != nil {
			log.WithFields(log.Fields{"source": contextMap[selectedContext]}).Fatalf("Failed to parse talosconfig file: %v", err)
		}

		// Update the current context
		log.Debugf("Updating talosconfig with new context: %s", selectedContext)
		talosconfig.Context = selectedContext

		// Write the updated talosconfig back to the file
		log.WithFields(log.Fields{"destination": destPath}).Debugf("Writing updated talosconfig file")
		err = talosconfig.Save(destPath)
		if err != nil {
			log.WithFields(log.Fields{"destination": destPath}).Fatalf("Failed to write to destination file: %v", err)
		}
		log.Debugf("Successfully wrote talosconfig file to: %s", destPath)

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
	rootCmd.PersistentFlags().StringVarP(&configDir, "talosconfig-dir", "", "", "Directory containing talosconfig files")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "", "", "The logging level. Acceptable values are panic, fatal, error, warn, info, debug, trace.")
	rootCmd.PersistentFlags().StringVarP(&logFormat, "log-format", "", "", "The log output format. Acceptable values are json and text.")
}
