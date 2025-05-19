package cmd

import (
	"fmt"
	"os"

	_ "embed"

	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/ByteTheCookies/cookieclient/internal/logger"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cookieclient",
	Short: "The client cli for CookieFarm",
	Long: `CookieFarm is a exploiter writed by the team ByteTheCookies for CyberChallenge
	competition. This is the client cli for the CookieFarm server for attack the teams with exploits.`, // Da migliorare
	Version: "v1.1.0",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

//go:embed banner.txt
var banner string

func init() {
	rootCmd.PersistentFlags().BoolVarP(&config.Args.Debug, "debug", "D", false, "Enable debug logging")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if config.Args.Debug {
			logPath = logger.Setup("debug")
		} else {
			logPath = logger.Setup("info")
		}
		fmt.Println(banner)
	}
}
