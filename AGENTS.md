# AGENTS.md

This file contains guidelines and commands for agentic coding agents working on this Go authentication service project.

## Development Commands

### Build & Test Commands
```bash
# Build the entire project
go build ./...

# Run all tests (unit tests only)
make test
go test -v ./...

# Run only unit tests (excludes integration tests)
make test-unit
go test -v ./... -tags=!integration

# Run integration tests (requires Docker)
make test-integration
go test -v ./... -tags=integration

# Run all tests including integration
make test-all
go test -v ./... -tags=integration

# Run tests with coverage
make test-coverage
go test -v -coverprofile=coverage.out ./... -tags=!integration
go tool cover -html=coverage.out -o coverage.html

# Run specific test function
go test -v ./path/to/package -run TestFunctionName
go test -v ./internal/application/user/usecases -run TestLoginUserUseCase_Execute_Success

# Run specific test file
go test -v ./path/to/package/file_test.go
go test -v ./internal/application/user/usecases/login_user_test.go
```

### Development & Docker Commands
```bash
# Run the application locally
make run
go run cmd/api/main.go

# Setup development database
make setup-db

# Docker commands
make docker-build
make docker-run
make docker-dev
make docker-stop
make docker-clean

# Integration test with Docker
make test-integration-docker
make clean-test
```

### Linting & Formatting
```bash
# Format code (no custom linting configuration - use Go defaults)
go fmt ./...
goimports -w .  # If goimports is available

# Vet for potential issues
go vet ./...

# Tidy dependencies
go mod tidy
```

## Code Style Guidelines

### 1. Project Structure
```
/
├── cmd/api/           # Application entry point
├── internal/          # Private application code
│   ├── application/   # Use cases and business logic
│   ├── domain/        # Domain entities and interfaces
│   ├── infrastructure/ # External implementations
│   └── presentation/ # HTTP handlers and routing
├── pkg/auth/          # Reusable authentication middleware
├── api/proto/         # Protocol buffer definitions
└── test/             # Integration tests
```

### 2. Import Organization
- Use standard Go import grouping:
  1. Standard library
  2. Third-party packages
  3. Internal packages (`github.com/EduardoPPCaldas/auth-service/...`)
- No blank lines between import groups
- Use named imports only when necessary for disambiguation

### 3. Naming Conventions
- **Packages**: `lowercase`, `short`, `descriptive` (e.g., `user`, `auth`, `handlers`)
- **Files**: `snake_case.go` (e.g., `login_user.go`, `auth_handler.go`)
- **Constants**: `PascalCase` for exported constants (Go convention), `camelCase` for unexported
- **Variables**: `camelCase` for local, `PascalCase` for exported
- **Functions**: `PascalCase` for exported, `camelCase` for unexported
- **Interfaces**: Often end with `-er` suffix (e.g., `UserRepository`, `TokenGenerator`)
- **Structs**: `PascalCase` with descriptive field names

### 4. Types & Modern Go
- Use `any` instead of `interface{}` (Go 1.18+)
- Use `context.Context` as first parameter for functions that need it
- Prefer typed nil checks for pointers (`*string` vs `string`)
- Use `uuid.UUID` for ID fields
- Use `time.Time` for timestamps

### 5. Error Handling
- Always wrap errors with context: `fmt.Errorf("operation failed: %w", err)`
- Use structured errors with specific types in domain layers
- Return errors as the last return value
- Handle errors immediately or wrap and return up the call stack
- Use `errors.Is()` and `errors.As()` for error checking

### 6. Testing Patterns
- Use table-driven tests for multiple scenarios
- Follow Arrange-Act-Assert pattern:
  ```go
  func TestLoginUserUseCase_Execute_Success(t *testing.T) {
      // Arrange
      mockRepo := new(mocks.MockUserRepository)
      // ... setup mocks and test data
      
      // Act
      result, err := useCase.Execute(email, password)
      
      // Assert
      assert.NoError(t, err)
      assert.Equal(t, expected, result)
      mockRepo.AssertExpectations(t)
  }
  ```

- Use testify/mock for interface mocking
- Mock external dependencies, not internal logic
- Test both success and failure scenarios
- Use descriptive test function names: `TestStruct_Method_Scenario`

### 7. Dependency Injection
- Use constructor functions with dependency injection
- Define interfaces for external dependencies
- Pass interfaces to constructors, not concrete types
- Keep constructors simple and focused

### 8. HTTP Handlers (Echo)
- Use Echo framework with JSON responses
- Define DTOs in `application/user/dto/`
- Validate requests using struct tags
- Return consistent error responses
- Use HTTP status codes appropriately

### 9. Database (GORM)
- Use GORM for ORM functionality
- Define entities in `domain/` package
- Use pointer types for optional fields (`*string`)
- Follow GORM conventions for field tags
- Handle database errors appropriately

### 10. JWT & Security
- Use the `pkg/auth` middleware for authentication
- Store JWT secrets in environment variables
- Use proper token validation and expiration
- Never log sensitive information (tokens, passwords)
- Use bcrypt for password hashing

### 11. Context
- Always use context.Context as the first parameter for functions that need it
- Pass context through the call stack
- Use context for cancellation, deadlines, and passing request-scoped values
- Never use background context for application logic

## Environment Variables
Required environment variables:
- `JWT_SECRET`: JWT signing secret
- `DATABASE_URL`: PostgreSQL connection string
- `PORT`: Server port (default: 8080)
- `GOOGLE_CLIENT_ID`: Google OAuth client ID (optional)

## Testing Requirements
- All new code must have unit tests
- Aim for >80% code coverage
- Run `make test-unit` before committing
- Integration tests should run with `make test-integration`

## Git Workflow
- Use conventional commit messages
- Run tests before committing
- Ensure `go build ./...` succeeds
- Keep changes focused and atomic