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
	formattedBody, err := json.Marshal(map[string][]models.Flag{"flags": flags})
	if err != nil {
		fmt.Println("Errore marshalling flags:", err)
		return
	}

	resp, err := http.Post(config.BASE_URL_SERVER+"/submit-flags", "application/json", bytes.NewReader(formattedBody))
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
	resp, err := http.Get(config.BASE_URL_SERVER + "/config")
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
