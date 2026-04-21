package cmd

import (
	"logger"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ckc",
	Short: "CLI client for interacting with the CookieFarm exploit framework",
	Long:  `CookieFarm is an automated exploitation framework developed by the ByteTheCookies team for the CyberChallenge competition. This CLI client connects to the CookieFarm server to deploy and manage exploits against target teams. To launch the terminal-based user interface (TUI), simply run the command "ckc" without any arguments.`, //nolint:revive
}

func buildCmd(useBanner *bool, useTUI *bool) *cobra.Command {
	var debug bool

	if len(os.Args) != 1 {
		*useTUI = false
	}

	rootCmd.AddCommand(buildConfigCmd())
	rootCmd.AddCommand(buildExploitCmd())

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "Enable debug logging")
	rootCmd.PersistentFlags().BoolVarP(useBanner, "banner", "B", true, "Show banner on startup")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if debug {
			logger.Setup("debug", true)
		} else {
			logger.Setup("info", true)
		}
	}

	return rootCmd
}
