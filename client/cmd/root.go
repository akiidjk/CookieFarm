package cmd

import (
	"os"

	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/ByteTheCookies/cookieclient/internal/logger"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "cookieclient",
	Short:   "The client cli for CookieFarm",
	Long:    `CookieFarm is a Attack/Defense CTF framework inspired by DestructiveFarm, developed by the Italian team ByteTheCookies. What sets CookieFarm apart is its hybrid Go+Python architecture and "zero distraction".`,
	Version: "v1.1.0",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&config.Args.Debug, "debug", "D", false, "Enable debug logging")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if config.Args.Debug {
			logger.Setup("debug")
		} else {
			logger.Setup("info")
		}
	}
}
