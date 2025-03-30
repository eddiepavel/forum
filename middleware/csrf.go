package middleware

import (
	"fmt"
	"forum-app/app"
	"forum-app/helpers"
	"net/http"
)

func CsrfTokenMiddlware(next http.HandlerFunc, app *app.Application) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			next(w, r)
			return
		}

		session, exists := app.Session.GetSession(cookie.Value)
		if !exists || session == nil {
			fmt.Printf("%+v\n", cookie)
			http.Error(w, "Session not found or expired", http.StatusUnauthorized)
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
				fmt.Println("form ", formCsrfToken)
				fmt.Println("session ", sessionCsrfToken)
				http.Error(w, "Invalid CSRF token", http.StatusForbidden)
				return
			}
		}

		next(w, r)
	})
}
