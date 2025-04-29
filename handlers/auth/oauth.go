package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"forum-app/app"
	"forum-app/render"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

var supportedAuth = []string{"google", "github"}
var google string = "https://accounts.google.com/o/oauth2/v2/auth"
var googleClientId string = ""
var googleRedirect string = "http://localhost:8080/login/oauth/callback/google"
var responseType string = "code"
var scopeGoogle string = "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile"
var accessType string = "offline"
var googleSercret string = ""

type ResponseGoogle struct {
	AcccessToken string `json:"access_token"`
}

func LoginOAuth(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		stripPathname := strings.Split(r.URL.Path, "/")
		authType := stripPathname[3]

		if !slices.Contains(supportedAuth, authType) && len(stripPathname) > 4 {
			render.RenderError(w, r, errors.New("page not found"))
		}

		if authType == "google" {
			redirectUrl := fmt.Sprintf(
				"%s?client_id=%s&redirect_uri=%s&response_type=%s&scope=%s&access_type=%s",
				google,
				url.QueryEscape(googleClientId),
				url.QueryEscape(googleRedirect),
				url.QueryEscape(responseType),
				url.QueryEscape(scopeGoogle),
				url.QueryEscape(accessType),
			)
			http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)
		}

	}
}

func LoginOAuthCallback(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		stripPathname := strings.Split(r.URL.Path, "/")
		authType := stripPathname[4]

		fmt.Println(stripPathname)
		fmt.Println(r.URL.RawQuery)

		if !slices.Contains(supportedAuth, authType) && len(stripPathname) > 5 {
			render.RenderError(w, r, errors.New("page not found"))
		}

		params := url.Values{}

		if authType == "google" && r.URL.Query().Has("code") {

			params.Add("code", r.URL.Query().Get("code"))
			params.Add("client_id", googleClientId)
			params.Add("client_secret", googleSercret)
			params.Add("redirect_uri", googleRedirect)
			params.Add("grant_type", "authorization_code")

			resp, err := http.PostForm("https://oauth2.googleapis.com/token", params)
			if err != nil {
				fmt.Println(err)
				return
			}

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)

			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(string(body))
			accessToken := ResponseGoogle{}
			err = json.Unmarshal(body, &accessToken)

			if err != nil {
				fmt.Println("dsadask")
				return
			}

			clientInit := &http.Client{}
			req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
			req.Header.Set("Authorization", "Bearer "+accessToken.AcccessToken)

			resp, _ = clientInit.Do(req)

			defer resp.Body.Close()

			googleUser, _ := io.ReadAll(resp.Body)

			w.Write(googleUser)
		}

	}
}
