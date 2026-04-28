package ckp

import (
	"encoding/binary"
	"path/filepath"
	"server/database"
)

const (
	SIZE_TIMESTAMP = 4
	SIZE_PORT      = 2
	SIZE_TEAMID    = 2
)

// Protocol: 4 bytes timestamp | 2 bytes port | 2 bytes teamID | null-terminated flag code | null-terminated exploit name
func buildPayload(flag database.Flag) []byte {
	exploitName := filepath.Base(flag.ExploitName)
	size := SIZE_TIMESTAMP + SIZE_PORT + SIZE_TEAMID +
		len(flag.FlagCode) + 1 +
		len(exploitName) + 1

	payload := make([]byte, size)

	binary.LittleEndian.PutUint32(payload[0:4], uint32(flag.SubmitTime))
	binary.LittleEndian.PutUint16(payload[4:6], flag.PortService)
	binary.LittleEndian.PutUint16(payload[6:8], uint16(flag.TeamID))

	offset := SIZE_TIMESTAMP + SIZE_PORT + SIZE_TEAMID
	offset += copy(payload[offset:], flag.FlagCode)
	payload[offset] = 0 // null terminator
	offset++
	offset += copy(payload[offset:], exploitName)
	payload[offset] = 0 // null terminator

	return payload
}
