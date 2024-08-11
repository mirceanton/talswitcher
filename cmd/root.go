package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var configDir string // Global flag for the talosconfig directory

var rootCmd = &cobra.Command{
	Use:   "talswitcher",
	Short: "A tool to switch between Talos contexts",
	Long:  `talswitcher is a CLI tool to switch between Talos contexts from multiple talosconfig files.`,
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
