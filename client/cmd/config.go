// Package cmd contains commands for the CookieFarm client
package cmd

import (
	"os"
	"path/filepath"

	"github.com/ByteTheCookies/cookieclient/internal/api"
	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/ByteTheCookies/cookieclient/internal/logger"
	"github.com/spf13/cobra"
)

var Password string

// ===== CONFIG COMMAND DEFINITIONS =====

// configCmd represents the main config command
// Exported for TUI usage
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage client configuration",
	Long:  `This command allows you to manage the client configuration, including setting the server host, port, and other parameters.`,
}

// resetConfigCmd represents the config reset command
var resetConfigCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the client configuration",
	Long:  `This command resets the client configuration to its default state, removing any custom settings that have been applied.`,
	Run:   reset,
}

// editConfigCmd represents the config update command
var editConfigCmd = &cobra.Command{
	Use:   "update",
	Short: "Update client configuration",
	Long:  `This command allows you to edit the client configuration interactively. It opens the configuration file in your default text editor, enabling you to make changes to settings such as server host, port, and other parameters.`,
	Run:   update,
}

// loginConfigCmd represents the config login command
var loginConfigCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the client",
	Long:  `This command allows you to log in to the client, providing your credentials to access protected resources. It will prompt for your username and password, and store the session information securely.`,
	Run:   login,
}

// logoutConfigCmd represents the config logout command
var logoutConfigCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove client session",
	Long:  `This command removes the current client session, effectively logging you out of the client. It clears any stored session information, ensuring that subsequent requests will require re-authentication.`,
	Run:   logout,
}

// showConfigCmd represents the config show command
var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the current client configuration",
	Long:  `This command displays the current client configuration settings, including server host, port, username, and other parameters.`,
	Run:   show,
}

// ===== CONFIG COMMAND FUNCTIONS =====

// reset resets the configuration to defaults
func reset(cmd *cobra.Command, args []string) {
	_, err := config.ResetLocalConfig()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Reset configuration failed")
		return
	}
	logger.Log.Info().Msg("Configuration file reset successfully.")
}

// update updates the configuration with new values
func update(cmd *cobra.Command, args []string) {
	configPath, err := config.UpdateLocalConfig(config.LocalConfig)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Update configuration failed")
		return
	}
	logger.Log.Info().Str("path", configPath).Msg("Configuration created or updated successfully. ")
}

// login handles user login
func login(cmd *cobra.Command, args []string) {
	sessionPath, err := LoginHandler(Password)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Login failed")
		return
	}
	logger.Log.Info().Str("path", sessionPath).Msg("Session token stored.")
}

// logout handles user logout
func logout(cmd *cobra.Command, args []string) {
	_, err := config.Logout()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Logout failed")
		return
	}
	logger.Log.Info().Msg("Logged out successfully. Session file removed.")
}

func show(cmd *cobra.Command, args []string) {
	content, err := config.ShowLocalConfig()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Show configuration failed")
		return
	}
	logger.Log.Info().Msg("Current configuration: \n```yaml\n" + content + "```")
}

// ===== COMMAND INITIALIZATION =====

func init() {
	ConfigCmd.AddCommand(resetConfigCmd)
	ConfigCmd.AddCommand(editConfigCmd)
	ConfigCmd.AddCommand(loginConfigCmd)
	ConfigCmd.AddCommand(logoutConfigCmd)
	ConfigCmd.AddCommand(showConfigCmd)

	editConfigCmd.Flags().StringVarP(&config.LocalConfig.Address, "host", "H", "localhost", "Server host to connect to")
	editConfigCmd.Flags().Uint16VarP(&config.LocalConfig.Port, "port", "p", 8080, "Server port to connect to")
	editConfigCmd.Flags().StringVarP(&config.LocalConfig.Username, "username", "u", "cookieguest", "Username for authenticating to the server")
	editConfigCmd.Flags().BoolVarP(&config.LocalConfig.HTTPS, "https", "s", false, "Use HTTPS for secure communication with the server")

	loginConfigCmd.Flags().StringVarP(&Password, "password", "P", "", "Password for authenticating to the server")
	loginConfigCmd.MarkFlagRequired("password")
}

// LoginHandler handles user login
func LoginHandler(password string) (string, error) {
	err := config.LoadLocalConfig()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error loading local configuration, try to run: `cookieclient config reset`")
		return "", err
	}

	config.Token, err = api.Login(password)
	if err != nil {
		return "", err
	}

	sessionPath := filepath.Join(config.DefaultConfigPath, "session")
	err = os.WriteFile(sessionPath, []byte(config.Token), 0o644)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error writing session token to file")
		return "", err
	}

	return sessionPath, nil
}
