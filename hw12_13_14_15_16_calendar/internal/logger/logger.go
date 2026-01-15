package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

type Logger struct {
	Level  LogLevel
	output io.Writer
}

func New(level string) *Logger {
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

func (l *Logger) log(level LogLevel, levelName, msg string) {
	if level < l.Level {
		return
	}
	timestamp := time.Now().UTC().Format(time.RFC3339)

	fmt.Fprintf(l.output, "%s [%s]: %s\n", timestamp, levelName, msg)
}

func (l *Logger) Debug(msg string) {
	l.log(DEBUG, "DEBUG", msg)
}

func (l *Logger) Info(msg string) {
	l.log(INFO, "INFO", msg)
}

func (l *Logger) Warn(msg string) {
	l.log(WARN, "WARN", msg)
}

func (l *Logger) Error(msg string) {
	l.log(ERROR, "ERROR", msg)
}
