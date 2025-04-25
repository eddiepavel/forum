package auth

import (
	"forum-app/app"
	"forum-app/render"

	"forum-app/helpers"
	"forum-app/helpers/validator"
	"net/http"
)

// GetLogin returns an HTTP handler function for rendering the login page.
func GetLogin(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view, err := render.PrepareView("login", r, app)
		if err != nil {
			render.RenderError(w, r, err)
			return
		}

		err = view.Render(w, r)
		if err != nil {
			render.RenderError(w, r, err)
			return
		}
	}
}

// PostLogin handles user login by validating credentials and creating a session.
func PostLogin(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			render.RenderError(w, r, err)
			return
		}

		inputs := map[string][]interface{}{
			"email":    {"required", "string", "email"},
			"password": {"required", "string", "login_attempt"},
		}

		valid, errors := validator.ValidateRequest(r, inputs, app)

		email := r.FormValue("email")

		if !valid {
			cookie, err := r.Cookie("session")
			if err != nil {
				render.RenderError(w, r, err)
				return
			}
			session, _ := app.Session.GetSession(cookie.Value)
			session.SetFlash("error", errors)
			session.SetFlash("old_email", email)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		user, err := app.DB.GetUserByEmail(email)
		if err != nil {
			render.RenderError(w, r, err)
			return
		}

		session, _ := app.DB.SessionInit(user.ID)

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

		redirectURL := r.FormValue("redirect")

		if redirectURL == "" {
			redirectURL = "/"
		}

		http.Redirect(w, r, redirectURL, http.StatusFound)
	}

}
