package config

import (
	"time"
)

// Default configuration values.
const (
	DefaultPort          = 8080
	DefaultReadTimeout   = 30 * time.Second
	DefaultWriteTimeout  = 30 * time.Second
	DefaultMaxOpenConns  = 25
	DefaultMaxIdleConns  = 5
	DefaultLogLevel      = "info"
	DefaultShutdownGrace = 15 * time.Second
	DefaultWorkerCount   = 10
)

// DefaultConfig returns a Config populated with default values.
func DefaultConfig() Config {
	return Config{
		Port:         DefaultPort,
		ReadTimeout:  DefaultReadTimeout,
		WriteTimeout: DefaultWriteTimeout,
		MaxOpenConns: DefaultMaxOpenConns,
		MaxIdleConns: DefaultMaxIdleConns,
		LogLevel:     DefaultLogLevel,
		WorkerCount:  DefaultWorkerCount,
	}
}

// WithPort returns a copy of the config with the port overridden.
func (c Config) WithPort(port int) Config {
	c.Port = port
	return c
}

// WithDatabaseURL returns a copy of the config with the database URL set.
func (c Config) WithDatabaseURL(url string) Config {
	c.DatabaseURL = url
	return c
}

// WithLogLevel returns a copy of the config with the log level set.
func (c Config) WithLogLevel(level string) Config {
	c.LogLevel = level
	return c
}

// WithWorkerCount returns a copy of the config with the worker count set.
func (c Config) WithWorkerCount(count int) Config {
	c.WorkerCount = count
	return c
}
