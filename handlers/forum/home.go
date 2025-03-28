package forum

import (
	"forum-app/app"
	"forum-app/middleware"
	"forum-app/models"
	"forum-app/session"
	"net/http"
	"text/template"
	"time"
)

func GetHome(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.UserKey).(*models.Users)
		session := r.Context().Value("user_session").(*session.Session)

		// Initialize the posts slice with one element
		posts := []models.Post{
			{
				ID:       1,
				Author:   "Edouardos",
				Category: "General",
				Content:  "This is a test post",
				Title:    "Test Post",
			},
		}

		// Parse the time for the first post
		var err error
		posts[0].Time, err = time.Parse("2006-01-02 15:04:05", "2021-07-01 12:00:00")
		if err != nil {
			http.Error(w, "Invalid time format", http.StatusInternalServerError)
			return
		}

		data := PageData{}
		if ok && user != nil {
			data = PageData{User: user, Session: session, Posts: posts}
		} else {
			data = PageData{User: nil, Session: session, Posts: posts}
		}

		t, err := template.ParseFiles("./assets/home.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		t.Execute(w, data)
	}
}
