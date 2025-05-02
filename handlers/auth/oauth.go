package auth

import (
	"database/sql"
	"errors"
	"forum-app/app"
	"forum-app/helpers"
	"forum-app/middleware"
	"forum-app/render"
	oauth2 "forum-app/services"
	githubAuthService "forum-app/services/github"
	googleAuthService "forum-app/services/google"
	"forum-app/session"
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

		session := r.Context().Value(middleware.SessionKey).(*session.Session)
		state, _ := helpers.GenerateToken()
		session.SetFlash("state", state)
		var service oauth2.OAuthService

		if authType == "google" {
			service = googleAuthService.NewGoogleAuthConfig(state)
		} else {
			service = githubAuthService.NewGithubAuthConfig(state)
		}

		authURL := service.GetAuthURL()
		http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
	}
}

func LoginOAuthCallback(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		stripPathname := strings.Split(r.URL.Path, "/")
		authType := stripPathname[4]

		if !slices.Contains(supportedAuth, authType) && len(stripPathname) > 5 {
			render.RenderError(w, r, errors.New("page not found"))
		}

		session := r.Context().Value(middleware.SessionKey).(*session.Session)
		state, _ := session.GetFlash("state")
		var user oauth2.OAuthUser
		var service oauth2.OAuthService

		if authType == "google" {
			service = googleAuthService.NewGoogleAuthConfig(state.(string))
		} else {
			service = githubAuthService.NewGithubAuthConfig(state.(string))
		}

		user, error := service.HandleOAuthCallback(r)

		if error != nil || (user.GetName()) == "" {
			render.RenderError(w, r, errors.New("call back failed"))
			return
		}

		authenticateWithOAuth(app, w, r, user, authType)

	}
}

func authenticateWithOAuth(app *app.Application, w http.ResponseWriter, r *http.Request, user oauth2.OAuthUser, authType string) {

	var UserID int64

	dbUser, err := app.DB.GetUserByEmail(user.GetEmail())

	if err == sql.ErrNoRows {
		id, _ := app.DB.RegisterOauthUser(user.GetEmail(), user.GetName(), authType, user.GetPicture())

		UserID, _ = id.LastInsertId()
	} else {

		UserID = int64(dbUser.ID)

		if dbUser.Auth != authType {
			render.RenderError(w, r, errors.New("auth missmatch"))
			return
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

	app.Logger.Info("User logged in", "email", user.GetEmail())

	http.Redirect(w, r, "/home", http.StatusFound)
}
