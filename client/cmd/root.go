package cmd

import (
	"fmt"
	"os"

	_ "embed"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cookieclient",
	Short: "The client cli for CookieFarm",
	Long: `CookieFarm is a exploiter writed by the team ByteTheCookies for CyberChallenge
	competition. This is the client cli for the CookieFarm server for attack the teams with exploits.`, // Da migliorare
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
	fmt.Println(banner)
}
