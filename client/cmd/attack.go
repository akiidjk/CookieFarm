// Package cmd contains commands for the CookieFarm client,
// responsible for initializing configuration, validating input,
// and executing exploits in a loop.
package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/ByteTheCookies/cookieclient/internal/api"
	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/ByteTheCookies/cookieclient/internal/executor"
	"github.com/ByteTheCookies/cookieclient/internal/logger"
	"github.com/ByteTheCookies/cookieclient/internal/submitter"
	"github.com/ByteTheCookies/cookieclient/internal/utils"
	"github.com/ByteTheCookies/cookieclient/internal/websockets"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var attackCmd = &cobra.Command{
	Use:   "attack",
	Short: "Attack the other team with a exploit",
	Long:  `This command allows you to attack the other team with a exploit. You can specify the exploit path and the server host.`, // Da finire
	Run:   attack,
}

// init initializes all command-line flags and binds them to the args struct.
func init() {
	RootCmd.AddCommand(attackCmd)
	attackCmd.Flags().StringVarP(&config.ArgsAttack.ExploitPath, "exploit", "e", "", "Path to the exploit file to execute")
	attackCmd.Flags().Uint16VarP(&config.ArgsAttack.ServicePort, "port", "p", 0, "Service Port to attack")
	attackCmd.Flags().StringVarP(&config.ServerAddress, "host", "H", "", "Host of the cookieserver")
	attackCmd.Flags().BoolVarP(&config.ArgsAttack.Detach, "detach", "d", false, "Run the exploit in the background (detached mode)")
	attackCmd.Flags().IntVarP(&config.ArgsAttack.TickTime, "tick", "t", 120, "Interval in seconds between exploit executions")
	attackCmd.Flags().IntVarP(&config.ArgsAttack.ThreadCount, "thread", "T", 5, "Number of concurrent threads to run the exploit with")
	attackCmd.MarkFlagRequired("exploit")
	attackCmd.MarkFlagRequired("password")
	attackCmd.MarkFlagRequired("port")
	attackCmd.MarkFlagRequired("host")
}

// SetupClient handles the full initialization process:
// - Parse flags
// - Setup logging
// - Validate arguments
// - Authenticate with the server
// - Sync client configuration
func setupClient() error {
	var err error

	if config.ArgsAttack.Detach {
		fmt.Println(utils.Blue + "[INFO]" + utils.Reset + " | Detaching from terminal")
		utils.Detach()
	}

	config.ArgsAttack.ExploitPath, err = utils.NormalizeNamePathExploit(config.ArgsAttack.ExploitPath)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error normalizing exploit name")
		return err
	}

	if !utils.IsPath(config.ArgsAttack.ExploitPath) {
		defaultPath, err := utils.ExpandTilde(config.DefaultConfigPath)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Error expanding path")
			return err
		}
		config.ArgsAttack.ExploitPath = filepath.Join(defaultPath, config.ArgsAttack.ExploitPath)
	}

	logger.Log.Debug().Str("Exploit path", config.ArgsAttack.ExploitPath).Msg("Using default exploit path")

	err = utils.ValidateArgs(config.ArgsAttack)
	if err != nil {
		return fmt.Errorf("invalid arguments: %w", err)
	}

	logger.Log.Debug().
		Int("ThreadCount", config.ArgsAttack.ThreadCount).
		Int("Tick time", config.ArgsAttack.TickTime).
		Str("ExploitPath", config.ArgsAttack.ExploitPath).
		Str("HostServer", config.ServerAddress).
		Msg("Arguments validated")

	config.Token, err = utils.GetSession()
	if err != nil {
		return fmt.Errorf("failed to get session token: %w", err)
	}

	config.Current, err = api.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	logger.Log.Debug().Msgf("Current configuration: %+v", config.Current)

	if !config.Current.Configured {
		logger.Log.Fatal().Msg("Server not configured. Please run the configurator before using the client.")
	}

	return nil
}

// Main is the main execution flow of the CookieFarm client.
// It handles setup, starts the exploit, and manages the flag submission process.
func attack(cmd *cobra.Command, args []string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logger.Log.Info().Msg("Shutting down gracefully...")
		os.Exit(0)
	}()

	if err := setupClient(); err != nil {
		if logger.LogLevel != zerolog.Disabled {
			logger.Log.Fatal().Err(err).Msg("Initialization error")
			logger.Close()
		} else {
			fmt.Println("Error initializing:", err)
		}
		os.Exit(1)
	}
	defer logger.Close()

	logger.Log.Info().Msg("Client initialized successfully")

	result, err := executor.Start(config.ArgsAttack.ExploitPath, config.ArgsAttack.TickTime, config.ArgsAttack.ThreadCount, config.ArgsAttack.ServicePort)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to execute exploit")
	}
	logger.Log.Info().Msg("Exploit started successfully")

	websockets.OnNewConfig = func() {
		executor.RestartGlobal()
	}

	go submitter.Start(result.FlagsChan)

	if err := result.Cmd.Wait(); err != nil {
		logger.Log.Error().Err(err).Msg("Exploit process exited with error")
	}
}
