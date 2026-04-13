//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"models"
	"net"
	"protocols"
	"strings"
	"time"
)

func Submit(url string, teamToken string, flags []string) ([]protocols.ResponseProtocol, error) {
	_ = teamToken // plaintext TCP protocol does not use team token header

	address := strings.TrimPrefix(url, "tcp://")
	address = strings.TrimPrefix(address, "tcp:")
	address = strings.TrimPrefix(address, "//")

	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("error connecting to submission server: %w", err)
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(15 * time.Second)); err != nil {
		return nil, fmt.Errorf("error setting connection deadline: %w", err)
	}

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Optional welcome banner terminated by an empty line.
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if line == "\n" {
			break
		}
	}

	for _, flag := range flags {
		if _, err := writer.WriteString(flag + "\n"); err != nil {
			return nil, fmt.Errorf("error writing flag: %w", err)
		}
	}
	if err := writer.Flush(); err != nil {
		return nil, fmt.Errorf("error flushing flags: %w", err)
	}

	responsesParsed := make([]protocols.ResponseProtocol, len(flags))
	flagIndex := make(map[string][]int, len(flags))
	for i, f := range flags {
		flagIndex[f] = append(flagIndex[f], i)
	}

	for range len(flags) {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("error reading response: %w", err)
		}
		line = strings.TrimSuffix(line, "\n")

		matchedFlag := ""
		for f := range flagIndex {
			if rest, ok := strings.CutPrefix(line, f); ok {
				if rest == "" || rest[0] == ' ' || rest[0] == '\t' {
					matchedFlag = f
					break
				}
			}
		}
		if matchedFlag == "" {
			return nil, fmt.Errorf("malformed response line: %q", line)
		}

		indices := flagIndex[matchedFlag]
		if len(indices) == 0 {
			return nil, fmt.Errorf("unexpected duplicate response for flag %q", matchedFlag)
		}
		idx := indices[0]
		flagIndex[matchedFlag] = indices[1:]

		rest := strings.TrimLeft(line[len(matchedFlag):], " \t")
		parts := strings.Fields(rest)
		if len(parts) == 0 {
			return nil, fmt.Errorf("missing status code in response: %q", line)
		}
		code := parts[0]
		msg := ""
		if len(parts) > 1 {
			msg = strings.Join(parts[1:], " ")
		}

		responsesParsed[idx].Flag = matchedFlag
		responsesParsed[idx].Msg = msg

		switch code {
		case "OK":
			responsesParsed[idx].Status = models.StatusAccepted
		case "DUP", "OWN", "OLD":
			responsesParsed[idx].Status = models.StatusDenied
		case "ERR":
			responsesParsed[idx].Status = models.StatusError
		case "INV":
			responsesParsed[idx].Status = models.StatusNotValid
		default:
			responsesParsed[idx].Status = models.StatusError
		}
	}

	return responsesParsed, nil
}
