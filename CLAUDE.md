# Battleship Codebase Guidelines

## Build and Test Commands
- Build: `go build -o battleship cmd/battleship/main.go`
- Run: `go run cmd/battleship/main.go`
- Run CLI mode: `go run cmd/battleship/main.go --cli`
- Test all packages: `go test ./... -v`
- Test single package: `go test ./internal/game -v`
- Test single test: `go test ./internal/game -run TestNewBoard -v`
- Test with coverage: `make test-cover` or `go test ./... -v -coverprofile=coverage.out`
- Frontend lint: `cd frontend && npm run lint`
- Frontend format: `cd frontend && npm run format`

## Code Style Guidelines
- **Go Standards**: Follow Effective Go and Go Code Review Comments
- **Formatting**: Use `gofmt` for formatting Go code
- **Imports**: Standard library first, then external packages, then internal packages
- **Naming**: Prefix interfaces with 'I' (e.g., `IStorage`), descriptive variable names
- **Types**: Use interfaces for abstraction, define custom types when appropriate
- **Error Handling**: Always check errors, use `fmt.Errorf` with `%w` for wrapping, provide context
- **Testing**: Table-driven tests, comprehensive test cases, use `testify/assert` package
- **Frontend**: ESLint + Prettier for JS/TS, 2-space indentation
- **Documentation**: Document all exported functions, types, and variables
- **Architecture**: Clean Architecture pattern with internal packages