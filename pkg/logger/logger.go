package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Logger provides structured logging for the application.
type Logger struct {
	level  string
	output *log.Logger
}

// New creates a new Logger with the given level.
func New(level string) *Logger {
	return &Logger{
		level:  level,
		output: log.New(os.Stdout, "", 0),
	}
}

// Info logs an informational message with key-value pairs.
func (l *Logger) Info(msg string, kvs ...interface{}) {
	l.log("INFO", msg, kvs...)
}

// Error logs an error message with key-value pairs.
func (l *Logger) Error(msg string, kvs ...interface{}) {
	l.log("ERROR", msg, kvs...)
}

// Debug logs a debug message with key-value pairs.
func (l *Logger) Debug(msg string, kvs ...interface{}) {
	if l.level != "debug" {
		return
	}
	l.log("DEBUG", msg, kvs...)
}

// Warn logs a warning message with key-value pairs.
func (l *Logger) Warn(msg string, kvs ...interface{}) {
	l.log("WARN", msg, kvs...)
}

// log formats and outputs a log line.
func (l *Logger) log(level, msg string, kvs ...interface{}) {
	ts := time.Now().Format(time.RFC3339)
	fields := ""
	for i := 0; i+1 < len(kvs); i += 2 {
		fields += fmt.Sprintf(" %v=%v", kvs[i], kvs[i+1])
	}
	l.output.Printf("[%s] %s: %s%s", ts, level, msg, fields)
}
