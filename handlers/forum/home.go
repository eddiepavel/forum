package forum

import (
	"forum-app/app"
	"forum-app/middleware"
	"forum-app/models"
	"net/http"
	"text/template"
)

func GetHome(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.UserKey).(*models.Users)
		data := PageData{}
		if ok && user != nil {
			data = PageData{User: user}
		} else {
			data = PageData{User: nil}
		}
		t, err := template.ParseFiles("./assets/home.html")

		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		t.Execute(w, data)
	}
}
