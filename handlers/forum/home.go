package forum

import (
	"forum-app/app"
	"forum-app/middleware"
	"forum-app/models"
	"forum-app/session"
	"net/http"
	"text/template"
)

func GetHome(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.UserKey).(*models.Users)
		session := r.Context().Value("user_session").(*session.Session)

		data := PageData{}
		if ok && user != nil {
			data = PageData{User: user, Session: session}
		} else {
			data = PageData{User: nil, Session: session}
		}

		t, err := template.ParseFiles("./assets/home.html")

		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		t.Execute(w, data)
	}
}
