package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"unique"

	"github.com/ByteTheCookies/backend/internal/utils"
)

const (
	DebugLevel = iota
	InfoLevel
	SuccessLevel
	WarningLevel
	ErrorLevel
	FatalLevel
)

type Logger struct {
	Level         int
	lastLogged    unique.Handle[string]
	debugLogger   *log.Logger
	infoLogger    *log.Logger
	succeLogger   *log.Logger
	warnLogger    *log.Logger
	errorLogger   *log.Logger
	IncludeCaller bool // <- toggle per includere info sul chiamante
}

var logger *Logger
var logFile *os.File

func init() {
	multiWriter := io.Writer(os.Stdout)

	logger = &Logger{
		Level:         InfoLevel,
		IncludeCaller: true, // Abilitato di default
		debugLogger:   log.New(multiWriter, utils.Gray+"[=] DEBUG: "+utils.White, 0),
		infoLogger:    log.New(multiWriter, utils.Cyan+"[*] INFO: "+utils.White, 0),
		succeLogger:   log.New(multiWriter, utils.Green+"[+] SUCCESS: "+utils.White, 0),
		warnLogger:    log.New(multiWriter, utils.Yellow+"[/] WARN: "+utils.White, 0),
		errorLogger:   log.New(multiWriter, utils.Red+"[//] ERROR: "+utils.White, 0),
	}
}

func CloseLogFile() {
	if logFile != nil {
		logFile.Close()
	}
}

func ParseLevel(level string) int {
	switch strings.ToLower(level) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "success":
		return SuccessLevel
	case "warning":
		return WarningLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	default:
		return InfoLevel
	}
}

func SetLevel(level int) {
	logger.Level = level
}

func logWithCaller(l *log.Logger, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if logger.IncludeCaller {
		// Recupera chiamante
		_, file, line, ok := runtime.Caller(2)
		if ok {
			shortFile := file[strings.LastIndex(file, "/")+1:]
			msg = fmt.Sprintf("[%s:%d] %s", shortFile, line, msg)
		}
	}
	l.Println(msg)
}

func Debug(format string, args ...interface{}) {
	if logger.Level <= DebugLevel {
		logWithCaller(logger.debugLogger, format, args...)
	}
}

func Info(format string, args ...interface{}) {
	if logger.Level <= InfoLevel {
		logWithCaller(logger.infoLogger, format, args...)
	}
}

func Success(format string, args ...interface{}) {
	if logger.Level <= SuccessLevel {
		logWithCaller(logger.succeLogger, format, args...)
	}
}

func Warning(format string, args ...interface{}) {
	if logger.Level <= WarningLevel {
		logWithCaller(logger.warnLogger, format, args...)
	}
}

func Error(format string, args ...interface{}) {
	if logger.Level <= ErrorLevel {
		logWithCaller(logger.errorLogger, format, args...)
	}
}

func Fatal(format string, args ...interface{}) {
	if logger.Level <= FatalLevel {
		logWithCaller(logger.errorLogger, format, args...)
		os.Exit(1)
	}
}
