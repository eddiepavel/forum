package auth

import (
	"fmt"
	"forum-app/app"
	"forum-app/helpers"
	"forum-app/helpers/validator"
	"forum-app/render"
	"net/http"
)

// GetRegister returns an HTTP handler function for rendering the registration page.
func GetRegister(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view, err := render.PrepareView("register", r, app)
		if err != nil {
			render.RenderError(w, r, err)
			return
		}

		err = view.Render(w, r)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	}
}

// PostRegister handles user registration by validating input and creating a new user.
func PostRegister(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			render.RenderError(w, r, err)
			return
		}

		inputs := map[string][]interface{}{
			"email":            {"required", "string", "email", app.DB.CheckUserExists(r.FormValue("email"), r.FormValue("username"))},
			"username":         {"required", "string", app.DB.CheckUserExists(r.FormValue("email"), r.FormValue("username"))},
			"password":         {"required", "string", "password"},
			"confirm_password": {"required", "string", "same:password"},
		}

		valid, errors := validator.ValidateRequest(r, inputs, app)

		userEmail := r.FormValue("email")
		userName := r.FormValue("username")

		if !valid {
			cookie, err := r.Cookie("session")
			if err != nil {
				render.RenderError(w, r, err)
				return
			}
			session, _ := app.Session.GetSession(cookie.Value)
			session.SetFlash("error", errors)
			session.SetFlash("old_email", userEmail)
			session.SetFlash("old_username", userName)
			http.Redirect(w, r, "/register", http.StatusFound)
			return
		}

		userPassword, err := helpers.HashPassword(r.FormValue("password"))

		if err != nil {
			app.Logger.Info("Failed to hash password", "error", err)
			render.RenderError(w, r, err)
			return
		}

		err = app.DB.RegisterUser(userEmail, userName, userPassword)
		if err != nil {
			app.Logger.Info("Failed to hash password", "error", err)
			render.RenderError(w, r, err)
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)
	}

}
