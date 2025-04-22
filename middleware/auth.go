package middleware

import (
	"context"
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

		session, err := app.DB.GetSession("token", cookie.Value)

		if err != nil && session == nil {
			app.Logger.Error("Session not found", "error", err)
			expireCookie := http.Cookie{
				Name:   "auth-token",
				Value:  "",
				MaxAge: -1,
			}
			http.SetCookie(w, &expireCookie)
			next(w, r)
			return
		}

		if session.ExpiresAt.Before(time.Now()) {
			app.Logger.Error("Session expired", "error", err)
			app.DB.DeleteSession(session.ID)
			expireCookie := http.Cookie{
				Name:   "auth-token",
				Value:  "",
				MaxAge: -1,
			}
			http.SetCookie(w, &expireCookie)
			next(w, r)
			return
		}

		authUser, _ := app.DB.GetUserById(session.UserId)

		ctx := context.WithValue(r.Context(), UserKey, &authUser)

		next(w, r.WithContext(ctx))
	}
}
