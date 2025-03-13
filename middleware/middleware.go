package middleware

import (
	"fmt"
	"forum-app/app"
	"net/http"
)

type Middleware func(h http.HandlerFunc, app *app.Application) http.HandlerFunc

func ChainMiddleware(h http.HandlerFunc, k []string, app *app.Application) http.HandlerFunc {

	selectMiddle := map[string]Middleware{
		"auth":    AuthMiddleware,
		"headers": CommonHeaders,
		"logs":    LoggingMiddleware,
		"session": SessionMiddleware,
	}

	globalMiddle := []string{"logs", "headers", "session"}

	wrapped := h

	fullMiddlewareList := append(globalMiddle, k...)

	for i := len(fullMiddlewareList) - 1; i >= 0; i-- {
		key := fullMiddlewareList[i]
		if mw, exists := selectMiddle[key]; exists {
			wrapped = mw(wrapped, app)
		} else {
			fmt.Printf("Middleware %s not found\n", key)
		}
	}

	return wrapped
}
