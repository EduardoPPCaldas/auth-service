# Auth Service

A centralized authentication service built with Go that provides secure user authentication and authorization for microservices. This service will support both HTTP and gRPC protocols, making it easy to integrate with various projects.

## Features

### Current Features

- ✅ User registration with email and password
- ✅ User login with JWT token generation
- ✅ Secure password hashing using bcrypt
- ✅ Token-based authentication using JWT
- ✅ PostgreSQL database integration
- ✅ Clean Architecture implementation

### Planned Features

- 🚧 HTTP REST API endpoints
- 🚧 gRPC service implementation
- 🚧 Token refresh mechanism
- 🚧 Password reset functionality
- 🚧 User profile management
- 🚧 Role-based access control (RBAC)

## Architecture

This project follows **Clean Architecture** principles, organizing code into three main layers:

```
auth-service/
├── cmd/api/              # Application entry point
├── internal/
│   ├── domain/          # Business entities and interfaces
│   │   └── user/
│   ├── application/     # Use cases and business logic
│   │   └── user/usecases/
│   └── infrastructure/  # External dependencies (database, APIs)
│       └── postgres/
└── README.md
```

### Layers

- **Domain Layer**: Contains core business entities (`User`) and repository interfaces
- **Application Layer**: Implements use cases like `CreateUserUseCase` and `LoginUserUseCase`
- **Infrastructure Layer**: Provides concrete implementations (PostgreSQL repository)

## Technology Stack

- **Language**: Go 1.25+
- **Database**: PostgreSQL (via GORM)
- **Authentication**: JWT (github.com/golang-jwt/jwt/v5)
- **Password Hashing**: bcrypt (golang.org/x/crypto/bcrypt)
- **ORM**: GORM (gorm.io/gorm)
- **UUID**: google/uuid

## Prerequisites

- Go 1.25 or higher
- PostgreSQL database
- (Optional) Docker for containerized PostgreSQL

## Installation

1. Clone the repository:

```bash
git clone https://github.com/EduardoPPCaldas/auth-service.git
cd auth-service
```

2. Install dependencies:

```bash
go mod download
```

3. Set up environment variables:

```bash
export DATABASE_URL="postgres://user:password@localhost:5432/authdb?sslmode=disable"
export JWT_SECRET="your-secret-key-here"
```

4. Set up the database:

   - Create a PostgreSQL database named `authdb` (or your preferred name)
   - Update the `DATABASE_URL` accordingly

5. Run migrations (when available):

```bash
# Migration commands will be added here
```

## Configuration

The service requires the following environment variables:

| Variable       | Description                      | Required | Default |
| -------------- | -------------------------------- | -------- | ------- |
| `DATABASE_URL` | PostgreSQL connection string     | Yes      | -       |
| `JWT_SECRET`   | Secret key for JWT token signing | Yes      | -       |
| `PORT`         | HTTP server port                 | No       | `8080`  |
| `GRPC_PORT`    | gRPC server port                 | No       | `50051` |

## Usage

### HTTP API (Coming Soon)

The HTTP API will provide REST endpoints for authentication operations:

```bash
# Register a new user
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}

# Login
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}

# Response
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### gRPC API (Coming Soon)

The gRPC service will provide the same functionality through gRPC:

```protobuf
service AuthService {
  rpc Register(RegisterRequest) returns (RegisterResponse);
  rpc Login(LoginRequest) returns (LoginResponse);
  rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
}
```

## Use Cases

### CreateUserUseCase

- Validates that the email is not already registered
- Hashes the password using bcrypt
- Creates a new user in the database

### LoginUserUseCase

- Validates user credentials (email and password)
- Generates a JWT token with user ID and expiration (24 hours)
- Returns the token for use in subsequent requests

## Security

- Passwords are hashed using bcrypt with default cost factor
- JWT tokens are signed using HS256 algorithm
- Tokens expire after 24 hours (configurable)
- Email uniqueness is enforced at the application level

## Development

### Running the Service

```bash
go run cmd/api/main.go
```

### Running Tests

```bash
go test ./...
```

### Project Structure

- `internal/domain/user/`: User entity and repository interface
- `internal/application/user/usecases/`: Business logic for user operations
- `internal/infrastructure/postgres/repository/`: PostgreSQL implementation
- `cmd/api/main.go`: Application entry point

## Contributing

This is a personal project, but suggestions and improvements are welcome!

## License

[Add your license here]

## Roadmap

- [ ] Implement HTTP REST API with Gin or Echo framework
- [ ] Implement gRPC service with protocol buffers
- [ ] Add token refresh mechanism
- [ ] Add password reset functionality
- [ ] Add user profile endpoints
- [ ] Implement RBAC (Role-Based Access Control)
- [ ] Add comprehensive test coverage
- [ ] Add Docker support
- [ ] Add CI/CD pipeline
- [ ] Add API documentation (OpenAPI/Swagger)
- [ ] Add rate limiting and security middleware
