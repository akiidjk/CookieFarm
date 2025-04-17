package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

	req, err := http.NewRequest(http.MethodPost, config.BASE_URL_SERVER+"/api/v1/submit-flags", bytes.NewReader(formattedBody))
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
	req, err := http.NewRequest(http.MethodGet, config.BASE_URL_SERVER+"/api/v1/config", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+config.Token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errore download config:", err)
		return models.Config{}
	}
	defer resp.Body.Close()

	bodyContent, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Errore lettura risposta:", err)
		return models.Config{}
	}

	var config models.Config
	err = json.Unmarshal(bodyContent, &config)
	if err != nil {
		fmt.Println("Errore parsing config:", err)
		return models.Config{}
	}

	return config
}

func Login(password string) (string, error) {
	resp, err := http.Post(config.BASE_URL_SERVER+"/api/v1/auth/login", "application/x-www-form-urlencoded", bytes.NewBufferString(fmt.Sprintf(`password=%s`, password)))
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
