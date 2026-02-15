package auth

import (
	"context"
	"errors"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthMiddleware struct {
	jwtSecret       []byte
	tokenValidation TokenValidation
}

type TokenValidation struct {
	SkipExpirationCheck bool
	RequiredIssuer      string
	RequiredAudience    []string
}

type UserContext struct {
	UserID      uuid.UUID
	Claims      jwt.MapClaims
	Permissions []string
	Roles       []string
}

type CustomClaims struct {
	UserID      uuid.UUID `json:"sub"`
	Permissions []string  `json:"permissions,omitempty"`
	Roles       []string  `json:"roles,omitempty"`
	Issuer      string    `json:"iss,omitempty"`
	Audience    []string  `json:"aud,omitempty"`
	jwt.RegisteredClaims
}

type MiddlewareOption func(*AuthMiddleware)

func WithJWTSecret(secret string) MiddlewareOption {
	return func(am *AuthMiddleware) {
		am.jwtSecret = []byte(secret)
	}
}

func WithTokenValidation(validation TokenValidation) MiddlewareOption {
	return func(am *AuthMiddleware) {
		am.tokenValidation = validation
	}
}

func NewAuthMiddleware(opts ...MiddlewareOption) (*AuthMiddleware, error) {
	am := &AuthMiddleware{
		jwtSecret: []byte{},
		tokenValidation: TokenValidation{
			SkipExpirationCheck: false,
		},
	}

	for _, opt := range opts {
		opt(am)
	}

	if len(am.jwtSecret) == 0 {
		return nil, errors.New("JWT secret is required")
	}

	return am, nil
}

func (am *AuthMiddleware) CreateToken(userID uuid.UUID, expiresAt time.Time, customClaims map[string]any) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   userID.String(),
		},
	}

	if am.tokenValidation.RequiredIssuer != "" {
		claims.Issuer = am.tokenValidation.RequiredIssuer
	}

	if len(am.tokenValidation.RequiredAudience) > 0 {
		claims.Audience = am.tokenValidation.RequiredAudience
	}

	for key, value := range customClaims {
		switch key {
		case "permissions":
			if perms, ok := value.([]string); ok {
				claims.Permissions = perms
			}
		case "roles":
			if roles, ok := value.([]string); ok {
				claims.Roles = roles
			}
		default:
			// Add other custom claims as needed
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(am.jwtSecret)
}

func (am *AuthMiddleware) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Bearer token required", http.StatusUnauthorized)
			return
		}

		token, err := am.parseAndValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*CustomClaims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		if err := am.validateClaims(claims); err != nil {
			http.Error(w, "Token validation failed: "+err.Error(), http.StatusUnauthorized)
			return
		}

		userCtx := UserContext{
			UserID: claims.UserID,
			Claims: jwt.MapClaims{
				"sub": claims.Subject,
				"iss": claims.Issuer,
				"aud": claims.Audience,
				"exp": claims.ExpiresAt,
				"nbf": claims.NotBefore,
				"iat": claims.IssuedAt,
				"jti": claims.ID,
			},
			Permissions: claims.Permissions,
			Roles:       claims.Roles,
		}

		ctx := context.WithValue(r.Context(), "user", userCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (am *AuthMiddleware) parseAndValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return am.jwtSecret, nil
	})
}

func (am *AuthMiddleware) validateClaims(claims *CustomClaims) error {
	if !am.tokenValidation.SkipExpirationCheck {
		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			return errors.New("token has expired")
		}
	}

	if am.tokenValidation.RequiredIssuer != "" && claims.Issuer != am.tokenValidation.RequiredIssuer {
		return errors.New("invalid issuer")
	}

	if len(am.tokenValidation.RequiredAudience) > 0 {
		for _, requiredAud := range am.tokenValidation.RequiredAudience {
			found := slices.Contains(claims.Audience, requiredAud)
			if !found {
				return errors.New("invalid audience")
			}
		}
	}

	return nil
}

func (am *AuthMiddleware) RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userCtx, ok := r.Context().Value("user").(UserContext)
			if !ok {
				http.Error(w, "User context not found", http.StatusUnauthorized)
				return
			}

			if !slices.Contains(userCtx.Permissions, permission) {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (am *AuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userCtx, ok := r.Context().Value("user").(UserContext)
			if !ok {
				http.Error(w, "User context not found", http.StatusUnauthorized)
				return
			}

			if !slices.Contains(userCtx.Roles, role) {
				http.Error(w, "Insufficient role", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetUserFromContext(ctx context.Context) (UserContext, bool) {
	user, ok := ctx.Value("user").(UserContext)
	return user, ok
}
