package middleware

import (
	"forum-app/app"
	"forum-app/models"
	"html/template"
	"net/http"
)

func RateLimitMiddleware(next http.HandlerFunc, app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var key string

		// Determine the key: IP for unauthenticated users, username for authenticated users.
		user, ok := r.Context().Value(UserKey).(*models.Users)
		if ok && user != nil {
			key = user.Username
		} else {
			key = r.RemoteAddr
		}

		// Check rate limit.
		allowed, cooldown := app.RateLimiter.Allow(key)
		if !allowed {
			msg := "Too many requests. Banned for " + cooldown.String()
			tmpl, err := template.ParseFiles("./assets/error.html")
			if err != nil {
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusTooManyRequests)
			tmpl.Execute(w, msg)
			return
		}

		next(w, r)
	}
}
