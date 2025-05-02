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
	LogLevel zerolog.Level  // Log level for the logger.
	Log      zerolog.Logger // Logger instance for the CookieFarm client.
	logFile  *os.File       // Log file for the CookieFarm client.
)

// Setup initializes the logger with the specified log level.
func Setup(level string) string {
	parsedLevel, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		parsedLevel = zerolog.InfoLevel
	}
	LogLevel = parsedLevel

	_ = os.MkdirAll("./logs", 0755)
	logPath := filepath.Join("logs", "clientfarm-"+strconv.Itoa(int(time.Now().UnixMilli()))) + ".log"

	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("cannot create log file: " + err.Error())
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05",
		FormatLevel: func(i any) string {
			level := strings.ToLower(fmt.Sprintf("%s", i))
			switch level {
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
				return level
			}
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		},
		FormatFieldName: func(i interface{}) string {
			return utils.White + fmt.Sprintf("%s=", i)
		},
		FormatFieldValue: func(i interface{}) string {
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

// Close closes the log file.
func Close() {
	if logFile != nil {
		_ = logFile.Close()
	}
}
