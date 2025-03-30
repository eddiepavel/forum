package render

import (
	"fmt"
	"forum-app/middleware"
	"forum-app/models"
	"forum-app/session"
	"html/template"
	"net/http"
	"os"
	"time"
)

var files = []string{
	"./assets/base.html",
	"./assets/partials/nav.html",
	"./assets/partials/create.html",
	"./assets/partials/home.html",
	"./assets/partials/posts.html"}

func (view *View) Render(w http.ResponseWriter, r *http.Request) error {
	tmpl, err := template.ParseFiles(view.Path...)
	if err != nil {
		fmt.Println(os.Getwd())
		return err
	}
	err = tmpl.Execute(w, view.Data)
	if err != nil {
		return err
	}
	return nil
}

func PrepareView(source string, r *http.Request) View {
	user, ok := r.Context().Value(middleware.UserKey).(*models.Users)
	session := r.Context().Value("user_session").(*session.Session)

	posts := []models.Post{
		{
			ID:       1,
			Author:   "Edouardos",
			Category: "General",
			Content:  "This is a test post",
			Title:    "Test Post",
			Time:     time.Now().Add(-time.Hour).Format("2006-01-02 15:04:05"),
		},
		{
			ID:       2,
			Author:   "Edouardos",
			Category: "Announcements",
			Content:  "Welcome to the forum! Feel free to explore and contribute.",
			Title:    "Welcome Post",
			Time:     time.Now().Add(-2 * time.Hour).Format("2006-01-02 15:04:05"),
		},
	}

	data := models.PageData{}
	if ok && user != nil {
		data = models.PageData{User: user, Session: session}
	} else {
		data = models.PageData{User: nil, Session: session}
	}
	data.Data = make(map[string]interface{})
	data.Data["posts"] = posts
	data.Source = source

	view := View{
		Name: source,
		Data: &data,
	}

	if source == "create" || source == "home" {
		view.Path = files
	} else {
		view.Path = []string{"./assets/" + source + ".html"}
	}

	return view
}
