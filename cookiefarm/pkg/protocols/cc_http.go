//go:build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"models"
	"net/http"
	"strings"

	"protocols"
)

func Submit(url string, teamToken string, flags []string) ([]protocols.ResponseProtocol, error) {
	jsonData, err := json.Marshal(flags)
	if err != nil {
		return nil, fmt.Errorf("error during marshalling: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error during request creation: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Team-Token", teamToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error during request submission: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error during response reading: %w", err)
	}

	var responses []struct {
		Status string `json:"status"` // Status of the response (e.g., "0", "1", "2" ,"3") see enum in pkg/models/models.go
		Flag   string `json:"flag"`   // Flag string received from the flag checker service
		Msg    string `json:"msg"`
	}
	// logger.Debug("Raw body %s", string(body))
	if err := json.Unmarshal(body, &responses); err != nil {
		return nil, fmt.Errorf("error during response parsing: %w", err)
	}

	var responsesParsed []protocols.ResponseProtocol = make([]protocols.ResponseProtocol, len(responses))
	for i := range responses {
		responsesParsed[i].Flag = responses[i].Flag
		responsesParsed[i].Msg = strings.Split(responses[i].Msg, "]")[1]

		switch responses[i].Status {
		case "ACCEPTED":
			responsesParsed[i].Status = models.StatusAccepted
		case "DENIED":
			responsesParsed[i].Status = models.StatusDenied
		case "Error":
			responsesParsed[i].Status = models.StatusError
		default:
			responsesParsed[i].Status = models.StatusNotValid
		}
	}

	return responsesParsed, nil
}
