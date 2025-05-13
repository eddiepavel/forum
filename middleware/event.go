package middleware

import (
	"bytes"
	"fmt"
	"forum-app/app"
	"io"
	"net/http"
)

func EventMiddleware(next http.HandlerFunc, app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("------ Incoming Request ------")
		// Print request body if present
		fmt.Println("\nBody:")
		if r.Body != nil {
			// Read and restore body so it remains usable
			bodyBytes, _ := io.ReadAll(r.Body)
			fmt.Println(string(bodyBytes))
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		fmt.Println("\n-----------------------------")

		next(w, r)
	}
}
