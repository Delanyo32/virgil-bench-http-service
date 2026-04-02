package logger

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Format defines the output format for log lines.
type Format int

const (
	// FormatText produces human-readable log lines.
	FormatText Format = iota
	// FormatJSON produces structured JSON log lines.
	FormatJSON
)

// Formatter converts log entries into formatted strings.
type Formatter struct {
	format    Format
	timestamp bool
}

// NewFormatter creates a Formatter with the given output format.
func NewFormatter(f Format) *Formatter {
	return &Formatter{
		format:    f,
		timestamp: true,
	}
}

// SetTimestamp enables or disables timestamp inclusion.
func (f *Formatter) SetTimestamp(enabled bool) {
	f.timestamp = enabled
}

// FormatLine formats a log entry with level, message, and key-value pairs.
func (f *Formatter) FormatLine(level, msg string, kvs map[string]interface{}) string {
	switch f.format {
	case FormatJSON:
		return f.formatJSON(level, msg, kvs)
	default:
		return f.formatText(level, msg, kvs)
	}
}

// formatText produces a line like: [2024-01-01T00:00:00Z] INFO: message key=value
func (f *Formatter) formatText(level, msg string, kvs map[string]interface{}) string {
	var sb strings.Builder
	if f.timestamp {
		sb.WriteString(fmt.Sprintf("[%s] ", time.Now().UTC().Format(time.RFC3339)))
	}
	sb.WriteString(fmt.Sprintf("%s: %s", level, msg))

	for k, v := range kvs {
		sb.WriteString(fmt.Sprintf(" %s=%v", k, v))
	}
	return sb.String()
}

// formatJSON produces a JSON object with ts, level, msg, and extra fields.
func (f *Formatter) formatJSON(level, msg string, kvs map[string]interface{}) string {
	entry := make(map[string]interface{})
	if f.timestamp {
		entry["ts"] = time.Now().UTC().Format(time.RFC3339)
	}
	entry["level"] = level
	entry["msg"] = msg

	for k, v := range kvs {
		entry[k] = v
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Sprintf(`{"level":"%s","msg":"%s","error":"marshal failed"}`, level, msg)
	}
	return string(data)
}
