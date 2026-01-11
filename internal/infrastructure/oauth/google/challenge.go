package google

import (
	"os"

	"golang.org/x/oauth2"
	googleOAuth2 "golang.org/x/oauth2/google"
)

type GoogleOAuthChallengeService struct {
	config *oauth2.Config
}

func NewGoogleOAuthChallengeService(clientID, clientSecret, redirectURI string) *GoogleOAuthChallengeService {
	if clientID == "" {
		clientID = os.Getenv("GOOGLE_CLIENT_ID")
	}
	if clientSecret == "" {
		clientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	}
	if redirectURI == "" {
		redirectURI = os.Getenv("GOOGLE_REDIRECT_URI")
		if redirectURI == "" {
			redirectURI = "http://localhost:8080/api/v1/auth/google/callback"
		}
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes: []string{
			"openid",
			"profile",
			"email",
		},
		Endpoint: googleOAuth2.Endpoint,
	}

	return &GoogleOAuthChallengeService{config: config}
}

func (s *GoogleOAuthChallengeService) GetAuthURL() string {
	return s.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

func (s *GoogleOAuthChallengeService) GetConfig() *oauth2.Config {
	return s.config
}
