// Package cmd contains commands for the CookieFarm client,
// responsible for initializing configuration, validating input,
// and executing exploits in a loop.
package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/ByteTheCookies/cookieclient/internal/api"
	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/ByteTheCookies/cookieclient/internal/executor"
	"github.com/ByteTheCookies/cookieclient/internal/filesystem"
	"github.com/ByteTheCookies/cookieclient/internal/logger"
	"github.com/ByteTheCookies/cookieclient/internal/submitter"
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
	attackCmd.Flags().StringVarP(&config.ArgsAttackInstance.ExploitPath, "exploit", "e", "", "Path to the exploit file to execute")
	attackCmd.Flags().Uint16VarP(&config.ArgsAttackInstance.ServicePort, "port", "p", 0, "Service Port to attack")
	attackCmd.Flags().BoolVarP(&config.ArgsAttackInstance.Detach, "detach", "d", false, "Run the exploit in the background (detached mode)")
	attackCmd.Flags().IntVarP(&config.ArgsAttackInstance.TickTime, "tick", "t", 120, "Interval in seconds between exploit executions")
	attackCmd.Flags().IntVarP(&config.ArgsAttackInstance.ThreadCount, "thread", "T", 5, "Number of concurrent threads to run the exploit with")
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

	err = config.LoadLocalConfig()

	if config.ArgsAttackInstance.Detach {
		fmt.Println(logger.Blue + "[INFO]" + logger.Reset + " | Detaching from terminal")
		Detach()
	}

	config.ArgsAttackInstance.ExploitPath, err = filesystem.NormalizeNamePathExploit(config.ArgsAttackInstance.ExploitPath)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error normalizing exploit name")
		return err
	}

	if !filesystem.IsPath(config.ArgsAttackInstance.ExploitPath) {
		if err != nil {
			logger.Log.Error().Err(err).Msg("Error expanding path")
			return err
		}
		config.ArgsAttackInstance.ExploitPath = filepath.Join(config.DefaultConfigPath, config.ArgsAttackInstance.ExploitPath)
	}

	logger.Log.Debug().Str("Exploit path", config.ArgsAttackInstance.ExploitPath).Msg("Using default exploit path")

	err = ValidateArgs(config.ArgsAttackInstance)
	if err != nil {
		return fmt.Errorf("invalid arguments: %w", err)
	}

	logger.Log.Debug().
		Int("ThreadCount", config.ArgsAttackInstance.ThreadCount).
		Int("Tick time", config.ArgsAttackInstance.TickTime).
		Str("ExploitPath", config.ArgsAttackInstance.ExploitPath).
		Msg("Arguments validated")

	config.Token, err = config.GetSession()
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

	result, err := executor.Start(config.ArgsAttackInstance.ExploitPath, config.ArgsAttackInstance.TickTime, config.ArgsAttackInstance.ThreadCount, config.ArgsAttackInstance.ServicePort)
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

// Detach detaches the current process from the terminal re executing itself.
func Detach() {
	cmd := exec.Command(os.Args[0], os.Args[1:]...)

	filteredArgs := []string{}
	for _, arg := range os.Args[1:] {
		if arg != "--detach" && arg != "-d" {
			filteredArgs = append(filteredArgs, arg)
		}
	}
	cmd = exec.Command(os.Args[0], filteredArgs...)

	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	err := cmd.Start()
	if err != nil {
		fmt.Println(logger.Red+"[ERROR]"+logger.Reset+"| Error during detach:", err)
		os.Exit(1)
	}

	fmt.Println(logger.Yellow+"[WARN]"+logger.Reset+"| Process detached with PID:", cmd.Process.Pid)
	os.Exit(0)
}

// ValidateArgs validates the arguments passed to the program.
func ValidateArgs(args config.ArgsAttack) error {
	if args.TickTime < 1 {
		return errors.New("tick time must be at least 1")
	}

	exploitPath, err := filepath.Abs(args.ExploitPath)
	if err != nil {
		return fmt.Errorf("error resolving exploit path: %v", err)
	}

	if info, err := os.Stat(exploitPath); err == nil && info.Mode()&0o111 == 0 {
		return errors.New("exploit file is not executable")
	}

	if _, err := os.Stat(exploitPath); os.IsNotExist(err) {
		return errors.New("exploit not found in the exploits directory")
	}

	return nil
}
