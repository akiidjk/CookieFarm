package cmd

import (
	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/ByteTheCookies/cookieclient/internal/logger"
	"github.com/spf13/cobra"
)

// configCmd represents the main config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage client configuration",
	Long:  `This command allows you to manage the client configuration, including setting the server host, port, and other parameters.`,
}

// resetConfigCmd represents the config reset command
var resetConfigCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the client configuration",
	Long:  `This command resets the client configuration to its default state, removing any custom settings that have been applied.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Log.Info().Msg("Resetting client configuration...")
		return nil
	},
}

// editConfigCmd represents the config edit command
var editConfigCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit client configuration",
	Long:  `This command allows you to edit the client configuration interactively. It opens the configuration file in your default text editor, enabling you to make changes to settings such as server host, port, and other parameters.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Log.Info().Msg("Opening configuration file for editing...")
		return nil
	},
}

// loginConfigCmd represents the config login command
var loginConfigCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the client",
	Long:  `This command allows you to log in to the client, providing your credentials to access protected resources. It will prompt for your username and password, and store the session information securely.`,
}

var logoutConfigCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove client session",
	Long:  `This command removes the current client session, effectively logging you out of the client. It clears any stored session information, ensuring that subsequent requests will require re-authentication.`,
}

func init() {
	configCmd.AddCommand(resetConfigCmd)
	configCmd.AddCommand(editConfigCmd)
	configCmd.AddCommand(loginConfigCmd)
	configCmd.AddCommand(logoutConfigCmd)

	loginConfigCmd.Flags().StringVarP(&config.Args.Password, "password", "P", "", "Password for authenticating to the server")
}
