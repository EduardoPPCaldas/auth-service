.PHONY: test test-unit test-integration test-all

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

