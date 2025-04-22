package middleware

import (
	"forum-app/app"
	"net/http"
)

// LoggingMiddleware logs the details of each incoming HTTP request, such as method and path.
func LoggingMiddleware(next http.HandlerFunc, app *app.Application) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.Logger.Info("Request", "method", r.Method, "path", r.URL.Path)
		next(w, r)
	})
}
