package tui

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/ByteTheCookies/cookieclient/internal/logger"
	"github.com/charmbracelet/lipgloss"
)

// GetTerminalSize returns the terminal's width and height
func GetTerminalSize() (width, height int) {
	// Default values if detection fails
	width, height = 80, 24

	// Try to get terminal size using stty command
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err == nil {
		fmt.Sscanf(string(out), "%d %d", &height, &width)
	}

	return width, height
}

// FormatCommand formats a command string for display
func FormatCommand(command string) string {
	return HighlightStyle(command)
}

// FormatOutput formats output text with appropriate styling
func FormatOutput(text string) string {
	// Add line numbers
	var formattedLines []string
	scanner := bufio.NewScanner(strings.NewReader(text))
	lineNum := 1

	for scanner.Scan() {
		line := scanner.Text()

		// Highlight errors and infos
		loweredLine := strings.ToLower(line)
		switch {
		case strings.Contains(loweredLine, "error"):
			line = ErrorText(line)
		case strings.Contains(loweredLine, "info"):
			line = InfoText(line)
		case strings.Contains(loweredLine, "warn"):
			line = WarningText(line)
		}

		formattedLines = append(formattedLines, fmt.Sprintf("%3d │ %s", lineNum, line))
		lineNum++
	}

	return lipgloss.JoinVertical(lipgloss.Left, formattedLines...)
}

// CaptureOutput captures the output of stdout during execution of the given function
func CaptureOutput(fn func()) string {
	// Save and restore original stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Capture output from logger too if possible
	if logger.LogFile != nil {
		logger.LogFile.Close()
		tempFile, _ := os.CreateTemp("", "cookieclient-log-*.tmp")
		logger.LogFile = tempFile
		defer func() {
			tempFile.Close()
			os.Remove(tempFile.Name())
		}()
	}

	// Execute the function
	fn()

	// Restore stdout and collect output
	w.Close()
	os.Stdout = oldStdout

	var buf strings.Builder
	io.Copy(&buf, r)

	return buf.String()
}

// OpenEditor opens a text file in the user's preferred editor
func OpenEditor(filename string) error {
	// Try to get preferred editor from environment
	editor := os.Getenv("EDITOR")
	if editor == "" {
		// Default editors based on OS
		switch runtime.GOOS {
		case "windows":
			editor = "notepad"
		case "darwin":
			editor = "open"
		default:
			// Try common editors on Linux/Unix
			for _, ed := range []string{"nano", "vim", "vi", "emacs"} {
				if _, err := exec.LookPath(ed); err == nil {
					editor = ed
					break
				}
			}
			if editor == "" {
				return errors.New("no suitable text editor found")
			}
		}
	}

	cmd := exec.Command(editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// LoadConfigToForms loads the current configuration into form inputs
func LoadConfigToForms() error {
	err := config.LoadLocalConfig()
	if err != nil {
		return err
	}

	// Configuration is now available in config.ArgsConfigInstance
	return nil
}

// FormatError returns a formatted error message
func FormatError(err error) string {
	if err == nil {
		return ""
	}

	return ErrorStyle.Render("Error: " + err.Error())
}

// FormatSuccess returns a formatted success message
func FormatSuccess(message string) string {
	return SuccessStyle.Render(message)
}
