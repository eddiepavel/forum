package middleware

import (
	"forum-app/app"
	"forum-app/helpers"
	"net/http"
)

// CsrfTokenMiddlware protects against Cross-Site Request Forgery (CSRF) attacks.
// It ensures a valid CSRF token is present in requests.
func CsrfTokenMiddlware(next http.HandlerFunc, app *app.Application) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			next(w, r)
			return
		}

		session, exists := app.Session.GetSession(cookie.Value)
		if !exists || session == nil {
			http.Error(w, "Session not found or expired please refresh", http.StatusUnauthorized)
			return
		}

		if r.Method == "GET" {
			// Ensure CSRF token is available in the session for GET requests.
			if _, ok := session.Data["csrf"]; !ok {
				csrfToken, _ := helpers.GenerateToken()
				session.Data["csrf"] = csrfToken
			}
			next(w, r)
			return
		}

		if r.Method == "POST" {
			// Validate CSRF token for POST requests.
			formCsrfToken := r.FormValue("csrf")
			sessionCsrfToken, ok := session.Data["csrf"].(string)

			if !ok || formCsrfToken != sessionCsrfToken {
				session.SetFlash("csrf_error", "CSRF token mismatch")
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
		}

		next(w, r)
	})
}
