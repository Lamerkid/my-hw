package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota + 1
	INFO
	WARN
	ERROR
)

type Logger struct {
	Level  LogLevel
	output io.Writer
}

func NewLogger(level string) *Logger {
	switch level {
	case "DEBUG":
		return &Logger{Level: DEBUG, output: os.Stdout}
	case "INFO":
		return &Logger{Level: INFO, output: os.Stdout}
	case "WARN":
		return &Logger{Level: WARN, output: os.Stdout}
	case "ERROR":
		return &Logger{Level: ERROR, output: os.Stdout}
	default:
		return &Logger{Level: INFO, output: os.Stdout}
	}
}

func (l *Logger) log(level LogLevel, levelName, msg string, args ...any) {
	if level < l.Level {
		return
	}
	timestamp := time.Now().UTC().Format(time.RFC3339)

	formattedMsg := fmt.Sprintf(msg, args...)

	fmt.Fprintf(l.output, "%s [%s]: %s\n", timestamp, levelName, formattedMsg)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.log(DEBUG, "DEBUG", msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.log(INFO, "INFO", msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.log(WARN, "WARN", msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.log(ERROR, "ERROR", msg, args...)
}
