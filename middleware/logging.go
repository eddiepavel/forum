package middleware

import (
	"forum-app/app"
	"net/http"
)

func LoggingMiddleware(next http.HandlerFunc, app *app.Application) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.Logger.Info("Request", "method", r.Method, "path", r.URL.Path)
		next(w, r)
	})
}
