package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var contextCmd = &cobra.Command{
	Use:               "context",
	Aliases:           []string{"ctx"},
	Short:             "Switch the active Talos context",
	ValidArgsFunction: getContextCompletions,
	Args:              cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		contextNames := configManager.GetAllContexts()
		if len(contextNames) == 0 {
			log.Fatal("No Talos contexts found in the provided directory")
		}

		var selectedContext string
		if len(args) == 1 {
			if args[0] == "-" {
				if err := configManager.Restore(); err != nil {
					log.Fatalf("Failed to switch to previous config: %v", err)
				}
				return
			}
			selectedContext = args[0]
		} else {
			prompt := &survey.Select{
				Message: "Choose a context:",
				Options: contextNames,
			}
			if err := survey.AskOne(prompt, &selectedContext); err != nil {
				log.Fatalf("Failed to get user input: %v", err)
			}
		}

		if err := configManager.SwitchToContext(selectedContext); err != nil {
			log.Fatalf("Failed to switch context: %v", err)
		}

		log.Infof("Switched to context '%s'", selectedContext)
	},
}

func init() {
	rootCmd.AddCommand(contextCmd)
}

func getContextCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return configManager.GetAllContexts(), cobra.ShellCompDirectiveNoFileComp
}
