package ckp

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"logger"
	"strconv"

	"server/config"
	"server/database"
)

func readUntilDelimiter(r *bufio.Reader, delim []byte, maxSize int) ([]byte, error) {
	if len(delim) == 0 {
		return nil, errors.New("empty delimiter")
	}

	buf := make([]byte, 0, 1024)

	for {
		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF && len(buf) > 0 {
				return nil, io.ErrUnexpectedEOF
			}
			return nil, err
		}

		buf = append(buf, b)

		if maxSize > 0 && len(buf) > maxSize {
			return nil, errors.New("message too large")
		}

		if len(buf) >= len(delim) && bytes.Equal(buf[len(buf)-len(delim):], delim) {
			return buf[:len(buf)-len(delim)], nil
		}
	}
}

var DelimiterBytes = []byte{0xBB, 'T', 0xCC}

func handler(conn Connection) {
	reader := bufio.NewReader(conn.GetNetConn())

	for {
		data, err := readUntilDelimiter(reader, DelimiterBytes, 1024)
		if err != nil {
			return
		}

		flag, err := parse(data)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Failed to parse CKP flag")
			continue
		}

		logger.Log.Trace().
			Str("flag", flag.FlagCode).
			Int64("team_id", flag.TeamID).
			Msg("Received flag from CKP connection")

		if err := database.GetCollector().AddFlag(flag); err != nil {
			logger.Log.Error().Err(err).Msg("Failed to add CKP flag")
		}
	}
}

func HandlerConfig(conn Connection, cfg []byte) {
	err := write(conn, cfg)
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
	logger.Log.Trace().Bytes("raw_data", data).Msg("Parsing CKP flag data")
	if len(data) < 8 {
		return database.Flag{}, errors.New("invalid length")
	}

	port := binary.LittleEndian.Uint16(data[4:6])
	teamID := binary.LittleEndian.Uint16(data[6:8])
	msg := "Flag found for team: " + strconv.Itoa(int(teamID))

	result := database.Flag{
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

	result.ExploitName, _, err = findString(data[8+idx+1:])
	return result, err
}

func write(conn Connection, data []byte) error {
	writer := bufio.NewWriter(conn.GetNetConn())
	_, err := writer.Write(data)
	if err != nil {
		return err
	}
	return writer.Flush()
}
