package middleware

import (
	"forum-app/app"
	"net/http"
)

func SessionMiddleware(next http.HandlerFunc, app *app.Application) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getSessionCookie, err := r.Cookie("session")

		if err != nil {
			session := app.Session.CreateSession()
			sessionCookie := &http.Cookie{
				Name:     "session",
				Value:    session.ID,
				HttpOnly: true,
				MaxAge:   int(app.Session.SessionDuration.Seconds()),
			}

			http.SetCookie(w, sessionCookie)
			next(w, r)
			return
		}

		_, exists := app.Session.GetSession(getSessionCookie.Value)

		if !exists {
			app.Session.RefreshSession(getSessionCookie.Value)
		}

		next(w, r)
	})
}
