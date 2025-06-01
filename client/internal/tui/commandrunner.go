package tui

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/ByteTheCookies/cookieclient/cmd"
	"github.com/ByteTheCookies/cookieclient/internal/config"
)

// CommandRunner handles the execution of commands for the TUI
type CommandRunner struct {
	// Additional fields could be added here for configuration
}

// NewCommandRunner creates a new command runner
func NewCommandRunner() *CommandRunner {
	return &CommandRunner{}
}

var defaultFlags = []string{"--no-tui", "--no-banner"}

// ExecuteCommand executes a generic shell command and returns its output
func (r *CommandRunner) ExecuteCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	// Create buffers for stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute the command
	err := cmd.Run()

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
func (r *CommandRunner) ExecuteConfigCommand(subcommand string) (string, error) {
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
func (r *CommandRunner) ExecuteLogin(password string) (string, error) {
	cmd, err := cmd.LoginConfigFunc(password)
	if err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}
	return cmd, nil
}

// ExecuteConfigUpdate handles the config update command
func (r *CommandRunner) ExecuteConfigUpdate(host, port, username string, useHttps bool) (string, error) {
	configuration := config.ArgsConfig{
		Address:  host,
		Username: username,
		HTTPS:    useHttps,
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
func (r *CommandRunner) ExecuteExploitCommand(subcommand string) (string, error) {
	switch subcommand {
	case "list":
		return r.ExecuteCommand(os.Args[0], append([]string{"exploit", "list"}, defaultFlags...)...)
	default:
		return "", fmt.Errorf("unknown exploit subcommand: %s", subcommand)
	}
}

// ExecuteExploitRun handles running an exploit
func (r *CommandRunner) ExecuteExploitRun(exploitPath, servicePort string, detach bool, tickTime, threadCount string) (string, error) {
	// Parse the string arguments into the required types
	var tickTimeInt, threadCountInt int
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

	if servicePort != "" {
		port, err := strconv.Atoi(servicePort)
		if err != nil {
			return "", fmt.Errorf("invalid service port: %s", servicePort)
		}
		servicePortUint16 = uint16(port)
	}

	return cmd.AttackFunc(exploitPath, tickTimeInt, threadCountInt, servicePortUint16, detach)
}

// ExecuteExploitCreate handles creating an exploit template
func (r *CommandRunner) ExecuteExploitCreate(name string) (string, error) {
	return cmd.CreateFunc(name)
}

// ExecuteExploitRemove handles removing an exploit template
func (r *CommandRunner) ExecuteExploitRemove(name string) (string, error) {
	return cmd.RemoveFunc(name)
}

// ExecuteExploitStop handles stopping a running exploit
func (r *CommandRunner) ExecuteExploitStop(pid string) (string, error) {
	// Validate that pid is a number
	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		return "", fmt.Errorf("invalid process ID: %s", pid)
	}
	return cmd.StopFunc(pidInt)
}

// ExecuteWithTimeout executes a command with a timeout
func (r *CommandRunner) ExecuteWithTimeout(timeout time.Duration, command string, args ...string) (string, error) {
	// Create the command
	cmd := exec.Command(command, args...)

	// Create buffers for stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Start the command
	if err := cmd.Start(); err != nil {
		return "", err
	}

	// Create a channel for the command to finish
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	// Wait for either the command to finish or the timeout
	select {
	case <-time.After(timeout):
		if err := cmd.Process.Kill(); err != nil {
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
func (r *CommandRunner) ExecuteBackgroundCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	// Start the command without waiting for it to complete
	if err := cmd.Start(); err != nil {
		return "", err
	}

	// Return the PID for tracking
	return fmt.Sprintf("Background process started with PID: %d", cmd.Process.Pid), nil
}
