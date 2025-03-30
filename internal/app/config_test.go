package app

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, "mongo", cfg.Database.Driver)
	assert.Equal(t, "mongodb://localhost:27017", cfg.Database.URL)
	assert.Equal(t, "battleship", cfg.Database.Name)
	assert.Equal(t, 5*time.Second, cfg.Database.Timeout)

	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, 5*time.Second, cfg.Server.Timeout)
}

func TestLoadConfig(t *testing.T) {
	// Set environment variables
	os.Setenv("DB_DRIVER", "mongo")
	os.Setenv("DB_URL", "mongodb://testhost:27017")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_TIMEOUT", "10s")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("SERVER_TIMEOUT", "15s")
	defer func() {
		os.Unsetenv("DB_DRIVER")
		os.Unsetenv("DB_URL")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_TIMEOUT")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("SERVER_TIMEOUT")
	}()

	cfg, err := LoadConfig()
	assert.NoError(t, err)

	assert.Equal(t, "mongo", cfg.Database.Driver)
	assert.Equal(t, "mongodb://testhost:27017", cfg.Database.URL)
	assert.Equal(t, "testdb", cfg.Database.Name)
	assert.Equal(t, 10*time.Second, cfg.Database.Timeout)

	assert.Equal(t, 9090, cfg.Server.Port)
	assert.Equal(t, 15*time.Second, cfg.Server.Timeout)
}

func TestLoadConfigInvalidValues(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectedErr string
	}{
		{
			name: "invalid_database_timeout",
			envVars: map[string]string{
				"DB_TIMEOUT": "invalid",
			},
			expectedErr: "time: invalid duration \"invalid\"",
		},
		{
			name: "invalid_server_port",
			envVars: map[string]string{
				"SERVER_PORT": "invalid",
			},
			expectedErr: "strconv.Atoi: parsing \"invalid\": invalid syntax",
		},
		{
			name: "invalid_server_timeout",
			envVars: map[string]string{
				"SERVER_TIMEOUT": "invalid",
			},
			expectedErr: "time: invalid duration \"invalid\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				for k := range tt.envVars {
					os.Unsetenv(k)
				}
			}()

			_, err := LoadConfig()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError bool
	}{
		{
			name: "valid_config",
			config: Config{
				Database: DatabaseConfig{
					Driver:  "mongo",
					URL:     "mongodb://localhost:27017",
					Name:    "testdb",
					Timeout: 5 * time.Second,
				},
				Server: ServerConfig{
					Port:    8080,
					Timeout: 5 * time.Second,
				},
			},
			expectError: false,
		},
		{
			name: "missing_database_driver",
			config: Config{
				Database: DatabaseConfig{
					URL:     "mongodb://localhost:27017",
					Name:    "testdb",
					Timeout: 5 * time.Second,
				},
			},
			expectError: true,
		},
		{
			name: "missing_database_URL",
			config: Config{
				Database: DatabaseConfig{
					Driver:  "mongo",
					Name:    "testdb",
					Timeout: 5 * time.Second,
				},
			},
			expectError: true,
		},
		{
			name: "invalid_database_timeout",
			config: Config{
				Database: DatabaseConfig{
					Driver:  "mongo",
					URL:     "mongodb://localhost:27017",
					Name:    "testdb",
					Timeout: 0,
				},
			},
			expectError: true,
		},
		{
			name: "missing_database_name",
			config: Config{
				Database: DatabaseConfig{
					Driver:  "mongo",
					URL:     "mongodb://localhost:27017",
					Timeout: 5 * time.Second,
				},
			},
			expectError: true,
		},
		{
			name: "invalid_server_port",
			config: Config{
				Database: DatabaseConfig{
					Driver:  "mongo",
					URL:     "mongodb://localhost:27017",
					Name:    "testdb",
					Timeout: 5 * time.Second,
				},
				Server: ServerConfig{
					Port:    0,
					Timeout: 5 * time.Second,
				},
			},
			expectError: true,
		},
		{
			name: "invalid_server_timeout",
			config: Config{
				Database: DatabaseConfig{
					Driver:  "mongo",
					URL:     "mongodb://localhost:27017",
					Name:    "testdb",
					Timeout: 5 * time.Second,
				},
				Server: ServerConfig{
					Port:    8080,
					Timeout: 0,
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
