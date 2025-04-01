# Development Guide

This guide provides information for developers working on the Battleship game project.

## Getting Started

1. **Prerequisites**
   - Go 1.21 or later
   - MongoDB 6.0 or later
   - Make (optional, for using Makefile commands)

2. **Setup**
   ```bash
   # Clone the repository
   git clone https://github.com/Jagreen1970/battleship.git
   cd battleship

   # Install dependencies
   go mod download

   # Set up environment variables
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Running the Application**
   ```bash
   # Development mode
   go run cmd/battleship/main.go

   # Production mode
   go build -o battleship cmd/battleship/main.go
   ./battleship
   ```

## Project Structure

The project follows a clean architecture pattern with the following structure:

```
cmd/
└── battleship/          # Application entry point
internal/
├── app/                 # Application core
├── game/               # Game logic
├── server/            # HTTP server
└── storage/          # Data persistence
```

## Development Guidelines

### Code Style

1. **Go Standards**
   - Follow [Effective Go](https://golang.org/doc/effective_go)
   - Use `gofmt` for code formatting
   - Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

2. **Naming Conventions**
   - Use clear, descriptive names
   - Follow Go naming conventions
   - Use interfaces for abstraction
   - Prefix interfaces with 'I' (e.g., `IStorage`)

3. **Error Handling**
   - Always check errors
   - Use `fmt.Errorf` with `%w` for wrapping
   - Create custom error types when needed
   - Provide context in error messages

### Testing

1. **Unit Tests**
   - Write tests for all packages
   - Use table-driven tests
   - Mock external dependencies
   - Test edge cases

2. **Integration Tests**
   - Test package interactions
   - Use test containers
   - Clean up test data
   - Handle test timeouts

3. **Running Tests**
   ```bash
   # Run all tests
   go test ./...

   # Run tests with coverage
   go test -cover ./...

   # Run specific package tests
   go test ./internal/game
   ```

### Configuration

1. **Environment Variables**
   - Use `.env` for local development
   - Document all environment variables
   - Provide default values
   - Validate configuration

2. **Configuration Structure**
   ```go
   type Config struct {
       Database DatabaseConfig
       Server   ServerConfig
       Log      LogConfig
   }
   ```

### Database

1. **MongoDB**
   - Use indexes for performance
   - Implement proper error handling
   - Use transactions when needed
   - Follow MongoDB best practices

2. **Migrations**
   - Version database schema
   - Provide migration scripts
   - Test migrations
   - Document changes

### API Design

1. **RESTful Principles**
   - Use proper HTTP methods
   - Follow resource naming
   - Implement proper status codes
   - Version the API

2. **Request/Response**
   - Validate input
   - Use proper content types
   - Implement rate limiting
   - Handle errors consistently

## CLI Development

### Running CLI Mode
```bash
# Development mode with CLI
go run cmd/battleship/main.go --cli

# Production mode with CLI
go build -o battleship cmd/battleship/main.go
./battleship --cli
```

### CLI Commands
- `create-game`: Create a new game session
- `join-game <game-id>`: Join an existing game
- `place-ship <game-id> <x> <y> <direction>`: Place a ship
- `fire <game-id> <x> <y>`: Fire at coordinates
- `show-board <game-id>`: Display game board
- `exit`: Exit CLI mode

### Testing CLI Mode
```bash
# Run CLI tests
go test ./internal/cli/...

# Run with coverage
go test -cover ./internal/cli/...
```

## Deployment

1. **Docker**
   ```bash
   # Build image
   docker build -t battleship .

   # Run container
   docker run -p 3000:3000 battleship
   ```

2. **Kubernetes**
   - Use deployment manifests
   - Configure health checks
   - Set resource limits
   - Implement proper logging

## Monitoring

1. **Logging**
   - Use structured logging
   - Include correlation IDs
   - Log appropriate levels
   - Handle sensitive data

2. **Metrics**
   - Track key metrics
   - Use Prometheus format
   - Implement health checks
   - Monitor performance

## Contributing

1. **Workflow**
   - Create feature branch
   - Write tests
   - Update documentation
   - Submit pull request

2. **Code Review**
   - Review for style
   - Check test coverage
   - Verify documentation
   - Test changes

3. **Release Process**
   - Update version
   - Update changelog
   - Tag release
   - Deploy changes

## Troubleshooting

1. **Common Issues**
   - Database connection
   - Configuration errors
   - Test failures
   - Build problems

2. **Debugging**
   - Use logging
   - Enable debug mode
   - Use debugger
   - Check logs

## Support

- Create issues for bugs
- Use discussions for questions
- Follow security policy
- Check documentation 