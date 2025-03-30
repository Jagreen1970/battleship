# Configuration Guide

This guide explains the configuration system used in the Battleship game project.

## Overview

The configuration system is centralized in the `internal/app` package and provides:
- Type-safe configuration
- Environment variable support
- Default values
- Configuration validation
- Logging configuration

## Configuration Structure

### Main Configuration

```go
type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
    Log      LogConfig
}
```

### Database Configuration

```go
type DatabaseConfig struct {
    Driver   string        // Database driver (e.g., "mongo")
    URL      string        // Database connection URL
    Timeout  time.Duration // Connection timeout
    Name     string        // Database name
    User     string        // Database user
    Password string        // Database password
}
```

### Server Configuration

```go
type ServerConfig struct {
    Port      int           // Server port
    Timeout   time.Duration // Request timeout
    LogLevel  string        // Server log level
}
```

### Log Configuration

```go
type LogConfig struct {
    Level  string // Log level (debug, info, warn, error)
    Format string // Log format (json, text)
    Output string // Log output (stdout, file)
}
```

## Environment Variables

Configuration can be set using environment variables:

### Database Configuration

```bash
DB_DRIVER=mongo
DB_URL=localhost:27017
DB_TIMEOUT=10s
DB_NAME=battleship
DB_USER=root
DB_PASSWORD=battleship
```

### Server Configuration

```bash
SERVER_PORT=3000
SERVER_TIMEOUT=30s
SERVER_LOG_LEVEL=info
```

### Log Configuration

```bash
LOG_LEVEL=info
LOG_FORMAT=json
LOG_OUTPUT=stdout
```

## Default Values

If environment variables are not set, the following defaults are used:

### Database Defaults

```go
DatabaseConfig{
    Driver:   "mongo",
    URL:      "localhost:27017",
    Timeout:  10 * time.Second,
    Name:     "battleship",
    User:     "root",
    Password: "battleship",
}
```

### Server Defaults

```go
ServerConfig{
    Port:      3000,
    Timeout:   30 * time.Second,
    LogLevel:  "info",
}
```

### Log Defaults

```go
LogConfig{
    Level:  "info",
    Format: "json",
    Output: "stdout",
}
```

## Configuration Loading

Configuration is loaded in the following order:

1. Default values
2. Environment variables
3. Configuration file (if implemented)
4. Command-line flags (if implemented)

## Validation

The configuration is validated at startup:

1. Required fields are checked
2. Values are validated for correctness
3. Timeouts are checked for positive values
4. Ports are checked for valid ranges
5. Log levels are validated

## Usage Example

```go
// Load configuration
cfg, err := app.LoadConfig()
if err != nil {
    log.Fatalf("Failed to load configuration: %v", err)
}

// Validate configuration
if err := cfg.Validate(); err != nil {
    log.Fatalf("Invalid configuration: %v", err)
}

// Use configuration
db, err := storage.New(cfg.Database)
if err != nil {
    log.Fatalf("Failed to initialize database: %v", err)
}
```

## Development

### Adding New Configuration

1. Add new fields to the appropriate config struct
2. Add environment variable support
3. Set default values
4. Add validation rules
5. Update documentation

### Testing Configuration

```go
func TestConfig(t *testing.T) {
    // Test default values
    cfg := app.DefaultConfig()
    assert.Equal(t, "mongo", cfg.Database.Driver)
    assert.Equal(t, 3000, cfg.Server.Port)

    // Test environment variables
    os.Setenv("DB_DRIVER", "test-db")
    cfg, err := app.LoadConfig()
    assert.NoError(t, err)
    assert.Equal(t, "test-db", cfg.Database.Driver)

    // Test validation
    cfg.Database.Driver = ""
    err = cfg.Validate()
    assert.Error(t, err)
}
```

## Production Considerations

1. **Security**
   - Use secrets management
   - Encrypt sensitive values
   - Use secure defaults
   - Validate all inputs

2. **Monitoring**
   - Log configuration changes
   - Monitor validation errors
   - Track configuration usage
   - Alert on invalid configs

3. **Deployment**
   - Use environment-specific configs
   - Version configuration
   - Backup configurations
   - Document changes 