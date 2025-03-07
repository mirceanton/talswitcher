package cmd

import (
	"os"

	"github.com/spf13/cobra"

	applog "github.com/mirceanton/talswitcher/internal/log"
	"github.com/mirceanton/talswitcher/pkg/types"
)

var (
	configDir string // Global flag for the talosconfig directory
	logLevel  string // Global flag for the log level
	logFormat string // Global flag for the log format
	version   string // Version of the tool
)

const exampleString = `  # Switch to a specific context
  talswitcher context my-context

  # Switch to previous context
  talswitcher ctx -

  # Interactive mode (no argument)
  talswitcher context

  # Generate shell completions
  talswitcher completion bash > ~/.bash_completion.d/talswitcher`

var rootCmd = &cobra.Command{
	Use:     "talswitcher",
	Short:   "CLI tool to switch between Talos contexts",
	Long:    `talswitcher is a CLI tool to switch between Talos contexts from multiple talosconfig files.`,
	Version: version,
	Example: exampleString,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		applog.Setup(types.Config{
			LogLevel:  logLevel,
			LogFormat: logFormat,
		})
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
