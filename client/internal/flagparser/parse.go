// Package flagparser provides functions to parse flags from JSON output.
package flagparser

import (
	"fmt"
	"time"

	"github.com/ByteTheCookies/cookieclient/internal/models"
	"github.com/ByteTheCookies/cookieclient/internal/utils"
	json "github.com/bytedance/sonic"
)

// ParseLine parses a JSON line into a Flag struct.
func ParseLine(line string) (models.Flag, error) {
	var out models.ParsedFlagOutput
	if err := json.Unmarshal([]byte(line), &out); err != nil {
		return models.Flag{}, fmt.Errorf("invalid JSON format: %w", err)
	}

	if out.Status == "failed" {
		return models.Flag{}, fmt.Errorf("flag submission failed for team %d on the %s: %s", out.TeamID, utils.MapPortToService(uint16(out.ServicePort)), out.Message)
	}

	flag := models.Flag{
		FlagCode:     out.FlagCode,
		ServiceName:  utils.MapPortToService(uint16(out.ServicePort)),
		ServicePort:  out.ServicePort,
		SubmitTime:   uint64(time.Now().Unix()),
		ResponseTime: 0,
		Status:       "UNSUBMITTED",
		TeamID:       out.TeamID,
	}

	return flag, nil
}
