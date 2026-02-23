# Auth Service

A centralized authentication service built with Go that provides secure user authentication and authorization for microservices. This service will support both HTTP and gRPC protocols, making it easy to integrate with various projects.

## Features

### Implemented Features

- ✅ User registration with email and password
- ✅ User login with JWT token generation
- ✅ Secure password hashing using bcrypt
- ✅ Token-based authentication using JWT
- ✅ PostgreSQL database integration
- ✅ Clean Architecture implementation
- ✅ HTTP REST API (Echo framework)
- ✅ gRPC service implementation
- ✅ Token refresh mechanism
- ✅ Google OAuth login
- ✅ Role-based access control (RBAC)
- ✅ User logout (single device and all devices)
- ✅ Swagger/OpenAPI documentation
- ✅ JWT middleware for protected routes

## Architecture

This project follows **Clean Architecture** principles, organizing code into three main layers:

```
auth-service/
├── cmd/api/              # Application entry point
├── internal/
│   ├── config/           # Configuration management
│   ├── domain/           # Business entities and interfaces
│   │   ├── user/
│   │   └── role/
│   ├── application/      # Use cases and business logic
│   │   ├── user/
│   │   │   ├── usecases/
│   │   │   ├── dto/
│   │   │   └── services/
│   │   └── role/
│   │       ├── usecases/
│   │       └── dto/
│   └── infrastructure/   # External dependencies
│       ├── postgres/
│       ├── oauth/google/
│       └── redis/
├── pkg/auth/             # Reusable authentication middleware
├── api/proto/            # Protocol buffer definitions
└── test/                 # Integration tests
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
- **HTTP Framework**: Echo (github.com/labstack/echo/v5)
- **Protocol Buffers**: gRPC (google.golang.org/grpc)
- **OAuth**: Google OAuth 2.0

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

| Variable           | Description                      | Required | Default |
| ------------------ | -------------------------------- | -------- | ------- |
| `DATABASE_URL`     | PostgreSQL connection string     | Yes      | -       |
| `JWT_SECRET`       | Secret key for JWT token signing | Yes      | -       |
| `PORT`             | HTTP server port                 | No       | `8080`  |
| `GRPC_PORT`        | gRPC server port                 | No       | `50051` |
| `GOOGLE_CLIENT_ID` | Google OAuth client ID           | No       | -       |

## Usage

### HTTP REST API

The HTTP API provides REST endpoints for authentication operations:

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

# Login with Google
POST /api/v1/auth/login/google
Content-Type: application/json

{
  "idToken": "google-id-token"
}

# Refresh token
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refreshToken": "your-refresh-token"
}

# Logout (single device)
POST /api/v1/auth/logout
Content-Type: application/json

{
  "refreshToken": "your-refresh-token"
}

# Logout (all devices)
POST /api/v1/auth/logout-all
Authorization: Bearer <access-token>

# Get Google OAuth URL
GET /api/v1/auth/google/challenge
```

### Role Management (RBAC)

```bash
# Create a role
POST /api/v1/roles
Authorization: Bearer <access-token>
Content-Type: application/json

{
  "name": "admin",
  "permissions": ["read", "write", "delete"]
}

# Get role
GET /api/v1/roles/:id
Authorization: Bearer <access-token>

# List roles
GET /api/v1/roles
Authorization: Bearer <access-token>

# Update role
PUT /api/v1/roles/:id
Authorization: Bearer <access-token>
Content-Type: application/json

{
  "name": "moderator",
  "permissions": ["read", "write"]
}

# Delete role
DELETE /api/v1/roles/:id
Authorization: Bearer <access-token>

# Assign role to user
POST /api/v1/roles/assign
Authorization: Bearer <access-token>
Content-Type: application/json

{
  "userId": "user-uuid",
  "roleId": "role-uuid"
}
```

### gRPC API

The gRPC service provides the same functionality through gRPC:

```protobuf
service AuthService {
  rpc Register(RegisterRequest) returns (RegisterResponse);
  rpc Login(LoginRequest) returns (LoginResponse);
  rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
  rpc RefreshToken(RefreshTokenRequest) returns (RefreshTokenResponse);
  rpc Logout(LogoutRequest) returns (LogoutResponse);
}

service RoleService {
  rpc CreateRole(CreateRoleRequest) returns (CreateRoleResponse);
  rpc GetRole(GetRoleRequest) returns (GetRoleResponse);
  rpc ListRoles(ListRolesRequest) returns (ListRolesResponse);
  rpc UpdateRole(UpdateRoleRequest) returns (UpdateRoleResponse);
  rpc DeleteRole(DeleteRoleRequest) returns (DeleteRoleResponse);
  rpc AssignRoleToUser(AssignRoleToUserRequest) returns (AssignRoleToUserResponse);
}
```

## Use Cases

### User Use Cases

### CreateUserUseCase

- Validates that the email is not already registered
- Hashes the password using bcrypt
- Creates a new user in the database
- Returns an access token

### LoginUserUseCase

- Validates user credentials (email and password)
- Generates a JWT token with user ID and expiration (24 hours)
- Returns the access token

### LoginWithGoogleUseCase

- Validates Google ID token
- Creates or retrieves user from database
- Returns an access token

### RefreshTokenUseCase

- Validates the refresh token
- Generates new access and refresh tokens
- Revokes the old refresh token

### LogoutUseCase

- Revokes refresh tokens (single or all)
- Supports logout from single device or all devices

### Role Use Cases

### CreateRoleUseCase

- Creates a new role with specified permissions

### AssignRoleToUserUseCase

- Assigns a role to a user
- Validates that both user and role exist

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
# Run all unit tests
go test -v ./... -tags=!integration

# Run integration tests (requires Docker)
go test -v ./... -tags=integration

# Run tests with coverage
go test -v -coverprofile=coverage.out ./... -tags=!integration
go tool cover -html=coverage.out -o coverage.html
```

### Swagger Documentation

When the service is running, access the Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

### Project Structure

- `internal/domain/user/`: User entity and repository interface
- `internal/domain/role/`: Role entity and repository interface
- `internal/application/user/usecases/`: Business logic for user operations
- `internal/application/role/usecases/`: Business logic for role operations
- `internal/infrastructure/postgres/repository/`: PostgreSQL implementation
- `internal/presentation/http/handlers/`: HTTP handlers
- `internal/presentation/http/middleware/`: HTTP middleware
- `pkg/auth/`: Reusable authentication middleware
- `cmd/api/main.go`: Application entry point

## Contributing

This is a personal project, but suggestions and improvements are welcome!

## License

[Add your license here]

## Roadmap

- [x] Implement HTTP REST API with Echo framework
- [x] Implement gRPC service with protocol buffers
- [x] Add token refresh mechanism
- [x] Add Google OAuth login
- [x] Add role management endpoints
- [x] Implement RBAC (Role-Based Access Control)
- [x] Add comprehensive test coverage
- [x] Add Docker support
- [x] Add Swagger/OpenAPI documentation
- [ ] Add password reset functionality
- [ ] Add user profile endpoints
- [ ] Add rate limiting and security middleware
- [ ] Add CI/CD pipeline
