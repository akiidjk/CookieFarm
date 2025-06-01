package tui

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ByteTheCookies/cookieclient/internal/logger"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog"
)

// TUIWriter is a custom writer for capturing log output in the TUI
type TUIWriter struct {
	mu       sync.Mutex
	messages []string
	maxLines int
}

// NewTUIWriter creates a new TUI log writer
func NewTUIWriter(maxLines int) *TUIWriter {
	return &TUIWriter{
		maxLines: maxLines,
	}
}

// Write implements io.Writer
func (w *TUIWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.messages = append(w.messages, string(p))
	if len(w.messages) > w.maxLines {
		w.messages = w.messages[1:]
	}

	return len(p), nil
}

// GetMessages returns all captured log messages
func (w *TUIWriter) GetMessages() []string {
	w.mu.Lock()
	defer w.mu.Unlock()

	return append([]string{}, w.messages...)
}

// ClearMessages clears all captured log messages
func (w *TUIWriter) ClearMessages() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.messages = []string{}
}

// SetupTUILogger initializes a logger that writes to both a file and the TUI
func SetupTUILogger(level string) *TUIWriter {
	tuiWriter := NewTUIWriter(100)

	// Create a multi-writer for zerolog
	writers := []io.Writer{tuiWriter}

	// Try to add file logger if possible
	logPath := logger.Setup(level)

	// Ensure we have a separate log file for TUI
	tuiLogDir := filepath.Join(os.TempDir(), "cookieclient-tui")
	_ = os.MkdirAll(tuiLogDir, 0755)
	tuiLogFile, err := os.OpenFile(
		filepath.Join(tuiLogDir, "tui.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)

	if err == nil {
		writers = append(writers, tuiLogFile)
	}

	if logger.LogFile != nil {
		writers = append(writers, logger.LogFile)
	}

	multiWriter := zerolog.MultiLevelWriter(writers...)

	// Set log level
	var logLevel zerolog.Level
	switch level {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	default:
		logLevel = zerolog.InfoLevel
	}

	// Setup the logger - use a separate logger instance to avoid
	// interfering with the main application logger
	tuiLogger := zerolog.New(multiWriter).
		With().
		Timestamp().
		Str("component", "tui").
		Logger().
		Level(logLevel)

	// Log initialization
	tuiLogger.Info().Str("logPath", logPath).Msg("TUI logger initialized")

	return tuiWriter
}

// FormatLogMessages formats log messages with appropriate colors
func FormatLogMessages(messages []string) string {
	var formattedLines []string

	for _, msg := range messages {
		// Skip empty messages
		if msg == "" {
			continue
		}

		switch {
		case containsAny(msg, []string{"error", "Error", "ERROR", "fail", "Fail", "FAIL"}):
			formattedLines = append(formattedLines, lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF0000")).
				Bold(true).
				Render(msg))
		case containsAny(msg, []string{"warning", "Warning", "WARNING", "warn", "Warn", "WARN"}):
			formattedLines = append(formattedLines, lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFC107")).
				Bold(true).
				Render(msg))
		case containsAny(msg, []string{"info", "Info", "INFO"}):
			formattedLines = append(formattedLines, lipgloss.NewStyle().
				Foreground(lipgloss.Color("#2196F3")).
				Bold(true).
				Render(msg))
		case containsAny(msg, []string{"debug", "Debug", "DEBUG"}):
			formattedLines = append(formattedLines, lipgloss.NewStyle().Faint(true).Render(msg))
		default:
			formattedLines = append(formattedLines, msg)
		}
	}

	if len(formattedLines) == 0 {
		return ""
	}

	return lipgloss.JoinVertical(lipgloss.Left, formattedLines...)
}

// containsAny checks if the string contains any of the substrings
func containsAny(s string, substrings []string) bool {
	for _, sub := range substrings {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// GetLogWriter returns a writer for logging that can be used with zerolog
func GetLogWriter() io.Writer {
	return NewTUIWriter(100)
}
