package middleware

import (
	"fmt"
	"forum-app/app"
	"net/http"
)

type Middleware func(h http.HandlerFunc, app *app.Application) http.HandlerFunc

func ChainMiddleware(h http.HandlerFunc, k []string, app *app.Application) http.HandlerFunc {

	selectMiddle := map[string]Middleware{
		"headers": CommonHeaders,
		"auth":    AuthMiddleware,
		"logs":    LoggingMiddleware,
	}

	if len(k) < 1 {
		return h
	}

	wrapped := h

	for i := len(k) - 1; i >= 0; i-- {
		key := k[i]
		if mw, exists := selectMiddle[key]; exists {
			wrapped = mw(wrapped, app)
		} else {
			fmt.Printf("Middlware %s not found\n", key)
		}
	}

	return wrapped

}
