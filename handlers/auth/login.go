package auth

import (
	"forum-app/app"
	"forum-app/handlers/forum"
	"forum-app/session"

	"forum-app/helpers"
	"forum-app/helpers/flash"
	"forum-app/helpers/validator"
	"html/template"
	"net/http"
)

func GetLogin(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("user_session").(*session.Session)
	data := forum.PageData{Session: session}

	t, err := template.ParseFiles("./assets/login.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	t.Execute(w, data)

}

func PostLogin(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		inputs := map[string][]interface{}{
			"email":    {"required", "string", "email"},
			"password": {"required", "string"},
		}

		valid, errors := validator.ValidateRequest(r, inputs)

		if !valid {
			flash.HandleMessages(w, r, errors, r.Header.Get("Referer"), "error")
			return
		}

		user, err := app.DB.GetUserByEmail(r.FormValue("email"))

		if err != nil {
			flash.HandleMessages(w, r, map[string]string{"user": "Credentials don't match our records"}, r.Header.Get("Referer"), "error")
			return
		}

		err = helpers.CompareHashAndPassword(user.Password, r.FormValue("password"))

		if err != nil {
			flash.HandleMessages(w, r, map[string]string{"user": "Credentials don't match our records"}, r.Header.Get("Referer"), "error")
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

		http.Redirect(w, r, "/", http.StatusFound)
	}

}
