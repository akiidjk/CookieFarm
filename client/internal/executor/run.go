package executor

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/ByteTheCookies/cookieclient/internal/flagparser"
	"github.com/ByteTheCookies/cookieclient/internal/logger"
	"github.com/ByteTheCookies/cookieclient/internal/models"
	"github.com/ByteTheCookies/cookieclient/internal/utils"
	json "github.com/bytedance/sonic"
)

type ExecutionResult struct {
	Cmd       *exec.Cmd
	FlagsChan chan models.Flag
}

func Start(exploitName, password string, tickTime int, logPath string, threadCount int) (*ExecutionResult, error) {
	exploitPath := filepath.Join(utils.GetExecutableDir(), "..", "exploits", exploitName)

	cmd := exec.Command(
		exploitPath,
		config.Current.ConfigClient.BaseUrlServer,
		password,
		strconv.Itoa(tickTime),
		config.Current.ConfigClient.RegexFlag,
		logPath,
		strconv.Itoa(threadCount),
	)

	logger.Log.Debug().
		Str("path", exploitPath).
		Int("tick", tickTime).
		Str("command", cmd.String()).
		Msg("Starting exploit process")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	flagsChan := make(chan models.Flag, 100)

	go readStdout(stdout, flagsChan)
	go readStderr(stderr)

	return &ExecutionResult{
		Cmd:       cmd,
		FlagsChan: flagsChan,
	}, nil
}

func readStdout(stdout io.Reader, flagsChan chan<- models.Flag) {
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		flag, err := flagparser.ParseLine(line)
		if err != nil {
			logger.Log.Warn().
				Err(err).
				Str("raw", line).
				Msg("Failed to parse JSON from stdout")
			continue
		}
		flagsChan <- flag

		if jsonFlag, err := json.Marshal(flag); err == nil {
			logger.Log.Debug().
				Str("flag", string(jsonFlag)).
				Msg("Parsed and queued flag")
		}
	}
	if err := scanner.Err(); err != nil {
		logger.Log.Error().Err(err).Msg("Errore lettura stdout scanner")
	}
}

func readStderr(stderr io.Reader) {
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		logger.Log.Warn().Str("stderr", scanner.Text()).Msg("Exploit stderr")
	}
	if err := scanner.Err(); err != nil {
		logger.Log.Error().Err(err).Msg("Errore lettura stderr scanner")
	}
}
