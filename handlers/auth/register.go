package auth

import (
	"forum-app/app"
	"forum-app/helpers"
	"forum-app/helpers/flash"
	"forum-app/helpers/validator"
	"html/template"
	"net/http"
)

func GetRegister(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Generate CSRF token
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "Session not found", http.StatusUnauthorized)
			return
		}

		session, exists := app.Session.GetSession(cookie.Value)
		if !exists {
			http.Error(w, "Session expired", http.StatusUnauthorized)
			return
		}

		csrfToken, _ := helpers.GenerateToken()
		session.Data["csrf"] = csrfToken

		// Pass CSRF token to the template
		data := map[string]interface{}{
			"csrf_token": csrfToken,
		}

		t, err := template.ParseFiles("./assets/register.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		t.Execute(w, data)
	}
}
func StoreRegister(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		inputs := map[string][]interface{}{
			"email":            {"required", "string", "email", app.DB.CheckUserExists(r.FormValue("email"), r.FormValue("username"))},
			"username":         {"required", "string", app.DB.CheckUserExists(r.FormValue("email"), r.FormValue("username"))},
			"password":         {"required", "string"},
			"confirm_password": {"required", "string", "same:password"},
		}

		valid, errors := validator.ValidateRequest(r, inputs)

		if !valid {
			flash.HandleMessages(w, r, errors, r.Header.Get("Referer"), "error")
			return
		}

		userEmail := r.FormValue("email")
		userName := r.FormValue("username")
		userPassword, err := helpers.HashPassword(r.FormValue("password"))

		if err != nil {
			app.Logger.Info("Failed to hash password", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = app.DB.RegisterUser(userEmail, userName, userPassword)
		if err != nil {
			app.Logger.Info("Failed to hash password", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)
	}

}
