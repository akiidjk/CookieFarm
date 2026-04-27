package ckp

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"

	"server/config"
	"server/database"
)

type Data struct {
	submitTime  uint32
	port        uint16
	teamId      uint16
	flag        string
	exploitName string
}

func handler(conn Connection) {
	data, err := read(conn)
	if err != nil {
		return
	}

	flagRaw, err := parse(data)
	if err != nil {
		return
	}

	flag := buildFlag(flagRaw)
	database.GetCollector().AddFlag(flag)
}

func buildFlag(flag Data) database.Flag {
	return database.Flag{
		FlagCode:    flag.flag,
		ServiceName: config.GetInstance().MapPortToService(flag.port),
		PortService: flag.port,
		SubmitTime:  uint64(flag.submitTime),
		TeamID:      flag.teamId,
		ExploitName: flag.exploitName,
		// username
	}
}

func findString(data []byte) (string, int, error) {
	idx := bytes.IndexByte(data, 0)
	if idx == -1 {
		return "", -1, errors.New("unterminted flag")
	}

	return string(data[:idx]), idx, nil
}

func parse(data []byte) (Data, error) {
	if len(data) < 8 {
		return Data{}, errors.New("Invalid length")
	}

	var result Data = Data{
		submitTime: binary.LittleEndian.Uint32(data[0:4]),
		port:       binary.LittleEndian.Uint16(data[4:6]),
		teamId:     binary.LittleEndian.Uint16(data[6:8]),
	}

	var idx int
	var err error

	result.flag, idx, err = findString(data[8:])
	if err != nil {
		return Data{}, err
	}

	result.exploitName, idx, err = findString(data[8+idx:])
	return result, err
}

func read(conn Connection) ([]byte, error) {
	reader := bufio.NewReader(conn.GetNetConn())
	return reader.ReadBytes('\n')
}
