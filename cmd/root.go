package cmd

import (
	"os"

	"github.com/mirceanton/talswitcher/internal/manager"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	configDir     string
	version       string
	configManager *manager.Manager
)

var rootCmd = &cobra.Command{
	Use:     "talswitcher",
	Short:   "A tool to switch Talos contexts",
	Long:    `talswitcher is a CLI tool to switch Talos contexts from multiple talosconfig files.`,
	Version: version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		configManager, err = manager.NewManager(configDir)
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 && args[0] == "-" {
			if err := configManager.Restore(); err != nil {
				log.Fatalf("Failed to switch to previous config: %v", err)
			}
			return nil
		}
		return cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configDir, "talosconfig-dir", "", "", "Directory containing talosconfig files")
}
