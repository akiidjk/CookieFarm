package ckp

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"logger"

	"server/config"
	"server/database"
)

func handler(conn Connection) {
	logger.Log.Debug().Msg("New connection received for CKP handler")
	data, err := read(conn)
	if err != nil {
		return
	}

	logger.Log.Debug().Bytes("raw data", data).Msg("Raw data received from CKP connection")

	flag, err := parse(data)
	if err != nil {
		return
	}

	logger.Log.Debug().Str("flag code", flag.FlagCode).
		Int64("team id", flag.TeamID).
		Uint16("port service", flag.PortService).
		Str("service name", flag.ServiceName).
		Str("exploit name", flag.ExploitName).
		Msg("Parsed flag from CKP connection")

	database.GetCollector().AddFlag(flag)
}

func findString(data []byte) (string, int, error) {
	idx := bytes.IndexByte(data, 0)
	if idx == -1 {
		return "", -1, errors.New("unterminted flag")
	}

	return string(data[:idx]), idx, nil
}

func parse(data []byte) (database.Flag, error) {
	if len(data) < 8 {
		return database.Flag{}, errors.New("Invalid length")
	}

	port := binary.LittleEndian.Uint16(data[4:6])
	msg := "Flag found for team: " + string(data[6:8])

	var result database.Flag = database.Flag{
		SubmitTime:  uint64(binary.LittleEndian.Uint32(data[0:4])),
		PortService: port,
		TeamID:      int64(binary.LittleEndian.Uint16(data[6:8])),
		ServiceName: config.GetInstance().MapPortToService(port),
		Msg:         msg,
	}

	var idx int
	var err error

	result.FlagCode, idx, err = findString(data[8:])
	if err != nil {
		return database.Flag{}, err
	}

	result.ExploitName, idx, err = findString(data[8+idx:])
	return result, err
}

func read(conn Connection) ([]byte, error) {
	reader := bufio.NewReader(conn.GetNetConn())
	return reader.ReadBytes('\n')
}
