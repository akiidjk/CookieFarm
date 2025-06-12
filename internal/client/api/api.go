// Package api provides functions to interact with the CookieFarm server API.
package api

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ByteTheCookies/CookieFarm/internal/client/config"
	"github.com/ByteTheCookies/CookieFarm/pkg/logger"
	"github.com/ByteTheCookies/CookieFarm/pkg/models"
	json "github.com/bytedance/sonic"
	"github.com/rs/zerolog/log"
)

// GetConfig retrieves the configuration from the CookieFarm server API.
func GetConfig() (models.ConfigShared, error) {
	cm := config.GetConfigManager()
	serverURL := "http://" + cm.GetLocalConfig().Host + ":" + strconv.Itoa(int(cm.GetLocalConfig().Port)) + "/api/v1/config"
	client := &http.Client{}

	_, err := url.Parse(serverURL)
	if err != nil {
		log.Fatal().Msg("Invalid base URL in config")
	}

	req, err := http.NewRequest(http.MethodGet, serverURL, nil)
	if err != nil {
		return models.ConfigShared{}, fmt.Errorf("error creating config request: %w", err)
	}
	req.Header.Set("Cookie", "token="+cm.GetToken())

	resp, err := client.Do(req)
	if err != nil {
		return models.ConfigShared{}, fmt.Errorf("error sending config request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.ConfigShared{}, fmt.Errorf("error reading config response: %w", err)
	}

	var parsedConfig models.ConfigShared
	if err := json.Unmarshal(respBody, &parsedConfig); err != nil {
		return models.ConfigShared{}, fmt.Errorf("error parsing config: %w", err)
	}

	logger.Log.Debug().Msgf("Configuration received correctly")

	return parsedConfig, nil
}

// Login sends a login request to the CookieFarm server API.
func Login(password string) (string, error) {
	cm := config.GetConfigManager()
	serverURL := "http://" + cm.GetLocalConfig().Host + ":" + strconv.Itoa(int(cm.GetLocalConfig().Port)) + "/api/v1/auth/login"

	_, err := url.Parse(serverURL)
	if err != nil {
		log.Fatal().Msg("Invalid base URL in config")
	}

	logger.Log.Debug().Str("url", serverURL).Msg("Login attempt")

	resp, err := http.Post(
		serverURL,
		"application/x-www-form-urlencoded",
		bytes.NewBufferString("username="+cm.GetLocalConfig().Username+"&password="+password),
	)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error sending login request")
		return "", err
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()
	for _, c := range cookies {
		if c.Name == "token" {
			logger.Log.Debug().Str("token", c.Value).Msg("Token found")
			logger.Log.Info().Msg("Login successfully")
			return c.Value, nil
		}
	}

	logger.Log.Warn().Msg("Token not found in Set-Cookie")
	return "", errors.New("token not found in Set-Cookie")
}
