package oauth2

import "net/http"

type OAuthUser interface {
	GetEmail() string
	GetName() string
	GetPicture() string
}

type OAuthService interface {
	GetAuthURL() string
	HandleOAuthCallback(r *http.Request) (OAuthUser, error)
}
