package flagparser

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ByteTheCookies/cookiefarm-client/internal/models"
	"github.com/ByteTheCookies/cookiefarm-client/internal/utils"
	"github.com/google/uuid"
)

func ParseLine(line string) (models.Flag, error) {
	var out models.ParsedFlagOutput
	if err := json.Unmarshal([]byte(line), &out); err != nil {
		return models.Flag{}, fmt.Errorf("invalid JSON format: %w", err)
	}

	flag := models.Flag{
		ID:           uuid.New().String(),
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
