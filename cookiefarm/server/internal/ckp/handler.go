package ckp

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"logger"
	"strconv"

	"server/config"
	"server/database"
)

func handler(conn Connection) {
	data, err := read(conn)
	if err != nil {
		return
	}

	flag, err := parse(data)
	if err != nil {
		return
	}

	database.GetCollector().AddFlag(flag)
}

func HandlerConfig(conn Connection, config []byte) {
	err := write(conn, config)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to write config data to CKP connection")
	}
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
	teamID := binary.LittleEndian.Uint16(data[6:8])
	msg := "Flag found for team: " + strconv.Itoa(int(teamID))

	var result = database.Flag{
		SubmitTime:  uint64(binary.LittleEndian.Uint32(data[0:4])),
		PortService: port,
		TeamID:      int64(teamID),
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

func write(conn Connection, data []byte) error {
	writer := bufio.NewWriter(conn.GetNetConn())
	_, err := writer.Write(data)
	if err != nil {
		return err
	}
	return writer.Flush()
}
