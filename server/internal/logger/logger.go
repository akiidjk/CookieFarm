package logger

import (
	"fmt"
	"os"
	"strings"

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
	LogLevel zerolog.Level
	Log      zerolog.Logger
	logFile  *os.File
)

func Setup(level string) {
	parsedLevel, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		parsedLevel = zerolog.InfoLevel
	}
	LogLevel = parsedLevel

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

	if level == "debug" {
		Log = zerolog.New(consoleWriter).With().Timestamp().Caller().Logger()
	} else {
		Log = zerolog.New(consoleWriter).With().Timestamp().Logger()
	}

	zerolog.SetGlobalLevel(parsedLevel)
}

func Close() {
	if logFile != nil {
		_ = logFile.Close()
	}
}
