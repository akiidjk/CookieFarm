package tui

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ByteTheCookies/cookieclient/cmd"
	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/ByteTheCookies/cookieclient/internal/exploit"
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
		return config.Show()
	case "reset":
		return config.Reset()
	case "logout":
		return config.Logout()
	default:
		return "", fmt.Errorf("unknown config subcommand: %s", subcommand)
	}
}

// ExecuteLogin handles the login command
func (*CommandRunner) ExecuteLogin(password string) (string, error) {
	cmd.Password = password
	pathSession, err := cmd.LoginHandler(password)
	if err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}
	return "Login successfully session saved at " + pathSession, nil
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

	path, err := config.Update(configuration)
	if err != nil {
		return "", fmt.Errorf("failed to update configuration: %w", err)
	}
	return "Configuration updated successfully. File saved at:" + path, nil
}

// Globale per condividere i dati della tabella
var ExploitTableData []ExploitProcess

// ExecuteExploitCommand executes exploit-related commands
func (*CommandRunner) ExecuteExploitCommand(subcommand string) (string, error) {
	switch subcommand {
	case "list":
		output, err := cmd.ListFunc()
		if err != nil {
			return "", err
		}

		// Parse the output and store data for table view
		formattedOutput := formatExploitListOutput(output)

		// Return original output for text view
		return formattedOutput, nil
	default:
		return "", fmt.Errorf("unknown exploit subcommand: %s", subcommand)
	}
}

// formatExploitListOutput formats the exploit list output as a table
// Returns both string output and structured data for table display
func formatExploitListOutput(output string) string {
	if strings.Contains(output, "No running exploits found") {
		return output
	}

	lines := strings.Split(output, "\n")
	var header, footer string
	var exploitLines []string

	// Extract header, footer, and exploit lines
	for _, line := range lines {
		switch {
		case strings.Contains(line, "===== Running Exploits ====="):
			header = line
		case strings.Contains(line, "Total:"):
			footer = line
		case strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "="):
			exploitLines = append(exploitLines, line)
		}
	}

	if len(exploitLines) == 0 {
		return output // Return original if no exploits found
	}

	// Parse exploit entries
	// Preallocate the array based on the number of exploit lines
	exploitData := make([]ExploitProcess, 0, len(exploitLines))
	for _, line := range exploitLines {
		// Extract data using string splitting
		parts := strings.Split(line, "Name: ")
		if len(parts) < 2 {
			continue
		}

		idPart := strings.TrimSpace(parts[0])
		namePidParts := strings.Split(parts[1], "PID: ")
		if len(namePidParts) < 2 {
			continue
		}

		idStr := strings.TrimSpace(strings.Trim(idPart, ". "))
		id, _ := strconv.Atoi(idStr)
		name := strings.TrimSpace(namePidParts[0])
		pidStr := strings.TrimSpace(namePidParts[1])
		pid, _ := strconv.Atoi(pidStr)

		exploitData = append(exploitData, ExploitProcess{
			ID:   id,
			Name: name,
			PID:  pid,
		})
	}

	// Create a formatted simple text table as fallback
	var result strings.Builder
	result.WriteString(header + "\n\n")
	result.WriteString(fmt.Sprintf("  %-5s | %-30s | %-10s\n", "ID", "NAME", "PID"))
	result.WriteString(fmt.Sprintf("  %s-+-%s-+-%s\n", strings.Repeat("-", 5), strings.Repeat("-", 30), strings.Repeat("-", 10)))

	for _, exploitS := range exploitData {
		result.WriteString(fmt.Sprintf("  %-5d | %-30s | %-10d\n", exploitS.ID, exploitS.Name, exploitS.PID))
	}
	result.WriteString("\n" + footer)

	// Store the parsed exploit data for the table component
	ExploitTableData = exploitData

	return result.String()
}

func (r *CommandRunner) ExecuteExploitRun(
	exploitPath, servicePort string,
	tickTime, threadCount string,
) (string, error) {
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

	result, err := exploit.Run(
		exploitPath,
		tickTimeInt,
		threadCountInt,
		servicePortUint16,
	)
	if err != nil {
		return "", err
	}

	r.currentExploitPID = result.PID

	go func() {
		for line := range result.OutputChan {
			select {
			case r.exploitOutputChan <- ExploitOutput{Content: line, PID: result.PID}:
			default:
			}
		}
	}()

	go func() {
		for err := range result.ErrorChan {
			select {
			case r.exploitOutputChan <- ExploitOutput{Error: err, PID: result.PID}:
			default:
			}
		}
	}()

	var initialOutput strings.Builder
	initialOutput.WriteString(fmt.Sprintf("Exploit started with PID: %d\n", result.PID))
	initialOutput.WriteString(fmt.Sprintf("Running with %d threads, tick time %d seconds\n", threadCountInt, tickTimeInt))
	initialOutput.WriteString("Output streaming started. Live updates will appear below...\n")

	return initialOutput.String(), nil
}

// ExecuteExploitCreate handles creating an exploit template
func (*CommandRunner) ExecuteExploitCreate(name string) (string, error) {
	return exploit.Create(name)
}

// ExecuteExploitRemove handles removing an exploit template
func (*CommandRunner) ExecuteExploitRemove(name string) (string, error) {
	return exploit.Remove(name)
}

// ExploitProcess represents a running exploit process
type ExploitProcess struct {
	ID   int
	Name string
	PID  int
}

// GetRunningExploits returns a list of running exploit processes
func (*CommandRunner) GetRunningExploits() ([]ExploitProcess, error) {
	// Load config to get current exploits
	if err := config.LoadLocalConfig(); err != nil {
		return nil, fmt.Errorf("error loading configuration: %w", err)
	}

	var processes []ExploitProcess
	id := 1

	// Check each exploit process
	for _, exploitS := range config.ArgsConfigInstance.Exploits {
		proc, err := os.FindProcess(exploitS.PID)
		if err != nil || proc == nil {
			continue // Skip invalid processes
		}

		// Check if process is running
		err = proc.Signal(syscall.Signal(0))
		if err == nil {
			processes = append(processes, ExploitProcess{
				ID:   id,
				Name: filepath.Base(exploitS.Name),
				PID:  exploitS.PID,
			})
			id++
		}
	}

	return processes, nil
}

// ExecuteExploitStop handles stopping a running exploit
func (*CommandRunner) ExecuteExploitStop(pid string) (string, error) {
	// Validate that pid is a number
	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		return "", fmt.Errorf("invalid process ID: %s", pid)
	}

	// Set the global Pid variable that StopFunc expects
	config.PID = pidInt

	return exploit.Stop(pidInt)
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
