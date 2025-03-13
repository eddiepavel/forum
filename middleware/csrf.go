package middleware

import (
	"forum-app/app"
	"forum-app/helpers"
	"net/http"
)

func CsrfTokenMiddlware(next http.HandlerFunc, app *app.Application) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session")

		if err != nil {
			next(w,r)
			return
		}

		session, exists := app.Session.GetSession(cookie.Value)

		if !exists && r.Method == "POST" {
			http.Error(w, "expiredss", 409)
			return
		}

		if r.Method != "POST" {
			csrfToken, _ := helpers.GenerateToken()
			session.Data["csrf"] = csrfToken
			next(w, r)
			return
		}

	})
}
