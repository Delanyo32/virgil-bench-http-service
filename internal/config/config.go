package config

import (
	"time"
)

// Config holds all application configuration.
type Config struct {
	// Server
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	// Database
	DatabaseURL    string
	MaxOpenConns   int
	MaxIdleConns   int
	ConnMaxLifetime time.Duration

	// Redis
	RedisAddr string

	// Auth
	JWTSecret string

	// Payment
	PaymentGatewayURL string

	// SMTP
	SMTPHost string
	SMTPPort int

	// Worker
	WorkerCount int

	// Logging
	LogLevel string
}

// Load reads configuration from environment variables with defaults.
// FLAW: magic numbers -- hardcoded default values scattered throughout.
func Load() *Config {
	return &Config{
		Port:              getEnvInt("PORT", 8080),                      // FLAW: magic number
		ReadTimeout:       time.Duration(getEnvInt("READ_TIMEOUT", 30)) * time.Second,
		WriteTimeout:      time.Duration(getEnvInt("WRITE_TIMEOUT", 30)) * time.Second,
		DatabaseURL:       getEnvStr("DATABASE_URL", "postgres://localhost:5432/ordersvc?sslmode=disable"),
		MaxOpenConns:      getEnvInt("MAX_OPEN_CONNS", 25),             // FLAW: magic number
		MaxIdleConns:      getEnvInt("MAX_IDLE_CONNS", 5),              // FLAW: magic number
		ConnMaxLifetime:   time.Duration(getEnvInt("CONN_MAX_LIFETIME", 300)) * time.Second,
		RedisAddr:         getEnvStr("REDIS_ADDR", "localhost:6379"),
		JWTSecret:         getEnvStr("JWT_SECRET", "super-secret-key-change-me"), // FLAW: hardcoded default secret
		PaymentGatewayURL: getEnvStr("PAYMENT_GATEWAY_URL", "http://localhost:9090"),
		SMTPHost:          getEnvStr("SMTP_HOST", "localhost"),
		SMTPPort:          getEnvInt("SMTP_PORT", 587),
		WorkerCount:       getEnvInt("WORKER_COUNT", 10),               // FLAW: magic number
		LogLevel:          getEnvStr("LOG_LEVEL", "info"),
	}
}

// Validate checks that required configuration values are set.
func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return errMissingConfig("DATABASE_URL")
	}
	if c.JWTSecret == "" || c.JWTSecret == "super-secret-key-change-me" {
		return errMissingConfig("JWT_SECRET")
	}
	return nil
}

type configError struct {
	field string
}

func (e *configError) Error() string {
	return "missing required configuration: " + e.field
}

func errMissingConfig(field string) error {
	return &configError{field: field}
}
