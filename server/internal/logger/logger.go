package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/ByteTheCookies/cookieserver/internal/utils"
	"github.com/rs/zerolog"
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
