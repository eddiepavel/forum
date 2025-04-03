package render

import (
	"forum-app/middleware"
	"forum-app/models"
	"forum-app/session"
	"html/template"
	"net/http"
)

var files = []string{
	"./assets/base.html",
	"./assets/partials/nav.html",
	"./assets/partials/create.html",
	"./assets/partials/home.html",
	"./assets/partials/posts.html",
	"./assets/partials/view.html",
	"./assets/partials/wip.html"}

var categories = []string{
	"General",
	"Technology",
	"Entertainment",
	"Sports",
	"News",
	"Anouncements",
	"Other"}

func (view *View) Render(w http.ResponseWriter, r *http.Request) error {
	tmpl, err := template.ParseFiles(view.Path...)
	if err != nil {
		return err
	}
	err = tmpl.Execute(w, view.Data)
	if err != nil {
		return err
	}
	return nil
}

func PrepareView(source string, r *http.Request) (View, error) {
	user, ok := r.Context().Value(middleware.UserKey).(*models.Users)
	session := r.Context().Value("user_session").(*session.Session)
	redirect := r.URL.Query().Get("redirect")

	posts, err := models.GetData(source, r)
	if err != nil {
		return View{}, err
	}

	data := models.PageData{}
	if ok && user != nil {
		data = models.PageData{User: user, Session: session}
	} else {
		data = models.PageData{User: nil, Session: session}
	}
	data.Data = make(map[string]interface{})
	if source == "home" {
		data.Data["posts"] = posts
	} else if source == "view" {
		data.Data["post"] = posts
	}
	data.Source = source

	if source == "create" || source == "home" {
		data.Data["categories"] = categories
	}

	view := View{
		Name: source,
		Data: &data,
	}

	if redirect != "" {
		data.Redirect = redirect
	}

	if source == "create" || source == "home" || source == "view" || source == "wip" {
		view.Path = files
	} else {
		view.Path = []string{"./assets/" + source + ".html"}
	}

	return view, nil
}
