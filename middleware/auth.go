package middleware

import (
	"context"
	"fmt"
	"forum-app/app"
	"net/http"
	"time"
)

type ContextKey string

const UserKey ContextKey = "user"

// AuthMiddleware authenticates users based on a session token stored in a cookie.
// It adds the authenticated user to the request context if the session is valid.
func AuthMiddleware(next http.HandlerFunc, app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("auth-token")

		if err != nil {
			next(w, r)
			return
		}

		// Retrieve the session
		session, err := app.DB.GetSession("token", cookie.Value)
		if err != nil || session == nil {
			fmt.Println("Session not found")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Use session.ExpiresAt as is (local time)
		expiresAt := session.ExpiresAt

		// Use local time for comparison
		now := time.Now().Add(time.Hour * 3).UTC() // Local time

		if expiresAt.Before(now) {
			app.DB.DeleteSession(session.ID)
			fmt.Println("Session expired")
			expireCookie := http.Cookie{
				Name:   "auth-token",
				Value:  "",
				Path:   "/",
				MaxAge: -1,
			}
			http.SetCookie(w, &expireCookie)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		authUser, _ := app.DB.GetUserById(session.UserId)

		ctx := context.WithValue(r.Context(), UserKey, &authUser)

		next(w, r.WithContext(ctx))
	}
}
