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
)

func SendFlag(flags ...models.Flag) {
	client := &http.Client{}
	formattedBody, err := json.Marshal(map[string][]models.Flag{"flags": flags})
	if err != nil {
		fmt.Println("Errore marshalling flags:", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, config.Current.ConfigClient.BaseUrlServer+"/api/v1/submit-flags", bytes.NewReader(formattedBody))
	if err != nil {
		fmt.Println("Errore creazione richiesta:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+config.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errore invio flags:", err)
		return
	}
	defer resp.Body.Close()

	bodyContent, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Errore lettura risposta:", err)
		return
	}

	logger.Info("Response %v", string(bodyContent))
}

func GetConfig() models.Config {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, config.Current.ConfigClient.BaseUrlServer+"/api/v1/config", nil)
	if err != nil {
		logger.Fatal("Error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+config.Token)

	resp, err := client.Do(req)
	if err != nil {
		logger.Fatal("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	bodyContent, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Fatal("Error reading response: %v", err)
	}

	var config models.Config
	err = json.Unmarshal(bodyContent, &config)
	if err != nil {
		logger.Fatal("Error parsing config: %v", err)
	}

	logger.Debug("%v", config)

	return config
}

func Login(password string) (string, error) {
	logger.Debug("%s", config.Current.ConfigClient.BaseUrlServer+"/api/v1/auth/login")
	resp, err := http.Post(config.Current.ConfigClient.BaseUrlServer+"/api/v1/auth/login", "application/x-www-form-urlencoded", bytes.NewBufferString(fmt.Sprintf(`password=%s`, password)))
	if err != nil {
		fmt.Println("Errore login:", err)
		return "", err
	}
	defer resp.Body.Close()

	bodyContent, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Errore lettura risposta:", err)
		return "", err
	}

	var result models.TokenResponse
	err = json.Unmarshal(bodyContent, &result)

	logger.Info("Response %v", string(bodyContent))
	return string(result.Token), nil
}
