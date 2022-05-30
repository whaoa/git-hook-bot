package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

func New(level, timeFormat string) zerolog.Logger {
	// set timestamp format
	if timeFormat == "" {
		timeFormat = "2006-01-02 15:04:05"
	}
	zerolog.TimeFieldFormat = timeFormat

	// set logger level
	var l zerolog.Level
	switch strings.ToUpper(level) {
	case "ERROR":
		l = zerolog.ErrorLevel
	case "WARN":
		l = zerolog.WarnLevel
	case "INFO":
		l = zerolog.InfoLevel
	case "DEBUG", "TRACE":
		l = zerolog.DebugLevel
	default:
		l = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(l)

	loggerContext := zerolog.
		New(os.Stdout).
		With().
		Timestamp().
		Stack()

	// add caller field when level is debug
	if l == zerolog.DebugLevel {
		loggerContext = loggerContext.Caller()
	}

	return loggerContext.Logger()
}
