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
		return r.ExecuteCommand(os.Args[0], append([]string{"config", "reset"}, defaultFlags...)...)
	case "logout":
		return r.ExecuteCommand(os.Args[0], append([]string{"config", "logout"}, defaultFlags...)...)
	default:
		return "", fmt.Errorf("unknown config subcommand: %s", subcommand)
	}
}

// ExecuteLogin handles the login command
func (r *CommandRunner) ExecuteLogin(username, password string) (string, error) {
	cmd := exec.Command(os.Args[0], append([]string{"config", "login"}, defaultFlags...)...)

	// Create input for sending username and password
	input := username + "\n" + password + "\n"

	// Create buffers for stdin, stdout and stderr
	var stdin bytes.Buffer
	var stdout, stderr bytes.Buffer

	// Write to stdin
	stdin.WriteString(input)

	cmd.Stdin = &stdin
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

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

// ExecuteConfigUpdate handles the config update command
func (r *CommandRunner) ExecuteConfigUpdate(host, port, username string, useHttps bool) (string, error) {
	args := append([]string{"config", "update"}, defaultFlags...)

	if host != "" {
		args = append(args, "--host", host)
	}

	if port != "" {
		args = append(args, "--port", port)
	}

	if username != "" {
		args = append(args, "--username", username)
	}

	if useHttps {
		args = append(args, "--https", "true")
	} else {
		args = append(args, "--https", "false")
	}

	return r.ExecuteCommand(os.Args[0], args...)
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
	args := append([]string{"exploit", "run"}, defaultFlags...)
	args = append(args, "--exploit", exploitPath, "--port", servicePort)

	if detach {
		args = append(args, "--detach")
	}

	if tickTime != "" {
		args = append(args, "--tick", tickTime)
	}

	if threadCount != "" {
		args = append(args, "--thread", threadCount)
	}

	return r.ExecuteCommand(os.Args[0], args...)
}

// ExecuteExploitCreate handles creating an exploit template
func (r *CommandRunner) ExecuteExploitCreate(name string) (string, error) {
	return r.ExecuteCommand(os.Args[0], append([]string{"exploit", "create"}, append(defaultFlags, "--name", name)...)...)
}

// ExecuteExploitRemove handles removing an exploit template
func (r *CommandRunner) ExecuteExploitRemove(name string) (string, error) {
	return r.ExecuteCommand(os.Args[0], append([]string{"exploit", "remove"}, append(defaultFlags, "--name", name)...)...)
}

// ExecuteExploitStop handles stopping a running exploit
func (r *CommandRunner) ExecuteExploitStop(pid string) (string, error) {
	// Validate that pid is a number
	_, err := strconv.Atoi(pid)
	if err != nil {
		return "", fmt.Errorf("invalid process ID: %s", pid)
	}

	return r.ExecuteCommand(os.Args[0], append([]string{"exploit", "stop"}, append(defaultFlags, "--pid", pid)...)...)
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
