package googleAuthService

import (
	"encoding/json"
	"errors"
	"fmt"
	"forum-app/helpers/envutil"
	oauth2 "forum-app/services"
	"io"
	"net/http"
	"net/url"
)

type GoogleAuthConfig struct {
	GoogleRedirect string
	GoogleClientId string
	GoogleCallback string
	ResponseType   string
	ScopeGoogle    string
	AccessType     string
	GoogleSecret   string
	State          string
}

type GoogleUser struct {
	FamilyName string `json:"family_name"`
	Name       string `json:"name"`
	Picture    string `json:"picture"`
	Email      string `json:"email"`
	GivenName  string `json:"give_name"`
}

type ResponseGoogle struct {
	AccessToken string `json:"access_token"`
}

func (u GoogleUser) GetEmail() string   { return u.Email }
func (u GoogleUser) GetName() string    { return u.Name }
func (u GoogleUser) GetPicture() string { return u.Picture }

func NewGoogleAuthConfig(state string) *GoogleAuthConfig {
	return &GoogleAuthConfig{
		GoogleRedirect: "https://accounts.google.com/o/oauth2/v2/auth",
		GoogleClientId: envutil.GetEnvString("GOOGLE_PUBLIC_KEY"),
		GoogleCallback: envutil.GetEnvString("GOOGLE_AUTH_CALLBACK"),
		ResponseType:   "code",
		ScopeGoogle:    "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile",
		AccessType:     "offline",
		GoogleSecret:   envutil.GetEnvString("GOOGLE_SECRET_KEY"),
		State:          state,
	}
}

func (g *GoogleAuthConfig) GetAuthURL() string {
	params := url.Values{}
	params.Add("client_id", g.GoogleClientId)
	params.Add("redirect_uri", g.GoogleCallback)
	params.Add("response_type", g.ResponseType)
	params.Add("scope", g.ScopeGoogle)
	params.Add("access_type", g.AccessType)
	params.Add("state", g.State)

	return fmt.Sprintf("%s?%s", g.GoogleRedirect, params.Encode())
}

func (g *GoogleAuthConfig) HandleOAuthCallback(r *http.Request) (oauth2.OAuthUser, error) {

	if g.State != r.URL.Query().Get("state") {
		return GoogleUser{}, errors.New("missingstate")
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		return GoogleUser{}, errors.New("missing code in callback")
	}

	params := url.Values{}
	params.Add("code", code)
	params.Add("client_id", g.GoogleClientId)
	params.Add("client_secret", g.GoogleSecret)
	params.Add("redirect_uri", g.GoogleCallback)
	params.Add("grant_type", "authorization_code")

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", params)
	if err != nil {
		return GoogleUser{}, fmt.Errorf("error posting form: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GoogleUser{}, fmt.Errorf("error reading token response: %w", err)
	}

	var accessTokenResp ResponseGoogle
	if err := json.Unmarshal(body, &accessTokenResp); err != nil {
		return GoogleUser{}, fmt.Errorf("error unmarshalling token response: %w", err)
	}

	if accessTokenResp.AccessToken == "" {
		return GoogleUser{}, errors.New("access token not found in response")
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+accessTokenResp.AccessToken)

	userInfoResp, err := client.Do(req)
	if err != nil {
		return GoogleUser{}, fmt.Errorf("error fetching user info: %w", err)
	}
	defer userInfoResp.Body.Close()

	userInfo, err := io.ReadAll(userInfoResp.Body)

	if err != nil {
		return GoogleUser{}, fmt.Errorf("error reading user info: %w", err)
	}

	var user GoogleUser

	if err := json.Unmarshal(userInfo, &user); err != nil {
		return GoogleUser{}, fmt.Errorf("error unmarshaling user info: %w", err)
	}

	return user, nil
}
