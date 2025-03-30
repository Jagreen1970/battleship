package app

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Log      LogConfig
}

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	Driver   string
	URL      string
	Timeout  time.Duration
	Name     string
	User     string
	Password string
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port     int
	Timeout  time.Duration
	LogLevel string
}

// LogConfig holds logging-specific configuration
type LogConfig struct {
	Level  string
	Format string
	Output string
}

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Driver:  "mongo",
			URL:     "mongodb://localhost:27017",
			Name:    "battleship",
			Timeout: 5 * time.Second,
		},
		Server: ServerConfig{
			Port:    8080,
			Timeout: 5 * time.Second,
		},
	}
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := DefaultConfig()

	// Database configuration
	if driver := os.Getenv("DB_DRIVER"); driver != "" {
		cfg.Database.Driver = driver
	}
	if url := os.Getenv("DB_URL"); url != "" {
		cfg.Database.URL = url
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		cfg.Database.Name = name
	}
	if timeout := os.Getenv("DB_TIMEOUT"); timeout != "" {
		duration, err := time.ParseDuration(timeout)
		if err != nil {
			return nil, err
		}
		cfg.Database.Timeout = duration
	}
	if user := os.Getenv("DB_USER"); user != "" {
		cfg.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		cfg.Database.Password = password
	}

	// Server configuration
	if port := os.Getenv("SERVER_PORT"); port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			return nil, err
		}
		cfg.Server.Port = p
	}
	if timeout := os.Getenv("SERVER_TIMEOUT"); timeout != "" {
		duration, err := time.ParseDuration(timeout)
		if err != nil {
			return nil, err
		}
		cfg.Server.Timeout = duration
	}
	if logLevel := os.Getenv("SERVER_LOG_LEVEL"); logLevel != "" {
		cfg.Server.LogLevel = logLevel
	}

	// Log configuration
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		cfg.Log.Level = level
	}
	if format := os.Getenv("LOG_FORMAT"); format != "" {
		cfg.Log.Format = format
	}
	if output := os.Getenv("LOG_OUTPUT"); output != "" {
		cfg.Log.Output = output
	}

	return cfg, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Database.Driver == "" {
		return fmt.Errorf("database driver is required")
	}
	if c.Database.URL == "" {
		return fmt.Errorf("database URL is required")
	}
	if c.Database.Timeout <= 0 {
		return fmt.Errorf("database timeout must be positive")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("database name is required")
	}
	if c.Server.Port <= 0 {
		return fmt.Errorf("server port must be positive")
	}
	if c.Server.Timeout <= 0 {
		return fmt.Errorf("server timeout must be positive")
	}
	return nil
}
