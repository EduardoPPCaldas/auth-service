.PHONY: test test-unit test-integration test-all run setup-db docker-build docker-run docker-dev docker-stop docker-clean

# Run all tests (unit tests only, integration tests require docker)
test:
	go test -v ./...

# Run only unit tests (excludes integration tests)
test-unit:
	go test -v ./... -tags=!integration

# Run integration tests (requires Docker)
test-integration:
	go test -v ./... -tags=integration

# Run all tests including integration
test-all:
	go test -v ./... -tags=integration

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./... -tags=!integration
	go tool cover -html=coverage.out -o coverage.html

# Run integration tests with docker-compose
test-integration-docker:
	docker-compose -f docker-compose.test.yml up -d
	sleep 5
	DATABASE_URL=postgres://testuser:testpass@localhost:5433/testdb?sslmode=disable go test -v ./... -tags=integration
	docker-compose -f docker-compose.test.yml down

# Clean up test containers
clean-test:
	docker-compose -f docker-compose.test.yml down -v

# Run the application
run:
	go run cmd/api/main.go

# Setup database (start PostgreSQL container)
setup-db:
	@echo "Starting PostgreSQL database container..."
	docker-compose up -d
	@echo "Waiting for database to be ready..."
	@sleep 3
	@echo "✅ Database is ready!"
	@echo "Connection string: postgres://postgres:postgres@localhost:5432/authdb?sslmode=disable"

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t auth-service:latest .

# Run application in Docker
docker-run:
	@echo "Starting application in Docker..."
	docker-compose up -d postgres
	@sleep 3
	docker run --rm -it \
		--name auth-service-app \
		--network auth-service_default \
		-p 8080:8080 \
		-e DATABASE_URL=postgres://postgres:postgres@postgres:5432/authdb?sslmode=disable \
		-e PORT=8080 \
		-e JWT_SECRET=$${JWT_SECRET:-your-secret-key} \
		-e GOOGLE_CLIENT_ID=$${GOOGLE_CLIENT_ID:-} \
		auth-service:latest

# Run application in Docker with debugging enabled
docker-dev:
	@echo "Starting application in Docker with debugging..."
	docker-compose -f docker-compose.dev.yml up --build

# Stop Docker containers
docker-stop:
	@echo "Stopping Docker containers..."
	docker-compose down
	docker-compose -f docker-compose.dev.yml down

# Clean Docker containers and volumes
docker-clean:
	@echo "Cleaning Docker containers and volumes..."
	docker-compose down -v
	docker-compose -f docker-compose.dev.yml down -v
	docker rmi auth-service:latest 2>/dev/null || true
