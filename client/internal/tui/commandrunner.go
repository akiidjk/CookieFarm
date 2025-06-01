package tui

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/ByteTheCookies/cookieclient/cmd"
	"github.com/ByteTheCookies/cookieclient/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

// ExploitOutput represents output from a running exploit
type ExploitOutput struct {
	Content string
	Error   error
	PID     int
}

// CommandRunner handles the execution of commands for the TUI
type CommandRunner struct {
	// Channel for streaming exploit output
	exploitOutputChan chan ExploitOutput
	currentExploitPID int
}

// NewCommandRunner creates a new command runner
func NewCommandRunner() *CommandRunner {
	return &CommandRunner{
		exploitOutputChan: make(chan ExploitOutput, 100),
	}
}

// GetExploitOutputCmd returns a tea.Cmd that streams exploit output
func (r *CommandRunner) GetExploitOutputCmd() tea.Cmd {
	return func() tea.Msg {
		select {
		case output, ok := <-r.exploitOutputChan:
			if !ok {
				return nil
			}
			return output
		case <-time.After(100 * time.Millisecond):
			// Periodically refresh to ensure UI updates
			return nil
		}
	}
}

// CloseExploitOutput closes the exploit output channel
func (r *CommandRunner) CloseExploitOutput() {
	if r.exploitOutputChan != nil {
		close(r.exploitOutputChan)
	}
}

// GetCurrentExploitPID returns the PID of the currently running exploit
func (r *CommandRunner) GetCurrentExploitPID() int {
	return r.currentExploitPID
}

// ExecuteCommand executes a generic shell command and returns its output
func (*CommandRunner) ExecuteCommand(command string, args ...string) (string, error) {
	cmdH := exec.Command(command, args...)

	// Create buffers for stdout and stderr
	var stdout, stderr bytes.Buffer
	cmdH.Stdout = &stdout
	cmdH.Stderr = &stderr

	// Execute the command
	err := cmdH.Run()

	// Combine stdout and stderr with appropriate formatting
	var output strings.Builder

	if stdout.Len() > 0 {
		output.WriteString(stdout.String())
	}

	if stderr.Len() > 0 {
		if output.Len() > 0 {
			output.WriteString("\n\n")
		}
		output.WriteString("Error output:\n")
		output.WriteString(stderr.String())
	}

	// If there's an error but we have output, return the output and the error
	if err != nil && output.Len() == 0 {
		return "", err
	}

	return output.String(), err
}

// ExecuteConfigCommand executes configuration-related commands
func (*CommandRunner) ExecuteConfigCommand(subcommand string) (string, error) {
	switch subcommand {
	case "show":
		return cmd.ShowConfigFunc()
	case "reset":
		return cmd.ResetConfigFunc()
	case "logout":
		return cmd.LogoutConfigFunc()
	default:
		return "", fmt.Errorf("unknown config subcommand: %s", subcommand)
	}
}

// ExecuteLogin handles the login command
func (*CommandRunner) ExecuteLogin(password string) (string, error) {
	cmd.Password = password
	output, err := cmd.LoginConfigFunc(password)
	if err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}
	return output, nil
}

// ExecuteConfigUpdate handles the config update command
func (*CommandRunner) ExecuteConfigUpdate(host, port, username string, useHTTPS bool) (string, error) {
	configuration := config.ArgsConfig{
		Address:  host,
		Username: username,
		HTTPS:    useHTTPS,
	}

	if port != "" {
		portNum, err := strconv.Atoi(port)
		if err == nil {
			configuration.Port = uint16(portNum)
		}
	}

	return cmd.UpdateConfigFunc(configuration)
}

// ExecuteExploitCommand executes exploit-related commands
func (*CommandRunner) ExecuteExploitCommand(subcommand string) (string, error) {
	switch subcommand {
	case "list":
		return cmd.ListFunc()
	default:
		return "", fmt.Errorf("unknown exploit subcommand: %s", subcommand)
	}
}

// ExecuteExploitRun gestisce l’esecuzione di un exploit e restituisce
// tutto ciò che RunFuncTui ha scritto su stdout/stderr.
func (r *CommandRunner) ExecuteExploitRun(
	exploitPath, servicePort string,
	detach bool,
	tickTime, threadCount string,
) (string, error) {
	// 1) Parsing degli argomenti string -> tipi corretti
	tickTimeInt := 120  // default
	threadCountInt := 5 // default
	var servicePortUint16 uint16
	var err error

	if tickTime != "" {
		tickTimeInt, err = strconv.Atoi(tickTime)
		if err != nil {
			return "", fmt.Errorf("invalid tick time: %s", tickTime)
		}
	}

	if threadCount != "" {
		threadCountInt, err = strconv.Atoi(threadCount)
		if err != nil {
			return "", fmt.Errorf("invalid thread count: %s", threadCount)
		}
	}

	if servicePort == "" {
		return "", errors.New("service port is required")
	}
	port, err := strconv.Atoi(servicePort)
	if err != nil {
		return "", fmt.Errorf("invalid service port: %s", servicePort)
	}
	servicePortUint16 = uint16(port)

	// 2) Avvia l'exploit con streaming in tempo reale
	result, err := cmd.RunFuncTui(
		exploitPath,
		tickTimeInt,
		threadCountInt,
		servicePortUint16,
		detach,
	)

	if err != nil {
		return "", err
	}

	// Salva il PID dell'exploit in esecuzione
	r.currentExploitPID = result.PID

	// Avvia un goroutine per inviare l'output al canale dell'exploit
	go func() {
		for line := range result.OutputChan {
			select {
			case r.exploitOutputChan <- ExploitOutput{Content: line, PID: result.PID}:
				// Output inviato correttamente
			default:
				// Il canale è pieno, continua senza bloccarsi
			}
		}
	}()

	// Avvia un goroutine per monitorare gli errori
	go func() {
		for err := range result.ErrorChan {
			select {
			case r.exploitOutputChan <- ExploitOutput{Error: err, PID: result.PID}:
				// Errore inviato correttamente
			default:
				// Il canale è pieno, continua senza bloccarsi
			}
		}
	}()

	// Costruisci un output iniziale per la visualizzazione immediata
	var initialOutput strings.Builder
	initialOutput.WriteString(fmt.Sprintf("Exploit started with PID: %d\n", result.PID))
	initialOutput.WriteString(fmt.Sprintf("Running with %d threads, tick time %d seconds\n", threadCountInt, tickTimeInt))
	initialOutput.WriteString("Output streaming started. Live updates will appear below...\n")

	return initialOutput.String(), nil
}

// ExecuteExploitCreate handles creating an exploit template
func (*CommandRunner) ExecuteExploitCreate(name string) (string, error) {
	return cmd.CreateFunc(name)
}

// ExecuteExploitRemove handles removing an exploit template
func (*CommandRunner) ExecuteExploitRemove(name string) (string, error) {
	return cmd.RemoveFunc(name)
}

// ExecuteExploitStop handles stopping a running exploit
func (*CommandRunner) ExecuteExploitStop(pid string) (string, error) {
	// Validate that pid is a number
	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		return "", fmt.Errorf("invalid process ID: %s", pid)
	}

	// Set the global Pid variable that StopFunc expects
	cmd.Pid = pidInt

	return cmd.StopFunc(pidInt)
}

// ExecuteWithTimeout executes a command with a timeout
func (*CommandRunner) ExecuteWithTimeout(timeout time.Duration, command string, args ...string) (string, error) {
	// Create the command
	cmdH := exec.Command(command, args...)

	// Create buffers for stdout and stderr
	var stdout, stderr bytes.Buffer
	cmdH.Stdout = &stdout
	cmdH.Stderr = &stderr

	// Start the command
	if err := cmdH.Start(); err != nil {
		return "", err
	}

	// Create a channel for the command to finish
	done := make(chan error, 1)
	go func() {
		done <- cmdH.Wait()
	}()

	// Wait for either the command to finish or the timeout
	select {
	case <-time.After(timeout):
		if err := cmdH.Process.Kill(); err != nil {
			return "", fmt.Errorf("timeout reached but failed to kill process: %v", err)
		}
		return "", fmt.Errorf("command timed out after %v", timeout)
	case err := <-done:
		// Command completed within the timeout
		var output strings.Builder

		if stdout.Len() > 0 {
			output.WriteString(stdout.String())
		}

		if stderr.Len() > 0 {
			if output.Len() > 0 {
				output.WriteString("\n\n")
			}
			output.WriteString("Error output:\n")
			output.WriteString(stderr.String())
		}

		// If there's an error but we have output, return the output and the error
		if err != nil && output.Len() == 0 {
			return "", err
		}

		return output.String(), err
	}
}

// ExecuteBackgroundCommand starts a command in the background and returns immediately
func (*CommandRunner) ExecuteBackgroundCommand(command string, args ...string) (string, error) {
	cmdH := exec.Command(command, args...)

	// Start the command without waiting for it to complete
	if err := cmdH.Start(); err != nil {
		return "", err
	}

	// Return the PID for tracking
	return fmt.Sprintf("Background process started with PID: %d", cmdH.Process.Pid), nil
}
