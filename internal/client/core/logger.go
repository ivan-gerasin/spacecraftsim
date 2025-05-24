package core

import (
	"fmt"
	"io"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	LevelInfo LogLevel = iota
	LevelSuccess
	LevelError
)

// Logger handles message logging
type Logger struct {
	writer io.Writer
}

// NewLogger creates a new logger
func NewLogger(writer io.Writer) *Logger {
	return &Logger{writer: writer}
}

// Log writes a message to the log
func (l *Logger) Log(level LogLevel, format string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	message := fmt.Sprintf(format, args...)

	var prefix string
	switch level {
	case LevelSuccess:
		prefix = "✓"
	case LevelError:
		prefix = "✗"
	default:
		prefix = "•"
	}

	fmt.Fprintf(l.writer, "[%s] %s %s\n", timestamp, prefix, message)
}

// LogInfo logs an informational message
func (l *Logger) LogInfo(format string, args ...interface{}) {
	l.Log(LevelInfo, format, args...)
}

// LogSuccess logs a success message
func (l *Logger) LogSuccess(format string, args ...interface{}) {
	l.Log(LevelSuccess, format, args...)
}

// LogError logs an error message
func (l *Logger) LogError(format string, args ...interface{}) {
	l.Log(LevelError, format, args...)
}

// SetOutput sets the output writer for the logger
func (l *Logger) SetOutput(writer io.Writer) {
	l.writer = writer
}
