package protocols

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var HOST string = "http://localhost:3000"
var TEAM_TOKEN string = "4242424242424242"

func Submit(flags []string) (map[string]interface{}, error) {
	// Codifica le flags in JSON
	jsonData, err := json.Marshal(flags)
	if err != nil {
		return nil, fmt.Errorf("errore nel marshalling: %w", err)
	}

	// Crea la request
	req, err := http.NewRequest(http.MethodPut, HOST+"/submit", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("errore nella creazione della richiesta: %w", err)
	}

	// Aggiungi header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Team-Token", TEAM_TOKEN)

	// Esegui la richiesta
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("errore nell'invio: %w", err)
	}
	defer resp.Body.Close()

	// Leggi la risposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("errore nella lettura della risposta: %w", err)
	}

	// Decodifica JSON di risposta
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("errore nel parse della risposta: %w", err)
	}

	return response, nil
}
