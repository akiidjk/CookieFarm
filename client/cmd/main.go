package main

import (
	"fmt"
	"os"

	"github.com/ByteTheCookies/cookiefarm-client/internal/api"
	"github.com/ByteTheCookies/cookiefarm-client/internal/config"
	"github.com/ByteTheCookies/cookiefarm-client/internal/executor"
	"github.com/ByteTheCookies/cookiefarm-client/internal/logger"
	"github.com/ByteTheCookies/cookiefarm-client/internal/models"
	"github.com/ByteTheCookies/cookiefarm-client/internal/submitter"
	"github.com/ByteTheCookies/cookiefarm-client/internal/utils"
	"github.com/rs/zerolog"

	"github.com/spf13/pflag"
)

var (
	args    models.Args = models.Args{}
	logPath string
)

func init() {
	args.ExploitName = pflag.StringP("exploit", "e", "", "Name of the exploit file")
	args.Debug = pflag.Bool("debug", false, "Enable debug log level")
	args.Password = pflag.StringP("password", "p", "", "Password for authentication")
	args.BaseURLServer = pflag.StringP("base_url_server", "b", "", "Base URL of the target server")
	args.Detach = pflag.BoolP("detach", "d", false, "Run the exploit in the background")
	args.TickTime = pflag.IntP("tick", "t", 120, "Interval in seconds between run exploits")
}

func setupClient() error {
	pflag.Parse()

	if *args.Detach {
		utils.Detach()
	}

	if *args.Debug {
		logPath = logger.Setup("debug")
	} else {
		logPath = logger.Setup("info")
	}

	err := utils.ValidateArgs(args)
	if err != nil {
		return fmt.Errorf("invalid arguments: %w", err)
	}

	config.Current.ConfigClient.BaseUrlServer = *args.BaseURLServer

	config.Token, err = api.Login(*args.Password)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	config.Current, err = api.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	logger.Log.Info().Msgf("Configurazione corrente: %+v", config.Current)

	if !config.Current.Configured {
		logger.Log.Fatal().Msg("Client not configured. Please run the configurator before using the client")
	}

	logger.Log.Info().Msg("Client initialized successfully")
	return nil
}

func main() {

	if err := setupClient(); err != nil {
		if logger.LogLevel != zerolog.Disabled {
			logger.Log.Fatal().Err(err).Msg("Initialization error")
			logger.Close()
		} else {
			fmt.Println("Errore inizializzazione:", err)
		}
		os.Exit(1)
	}
	defer logger.Close()

	result, err := executor.Start(*args.ExploitName, *args.Password, *args.TickTime, logPath)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to execute exploit")
	}

	go submitter.Start(result.FlagsChan)

	if err := result.Cmd.Wait(); err != nil {
		logger.Log.Error().Err(err).Msg("Exploit process exited with error")
	}
}
