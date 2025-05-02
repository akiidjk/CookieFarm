// Package api provides functions to interact with the CookieFarm server API.
package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/ByteTheCookies/cookieclient/internal/logger"
	"github.com/ByteTheCookies/cookieclient/internal/models"
	json "github.com/bytedance/sonic"
	"github.com/rs/zerolog/log"
)

var (
	client = &http.Client{}
)

// SendFlag sends flags to the CookieFarm server API.
func SendFlag(flags ...models.Flag) error {
	body, err := json.Marshal(map[string][]models.Flag{"flags": flags})
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error during flags marshalling")
		return err
	}

	url := *config.BaseURLServer + "/api/v1/submit-flags"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("Error creating request")
		return err
	}

	req.Header.Set("Cookie", "token="+config.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		logger.Log.Error().Err(err).Str("url", url).Msg("Error during flags submission request")
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error reading response")
		return err
	}

	logger.Log.Debug().
		Int("status", resp.StatusCode).
		Msgf("Server response from submit-flags: %s", string(respBody))

	return nil
}

// GetConfig retrieves the configuration from the CookieFarm server API.
func GetConfig() (models.Config, error) {
	url := *config.BaseURLServer + "/api/v1/config"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return models.Config{}, fmt.Errorf("Error creating config request: %w", err)
	}
	req.Header.Set("Cookie", "token="+config.Token)

	resp, err := client.Do(req)
	if err != nil {
		return models.Config{}, fmt.Errorf("Error sending config request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Config{}, fmt.Errorf("Error reading config response: %w", err)
	}

	var parsedConfig models.Config
	if err := json.Unmarshal(respBody, &parsedConfig); err != nil {
		return models.Config{}, fmt.Errorf("Error parsing config: %w", err)
	}

	logger.Log.Debug().Msgf("Configuration received correctly")

	return parsedConfig, nil
}

// Login sends a login request to the CookieFarm server API.
func Login(password string) (string, error) {
	url := *config.BaseURLServer + "/api/v1/auth/login"

	logger.Log.Debug().Str("url", url).Msg("Tentativo di login")

	resp, err := http.Post(
		url,
		"application/x-www-form-urlencoded",
		bytes.NewBufferString(fmt.Sprintf(`password=%s`, password)),
	)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error sending login request")
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

	logger.Log.Warn().Msg("Token non trovato nei cookie")
	return "", fmt.Errorf("token non trovato nel Set-Cookie")
}
