package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mirceanton/talswitcher/internal/config"
	applog "github.com/mirceanton/talswitcher/internal/log"
	"github.com/mirceanton/talswitcher/internal/switcher"
	"github.com/mirceanton/talswitcher/pkg/types"
)

var contextCmd = &cobra.Command{
	Use:     "context [context name]",
	Aliases: []string{"ctx"},
	Short:   "Switch to a specified Talos context",
	Long:    `Switch to a specified Talos context or select one interactively.`,
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resolvedConfigDir := config.GetConfigDirectory(configDir)
		// Create a new switcher
		s := switcher.NewSwitcher(resolvedConfigDir)

		// Get the context argument (if any)
		contextArg := ""
		if len(args) == 1 {
			contextArg = args[0]
		}

		// Execute the switch
		s.Switch(contextArg)
	},
	ValidArgsFunction: contextCompletion,
}

// contextCompletion provides autocompletion for available Talos contexts
func contextCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// If we already have an argument, don't provide more completions
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	// Get config directory - handle errors gracefully
	resolvedConfigDir := safeGetConfigDir()
	if resolvedConfigDir == "" {
		// Return just the previous context symbol if we can't determine config dir
		return []string{types.PreviousContextSymbol}, cobra.ShellCompDirectiveNoFileComp
	}

	// Try to load all available contexts, but don't crash completion if it fails
	contextMap, contextNames := safeLoadContexts(resolvedConfigDir)

	// If no contexts found, just return the previous context symbol
	if len(contextMap) == 0 {
		return []string{types.PreviousContextSymbol}, cobra.ShellCompDirectiveNoFileComp
	}

	// Add previous context symbol
	contexts := append([]string{types.PreviousContextSymbol}, contextNames...)

	return contexts, cobra.ShellCompDirectiveNoFileComp
}

// safeGetConfigDir tries to get the config directory without failing
func safeGetConfigDir() string {
	defer func() {
		// Recover from any panic during directory resolution
		if r := recover(); r != nil {
			return
		}
	}()

	// Override logging for completions to prevent output
	applog.Setup(types.Config{
		LogLevel:  "error",
		LogFormat: "text",
	})

	return config.GetConfigDirectory(configDir)
}

// safeLoadContexts attempts to load contexts but handles errors gracefully
func safeLoadContexts(configDir string) (map[string]string, []string) {
	defer func() {
		// Recover from any panic during context loading
		if r := recover(); r != nil {
			return
		}
	}()

	// Try to safely load contexts
	return config.LoadContexts(configDir)
}

func init() {
	rootCmd.AddCommand(contextCmd)
}
