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

	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/ByteTheCookies/cookieclient/internal/logger"
	json "github.com/bytedance/sonic"
	"github.com/rs/zerolog/log"
)

// Flag represents a single flag captured during a CTF round.
// It includes metadata about the submission and the service context.
type Flag struct {
	SubmitTime   uint64 `json:"submit_time"`   // UNIX timestamp when the flag was submitted
	ResponseTime uint64 `json:"response_time"` // UNIX timestamp when a response was received
	PortService  uint16 `json:"port_service"`  // Port of the vulnerable service
	TeamID       uint16 `json:"team_id"`       // ID of the team the flag was captured from
	Status       string `json:"status"`        // Status of the submission (e.g., "unsubmitted", "accepted", "denied")
	FlagCode     string `json:"flag_code"`     // Actual flag string
	ServiceName  string `json:"service_name"`  // Human-readable name of the service
	Msg          string `json:"msg"`           // Message from the flag checker service
}

var client = &http.Client{}

// SendFlag sends flags to the CookieFarm server API.
func SendFlag(flags ...Flag) error {
	body, err := json.Marshal(map[string][]Flag{"flags": flags})
	if err != nil {
		logger.Log.Error().Err(err).Msg("error during flags marshalling")
		return err
	}

	ServerURL := "http://" + config.ArgsConfigInstance.Address + ":" + strconv.Itoa(int(config.ArgsConfigInstance.Port)) + "/api/v1/submit-flags"
	req, err := http.NewRequest(http.MethodPost, ServerURL, bytes.NewReader(body))
	if err != nil {
		log.Error().Err(err).Str("url", ServerURL).Msg("error creating request")
		return err
	}

	req.Header.Set("Cookie", "token="+config.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		logger.Log.Error().Err(err).Str("url", ServerURL).Msg("error during flags submission request")
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error reading response")
		return err
	}

	logger.Log.Debug().
		Int("status", resp.StatusCode).
		Msgf("Server response from submit-flags: %s", string(respBody))

	return nil
}

// GetConfig retrieves the configuration from the CookieFarm server API.
func GetConfig() (config.Config, error) {
	ServerURL := "http://" + config.ArgsConfigInstance.Address + ":" + strconv.Itoa(int(config.ArgsConfigInstance.Port)) + "/api/v1/config"
	req, err := http.NewRequest(http.MethodGet, ServerURL, nil)
	if err != nil {
		return config.Config{}, fmt.Errorf("error creating config request: %w", err)
	}
	req.Header.Set("Cookie", "token="+config.Token)

	resp, err := client.Do(req)
	if err != nil {
		return config.Config{}, fmt.Errorf("error sending config request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return config.Config{}, fmt.Errorf("error reading config response: %w", err)
	}

	var parsedConfig config.Config
	if err := json.Unmarshal(respBody, &parsedConfig); err != nil {
		return config.Config{}, fmt.Errorf("error parsing config: %w", err)
	}

	logger.Log.Debug().Msgf("Configuration received correctly")

	return parsedConfig, nil
}

// Login sends a login request to the CookieFarm server API.
func Login(password string) (string, error) {
	ServerURL := "http://" + config.ArgsConfigInstance.Address + ":" + strconv.Itoa(int(config.ArgsConfigInstance.Port)) + "/api/v1/auth/login"

	_, err := url.Parse(ServerURL)
	if err != nil {
		log.Fatal().Msg("Invalid base URL in config")
	}

	logger.Log.Debug().Str("url", ServerURL).Msg("Login attempt")

	resp, err := http.Post(
		ServerURL,
		"application/x-www-form-urlencoded",
		bytes.NewBufferString("username="+config.ArgsConfigInstance.Nickname+"&password="+password),
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
