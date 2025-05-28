package githubAuthService

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

type GithubAuthConfig struct {
	GithubRedirect string
	GithubClientId string
	GithubCallback string
	ResponseType   string
	ScopeGithub    string
	GithubSecret   string
	State          string
}

type GithubUser struct {
	FamilyName string `json:"name"`
	Name       string `json:"login"`
	Picture    string `json:"avatar_url"`
	Email      string `json:"email"`
}

type ResponseGithub struct {
	AccessToken string
}

func (u GithubUser) GetEmail() string   { return u.Email }
func (u GithubUser) GetName() string    { return u.Name }
func (u GithubUser) GetPicture() string { return u.Picture }

func NewGithubAuthConfig(state string) *GithubAuthConfig {
	return &GithubAuthConfig{
		GithubRedirect: "https://github.com/login/oauth/authorize",
		GithubClientId: envutil.GetEnvString("GITHUB_PUBLIC_KEY"),
		GithubCallback: envutil.GetEnvString("GITHUB_AUTH_CALLBACK"),
		ScopeGithub:    "user user:email",
		GithubSecret:   envutil.GetEnvString("GITHUB_SECRET_KEY"),
		State:          state,
	}
}

func (g *GithubAuthConfig) GetAuthURL() string {
	params := url.Values{}
	params.Add("client_id", g.GithubClientId)
	params.Add("redirect_uri", g.GithubCallback)
	params.Add("response_type", g.ResponseType)
	params.Add("scope", g.ScopeGithub)
	params.Add("state", g.State)

	return fmt.Sprintf("%s?%s", g.GithubRedirect, params.Encode())
}

func (g *GithubAuthConfig) HandleOAuthCallback(r *http.Request) (oauth2.OAuthUser, error) {

	if g.State != r.URL.Query().Get("state") {
		return GithubUser{}, errors.New("missing state")
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		return GithubUser{}, errors.New("missing code in callback")
	}

	params := url.Values{}
	params.Add("code", code)
	params.Add("client_id", g.GithubClientId)
	params.Add("client_secret", g.GithubSecret)
	params.Add("redirect_uri", g.GithubCallback)

	resp, err := http.PostForm("https://github.com/login/oauth/access_token", params)
	if err != nil {
		return GithubUser{}, fmt.Errorf("error posting form: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GithubUser{}, fmt.Errorf("error reading token response: %w", err)
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		return GithubUser{}, fmt.Errorf("error parsing token response: %w", err)
	}

	accessToken := values.Get("access_token")
	if accessToken == "" {
		return GithubUser{}, errors.New("access token not found in response")
	}

	client := &http.Client{}

	// Fetch user profile
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	userInfoResp, err := client.Do(req)
	if err != nil {
		return GithubUser{}, fmt.Errorf("error fetching user info: %w", err)
	}
	defer userInfoResp.Body.Close()

	userInfo, err := io.ReadAll(userInfoResp.Body)
	if err != nil {
		return GithubUser{}, fmt.Errorf("error reading user info: %w", err)
	}

	var user GithubUser
	if err := json.Unmarshal(userInfo, &user); err != nil {
		return GithubUser{}, fmt.Errorf("error unmarshaling user info: %w", err)
	}

	// Fetch user email
	emailReq, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	emailReq.Header.Set("Authorization", "Bearer "+accessToken)

	emailResp, err := client.Do(emailReq)
	if err != nil {
		return GithubUser{}, fmt.Errorf("error fetching user emails: %w", err)
	}
	defer emailResp.Body.Close()

	emailBody, err := io.ReadAll(emailResp.Body)
	if err != nil {
		return GithubUser{}, fmt.Errorf("error reading email response: %w", err)
	}

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	if err := json.Unmarshal(emailBody, &emails); err != nil {
		return GithubUser{}, fmt.Errorf("error unmarshaling emails: %w", err)
	}


	// Pick primary verified email
	for _, e := range emails {
		if e.Primary && e.Verified {
			user.Email = e.Email
			break
		}
	}

	return user, nil
}
