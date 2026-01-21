# API Documentation

This directory contains OpenAPI/Swagger specifications for the Auth Service API.

## Files

- `swagger.json` - OpenAPI 3.0 specification in JSON format
- `swagger.yaml` - OpenAPI 3.0 specification in YAML format

## Accessing Swagger Documentation

### Swagger UI (Interactive Documentation)

The interactive Swagger UI is available at:
```
GET /swagger
GET /swagger/
```

Simply start the server and navigate to `http://localhost:8080/swagger` in your browser to view and interact with the API documentation.

### Swagger JSON Specification

The raw Swagger JSON specification is served at:
```
GET /swagger.json
```

This endpoint can be used by API clients, code generators, or other tools that consume OpenAPI specifications.

### Alternative Viewing Methods

You can also view the API documentation using:

1. **Swagger Editor (Online):**
   - Go to https://editor.swagger.io/
   - Copy the contents of `swagger.yaml` or `swagger.json`
   - Paste it into the editor

2. **Docker:**
   ```bash
   docker run -p 8081:8080 -e SWAGGER_JSON=/api/swagger.json -v $(pwd)/api:/api swaggerapi/swagger-ui
   ```
   Then open http://localhost:8081 in your browser

## API Endpoints

### Authentication

- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login with email and password
- `POST /api/v1/auth/login/google` - Login with Google OAuth
- `GET /api/v1/auth/google/challenge` - Get Google OAuth challenge URL

## Bruno Collection

A Bruno API collection is available in the `bruno/` directory. To use it:

1. Install Bruno: https://www.usebruno.com/
2. Open Bruno and click "Open Collection"
3. Navigate to the `bruno` directory in this project
4. Select the collection to load all requests

The collection includes:
- Pre-configured requests for all endpoints
- Environment variables for easy server switching
- Example request bodies
- Documentation for each endpoint

