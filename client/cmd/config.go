package cmd

import (
	"os"
	"path/filepath"

	"github.com/ByteTheCookies/cookieclient/internal/api"
	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/ByteTheCookies/cookieclient/internal/logger"
	"github.com/ByteTheCookies/cookieclient/internal/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var Password string

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

func resetConfigFunc(cmd *cobra.Command, args []string) {
	var err error
	expandendPath, err := utils.ExpandTilde(config.DefaultConfigPath)
	err = os.MkdirAll(expandendPath, 0o755)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error creating config directory")
	}

	configPath := filepath.Join(expandendPath, "config.yml")

	file, err := os.Create(configPath)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error opening configuration file")
		return
	}
	defer file.Close()

	err = yaml.Unmarshal([]byte(config.ConfigTemplate), &config.ArgsConfig)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error unmarshalling default configuration")
		return
	}

	err = yaml.NewEncoder(file).Encode(config.ArgsConfig)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error encoding configuration to YAML")
		return
	}
	logger.Log.Info().Msg("Configuration file reset successfully.")
}

// editConfigCmd represents the config edit command
var editConfigCmd = &cobra.Command{
	Use:   "update",
	Short: "Update client configuration",
	Long:  `This command allows you to edit the client configuration interactively. It opens the configuration file in your default text editor, enabling you to make changes to settings such as server host, port, and other parameters.`,
	Run:   updateConfigFunc,
}

func updateConfigFunc(cmd *cobra.Command, args []string) {
	var err error
	expandendPath, err := utils.ExpandTilde(config.DefaultConfigPath)
	err = os.MkdirAll(expandendPath, 0o755)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error creating config directory")
	}

	configPath := filepath.Join(expandendPath, "config.yml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		logger.Log.Warn().Msg("Configuration file does not exist, creating a new one with default settings")
		os.WriteFile(configPath, []byte(config.ConfigTemplate), 0o644)
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

	err = yaml.NewEncoder(file).Encode(config.ArgsConfig)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error encoding configuration to YAML")
		return
	}

	logger.Log.Info().Str("path", configPath).Msg("Configuration created or updated successfully. ")
}

// loginConfigCmd represents the config login command
var loginConfigCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the client",
	Long:  `This command allows you to log in to the client, providing your credentials to access protected resources. It will prompt for your username and password, and store the session information securely.`,
	Run:   loginConfigFunc,
}

func loginConfigFunc(cmd *cobra.Command, args []string) {
	err := utils.LoadLocalConfig()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error loading local configuration")
		return
	}

	config.Token, err = api.Login(Password)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Login failed")
		return
	}

	expandendPath, err := utils.ExpandTilde(config.DefaultConfigPath)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error expanding path for config directory")
		return
	}
	sessionPath := filepath.Join(expandendPath, "session")
	err = os.WriteFile(sessionPath, []byte(config.Token), 0o644)
	logger.Log.Info().Str("path", sessionPath).Msg("Session token stored.")
}

var logoutConfigCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove client session",
	Long:  `This command removes the current client session, effectively logging you out of the client. It clears any stored session information, ensuring that subsequent requests will require re-authentication.`,
	Run:   logoutConfigFunc,
}

func logoutConfigFunc(cmd *cobra.Command, args []string) {
	expandendPath, err := utils.ExpandTilde(config.DefaultConfigPath)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error expanding path for config directory")
		return
	}
	sessionPath := filepath.Join(expandendPath, "session")
	err = os.Remove(sessionPath)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error removing session file")
		return
	}
	logger.Log.Info().Msg("Logged out successfully. Session file removed.")
}

func init() {
	configCmd.AddCommand(resetConfigCmd)
	configCmd.AddCommand(editConfigCmd)
	configCmd.AddCommand(loginConfigCmd)
	configCmd.AddCommand(logoutConfigCmd)

	editConfigCmd.Flags().StringVarP(&config.ArgsConfig.Address, "host", "H", "", "Server host to connect to")
	editConfigCmd.Flags().Uint16VarP(&config.ArgsConfig.Port, "port", "p", 0, "Server port to connect to")
	editConfigCmd.Flags().StringVarP(&config.ArgsConfig.Nickname, "username", "u", "", "Username for authenticating to the server")
	editConfigCmd.Flags().BoolVarP(&config.ArgsConfig.Https, "https", "s", false, "Use HTTPS for secure communication with the server")
	editConfigCmd.MarkFlagRequired("host")
	editConfigCmd.MarkFlagRequired("port")
	editConfigCmd.MarkFlagRequired("username")

	loginConfigCmd.Flags().StringVarP(&Password, "password", "P", "", "Password for authenticating to the server")
	loginConfigCmd.MarkFlagRequired("password")
}
