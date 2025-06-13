package core

import (
	"context"
	"os"

	"github.com/ByteTheCookies/CookieFarm/internal/server/config"
	"github.com/ByteTheCookies/CookieFarm/pkg/logger"
	"gopkg.in/yaml.v3"
)

// LoadConfig loads the configuration from the given path.
func LoadConfig(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logger.Log.Error().Err(err).Msg("Configuration file does not exist")
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to read configuration file")
		return err
	}

	err = yaml.Unmarshal(data, &config.SharedConfig)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to parse configuration file")
		return err
	}

	if !config.SharedConfig.Configured {
		config.SharedConfig.Configured = true
	}

	ctx, cancel := context.WithCancel(context.Background())
	if shutdownCancel != nil {
		shutdownCancel()
	}
	shutdownCancel = cancel

	go StartFlagProcessingLoop(ctx)

	return nil
}
