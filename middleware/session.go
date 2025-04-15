package middleware

import (
	"context"
	"forum-app/app"
	"net/http"
)

func SessionMiddleware(next http.HandlerFunc, app *app.Application) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getSessionCookie, err := r.Cookie("session")
		if err != nil {
			session := app.Session.CreateSession()
			session.Data = make(map[string]interface{})
			sessionCookie := &http.Cookie{
				Name:     "session",
				Value:    session.ID,
				HttpOnly: true,
				MaxAge:   int(app.Session.SessionDuration.Seconds()),
			}

			http.SetCookie(w, sessionCookie)
			context := context.WithValue(r.Context(), "user_session", session)

			next(w, r.WithContext(context))
			return
		} else {
			session, exists := app.Session.GetSession(getSessionCookie.Value)
			if !exists {
				session = app.Session.CreateSession()
				session.Data = make(map[string]interface{})
			}
		}

		session, _ := app.Session.GetSession(getSessionCookie.Value)

		// Ensure Data map is initialized
		if session.Data == nil {
			session.Data = make(map[string]interface{})
		}

		context := context.WithValue(r.Context(), "user_session", session)

		next(w, r.WithContext(context))
	})
}
