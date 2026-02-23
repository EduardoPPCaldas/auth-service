package google

import (
	"golang.org/x/oauth2"
	googleOAuth2 "golang.org/x/oauth2/google"
)

type GoogleOAuthChallengeService struct {
	config *oauth2.Config
}

func NewGoogleOAuthChallengeService(clientID, clientSecret, redirectURI string) *GoogleOAuthChallengeService {
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
