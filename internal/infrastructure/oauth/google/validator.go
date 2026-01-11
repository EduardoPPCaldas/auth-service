package google

import (
	"context"
	"fmt"
	"os"

	"github.com/EduardoPPCaldas/auth-service/internal/application/user/services/oauth"
	"google.golang.org/api/idtoken"
)

type GoogleTokenValidator struct {
	clientID string
}

func NewGoogleTokenValidator(clientID string) *GoogleTokenValidator {
	return &GoogleTokenValidator{clientID: clientID}
}

func (v *GoogleTokenValidator) Validate(ctx context.Context, token string) (*oauth.GoogleUser, error) {
	audience := v.clientID
	if audience == "" {
		audience = os.Getenv("GOOGLE_CLIENT_ID")
	}


	payload, err := idtoken.Validate(ctx, token, audience)
	if err != nil {
		return nil, fmt.Errorf("failed to validate google token: %w", err)
	}

	email, ok := payload.Claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("token does not contain email claim")
	}

	name, ok := payload.Claims["name"].(string)
	if !ok {
		// Name is optional in some contexts, but usually present.
		name = ""
	}

	return &oauth.GoogleUser{
		Email: email,
		Name:  name,
	}, nil
}
