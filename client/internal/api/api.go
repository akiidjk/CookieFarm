package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ByteTheCookies/cookiefarm-client/internal/logger"
	"github.com/ByteTheCookies/cookiefarm-client/internal/models"
)

const HOST = "http://localhost:8080"

func SendFlag(flags ...models.Flag) {
	formattedBody, err := json.Marshal(map[string][]models.Flag{"flags": flags})
	if err != nil {
		fmt.Println("Errore marshalling flags:", err)
		return
	}

	resp, err := http.Post(HOST+"/submit-flags", "application/json", bytes.NewReader(formattedBody))
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
