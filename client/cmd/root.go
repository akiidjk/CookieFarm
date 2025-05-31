package cmd

import (
	"os"

	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/ByteTheCookies/cookieclient/internal/logger"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cookieclient",
	Short: "The client cli for CookieFarm",
	Long: `CookieFarm is a exploiter writed by the team ByteTheCookies for CyberChallenge
	competition. This is the client cli for the CookieFarm server for attack the teams with exploits.`, // Da migliorare
	Version: "v1.1.0",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(configCmd)
	RootCmd.PersistentFlags().BoolVarP(&config.ArgsAttackInstance.Debug, "debug", "D", false, "Enable debug logging")

	RootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if config.ArgsAttackInstance.Debug {
			logger.Setup("debug")
		} else {
			logger.Setup("info")
		}
	}
}
