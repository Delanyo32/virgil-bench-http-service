package logger

import (
	"fmt"
	"strings"
)

// Level represents a log severity level.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// String returns the level name as an uppercase string.
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// ParseLevel converts a string level name to a Level value.
// Returns an error for unrecognized level names.
func ParseLevel(s string) (Level, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return LevelDebug, nil
	case "info":
		return LevelInfo, nil
	case "warn", "warning":
		return LevelWarn, nil
	case "error":
		return LevelError, nil
	default:
		return LevelInfo, fmt.Errorf("unknown log level: %q", s)
	}
}

// ShouldLog returns true if the message level meets or exceeds the threshold.
func ShouldLog(messageLevel, threshold Level) bool {
	return messageLevel >= threshold
}

// AllLevels returns all defined log levels in order of severity.
func AllLevels() []Level {
	return []Level{LevelDebug, LevelInfo, LevelWarn, LevelError}
}
