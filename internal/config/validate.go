package config

import (
	"fmt"
	"strings"
)

// ValidationError collects multiple configuration validation failures.
type ValidationError struct {
	Errors []string
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("config validation failed: %s", strings.Join(e.Errors, "; "))
}

// HasErrors returns true if any validation errors were recorded.
func (e *ValidationError) HasErrors() bool {
	return len(e.Errors) > 0
}

// ValidateConfig checks the Config for required fields and valid value ranges.
func ValidateConfig(cfg Config) error {
	ve := &ValidationError{}

	if cfg.Port < 1 || cfg.Port > 65535 {
		ve.Errors = append(ve.Errors, fmt.Sprintf("port must be between 1 and 65535, got %d", cfg.Port))
	}

	if strings.TrimSpace(cfg.DatabaseURL) == "" {
		ve.Errors = append(ve.Errors, "database_url is required")
	}

	if cfg.ReadTimeout <= 0 {
		ve.Errors = append(ve.Errors, "read_timeout must be positive")
	}

	if cfg.WriteTimeout <= 0 {
		ve.Errors = append(ve.Errors, "write_timeout must be positive")
	}

	if cfg.MaxOpenConns < 1 {
		ve.Errors = append(ve.Errors, "max_open_conns must be at least 1")
	}

	if cfg.MaxIdleConns < 0 {
		ve.Errors = append(ve.Errors, "max_idle_conns must not be negative")
	}

	if cfg.MaxIdleConns > cfg.MaxOpenConns {
		ve.Errors = append(ve.Errors, "max_idle_conns must not exceed max_open_conns")
	}

	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[strings.ToLower(cfg.LogLevel)] {
		ve.Errors = append(ve.Errors, fmt.Sprintf("invalid log_level: %s", cfg.LogLevel))
	}

	if ve.HasErrors() {
		return ve
	}
	return nil
}
