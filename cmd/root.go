package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/mirceanton/talswitcher/internal/config"
	applog "github.com/mirceanton/talswitcher/internal/log"
	"github.com/mirceanton/talswitcher/internal/switcher"
	"github.com/mirceanton/talswitcher/pkg/types"
)

var (
	configDir string // Global flag for the talosconfig directory
	logLevel  string // Global flag for the log level
	logFormat string // Global flag for the log format
	version   string // Version of the tool
)

var rootCmd = &cobra.Command{
	Use:   "talswitcher [context]",
	Short: "Select or specify a Talos context to switch to.",
	Long:  `talswitcher is a CLI tool to switch between Talos contexts from multiple talosconfig files.`,
	Example: `
		# Switch to a specific context
		talswitcher my-context
		
		# Switch to previous context
		talswitcher -
		
		# Interactive mode (no argument)
		talswitcher
	`,
	Version: version,
	Args:    cobra.MaximumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize logger
		applog.Setup(types.Config{
			LogLevel:  logLevel,
			LogFormat: logFormat,
		})

		// Resolve config directory
		configDir = config.GetConfigDirectory(configDir)
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Create a new switcher
		s := switcher.NewSwitcher(configDir)

		// Get the context argument (if any)
		contextArg := ""
		if len(args) == 1 {
			contextArg = args[0]
		}

		// Execute the switch
		s.Switch(contextArg)
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
