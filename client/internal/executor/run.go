// Package executor provides functions to execute exploits and manage their output.
package executor

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"sync"

	"github.com/ByteTheCookies/cookieclient/internal/api"
	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/ByteTheCookies/cookieclient/internal/flagparser"
	"github.com/ByteTheCookies/cookieclient/internal/logger"
	"github.com/ByteTheCookies/cookieclient/internal/submitter"
)

type ExecutionResult struct {
	Cmd        *exec.Cmd
	FlagsChan  chan api.Flag
	stopReader chan struct{}
	done       chan struct{}
}

var (
	GlobalResult *ExecutionResult
	MutexResult  sync.Mutex
)

// Start starts the exploit_manager and listens for flags in stdout.
func Start(exploitPath string, tickTime int, threadCount int, port uint16) (*ExecutionResult, error) {
	cmd := exec.Command(
		exploitPath,
		config.ServerAddress,
		strconv.Itoa(tickTime),
		strconv.Itoa(threadCount),
		strconv.Itoa(int(port)),
		config.MapPortToService(port),
	)

	logger.Log.Debug().
		Str("full path exploit", exploitPath).
		Int("tick time", tickTime).
		Str("command executed", cmd.String())

	logger.Log.Info().Msg("Starting exploiting process with exploit_manager...")

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

	flagsChan := make(chan api.Flag, 500)
	stopReader := make(chan struct{})
	done := make(chan struct{})

	go readStdout(stdout, flagsChan, stopReader, done)
	go readStderr(stderr, stopReader, done)

	return &ExecutionResult{
		Cmd:        cmd,
		FlagsChan:  flagsChan,
		stopReader: stopReader,
		done:       done,
	}, nil
}

func RestartGlobal() {
	MutexResult.Lock()
	defer MutexResult.Unlock()

	if GlobalResult != nil {
		logger.Log.Info().Msg("Stopping existing exploit...")
		_ = GlobalResult.Shutdown()
	}

	logger.Log.Info().Msg("Starting new exploit process...")

	result, err := Start(
		config.ArgsAttackInstance.ExploitPath,
		config.ArgsAttackInstance.TickTime,
		config.ArgsAttackInstance.ThreadCount,
		config.ArgsAttackInstance.ServicePort,
	)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to start new exploit")
	}

	GlobalResult = result

	go func() {
		if err := result.Cmd.Wait(); err != nil {
			logger.Log.Error().Err(err).Msg("Exploit process exited with error")
		}
	}()

	go func() {
		err := submitter.Start(result.FlagsChan)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Submitter exited with error")
		}
	}()
}

// Shutdown end	the exploit process.
func (e *ExecutionResult) Shutdown() error {
	logger.Log.Info().Msg("Shutting down exploit...")

	if err := e.Cmd.Process.Kill(); err != nil {
		return fmt.Errorf("failed to kill exploit process: %w", err)
	}

	close(e.stopReader)
	<-e.done

	close(e.FlagsChan)

	logger.Log.Info().Msg("Exploit shutdown completed.")
	return nil
}

// LogParsedLineError logs an error based on the status and line.
func logParsedLineError(err error, status, line string) {
	switch status {
	case "fatal":
		logger.Log.Fatal().Err(err).Msg("Fatal error")
	case "info":
		logger.Log.Info().Err(err).Msg("Info")
	case "failed":
		logger.Log.Warn().Err(err).Msg("Failed to run exploit")
	default:
		logger.Log.Debug().Err(err).Msg("Parsing warning")
	}
	logger.Log.Debug().Str("raw", line).Msg("Raw line with error")
}

// Read the stdout and parse JSON lines into Flag structs.
func readStdout(stdout io.Reader, flagsChan chan<- api.Flag, stop <-chan struct{}, done chan<- struct{}) {
	defer func() { done <- struct{}{} }()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		select {
		case <-stop:
			logger.Log.Info().Msg("readStdout received shutdown signal")
			return
		default:
			line := scanner.Text()
			flag, status, err := flagparser.ParseLine(line)
			if err != nil {
				logParsedLineError(err, status, line)
				continue
			}
			flagsChan <- flag
			logger.Log.Info().
				Str("flag", flag.FlagCode).
				Int("team", int(flag.TeamID)).
				Str("service", flag.ServiceName).
				Uint16("port", flag.PortService).
				Msg("Parsed and queued flag")
		}
	}
	if err := scanner.Err(); err != nil {
		logger.Log.Error().Err(err).Msg("Error reading stdout scanner")
	}
}

// Read the stderr and log any errors.
func readStderr(stderr io.Reader, stop <-chan struct{}, done chan<- struct{}) {
	defer func() { done <- struct{}{} }()

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		select {
		case <-stop:
			logger.Log.Info().Msg("readStderr received shutdown signal")
			return
		default:
			logger.Log.Warn().Str("stderr", scanner.Text()).Msg("Exploit stderr")
		}
	}
	if err := scanner.Err(); err != nil {
		logger.Log.Error().Err(err).Msg("Error reading stderr scanner")
	}
}
