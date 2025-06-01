// Package cmd contains commands for the CookieFarm client
package cmd

import (
	"os"
	"path/filepath"

	"github.com/ByteTheCookies/cookieclient/internal/api"
	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/ByteTheCookies/cookieclient/internal/logger"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var Password string

// ===== CONFIG COMMAND DEFINITIONS =====

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
	Run:   resetConfigFunc,
}

// editConfigCmd represents the config update command
var editConfigCmd = &cobra.Command{
	Use:   "update",
	Short: "Update client configuration",
	Long:  `This command allows you to edit the client configuration interactively. It opens the configuration file in your default text editor, enabling you to make changes to settings such as server host, port, and other parameters.`,
	Run:   updateConfigFunc,
}

// loginConfigCmd represents the config login command
var loginConfigCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the client",
	Long:  `This command allows you to log in to the client, providing your credentials to access protected resources. It will prompt for your username and password, and store the session information securely.`,
	Run:   loginConfigFunc,
}

// logoutConfigCmd represents the config logout command
var logoutConfigCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove client session",
	Long:  `This command removes the current client session, effectively logging you out of the client. It clears any stored session information, ensuring that subsequent requests will require re-authentication.`,
	Run:   logoutConfigFunc,
}

// showConfigCmd represents the config show command
var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the current client configuration",
	Long:  `This command displays the current client configuration settings, including server host, port, username, and other parameters.`,
	Run:   showConfigFunc,
}

// ===== CONFIG COMMAND FUNCTIONS =====

// resetConfigFunc resets the configuration to defaults
func resetConfigFunc(cmd *cobra.Command, args []string) {
	var err error
	err = os.MkdirAll(config.DefaultConfigPath, 0o755)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error creating config directory")
	}

	configPath := filepath.Join(config.DefaultConfigPath, "config.yml")

	file, err := os.Create(configPath)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error opening configuration file")
		return
	}
	defer file.Close()

	err = yaml.Unmarshal(config.ConfigTemplate, &config.ArgsConfigInstance)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error unmarshalling default configuration")
		return
	}

	err = yaml.NewEncoder(file).Encode(config.ArgsConfigInstance)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error encoding configuration to YAML")
		return
	}
	logger.Log.Info().Msg("Configuration file reset successfully.")
}

// updateConfigFunc updates the configuration with new values
func updateConfigFunc(cmd *cobra.Command, args []string) {
	var err error
	err = os.MkdirAll(config.DefaultConfigPath, 0o755)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error creating config directory")
	}

	configPath := filepath.Join(config.DefaultConfigPath, "config.yml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		logger.Log.Warn().Msg("Configuration file does not exist, creating a new one with default settings")
		os.WriteFile(configPath, config.ConfigTemplate, 0o644)
	} else if err != nil {
		logger.Log.Error().Err(err).Msg("Error checking configuration file")
		return
	}

	file, err := os.Create(configPath)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error creating or opening configuration file")
		return
	}
	defer file.Close()

	err = yaml.NewEncoder(file).Encode(config.ArgsConfigInstance)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error encoding configuration to YAML")
		return
	}

	logger.Log.Info().Str("path", configPath).Msg("Configuration created or updated successfully. ")
}

// loginConfigFunc handles user login
func loginConfigFunc(cmd *cobra.Command, args []string) {
	err := config.LoadLocalConfig()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error loading local configuration, try to run: `cookieclient config reset`")
		return
	}

	config.Token, err = api.Login(Password)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Login failed")
		return
	}

	sessionPath := filepath.Join(config.DefaultConfigPath, "session")
	err = os.WriteFile(sessionPath, []byte(config.Token), 0o644)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error writing session token to file")
		return
	}
	logger.Log.Info().Str("path", sessionPath).Msg("Session token stored.")
}

// logoutConfigFunc handles user logout
func logoutConfigFunc(cmd *cobra.Command, args []string) {
	sessionPath := filepath.Join(config.DefaultConfigPath, "session")
	err := os.Remove(sessionPath)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error removing session file")
		return
	}
	logger.Log.Info().Msg("Logged out successfully. Session file removed.")
}

// showConfigFunc displays the current configuration
func showConfigFunc(cmd *cobra.Command, args []string) {
	configPath := filepath.Join(config.DefaultConfigPath, "config.yml")

	content, err := os.ReadFile(configPath)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error reading configuration file")
		return
	}

	logger.Log.Info().Msg("Current configuration: \n```yaml\n" + string(content) + "```")
}

// ===== COMMAND INITIALIZATION =====

func init() {
	// Add subcommands to config command
	configCmd.AddCommand(resetConfigCmd)
	configCmd.AddCommand(editConfigCmd)
	configCmd.AddCommand(loginConfigCmd)
	configCmd.AddCommand(logoutConfigCmd)
	configCmd.AddCommand(showConfigCmd)

	// Setup flags for editConfigCmd
	editConfigCmd.Flags().StringVarP(&config.ArgsConfigInstance.Address, "host", "H", "localhost", "Server host to connect to")
	editConfigCmd.Flags().Uint16VarP(&config.ArgsConfigInstance.Port, "port", "p", 8080, "Server port to connect to")
	editConfigCmd.Flags().StringVarP(&config.ArgsConfigInstance.Nickname, "username", "u", "cookieguest", "Username for authenticating to the server")
	editConfigCmd.Flags().BoolVarP(&config.ArgsConfigInstance.HTTPS, "https", "s", false, "Use HTTPS for secure communication with the server")

	// Setup flags for loginConfigCmd
	loginConfigCmd.Flags().StringVarP(&Password, "password", "P", "", "Password for authenticating to the server")
	loginConfigCmd.MarkFlagRequired("password")
}
