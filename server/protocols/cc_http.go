//go:build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ByteTheCookies/cookieserver/internal/models"
)

func Submit(host string, team_token string, flags []string) ([]models.ResponseProtocol, error) {
	jsonData, err := json.Marshal(flags)
	if err != nil {
		return nil, fmt.Errorf("errore nel marshalling: %w", err)
	}

	url := "http://" + host + "/submit"
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("errore nella creazione della richiesta: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Team-Token", team_token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("errore nell'invio: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("errore nella lettura della risposta: %w", err)
	}

	var response []models.ResponseProtocol
	// logger.Debug("Raw body %s", string(body))
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("errore nel parse della risposta: %w", err)
	}

	return response, nil
}
