package auth

import (
	"errors"
	"fmt"
	"forum-app/app"
	"forum-app/render"
	googleAuthService "forum-app/services"
	"net/http"
	"slices"
	"strings"
)

var supportedAuth = []string{"google", "github"}

func LoginOAuth(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		stripPathname := strings.Split(r.URL.Path, "/")
		authType := stripPathname[3]

		if !slices.Contains(supportedAuth, authType) && len(stripPathname) > 4 {
			render.RenderError(w, r, errors.New("page not found"))
		}

		if authType == "google" {
			authService := googleAuthService.NewGoogleAuthConfig()
			authURL := authService.GetAuthURL()
			fmt.Println(authURL)

			http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
		}

	}
}

func LoginOAuthCallback(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		stripPathname := strings.Split(r.URL.Path, "/")
		authType := stripPathname[4]

		if !slices.Contains(supportedAuth, authType) && len(stripPathname) > 5 {
			render.RenderError(w, r, errors.New("page not found"))
		}

		if authType == "google" {

			authService := googleAuthService.NewGoogleAuthConfig()
			user, error := authService.HandleOAuthCallback(r)

			if error != nil {
				render.RenderError(w, r, errors.New("call back fail"))
			}

			w.Write(user)

		}

	}
}
