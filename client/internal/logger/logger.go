// Package logger provides functions to manage the CookieFarm client logging.
package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
)

var (
	LogLevel zerolog.Level  // LogLevel defines the log level for the logger.
	Log      zerolog.Logger // Log is the logger instance for the CookieFarm client.
	LogFile  *os.File       // logFile represents the log file for the CookieFarm client.
	useTUI   bool           // NoTUI indicates whether to disable the TUI mode for logging
)

// SetTUI sets the useTUI variable to enable or disable TUI mode for logging.
func SetTUI(value bool) {
	useTUI = value
}

// Setup configures the logger with the specified log level and returns the log file path.
func Setup(level string) string {
	parsedLevel, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		parsedLevel = zerolog.InfoLevel
	}
	LogLevel = parsedLevel

	_ = os.MkdirAll("/tmp/cookielogs", 0o755)
	logPath := filepath.Join("/", "tmp", "cookielogs", "clientfarm-"+strconv.Itoa(int(time.Now().UnixMilli()))) + ".log"

	LogFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_SYNC, 0o666)
	if err != nil {
		panic("cannot create log file: " + err.Error())
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05",
		FormatLevel: func(i any) string {
			lvl := strings.ToLower(fmt.Sprintf("%s", i))
			switch lvl {
			case "debug":
				return Gray + "[DEBUG]" + Reset
			case "info":
				return Cyan + "[INFO]" + Reset
			case "warn":
				return Yellow + "[WARN]" + Reset
			case "error":
				return Red + "[ERROR]" + Reset
			case "fatal":
				return Magenta + "[FATAL]" + Reset
			default:
				return lvl
			}
		},
		FormatMessage: func(i any) string {
			return fmt.Sprintf("%s", i)
		},
		FormatFieldName: func(i any) string {
			return White + fmt.Sprintf("%s=", i)
		},
		FormatFieldValue: func(i any) string {
			return fmt.Sprintf("%v", i) + Reset
		},
	}

	multi := zerolog.MultiLevelWriter(consoleWriter, LogFile)

	if useTUI {
		multi = zerolog.MultiLevelWriter(LogFile)
	}

	if level == "debug" {
		Log = zerolog.New(multi).With().Timestamp().Caller().Logger()
	} else {
		Log = zerolog.New(multi).With().Timestamp().Logger()
	}

	zerolog.SetGlobalLevel(parsedLevel)

	return logPath
}

// Close shuts down the logger by closing the log file if it is open.
func Close() {
	if LogFile != nil {
		_ = LogFile.Close()
	}
}
