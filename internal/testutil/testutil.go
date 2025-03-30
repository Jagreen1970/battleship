package testutil

import (
	"context"
	"time"

	"github.com/Jagreen1970/battleship/internal/app"
)

// NewTestConfig creates a configuration suitable for testing
func NewTestConfig() app.Config {
	return app.Config{
		Server: app.ServerConfig{
			Port:    8080,
			Timeout: 5 * time.Second,
		},
		Database: app.DatabaseConfig{
			Driver:  "mongo",
			URL:     "mongodb://localhost:27017",
			Name:    "battleship_test",
			Timeout: 5 * time.Second,
		},
	}
}

// NewTestContext creates a context with timeout for testing
func NewTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}
