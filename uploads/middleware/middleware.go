package middleware

import (
	"fmt"
	"forum-app/app"
	"net/http"
)

// Middleware represents a function that wraps an HTTP handler with additional functionality.
type Middleware func(h http.HandlerFunc, app *app.Application) http.HandlerFunc

// ChainMiddleware applies a sequence of middlewares to an HTTP handler.
// It combines global middlewares with route-specific middlewares.
func ChainMiddleware(h http.HandlerFunc, k []string, app *app.Application) http.HandlerFunc {

	selectMiddle := map[string]Middleware{
		"auth":      AuthMiddleware,
		"headers":   CommonHeaders,
		"logs":      LoggingMiddleware,
		"session":   SessionMiddleware,
		"csrf":      CsrfTokenMiddlware,
		"ratelimit": RateLimitMiddleware,
	}

	globalMiddle := []string{"logs", "headers", "csrf", "session", "ratelimit"}

	wrapped := h

	fullMiddlewareList := append(globalMiddle, k...)

	for i := 0; i <= len(fullMiddlewareList)-1; i++ {
		key := fullMiddlewareList[i]
		if mw, exists := selectMiddle[key]; exists {
			wrapped = mw(wrapped, app)
		} else {
			fmt.Printf("Middleware %s not found\n", key)
		}
	}

	return wrapped
}
