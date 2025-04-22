package middleware

import (
	"context"
	"forum-app/app"
	"net/http"
	"time"
)

// SessionMiddleware ensures that each request has a valid session.
// It creates a new session if none exists or refreshes the expiration time of an existing session.
func SessionMiddleware(next http.HandlerFunc, app *app.Application) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getSessionCookie, err := r.Cookie("session")
		if err != nil || getSessionCookie.Value == "" {
			// No valid cookie, create a new session
			session := app.Session.CreateSession()
			session.Data = make(map[string]interface{})
			sessionCookie := &http.Cookie{
				Name:     "session",
				Value:    session.ID,
				HttpOnly: true,
				MaxAge:   int(app.Session.SessionDuration.Seconds()),
			}

			http.SetCookie(w, sessionCookie)
			ctx := context.WithValue(r.Context(), "user_session", session)

			next(w, r.WithContext(ctx))
			return
		}

		// Retrieve the session using the cookie value
		session, exists := app.Session.GetSession(getSessionCookie.Value)
		if !exists || session == nil || session.ExpiresAt.Before(time.Now()) {
			// Session is expired or does not exist, create a new session
			if exists {
				app.Session.RemoveSession(getSessionCookie.Value)
			}

			newSession := app.Session.CreateSession()
			newSession.Data = make(map[string]interface{})
			newSessionCookie := &http.Cookie{
				Name:     "session",
				Value:    newSession.ID,
				HttpOnly: true,
				MaxAge:   int(app.Session.SessionDuration.Seconds()),
			}

			http.SetCookie(w, newSessionCookie)
			ctx := context.WithValue(r.Context(), "user_session", newSession)

			next(w, r.WithContext(ctx))
			return
		}

		// Session is valid, refresh its expiration time
		app.Session.RefreshSession(session.ID)

		// Ensure Data map is initialized
		if session.Data == nil {
			session.Data = make(map[string]interface{})
		}

		ctx := context.WithValue(r.Context(), "user_session", session)

		next(w, r.WithContext(ctx))
	})
}
