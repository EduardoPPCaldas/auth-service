package main

import (
	"fmt"
	"log"

	roleusecases "github.com/EduardoPPCaldas/auth-service/internal/application/role/usecases"
	"github.com/EduardoPPCaldas/auth-service/internal/application/user/services/token"
	"github.com/EduardoPPCaldas/auth-service/internal/application/user/usecases"
	"github.com/EduardoPPCaldas/auth-service/internal/config"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/role"
	tokenDomain "github.com/EduardoPPCaldas/auth-service/internal/domain/token"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"github.com/EduardoPPCaldas/auth-service/internal/infrastructure/oauth/google"
	postgresRepo "github.com/EduardoPPCaldas/auth-service/internal/infrastructure/postgres/repository"
	"github.com/EduardoPPCaldas/auth-service/internal/presentation/http"
	"github.com/EduardoPPCaldas/auth-service/internal/presentation/http/handlers"
	"github.com/EduardoPPCaldas/auth-service/pkg/auth"
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

	cfg := config.Load()

	// Initialize database
	db, err := initDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories
	userRepo := postgresRepo.NewUserRepository(db)
	roleRepo := postgresRepo.NewRoleRepository(db)
	refreshTokenRepo := postgresRepo.NewRefreshTokenRepository(db)

	// Initialize services
	googleValidator := google.NewGoogleTokenValidator(cfg.GoogleClientID)
	tokenGenerator := token.NewTokenGenerator()
	googleOAuthService := google.NewGoogleOAuthChallengeService(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURI)

	// Initialize additional services
	refreshTokenService := token.NewRefreshTokenService(refreshTokenRepo, userRepo, cfg.JWTRefreshExpiry)

	// Initialize auth middleware
	authMiddleware, err := auth.NewAuthMiddleware(auth.WithJWTSecret(cfg.JWTSecret))
	if err != nil {
		log.Fatalf("Failed to initialize auth middleware: %v", err)
	}

	// Initialize use cases
	createUserUseCase := usecases.NewCreateUserUseCase(userRepo, roleRepo, tokenGenerator)
	loginUserUseCase := usecases.NewLoginUserUseCase(userRepo, tokenGenerator)
	loginWithGoogleUseCase := usecases.NewLoginWithGoogleUseCase(userRepo, roleRepo, tokenGenerator, googleValidator)
	refreshTokenUseCase := usecases.NewRefreshTokenUseCase(userRepo, tokenGenerator, refreshTokenService, cfg.JWTAccessExpiry)
	logoutUseCase := usecases.NewLogoutUseCase(userRepo, refreshTokenService)

	// Initialize role management use cases
	createRoleUseCase := roleusecases.NewCreateRoleUseCase(roleRepo, userRepo)
	updateRoleUseCase := roleusecases.NewUpdateRoleUseCase(roleRepo, userRepo)
	deleteRoleUseCase := roleusecases.NewDeleteRoleUseCase(roleRepo, userRepo)
	listRolesUseCase := roleusecases.NewListRolesUseCase(roleRepo)
	getRoleUseCase := roleusecases.NewGetRoleUseCase(roleRepo)
	assignRoleToUserUseCase := roleusecases.NewAssignRoleToUserUseCase(roleRepo, userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(
		createUserUseCase,
		loginUserUseCase,
		loginWithGoogleUseCase,
		refreshTokenUseCase,
		logoutUseCase,
		googleOAuthService,
	)

	roleHandler := handlers.NewRoleHandler(
		createRoleUseCase,
		updateRoleUseCase,
		deleteRoleUseCase,
		listRolesUseCase,
		getRoleUseCase,
		assignRoleToUserUseCase,
	)

	// Initialize Echo
	e := echo.New()

	// Middleware
	http.SetupMiddleware(e)
	e.Validator = &CustomValidator{validator: validator.New()}

	// Routes
	http.SetupRoutes(
		e,
		authHandler,
		roleHandler,
		func() echo.MiddlewareFunc { return authMiddleware.EchoMiddleware() },
		func() echo.MiddlewareFunc { return authMiddleware.EchoRequireRole(role.RoleAdmin) },
	)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initDatabase(dbURL string) (*gorm.DB, error) {
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Auto-migrate entities
	if err := db.AutoMigrate(&user.User{}, &tokenDomain.RefreshToken{}, &role.Role{}, &role.Permission{}); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate: %w", err)
	}

	return db, nil
}

// CustomValidator is a custom validator for Echo
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}
