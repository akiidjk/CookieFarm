package cmd

import (
	"logger"
	"os"
	"path/filepath"

	"client/api"
	"client/config"

	"github.com/spf13/cobra"
)

var (
	port     uint16
	https    bool
	host     string
	username string
	password string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage client configuration",
	Long:  `This command allows you to manage the client configuration, including setting the server host, port, and other parameters.`,
}

var resetConfigCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the client configuration",
	Long:  `This command resets the client configuration to its default state, removing any custom settings that have been applied.`,
	Run:   reset,
}

var editConfigCmd = &cobra.Command{
	Use:   "update",
	Short: "Update client configuration",
	Long:  `This command allows you to edit the client configuration interactively. It opens the configuration file in your default text editor, enabling you to make changes to settings such as server host, port, and other parameters.`, //nolint:revive
	Run:   update,
}

var loginConfigCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the client",
	Long:  `This command allows you to log in to the client, providing your credentials to access protected resources. It will prompt for your username and password, and store the session information securely.`, //nolint:revive
	Run:   login,
}

var logoutConfigCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove client session",
	Long:  `This command removes the current client session, effectively logging you out of the client. It clears any stored session information, ensuring that subsequent requests will require re-authentication.`, //nolint:revive
	Run:   logout,
}

var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the current client configuration",
	Long:  `This command displays the current client configuration settings, including server host, port, username, and other parameters.`,
	Run:   show,
}

func buildConfigCmd() *cobra.Command {
	editConfigCmd.Flags().StringVarP(&host, "host", "H", "localhost", "Server host to connect to")
	editConfigCmd.Flags().Uint16VarP(&port, "port", "p", 8080, "Server port to connect to")
	editConfigCmd.Flags().StringVarP(&username, "username", "u", "cookieguest", "Username for authenticating to the server")
	editConfigCmd.Flags().BoolVarP(&https, "https", "s", false, "Use HTTPS for secure communication with the server")

	loginConfigCmd.Flags().StringVarP(&host, "host", "H", "localhost", "Server host to connect to")
	loginConfigCmd.Flags().Uint16VarP(&port, "port", "p", 8080, "Server port to connect to")
	loginConfigCmd.Flags().StringVarP(&username, "username", "u", "cookieguest", "Username for authenticating to the server")
	loginConfigCmd.Flags().BoolVarP(&https, "https", "s", false, "Use HTTPS for secure communication with the server")
	loginConfigCmd.Flags().StringVarP(&password, "password", "P", "", "Password for authenticating to the server")
	loginConfigCmd.MarkFlagRequired("password")

	configCmd.AddCommand(resetConfigCmd)
	configCmd.AddCommand(editConfigCmd)
	configCmd.AddCommand(loginConfigCmd)
	configCmd.AddCommand(logoutConfigCmd)
	configCmd.AddCommand(showConfigCmd)
	
	return configCmd
}

func reset(cmd *cobra.Command, args []string) {
	cm := config.GetInstance()
	err := cm.Reset()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Reset configuration failed")
		return
	}
	logger.Log.Info().Msg("Configuration file reset successfully.")
}

func update(cmd *cobra.Command, args []string) {
	cm := config.GetInstance()
	if host == "localhost" && port == 8080 && username == "cookieguest" && !https {
		logger.Log.Warn().Msg("All default args detected. Update skipped. For available options, run `ckc config update --help`")
		return
	}

	cm.SetHost(host)
	cm.SetPort(port)
	cm.SetUsername(username)
	cm.SetHTTPS(https)
	
	err := cm.Write()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Update configuration failed")
		return
	}
}

func login(cmd *cobra.Command, args []string) {
	update(cmd, args)
	sessionPath, err := LoginHandler(password)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Login failed")
		return
	}
	logger.Log.Info().Str("path", sessionPath).Msg("Session token stored.")
}

func logout(cmd *cobra.Command, args []string) {
	cm := config.GetInstance()
	_, err := cm.Logout()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Logout failed")
		return
	}
	logger.Log.Info().Msg("Logged out successfully. Session file removed.")
}

func show(cmd *cobra.Command, args []string) {
	cm := config.GetInstance()
	content, err := cm.ShowLocalConfigContent()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Show configuration failed")
		return
	}
	logger.Log.Info().Msg("Current configuration: \n```yaml\n" + content + "```")
}

func LoginHandler(password string) (string, error) {
	cm := config.GetInstance()
	
	err := api.Login(cm.GetUsername(), password)
	if err != nil {
		return "", err
	}
	
	sessionPath := filepath.Join(config.DefaultPath, "session")
	err = os.WriteFile(sessionPath, []byte(cm.GetToken()), 0o644)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error writing session token to file")
		return "", err
	}

	return sessionPath, nil
}
