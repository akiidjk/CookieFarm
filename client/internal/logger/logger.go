// Package logger provides functions to manage the CookieFarm client logging.
package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ByteTheCookies/cookieclient/internal/utils"
	"github.com/rs/zerolog"
)

var (
	// LogLevel defines the log level for the logger.
	LogLevel zerolog.Level
	// Log is the logger instance for the CookieFarm client.
	Log zerolog.Logger
	// logFile represents the log file for the CookieFarm client.
	logFile *os.File
)

// Setup configures the logger with the specified log level and returns the log file path.
func Setup(level string) string {
	parsedLevel, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		parsedLevel = zerolog.InfoLevel
	}
	LogLevel = parsedLevel

	_ = os.MkdirAll("./logs", 0o755)
	logPath := filepath.Join("logs", "clientfarm-"+strconv.Itoa(int(time.Now().UnixMilli()))) + ".log"

	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
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
				return utils.Gray + "[DEBUG]" + utils.Reset
			case "info":
				return utils.Cyan + "[INFO]" + utils.Reset
			case "warn":
				return utils.Yellow + "[WARN]" + utils.Reset
			case "error":
				return utils.Red + "[ERROR]" + utils.Reset
			case "fatal":
				return utils.Magenta + "[FATAL]" + utils.Reset
			default:
				return lvl
			}
		},
		FormatMessage: func(i any) string {
			return fmt.Sprintf("%s", i)
		},
		FormatFieldName: func(i any) string {
			return utils.White + fmt.Sprintf("%s=", i)
		},
		FormatFieldValue: func(i any) string {
			return fmt.Sprintf("%v", i) + utils.Reset
		},
	}

	multi := zerolog.MultiLevelWriter(consoleWriter, logFile)

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
	if logFile != nil {
		_ = logFile.Close()
	}
}
