package main

import (
	"fmt"
	"log"
	"os"

	"github.com/EduardoPPCaldas/auth-service/internal/application/user/services/token"
	"github.com/EduardoPPCaldas/auth-service/internal/application/user/usecases"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"github.com/EduardoPPCaldas/auth-service/internal/infrastructure/oauth/google"
	postgresRepo "github.com/EduardoPPCaldas/auth-service/internal/infrastructure/postgres/repository"
	"github.com/EduardoPPCaldas/auth-service/internal/presentation/http"
	"github.com/EduardoPPCaldas/auth-service/internal/presentation/http/handlers"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		// .env file is optional, so we just log a warning if it's not found
		log.Println("Warning: .env file not found, using system environment variables")
	}
	// Initialize database
	db, err := initDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories
	userRepo := postgresRepo.NewUserRepository(db)

	// Initialize services
	tokenGenerator := token.NewTokenGenerator()
	googleValidator := google.NewGoogleTokenValidator(os.Getenv("GOOGLE_CLIENT_ID"))
	googleOAuthService := google.NewGoogleOAuthChallengeService("", "", "")

	// Initialize use cases
	createUserUseCase := usecases.NewCreateUserUseCase(userRepo, tokenGenerator)
	loginUserUseCase := usecases.NewLoginUserUseCase(userRepo, tokenGenerator)
	loginWithGoogleUseCase := usecases.NewLoginWithGoogleUseCase(userRepo, tokenGenerator, googleValidator)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(
		createUserUseCase,
		loginUserUseCase,
		loginWithGoogleUseCase,
		googleOAuthService,
	)

	// Initialize Echo
	e := echo.New()

	// Middleware
	http.SetupMiddleware(e)
	e.Validator = &CustomValidator{validator: validator.New()}

	// Routes
	http.SetupRoutes(e, authHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initDatabase() (*gorm.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Auto-migrate User entity
	if err := db.AutoMigrate(&user.User{}); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate: %w", err)
	}

	return db, nil
}

// CustomValidator is a custom validator for Echo
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
