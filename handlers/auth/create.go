package auth

import (
	"forum-app/app"
	"forum-app/helpers"
	"html/template"
	"net/http"
)

// GetCreate is a handler function that returns the create forum page.
func GetCreate(app *app.Application) http.HandlerFunc {
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
		t, err := template.ParseFiles("./assets/create.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		t.Execute(w, data)
	}
}

func PostCreate(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
