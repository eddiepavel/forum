package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"forum-app/app"
	"forum-app/helpers"
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

			if error != nil || (googleAuthService.GoogleUser{}) == user {
				render.RenderError(w, r, errors.New("call back fail"))
			}

			var UserID int64

			dbUser, err := app.DB.GetUserByEmail(user.Email)

			if err == sql.ErrNoRows {
				id, _ := app.DB.RegisterOauthUser(user.Email, user.Name, "google", user.Picture)

				UserID, _ = id.LastInsertId()
			} else {

				UserID = int64(dbUser.ID)

				if dbUser.Auth != "google" {
					render.RenderError(w, r, errors.New("auth missmatch"))
				}
			}

			session, err := app.DB.SessionInit(int(UserID))

			if err != nil {
				render.RenderError(w, r, err)
				return
			}

			maxAge := helpers.DdSessionTimeSeconds(session.ExpiresAt.Format("2006-01-02 15:04:05"))

			cookie := http.Cookie{
				Name:     "auth-token",
				Value:    session.Token,
				Path:     "/",
				MaxAge:   maxAge,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			}

			http.SetCookie(w, &cookie)

			app.Logger.Info("User logged in", "email", user.Email)

			http.Redirect(w, r, "/home", http.StatusFound)

		}

	}
}
