package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

// RecoverMiddleware is a middleware to catch panics and return a "500 Internal Server Error".
func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// Recover from panics
			if rec := recover(); rec != nil {
				// Log the recovered panic and stack trace
				log.Printf("Panic: %v\nStack Trace:\n%s", rec, debug.Stack())

				// Respond with an Internal Server Error
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
