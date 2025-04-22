package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ByteTheCookies/cookiefarm-client/internal/utils"
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

	_ = os.MkdirAll("./logs", 0755)
	logPath := filepath.Join("logs", "clientfarm-"+time.Now().Format("20060102-150405")+".log")

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

	Log = zerolog.New(multi).With().Timestamp().Caller().Logger()
	zerolog.SetGlobalLevel(parsedLevel)
}

func Close() {
	if logFile != nil {
		_ = logFile.Close()
	}
}
