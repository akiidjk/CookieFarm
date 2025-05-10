// Package flagparser provides functions to parse flags from JSON output.
package flagparser

import (
	"errors"
	"fmt"
	"time"

	"github.com/ByteTheCookies/cookieclient/internal/models"
	"github.com/ByteTheCookies/cookieclient/internal/utils"
	json "github.com/bytedance/sonic"
)

// ParseLine parses a JSON line into a Flag struct.
func ParseLine(line string) (models.Flag, string, error) {
	var out models.ParsedFlagOutput
	if err := json.Unmarshal([]byte(line), &out); err != nil {
		return models.Flag{}, "invalid", fmt.Errorf("invalid JSON format: %w", err)
	}

	switch out.Status {
	case "info":
		return models.Flag{}, out.Status, errors.New(out.Message)
	case "failed":
		return models.Flag{}, out.Status, fmt.Errorf("flag submission failed for team %d on the %s: %s",
			out.TeamID, utils.MapPortToService(out.PortService), out.Message)
	case "error":
		return models.Flag{}, out.Status, fmt.Errorf("flag submission error: %s", out.Message)
	case "fatal":
		return models.Flag{}, out.Status, fmt.Errorf("fatal error in the exploiter: %s", out.Message)
	case "success":
		return models.Flag{
			FlagCode:     out.FlagCode,
			ServiceName:  utils.MapPortToService(out.PortService),
			PortService:  out.PortService,
			SubmitTime:   uint64(time.Now().Unix()),
			ResponseTime: 0,
			Status:       "UNSUBMITTED",
			TeamID:       out.TeamID,
		}, out.Status, nil
	default:
		return models.Flag{}, "unknown", fmt.Errorf("unhandled status: %s", out.Status)
	}
}
