package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ByteTheCookies/cookiefarm-client/internal/config"
	"github.com/ByteTheCookies/cookiefarm-client/internal/logger"
	"github.com/ByteTheCookies/cookiefarm-client/internal/models"
	"github.com/ByteTheCookies/cookiefarm-client/internal/utils"
	"github.com/rs/zerolog/log"
)

var (
	client = &http.Client{}
)

func SendFlag(flags ...models.Flag) error {
	err := utils.CheckUrl()
	if err != nil {
		return err
	}
	body, err := json.Marshal(map[string][]models.Flag{"flags": flags})
	if err != nil {
		logger.Log.Error().Err(err).Msg("Errore durante il marshalling delle flags")
		return err
	}

	url := config.Current.ConfigClient.BaseUrlServer + "/api/v1/submit-flags"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("Errore creazione richiesta")
		return err
	}

	req.Header.Set("Cookie", "token="+config.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		logger.Log.Error().Err(err).Str("url", url).Msg("Errore durante la richiesta di invio flags")
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Errore lettura risposta server")
		return err
	}

	logger.Log.Info().
		Int("status", resp.StatusCode).
		Msgf("Risposta server invio flags: %s", string(respBody))
	return nil
}

func GetConfig() (models.Config, error) {
	err := utils.CheckUrl()
	if err != nil {
		return models.Config{}, err
	}
	url := config.Current.ConfigClient.BaseUrlServer + "/api/v1/config"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return models.Config{}, fmt.Errorf("errore creazione richiesta config: %w", err)
	}
	req.Header.Set("Cookie", "token="+config.Token)

	resp, err := client.Do(req)
	if err != nil {
		return models.Config{}, fmt.Errorf("errore invio richiesta config: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Config{}, fmt.Errorf("errore lettura risposta config: %w", err)
	}

	var parsedConfig models.Config
	if err := json.Unmarshal(respBody, &parsedConfig); err != nil {
		return models.Config{}, fmt.Errorf("errore parsing config: %w", err)
	}

	if jsonOut, err := json.Marshal(parsedConfig); err == nil {
		logger.Log.Debug().Msgf("Configurazione ricevuta dal server: %s", string(jsonOut))
	}

	return parsedConfig, nil
}

func Login(password string) (string, error) {
	err := utils.CheckUrl()
	if err != nil {
		return "", err
	}
	url := config.Current.ConfigClient.BaseUrlServer + "/api/v1/auth/login"
	logger.Log.Debug().Str("url", url).Msg("Tentativo di login")

	resp, err := http.Post(
		url,
		"application/x-www-form-urlencoded",
		bytes.NewBufferString(fmt.Sprintf(`password=%s`, password)),
	)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Errore invio login")
		return "", err
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()
	for _, c := range cookies {
		if c.Name == "token" {
			logger.Log.Info().Str("token", c.Value).Msg("Login effettuato con successo via cookie")
			return c.Value, nil
		}
	}

	logger.Log.Warn().Msg("Token non trovato nei cookie")
	return "", fmt.Errorf("token non trovato nel Set-Cookie")
}
