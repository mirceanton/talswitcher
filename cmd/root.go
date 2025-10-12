package cmd

import (
	"os"

	"github.com/mirceanton/talswitcher/internal/config"
	"github.com/mirceanton/talswitcher/internal/manager"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version       string
	configManager *manager.Manager
)

var rootCmd = &cobra.Command{
	Use:     "talswitcher",
	Short:   "A tool to switch Talos contexts",
	Long:    `talswitcher is a CLI tool to switch Talos contexts from multiple talosconfig files.`,
	Version: version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		// Set up logging
		log.SetLevel(cfg.LogLevel)
		log.SetFormatter(cfg.LogFormat)

		// Create manager with config
		configManager, err = manager.NewManager(cfg.Talosconfig, cfg.TalosconfigDir)
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Initialize Viper
	cobra.OnInitialize(config.Init)

	// Bind flags to Viper
	rootCmd.PersistentFlags().String("talosconfig-dir", "", "Directory containing talosconfig files (env: TALOSCONFIG_DIR)")
	err := viper.BindPFlag("talosconfig-dir", rootCmd.PersistentFlags().Lookup("talosconfig-dir"))
	if err != nil {
		log.Fatalf("Failed to bind flag: %v", err)
	}

	rootCmd.PersistentFlags().String("talosconfig", "", "Currently active talosconfig file (env: TALOSCONFIG)")
	err = viper.BindPFlag("talosconfig", rootCmd.PersistentFlags().Lookup("talosconfig"))
	if err != nil {
		log.Fatalf("Failed to bind flag: %v", err)
	}

	rootCmd.PersistentFlags().String("log-level", "info", "Log level (trace, debug, info, warn, error, fatal, panic) (env: LOG_LEVEL)")
	err = viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	if err != nil {
		log.Fatalf("Failed to bind flag: %v", err)
	}

	rootCmd.PersistentFlags().String("log-format", "text", "Log format (text, json) (env: LOG_FORMAT)")
	err = viper.BindPFlag("log-format", rootCmd.PersistentFlags().Lookup("log-format"))
	if err != nil {
		log.Fatalf("Failed to bind flag: %v", err)
	}
}
