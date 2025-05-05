package protocols

import (
	"fmt"
	"path"
	"plugin"

	"github.com/ByteTheCookies/cookieserver/internal/logger"
	"github.com/ByteTheCookies/cookieserver/internal/models"
	"github.com/ByteTheCookies/cookieserver/internal/utils"
)

type SubmitFunc = func(string, string, []string) ([]models.ResponseProtocol, error)

func LoadProtocol(protocolName string) (SubmitFunc, error) {
	pluginPath := path.Join(utils.GetExecutableDir(), "..", "protocols", protocolName+".so")

	logger.Log.Debug().Str("plugin", pluginPath).Msg("Loading protocol plugin")

	plug, err := plugin.Open(pluginPath)
	if err != nil {
		logger.Log.Error().Err(err).Str("plugin", pluginPath).Msg("Failed to open plugin")
		return nil, fmt.Errorf("failed to open plugin: %w", err)
	}

	submitSymbol, err := plug.Lookup("Submit")
	if err != nil {
		logger.Log.Error().Err(err).Str("plugin", pluginPath).Msg("Submit symbol not found")
		return nil, fmt.Errorf("failed to lookup 'Submit': %w", err)
	}

	submitFunc, ok := submitSymbol.(SubmitFunc)
	if !ok {
		logger.Log.Error().Str("plugin", pluginPath).Msg("Invalid Submit function signature")
		return nil, fmt.Errorf("plugin 'Submit' has invalid signature")
	}

	logger.Log.Info().Str("protocol", protocolName).Msg("Protocol loaded successfully")

	return submitFunc, nil
}
