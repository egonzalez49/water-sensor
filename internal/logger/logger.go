package logger

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

type Level int8

const (
	LevelDebug Level = iota
	LevelInfo
	LevelError
	LevelFatal
	LevelOff
)

// Return a human-friendly string for the severity level.
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

// Return the corresponding severity level from the human-friendly string.
// Defaults to off given no matching level.
func LevelFromString(value string) Level {
	switch strings.ToUpper(value) {
	case "DEBUG":
		return LevelDebug
	case "INFO":
		return LevelInfo
	case "ERROR":
		return LevelError
	case "FATAL":
		return LevelFatal
	default:
		return LevelOff
	}
}

type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

// Return a new Logger instance which writes log entries at or above a minimum severity
// level to a specific output destination.
func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

// Declare some helper methods for writing log entries at the different levels.
func (l *Logger) Debug(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

func (l *Logger) Info(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

func (l *Logger) Error(err error, properties map[string]string) {
	l.print(LevelError, err.Error(), properties)
}

func (l *Logger) Fatal(err error, properties map[string]string) {
	l.print(LevelFatal, err.Error(), properties)
	os.Exit(1) // For entries at the FATAL level, we also terminate the application.
}

func (l *Logger) print(level Level, message string, properties map[string]string) (int, error) {
	// If the severity level of the log entry is below the minimum severity for the
	// logger, then return with no further action.
	if level < l.minLevel {
		return 0, nil
	}

	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	// Include a stack trace for entries at the ERROR and FATAL levels.
	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	var line []byte
	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message: " + err.Error())
	}

	// Lock the mutex so that no two writes to
	// the output destination can happen concurrently.
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(line, '\n'))
}

// Implement a Write() method to satisfy the
// io.Writer interface. This writes a log entry
// at the ERROR level with no additional properties.
func (l *Logger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}
