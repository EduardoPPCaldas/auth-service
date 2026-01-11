package oauth

import "context"

type GoogleUser struct {
	Email string
	Name  string
}

type GoogleTokenValidator interface {
	Validate(ctx context.Context, idToken string) (*GoogleUser, error)
}
