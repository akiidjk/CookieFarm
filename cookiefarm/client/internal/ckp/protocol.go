package ckp

import (
	"encoding/binary"
	"encoding/json"
	"path/filepath"
	"server/database"
	"sharedconfig"

	"client/config"
)

var OnNewConfig func()

const (
	SizeTimestamp = 4
	SizePort      = 2
	SizeTeamID    = 2
)

// Protocol: 4 bytes timestamp | 2 bytes port | 2 bytes teamID | null-terminated flag code | null-terminated exploit name
func buildPayload(flag database.Flag) []byte {
	exploitName := filepath.Base(flag.ExploitName)
	size := SizeTimestamp + SizePort + SizeTeamID +
		len(flag.FlagCode) + 1 +
		len(exploitName) + 1 + 1

	payload := make([]byte, size)

	binary.LittleEndian.PutUint32(payload[0:4], uint32(flag.SubmitTime))
	binary.LittleEndian.PutUint16(payload[4:6], flag.PortService)
	binary.LittleEndian.PutUint16(payload[6:8], uint16(flag.TeamID))

	offset := SizeTimestamp + SizePort + SizeTeamID
	offset += copy(payload[offset:], flag.FlagCode)
	payload[offset] = 0 // null terminator
	offset++
	offset += copy(payload[offset:], exploitName)
	payload[offset] = 0 // null terminator

	payload[offset+1] = '\n'

	return payload
}

func handleConfig(payload []byte) error {
	var configReceived sharedconfig.Shared

	if err := json.Unmarshal(payload, &configReceived); err != nil {
		return err
	}

	cm := config.GetInstance()
	cm.Get().Shared.Set(configReceived)

	if OnNewConfig != nil {
		go OnNewConfig()
	}

	return nil
}
