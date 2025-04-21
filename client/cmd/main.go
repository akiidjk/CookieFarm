package main

import (
	"fmt"
	"os"

	"github.com/ByteTheCookies/cookiefarm-client/internal/api"
	"github.com/ByteTheCookies/cookiefarm-client/internal/config"
	"github.com/ByteTheCookies/cookiefarm-client/internal/executor"
	"github.com/ByteTheCookies/cookiefarm-client/internal/logger"
	"github.com/ByteTheCookies/cookiefarm-client/internal/submitter"
	"github.com/ByteTheCookies/cookiefarm-client/internal/utils"
	"github.com/rs/zerolog"

	"github.com/spf13/pflag"
)

var (
	exploitName   = pflag.StringP("exploit", "e", "", "Name of the exploit file")
	debug         = pflag.Bool("debug", false, "Enable debug log level")
	password      = pflag.StringP("password", "p", "", "Password for authentication")
	baseURLServer = pflag.StringP("base_url_server", "b", "", "Base URL of the target server")
	detach        = pflag.BoolP("detach", "d", false, "Run the exploit in the background")
	threadsNumber = pflag.IntP("threads", "t", 1, "Number of threads to use")
	tickTime      = pflag.IntP("tick", "i", 120, "Interval in seconds between run exploits")
)

func setupClient() error {
	pflag.Parse()

	if *detach {
		utils.Detach()
	}

	if *exploitName == "" {
		return fmt.Errorf("missing required --exploit argument")
	}
	if *baseURLServer == "" {
		return fmt.Errorf("missing required --base_url_server argument")
	}
	if *password == "" {
		return fmt.Errorf("missing required --password argument")
	}

	if *debug {
		logger.Setup("debug")
	} else {
		logger.Setup("info")
	}

	config.Current.ConfigClient.BaseUrlServer = *baseURLServer

	var err error
	config.Token, err = api.Login(*password)
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

	result, err := executor.Start(*exploitName, *password, *tickTime, *threadsNumber)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to execute exploit")
	}

	go submitter.Start(result.FlagsChan)

	if err := result.Cmd.Wait(); err != nil {
		logger.Log.Error().Err(err).Msg("Exploit process exited with error")
	}
}
